package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
    "lemfi/simplebank/util"
)

func TestGetUser(t *testing.T) {
	// Create a user using helper
	createdUser := createRandomUser(t)

	// Get the user
	retrievedUser, err := testQueries.GetUser(context.Background(), createdUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedUser)

	// Verify the retrieved user data matches the created user
	require.Equal(t, createdUser.Username, retrievedUser.Username)
	require.Equal(t, createdUser.FullName, retrievedUser.FullName)
	require.Equal(t, createdUser.Email, retrievedUser.Email)
	require.Equal(t, createdUser.CreatedAt, retrievedUser.CreatedAt)
}

func TestGetUserNotFound(t *testing.T) {
	// Try to get a non-existent user
	nonExistentUsername := "non_existent_user_get"
	user, err := testQueries.GetUser(context.Background(), nonExistentUsername)

	require.Error(t, err)
	require.Empty(t, user)
	require.Contains(t, err.Error(), "no rows in result set")
}

func TestGetUserEmptyUsername(t *testing.T) {
	// Try to get user with empty username
	user, err := testQueries.GetUser(context.Background(), "")

	require.Error(t, err)
	require.Empty(t, user)
	require.Contains(t, err.Error(), "no rows in result set")
}

func TestGetUserSpecialCharacters(t *testing.T) {
	// Create a user with special characters in username
    username := "user-with_special.chars_get123_" + util.RandomOwner()
	hashedPassword := "hashedpassword123"
	fullName := "Special User Get"
    email := util.RandomEmail()

	// Create the user
	createdUser, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       fullName,
		Email:          email,
	})
	require.NoError(t, err)

	// Get the user with special characters
	retrievedUser, err := testQueries.GetUser(context.Background(), username)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedUser)

	// Verify the retrieved user data
	require.Equal(t, username, retrievedUser.Username)
	require.Equal(t, fullName, retrievedUser.FullName)
	require.Equal(t, email, retrievedUser.Email)
	require.Equal(t, createdUser.CreatedAt, retrievedUser.CreatedAt)
}

func TestGetUserCaseSensitivity(t *testing.T) {
	// Create a user with specific case
    username := "TestUserCase_" + util.RandomOwner()
	hashedPassword := "hashedpassword123"
	fullName := "Case Test User"
    email := util.RandomEmail()

	// Create the user
	_, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       fullName,
		Email:          email,
	})
	require.NoError(t, err)

	// Try to get user with different case (should fail if case-sensitive)
	differentCaseUsername := "testusercase"
    user, err := testQueries.GetUser(context.Background(), differentCaseUsername)

	// PostgreSQL is case-sensitive by default, so this should fail
	require.Error(t, err)
	require.Empty(t, user)
	require.Contains(t, err.Error(), "no rows in result set")

	// Get user with correct case
	correctCaseUser, err := testQueries.GetUser(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, username, correctCaseUser.Username)
}

func TestGetUserAfterUpdate(t *testing.T) {
	// Create a user using helper
	user := createRandomUser(t)

	// Get user before update
	userBeforeUpdate, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.FullName, userBeforeUpdate.FullName)

	// Update the user
	updatedFullName := "Updated Name"
	_, err = testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		FullName: pgtype.Text{String: updatedFullName, Valid: true},
	})
	require.NoError(t, err)

	// Get user after update
	userAfterUpdate, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, updatedFullName, userAfterUpdate.FullName)
	require.NotEqual(t, userBeforeUpdate.FullName, userAfterUpdate.FullName)
}

func TestGetUserMultipleRetrievals(t *testing.T) {
	// Create a user using helper
	createdUser := createRandomUser(t)

	// Get the user multiple times
	for i := 0; i < 5; i++ {
		retrievedUser, err := testQueries.GetUser(context.Background(), createdUser.Username)
		require.NoError(t, err)
		require.NotEmpty(t, retrievedUser)
		require.Equal(t, createdUser.Username, retrievedUser.Username)
		require.Equal(t, createdUser.FullName, retrievedUser.FullName)
		require.Equal(t, createdUser.Email, retrievedUser.Email)
	}
}
