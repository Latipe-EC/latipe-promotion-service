package rabbitclient

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-promotion-services/config"
)

func NewRabbitClientConnection(globalCfg *config.Config) *amqp.Connection {
	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": globalCfg.RabbitMQ.ServiceName,
		},
	}

	conn, err := amqp.DialConfig(globalCfg.RabbitMQ.Connection, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ cause:%v", err)
	}

	log.Info("Comsumer has been connected")
	return conn
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ParseOrderToByte(request interface{}) ([]byte, error) {
	jsonObj, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}
	return jsonObj, err
}
