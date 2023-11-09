// domain/models/user.go

package models

// User represents a user.
type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
