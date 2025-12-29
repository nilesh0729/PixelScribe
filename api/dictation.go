package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
	"github.com/nilesh0729/PixelScribe/token"
)

type createDictationRequest struct {
	UserID   int64  `json:"user_id"`
	Title    string `json:"title" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=text audio"`
	Content  string `json:"content"`   // Required for text
	AudioURL string `json:"audio_url"` // Required for audio
	Language string `json:"language" binding:"required"`
}
// ... func newDictationResponse ...
// ... func createDictation ... (will be replaced in next call or same call if contiguous)

// Actually I can only replace one contiguous block. 
// The structs are far apart. I will do them separately or use multi_replace.
// I'll use multi_replace to be efficient.

type dictationResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Content   string    `json:"content,omitempty"`
	AudioURL  string    `json:"audio_url,omitempty"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
}

func newDictationResponse(d db.Dictation) dictationResponse {
	return dictationResponse{
		ID:        d.ID,
		UserID:    d.UserID.Int64,
		Title:     d.Title.String,
		Type:      d.Type.String,
		Content:   d.Content.String,
		AudioURL:  d.AudioUrl.String,
		Language:  d.Language.String,
		CreatedAt: d.CreatedAt,
	}
}

func (server *Server) createDictation(ctx *gin.Context) {
	var req createDictationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUsers(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var dictation db.Dictation

	if req.Type == "text" {
		if req.Content == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "content is required for text dictation"})
			return
		}
		arg := db.CreateTextDictationsParams{
			UserID:   sql.NullInt64{Int64: user.ID, Valid: true},
			Title:    sql.NullString{String: req.Title, Valid: true},
			Content:  sql.NullString{String: req.Content, Valid: true},
			Language: sql.NullString{String: req.Language, Valid: true},
		}
		dictation, err = server.store.CreateTextDictations(ctx, arg)
	} else if req.Type == "audio" {
		if req.AudioURL == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "audio_url is required for audio dictation"})
			return
		}
		arg := db.CreateAudioDictationsParams{
			UserID:   sql.NullInt64{Int64: user.ID, Valid: true},
			Title:    sql.NullString{String: req.Title, Valid: true},
			AudioUrl: sql.NullString{String: req.AudioURL, Valid: true},
			Language: sql.NullString{String: req.Language, Valid: true},
		}
		dictation, err = server.store.CreateAudioDictations(ctx, arg)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newDictationResponse(dictation))
}

type listDictationsRequest struct {
	UserID int64  `form:"user_id"`
	Type   string `form:"type" binding:"omitempty,oneof=text audio"`
}

func (server *Server) listDictations(ctx *gin.Context) {
	var req listDictationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUsers(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var dictations []db.Dictation

	// Note: Generated queries are split by type: ListTextDictations, ListAudioDictations
	// Or ListDictationsByUser.
	// Let's use ListDictationsByUser if no type specified, or filter manually?
	// Looking at querier.go:
	// ListTextDictations(ctx, userID)
	// ListAudioDictations(ctx, userID)
	// ListDictationsByUser(ctx, userID)

	if req.Type == "text" {
		dictations, err = server.store.ListTextDictations(ctx, sql.NullInt64{Int64: user.ID, Valid: true})
	} else if req.Type == "audio" {
		dictations, err = server.store.ListAudioDictations(ctx, sql.NullInt64{Int64: user.ID, Valid: true})
	} else {
		dictations, err = server.store.ListDictationsByUser(ctx, sql.NullInt64{Int64: user.ID, Valid: true})
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := make([]dictationResponse, len(dictations))
	for i, d := range dictations {
		rsp[i] = newDictationResponse(d)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type deleteDictationRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteDictation(ctx *gin.Context) {
	var req deleteDictationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Use Transaction for deleting dictation (cascading)
	err := server.store.DeleteDictationTx(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
