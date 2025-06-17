package v1

import "github.com/gofiber/fiber/v2"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Routes(instance *fiber.App) *fiber.App {
	v1 := instance.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			signUp := auth.Group("/sign-up")
			{
				signUp.Post("/send-otp", h.sendOtp)
				signUp.Post("/verify-otp", h.verifyOtp)
				signUp.Post("/", h.signUp)
			}
			auth.Post("/sign-in", h.signIn)
			auth.Post("/sign-out", h.signOut)
			auth.Post("/refresh", h.refresh)
		}

		protected := v1.Group("", h.AuthMiddleware())
		{
			media := protected.Group("/media")
			{
				media.Post("/", h.createMedia)
				media.Get("/:id", h.getMedia)
				media.Put("/:id", h.updateMedia)
				media.Delete("/:id", h.deleteMedia)
				media.Get("/", h.listMedia)

				streaming := media.Group("/stream")
				{
					streaming.Post("/upload/", h.uploadFile)
					streaming.Get("/download", h.downloadFile)
				}
			}
		}
	}
	return instance
}
