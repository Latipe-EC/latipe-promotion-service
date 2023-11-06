package middleware

import (
	"github.com/gofiber/fiber/v2"
	"latipe-promotion-services/adapter/userserv"
	"latipe-promotion-services/adapter/userserv/dto"
	"net/http"
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
			return ctx.Status(http.StatusUnauthorized).SendString("Unauthenticated")
		}

		str := strings.Split(bearToken, " ")
		if len(str) < 2 {
			return ctx.Status(http.StatusUnauthorized).SendString("Unauthenticated")
		}

		bearToken = str[1]
		req := dto.AuthorizationRequest{}
		req.Token = bearToken

		resp, err := auth.userServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}

		for _, i := range roles {
			if i == resp.Role {
				return ctx.Next()
			}
		}
		return ctx.Status(http.StatusForbidden).SendString("Permission Denied")
	}
}

func (auth AuthMiddleware) RequiredAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return ctx.Status(http.StatusUnauthorized).SendString("Unauthenticated")
		}

		str := strings.Split(bearToken, " ")
		if len(str) < 2 {
			return ctx.Status(http.StatusUnauthorized).SendString("Unauthenticated")
		}

		bearToken = str[1]
		req := dto.AuthorizationRequest{
			Token: bearToken,
		}
		resp, err := auth.userServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}

		ctx.Locals("user_name", resp.Email)
		ctx.Locals("user_id", resp.Id)
		ctx.Locals("bearer_token", bearToken)
		return ctx.Next()
	}
}
