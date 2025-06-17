package v1

import (
	authv1 "github.com/co1seam/ember-backend-api-contracts/gen/go/auth"
	models "github.com/co1seam/ember_backend_api_gateway"
	"github.com/co1seam/ember_backend_api_gateway/http/rpc"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (h *Handler) sendOtp(ctx *fiber.Ctx) error {
	reqCtx := ctx.Context()

	var req struct {
		Email string `json:"email"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	resp, err := rpc.AuthClient.SendOTP(reqCtx, &authv1.SendOTPRequest{
		Email: req.Email,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if !resp.Success {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid OTP"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": resp.Success})
}

func (h *Handler) verifyOtp(ctx *fiber.Ctx) error {
	reqCtx := ctx.Context()

	var req struct {
		Otp string `json:"otp"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	resp, err := rpc.AuthClient.VerifyOTP(reqCtx, &authv1.VerifyOTPRequest{
		Otp: req.Otp,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"email": resp.Email})
}

func (h *Handler) signUp(ctx *fiber.Ctx) error {
	reqCtx := ctx.Context()

	var req models.SignUpRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := rpc.AuthClient.SignUp(reqCtx, &authv1.SignUpRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	h.createCookies(ctx, resp.RefreshToken, resp.AccessToken)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *Handler) signIn(ctx *fiber.Ctx) error {
	reqCtx := ctx.Context()
	var req models.SignInRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := rpc.AuthClient.SignIn(reqCtx, &authv1.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	h.createCookies(ctx, resp.RefreshToken, resp.AccessToken)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *Handler) signOut(ctx *fiber.Ctx) error {

	h.deleteCookies(ctx)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *Handler) refresh(ctx *fiber.Ctx) error {
	reqCtx := ctx.Context()
	var req models.RefreshTokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := rpc.AuthClient.RefreshToken(reqCtx, &authv1.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	h.createCookies(ctx, resp.RefreshToken, resp.AccessToken)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *Handler) createCookies(ctx *fiber.Ctx, refreshToken, accessToken string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
	})
}

func (h *Handler) deleteCookies(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})
}
