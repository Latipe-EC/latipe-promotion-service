package router

import (
	"github.com/gofiber/fiber/v2"
	"latipe-promotion-services/internal/api"
	"latipe-promotion-services/internal/middleware"
)

type VoucherRouter interface {
	Init(root *fiber.Router)
}

type voucherRouter struct {
	voucherHandler *api.VoucherHandle
	middleware     *middleware.AuthMiddleware
}

func NewVoucherRouter(voucherHandler *api.VoucherHandle, middleware *middleware.AuthMiddleware) VoucherRouter {
	return voucherRouter{
		voucherHandler: voucherHandler,
		middleware:     middleware,
	}
}

func (o voucherRouter) Init(root *fiber.Router) {
	voucher := (*root).Group("/vouchers")
	admin := voucher.Group("/admin")
	{
		admin.Get("", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.FindAll)
		admin.Post("", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.CreateNewVoucher)
		admin.Get("/:id", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.GetById)
		admin.Get("/code/:code", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.GetByCode)
		admin.Patch("code/:code", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.UpdateStatusVoucher)
	}

	user := voucher.Group("/user", o.middleware.RequiredAuthentication())
	{
		user.Get("/foryou", o.voucherHandler.FindVoucherForUser)
		user.Get("/code/:code", o.voucherHandler.GetByCode)
	}

	store := voucher.Group("/store")
	{
		store.Get("", o.middleware.RequiredStoreAuthentication(), o.voucherHandler.GetAllVoucherForStore)
		store.Post("", o.middleware.RequiredStoreAuthentication(), o.voucherHandler.StoreCreateNewVoucher)
		store.Patch("/code/:code", o.middleware.RequiredStoreAuthentication(), o.voucherHandler.UpdateStatusVoucher)
		store.Get("/code/:code", o.middleware.RequiredStoreAuthentication(), o.voucherHandler.GetByCode)
	}

	voucher.Post("/checking", o.middleware.RequiredAuthentication(), o.voucherHandler.CheckingVoucher)
	voucher.Post("/checkout", o.middleware.RequiredAuthentication(), o.voucherHandler.CheckoutPurchaseVoucher)
	voucher.Get("/coming", o.middleware.RequiredAuthentication(), o.voucherHandler.GetComingVoucher)
}
