package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mychat/chat"
	"mychat/db"

	"github.com/gin-gonic/gin"
)

var messageClient db.IMessageClient

func InjectClient(client db.IMessageClient) {
	messageClient = client
}

func MessageDelete(c *gin.Context) {
	message, ok := messageClient.Delete(c.Param("id"))

	if !ok {
		c.JSON(404, gin.H{
			"message": "Resource not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}

func MessageDeleteAll(c *gin.Context) {
	rowsDeleted := messageClient.DeleteAll()
	c.JSON(200, gin.H{
		"rowsDeleted": rowsDeleted,
	})
}

func MessageGet(c *gin.Context) {
	message, ok := messageClient.Read(c.Param("id"))

	if !ok {
		c.JSON(404, gin.H{
			"message": "Resource not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}

func MessageGetAll(c *gin.Context) {
	messages := messageClient.ReadAll()
	c.JSON(200, gin.H{
		"messages": messages,
	})
}

func MessagePostAndBroadcast(h *chat.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := MessagePost(c)
		jsonData, err := json.Marshal(msg)
		if err != nil {
			slog.Error("JSON serialization failed for object text:"+msg.Text, "error", err, "id", msg.ID)
		}
		h.Broadcast(jsonData)
	}
}

func MessagePost(c *gin.Context) *db.Message {
	var msg db.Message
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
