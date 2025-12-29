package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/nilesh0729/PixelScribe/internal/token"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
    userID int64,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(userID, username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
