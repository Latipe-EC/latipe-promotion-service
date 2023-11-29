package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"latipe-promotion-services/domain/dto"
	"latipe-promotion-services/pkgs/pagable"
	"latipe-promotion-services/pkgs/valid"
	"latipe-promotion-services/service/voucherserv"
	"net/http"
)

type VoucherHandle struct {
	service *voucherserv.VoucherService
}

func NewVoucherHandler(service *voucherserv.VoucherService) VoucherHandle {
	return VoucherHandle{
		service: service,
	}
}

func (api VoucherHandle) CreateNewVoucher(ctx *fiber.Ctx) error {
	var request dto.CreateVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusBadRequest).SendString("Parse body was failed")
	}

	if err := valid.GetValidator().Validate(request); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	dataResp, err := api.service.CreateVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusBadRequest).SendString("Invalid params")
	}

	return ctx.JSON(dataResp)
}

func (api VoucherHandle) GetById(ctx *fiber.Ctx) error {
	var request dto.GetVoucherByIdRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusBadRequest).SendString("Parse id was failed")
	}

	dataResp, err := api.service.GetVoucherById(ctx.Context(), request.Id)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return ctx.Status(http.StatusNotFound).SendString("Not found")
		default:
			return ctx.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}

	if dataResp == nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusNotFound).SendString("The voucher was not found")
	}

	return ctx.JSON(dataResp)
}

func (api VoucherHandle) GetByCode(ctx *fiber.Ctx) error {
	var request dto.GetVoucherByCodeRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusBadRequest).SendString("Parse code was failed")
	}

	dataResp, err := api.service.GetVoucherByCode(ctx.Context(), request.Code)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return ctx.Status(http.StatusNotFound).SendString("Not found")
		default:
			return ctx.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}

	if dataResp == nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusNotFound).SendString("The voucher was not found")
	}

	return ctx.JSON(dataResp)
}

func (api VoucherHandle) UseVoucher(ctx *fiber.Ctx) error {
	var request dto.UseVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusBadRequest).SendString("Parse code was failed")
	}

	if err := valid.GetValidator().Validate(request); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	dataResp, err := api.service.UseVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return ctx.Status(http.StatusNotFound).SendString("Not found")
		default:
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	return ctx.JSON(dataResp)
}

func (api VoucherHandle) RollBack(ctx *fiber.Ctx) error {
	var request dto.UseVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusBadRequest).SendString("Parse code was failed")
	}

	if err := valid.GetValidator().Validate(request); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err := api.service.RollBackVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return ctx.Status(http.StatusNotFound).SendString("Not found")
		default:
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	dataResp := make(map[string]interface{})
	dataResp["status"] = "success"

	return ctx.JSON(dataResp)
}

func (api VoucherHandle) FindAll(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid query")
	}
	request := dto.GetVoucherListRequest{
		Query: query,
	}

	dataResp, err := api.service.GetAllVoucher(context, request.Query)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return ctx.Status(http.StatusNotFound).SendString("Not found")
		default:
			return ctx.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}

	if dataResp == nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusNotFound).SendString("The voucher was not found")
	}

	return ctx.JSON(dataResp)
}

func (api VoucherHandle) FindVoucherForUser(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid query")
	}
	request := dto.GetVoucherListRequest{
		Query: query,
	}

	dataResp, err := api.service.GetUserVoucher(context, request.Query)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return ctx.Status(http.StatusNotFound).SendString("Not found")
		default:
			return ctx.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}

	if dataResp == nil {
		log.Errorf("%v", err)
		return ctx.Status(http.StatusNotFound).SendString("The voucher was not found")
	}

	return ctx.JSON(dataResp)
}
