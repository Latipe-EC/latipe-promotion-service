package healthService

import (
	"github.com/hellofresh/health-go/v5"
	healthMongoDB "github.com/hellofresh/health-go/v5/checks/mongo"
	healthRabbit "github.com/hellofresh/health-go/v5/checks/rabbitmq"
	"latipe-promotion-services/config"
	"time"
)

func NewHealthCheckService(config *config.Config) (*health.Health, error) {
	// add some checks on instance creation
	h, err := health.New(health.WithComponent(health.Component{
		Name:    "order-service-v2",
		Version: "v2.0",
	}))
	if err != nil {
		return nil, err
	}

	//mysql check
	err = h.Register(health.Config{
		Name:      "mongoDB",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthMongoDB.New(healthMongoDB.Config{
			DSN: config.Mongodb.ConnectionString,
		}),
	})
	if err != nil {
		return nil, err
	}

	//rabbitMQ check
	err = h.Register(health.Config{
		Name:      "rabbitMQ",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthRabbit.New(healthRabbit.Config{
			DSN: config.RabbitMQ.Connection,
		}),
	})
	if err != nil {
		return nil, err
	}

	return h, nil
}
