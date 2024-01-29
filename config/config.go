package config

import (
	"errors"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var Set = wire.NewSet(NewConfig)

type Config struct {
	Server         Server
	AdapterService AdapterService
	RabbitMQ       RabbitMQ
	Mongodb        Mongodb
}

type Server struct {
	Name                string
	ApiHeaderKey        string
	AppVersion          string
	RestPort            string
	GrpcPort            string
	BaseURI             string
	Mode                string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	SSL                 bool
	CtxDefaultTimeout   int
	CSRF                bool
	Debug               bool
	MaxCountRequest     int           // max count of connections
	ExpirationLimitTime time.Duration //  expiration time of the limit
}

type Mongodb struct {
	ConnectionString string
	Address          string
	Username         string
	Password         string
	DbName           string
	ConnectTimeout   time.Duration
	MaxConnIdleTime  int
	MinPoolSize      uint64
	MaxPoolSize      uint64
}

type RabbitMQ struct {
	Connection          string
	EmailEvent          EmailEvent
	CreatePurchaseEvent PurchaseEvent
	ServiceName         string
}

type PurchaseEvent struct {
	Exchange           string
	CommitRoutingKey   string
	RollbackRoutingKey string
	ReplyRoutingKey    string
}

type EmailEvent struct {
	Connection string
	Exchange   string
	RoutingKey string
	Queue      string
}

type AdapterService struct {
	UserService  UserService
	EmailService EmailService
}

type UserService struct {
	AuthURL     string
	UserURL     string
	InternalKey string
}

type ProductService struct {
	BaseURL     string
	InternalKey string
}

type EmailService struct {
	Email string
	Host  string
	Key   string
}

// Get config path for local or docker
func getDefaultConfig() string {
	return "./config/config"
}

// Load config file from given path
func NewConfig() (*Config, error) {
	config := Config{}
	path := os.Getenv("cfgPath")
	if path == "" {
		path = getDefaultConfig()
	}

	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &config, nil
}
