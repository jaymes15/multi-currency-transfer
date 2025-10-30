package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
    "lemfi/simplebank/util"
)

func TestListUsers(t *testing.T) {
	// Create multiple users using helper
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	user3 := createRandomUser(t)

	// List all users
	users, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)

	// Verify our created users are in the list
	userMap := make(map[string]ListUsersRow)
	for _, user := range users {
		userMap[user.Username] = user
	}

	// Check that our created users exist in the list
	require.Contains(t, userMap, user1.Username)
	require.Contains(t, userMap, user2.Username)
	require.Contains(t, userMap, user3.Username)

	// Verify user details
	require.Equal(t, user1.FullName, userMap[user1.Username].FullName)
	require.Equal(t, user1.Email, userMap[user1.Username].Email)
	require.Equal(t, user2.FullName, userMap[user2.Username].FullName)
	require.Equal(t, user2.Email, userMap[user2.Username].Email)
	require.Equal(t, user3.FullName, userMap[user3.Username].FullName)
	require.Equal(t, user3.Email, userMap[user3.Username].Email)
}

func TestListUsersEmpty(t *testing.T) {
	// This test assumes we're starting with a clean database
	// or we need to clean up existing users first

	// For now, let's just test that the function doesn't error
	// even if there are no users
	users, err := testQueries.ListUsers(context.Background())

	// The function should work even with no users
	require.NoError(t, err)
	// users could be empty slice or contain existing users
	require.NotNil(t, users)
}

func TestListUsersOrdering(t *testing.T) {
	// Create users with specific usernames to test ordering
    username1 := "a_user_first_" + util.RandomOwner()
    username2 := "b_user_second_" + util.RandomOwner()
    username3 := "c_user_third_" + util.RandomOwner()

	// Create users in reverse order to test SQL ordering
    _, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username3,
		HashedPassword: "hashedpassword123",
		FullName:       "Third User",
        Email:          util.RandomEmail(),
	})
	require.NoError(t, err)

    _, err = testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username1,
		HashedPassword: "hashedpassword123",
		FullName:       "First User",
        Email:          util.RandomEmail(),
	})
	require.NoError(t, err)

    _, err = testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username2,
		HashedPassword: "hashedpassword123",
		FullName:       "Second User",
        Email:          util.RandomEmail(),
	})
	require.NoError(t, err)

	// List users
	users, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)

	// Find our test users in the list
	var foundUser1, foundUser2, foundUser3 *ListUsersRow
	for i, user := range users {
		if user.Username == username1 {
			foundUser1 = &users[i]
		}
		if user.Username == username2 {
			foundUser2 = &users[i]
		}
		if user.Username == username3 {
			foundUser3 = &users[i]
		}
	}

	// Verify all users were found
	require.NotNil(t, foundUser1)
	require.NotNil(t, foundUser2)
	require.NotNil(t, foundUser3)

	// Verify the ordering (should be alphabetical by username)
	// Find the positions of our users in the list
	var pos1, pos2, pos3 int
	for i, user := range users {
		if user.Username == username1 {
			pos1 = i
		}
		if user.Username == username2 {
			pos2 = i
		}
		if user.Username == username3 {
			pos3 = i
		}
	}

	// Verify alphabetical ordering: a_user_first < b_user_second < c_user_third
	require.Less(t, pos1, pos2)
	require.Less(t, pos2, pos3)
}

func TestListUsersAfterDelete(t *testing.T) {
	// Create a user
	user := createRandomUser(t)

	// Verify user exists in list
	usersBeforeDelete, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)

	userExists := false
	for _, u := range usersBeforeDelete {
		if u.Username == user.Username {
			userExists = true
			break
		}
	}
	require.True(t, userExists, "User should exist in list before deletion")

	// Delete the user
	err = testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	// Verify user no longer exists in list
	usersAfterDelete, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)

	userExists = false
	for _, u := range usersAfterDelete {
		if u.Username == user.Username {
			userExists = true
			break
		}
	}
	require.False(t, userExists, "User should not exist in list after deletion")
}

func TestListUsersMultipleCalls(t *testing.T) {
	// Create a user
	user := createRandomUser(t)

	// Call ListUsers multiple times
	users1, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)

	users2, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)

	users3, err := testQueries.ListUsers(context.Background())
	require.NoError(t, err)

	// All calls should return the same result
	require.Equal(t, len(users1), len(users2))
	require.Equal(t, len(users2), len(users3))

	// Verify our user exists in all results
	userExistsIn1 := false
	userExistsIn2 := false
	userExistsIn3 := false

	for _, u := range users1 {
		if u.Username == user.Username {
			userExistsIn1 = true
			break
		}
	}

	for _, u := range users2 {
		if u.Username == user.Username {
			userExistsIn2 = true
			break
		}
	}

	for _, u := range users3 {
		if u.Username == user.Username {
			userExistsIn3 = true
			break
		}
	}

	require.True(t, userExistsIn1)
	require.True(t, userExistsIn2)
	require.True(t, userExistsIn3)
}
