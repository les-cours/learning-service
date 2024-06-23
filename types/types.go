package types

import (
	"time"
)

type Comment struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RepliedTo  string    `json:"replied_to"`
	Content    string    `json:"content"`
	DocumentID string    `json:"document_id"`
	Timestamp  time.Time `json:"timestamp"`
	Edited     bool      `json:"edited"`
	IsTeacher  bool      `json:"is_teacher"`
}

type Message struct {
	ID        string
	RoomID    string
	Message   string
	Timestamp string
	IsTeacher bool
	Owner     *User
}

type Room struct {
	ID       string     `json:"roomID"`
	Name     string     `json:"name"`
	Teacher  *User      `json:"teacher"`
	Users    []*User    `json:"users"`
	Messages []*Message `json:"messages"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Avatar    string `json:"avatar"`
	Paid      bool   `json:"isActive"`
}
