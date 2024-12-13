package handler

import (
	"errors"
	"gitlab/live/be-live-api/service"
	"gitlab/live/be-live-api/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type streamHandler struct {
	Handler
	r   *echo.Group
	srv *service.Service
}

func newStreamHandler(r *echo.Group, srv *service.Service) *streamHandler {
	statistics := &streamHandler{
		r:   r,
		srv: srv,
	}

	statistics.register()

	return statistics
}

func (h *streamHandler) register() {
	group := h.r.Group("api/streams")

	group.Use(h.JWTMiddleware())
	group.Use(h.RoleGuardMiddleware())
	group.GET("/statistics/:page/:limit", h.getLiveStreamStatisticsData)
	group.GET("/statistics/total", h.getTotalLiveStream)

}

func (h *streamHandler) getTotalLiveStream(c echo.Context) error {

	data, err := h.srv.Stream.GetStatisticsTotalLiveStreamData()
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", data)
}

func (h *streamHandler) getLiveStreamStatisticsData(c echo.Context) error {

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

	data, err := h.srv.Stream.GetStreamAnalyticsData(page, limit)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", data)

}
