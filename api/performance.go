package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
)

func (server *Server) listPerformance(ctx *gin.Context) {
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

	summaries, err := server.store.ListPerformanceSummaryByUser(ctx, sql.NullInt64{Int64: userID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, summaries)
}

func (server *Server) getOverallPerformance(ctx *gin.Context) {
	// This uses UserAggregatePerformance query which returns all users.
	// That might not be efficient for a single user query but schema doesn't have "GetUserAggregate".
	// We can implement strict filtering in Go for now or just return summary for everyone?
	// Actually query is `AVG... GROUP BY user_id`.
	// For now, let's just return ListPerformance as the main dashboard data.
	// Or we use `RecentAttemptsByUser` to get a feed.
	
	// Let's implement RecentAttemptsByUser as "feed"
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

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	recent, err := server.store.RecentAttemptsByUser(ctx, db.RecentAttemptsByUserParams{
		UserID: sql.NullInt64{Int64: userID, Valid: true},
		Limit:  int32(limit),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, recent)
}
