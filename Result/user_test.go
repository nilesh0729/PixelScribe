package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nilesh0729/PixelScribe/util"
	"github.com/stretchr/testify/require"
)

func RandomUser(t *testing.T) User {
	arg := CreateUsersParams{
		Name:         sql.NullString{String: util.RandomString(8), Valid: true},
		Username:     util.RandomUsername(),
		Email:        util.RandomEmail(),
		PasswordHash: sql.NullString{String: util.RandomString(8), Valid: true},
	}
	user, err := testQueries.CreateUsers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Name, arg.Name)
	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.PasswordHash, arg.PasswordHash)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	RandomUser(t)
}

func TestGetuser(t *testing.T) {
	user := RandomUser(t)
	user1, err := testQueries.GetUsers(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user1)

	require.Equal(t, user.ID, user1.ID)
	require.Equal(t, user.Name, user1.Name)
	require.Equal(t, user.Username, user1.Username)
	require.Equal(t, user.Email, user1.Email)
	require.Equal(t, user.PasswordHash, user1.PasswordHash)

	require.WithinDuration(t, user.CreatedAt, user1.CreatedAt, time.Second)
}

func TestListUser(t *testing.T){

	for i:= 0; i<10; i++{
		RandomUser(t)
	}
	arg := ListUsersParams{
		Limit: 5,
		Offset: 5,
	}

	users, err:= testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t,users)
	require.Len(t, users, 5)

}

func TestUpdateUser(t *testing.T){
	user1 := RandomUser(t)
	arg := UpdateUsersParams{
		ID: user1.ID,
		Name: user1.Name,
		Email: util.RandomEmail(),
		PasswordHash: sql.NullString{String: util.RandomString(8), Valid :true},
	}
	user2, err := testQueries.UpdateUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.ID, user1.ID)
	require.Equal(t, user2.Name, user1.Name)
	require.Equal(t, user2.Email, arg.Email)
	require.Equal(t, user2.PasswordHash, arg.PasswordHash)
	require.Equal(t, user1.Username, user2.Username)

	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.WithinDuration(t, user2.UpdatedAt, time.Now(), time.Second)
	
}

func TestDeleteUser(t *testing.T){
	user := RandomUser(t)
	err := testQueries.DeleteUsers(context.Background(), user.Username)
	require.NoError(t, err)

	user2, err := testQueries.GetUsers(context.Background(), user.Username)
	require.EqualError(t, err ,sql.ErrNoRows.Error())
	require.Zero(t, user2)
}
