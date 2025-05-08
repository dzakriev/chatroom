package main

import (
	"context"
	"log/slog"
)

type User struct {
	Name     string
	Username string
	Id       int
}

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Room struct {
	Name string
	Id   int
}

type UserInRoom struct {
	UserId int
	RoomId int
}

func Write() string {
	return "string"
}

type MessageClient interface {
	Create(id string) Message
	Read(id string) Message
	ReadAll() []Message
	Delete(id string) Message
	DeleteAll() int
}

var messageClient MessageClient

type MessageClientDb struct{}

// type MessageClientMock struct{}

func (client MessageClientDb) Create(text string) Message {
	var id int
	var err = dbPool.QueryRow(context.Background(), "insert into message (text) values ($1) returning id", text).Scan(&id)
	if err != nil {
		slog.Error("Creating message failed", "error", err)
		return Message{}
	}

	msg := Message{ID: id, Text: text}
	return msg
}

func (client MessageClientDb) Delete(id string) Message {
	var msg Message
	err := dbPool.QueryRow(context.Background(), "delete from message where id=$1 returning id, text", id).Scan(&msg.ID, &msg.Text)
	if err != nil {
		slog.Error("Deleting message failed", "error", err)
	}

	return msg
}

func (client MessageClientDb) DeleteAll() int {
	commandTag, err := dbPool.Exec(context.Background(), "delete from message")
	if err != nil {
		slog.Error("Deleting message failed", "error", err)
	}

	return int(commandTag.RowsAffected())
}

func (client MessageClientDb) Read(id string) Message {
	var msg Message
	err := dbPool.QueryRow(context.Background(), "select id, text from message where id=$1", id).Scan(&msg.ID, &msg.Text)
	if err != nil {
		slog.Error("Reading message failed", "error", err, "id", id)
	}

	return msg
}

func (client MessageClientDb) ReadAll() []Message {
	rows, err := dbPool.Query(context.Background(), "select id, text from message")
	if err != nil {
		slog.Error("Reading failed", "error", err)
		return nil
	}

	messages := []Message{}

	for rows.Next() {
		msg := Message{}
		err = rows.Scan(&msg.ID, &msg.Text)
		if err != nil {
			slog.Error("Scan error", "error", err)
			return nil
		}
		messages = append(messages, msg)
	}

	return messages
}
