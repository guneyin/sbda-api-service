package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/sbda-api-service/config"
	"github.com/guneyin/sbda-api-service/middleware"
	sdk "github.com/guneyin/sbda-sdk"
	pb "github.com/guneyin/sbda-sdk/pb"
	"net/http"
	"strings"
	"time"
)

var _ IServiceHandler = (*AuthHandler)(nil)

type AuthHandler struct {
	cfg *config.Config
	ds  *sdk.DiscoveryService
}

func (ah *AuthHandler) Register(r fiber.Router) {
	route := r.Group("/auth")
	route.Get("/init", ah.Init)
	route.Get("/callback", ah.Callback)
}

func NewAuthHandler(cfg *config.Config, ds *sdk.DiscoveryService) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
		ds:  ds,
	}
}

func (ah *AuthHandler) Init(c *fiber.Ctx) error {
	res, err := ah.init(c)
	if err != nil {
		return middleware.HttpError(c, err)
	}

	c.Cookies("state", res.State)

	return c.Redirect(res.Url, http.StatusFound)
}

func (ah *AuthHandler) init(c *fiber.Ctx) (*pb.InitAuthResponse, error) {
	conn, err := ah.ds.GetServiceConn(sdk.AuthService)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	svc, err := ah.ds.GetAuthServiceClient(conn)
	if err != nil {
		return nil, err
	}

	hostName := ah.cfg.NetworkAlias
	if strings.TrimSpace(hostName) == "" {
		hostName = c.IP()
	}

	callbackUrl := fmt.Sprintf("http://%s:%d/api/auth/callback", hostName, ah.cfg.HttpPort)

	res, err := svc.InitAuth(c.Context(), &pb.InitAuthRequest{CallbackUrl: callbackUrl})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ah *AuthHandler) Callback(c *fiber.Ctx) error {
	res, err := ah.callback(c)
	if err != nil {
		return middleware.HttpError(c, err)
	}

	exp, err := time.Parse(time.DateTime, res.Token.Expiry)
	if err != nil {
		exp = time.Now().Add(time.Hour * 7 * 24)
	}

	c.Cookie(&fiber.Cookie{
		Name:        "auth-token",
		Value:       res.Token.AccessToken,
		Path:        "",
		Domain:      "",
		MaxAge:      0,
		Expires:     exp,
		Secure:      true,
		HTTPOnly:    true,
		SameSite:    "",
		SessionOnly: false,
	})

	return middleware.HttpSuccess(c, "login successful", res)
}

func (ah *AuthHandler) callback(c *fiber.Ctx) (*pb.CallbackResponse, error) {
	conn, err := ah.ds.GetServiceConn(sdk.AuthService)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	svc, err := ah.ds.GetAuthServiceClient(conn)
	if err != nil {
		return nil, err
	}

	code := c.FormValue("code")

	return svc.Callback(c.Context(), &pb.CallbackRequest{Code: code})
}
