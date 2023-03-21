package rabbit_client

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	publishExchange string
	subscribeQueue  amqp.Queue
	deliveries      <-chan amqp.Delivery
	connErr         chan *amqp.Error
}

type DeliveryHandler interface {
	Handle(ctx context.Context, data []byte) error
}

func CreateRabbitClient(url, publishExchange, subscribeQueue string) (*RabbitClient, error) {
	rabbitClient := &RabbitClient{
		conn:            nil,
		channel:         nil,
		publishExchange: publishExchange,
		deliveries:      make(chan amqp.Delivery),
		connErr:         make(chan *amqp.Error),
	}
	var err error
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	rabbitClient.conn, err = amqp.DialConfig(url, config)
	if err != nil {
		return nil, fmt.Errorf("rabbit dial: %s", err)
	}
	rabbitClient.channel, err = rabbitClient.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbit channel: %s", err)
	}
	if err = rabbitClient.channel.ExchangeDeclarePassive(
		rabbitClient.publishExchange, // name of the exchange
		"topic",                      // type
		true,                         // durable
		false,                        // delete when complete
		false,                        // internal
		false,                        // noWait
		nil,                          // arguments
	); err != nil {
		return nil, fmt.Errorf("rabbit exchange declare: %s", err)
	}
	rabbitClient.subscribeQueue, err = rabbitClient.channel.QueueDeclarePassive(
		subscribeQueue, // queue name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("rabbit queue declare: %s", err)
	}
	rabbitClient.deliveries, err = rabbitClient.channel.Consume(
		rabbitClient.subscribeQueue.Name, // name
		"",                               // consumerTag,
		false,                            // autoAck
		false,                            // exclusive
		false,                            // noLocal
		false,                            // noWait
		nil,                              // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("rabbit queue consume: %s", err)
	}
	rabbitClient.conn.NotifyClose(rabbitClient.connErr)
	return rabbitClient, nil
}

func (c *RabbitClient) ProcessMessages(ctx context.Context, handler DeliveryHandler) error {
	for {
		select {
		case d, ok := <-c.deliveries:
			if !ok {
				return nil
			}
			if err := handler.Handle(ctx, d.Body); err != nil {
				return err
			}
			if err := d.Ack(false); err != nil {
				return err
			}
		case err := <-c.connErr:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *RabbitClient) PublishMessage(ctx context.Context, routingKey string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.channel.PublishWithContext(
		ctx,
		c.publishExchange,
		routingKey,
		true,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Body:            body,
		},
	)
}

func (c *RabbitClient) Shutdown() error {
	if !c.channel.IsClosed() {
		if err := c.channel.Close(); err != nil {
			return err
		}
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
