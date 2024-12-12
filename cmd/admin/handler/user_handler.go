package handler

import (
	"errors"
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/service"
	"gitlab/live/be-live-api/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	Handler
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

	group.Use(h.JWTMiddleware())
	group.Use(h.RoleGuardMiddleware())
	group.GET("/:page/:limit", h.page)

}

func (h *userHandler) page(c echo.Context) error {
	var page, limit int
	var err error

	page = utils.DEFAULT_PAGE
	limit = utils.DEFAULT_LIMIT

	if c.Param("page") != "" {
		page, err = strconv.Atoi(c.Param("page"))
		if err != nil {
			return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid page parameter"), nil)
		}
	}

	if c.Param("limit") != "" {
		limit, err = strconv.Atoi(c.Param("limit"))
		if err != nil {
			return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid limit parameter"), nil)
		}
	}

	var req dto.UserQuery
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	data, err := h.srv.User.GetUserList(&req, page, limit)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	return utils.BuildSuccessResponseWithData(c, http.StatusOK, data)
}
