package api

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	server := NewServer(store)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// requireBodyMatchUser checks if the response body matches the user
func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser userResponse
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Name.String, gotUser.Name)
	require.Equal(t, user.Email, gotUser.Email)
}
