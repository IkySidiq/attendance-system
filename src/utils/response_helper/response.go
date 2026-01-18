package response

import (
	"github.com/gin-gonic/gin"
)

// Response struct umum
type Response struct {
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status,omitempty"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// --------------------
// Helper umum
// --------------------

// Success response
func Success(c *gin.Context, data interface{}, message string, statusCode int, meta interface{}) {
	if message == "" {
		message = "Success"
	}
	if statusCode == 0 {
		statusCode = 200
	}

	resp := Response{
		StatusCode: statusCode,
		Status:     "OK",
		Message:    message,
		Data:       data,
		Meta:       meta,
	}

	c.JSON(statusCode, resp)
}

// Created (201)
func Created(c *gin.Context, data interface{}, message string, meta interface{}) {
	if message == "" {
		message = "Resource created successfully"
	}
	Success(c, data, message, 201, meta)
}

// No Content (204)
func NoContent(c *gin.Context, message string) {
	c.Status(204)
}

// --------------------
// Error helpers
// --------------------

// Error generic
func Error(c *gin.Context, statusCode int, errMsg, message string, data interface{}) {
	if statusCode == 0 {
		statusCode = 500
	}
	if errMsg == "" {
		errMsg = "Internal Server Error"
	}
	if message == "" {
		message = "An error occurred"
	}

	resp := Response{
		StatusCode: statusCode,
		Error:      errMsg,
		Message:    message,
		Data:       data,
	}

	c.JSON(statusCode, resp)
}

// 400 Bad Request
func BadRequest(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Bad Request"
	}
	Error(c, 400, "Bad Request", message, data)
}

// 401 Unauthorized
func Unauthorized(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Unauthorized access"
	}
	Error(c, 401, "Unauthorized", message, data)
}

// 403 Forbidden
func Forbidden(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Access forbidden"
	}
	Error(c, 403, "Forbidden", message, data)
}

// 404 Not Found
func NotFound(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, 404, "Not Found", message, data)
}

// 409 Conflict
func Conflict(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Resource conflict"
	}
	Error(c, 409, "Conflict", message, data)
}

// 422 Unprocessable Entity
func UnprocessableEntity(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Validation failed"
	}
	Error(c, 422, "Unprocessable Entity", message, data)
}

// 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, 500, "Internal Server Error", message, data)
}

// 503 Service Unavailable
func ServiceUnavailable(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	Error(c, 503, "Service Unavailable", message, data)
}

// --------------------
// Validation Error helper
// --------------------
func ValidationError(c *gin.Context, validationDetails interface{}, message string) {
	if message == "" {
		message = "Validation failed"
	}
	UnprocessableEntity(c, message, validationDetails)
}

// --------------------
// Pagination helper
// --------------------
func Paginated(c *gin.Context, data interface{}, pagination interface{}, message string) {
	if message == "" {
		message = "Data retrieved successfully"
	}

	resp := Response{
		StatusCode: 200,
		Status:     "OK",
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	c.JSON(200, resp)
}
