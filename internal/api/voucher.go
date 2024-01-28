package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	dto2 "latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/services/voucherserv"
	"latipe-promotion-services/pkgs/pagable"
	responses2 "latipe-promotion-services/pkgs/response"
	"latipe-promotion-services/pkgs/valid"
	"net/http"
)

type VoucherHandle struct {
	service *voucherserv.VoucherService
}

func NewVoucherHandler(service *voucherserv.VoucherService) *VoucherHandle {
	return &VoucherHandle{
		service: service,
	}
}

func (api VoucherHandle) CreateNewVoucher(ctx *fiber.Ctx) error {
	var request dto2.CreateVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.CreateVoucher(ctx.Context(), &request)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return responses2.ErrExistVoucherCode
		}
		return responses2.ErrInternalServer
	}

	resp := responses2.DefaultSuccess
	resp.Message = "the voucher was created"
	resp.Data = dataResp

	return resp.JSON(ctx)
}

func (api VoucherHandle) UpdateStatusVoucher(ctx *fiber.Ctx) error {
	var request dto2.UpdateVoucherRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	err := api.service.UpdateVoucherStatus(ctx.Context(), &request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return responses2.ErrNotFound
		}

	}

	resp := responses2.DefaultSuccess
	resp.Message = "the voucher was updated"

	return resp.JSON(ctx)
}

func (api VoucherHandle) GetById(ctx *fiber.Ctx) error {
	var request dto2.GetVoucherByIdRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.GetVoucherById(ctx.Context(), request.Id)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses2.ErrNotFoundRecord
		default:
			return responses2.ErrInternalServer
		}
	}

	resp := responses2.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) GetByCode(ctx *fiber.Ctx) error {
	var request dto2.GetVoucherByCodeRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.GetVoucherByCode(ctx.Context(), request.Code)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses2.ErrNotFoundRecord
		default:
			return responses2.ErrInternalServer
		}
	}

	resp := responses2.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) UseVoucher(ctx *fiber.Ctx) error {
	var request dto2.UseVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.UseVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch {
		case responses2.Is(err, mongo.ErrNoDocuments):
			return responses2.ErrNotFoundRecord
		}

		return err
	}
	resp := responses2.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) CheckingVoucher(ctx *fiber.Ctx) error {
	var request dto2.CheckingVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.CheckingVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses2.ErrNotFoundRecord
		}
		return err
	}
	resp := responses2.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) RollBack(ctx *fiber.Ctx) error {
	var request dto2.UseVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(request); err != nil {
		log.Errorf("%v", err)
		resp := responses2.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	err := api.service.RollBackVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses2.ErrNotFoundRecord
		default:
			return responses2.ErrInternalServer
		}
	}
	dataResp := make(map[string]interface{})
	dataResp["status"] = "success"

	resp := responses2.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) FindAll(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return responses2.ErrInvalidParameters
	}
	request := dto2.GetVoucherListRequest{
		Query: query,
	}
	search := ctx.Query("search")

	dataResp, err := api.service.GetAllVoucher(context, search, request.Query)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses2.ErrNotFoundRecord
		default:
			return responses2.ErrInternalServer
		}
	}

	resp := responses2.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) FindVoucherForUser(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return responses2.ErrInvalidParameters
	}
	request := dto2.GetVoucherListRequest{
		Query: query,
	}

	voucherCode := ctx.Query("search")

	dataResp, err := api.service.GetUserVoucher(context, voucherCode, request.Query)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses2.ErrNotFoundRecord
		default:
			return responses2.ErrInternalServer
		}
	}

	resp := responses2.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}
