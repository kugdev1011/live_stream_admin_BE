package handler

import (
	"errors"
	"fmt"
	"gitlab/live/be-live-api/conf"
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/service"
	"gitlab/live/be-live-api/utils"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
)

type streamHandler struct {
	Handler
	r                     *echo.Group
	srv                   *service.Service
	thumbnailFolder       string
	rtmpURL               string
	hlsURL                string
	liveFolder            string
	scheduledVideosFolder string
	ApiURL                string
}

func newStreamHandler(r *echo.Group, srv *service.Service) *streamHandler {

	fileStorageConfig := conf.GetFileStorageConfig()
	streamConfig := conf.GetStreamServerConfig()

	stream := &streamHandler{
		r:                     r,
		srv:                   srv,
		thumbnailFolder:       fileStorageConfig.ThumbnailFolder,
		rtmpURL:               streamConfig.RTMPURL,
		hlsURL:                streamConfig.HLSURL,
		liveFolder:            fileStorageConfig.LiveFolder,
		scheduledVideosFolder: fileStorageConfig.ScheduledVideosFolder,
		ApiURL:                conf.GetApiFileConfig().Url,
	}

	stream.register()

	return stream
}

func (h *streamHandler) register() {
	group := h.r.Group("api/streams")

	group.Use(h.JWTMiddleware())
	group.Use(h.RoleGuardMiddleware())
	group.GET("/statistics", h.getLiveStreamStatisticsData)
	group.GET("/live-statistics", h.getLiveStatData)
	group.GET("/statistics/total", h.getTotalLiveStream)
	group.GET("", h.getLiveStreamWithPagination)
	group.GET("/:id", h.getLiveStreamBroadCastByID)
	group.POST("", h.createLiveStreamByAdmin)
	group.PUT("/:id", h.updateLiveStreamByAdmin)
	group.DELETE("/:id", h.deleteLiveStream)

}

func (h *streamHandler) deleteLiveStream(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid id parameter"), nil)
	}

	deletedStream, err := h.srv.Stream.GetLiveStreamByID(id)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	if deletedStream == nil {
		return utils.BuildErrorResponse(c, http.StatusNotFound, errors.New("not found"), nil)
	}

	if slices.Contains([]model.StreamStatus{model.STARTED, model.ENDED}, deletedStream.Stream.Status) {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("currently, started live-stream, ended live-stream don't allow deleting"), nil)
	}

	// remove thumbnail
	thumbnailPath := fmt.Sprintf("%s%s", h.thumbnailFolder, deletedStream.Stream.ThumbnailFileName)
	thumbnailsToRemove := []string{thumbnailPath}
	go utils.RemoveFiles(thumbnailsToRemove)

	//remove video
	if deletedStream.Stream.StreamType == model.PRERECORDSTREAM && deletedStream.ScheduleStream != nil {
		videoPath := fmt.Sprintf("%s%s", h.scheduledVideosFolder, deletedStream.ScheduleStream.VideoName)
		videosToRemove := []string{videoPath}
		go utils.RemoveFiles(videosToRemove)
	}
	if err := h.srv.Stream.DeleteLiveStream(id); err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	currentUser := c.Get("user").(*utils.Claims)
	adminLog := h.srv.Admin.MakeAdminLogModel(currentUser.ID, model.DeleteLiveStreamByAdmin, fmt.Sprintf(" delete live_stream_broad_cast id %d", id))

	err = h.srv.Admin.CreateLog(adminLog)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", nil)
}

func (h *streamHandler) getLiveStreamBroadCastByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid id parameter"), nil)
	}

	data, err := h.srv.Stream.GetLiveStreamBroadCastByID(id, h.ApiURL, h.rtmpURL, h.hlsURL)

	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	adminLog := h.srv.Admin.MakeAdminLogModel(data.User.ID, model.LiveBroadCastByID, fmt.Sprintf(" %s live_stream_broad_cast request", data.User.DisplayName))

	err = h.srv.Admin.CreateLog(adminLog)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}
	return utils.BuildSuccessResponseWithData(c, http.StatusOK, data)

}

func (h *streamHandler) removeFileLiveStreamAdminByID(liveStream *dto.StreamAndStreamScheduleDto) error {
	// remove thumbnail
	thumbnailPath := fmt.Sprintf("%s%s", h.thumbnailFolder, liveStream.Stream.ThumbnailFileName)
	utils.RemoveFiles([]string{thumbnailPath})

	//remove video
	if liveStream.ScheduleStream != nil {
		videoPath := fmt.Sprintf("%s%s", h.scheduledVideosFolder, liveStream.ScheduleStream.VideoName)
		utils.RemoveFiles([]string{videoPath})
	}
	return nil
}

func (h *streamHandler) updateLiveStreamByAdmin(c echo.Context) error {
	var req dto.StreamRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid id parameter"), nil)
	}

	streamer, err := h.srv.User.CheckUserTypeByID(int(req.UserID))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	if streamer == nil || streamer.Role.Type != model.STREAMER {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("user is not a streamer"), nil)
	}

	if !utils.IsValidSchedule(req.ScheduledAt) {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid schedule"), nil)
	}

	liveStreamForRemoveFiles, err := h.srv.Stream.GetLiveStreamByID(id)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, fmt.Sprintf("get live stream broadCast by ID failed: %s", err.Error()))
	}

	//
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("thumbnail field is required: %s", err.Error()))
	}

	claims := c.Get("user").(*utils.Claims)
	isImage, err := utils.IsImage(file)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	if !isImage {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("file is not an image"), nil)
	}

	if file.Size > utils.MAX_IMAGE_SIZE {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, nil, "Image size exceeds the 1MB limit")
	}
	// save thumbnail
	fileExt := utils.GetFileExtension(file)
	req.ThumbnailFileName = fmt.Sprintf("%d_%s%s", req.UserID, utils.MakeUniqueIDWithTime(), fileExt)
	thumbnailPath := fmt.Sprintf("%s%s", h.thumbnailFolder, req.ThumbnailFileName)

	filesToRemove := []string{thumbnailPath}

	src, err := file.Open()
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	defer src.Close()

	dst, err := os.Create(thumbnailPath)
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	//save recording
	video, err := c.FormFile("video")
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("video field is required: %s", err.Error()))
	}

	if video.Size > utils.MAX_VIDEO_SIZE {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, nil, "Video size exceeds the 2GB limit")
	}

	isVideo, err := utils.IsVideoFile(video)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	if !isVideo {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("file is not a supported video format"), nil)
	}

	fileVideoExt := utils.GetFileExtension(video)
	req.VideoFileName = fmt.Sprintf("%d_%s%s", req.UserID, utils.MakeUniqueIDWithTime(), fileVideoExt)
	videoPath := fmt.Sprintf("%s%s", h.scheduledVideosFolder, req.VideoFileName)

	filesToRemove = append(filesToRemove, videoPath)

	videoSrc, err := video.Open()
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	defer videoSrc.Close()

	dstRecord, err := os.Create(videoPath)
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	defer dstRecord.Close()

	if _, err = io.Copy(dstRecord, videoSrc); err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	//

	stream, err := h.srv.Stream.UpdateStreamByAdmin(id, &req)
	if err != nil {
		go utils.RemoveFiles(filesToRemove)

		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	// remove old files
	if err := h.removeFileLiveStreamAdminByID(liveStreamForRemoveFiles); err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	currentUser := c.Get("user").(*utils.Claims)
	adminLog := h.srv.Admin.MakeAdminLogModel(currentUser.ID, model.UpdateStreamByAdmin, fmt.Sprintf(" %s update_live_stream_by_admin request", claims.Email))

	err = h.srv.Admin.CreateLog(adminLog)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", map[string]any{
		"id":            stream.ID,
		"title":         stream.Title,
		"description":   stream.Description,
		"thumbnail_url": utils.MakeThumbnailURL(h.ApiURL, stream.ThumbnailFileName),
	})
}

func (h *streamHandler) createLiveStreamByAdmin(c echo.Context) error {
	var req dto.StreamRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	streamer, err := h.srv.User.CheckUserTypeByID(int(req.UserID))
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	if streamer == nil || streamer.Role.Type != model.STREAMER {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("user is not a streamer"), nil)
	}

	if !utils.IsValidSchedule(req.ScheduledAt) {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid schedule"), nil)
	}

	file, err := c.FormFile("thumbnail")
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("thumbnail field is required: %s", err.Error()))
	}

	claims := c.Get("user").(*utils.Claims)
	isImage, err := utils.IsImage(file)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	if !isImage {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("file is not an image"), nil)
	}

	if file.Size > utils.MAX_IMAGE_SIZE {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, nil, "Image size exceeds the 1MB limit")
	}
	// save thumbnail
	fileExt := utils.GetFileExtension(file)
	req.ThumbnailFileName = fmt.Sprintf("%d_%s%s", req.UserID, utils.MakeUniqueIDWithTime(), fileExt)
	thumbnailPath := fmt.Sprintf("%s%s", h.thumbnailFolder, req.ThumbnailFileName)

	filesToRemove := []string{thumbnailPath}

	src, err := file.Open()
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	defer src.Close()

	dst, err := os.Create(thumbnailPath)
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	//save recording
	video, err := c.FormFile("video")
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("video field is required: %s", err.Error()))
	}

	if video.Size > utils.MAX_VIDEO_SIZE {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, nil, "Video size exceeds the 2GB limit")
	}

	isVideo, err := utils.IsVideoFile(video)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	if !isVideo {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("file is not a supported video format"), nil)
	}

	fileVideoExt := utils.GetFileExtension(video)
	req.VideoFileName = fmt.Sprintf("%d_%s%s", req.UserID, utils.MakeUniqueIDWithTime(), fileVideoExt)
	videoPath := fmt.Sprintf("%s%s", h.scheduledVideosFolder, req.VideoFileName)

	filesToRemove = append(filesToRemove, videoPath)

	videoSrc, err := video.Open()
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	defer videoSrc.Close()

	dstRecord, err := os.Create(videoPath)
	if err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	defer dstRecord.Close()

	if _, err = io.Copy(dstRecord, videoSrc); err != nil {
		go utils.RemoveFiles(filesToRemove)
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}
	//

	stream, err := h.srv.Stream.CreateStreamByAdmin(&req)
	if err != nil {
		go utils.RemoveFiles(filesToRemove)

		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	adminLog := h.srv.Admin.MakeAdminLogModel(req.UserID, model.LiveStreamByAdmin, fmt.Sprintf(" %s create_live_stream_by_admin request", claims.Email))

	err = h.srv.Admin.CreateLog(adminLog)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	return utils.BuildSuccessResponse(c, http.StatusCreated, "Successfully", map[string]any{
		"id":            stream.ID,
		"title":         stream.Title,
		"description":   stream.Description,
		"thumbnail_url": utils.MakeThumbnailURL(h.ApiURL, stream.ThumbnailFileName),
	})
}

func (h *streamHandler) getTotalLiveStream(c echo.Context) error {

	data, err := h.srv.Stream.GetStatisticsTotalLiveStreamData()

	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", data)
}

func (h *streamHandler) getLiveStreamStatisticsData(c echo.Context) error {

	var req dto.StatisticsQuery
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	data, err := h.srv.Stream.GetStreamAnalyticsData(&req)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", data)

}

func (h *streamHandler) getLiveStatData(c echo.Context) error {

	var req dto.LiveStatQuery
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	data, err := h.srv.Stream.GetLiveStatWithPagination(&req)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", data)

}

func (h *streamHandler) getLiveStreamWithPagination(c echo.Context) error {

	var req dto.LiveStreamBroadCastQueryDTO
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	data, err := h.srv.Stream.GetLiveStreamBroadCastWithPagination(req.Page, req.Limit, &req, h.ApiURL, h.rtmpURL, h.hlsURL)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, "Successfully", data)

}
