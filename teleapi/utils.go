package teleapi

import (
	"log"
)

// TagsFromMessage ...
func TagsFromMessage(message *Message) []string {
	if message == nil {
		return make([]string, 0, 0)
	}
	tags := make([]string, 0, len(message.Entities))
	log.Printf("[Info] message is: %v\n", message)
	for _, e := range message.Entities {
		if e.Type != "hashtag" {
			continue
		}
		tag := message.Text[e.Offset: len(message.Text) - e.Length]
		tags = append(tags, tag)
	}
	log.Printf("[Info] tags is: %v\n", tags)
	return tags
}