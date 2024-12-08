package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/service"
	"gitlab/live/be-live-api/utils"
	"net/http"
	"time"
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
	group.POST("/register", h.signUp)
	group.POST("/forgetPassword", h.forgetPassword)

	group.Use(h.JWTMiddleware())
	group.POST("/resetPassword", h.resetPassword)

}

func (h *authHandler) signUp(c echo.Context) error {

	var registerDTO dto.RegisterDTO

	if err := c.Bind(&registerDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if err := c.Validate(&registerDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	roleType := registerDTO.RoleType
	if roleType != model.ADMINROLE && roleType != model.USERROLE && roleType != model.GUESTROLE {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role type"})
	}

	// Find the Role by Type
	role, err := h.srv.Role.GetRoleByType(string(roleType))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate role"})
	}
	if role == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role type"})
	}

	hashedPassword, err := utils.HashPassword(registerDTO.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not hash password"})
	}

	user := &model.User{
		Username:     registerDTO.Username,
		Email:        registerDTO.Email,
		PasswordHash: hashedPassword,
		RoleID:       role.ID,
	}

	if err := h.srv.User.Create(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})

}

func (h *authHandler) login(c echo.Context) error {
	var loginDTO dto.LoginDTO
	if err := c.Bind(&loginDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if err := c.Validate(&loginDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Find user by username
	user, err := h.srv.User.FindByEmail(loginDTO.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	// Check password
	if !utils.CheckPasswordHash(loginDTO.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	// Generate access token

	roleType := model.RoleType(user.Role.Type)
	token, err := utils.GenerateAccessToken(user.Email, roleType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate access token"})
	}

	//
	adminLog := &model.AdminLog{

		UserID:  user.ID,
		Action:  string(model.LoginAction),
		Details: fmt.Sprintf("User %s logged in", user.Email),
	}

	err = h.srv.Admin.Create(adminLog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	response := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role.Type,
		"token":    token,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *authHandler) forgetPassword(c echo.Context) error {
	var forgetPasswordDTO dto.ForgetPasswordDTO
	if err := c.Bind(&forgetPasswordDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if err := c.Validate(&forgetPasswordDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Check if the user exists
	user, err := h.srv.User.FindByEmail(forgetPasswordDTO.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Email not found"})
	}

	// Generate a password reset token (you could also send an email here)
	//resetToken, err := utils.GenerateAccessToken(user.Email, model.RoleType(user.Role.Type))
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate reset token"})
	//}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate OTP"})
	}

	//hashedOTP, err := utils.HashOTP(otp)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash OTP"})
	//}
	user.OTP = otp
	user.OTPExpiresAt = time.Now().Add(15 * time.Minute)

	err = h.srv.User.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}
	adminLog := &model.AdminLog{

		UserID:  user.ID,
		Action:  string(model.LoginAction),
		Details: fmt.Sprintf("User %s logged in", user.Email),
	}

	err = h.srv.Admin.Create(adminLog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to created admin log"})
	}

	return c.JSON(http.StatusOK, map[string]string{"otp": otp})
}

func (h *authHandler) resetPassword(c echo.Context) error {

	var resetPasswordDTO dto.ResetPasswordDTO

	if err := c.Bind(&resetPasswordDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&resetPasswordDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if resetPasswordDTO.NewPassword != resetPasswordDTO.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Passwords do not match"})
	}

	claims := c.Get("user").(*utils.Claims)
	email := claims.Email

	user, err := h.srv.User.FindByEmail(email)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Email not found"})
	}

	if time.Now().After(user.OTPExpiresAt) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "OTP expired"})
	}
	if user.OTP != resetPasswordDTO.OTP {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid OTP"})
	}

	hashedPassword, err := utils.HashPassword(resetPasswordDTO.NewPassword)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not hash password"})
	}
	err = h.srv.User.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
	}

	// Clear the OTP and expiration time
	err = h.srv.User.ClearOTP(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to clear OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}
