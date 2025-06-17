package v1

import (
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	api_gateway "github.com/co1seam/ember_backend_api_gateway"
	"github.com/co1seam/ember_backend_api_gateway/http/rpc"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *Handler) createMedia(ctx *fiber.Ctx) error {
	var req api_gateway.CreateMediaRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	sub, ok := GetSubject(ctx)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "subject not found")
	}

	resp, err := rpc.MediaClient.CreateMedia(ctx.Context(), &mediav1.CreateMediaRequest{
		Title:       req.Title,
		Description: req.Description,
		ContentType: req.ContentType,
		OwnerId:     sub,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp.Media)
}

func (h *Handler) getMedia(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")
	if ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}

	sub, ok := GetSubject(ctx)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "subject not found")
	}

	resp, err := rpc.MediaClient.GetMedia(ctx.Context(), &mediav1.GetMediaRequest{
		Id: ID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if resp.Media.OwnerId != sub {
		return fiber.NewError(fiber.StatusForbidden, "not allowed")
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.Media)
}

func (h *Handler) updateMedia(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "media ID is required")
	}

	current, err := rpc.MediaClient.GetMedia(ctx.Context(), &mediav1.GetMediaRequest{Id: id})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userID, ok := GetSubject(ctx)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
	}

	if current.Media.OwnerId != userID {
		return fiber.NewError(fiber.StatusForbidden, "access denied")
	}

	var req api_gateway.UpdateMediaRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := rpc.MediaClient.UpdateMedia(ctx.Context(), &mediav1.UpdateMediaRequest{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.Media)
}

func (h *Handler) deleteMedia(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "media ID is required")
	}

	current, err := rpc.MediaClient.GetMedia(ctx.Context(), &mediav1.GetMediaRequest{Id: id})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userID, ok := GetSubject(ctx)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
	}

	if current.Media.OwnerId != userID {
		return fiber.NewError(fiber.StatusForbidden, "access denied")
	}

	_, err = rpc.MediaClient.DeleteMedia(ctx.Context(), &mediav1.DeleteMediaRequest{
		Id: id,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) listMedia(ctx *fiber.Ctx) error {
	limitStr := ctx.Query("limit", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	sub, ok := GetSubject(ctx)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "subject not found")
	}

	resp, err := rpc.MediaClient.ListMedia(ctx.Context(), &mediav1.ListMediaRequest{
		OwnerId: sub,
		Limit:   int32(limit),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.Media)
}
