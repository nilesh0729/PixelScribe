package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "github.com/nilesh0729/PixelScribe/Result"
	mockdb "github.com/nilesh0729/PixelScribe/Result/mock"
	"github.com/nilesh0729/PixelScribe/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListPerformance(t *testing.T) {
	summaries := []db.PerformanceSummary{
		{
			ID:              1,
			UserID:          sql.NullInt64{Int64: 1, Valid: true},
			DictationID:     sql.NullInt64{Int64: 1, Valid: true},
			AverageAccuracy: sql.NullFloat64{Float64: 95.5, Valid: true},
			LastAttemptAt:   sql.NullTime{Time: time.Now(), Valid: true},
		},
	}

	// User for auth
	user, _ := randomUserForLogin(t)

	testCases := []struct {
		name          string
		userID        int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: 1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "bearer", user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPerformanceSummaryByUser(gomock.Any(), gomock.Eq(sql.NullInt64{Int64: 1, Valid: true})).
					Times(1).
					Return(summaries, nil)
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

			url := fmt.Sprintf("/performance?user_id=%d", tc.userID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
