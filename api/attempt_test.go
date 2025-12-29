package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
	mockdb "github.com/nilesh0729/PixelScribe/Result/mock"
	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSubmitAttempt(t *testing.T) {
	// Sample data
	attempt := db.Attempt{
		ID:          1,
		UserID:      sql.NullInt64{Int64: 1, Valid: true},
		DictationID: sql.NullInt64{Int64: 1, Valid: true},
		TypedText:   sql.NullString{String: "Hello world", Valid: true},
		Accuracy:    sql.NullFloat64{Float64: 100.0, Valid: true},
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
	}
	summary := db.PerformanceSummary{
		ID:              1,
		UserID:          sql.NullInt64{Int64: 1, Valid: true},
		DictationID:     sql.NullInt64{Int64: 1, Valid: true},
		AverageAccuracy: sql.NullFloat64{Float64: 100.0, Valid: true},
	}
	
	result := db.SubmitAttemptTxResult{
		Attempt:            attempt,
		PerformanceSummary: summary,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_id":            1,
				"dictation_id":       1,
				"typed_text":         "Hello world",
				"total_words":        2,
				"correct_words":      2,
				"grammatical_errors": 0,
				"spelling_errors":    0,
				"case_errors":        0,
				"accuracy":           100.0,
				"comparison_data":    []interface{}{}, 
				"time_spent":         10.5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAttemptsParams{
					UserID:            sql.NullInt64{Int64: 1, Valid: true},
					DictationID:       sql.NullInt64{Int64: 1, Valid: true},
					TypedText:         sql.NullString{String: "Hello world", Valid: true},
					TotalWords:        sql.NullInt32{Int32: 2, Valid: true},
					CorrectWords:      sql.NullInt32{Int32: 2, Valid: true},
					GrammaticalErrors: sql.NullInt32{Int32: 0, Valid: true},
					SpellingErrors:    sql.NullInt32{Int32: 0, Valid: true},
					CaseErrors:        sql.NullInt32{Int32: 0, Valid: true},
					Accuracy:          sql.NullFloat64{Float64: 100.0, Valid: true},
					ComparisonData:    pqtype.NullRawMessage{RawMessage: json.RawMessage("[]"), Valid: true},
					TimeSpent:         sql.NullFloat64{Float64: 10.5, Valid: true},
				}
				store.EXPECT().
					SubmitAttemptTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/attempts", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
