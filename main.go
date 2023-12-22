package main

import (
	"context"
	"fmt"
	"github.com/ansrivas/fiberprometheus/v2"
	"go.mongodb.org/mongo-driver/event"
	"latipe-promotion-services/adapter/userserv"
	handler "latipe-promotion-services/api"
	"latipe-promotion-services/domain/repos"
	"latipe-promotion-services/middleware"
	responses "latipe-promotion-services/response"
	"latipe-promotion-services/service/voucherserv"

	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func main() {
	//read env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, e *event.CommandStartedEvent) {
			fmt.Println(e.Command)
		},
		Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			fmt.Println(failedEvent.Failure)
		},
	}

	opts := options.Client().SetMonitor(monitor)

	//create connect to mongo
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri), opts)
	db := client.Database("latipe_promotion_db")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//create instance fiber
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
		ErrorHandler: responses.CustomErrorHandler,
	})
	// cai thang nay` wrapper dong code kia roi
	prometheus := fiberprometheus.New("promotion-service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	app.Use(logger.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "Promotion service was developed by TienDat",
			Version: "v0.0.1",
		}
		return ctx.JSON(s)
	})

	//create instance resty-go
	cli := resty.New().
		SetTimeout(5 * time.Second)

	//repository

	voucherRepos := repos.NewVoucherRepos(db)

	//service
	userServ := userserv.NewUserService(cli)
	voucherServ := voucherserv.NewVoucherService(&voucherRepos)

	//api handler
	voucherApi := handler.NewVoucherHandler(&voucherServ)

	//middleware
	authMiddleware := middleware.NewAuthMiddleware(&userServ)

	api := app.Group("/api")

	v1 := api.Group("/v1")

	voucher := v1.Group("/vouchers")
	voucher.Post("", authMiddleware.RequiredRoles([]string{"ADMIN"}), voucherApi.CreateNewVoucher)
	voucher.Get("", authMiddleware.RequiredRoles([]string{"ADMIN"}), voucherApi.FindAll)
	voucher.Get("/user/foryou", voucherApi.FindVoucherForUser)
	voucher.Get("/:id", authMiddleware.RequiredRoles([]string{"ADMIN"}), voucherApi.GetById)
	voucher.Patch("code/:code", authMiddleware.RequiredRoles([]string{"ADMIN"}), voucherApi.UpdateStatusVoucher)
	voucher.Get("/code/:code", authMiddleware.RequiredRoles([]string{"ADMIN"}), voucherApi.GetByCode)
	voucher.Post("/apply", authMiddleware.RequiredAuthentication(), voucherApi.UseVoucher)
	voucher.Post("/rollback", authMiddleware.RequiredAuthentication(), voucherApi.UseVoucher)
	voucher.Post("/checking", authMiddleware.RequiredAuthentication(), voucherApi.CheckingVoucher)
	err = app.Listen(":5010")
	if err != nil {
		return
	}
}
