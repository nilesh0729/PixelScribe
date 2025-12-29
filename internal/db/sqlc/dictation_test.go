package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nilesh0729/PixelScribe/internal/util"
	"github.com/stretchr/testify/require"
)

func RandomAudioDictation(t *testing.T, User User)Dictation {
	
	arg := CreateAudioDictationsParams{
		UserID: sql.NullInt64{Int64: User.ID, Valid: true},
		Title: sql.NullString{String: util.RandomString(4), Valid: true},
		AudioUrl: sql.NullString{String: util.RandomString(6), Valid: true},
		Language: sql.NullString{String: util.RandomString(5), Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	dictation, err := testQueries.CreateAudioDictations(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, dictation)

	require.Equal(t, dictation.UserID, arg.UserID)
	require.Equal(t, dictation.Title, arg.Title)
	require.Equal(t, dictation.AudioUrl, arg.AudioUrl)
	require.Equal(t, dictation.Language, arg.Language)
	
	require.WithinDuration(t, dictation.CreatedAt, User.CreatedAt, time.Second)
	require.WithinDuration(t,dictation.UpdatedAt, User.UpdatedAt, time.Second)
	
	return dictation
}

func RandomTextDictation(t *testing.T, User User)Dictation{

	arg := CreateTextDictationsParams{
		UserID: sql.NullInt64{Int64: User.ID, Valid: true},
		Title: sql.NullString{String: util.RandomString(5), Valid: true},
		Content: sql.NullString{String: util.RandomString(8),Valid: true},
		Language: sql.NullString{String: util.RandomString(5), Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	dictation, err := testQueries.CreateTextDictations(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, dictation)
	require.Equal(t, dictation.UserID.Int64, User.ID)
	require.Equal(t, dictation.Title, arg.Title)
	require.Equal(t, dictation.Content, arg.Content)
	require.Equal(t, "text", dictation.Type.String)
	require.Equal(t, dictation.Language, arg.Language)
	
	require.WithinDuration(t, dictation.CreatedAt, arg.CreatedAt, time.Second)
	require.WithinDuration(t, dictation.UpdatedAt, arg.UpdatedAt, time.Second)

	return dictation
}

func TestCreateAudioDictation(t *testing.T){
	user := RandomUser(t)
	RandomAudioDictation(t, user)
}

func TestCreateTextDictation(t *testing.T){
	user := RandomUser(t)
	RandomTextDictation(t, user)
}

func TestGetDictationByTitle(t *testing.T) {
	user := RandomUser(t)
	d1 := RandomTextDictation(t, user)

	d2, err := testQueries.GetDictationsByTitle(context.Background(), d1.Title)
	require.NoError(t, err)
	require.NotEmpty(t, d2)

	require.Equal(t, d1.ID, d2.ID)
	require.Equal(t, d1.Title, d2.Title)
	require.Equal(t, d1.Type, d2.Type)
}

func TestUpdateDictation(t *testing.T) {
	user := RandomUser(t)
	d1 := RandomTextDictation(t, user)

	arg := UpdateDictationParams{
		Title:    sql.NullString{String: util.RandomString(15), Valid: true},
		Content:  sql.NullString{String: util.RandomString(40), Valid: true},
		AudioUrl: sql.NullString{Valid: false}, 
		ID:       d1.ID,
		UserID:   d1.UserID,
	}

	d2, err := testQueries.UpdateDictation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, d2)

	require.Equal(t, d1.ID, d2.ID)
	require.Equal(t, arg.Title, d2.Title)
	require.Equal(t, arg.Content, d2.Content)
	require.WithinDuration(t, d2.UpdatedAt, time.Now(), time.Second)
	require.WithinDuration(t, d1.CreatedAt, d2.CreatedAt, time.Second)
}

func TestDeleteDictation(t *testing.T) {
	user := RandomUser(t)
	d := RandomTextDictation(t, user)

	err := testQueries.DeleteDictations(context.Background(), d.Title)
	require.NoError(t, err)

	d2, err := testQueries.GetDictationsByTitle(context.Background(), d.Title)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, d2)
}


func TestListAudioDictations(t *testing.T) {
	user := RandomUser(t)

	for i := 0; i < 3; i++ {
		RandomAudioDictation(t, user)
	}
	for i := 0; i < 2; i++ {
		RandomTextDictation(t, user)
	}

	list, err := testQueries.ListAudioDictations(context.Background(), sql.NullInt64{Int64: user.ID, Valid: true})
	require.NoError(t, err)
	require.Len(t, list, 3)

	for _, d := range list {
		require.Equal(t, "audio", d.Type.String)
	}
}

func TestListTextDictations(t *testing.T) {
	user := RandomUser(t)

	for i := 0; i < 4; i++ {
		RandomTextDictation(t, user)
	}
	RandomAudioDictation(t, user)

	list, err := testQueries.ListTextDictations(context.Background(), sql.NullInt64{Int64: user.ID, Valid: true})
	require.NoError(t, err)
	require.Len(t, list, 4)

	for _, d := range list {
		require.Equal(t, "text", d.Type.String)
	}
}

func TestListDictationsByUser(t *testing.T) {
	user := RandomUser(t)

	RandomTextDictation(t, user)
	RandomAudioDictation(t, user)

	list, err := testQueries.ListDictationsByUser(context.Background(), sql.NullInt64{Int64: user.ID, Valid: true})
	require.NoError(t, err)
	require.Len(t, list, 2)
}


