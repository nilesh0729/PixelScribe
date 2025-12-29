package api

import (
	"context"
	"database/sql"
	"encoding/json"
    "fmt"
	"net/http"
    "strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sqlc-dev/pqtype"
	db "github.com/nilesh0729/PixelScribe/internal/db/sqlc"
    "github.com/nilesh0729/PixelScribe/internal/token"
)

type submitAttemptRequest struct {
	DictationID       int64           `json:"dictation_id" binding:"required"`
	TypedText         string          `json:"typed_text"`
	TimeSpent         float64         `json:"time_spent"`
    // Optional / Calculated server-side fields
	TotalWords        int32           `json:"total_words"`
	CorrectWords      int32           `json:"correct_words"`
	GrammaticalErrors int32           `json:"grammatical_errors"`
	SpellingErrors    int32           `json:"spelling_errors"`
	CaseErrors        int32           `json:"case_errors"`
	Accuracy          float64         `json:"accuracy"`
	ComparisonData    json.RawMessage `json:"comparison_data"`
}

type attemptResponse struct {
	ID                int64           `json:"id"`
	UserID            int64           `json:"user_id"`
	DictationID       int64           `json:"dictation_id"`
	TypedText         string          `json:"typed_text"`
	AttemptNo         int32           `json:"attempt_no"`
	Accuracy          float64         `json:"accuracy"`
	TimeSpent         float64         `json:"time_spent"`
	CreatedAt         time.Time       `json:"created_at"`
	PerformanceUpdate *performanceSum `json:"performance_update,omitempty"`
}

type performanceSum struct {
	TotalAttempts   int32   `json:"total_attempts"`
	BestAccuracy    float64 `json:"best_accuracy"`
	AverageAccuracy float64 `json:"average_accuracy"`
	AverageTime     float64 `json:"average_time"`
}

func (server *Server) submitAttempt(ctx *gin.Context) {
	var req submitAttemptRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Fetch original dictation for verification
    dictation, err := server.store.GetDictation(ctx, req.DictationID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("dictation not found")))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    // Server-side calculation
    originalText := dictation.Content.String
    typedText := req.TypedText

    originalWords := strings.Fields(originalText)
    typedWords := strings.Fields(typedText)
    
    totalWords := int32(len(originalWords))
    correctWords := int32(0)
    
    // Simple verification algorithm: Word-by-word match
    // Note: This is a basic comparison. For more advanced diffing (insertions/deletions),
    // we would need a diff library (e.g. sergi/go-diff), but this suffices for exact matching verification.
    limit := len(originalWords)
    if len(typedWords) < limit {
        limit = len(typedWords)
    }

    for i := 0; i < limit; i++ {
        // Case-insensitive comparison could be used, but strict typing usually requires exact case
        if originalWords[i] == typedWords[i] {
            correctWords++
        }
    }
    
    // Calculate accuracy
    accuracy := float64(0)
    if totalWords > 0 {
        accuracy = (float64(correctWords) / float64(totalWords)) * 100
    }

    // Simplified error category estimation (can be improved with proper diffing later)
    errors := totalWords - correctWords
    
    // Get UserID from auth payload
    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAttemptsParams{
		UserID:            sql.NullInt64{Int64: authPayload.UserID, Valid: true},
		DictationID:       sql.NullInt64{Int64: req.DictationID, Valid: true},
		TypedText:         sql.NullString{String: req.TypedText, Valid: true},
		TotalWords:        sql.NullInt32{Int32: totalWords, Valid: true},
		CorrectWords:      sql.NullInt32{Int32: correctWords, Valid: true},
		GrammaticalErrors: sql.NullInt32{Int32: 0, Valid: true}, // Placeholder
		SpellingErrors:    sql.NullInt32{Int32: errors, Valid: true}, // Lump all errors here for now
		CaseErrors:        sql.NullInt32{Int32: 0, Valid: true}, // Placeholder
		Accuracy:          sql.NullFloat64{Float64: accuracy, Valid: true},
		ComparisonData:    pqtype.NullRawMessage{RawMessage: req.ComparisonData, Valid: len(req.ComparisonData) > 0},
		TimeSpent:         sql.NullFloat64{Float64: req.TimeSpent, Valid: true},
	}

	// Use Transaction
	result, err := server.store.SubmitAttemptTx(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := attemptResponse{
		ID:          result.Attempt.ID,
		UserID:      result.Attempt.UserID.Int64,
		DictationID: result.Attempt.DictationID.Int64,
		TypedText:   result.Attempt.TypedText.String,
		AttemptNo:   result.Attempt.AttemptNo.Int32,
		Accuracy:    result.Attempt.Accuracy.Float64,
		TimeSpent:   result.Attempt.TimeSpent.Float64,
		CreatedAt:   result.Attempt.CreatedAt.Time,
		PerformanceUpdate: &performanceSum{
			TotalAttempts:   result.PerformanceSummary.TotalAttempts.Int32,
			BestAccuracy:    result.PerformanceSummary.BestAccuracy.Float64,
			AverageAccuracy: result.PerformanceSummary.AverageAccuracy.Float64,
			AverageTime:     result.PerformanceSummary.AverageTime.Float64,
		},
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listAttemptsRequest struct {
	UserID      int64 `form:"user_id"`
	DictationID int64 `form:"dictation_id"`
}

func (server *Server) listAttempts(ctx *gin.Context) {
	var req listAttemptsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var attempts []db.Attempt
	var err error

	// Get authenticated user
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.DictationID != 0 {
		attempts, err = server.store.ListAttemptsByDictation(ctx, sql.NullInt64{Int64: req.DictationID, Valid: true})
	} else {
		// Default to authenticated user if no specific user requested (or if specific user requested, we could enforce permissions)
		// For now, let's just use the requested user OR the authenticated user
		targetUserID := req.UserID
		if targetUserID == 0 {
			targetUserID = authPayload.UserID
		}
		attempts, err = server.store.ListAttemptsByUser(ctx, sql.NullInt64{Int64: targetUserID, Valid: true})
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Simplifying response for list (lightweight)
	var rsp []attemptResponse
	for _, attempt := range attempts {
		rsp = append(rsp, attemptResponse{
			ID:          attempt.ID,
			UserID:      attempt.UserID.Int64,
			DictationID: attempt.DictationID.Int64,
			TypedText:   attempt.TypedText.String,
			AttemptNo:   attempt.AttemptNo.Int32,
			Accuracy:    attempt.Accuracy.Float64,
			TimeSpent:   attempt.TimeSpent.Float64,
			CreatedAt:   attempt.CreatedAt.Time,
		})
	}
	ctx.JSON(http.StatusOK, rsp)
}

type getAttemptRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAttempt(ctx *gin.Context) {
	var req getAttemptRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	attempt, err := server.store.GetAttemptById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

    // Security check: Ensure the attempt belongs to the user
    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    if attempt.UserID.Int64 != authPayload.UserID {
        err := fmt.Errorf("account doesn't belong to the authenticated user")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    // Construct response
	rsp := attemptResponse{
		ID:          attempt.ID,
		UserID:      attempt.UserID.Int64,
		DictationID: attempt.DictationID.Int64,
		TypedText:   attempt.TypedText.String,
		AttemptNo:   attempt.AttemptNo.Int32,
		Accuracy:    attempt.Accuracy.Float64,
		TimeSpent:   attempt.TimeSpent.Float64,
		CreatedAt:   attempt.CreatedAt.Time,
	}

	ctx.JSON(http.StatusOK, rsp)
}
