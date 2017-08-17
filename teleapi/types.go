package teleapi

// Update ...
type Update struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

// Message ...
type Message struct {
	MessageID int64  `json:"message_id"`
	From      User   `json:"from"`
	Date      int64  `json:"date"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	Entities []Entity `json:"entities"`
}

// User ...
type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// Chat ...
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// EntityType ...
type EntityType int
func (e EntityType) String() string {
	return entities[e]
}
var entities = [...]string{
	"mention",
	"hashtag",
	"bot_command",
	"url",
	"email",
	"bold",
	"italic",
	"code",
	"pre",
	"text_link",
	"text_mention",
}

// ...
const (
	MentionEntity EntityType = iota
	HashtagEntity
	BotCommandEntity
	URLEntity
	EmailEntity
	BoldEntity
	ItalicEntity
	CodeEntity
	PreEntity
	TextLinkEntity
	TextMentionEntity
)
// Entity ...
type Entity struct {
	Type string `json:"type"`
	Offset int `json:"offset"`
	Length int `json:"length"`
	URL string `json:"url"`
	User User `json:"user"`
}

type getUpdatesResp struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}

type sendMessageReq struct {
	ChatID					interface{}	`json:"chat_id"` // string || number
	Text					string		`json:"text"`
	ParseMode				bool		`json:"parse_mode,omitempty"`
	DisableWebPagePreview	bool		`json:"disable_web_page_preview,omitempty"`
	DisableNotification		bool		`json:"disable_notification,omitempty"`
	ReplyToMessageID		int64		`json:"reply_to_message_id,omitempty"`
	ReplyMarkup				interface{}	`json:"reply_markup,omitempty"`
}