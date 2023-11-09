package models

// Friend represents a user's friend.
type Friend struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
