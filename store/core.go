package store

// Store ...
type Store interface {
	GetUser(userID int64) (*User, error)
	SaveMemo(userID int64, memo string) error
}