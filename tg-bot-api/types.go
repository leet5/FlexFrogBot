package tg_bot_api

type Update struct {
	UpdateID int            `json:"update_id"`
	Message  *Message       `json:"message,omitempty"`
	Callback *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	MessageID      int         `json:"message_id"`
	Text           string      `json:"text"`
	Chat           Chat        `json:"chat"`
	NewChatMembers *[]User     `json:"new_chat_members,omitempty"`
	Photo          []PhotoSize `json:"photo,omitempty"`
	From           *User       `json:"from,omitempty"`
	Document       *Document   `json:"document,omitempty"`
}

type Document struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
}

type PhotoSize struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size,omitempty"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	Data    string   `json:"data"`
	Message *Message `json:"message"`
	From    *User    `json:"from,omitempty"`
}

type GetUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type InlineKeyboardButton struct {
	Text         string  `json:"text"`
	CallbackData string  `json:"callback_data"`
	WebApp       *WebApp `json:"web_app,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type MessagePayload struct {
	ChatID      int64                `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type WebApp struct {
	URL string `json:"url"`
}

type FileInfo struct {
	FilePath string `json:"file_path"`
}
