package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sqlc-dev/pqtype"
	db "github.com/nilesh0729/PixelScribe/Result"
)

type submitAttemptRequest struct {
	UserID            int64           `json:"user_id" binding:"required"`
	DictationID       int64           `json:"dictation_id" binding:"required"`
	TypedText         string          `json:"typed_text"`
	AttemptNo         int32           `json:"attempt_no"`
	TotalWords        int32           `json:"total_words"`
	CorrectWords      int32           `json:"correct_words"`
	GrammaticalErrors int32           `json:"grammatical_errors"`
	SpellingErrors    int32           `json:"spelling_errors"`
	CaseErrors        int32           `json:"case_errors"`
	Accuracy          float64         `json:"accuracy"`
	ComparisonData    json.RawMessage `json:"comparison_data"`
	TimeSpent         float64         `json:"time_spent"`
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

	arg := db.CreateAttemptsParams{
		UserID:            sql.NullInt64{Int64: req.UserID, Valid: true},
		DictationID:       sql.NullInt64{Int64: req.DictationID, Valid: true},
		TypedText:         sql.NullString{String: req.TypedText, Valid: true},
		TotalWords:        sql.NullInt32{Int32: req.TotalWords, Valid: true},
		CorrectWords:      sql.NullInt32{Int32: req.CorrectWords, Valid: true},
		GrammaticalErrors: sql.NullInt32{Int32: req.GrammaticalErrors, Valid: true},
		SpellingErrors:    sql.NullInt32{Int32: req.SpellingErrors, Valid: true},
		CaseErrors:        sql.NullInt32{Int32: req.CaseErrors, Valid: true},
		Accuracy:          sql.NullFloat64{Float64: req.Accuracy, Valid: true},
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

	if req.DictationID != 0 {
		attempts, err = server.store.ListAttemptsByDictation(ctx, sql.NullInt64{Int64: req.DictationID, Valid: true})
	} else if req.UserID != 0 {
		attempts, err = server.store.ListAttemptsByUser(ctx, sql.NullInt64{Int64: req.UserID, Valid: true})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id or dictation_id required"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Simplifying response for list (lightweight)
	ctx.JSON(http.StatusOK, attempts)
}
