package handlers

import (
	"attandance-system/src/modules/employees/model"
	"attandance-system/src/utils/response_helper"
	"attandance-system/src/modules/employees/dto"
	"attandance-system/src/exceptions"

	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"

	"net/http"
)

type EmployeeHandler struct {
	service  *model.EmployeeService
	validate *validator.Validate
}

func NewEmployeeHandler(s *model.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service:  s,
	}
}

func (h *EmployeeHandler) RegisterEmployee(ctx *gin.Context) {
	var payload dto.CreateEmployeeDTO

	// 1️Bind JSON
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		// Pesan aman untuk client
		response.BadRequest(ctx, "Invalid request body", nil)
		return
	}

	// 2️Validasi payload (validator)
	if err := h.validate.Struct(payload); err != nil {
		// Bisa ditingkatkan menjadi field-level error nanti
		response.ValidationError(ctx, err.Error(), "Validation failed")
		return
	}

	// 3️Buat user melalui service
	user, err := h.service.CreateEmployee(&payload, "")
	if err != nil {

		// Jika error implement HTTPError, kirim status sesuai^ code
		if httpErr, ok := err.(exceptions.HTTPError); ok {
			response.Error(
				ctx,
				httpErr.StatusCode(),
				http.StatusText(httpErr.StatusCode()),
				httpErr.Error(),
				nil,
			)
			return
		}

		// Error lain = 500
		response.InternalServerError(ctx, "Internal server error", nil)
		return
	}

	// 4Ambil userID untuk session
	userID := user.ID

	// Ambil device info
	deviceInfo := model.DeviceInfo{
		IP:        ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
	}

	// Buat session (biarkan untuk future feature)
	sessionData, err := h.service.CreateSession(userID, deviceInfo)
	if err != nil {
		response.InternalServerError(ctx, "Failed to create session", nil)
		return
	}

	// Build response user
	userResponse := map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"name":     user.Name,
	}

	// Kirim response
	response.Created(ctx, map[string]interface{}{
		"user":    userResponse,
		"session": sessionData,
	}, "User registered successfully", nil)
}

func (h *EmployeeHandler) LoginEmployee(ctx *gin.Context) {
	
}