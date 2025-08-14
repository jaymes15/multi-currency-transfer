package users

import (
	responses "lemfi/simplebank/internal/apps/users/responses"
)

func (s *UserService) GetUser(username string) (responses.GetUserResponse, error) {
	// Get user from repository
	user, err := s.userRespository.GetUser(username)
	if err != nil {
		return responses.GetUserResponse{}, err
	}

	// Convert to response (GetUserRow only has limited fields)
	response := responses.GetUserResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		// Note: ID, Role, and UpdatedAt are not available from GetUserRow
		// These would need to be added to the SQL query if needed
	}

	return response, nil
}
