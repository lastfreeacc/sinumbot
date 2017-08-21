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

// NewMemo ...
func NewMemo(feed string, tags []string) Memo {
	memo := Memo{
		Feed: feed,
		Tags: tags,
	}
	return memo
}