package middleware

import (
	"github.com/gofiber/fiber/v2"
	"latipe-promotion-services/internal/adapter/storeserv"
	storeDTO "latipe-promotion-services/internal/adapter/storeserv/dto"
	"latipe-promotion-services/internal/adapter/userserv"
	"latipe-promotion-services/internal/adapter/userserv/dto"
	"latipe-promotion-services/pkgs/response"
	"strings"
)

type AuthMiddleware struct {
	userServ  *userserv.UserService
	storeServ storeserv.StoreService
}

func NewAuthMiddleware(service *userserv.UserService, storeServ storeserv.StoreService) *AuthMiddleware {
	return &AuthMiddleware{userServ: service, storeServ: storeServ}
}

func (auth AuthMiddleware) RequiredRoles(roles []string, option ...int) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" || len(strings.Split(bearToken, " ")) < 2 {
			return responses.ErrUnauthenticated
		}

		str := strings.Split(bearToken, " ")
		if len(str) < 2 {
			return responses.ErrUnauthenticated
		}

		bearToken = str[1]
		req := dto.AuthorizationRequest{}
		req.Token = bearToken

		resp, err := auth.userServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return responses.ErrInternalServer
		}

		for _, i := range roles {
			if i == resp.Role {
				return ctx.Next()
			}
		}
		return responses.ErrPermissionDenied
	}
}

func (auth AuthMiddleware) RequiredAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return responses.ErrUnauthenticated
		}

		str := strings.Split(bearToken, " ")
		if len(str) < 2 {
			return responses.ErrUnauthenticated
		}

		bearToken = str[1]
		req := dto.AuthorizationRequest{
			Token: bearToken,
		}
		resp, err := auth.userServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return responses.ErrInternalServer
		}

		ctx.Locals("user_name", resp.Email)
		ctx.Locals("user_id", resp.Id)
		ctx.Locals("bearer_token", bearToken)
		return ctx.Next()
	}
}

func (a AuthMiddleware) RequiredStoreAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return responses.ErrUnauthenticated
		}

		str := strings.Split(bearToken, " ")
		if len(str) < 2 {
			return responses.ErrUnauthenticated
		}

		bearToken = str[1]
		req := dto.AuthorizationRequest{
			Token: bearToken,
		}

		resp, err := a.userServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return err
		}

		//validate store
		storeReq := storeDTO.GetStoreIdByUserRequest{UserID: resp.Id}
		storeReq.BaseHeader.BearToken = bearToken

		storeResp, err := a.storeServ.GetStoreByUserId(ctx.Context(), &storeReq)
		if err != nil {
			return err
		}

		if storeResp.StoreID == "" {
			return responses.ErrPermissionDenied
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		ctx.Locals(STORE_ID, storeResp.StoreID)

		return ctx.Next()
	}
}
