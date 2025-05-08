package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mychat/chat"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func initDB(conString string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	dbPool, err = pgxpool.New(ctxTimeout, conString)
	if err != nil {
		return err
	}

	return dbPool.Ping(ctxTimeout)
}

func messageDelete(c *gin.Context) {
	message := messageClient.Delete(c.Param("id"))
	c.JSON(200, gin.H{
		"message": message,
	})
}

func messageDeleteAll(c *gin.Context) {
	rowsDeleted := messageClient.DeleteAll()
	c.JSON(200, gin.H{
		"rowsDeleted": rowsDeleted,
	})
}

func messageGet(c *gin.Context) {
	message := messageClient.Read(c.Param("id"))
	c.JSON(200, gin.H{
		"message": message,
	})
}

func messageGetAll(c *gin.Context) {
	messages := messageClient.ReadAll()
	c.JSON(200, gin.H{
		"messages": messages,
	})
}

func messagePostAndBroadcast(h *chat.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := messagePost(c)
		jsonData, err := json.Marshal(msg)
		if err != nil {
			slog.Error("JSON serialization failed for object text:"+msg.Text, "error", err, "id", msg.ID)
		}
		h.Broadcast(jsonData)
	}
}

func messagePost(c *gin.Context) *Message {
	var msg Message
	if err := c.BindJSON(&msg); err != nil {
		fmt.Print(err)
		return nil
	}
	msg = messageClient.Create(msg.Text)
	c.JSON(200, gin.H{
		"id":   msg.ID,
		"text": msg.Text,
	})
	return &msg
}
