package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success writes a 2xx JSON response.
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes a 4xx/5xx JSON response.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{
		Success: false,
		Message: message,
	})
}
