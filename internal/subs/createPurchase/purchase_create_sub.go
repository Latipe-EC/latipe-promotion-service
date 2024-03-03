package createPurchase

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-promotion-services/config"
	"latipe-promotion-services/internal/domain/message"
	"latipe-promotion-services/internal/services/voucherserv"
	"sync"
	"time"
)

type PurchaseCreateSubscriber struct {
	config      *config.Config
	voucherServ *voucherserv.VoucherService
	conn        *amqp.Connection
}

func NewPurchaseCreateSubscriber(cfg *config.Config,
	voucherServ *voucherserv.VoucherService, conn *amqp.Connection) *PurchaseCreateSubscriber {
	return &PurchaseCreateSubscriber{
		config:      cfg,
		voucherServ: voucherServ,
		conn:        conn,
	}
}

func (orch PurchaseCreateSubscriber) ListenProductPurchaseCreate(wg *sync.WaitGroup) {
	channel, err := orch.conn.Channel()
	defer channel.Close()

	// define an exchange type "topic"
	err = channel.ExchangeDeclare(
		orch.config.RabbitMQ.CreatePurchaseEvent.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare exchange: %v", err)
	}

	// create queue
	q, err := channel.QueueDeclare(
		"purchase_promotion_commit",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare queue: %v", err)
	}

	err = channel.QueueBind(
		q.Name,
		orch.config.RabbitMQ.CreatePurchaseEvent.CommitRoutingKey,
		orch.config.RabbitMQ.CreatePurchaseEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                           // queue
		orch.config.RabbitMQ.ServiceName, // consumer
		true,                             // auto ack
		false,                            // exclusive
		false,                            // no local
		false,                            // no wait
		nil,                              //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := orch.handleMessage(&msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}
	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")
}

func (orch PurchaseCreateSubscriber) handleMessage(msg *amqp.Delivery) error {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageDTO := message.CreatePurchaseMessage{}

	if err := sonic.Unmarshal(msg.Body, &messageDTO); err != nil {
		log.Infof("Parse message to order failed cause: %s", err)
		return err
	}

	err := orch.voucherServ.CommitVoucherTransaction(ctx, &messageDTO)
	if err != nil {
		log.Infof("Handling message was failed cause: %s", err)
		return err
	}

	endTime := time.Now()
	log.Infof("The message [%v]  was processed successfully - duration:%v", messageDTO.CheckoutID, endTime.Sub(startTime))
	return nil
}
