package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
	mockdb "github.com/nilesh0729/PixelScribe/Result/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetSettings(t *testing.T) {
	setting := db.Setting{
		ID:                    1,
		UserID:                sql.NullInt64{Int64: 1, Valid: true},
		DefaultVoice:          sql.NullString{String: "en-US", Valid: true},
		DefaultSpeed:          sql.NullFloat64{Float64: 1.0, Valid: true},
		HighlightColorGrammar: sql.NullString{String: "#FFA500", Valid: true},
		CreatedAt:             sql.NullTime{Time: time.Now(), Valid: true},
	}

	testCases := []struct {
		name          string
		userID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSettingByUserID(gomock.Any(), gomock.Eq(sql.NullInt64{Int64: 1, Valid: true})).
					Times(1).
					Return(setting, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			userID: 2,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSettingByUserID(gomock.Any(), gomock.Eq(sql.NullInt64{Int64: 2, Valid: true})).
					Times(1).
					Return(db.Setting{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprintf("/settings?user_id=%d", tc.userID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateSettings(t *testing.T) {
	setting := db.Setting{
		ID:           1,
		UserID:       sql.NullInt64{Int64: 1, Valid: true},
		DefaultVoice: sql.NullString{String: "old-voice", Valid: true},
	}
	
	updatedSetting := db.Setting{
		ID:           1,
		UserID:       sql.NullInt64{Int64: 1, Valid: true},
		DefaultVoice: sql.NullString{String: "new-voice", Valid: true},
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
				"user_id":       1,
				"default_voice": "new-voice",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// 1. Get existing
				store.EXPECT().
					GetSettingByUserID(gomock.Any(), gomock.Eq(sql.NullInt64{Int64: 1, Valid: true})).
					Times(1).
					Return(setting, nil)
				
				// 2. Update (expects merge of old+new)
				arg := db.UpdateSettingParams{
					ID:                     1,
					DefaultVoice:           sql.NullString{String: "new-voice", Valid: true},
					DefaultSpeed:           setting.DefaultSpeed,
					HighlightColorGrammar:  setting.HighlightColorGrammar,
					HighlightColorSpelling: setting.HighlightColorSpelling,
					HighlightColorCase:     setting.HighlightColorCase,
				}
				store.EXPECT().
					UpdateSetting(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedSetting, nil)
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

			request, err := http.NewRequest(http.MethodPut, "/settings", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
