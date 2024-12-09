package handler

import (
	"errors"
	"fmt"
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/service"
	"gitlab/live/be-live-api/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type authHandler struct {
	Handler
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

	group.Use(h.JWTMiddleware())
	group.POST("/register", h.signUp)
	group.POST("/resetPassword", h.resetPassword)
	group.POST("/forgetPassword", h.forgetPassword)

}

func (h *authHandler) signUp(c echo.Context) error {

	var registerDTO dto.RegisterDTO

	if err := utils.BindAndValidate(c, &registerDTO); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)

	}

	roleType := registerDTO.RoleType
	if !model.IsValidRoleType(roleType) {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid role type"), nil)
	}

	role, err := h.srv.Role.GetRoleByType(string(roleType))

	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	if role == nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid role type"), nil)
	}

	hashedPassword, err := utils.HashPassword(registerDTO.Password)

	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	claims := c.Get("user").(*utils.Claims)

	user := &model.User{
		Username:     registerDTO.Username,
		Email:        registerDTO.Email,
		PasswordHash: hashedPassword,
		RoleID:       role.ID,
		Role:         *role,
		CreatedByID:  &claims.CreatedByID,
		UpdatedByID:  &claims.CreatedByID,
		OTPExpiresAt: nil,
	}

	if err := h.srv.User.Create(user); err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	return utils.BuildSuccessResponse(c, http.StatusCreated, fmt.Sprintf("%s created successfully", role.Type), user)

}

func (h *authHandler) login(c echo.Context) error {
	var loginDTO dto.LoginDTO

	if err := utils.BindAndValidate(c, &loginDTO); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	user, err := h.srv.User.FindByEmail(loginDTO.Email)

	if err != nil || !utils.CheckPasswordHash(loginDTO.Password, user.PasswordHash) {
		return utils.BuildErrorResponse(c, http.StatusUnauthorized, errors.New("invalid username or password"), nil)
	}

	roleType := model.RoleType(user.Role.Type)
	token, err := utils.GenerateAccessToken(user.Email, roleType, user.ID)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	adminLog := service.CreateAdminLog(user.ID, model.LoginAction, fmt.Sprintf("User %s logged in", user.Email))
	err = h.srv.Admin.CreateLog(adminLog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	response := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role.Type,
		"token":    token,
	}
	return utils.BuildSuccessResponse(c, http.StatusOK, "Login successful", response)
}

func (h *authHandler) forgetPassword(c echo.Context) error {
	var forgetPasswordDTO dto.ForgetPasswordDTO

	if err := utils.BindAndValidate(c, &forgetPasswordDTO); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	// Check if the user exists
	user, err := h.srv.User.FindByEmail(forgetPasswordDTO.Email)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusNotFound, errors.New("email not found"), nil)

	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	user.OTP = otp
	//user.OTPExpiresAt = time.Now().Add(15 * time.Minute)

	user.OTPExpiresAt = func(t time.Time) *time.Time {
		return &t
	}(time.Now().Add(15 * time.Minute))

	err = h.srv.User.Update(user)
	if err != nil {

		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)

	}
	return utils.BuildSuccessResponse(c, http.StatusOK, "OTP generated successfully", map[string]string{"otp": otp})

}

func (h *authHandler) resetPassword(c echo.Context) error {

	var resetPasswordDTO dto.ResetPasswordDTO
	if err := utils.BindAndValidate(c, &resetPasswordDTO); err != nil {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
	}

	if resetPasswordDTO.NewPassword != resetPasswordDTO.ConfirmPassword {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("passwords do not match"), nil)
	}

	claims := c.Get("user").(*utils.Claims)

	user, err := h.srv.User.FindByEmail(claims.Email)
	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusNotFound, errors.New("email not found"), nil)
	}

	if time.Now().After(*user.OTPExpiresAt) {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("OTP expired"), nil)

	}
	if user.OTP != resetPasswordDTO.OTP {
		return utils.BuildErrorResponse(c, http.StatusBadRequest, errors.New("invalid OTP"), nil)

	}

	hashedPassword, err := utils.HashPassword(resetPasswordDTO.NewPassword)

	if err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	if err := h.srv.User.UpdatePassword(user.ID, hashedPassword); err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}
	if err := h.srv.User.ClearOTP(user.ID); err != nil {
		return utils.BuildErrorResponse(c, http.StatusInternalServerError, err, nil)
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, "Password reset successfully", nil)

}
