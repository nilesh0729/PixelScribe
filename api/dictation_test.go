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
	"github.com/nilesh0729/PixelScribe/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteDictation(t *testing.T) {
	dictationID := int64(10)

	// User for auth
	user, _ := randomUserForLogin(t)

	testCases := []struct {
		name          string
		dictationID   int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			dictationID: dictationID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "bearer", user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDictationTx(gomock.Any(), gomock.Eq(dictationID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "InternalError",
			dictationID: dictationID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "bearer", user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDictationTx(gomock.Any(), gomock.Eq(dictationID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "InvalidID",
			dictationID: 0, 
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "bearer", user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No calls expected
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/dictations/%d", tc.dictationID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateDictation(t *testing.T) {
	// Simple test for text dictation creation
	d := db.Dictation{
		ID:        1,
		UserID:    sql.NullInt64{Int64: 1, Valid: true},
		Title:     sql.NullString{String: "Test Dictation", Valid: true},
		Type:      sql.NullString{String: "text", Valid: true},
		Content:   sql.NullString{String: "Content", Valid: true},
		Language:  sql.NullString{String: "en-US", Valid: true},
		CreatedAt: time.Now(),
	}

	// User for auth
	user, _ := randomUserForLogin(t)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK_Text",
			body: gin.H{
				"user_id":  1,
				"title":    "Test Dictation",
				"type":     "text",
				"content":  "Content",
				"language": "en-US",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "bearer", user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTextDictationsParams{
					UserID:   sql.NullInt64{Int64: 0, Valid: true},
					Title:    sql.NullString{String: "Test Dictation", Valid: true},
					Content:  sql.NullString{String: "Content", Valid: true},
					Language: sql.NullString{String: "en-US", Valid: true},
				}
				// Mock GetUsers call to fetch user ID
				store.EXPECT().
					GetUsers(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				// Mock CreateTextDictations call
				store.EXPECT().
					CreateTextDictations(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(d, nil)
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

			request, err := http.NewRequest(http.MethodPost, "/dictations", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
