package middleware

import (
	"github.com/gofiber/fiber/v2"
	"latipe-promotion-services/adapter/userserv"
	"latipe-promotion-services/adapter/userserv/dto"
	responses "latipe-promotion-services/response"
	"strings"
)

type AuthMiddleware struct {
	userServ *userserv.UserService
}

func NewAuthMiddleware(service *userserv.UserService) *AuthMiddleware {
	return &AuthMiddleware{userServ: service}
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
