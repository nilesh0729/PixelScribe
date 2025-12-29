package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
)

type updateSettingsRequest struct {
	UserID                 int64   `json:"user_id" binding:"required"`
	DefaultVoice           string  `json:"default_voice"`
	DefaultSpeed           float64 `json:"default_speed"`
	HighlightColorGrammar  string  `json:"highlight_color_grammar"`
	HighlightColorSpelling string  `json:"highlight_color_spelling"`
	HighlightColorCase     string  `json:"highlight_color_case"`
}

type settingResponse struct {
	ID                     int64   `json:"id"`
	UserID                 int64   `json:"user_id"`
	DefaultVoice           string  `json:"default_voice"`
	DefaultSpeed           float64 `json:"default_speed"`
	HighlightColorGrammar  string  `json:"highlight_color_grammar"`
	HighlightColorSpelling string  `json:"highlight_color_spelling"`
	HighlightColorCase     string  `json:"highlight_color_case"`
}

func newSettingResponse(s db.Setting) settingResponse {
	return settingResponse{
		ID:                     s.ID,
		UserID:                 s.UserID.Int64,
		DefaultVoice:           s.DefaultVoice.String,
		DefaultSpeed:           s.DefaultSpeed.Float64,
		HighlightColorGrammar:  s.HighlightColorGrammar.String,
		HighlightColorSpelling: s.HighlightColorSpelling.String,
		HighlightColorCase:     s.HighlightColorCase.String,
	}
}

func (server *Server) getSettings(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	setting, err := server.store.GetSettingByUserID(ctx, sql.NullInt64{Int64: userID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			// This shouldn't typically happen if user was created via API, but handle it anyway
			ctx.JSON(http.StatusNotFound, gin.H{"error": "settings not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newSettingResponse(setting))
}

func (server *Server) updateSettings(ctx *gin.Context) {
	var req updateSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// First get existing settings to find ID being updated (or we could assume 1:1 user:settings mapping logic)
	// Query GetSettingByUserID is easiest.
	existing, err := server.store.GetSettingByUserID(ctx, sql.NullInt64{Int64: req.UserID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "settings not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateSettingParams{
		ID:                     existing.ID,
		DefaultVoice:           sql.NullString{String: req.DefaultVoice, Valid: req.DefaultVoice != ""},
		DefaultSpeed:           sql.NullFloat64{Float64: req.DefaultSpeed, Valid: req.DefaultSpeed != 0},
		HighlightColorGrammar:  sql.NullString{String: req.HighlightColorGrammar, Valid: req.HighlightColorGrammar != ""},
		HighlightColorSpelling: sql.NullString{String: req.HighlightColorSpelling, Valid: req.HighlightColorSpelling != ""},
		HighlightColorCase:     sql.NullString{String: req.HighlightColorCase, Valid: req.HighlightColorCase != ""},
	}
	
	// If a field is not provided in update request (e.g. empty string), UpdateSetting (as generated) updates it to NULL or value?
	// Generated UpdateSetting sets:
	// default_voice = $2 ... 
	// If we pass NULL, it sets NULL.
	// But Wait! The Generated Query was:
	// Title = COALESCE($1, title) ... (in Dictations)
	// But in Settings:
	// default_voice = $2
	// So it OVERWRITES with NULL if we pass NULL.
	// We need to preserve old values if not provided.
	
	if req.DefaultVoice == "" { arg.DefaultVoice = existing.DefaultVoice }
	if req.DefaultSpeed == 0 { arg.DefaultSpeed = existing.DefaultSpeed }
	if req.HighlightColorGrammar == "" { arg.HighlightColorGrammar = existing.HighlightColorGrammar }
	if req.HighlightColorSpelling == "" { arg.HighlightColorSpelling = existing.HighlightColorSpelling }
	if req.HighlightColorCase == "" { arg.HighlightColorCase = existing.HighlightColorCase }

	updated, err := server.store.UpdateSetting(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newSettingResponse(updated))
}
