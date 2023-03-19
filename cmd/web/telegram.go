package main

type Update struct {
	UpdateId int64     `json:"update_id""`
	Message  TgMessage `json:"message""`
}

type TgMessage struct {
	MessageId       int64  `json:"message_id"`
	MessageThreadId int64  `json:"message_thread_id"`
	Chat            TgChat `json:"chat"`
	Text            string `json:"text"`
	Date            int64  `json:"date"`
	From            TgFrom `json:"from"`
}

type TgChat struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	ChatType  string `json:"type"` // private,group,supergroup,supergroup
	UserName  string `json:"username"`
}

type TgFrom struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

const (
	ChatTypePrivate    = "private"
	ChatTypeGroup      = "group"
	ChatTypeSuperGroup = "supergroup"
	ChatTypeChannel    = "supergroup"
)
