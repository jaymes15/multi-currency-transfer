package db

import (
	"context"
	"testing"

	"lemfi/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	// Use helper function to create a random user
	user := createRandomUser(t)

	// Verify the returned user data
	require.NotEmpty(t, user.Username)
	require.NotEmpty(t, user.FullName)
	require.NotEmpty(t, user.Email)
	require.NotZero(t, user.CreatedAt)

	// Verify user was actually created in database
	createdUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.FullName, createdUser.FullName)
	require.Equal(t, user.Email, createdUser.Email)
	require.Equal(t, user.CreatedAt, createdUser.CreatedAt)
}

func TestCreateUserDuplicateUsername(t *testing.T) {
	// Create first user using helper
	firstUser := createRandomUser(t)

	// Try to create another user with same username
	duplicateUser, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       firstUser.Username,
		HashedPassword: "differentpassword",
		FullName:       "Second User",
		Email:          "second@example.com",
	})

	require.Error(t, err)
	require.Empty(t, duplicateUser)
	require.Contains(t, err.Error(), "duplicate key value")
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	// Create first user using helper
	firstUser := createRandomUser(t)

	// Try to create another user with same email
	duplicateUser, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       "user2",
		HashedPassword: "differentpassword",
		FullName:       "Second User",
		Email:          firstUser.Email,
	})

	require.Error(t, err)
	require.Empty(t, duplicateUser)
	require.Contains(t, err.Error(), "duplicate key value")
}

func TestCreateUserEmptyFields(t *testing.T) {
	// Test with empty username
	_, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       "",
		HashedPassword: "hashedpassword123",
		FullName:       "Test User",
		Email:          "test@example.com",
	})
	require.Error(t, err)

	// Test with empty hashed password
	_, err = testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       "testuser",
		HashedPassword: "",
		FullName:       "Test User",
		Email:          "test@example.com",
	})
	require.Error(t, err)

	// Test with empty full name
	_, err = testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       "testuser",
		HashedPassword: "hashedpassword123",
		FullName:       "",
		Email:          "test@example.com",
	})
	require.Error(t, err)

	// Test with empty email
	_, err = testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       "testuser",
		HashedPassword: "hashedpassword123",
		FullName:       "Test User",
		Email:          "",
	})
	require.Error(t, err)
}

func TestCreateUserSpecialCharacters(t *testing.T) {
	// Test with special characters in username
	username := "test-user_123_" + util.RandomOwner()
	hashedPassword := "hashedpassword123"
	fullName := "Test User"
	email := util.RandomEmail()

	user, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       fullName,
		Email:          email,
	})

	require.NoError(t, err)
	require.Equal(t, username, user.Username)

	// Test with special characters in full name
	username2 := "testuser2_" + util.RandomOwner()
	fullName2 := "Test User Jr. (II)"

	user2, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username2,
		HashedPassword: hashedPassword,
		FullName:       fullName2,
		Email:          util.RandomEmail(),
	})

	require.NoError(t, err)
	require.Equal(t, fullName2, user2.FullName)
}
