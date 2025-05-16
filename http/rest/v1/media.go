package v1

import (
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"github.com/co1seam/ember_backend_api_gateway"
	"github.com/co1seam/ember_backend_api_gateway/http/rpc"
	"github.com/gofiber/fiber/v2"
	"io"
)

func (h *Handler) uploadMedia(ctx *fiber.Ctx) error {
	reqCtx := ctx.UserContext()
	mediaContent, err := ctx.FormFile("mediaContent")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	mediaPoster, err := ctx.FormFile("mediaPoster")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	fileContent, err := mediaContent.Open()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	defer fileContent.Close()

	filePoster, err := mediaPoster.Open()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	defer filePoster.Close()

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	poster, err := io.ReadAll(filePoster)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var req ember_backend_api_gateway.SendMediaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := rpc.MediaClient.SendMedia(reqCtx, mediav1.SendMediaRequest{
		Author:   req.Author,
		Type:     req.Type,
		Content:  &ember_backend_api_gateway.MediaFile{Content: content, Filename: mediaContent.Filename, MimeType: mediaContent.Header.Get("Content-Type")},
		Poster:   poster,
		Duration: req.Duration,
		IsActive: req.IsActive,
	})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
