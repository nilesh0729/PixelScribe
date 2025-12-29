package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// helper function to create a random setting for testing
func createRandomSetting(t *testing.T, userID int64) Setting {
	arg := CreateSettingParams{
		UserID:               sql.NullInt64{Int64: userID, Valid: true},
		DefaultVoice:         sql.NullString{String: "en-US", Valid: true},
		DefaultSpeed:         sql.NullFloat64{Float64: 1.0, Valid: true},
		HighlightColorGrammar: sql.NullString{String: "#FF0000", Valid: true},
		HighlightColorSpelling: sql.NullString{String: "#00FF00", Valid: true},
		HighlightColorCase:    sql.NullString{String: "#0000FF", Valid: true},
	}

	setting, err := testQueries.CreateSetting(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, setting)

	require.Equal(t, arg.UserID, setting.UserID)
	require.Equal(t, arg.DefaultVoice, setting.DefaultVoice)
	require.Equal(t, arg.DefaultSpeed, setting.DefaultSpeed)
	require.Equal(t, arg.HighlightColorGrammar, setting.HighlightColorGrammar)
	require.Equal(t, arg.HighlightColorSpelling, setting.HighlightColorSpelling)
	require.Equal(t, arg.HighlightColorCase, setting.HighlightColorCase)

	require.NotZero(t, setting.ID)
	require.WithinDuration(t, time.Now(), setting.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), setting.UpdatedAt.Time, time.Second)

	return setting
}

func TestCreateSetting(t *testing.T) {
	user := RandomUser(t)
	createRandomSetting(t, user.ID)
}

func TestGetSettingByID(t *testing.T) {
	user := RandomUser(t)
	setting1 := createRandomSetting(t, user.ID)
	setting2, err := testQueries.GetSettingByID(context.Background(), setting1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, setting2)

	require.Equal(t, setting1.ID, setting2.ID)
	require.Equal(t, setting1.UserID, setting2.UserID)
}

func TestGetSettingByUserID(t *testing.T) {
	user := RandomUser(t)
	userID := sql.NullInt64{Int64: user.ID, Valid: true}
	setting1 := createRandomSetting(t, userID.Int64)
	setting2, err := testQueries.GetSettingByUserID(context.Background(), userID)
	require.NoError(t, err)
	require.NotEmpty(t, setting2)

	require.Equal(t, setting1.ID, setting2.ID)
	require.Equal(t, setting1.UserID, setting2.UserID)
}

func TestUpdateSetting(t *testing.T) {
	user := RandomUser(t)
	setting1 := createRandomSetting(t, user.ID)

	arg := UpdateSettingParams{
		ID:                    setting1.ID,
		DefaultVoice:          sql.NullString{String: "en-GB", Valid: true},
		DefaultSpeed:          sql.NullFloat64{Float64: 1.5, Valid: true},
		HighlightColorGrammar: sql.NullString{String: "#111111", Valid: true},
		HighlightColorSpelling: sql.NullString{String: "#222222", Valid: true},
		HighlightColorCase:     sql.NullString{String: "#333333", Valid: true},
	}

	setting2, err := testQueries.UpdateSetting(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, setting2)

	require.Equal(t, arg.DefaultVoice, setting2.DefaultVoice)
	require.Equal(t, arg.DefaultSpeed, setting2.DefaultSpeed)
	require.Equal(t, arg.HighlightColorGrammar, setting2.HighlightColorGrammar)
	require.Equal(t, arg.HighlightColorSpelling, setting2.HighlightColorSpelling)
	require.Equal(t, arg.HighlightColorCase, setting2.HighlightColorCase)
}

func TestDeleteSetting(t *testing.T) {
	user := RandomUser(t)
	setting := createRandomSetting(t, user.ID)
	err := testQueries.DeleteSetting(context.Background(), setting.ID)
	require.NoError(t, err)

	setting2, err := testQueries.GetSettingByID(context.Background(), setting.ID)
	require.Error(t, err)
	require.Empty(t, setting2)
}
