package store

// User ...
type User struct {
	ID int64
	Memos []*Memo
}

// Memo ...
type Memo struct {
	Feed string
	Preview string
	Tags []string
}