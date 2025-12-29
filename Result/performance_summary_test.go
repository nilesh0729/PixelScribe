package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nilesh0729/PixelScribe/util"
	"github.com/stretchr/testify/require"
)

func createRandomPerformanceSummary(t *testing.T) PerformanceSummary {
	user := RandomUser(t)
	arg := CreatePerformanceSummaryParams{
		UserID:        sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID:   sql.NullInt64{Int64: util.RandomInt(1,9), Valid: true}, // Dictation creation optional for this test unless FK needed
		TotalAttempts: sql.NullInt32{Int32: 1, Valid: true},
		BestAccuracy: sql.NullFloat64{
			Float64: 92.5,
			Valid:   true,
		},
		AverageAccuracy: sql.NullFloat64{
			Float64: 88.0,
			Valid:   true,
		},
		AverageTime: sql.NullFloat64{
			Float64: 12.3,
			Valid:   true,
		},
		LastAttemptAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	ps, err := testQueries.CreatePerformanceSummary(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ps)

	require.Equal(t, arg.UserID, ps.UserID)
	require.Equal(t, arg.DictationID, ps.DictationID)
	require.Equal(t, arg.TotalAttempts, ps.TotalAttempts)
	require.Equal(t, arg.BestAccuracy.Float64, ps.BestAccuracy.Float64)
	require.Equal(t, arg.AverageAccuracy.Float64, ps.AverageAccuracy.Float64)
	require.Equal(t, arg.AverageTime.Float64, ps.AverageTime.Float64)
	require.True(t, ps.LastAttemptAt.Valid)

	return ps
}

func TestCreatePerformanceSummary(t *testing.T) {
	createRandomPerformanceSummary(t)
}

func TestGetPerformanceSummaryByID(t *testing.T) {
	ps1 := createRandomPerformanceSummary(t)

	ps2, err := testQueries.GetPerformanceSummaryByID(context.Background(), ps1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ps2)

	require.Equal(t, ps1.ID, ps2.ID)
	require.Equal(t, ps1.UserID, ps2.UserID)
	require.Equal(t, ps1.DictationID, ps2.DictationID)
}

func TestUpdatePerformanceSummary(t *testing.T) {
	ps := createRandomPerformanceSummary(t)

	arg := UpdatePerformanceSummaryParams{
		TotalAttempts: sql.NullInt32{Int32: 1, Valid: true},
		BestAccuracy: sql.NullFloat64{
			Float64: 95.0,
			Valid:   true,
		},
		AverageAccuracy: sql.NullFloat64{
			Float64: 90.5,
			Valid:   true,
		},
		AverageTime: sql.NullFloat64{
			Float64: 10.8,
			Valid:   true,
		},
		LastAttemptAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: ps.ID,
	}

	updated, err := testQueries.UpdatePerformanceSummary(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, arg.TotalAttempts, updated.TotalAttempts)
	require.Equal(t, arg.BestAccuracy.Float64, updated.BestAccuracy.Float64)
	require.Equal(t, arg.AverageAccuracy.Float64, updated.AverageAccuracy.Float64)
	require.Equal(t, arg.AverageTime.Float64, updated.AverageTime.Float64)
}

func TestListPerformanceSummaryByUser(t *testing.T) {
	// Create a real user first
	user := RandomUser(t)
	userID := sql.NullInt64{Int64: user.ID, Valid: true}

	// Create a performance summary specifically for this user
	arg := CreatePerformanceSummaryParams{
		UserID:        userID,
		DictationID:   sql.NullInt64{Int64: util.RandomInt(1, 9), Valid: true},
		TotalAttempts: sql.NullInt32{Int32: 1, Valid: true},
		BestAccuracy:  sql.NullFloat64{Float64: 90.0, Valid: true},
		AverageAccuracy: sql.NullFloat64{Float64: 90.0, Valid: true},
		AverageTime:     sql.NullFloat64{Float64: 10.0, Valid: true},
		LastAttemptAt:   sql.NullTime{Time: time.Now(), Valid: true},
	}
	_, err := testQueries.CreatePerformanceSummary(context.Background(), arg)
	require.NoError(t, err)

	list, err := testQueries.ListPerformanceSummaryByUser(context.Background(), userID)
	require.NoError(t, err)
	require.NotEmpty(t, list)

	require.GreaterOrEqual(t, len(list), 1)
}

func TestRecentAttemptsByUser(t *testing.T) {
	// Create a real user first
	user := RandomUser(t)
	userID := sql.NullInt64{Int64: user.ID, Valid: true}

	// Create 2 records specifically for this user
	for i := 0; i < 2; i++ {
		arg := CreatePerformanceSummaryParams{
			UserID:        userID,
			DictationID:   sql.NullInt64{Int64: util.RandomInt(1, 9), Valid: true},
			TotalAttempts: sql.NullInt32{Int32: 1, Valid: true},
			BestAccuracy:  sql.NullFloat64{Float64: 90.0, Valid: true},
			AverageAccuracy: sql.NullFloat64{Float64: 90.0, Valid: true},
			AverageTime:     sql.NullFloat64{Float64: 10.0, Valid: true},
			LastAttemptAt:   sql.NullTime{Time: time.Now(), Valid: true},
		}
		_, err := testQueries.CreatePerformanceSummary(context.Background(), arg)
		require.NoError(t, err)
	}

	recent, err := testQueries.RecentAttemptsByUser(context.Background(), RecentAttemptsByUserParams{
		UserID: userID,
		Limit:  2,
	})
	require.NoError(t, err)
	require.Len(t, recent, 2)
}

func TestDeletePerformanceSummary(t *testing.T) {
	ps := createRandomPerformanceSummary(t)

	err := testQueries.DeletePerformanceSummary(context.Background(), ps.ID)
	require.NoError(t, err)

	_, err = testQueries.GetPerformanceSummaryByID(context.Background(), ps.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestGetPerformanceSummaryByUserAndDictation(t *testing.T) {
    ctx := context.Background()

    // First, create a random PerformanceSummary
    ps := createRandomPerformanceSummary(t)

    // Prepare the query argument
    arg := GetPerformanceSummaryByUserAndDictationParams{
        UserID:      ps.UserID,
        DictationID: ps.DictationID,
    }

    // Call the function
    got, err := testQueries.GetPerformanceSummaryByUserAndDictation(ctx, arg)
    require.NoError(t, err)
    require.NotEmpty(t, got)

    // Compare the inserted row with what we got
    require.Equal(t, ps.ID, got.ID)
    require.Equal(t, ps.UserID.Int64, got.UserID.Int64)
    require.Equal(t, ps.DictationID.Int64, got.DictationID.Int64)
    require.Equal(t, ps.TotalAttempts.Int32, got.TotalAttempts.Int32)
    require.Equal(t, ps.BestAccuracy.Float64, got.BestAccuracy.Float64)
    require.Equal(t, ps.AverageAccuracy.Float64, got.AverageAccuracy.Float64)
    require.Equal(t, ps.AverageTime.Float64, got.AverageTime.Float64)

    // For time comparison, check if times are close
    require.WithinDuration(t, ps.LastAttemptAt.Time, got.LastAttemptAt.Time, time.Second)
}


func TestUserAggregatePerformance(t *testing.T) {
	// Create a few performance summaries for same and different users
	for i := 0; i < 3; i++ {
		createRandomPerformanceSummary(t)
	}

	aggregates, err := testQueries.UserAggregatePerformance(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, aggregates)

	for _, agg := range aggregates {
		require.True(t, agg.OverallAvgAccuracy >= 0)
		require.True(t, agg.OverallAvgTime >= 0)
		require.True(t, agg.UserID.Valid)
	}
}
