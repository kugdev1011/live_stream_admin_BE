package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/labstack/echo/v4"
)

func HashingPassword(content, salt_key string) string {
	data := salt_key + content
	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func VerifyPassword(storedHash, content, salt_key string) bool {
	return HashingPassword(content, salt_key) == storedHash
}

// type Response struct {
// 	Data    interface{} `json:"data,omitempty"`
// 	Message interface{} `json:"message,omitempty"`
// 	Code    interface{} `json:"code,omitempty"`
// }

// func BuildErrorResponse(ctx echo.Context, status int, err error, body interface{}) error {
// 	return BuildStandardResponse(ctx, status, Response{Message: err.Error(), Code: status, Data: body})
// }

func BuildSuccessResponseWithData(ctx echo.Context, status int, body interface{}) error {
	return BuildStandardResponse(ctx, status, Response{Message: "Successfully", Code: status, Data: body})
}

// func BuildStandardResponse(ctx echo.Context, status int, resp Response) error {
// 	return ctx.JSON(status, Response{Data: resp.Data, Code: resp.Code, Message: resp.Message})
// }
