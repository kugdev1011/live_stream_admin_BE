package handler

import (
	"gitlab/live/be-live-api/service"

	"github.com/labstack/echo/v4"
)

type authHandler struct {
	r   *echo.Group
	srv *service.Service
}

func newAuthHandler(r *echo.Group, srv *service.Service) *authHandler {
	auth := &authHandler{
		r:   r,
		srv: srv,
	}

	auth.register()

	return auth
}

func (h *authHandler) register() {
	group := h.r.Group("api/auth")

	group.POST("/login", h.login)

}

func (h *authHandler) login(c echo.Context) error {
	return c.JSON(200, nil)
}
