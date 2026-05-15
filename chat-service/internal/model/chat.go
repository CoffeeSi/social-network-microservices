package model

import "time"

type Chat struct {
	Id             string
	Name           string
	IsGroup        bool
	ParticipantIds []string
	CreatorId      string
	CreatedAt      time.Time
	LastMessage    *ChatMessage
}

type ChatMessage struct {
	Id        string
	ChatId    string
	SenderId  string
	Content   string
	IsEdited  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
