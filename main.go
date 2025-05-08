package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mychat/chat"
	"mychat/config"
	"mychat/db"
	"mychat/handlers"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

var cfg config.Config
var router *gin.Engine

func main() {
	mustInitConfig()
	mustInitDb()

	initRouter()
	runServer()

	waitForShutdown()
}

func mustInitConfig() {
	err := cleanenv.ReadConfig("config/local.yaml", &cfg)
	if err != nil {
		slog.Error("failed to read config", "error", err)
		os.Exit(1)
	}
}

func mustInitDb() {
	handlers.InjectClient(db.MessageClientDb{})
	err := db.InitDB(cfg.DbConnectionString)
	if err != nil {
		slog.Error("failed to start db", "error", err)
		os.Exit(1)
	}
}

func initRouter() {
	router = gin.New()
	router.Use(corsMiddleware())
	router.Use(slogMiddleware())
	router.Use(gin.Recovery())
	router.LoadHTMLFiles(cfg.IndexFilePath)
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, cfg.IndexFilePath, nil)
	})

	hub := chat.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		chat.ServeWs(hub, c.Writer, c.Request)
	})

	router.GET("/message", handlers.MessageGetAll)
	router.GET("/message/:id", handlers.MessageGet)
	router.POST("/message", handlers.MessagePostAndBroadcast(hub))
	router.DELETE("/message/:id", handlers.MessageDelete)
	router.DELETE("/message", handlers.MessageDeleteAll)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		c.Next()
	}
}

func slogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		slog.Info("HTTP request",
			"status", status,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency.String(),
			"error", errorMessage,
		)
	}
}

func runServer() {
	go func() {
		slog.Info("Server is starting")
		err := router.Run("localhost:8080")
		if err != nil {
			slog.Error("Server starting failed", "error", err)
			os.Exit(1)
		}
	}()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	slog.Info("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-shutdown():
		slog.Info("Server gracefully stopped")
	case <-ctx.Done():
		slog.Info("Server forced to shutdown")
	}
}

func shutdown() chan struct{} {
	stopped := make(chan struct{})
	go func() {
		db.DbPool.Close()
		stopped <- struct{}{}
	}()
	return stopped
}
