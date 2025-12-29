package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/require"
	"github.com/nilesh0729/PixelScribe/internal/util"
)

func TestSubmitAttemptTx(t *testing.T) {
	store := NewStore(testDB)

	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	// round 1: First Attempt
	arg1 := CreateAttemptsParams{
		UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
		TypedText:   sql.NullString{String: "hello world", Valid: true},
		TotalWords:  sql.NullInt32{Int32: 2, Valid: true},
		CorrectWords: sql.NullInt32{
			Int32: 2,
			Valid: true,
		},
		GrammaticalErrors: sql.NullInt32{Int32: 0, Valid: true},
		SpellingErrors:    sql.NullInt32{Int32: 0, Valid: true},
		CaseErrors:        sql.NullInt32{Int32: 0, Valid: true},
		Accuracy:          sql.NullFloat64{Float64: 100, Valid: true},
		ComparisonData:    pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
		TimeSpent:         sql.NullFloat64{Float64: 10.0, Valid: true},
	}

	result1, err := store.SubmitAttemptTx(context.Background(), arg1)
	require.NoError(t, err)

	require.NotZero(t, result1.Attempt.ID)
	require.Equal(t, float64(100), result1.PerformanceSummary.AverageAccuracy.Float64)
	require.Equal(t, float64(100), result1.PerformanceSummary.BestAccuracy.Float64)
	require.Equal(t, int32(1), result1.PerformanceSummary.TotalAttempts.Int32)
	require.Equal(t, float64(10.0), result1.PerformanceSummary.AverageTime.Float64)

	// round 2: Second Attempt (Low Accuracy)
	arg2 := CreateAttemptsParams{
		UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
		TypedText:   sql.NullString{String: "hello", Valid: true}, // Missing word
		TotalWords:  sql.NullInt32{Int32: 2, Valid: true},
		CorrectWords: sql.NullInt32{
			Int32: 1,
			Valid: true,
		},
		GrammaticalErrors: sql.NullInt32{Int32: 0, Valid: true},
		SpellingErrors:    sql.NullInt32{Int32: 1, Valid: true},
		CaseErrors:        sql.NullInt32{Int32: 0, Valid: true},
		Accuracy:          sql.NullFloat64{Float64: 50, Valid: true},
		ComparisonData:    pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
		TimeSpent:         sql.NullFloat64{Float64: 20.0, Valid: true},
	}

	result2, err := store.SubmitAttemptTx(context.Background(), arg2)
	require.NoError(t, err)

	require.NotZero(t, result2.Attempt.ID)
	require.NotEqual(t, result1.Attempt.ID, result2.Attempt.ID)
	
	// Check aggregate stats
	// Best Accuracy should still be 100
	require.Equal(t, float64(100), result2.PerformanceSummary.BestAccuracy.Float64)
	
	// Total attempts should be 2
	require.Equal(t, int32(2), result2.PerformanceSummary.TotalAttempts.Int32)
	
	// Average Accuracy: (100 + 50) / 2 = 75
	require.Equal(t, float64(75), result2.PerformanceSummary.AverageAccuracy.Float64)

	// Average Time: (10 + 20) / 2 = 15
	require.Equal(t, float64(15), result2.PerformanceSummary.AverageTime.Float64)
}

func TestCreateUserTx(t *testing.T) {
	store := NewStore(testDB)

	
	arg := CreateUsersParams{
		Name:         sql.NullString{String: util.RandomName(6), Valid: true},
		Username:     util.RandomUsername(),
		Email:        util.RandomEmail(),
		PasswordHash: sql.NullString{String: "hashed_secret", Valid: true},
	}

	result, err := store.CreateUserTx(context.Background(), arg)
	require.NoError(t, err)

	// Verify User
	require.NotZero(t, result.User.ID)
	require.Equal(t, arg.Username, result.User.Username)
	require.Equal(t, arg.Email, result.User.Email)

	// Verify Settings created automatically
	require.NotZero(t, result.Setting.ID)
	require.Equal(t, result.User.ID, result.Setting.UserID.Int64)
	require.Equal(t, "en-US-Neural2-F", result.Setting.DefaultVoice.String)
	require.Equal(t, "#FFA500", result.Setting.HighlightColorGrammar.String)
}

func TestDeleteDictationTx(t *testing.T) {
	store := NewStore(testDB)
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	// Create some attempts and summary first
	attemptArg := CreateAttemptsParams{
		UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
		TypedText:   sql.NullString{String: "test", Valid: true},
		Accuracy:    sql.NullFloat64{Float64: 100, Valid: true},
	}
	_, err := store.SubmitAttemptTx(context.Background(), attemptArg)
	require.NoError(t, err)

	// Verify data exists
	count, err := testQueries.CountAttemptsByDictation(context.Background(), CountAttemptsByDictationParams{
		UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	// Execute Transaction
	err = store.DeleteDictationTx(context.Background(), dict.ID)
	require.NoError(t, err)

	// Verify Data Validations
	// 1. Attempts gone
	count, err = testQueries.CountAttemptsByDictation(context.Background(), CountAttemptsByDictationParams{
		UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
	})
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	// 2. Summary gone
	_, err = testQueries.GetPerformanceSummaryByUserAndDictation(context.Background(), GetPerformanceSummaryByUserAndDictationParams{
		UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
	})
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)

	// 3. Dictation gone - We need to check by ID, but generated query is by Title. 
	// We can use ListAudioDictations or similar, or just trust the previous steps passed and the FK check didn't fail.
	// Actually we can check ListDictationsByUser
	dictations, err := testQueries.ListDictationsByUser(context.Background(), sql.NullInt64{Int64: user.ID, Valid: true})
	require.NoError(t, err)
	for _, d := range dictations {
		require.NotEqual(t, dict.ID, d.ID)
	}
}
