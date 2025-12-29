package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type generateTTSRequest struct {
	Text string `json:"text" binding:"required"`
}

func (server *Server) generateTTS(ctx *gin.Context) {
	var req generateTTSRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Prepare OpenAI Request
	openAIUrl := "https://api.openai.com/v1/audio/speech"
	payload := map[string]interface{}{
		"model":           "tts-1",
		"input":           req.Text,
		"voice":           "alloy",
		"response_format": "mp3",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	request, err := http.NewRequest("POST", openAIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+server.config.OpenAIKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error":       "OpenAI API failed",
			"details":     string(bodyBytes),
			"status_code": response.StatusCode,
		})
		return
	}

	// Stream response back to client
	ctx.Header("Content-Type", "audio/mpeg")
	ctx.Header("Transfer-Encoding", "chunked")

	// Copy the stream
	_, err = io.Copy(ctx.Writer, response.Body)
	if err != nil {
		// Cannot write JSON error here as headers likely sent
		return
	}
}
