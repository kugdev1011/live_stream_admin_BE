package handler

import (
	"errors"
	"gitlab/live/be-live-api/pkg/utils"
	"gitlab/live/be-live-api/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	r   *echo.Group
	srv *service.Service
}

func newUserHandler(r *echo.Group, srv *service.Service) *userHandler {
	user := &userHandler{
		r:   r,
		srv: srv,
	}

	user.register()

	return user
}

func (h *userHandler) register() {
	group := h.r.Group("api/users")

	group.POST("/:page/:limit", h.page)

}

func (h *userHandler) page(c echo.Context) error {
	var page, limit int
	var err error

	page = utils.DEFAULT_PAGE
	limit = utils.DEFAULT_LIMIT

	if c.Param("page") != "" {
		page, err = strconv.Atoi(c.Param("page"))
		if err != nil {
			return c.JSON(400, errors.New("invalid page parameter").Error())
		}
	}

	if c.Param("limit") != "" {
		page, err = strconv.Atoi(c.Param("limit"))
		if err != nil {
			return c.JSON(400, errors.New("invalid limit parameter").Error())
		}
	}

	data, err := h.srv.User.GetUserList(page, limit)
	if err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, data)
}
