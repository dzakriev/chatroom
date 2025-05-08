package tests

import (
	"encoding/json"
	"mychat/db"
	"mychat/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestDbConnection(t *testing.T) {
}

func TestMessageGet(t *testing.T) {
	mock := NewMessageClientMock()
	mock.Messages["123"] = db.Message{ID: 123, Text: "found me"}
	handlers.InjectClient(mock)

	router := gin.Default()
	router.GET("/messages/:id", handlers.MessageGet)

	req, _ := http.NewRequest("GET", "/messages/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	message := resp["message"].(map[string]interface{})
	assert.Equal(t, float64(123), message["id"])
	assert.Equal(t, "found me", message["text"])
}
