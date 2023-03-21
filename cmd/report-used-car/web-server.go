package main

import (
	"context"
	"github.com/gin-gonic/gin"
	rabbitclient "github.com/kiaplayer/rabbitmq-report-system/internal/rabbit-client"
	reportsusedcar "github.com/kiaplayer/rabbitmq-report-system/internal/reports/used-car"
	"log"
	"os"
)

func setupRouter(rabbitClient *rabbitclient.RabbitClient) *gin.Engine {
	r := gin.Default()
	reportsusedcar.SetupRoutes(r, rabbitClient)
	return r
}

func main() {
	amqpUri := os.Getenv("AMQP_URI")
	if amqpUri == "" {
		amqpUri = "amqp://guest:guest@127.0.0.1:5672/"
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	inputQueue := "reports.used-car.results"
	rabbitClient, err := rabbitclient.CreateRabbitClient(amqpUri, "reports.tasks", inputQueue)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := rabbitClient.Shutdown(); err != nil {
			logger.Printf("Client shutdown error: %s\n", err.Error())
		}
	}()
	go func() {
		logger.Printf("Listening queue \"%s\"...", inputQueue)
		if err := rabbitClient.ProcessMessages(context.Background(), &reportsusedcar.ReportResultsHandler{}); err != nil {
			logger.Printf("Error while hadling deliveries: %s\n", err.Error())
		}
	}()
	r := setupRouter(rabbitClient)
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
