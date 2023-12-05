package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"latipe-promotion-services/domain/dto"
	"latipe-promotion-services/pkgs/pagable"
	"latipe-promotion-services/pkgs/valid"
	responses "latipe-promotion-services/response"
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
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.CreateVoucher(ctx.Context(), &request)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return responses.ErrExistVoucherCode
		}
		return responses.ErrInternalServer
	}

	resp := responses.DefaultSuccess
	resp.Message = "the voucher was created"
	resp.Data = dataResp

	return resp.JSON(ctx)
}

func (api VoucherHandle) UpdateStatusVoucher(ctx *fiber.Ctx) error {
	var request dto.UpdateVoucherRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	err := api.service.UpdateVoucherStatus(ctx.Context(), &request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return responses.ErrNotFound
		}

	}

	resp := responses.DefaultSuccess
	resp.Message = "the voucher was updated"

	return resp.JSON(ctx)
}

func (api VoucherHandle) GetById(ctx *fiber.Ctx) error {
	var request dto.GetVoucherByIdRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.GetVoucherById(ctx.Context(), request.Id)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		default:
			return responses.ErrInternalServer
		}
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) GetByCode(ctx *fiber.Ctx) error {
	var request dto.GetVoucherByCodeRequest

	if err := ctx.ParamsParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.GetVoucherByCode(ctx.Context(), request.Code)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		default:
			return responses.ErrInternalServer
		}
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) UseVoucher(ctx *fiber.Ctx) error {
	var request dto.UseVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.UseVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch {
		case responses.Is(err, mongo.ErrNoDocuments):
			return responses.ErrNotFoundRecord
		}

		return err
	}
	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) CheckingVoucher(ctx *fiber.Ctx) error {
	var request dto.CheckingVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.CheckingVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		}
		return err
	}
	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) RollBack(ctx *fiber.Ctx) error {
	var request dto.UseVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	if err := valid.GetValidator().Validate(request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	err := api.service.RollBackVoucher(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		default:
			return responses.ErrInternalServer
		}
	}
	dataResp := make(map[string]interface{})
	dataResp["status"] = "success"

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) FindAll(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return responses.ErrInvalidParameters
	}
	request := dto.GetVoucherListRequest{
		Query: query,
	}

	dataResp, err := api.service.GetAllVoucher(context, request.Query)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		default:
			return responses.ErrInternalServer
		}
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (api VoucherHandle) FindVoucherForUser(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return responses.ErrInvalidParameters
	}
	request := dto.GetVoucherListRequest{
		Query: query,
	}

	dataResp, err := api.service.GetUserVoucher(context, request.Query)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		default:
			return responses.ErrInternalServer
		}
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}
