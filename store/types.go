package store

// User ...
type User struct {
	ID int64
	Memos []*Memo
}

// Memo ...
type Memo struct {
	Entry string
	Preview string
	Tags []string
}