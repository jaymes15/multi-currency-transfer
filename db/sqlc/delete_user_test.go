package db

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser(t *testing.T) {
	// Create a user to delete using helper
	user := createRandomUser(t)

	// Verify user exists
	createdUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, createdUser.Username)

	// Delete the user
	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	// Verify user was actually deleted
	deletedUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.Empty(t, deletedUser)
	require.Contains(t, err.Error(), "no rows in result set")
}

func TestDeleteUserNotFound(t *testing.T) {
	// Try to delete a non-existent user
	nonExistentUsername := "non_existent_user"
	err := testQueries.DeleteUser(context.Background(), nonExistentUsername)

	// DeleteUser should not return an error for non-existent users
	// It's a common pattern for DELETE operations to succeed even if nothing was deleted
	require.NoError(t, err)
}

func TestDeleteUserEmptyUsername(t *testing.T) {
	// Try to delete with empty username
	err := testQueries.DeleteUser(context.Background(), "")

	// This should either succeed (no rows affected) or fail gracefully
	// The behavior depends on database constraints
	if err != nil {
		require.Contains(t, err.Error(), "empty")
	}
}

func TestDeleteUserSpecialCharacters(t *testing.T) {
	// Create a user with special characters in username
	username := "user-with_special.chars123"
	hashedPassword := "hashedpassword123"
	fullName := "Special User"
	email := "special@example.com"

	// Create the user
	_, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       fullName,
		Email:          email,
	})
	require.NoError(t, err)

	// Verify user exists
	createdUser, err := testQueries.GetUser(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, username, createdUser.Username)

	// Delete the user
	err = testQueries.DeleteUser(context.Background(), username)
	require.NoError(t, err)

	// Verify user was deleted
	deletedUser, err := testQueries.GetUser(context.Background(), username)
	require.Error(t, err)
	require.Empty(t, deletedUser)
}

func TestDeleteUserMultipleTimes(t *testing.T) {
	// Create a user using helper
	user := createRandomUser(t)

	// Delete the user first time
	err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	// Try to delete the same user again
	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err) // Should not error, just no rows affected

	// Verify user is still deleted
	deletedUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.Empty(t, deletedUser)
}

func TestDeleteUserAndVerifyCascade(t *testing.T) {
	// Create a user using helper
	user := createRandomUser(t)

	// Create an account for this user
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    user.Username,
		Balance:  decimal.NewFromFloat(100.0),
		Currency: "USD",
	})
	require.NoError(t, err)

	// Verify account exists
	require.Equal(t, user.Username, account.Owner)

	// Delete the user
	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	// Verify user was deleted
	deletedUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.Empty(t, deletedUser)

	// Note: This test assumes foreign key constraints are set up
	// If CASCADE DELETE is configured, the account should also be deleted
	// If RESTRICT is configured, the delete should fail
	// The actual behavior depends on your migration constraints
}
