package main

import (
	"context"
	rabbitclient "github.com/kiaplayer/rabbitmq-report-system/internal/rabbit-client"
	reportsusedcar "github.com/kiaplayer/rabbitmq-report-system/internal/reports/used-car"
	"log"
	"os"
)

func main() {
	amqpUri := os.Getenv("AMQP_URI")
	if amqpUri == "" {
		amqpUri = "amqp://guest:guest@127.0.0.1:5672/"
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	inputQueue := "reports.used-car.limits-wanted-info"
	rabbitClient, err := rabbitclient.CreateRabbitClient(amqpUri, "reports.results", inputQueue)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := rabbitClient.Shutdown(); err != nil {
			logger.Printf("Client shutdown error: %s\n", err.Error())
		}
	}()
	logger.Printf("Listening queue \"%s\"...", inputQueue)
	handler := reportsusedcar.SubreportLimitsWantedInfoHandler{RabbitClient: rabbitClient}
	if err := rabbitClient.ProcessMessages(context.Background(), &handler); err != nil {
		logger.Printf("Error while hadling deliveries: %s\n", err.Error())
	}
}
