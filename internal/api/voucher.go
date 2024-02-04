package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	dto "latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/middleware"
	"latipe-promotion-services/internal/services/voucherserv"
	"latipe-promotion-services/pkgs/pagable"
	responses "latipe-promotion-services/pkgs/response"
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
		if errors.Is(err, responses.ErrInvalidVoucherData) {
			return responses.ErrInvalidVoucherData
		}
		return responses.ErrInternalServer
	}

	resp := responses.DefaultSuccess
	resp.Message = "the voucher was created"
	resp.Data = dataResp

	return resp.JSON(ctx)
}

func (api VoucherHandle) StoreCreateNewVoucher(ctx *fiber.Ctx) error {
	var request dto.CreateVoucherRequest

	storeId := fmt.Sprintf("%s", ctx.Locals(middleware.STORE_ID))
	if storeId == "" {
		return responses.ErrUnauthenticated
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

	request.StoreID = storeId
	dataResp, err := api.service.CreateVoucher(ctx.Context(), &request)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return responses.ErrExistVoucherCode
		}
		if errors.Is(err, responses.ErrInvalidVoucherData) {
			return responses.ErrInvalidVoucherData
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

func (api VoucherHandle) ApplyVoucher(ctx *fiber.Ctx) error {
	var request dto.ApplyVoucherRequest

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

	err := api.service.UsingVoucherToOrder(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch {
		case responses.Is(err, mongo.ErrNoDocuments):
			return responses.ErrNotFoundRecord
		}

		return err
	}
	resp := responses.DefaultSuccess
	return resp.JSON(ctx)
}

func (api VoucherHandle) CheckingVoucher(ctx *fiber.Ctx) error {
	var request dto.PurchaseVoucherRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	userId := fmt.Sprintf("%s", ctx.Locals("user_id"))
	if userId == "" {
		return responses.ErrUnauthenticated
	}

	request.UserId = userId
	if err := valid.GetValidator().Validate(&request); err != nil {
		log.Errorf("%v", err)
		resp := responses.DefaultError
		resp.Message = err.Error()
		resp.Code = http.StatusBadRequest

		return resp.JSON(ctx)
	}

	dataResp, err := api.service.CheckoutVoucher(ctx.Context(), &request)
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
	var request dto.RollbackVoucherRequest

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
	search := ctx.Query("search")

	dataResp, err := api.service.GetAllVoucher(context, search, request.Query)
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

	userId := fmt.Sprintf("%s", ctx.Locals("user_id"))
	if userId == "" {
		return responses.ErrUnauthenticated
	}

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return responses.ErrInvalidParameters
	}
	request := dto.GetVoucherListRequest{
		Query: query,
	}

	voucherCode := ctx.Query("search")

	dataResp, err := api.service.GetUserVoucher(context, voucherCode, userId, request.Query)
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

func (api VoucherHandle) GetAllVoucherForStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	storeId := fmt.Sprintf("%s", ctx.Locals(middleware.STORE_ID))
	if storeId == "" {
		return responses.ErrUnauthenticated
	}

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return responses.ErrInvalidParameters
	}
	request := dto.GetVoucherListRequest{
		Query: query,
	}

	voucherCode := ctx.Query("search")

	dataResp, err := api.service.GetAllVoucherOfStore(context, voucherCode, storeId, request.Query)
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

func (api VoucherHandle) StoreCancelStatusVoucher(ctx *fiber.Ctx) error {
	var request dto.CancelUpdateVoucherRequest

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

	err := api.service.StoreCancelVoucher(ctx.Context(), &request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return responses.ErrNotFound
		}

	}

	resp := responses.DefaultSuccess
	resp.Message = "the voucher was updated"

	return resp.JSON(ctx)
}
