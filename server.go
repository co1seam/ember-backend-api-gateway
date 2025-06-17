package ember_backend_api_gateway

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Server *fiber.App
}

func NewServer() *Server {
	return &Server{
		Server: fiber.New(fiber.Config{
			AppName:     "API Gateway",
			Prefork:     false,
			JSONDecoder: json.Unmarshal,
			JSONEncoder: json.Marshal,
			BodyLimit:   100 * 1024 * 1024,
		}),
	}
}

func (s *Server) Run(port string) error {
	return s.Server.Listen(":" + port)
}

func (s *Server) Shutdown() error {
	return s.Server.Shutdown()
}
