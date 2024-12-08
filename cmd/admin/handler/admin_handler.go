package handler

import (
	"errors"
	model "gitlab/live/be-live-api/model/api-model"
	"gitlab/live/be-live-api/pkg/utils"
	"gitlab/live/be-live-api/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type adminHandler struct {
	r   *echo.Group
	srv *service.Service
}

func newAdminHandler(r *echo.Group, srv *service.Service) *adminHandler {
	admin := &adminHandler{
		r:   r,
		srv: srv,
	}

	admin.register()

	return admin
}

func (h *adminHandler) register() {
	group := h.r.Group("api/admins")

	group.POST("", h.createAdmin)
	group.GET("/:id", h.byId)

}

func (h *adminHandler) byId(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid id parameter"), nil)
	}

	data, err := h.srv.Admin.ById(uint(id))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	return utils.BuildSuccessResponse(c, http.StatusAccepted, data)

}

func (h *adminHandler) createAdmin(c echo.Context) error {
	var err error
	var req model.CreateAdminRequest
	if err := c.Bind(&req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid request"), nil)
	}

	data, err := h.srv.Admin.CreateAdmin(&req)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	return utils.BuildSuccessResponse(c, http.StatusAccepted, data)
}
