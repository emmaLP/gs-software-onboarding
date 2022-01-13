package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type client struct {
	logger      *zap.Logger
	amqpConn    *amqp.Connection
	queue       amqp.Queue
	amqpChannel *amqp.Channel
}

type Client interface {
	SendMessage(item commonModel.Item) error
	ReceiveMessage(msgChan chan *commonModel.Item) error
	CloseConnection()
}

func New(logger *zap.Logger, ctx context.Context, amqpConfig *model.RabbitMqConfig) (*client, error) {
	amqpConn, err := connect(ctx, logger, amqpConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to rabbitmq. %w", err)
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Unable to create channel. %w", err)
	}

	queue, err := channel.QueueDeclare(
		amqpConfig.QueueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare queue. %w", err)
	}

	return &client{
		logger:      logger,
		amqpConn:    amqpConn,
		amqpChannel: channel,
		queue:       queue,
	}, nil
}

func connect(ctx context.Context, logger *zap.Logger, amqpConfig *model.RabbitMqConfig) (*amqp.Connection, error) {
	amqpUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		url.QueryEscape(amqpConfig.Username),
		url.QueryEscape(amqpConfig.Password),
		url.QueryEscape(amqpConfig.Host),
		url.QueryEscape(amqpConfig.Port))
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", ctx.Err())
		default:
			conn, err := amqp.Dial(amqpUrl)
			if err != nil {
				logger.Error("Unable to establish connection to RabbitMQ", zap.Error(err))
				break
			}
			return conn, nil
		}
	}
}

func (c *client) SendMessage(item commonModel.Item) error {
	body, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("Failed to marshel item to json: %w", err)
	}
	err = c.amqpChannel.Publish(
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err == nil {
		c.logger.Info("Item successfully pushed to queue")
	}
	return err
}

func (c *client) ReceiveMessage(msgChan chan *commonModel.Item) error {
	messages, err := c.amqpChannel.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Consuming messages: %w.", err)
	}
	for message := range messages {
		item := commonModel.Item{}
		err := json.Unmarshal(message.Body, &item)
		if err != nil {
			c.logger.Error("Unable to unmarshal message body", zap.ByteString("message_body", message.Body), zap.Error(err))
			continue
		}
		if err := message.Ack(false); err != nil {
			c.logger.Error("Unable to acknowledge message", zap.Error(err))
			continue
		}
		msgChan <- &item
	}

	return nil
}

func (c *client) CloseConnection() {
	if err := c.amqpChannel.Close(); err != nil {
		c.logger.Error("Unable to close channel", zap.Error(err))
	}
	if err := c.amqpConn.Close(); err != nil {
		c.logger.Error("Unable to close rabbitmq connection", zap.Error(err))
	}
}
