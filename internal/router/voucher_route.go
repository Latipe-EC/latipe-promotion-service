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
	voucher.Post("", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.CreateNewVoucher)
	voucher.Get("", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.FindAll)
	voucher.Get("/user/foryou", o.voucherHandler.FindVoucherForUser)
	voucher.Get("/:id", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.GetById)
	voucher.Patch("code/:code", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.UpdateStatusVoucher)
	voucher.Get("/code/:code", o.middleware.RequiredRoles([]string{"ADMIN"}), o.voucherHandler.GetByCode)
	voucher.Post("/apply", o.middleware.RequiredAuthentication(), o.voucherHandler.UseVoucher)
	voucher.Post("/rollback", o.middleware.RequiredAuthentication(), o.voucherHandler.UseVoucher)
	voucher.Post("/checking", o.middleware.RequiredAuthentication(), o.voucherHandler.CheckingVoucher)
}
