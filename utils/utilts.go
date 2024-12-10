package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func BuildSuccessResponseWithData(ctx echo.Context, status int, body interface{}) error {
	return BuildStandardResponse(ctx, status, Response{Message: "Successfully", Code: status, Data: body})
}

func ConvertBytes(bytes int64) string {
	if bytes < 1024 {
		return strconv.FormatInt(bytes, 10) + "B" // bytes
	}

	// Define units for bytes, KB, MB, GB, etc.
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	// Convert the bytes to the appropriate unit
	var unitIndex int
	size := float64(bytes)

	// Iterate to find the appropriate size and unit
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	// Format the size to two decimal places
	return fmt.Sprintf("%.1f%s", size, units[unitIndex])
}

func ConvertTimestampToDuration(timestamp int64) string {
	timestampTime := time.Unix(timestamp, 0)
	duration := time.Since(timestampTime)
	hours := int(duration.Hours())          // Total hours
	minutes := int(duration.Minutes()) % 60 // Remainder minutes
	return fmt.Sprintf("%d hours %d minutes", hours, minutes)
}

const DATETIME_LAYOUT = "2024-12-09 20:56:08.408 +0700"

func ConvertDatetimeToTimestamp(datetimeStr, layout string) (int64, error) {
	timestampTime, err := time.Parse(layout, datetimeStr)
	if err != nil {
		return 0, err
	}
	return timestampTime.Unix(), nil
}
