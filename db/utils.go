package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DbPool *pgxpool.Pool

type MessageClientDb struct{}

func InitDB(conString string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DbPool, err = pgxpool.New(ctxTimeout, conString)
	if err != nil {
		return err
	}

	return DbPool.Ping(ctxTimeout)
}

func (client *MessageClientDb) Create(text string) Message {
	var id int
	var err = DbPool.QueryRow(context.Background(), "insert into message (text) values ($1) returning id", text).Scan(&id)
	if err != nil {
		slog.Error("Creating message failed", "error", err)
		return Message{}
	}

	msg := Message{ID: id, Text: text}
	return msg
}

func (client *MessageClientDb) Delete(id string) (Message, bool) {
	var msg Message
	found := true

	err := DbPool.QueryRow(context.Background(), "delete from message where id=$1 returning id, text", id).Scan(&msg.ID, &msg.Text)
	if err != nil {
		slog.Error("Deleting message failed", "error", err)
		found = false
	}

	return msg, found
}

func (client *MessageClientDb) DeleteAll() int {
	commandTag, err := DbPool.Exec(context.Background(), "delete from message")
	if err != nil {
		slog.Error("Deleting message failed", "error", err)
	}

	return int(commandTag.RowsAffected())
}

func (client *MessageClientDb) Read(id string) (Message, bool) {
	var msg Message
	found := true

	err := DbPool.QueryRow(context.Background(), "select id, text from message where id=$1", id).Scan(&msg.ID, &msg.Text)
	if err != nil {
		slog.Error("Reading message failed", "error", err, "id", id)
		found = false
	}

	return msg, found
}

func (client *MessageClientDb) ReadAll() []Message {
	rows, err := DbPool.Query(context.Background(), "select id, text from message")
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
