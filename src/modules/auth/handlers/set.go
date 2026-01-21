package handlers

import (
	"attandance-system/src/modules/auth/model"
	"attandance-system/src/utils/response_helper"
	"attandance-system/src/modules/auth/dto"
	"attandance-system/src/exceptions"

	"github.com/gin-gonic/gin"

	"net/http"
)

	type EmployeeResponse struct{
		Id string `json:"id"`
		Email string `json:"email"`
		Username string `json:"username"`
		Name string `json:"name"`
		Session *model.CreatedSession `json:"session"`
		RememberMe bool `json:"remember_me"`
	}

type EmployeeHandler struct {
	service  *model.EmployeeService
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
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	// 3️Buat user melalui service
	employee, err := h.service.CreateEmployee(&payload)
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
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	// 4Ambil userID untuk session
	employeeId := employee.ID

	// Ambil device info
	deviceInfo := model.DeviceInfo{
		IP:        ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
	}

	// Buat session (biarkan untuk future feature)
	sessionData, err := h.service.CreateSession(employeeId, deviceInfo, payload.RememberMe)
	if err != nil {
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	employeeResponse := EmployeeResponse{
		Id: employee.ID,
		Email: employee.Email,
		Username: employee.Username,
		Name: employee.Name,
		Session: sessionData,
	}

	// Kirim response
	response.Created(ctx, map[string]interface{}{
		"user":    employeeResponse,
		"session": sessionData,
	}, "User registered successfully", nil)
}

func (h *EmployeeHandler) Login(ctx *gin.Context) {
	var payload dto.LoginDTO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(ctx, err.Error(), nil)
	}

	data, err := h.service.Login(&payload)
	if err != nil {
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
	}

	response.Created(ctx, data, "Login successfully", nil)
}

func (h *EmployeeHandler) GetAllEmployees(ctx *gin.Context) () {
	
}