package v1

import (
	authv1 "github.com/co1seam/ember-backend-api-contracts/gen/go/auth"
	"github.com/co1seam/ember_backend_api_gateway/http/rpc"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) AuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Cookies("access_token")
		if token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token required"})
		}

		sub, err := rpc.AuthClient.ValidateToken(ctx.Context(), &authv1.ValidateTokenRequest{AccessToken: token})
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		ctx.Locals("sub", sub.Subject)

		return ctx.Next()
	}
}

func GetSubject(ctx *fiber.Ctx) (string, bool) {
	sub := ctx.Locals("sub")
	if sub == nil {
		return "", false
	}

	return sub.(string), true
}
