package customerrors

import "fmt"

// NotFoundError represents a 404 error for a specific resource
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}
