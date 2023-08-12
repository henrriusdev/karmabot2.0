package main

import "strings"

type Updates struct {
	Ok     bool     `json:"ok"`
	Result []Result `json:"result"`
}

type Result struct {
	EditedMessage EditedMessage `json:"edited_message"`
	Message       Message       `json:"message"`
	UpdateID      int64         `json:"update_id"`
}

type EditedMessage struct {
	Chat            Chat                        `json:"chat"`
	Date            int64                       `json:"date"`
	EditDate        int64                       `json:"edit_date"`
	From            From                        `json:"from"`
	IsTopicMessage  bool                        `json:"is_topic_message"`
	MessageID       int64                       `json:"message_id"`
	MessageThreadID *int64                      `json:"message_thread_id"`
	ReplyToMessage  EditedMessageReplyToMessage `json:"reply_to_message"`
	Text            string                      `json:"text"`
}

type Chat struct {
	ID       int64  `json:"id"`
	IsForum  bool   `json:"is_forum"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Username string `json:"username"`
}

func (c *Chat) IsPrivate() bool {
	return c.Type == "private"
}

func (c *Chat) IsGroup() bool {
	return c.Type == "group" || c.Type == "supergroup"
}

func (c *Chat) IsChannel() bool {
	return c.Type == "channel"
}

type From struct {
	FirstName    string `json:"first_name"`
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	LanguageCode string `json:"language_code"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
}

type EditedMessageReplyToMessage struct {
	Chat            Chat   `json:"chat"`
	Date            int64  `json:"date"`
	EditDate        int64  `json:"edit_date"`
	From            From   `json:"from"`
	IsTopicMessage  bool   `json:"is_topic_message"`
	MessageID       int64  `json:"message_id"`
	MessageThreadID int64  `json:"message_thread_id"`
	Text            string `json:"text"`
}

type Message struct {
	Animation          Animation             `json:"animation"`
	Caption            string                `json:"caption"`
	CaptionEntities    []Entity              `json:"caption_entities"`
	Chat               Chat                  `json:"chat"`
	Date               int64                 `json:"date"`
	Document           Document              `json:"document"`
	Entities           []Entity              `json:"entities"`
	From               From                  `json:"from"`
	IsTopicMessage     bool                  `json:"is_topic_message"`
	MessageID          int64                 `json:"message_id"`
	MessageThreadID    *int64                `json:"message_thread_id"`
	NewChatMember      From                  `json:"new_chat_member"`
	NewChatMembers     []From                `json:"new_chat_members"`
	NewChatParticipant From                  `json:"new_chat_participant"`
	Photo              []Thumb               `json:"photo"`
	ReplyToMessage     MessageReplyToMessage `json:"reply_to_message"`
	Sticker            Sticker               `json:"sticker"`
	Text               string                `json:"text"`
}

func (m *Message) IsCommand() bool {
	entity := (m.Entities)[0]
	return entity.Offset == 0 && entity.IsCommand()
}

func (m *Message) Command() string {
	command := m.CommandWithAt()

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

func (m *Message) CommandWithAt() string {
	if !m.IsCommand() {
		return ""
	}

	entity := (m.Entities)[0]
	return m.Text[1:entity.Length]
}

type Animation struct {
	Duration     int64  `json:"duration"`
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	Height       int64  `json:"height"`
	MIMEType     string `json:"mime_type"`
	Thumb        Thumb  `json:"thumb"`
	Thumbnail    Thumb  `json:"thumbnail"`
	Width        int64  `json:"width"`
}

type Thumb struct {
	FileID       string `json:"file_id"`
	FileSize     int64  `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	Height       int64  `json:"height"`
	Width        int64  `json:"width"`
}

type Entity struct {
	Length int64  `json:"length"`
	Offset int64  `json:"offset"`
	Type   string `json:"type"`
}

func (e *Entity) IsCommand() bool {
	return e.Type == "bot_command"
}

type Document struct {
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	MIMEType     string `json:"mime_type"`
	Thumb        Thumb  `json:"thumb"`
	Thumbnail    Thumb  `json:"thumbnail"`
}

type MessageReplyToMessage struct {
	Caption           string            `json:"caption"`
	CaptionEntities   []Entity          `json:"caption_entities"`
	Chat              Chat              `json:"chat"`
	Date              int64             `json:"date"`
	EditDate          *int64            `json:"edit_date"`
	ForumTopicCreated ForumTopicCreated `json:"forum_topic_created"`
	From              From              `json:"from"`
	IsTopicMessage    bool              `json:"is_topic_message"`
	MessageID         int64             `json:"message_id"`
	MessageThreadID   *int64            `json:"message_thread_id"`
	Photo             []Thumb           `json:"photo"`
	Sticker           Sticker           `json:"sticker"`
	Text              string            `json:"text"`
}

type Sticker struct {
	Emoji        string `json:"emoji"`
	FileID       string `json:"file_id"`
	FileSize     int64  `json:"file_size"`
	FileUniqueID string `json:"file_unique_id"`
	Height       int64  `json:"height"`
	IsAnimated   bool   `json:"is_animated"`
	IsVideo      bool   `json:"is_video"`
	SetName      string `json:"set_name"`
	Thumb        Thumb  `json:"thumb"`
	Thumbnail    Thumb  `json:"thumbnail"`
	Type         string `json:"type"`
	Width        int64  `json:"width"`
}

type ForumTopicCreated struct {
	IconColor         int64  `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
	Name              string `json:"name"`
}
