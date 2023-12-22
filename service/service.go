package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/guneyin/sbda-api-service/config"
	"github.com/guneyin/sbda-api-service/handlers"
	sdk "github.com/guneyin/sbda-sdk"
	"time"
)

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
)

type ApiService struct {
	cfg *config.Config
	log *sdk.Logger
	app *fiber.App
	hnd *handlers.Handler
}

func NewApiService() (*ApiService, error) {
	cfg := config.GetConfig()
	ds, err := sdk.NewDiscoveryService(cfg.DiscoverySvcAddr)
	if err != nil {
		return nil, err
	}

	as := &ApiService{
		cfg: cfg,
		log: sdk.NewLogger(),
		app: createApp(cfg),
		hnd: handlers.NewHandler(cfg, ds),
	}
	as.registerHandlers()

	return as, nil
}

func createApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ServerHeader:      fmt.Sprintf("%s HTTP Server", cfg.ApiName),
		AppName:           cfg.ApiName,
		EnablePrintRoutes: true,
		ReadTimeout:       _defaultReadTimeout,
		WriteTimeout:      _defaultWriteTimeout,
	})

	return app
}

func (as *ApiService) Serve() error {
	return as.app.Listen(fmt.Sprintf(":%d", as.cfg.HttpPort))
}

func (as *ApiService) registerHandlers() {
	api := as.app.Use(cors.New()).Group("/api")
	api.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	as.hnd.RegisterService(api, sdk.AuthService)
}
