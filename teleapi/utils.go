package teleapi

// TagsFromMessage ...
func TagsFromMessage(message Message) []string {
	tags := make([]string, 0, len(message.Entities))
	for _, e := range message.Entities {
		if e.Type != "hashtag" {
			continue
		}
		f := e.Offset
		l := e.Offset + e.Length
		tag := message.Text[f:l]
		tags = append(tags, tag)
	}
	return tags
}