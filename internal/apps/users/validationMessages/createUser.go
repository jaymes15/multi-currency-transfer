package users

var CreateUserValidationMessages = map[string]string{
	"Username.required": "username is required.",
	"Password.required": "password is required.",
	"Password.min":      "password must be at least 6 characters long.",
	"FullName.required": "full name is required.",
	"Email.required":    "email is required.",
	"Email.email":       "email must be a valid email address.",
}
