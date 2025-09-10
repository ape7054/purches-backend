package utils

import (
	"net/http"
	"purches-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseOK 生成成功响应
func ResponseOK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// ResponseError 生成错误响应
func ResponseError(c *gin.Context, code int, message string, err string) {
	c.JSON(code, models.APIResponse{
		Code:      code,
		Message:   message,
		Data:      err,
		Timestamp: time.Now(),
	})
}
