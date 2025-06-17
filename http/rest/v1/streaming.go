package v1

import (
	"fmt"
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"github.com/co1seam/ember_backend_api_gateway/http/rpc"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"strconv"
	"strings"
)

func (h *Handler) uploadFile(ctx *fiber.Ctx) error {
	FileID := ctx.Query("file_id")
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to open file")
	}
	defer file.Close()

	stream, err := rpc.MediaClient.UploadFile(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	log.Println(FileID)

	if err := stream.Send(&mediav1.FileChunk{
		FileId:    FileID,
		FileName:  fileHeader.Filename,
		TotalSize: fileHeader.Size,
		IsFirst:   true,
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	buf := make([]byte, 64*1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fiber.NewError(fiber.StatusInternalServerError, "read error: "+err.Error())
		}

		if err := stream.Send(&mediav1.FileChunk{
			Content: buf[:n],
			FileId:  FileID,
		}); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "send error: "+err.Error())
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "close error: "+err.Error())
	}

	return ctx.JSON(fiber.Map{
		"file_id": resp.FileId,
		"url":     resp.Url,
	})
}

func (h *Handler) downloadFile(ctx *fiber.Ctx) error {
	fileID := ctx.Query("file_id")
	if fileID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "file ID required")
	}

	userID, ok := GetSubject(ctx)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
	}

	rangeHeader := ctx.Get("Range")

	media, err := rpc.MediaClient.GetMedia(ctx.Context(), &mediav1.GetMediaRequest{
		Id: fileID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if media.Media.OwnerId != userID {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	start, end := parseRangeHeader(rangeHeader, media.Media.Size)

	// Создаем FileRequest с учетом диапазона
	req := &mediav1.FileRequest{
		FileId:  media.Media.Id,
		OwnerId: media.Media.OwnerId,
		Start:   start,
		End:     end,
	}

	stream, err := rpc.MediaClient.DownloadFile(ctx.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Set("Content-Type", media.Media.ContentType)
	ctx.Set("Accept-Ranges", "bytes")
	ctx.Set("Cache-Control", "public, max-age=31536000")
	ctx.Set("Access-Control-Allow-Origin", "http://localhost:8080")
	ctx.Set("Access-Control-Expose-Headers", "Content-Range, Content-Length")

	// Устанавливаем статус и заголовки для частичного контента
	if rangeHeader != "" {
		actualEnd := end
		if end < 0 || end >= media.Media.Size {
			actualEnd = media.Media.Size - 1
		}
		contentLength := actualEnd - start + 1

		ctx.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, actualEnd, media.Media.Size))
		ctx.Set("Content-Length", strconv.FormatInt(contentLength, 10))
		ctx.Status(fiber.StatusPartialContent)
	} else {
		ctx.Set("Content-Length", strconv.FormatInt(media.Media.Size, 10))
	}

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if _, err := ctx.Write(chunk.Chunk); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func parseRangeHeader(header string, fileSize int64) (start, end int64) {
	if header == "" {
		return 0, -1 // -1 означает до конца файла
	}

	parts := strings.SplitN(header, "=", 2)
	if len(parts) != 2 || parts[0] != "bytes" {
		return 0, -1
	}

	rangeParts := strings.SplitN(parts[1], "-", 2)
	if rangeParts[0] == "" {
		return 0, -1
	}

	start, _ = strconv.ParseInt(rangeParts[0], 10, 64)

	if rangeParts[1] != "" {
		end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
	} else {
		end = -1 // До конца файла
	}

	return start, end
}
