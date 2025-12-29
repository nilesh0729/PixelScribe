package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/require"
)

func createRandomAttempt(t *testing.T, userID int64, dictationID int64) Attempt {
	arg := CreateAttemptsParams{
		UserID:      sql.NullInt64{Int64: userID, Valid: true},
		DictationID: sql.NullInt64{Int64: dictationID, Valid: true},
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
		TimeSpent:         sql.NullFloat64{Float64: 1.5, Valid: true},
	}

	attempt, err := testQueries.CreateAttempts(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, attempt.ID)
	require.Equal(t, userID, attempt.UserID.Int64)
	require.Equal(t, dictationID, attempt.DictationID.Int64)

	return attempt
}

func TestCreateAttempts(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	attempt := createRandomAttempt(t, user.ID, dict.ID)

	require.NotEmpty(t, attempt)
	require.Equal(t, user.ID, attempt.UserID.Int64)
	require.Equal(t, dict.ID, attempt.DictationID.Int64)
	require.True(t, attempt.TypedText.Valid)
	require.NotZero(t, attempt.AttemptNo)
}

func TestCountAttemptsByDictation(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	createRandomAttempt(t, user.ID, dict.ID)
	createRandomAttempt(t, user.ID, dict.ID)

	arg := CountAttemptsByDictationParams{
		UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
	}
	count, err := testQueries.CountAttemptsByDictation(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, int64(2), count)
}

func TestGetAttemptById(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)
	attempt := createRandomAttempt(t, user.ID, dict.ID)

	result, err := testQueries.GetAttemptById(context.Background(), attempt.ID)

	require.NoError(t, err)
	require.Equal(t, attempt.ID, result.ID)
}

func TestGetLatestAttempt(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	createRandomAttempt(t, user.ID, dict.ID)
	last := createRandomAttempt(t, user.ID, dict.ID)

	result, err := testQueries.GetLatestAttempt(context.Background(),
		GetLatestAttemptParams{
			UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
			DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
		},
	)

	require.NoError(t, err)
	require.Equal(t, last.ID, result.ID)
}

func TestListAttemptsByDictation(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	createRandomAttempt(t, user.ID, dict.ID)
	createRandomAttempt(t, user.ID, dict.ID)

	items, err := testQueries.ListAttemptsByDictation(context.Background(),
		sql.NullInt64{Int64: dict.ID, Valid: true},
	)

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(items), 2)
}

func TestListAttemptsByUser(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	createRandomAttempt(t, user.ID, dict.ID)
	createRandomAttempt(t, user.ID, dict.ID)

	items, err := testQueries.ListAttemptsByUser(context.Background(),
		sql.NullInt64{Int64: user.ID, Valid: true},
	)

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(items), 2)
}

func TestUpdateAttemptAccuracy(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)
	attempt := createRandomAttempt(t, user.ID, dict.ID)

	arg := UpdateAttemptAccuracyParams{
		ID:                attempt.ID,
		Accuracy:          sql.NullFloat64{Float64: 50, Valid: true},
		CorrectWords:      sql.NullInt32{Int32: 1, Valid: true},
		GrammaticalErrors: sql.NullInt32{Int32: 1, Valid: true},
		SpellingErrors:    sql.NullInt32{Int32: 1, Valid: true},
		CaseErrors:        sql.NullInt32{Int32: 1, Valid: true},
		ComparisonData:    pqtype.NullRawMessage{RawMessage: []byte(`{"x":1}`), Valid: true},
		TimeSpent:         sql.NullFloat64{Float64: 3.0, Valid: true},
	}

	updated, err := testQueries.UpdateAttemptAccuracy(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, updated.ID, attempt.ID)
	require.Equal(t, float64(50), updated.Accuracy.Float64)
}

func TestDeleteAttempt(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)
	attempt := createRandomAttempt(t, user.ID, dict.ID)

	err := testQueries.DeleteAttempt(context.Background(), attempt.ID)
	require.NoError(t, err)

	// Now check it is gone
	res, err := testQueries.GetAttemptById(context.Background(), attempt.ID)

	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows.Error(), err.Error())
	require.Empty(t, res)
}

func TestDeleteAttemptsByDictation(t *testing.T) {
	user := RandomUser(t)
	dict := RandomTextDictation(t, user)

	createRandomAttempt(t, user.ID, dict.ID)
	createRandomAttempt(t, user.ID, dict.ID)

	err := testQueries.DeleteAttemptsByDictation(context.Background(),
		DeleteAttemptsByDictationParams{
			DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
			UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
		},
	)
	require.NoError(t, err)

	count, err := testQueries.CountAttemptsByDictation(context.Background(),
		CountAttemptsByDictationParams{
			UserID:      sql.NullInt64{Int64: user.ID, Valid: true},
			DictationID: sql.NullInt64{Int64: dict.ID, Valid: true},
		},
	)

	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}
