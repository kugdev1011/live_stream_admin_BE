package handler

import (
	"errors"
	"fmt"
	"gitlab/live/be-live-api/conf"
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/service"
	"gitlab/live/be-live-api/utils"
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
	group.Use(h.RoleGuardMiddleware())
	group.POST("", h.createAdmin)
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
	adminLog := service.CreateAdminLog(data.ID, "byId", fmt.Sprintf(" %s by Id request", data.Email), "get_user_by_id")

	err = h.srv.Admin.CreateLog(adminLog)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	return utils.BuildSuccessResponseWithData(c, http.StatusOK, data)

}

func (h *adminHandler) createAdmin(c echo.Context) error {
	var err error
	var req dto.CreateAdminRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	currentUser := c.Get("user").(*utils.Claims)
	req.CreatedByID = &currentUser.CreatedByID
	data, err := h.srv.Admin.CreateAdmin(&req)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	adminLog := service.CreateAdminLog(data.ID, "createAdmin", fmt.Sprintf(" %s created admin", data.Email), "register")

	err = h.srv.Admin.CreateLog(adminLog)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	return utils.BuildSuccessResponseWithData(c, http.StatusAccepted, data)
}
