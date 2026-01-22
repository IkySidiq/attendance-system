package handlers

import (
	"attandance-system/src/exceptions"
	"attandance-system/src/modules/auth/dto"
	"attandance-system/src/modules/auth/model"
	"attandance-system/src/utils"
	"attandance-system/src/utils/response_helper"
	"os"
	"time"

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

func (h *EmployeeHandler) CreateEmployee(ctx *gin.Context) {
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
		return
	}

	deviceInfo := model.DeviceInfo{
		IP:        ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
	}

	employeeData, sessionData, err := h.service.Login(&payload, deviceInfo)
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

	dataToken := map[string]interface{}{
		"id": employeeData.Id,
		"username": employeeData.Username,
		"name": employeeData.Name,
		"email": employeeData.Email,
		"employee_code": employeeData.EmployeeCode,
		"status": employeeData.Status,
		"branch_name": employeeData.BranchName,
	}

	accessToken, err := utils.GenerateAccessToken(dataToken);
	if err != nil {
		exceptions.NewClientError(400, "access token gagal")
	}

	employee := h.buildEmployeeResponse(employeeData)
	token := h.buildTokenResponse(accessToken, sessionData.RefreshToken)
	session := h.buildSessionResponse(sessionData)

	response.Created(
		ctx,
		map[string]interface{}{
		"employee": employee,
		"token": token,
		"session": session,
		}, 
		"Login successfully", 
		nil,
	)
}

func (h *EmployeeHandler) GetAllEmployees(ctx *gin.Context) () {
	
}

//* INTERNAL FUNCTION

func (h *EmployeeHandler) buildTokenResponse(accessToken, refreshToken string) map[string]interface{} {
	return map[string]interface{}{
		"access_token": accessToken,
		"refresh_token": refreshToken,
		"token_type": "Bearer",
		"expires_in": os.Getenv("JWT_EXP"),
		"refresh_expires_in": os.Getenv("JWT_REFRESH_EXP"),
		"issued_at": time.Now(),
	}
}

func (h *EmployeeHandler) buildSessionResponse(session *model.CreatedSession) map[string]interface{} {
	return map[string]interface{}{
		"session_id": session.Id,
		"device_info": map[string]interface{}{
			"user_agent": session.DeviceInfo.UserAgent,
			"ip_address": session.DeviceInfo.IP,
			"device_type": session.DeviceInfo.DeviceType,
			"device_name": session.DeviceInfo.DeviceName,
			"platform":	session.DeviceInfo.Platform,
		},
		"session_data": session.LastActivity,
		"login_time": session.LogInTime,
	}
}

func (h *EmployeeHandler) buildEmployeeResponse(employee *model.ResponseLogin) map[string]interface{} {
	return map[string]interface{}{
		"id": employee.Id,
		"name": employee.Name,
		"status": employee.Status,
	}
}