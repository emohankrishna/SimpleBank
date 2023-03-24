package db

import (
	"context"
	"testing"
	"time"

	"github.com/emohankrishna/Simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(10),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, hashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	fetchedUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)
	require.Equal(t, user.Username, fetchedUser.Username)
	require.Equal(t, user.HashedPassword, fetchedUser.HashedPassword)
	require.Equal(t, user.FullName, fetchedUser.FullName)
	require.Equal(t, user.Email, fetchedUser.Email)
	require.WithinDuration(t, user.CreatedAt, fetchedUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, fetchedUser.PasswordChangedAt, time.Second)
}
