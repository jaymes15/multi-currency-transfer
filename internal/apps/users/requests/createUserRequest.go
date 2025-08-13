package users

type CreateUserRequest struct {
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required,min=6"`
	FullName       string `json:"full_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	HashedPassword string `json:"-"` // This field is used internally, not exposed in JSON
}
