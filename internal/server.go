//go:build wireinject
// +build wireinject

package server

import (
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/google/wire"
	"github.com/hellofresh/health-go/v5"
	"google.golang.org/grpc"
	"latipe-promotion-services/config"
	"latipe-promotion-services/internal/adapter"
	"latipe-promotion-services/internal/api"
	"latipe-promotion-services/internal/domain/repos"
	"latipe-promotion-services/internal/grpcservice"
	"latipe-promotion-services/internal/grpcservice/interceptor"
	"latipe-promotion-services/internal/grpcservice/vouchergrpc"
	"latipe-promotion-services/internal/middleware"
	"latipe-promotion-services/internal/publisher"
	"latipe-promotion-services/internal/router"
	"latipe-promotion-services/internal/services"
	subscriber "latipe-promotion-services/internal/subs"
	"latipe-promotion-services/internal/subs/createPurchase"
	healthService "latipe-promotion-services/pkgs/healthcheck"
	"latipe-promotion-services/pkgs/mongodb"
	"latipe-promotion-services/pkgs/rabbitclient"
	responses "latipe-promotion-services/pkgs/response"
)

type Server struct {
	app                  *fiber.App
	cfg                  *config.Config
	grpcServ             *grpc.Server
	purchaseSubs         *createPurchase.PurchaseCreateSubscriber
	rollbackPurchaseSubs *createPurchase.PurchaseRollbackSubscriber
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		mongodb.Set,
		publisher.Set,
		rabbitclient.Set,
		grpcservice.Set,
		router.Set,
		repos.Set,
		services.Set,
		adapter.Set,
		api.Set,
		middleware.Set,
		subscriber.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	voucherGrpc vouchergrpc.VoucherServiceServer,
	voucherRouter router.VoucherRouter,
	grpcInterceptor *interceptor.GrpcInterceptor,
	purchaseSubs *createPurchase.PurchaseCreateSubscriber,
	rollbackPurchaseSubs *createPurchase.PurchaseRollbackSubscriber) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
		ErrorHandler: responses.CustomErrorHandler,
	})

	app.Use(logger.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//providing basic authentication for metrics endpoints
	basicAuthConfig := basicauth.Config{
		Users: map[string]string{
			cfg.Metrics.Username: cfg.Metrics.Password,
		},
	}

	// Fiber prometheus
	prometheus := fiberprometheus.New("promotion-services")
	prometheus.RegisterAt(app, cfg.Metrics.Host, basicauth.New(basicAuthConfig))
	app.Use(prometheus.Middleware)
	app.Use(logger.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "Promotion services was developed by TienDat",
			Version: "v1.0.1",
		}
		return ctx.JSON(s)
	})

	// Healthcheck
	h, _ := healthService.NewHealthCheckService(cfg)
	app.Get("/health", basicauth.New(basicAuthConfig), adaptor.HTTPHandlerFunc(h.HandlerFunc))
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/liveness",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			result := h.Measure(c.Context())
			return result.Status == health.StatusOK
		},
		ReadinessEndpoint: "/readiness",
	}))

	//fiber dashboard
	app.Get(cfg.Metrics.FiberDashboard, basicauth.New(basicAuthConfig),
		monitor.New(monitor.Config{Title: "Promotion Services Metrics Page (Fiber)"})) // Healthcheck
	h, _ := healthService.NewHealthCheckService(cfg)
	app.Get("/health", basicauth.New(basicAuthConfig), adaptor.HTTPHandlerFunc(h.HandlerFunc))
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/liveness",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			result := h.Measure(c.Context())
			return result.Status == health.StatusOK
		},
		ReadinessEndpoint: "/readiness",
	}))

	//fiber dashboard
	app.Get(cfg.Metrics.FiberDashboard, basicauth.New(basicAuthConfig),
		monitor.New(monitor.Config{Title: "Promotion Services Metrics Page (Fiber)"}))

	voucherRouter.Init(&v1)

	//init grpc
	grpcServ := grpc.NewServer(grpc.UnaryInterceptor(grpcInterceptor.MiddlewareUnaryRequest))
	vouchergrpc.RegisterVoucherServiceServer(grpcServ, voucherGrpc)

	return &Server{
		cfg:                  cfg,
		app:                  app,
		purchaseSubs:         purchaseSubs,
		rollbackPurchaseSubs: rollbackPurchaseSubs,
		grpcServ:             grpcServ,
	}
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.cfg
}

func (serv Server) GrpcServ() *grpc.Server {
	return serv.grpcServ
}

func (serv Server) CommitPurchaseTransactionSubscriber() *createPurchase.PurchaseCreateSubscriber {
	return serv.purchaseSubs
}

func (serv Server) RollbackPurchaseTransactionSubscriber() *createPurchase.PurchaseRollbackSubscriber {
	return serv.rollbackPurchaseSubs
}
