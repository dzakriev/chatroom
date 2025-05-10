package grpc

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendAuth(c *gin.Context) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error", "error", err)
	}

	client := NewAuthServiceClient(conn)

	request := &AuthRequest{
		Credentials: []*Credential{
			{
				Name:     "in Duis mollit laborum",
				Password: "in elit aute id velit",
				Type:     AccountType_ACCOUNT_TYPE_ADMIN,
			},
		},
	}
	request.Credentials = append(request.Credentials, &Credential{
		Name:     "in Duis mollit laborum",
		Password: "in elit aute id velit",
		Type:     AccountType_ACCOUNT_TYPE_USER,
	})

	result, err := client.Auth(c, request)
	if err != nil {
		slog.Error("Error", "error", err)
	}

	slog.Info("Request", "request", request)
	slog.Info("Result", "result", result)
	c.JSON(200, gin.H{
		"result": result.Success,
	})
}
