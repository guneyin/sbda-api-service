package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/sbda-api-service/config"
	sdk "github.com/guneyin/sbda-sdk"
)

type IServiceHandler interface {
	Register(r fiber.Router)
}

type Handler struct {
	cfg      *config.Config
	ds       *sdk.DiscoveryService
	handlers map[sdk.ServiceEnum]IServiceHandler
}

func NewHandler(cfg *config.Config, ds *sdk.DiscoveryService) *Handler {
	return &Handler{
		cfg:      cfg,
		ds:       ds,
		handlers: make(map[sdk.ServiceEnum]IServiceHandler),
	}
}

func (h *Handler) RegisterService(r fiber.Router, svc sdk.ServiceEnum) {
	if _, ok := h.handlers[svc]; !ok {
		sh := h.getServiceHandler(svc)
		sh.Register(r)

		h.handlers[svc] = sh
	}
}

func (h *Handler) getServiceHandler(svc sdk.ServiceEnum) IServiceHandler {
	switch svc {
	case sdk.AuthService:
		return NewAuthHandler(h.cfg, h.ds)
	default:
		return nil
	}
}
