package purchaseCreate

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-promotion-services/config"
	"latipe-promotion-services/internal/domain/message"
	"latipe-promotion-services/pkgs/rabbitclient"

	"time"
)

type ReplyPurchaseTransactionPub struct {
	channel *amqp.Channel
	cfg     *config.Config
}

func NewReplyPurchaseTransactionPub(cfg *config.Config, conn *amqp.Connection) *ReplyPurchaseTransactionPub {
	producer := ReplyPurchaseTransactionPub{
		cfg: cfg,
	}

	ch, err := conn.Channel()
	if err != nil {
		rabbitclient.FailOnError(err, "Failed to open a channel")
		return nil
	}
	producer.channel = ch

	return &producer
}

func (pub *ReplyPurchaseTransactionPub) ReplyPurchaseMessage(message *message.ReplyPurchaseMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := rabbitclient.ParseOrderToByte(&message)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.CreatePurchaseEvent.Exchange,
		pub.cfg.RabbitMQ.CreatePurchaseEvent.ReplyRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.CreatePurchaseEvent.Exchange,
		pub.cfg.RabbitMQ.CreatePurchaseEvent.ReplyRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	rabbitclient.FailOnError(err, "Failed to publish a message")

	return nil
}
