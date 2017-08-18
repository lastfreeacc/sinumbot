package store

// NewMemo ...
func NewMemo(feed string, tags []string) Memo {
	memo := Memo{
		Feed: feed,
		Tags: tags,
	}
	return memo
}