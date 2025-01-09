package handler

import (
	"errors"
	"gitlab/live/be-live-admin/conf"
	"gitlab/live/be-live-admin/dto"
	"gitlab/live/be-live-admin/model"
	"gitlab/live/be-live-admin/service"
	"gitlab/live/be-live-admin/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type adminHandler struct {
	Handler
	r      *echo.Group
	srv    *service.Service
	apiURL string
}

func newAdminHandler(r *echo.Group, srv *service.Service) *adminHandler {
	admin := &adminHandler{
		r:      r,
		srv:    srv,
		apiURL: conf.GetApiFileConfig().Url,
	}

	admin.register()

	return admin
}

func (h *adminHandler) register() {
	group := h.r.Group("api/admins")

	group.Use(h.JWTMiddleware())
	group.GET("/logs", h.getAdminLogs)
	group.GET("/:id", h.byId)

}

func (h *adminHandler) byId(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid id parameter"), nil)
	}

	data, err := h.srv.Admin.ById(uint(id), h.apiURL)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	if data == nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("not found"), nil)
	}

	return utils.BuildSuccessResponseWithData(c, http.StatusOK, data)

}

func (h *adminHandler) getAdminLogs(c echo.Context) error {

	var req dto.AdminLogQuery
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	currentUser := c.Get("user").(*utils.Claims)
	if currentUser.RoleType == model.ADMINROLE {
		req.IsAdmin = true
	}
	req.UserID = currentUser.ID

	data, err := h.srv.Admin.GetAdminLogs(&req)

	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	return utils.BuildSuccessResponseWithData(c, http.StatusOK, data)
}
