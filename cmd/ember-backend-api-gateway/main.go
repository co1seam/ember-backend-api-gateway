package main

import (
	"flag"
	"fmt"
	"github.com/co1seam/ember_backend_api_gateway"
	"github.com/co1seam/ember_backend_api_gateway/config"
	"github.com/co1seam/ember_backend_api_gateway/http/rest/v1"
	"github.com/co1seam/ember_backend_api_gateway/http/rpc"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfgFlag := flag.String("config", "", "flag to add config path")

	flag.Parse()

	_, err := config.New(cfgFlag)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	err = rpc.NewAuthClient()
	if err != nil {
		panic(err)
	}

	err = rpc.NewMediaClient()
	if err != nil {
		panic(err)
	}

	handler := v1.NewHandler()

	server := ember_backend_api_gateway.NewServer()

	server.Server.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080",
		AllowMethods:     "GET, HEAD, OPTIONS",
		AllowHeaders:     "Origin, Range, Accept, Content-Type",
		ExposeHeaders:    "Content-Range, Content-Length, Accept-Ranges",
		AllowCredentials: true,
	}))

	handler.Routes(server.Server)
	if err := server.Run("8080"); err != nil {
		log.Error(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(); err != nil {
		log.Error("error: ", err)
	}
}
