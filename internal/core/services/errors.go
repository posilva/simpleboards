package services

import "fmt"

// LeaderboardNotFoundError ...
type LeaderboardNotFoundError struct {
	Name string
}

// Error interface implementation
func (e *LeaderboardNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Name)
}
