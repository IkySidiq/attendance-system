package model

import (
	"attandance-system/src/modules/auth/dto"
	"attandance-system/src/utils"
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"attandance-system/src/exceptions"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreatedEmployee struct{
	ID string
	Name string
	Username string
	Email string
	BranchId string
	CreatedAt time.Time
	CreatedBy string
}

type CreatedSession struct{
	Id string
	EmployeeId string
	AccessTokenExpiry string
	RefreshToken string
	RefreshTokenExpiry time.Time
	Ip string
	UserAgent string
	LastActivity string
	RememberMe bool
	SessionData string
	Revoked string
	CreatedAt time.Time
}

type ResponseLogin struct{
	Username string
	Name string
	Email string
	EmployeeCode string
	Status string
	BranchName string
}

type DeviceInfo struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

type DeviceInfo2 struct {
	DeviceType string
	Platform   string
	DeviceName string
}

type SessionData struct{
	LoginTime time.Time
	IPAddress string
	UserAgent string
}

type EmployeeService struct{
	db *sql.DB
}

func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{
		db: db,
	}
}

func (es *EmployeeService) CreateEmployee(dto *dto.CreateEmployeeDTO) (*CreatedEmployee, error) {
	normalizedName := strings.ToLower(strings.TrimSpace(dto.Name))
	normalizedUsername := strings.TrimSpace(dto.Username)
	normalizedEmail := strings.TrimSpace(dto.Email)
	normalizedPassword := strings.TrimSpace(dto.Password)
	normalizedEmployeeCode := strings.ToLower(strings.TrimSpace(dto.EmployeeCode))

	// EMAIL
	if exists, err := es.isFieldExists("email", normalizedEmail); err != nil {
		return nil, err
	} else if exists {
		return nil, exceptions.NewClientError(409, "Email already taken")
	}

	// USERNAME
	if exists, err := es.isFieldExists("username", normalizedUsername); err != nil {
		return nil, err
	} else if exists {
		return nil, exceptions.NewClientError(409, "Username already taken")
	}

	// EMPLOYEE CODE
	if exists, err := es.isFieldExists("employee_code", normalizedEmployeeCode); err != nil {
		return nil, err
	} else if exists {
		return nil, exceptions.NewClientError(409, "Employee code already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(normalizedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	employeeId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	var result CreatedEmployee

	query := `
		INSERT INTO employees
			(id, name, username, email, branch_id, employee_code, status, password, created_at, created_by)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING
			id, name, username, email, branch_id, created_at, created_by
	`

	if err := es.db.QueryRow(
		query,
		employeeId,
		normalizedName,
		normalizedUsername,
		normalizedEmail,
		dto.BranchId,
		normalizedEmployeeCode,
		dto.Status,
		string(hashedPassword),
		now,
		"system_admin",
	).Scan(
		&result.ID, &result.Name, &result.Username, &result.Email, &result.BranchId, &result.CreatedAt, &result.CreatedBy,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

func (es *EmployeeService) CreateSession(employeeId string, device DeviceInfo, rememberMe bool) (*CreatedSession, error) {
	tx, err := es.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback() // rollback jika error
	}()

	rememberMe = utils.ParseBooleanField(rememberMe)

	sessionId, _ := uuid.NewV7()
	now := time.Now()

	accessTokenExpiry := now.Add(parseExpiration(os.Getenv("JWT_EXP")))
	refreshTokenExpiry := now.Add(parseExpiration(os.Getenv("JWT_REFRESH_EXP")))

	deviceInfo2 := parseUserAgent(device.UserAgent)

	refreshTokenString, err := utils.GenerateRefreshToken(map[string]interface{}{
		"employee_id": employeeId,
		"exp": refreshTokenExpiry.Unix(),
	})
	if err != nil {
		return nil, err
	}

	sessionData := SessionData{
		LoginTime: now,
		IPAddress: device.IP,
		UserAgent: device.UserAgent,
	}

	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		return nil, err
	}

	// Query INSERT dengan RETURNING *
	query := `
		INSERT INTO sessions (
			id, employee_id, expires_at, refresh_token, refresh_expires_at, device_type, device_name, platform,
			ip, user_agent, last_activity, remember_me, session_data, revoked, created_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,false,CURRENT_TIMESTAMP
		)
		RETURNING id, employee_id, expires_at, refresh_token, refresh_expires_at,
		          ip, user_agent, last_activity, remember_me, session_data, revoked, created_at
	`

	var result CreatedSession

	err = tx.QueryRow(
		query,
		sessionId, employeeId, accessTokenExpiry, refreshTokenString, refreshTokenExpiry, deviceInfo2.DeviceType,
		deviceInfo2.DeviceName, deviceInfo2.Platform, device.IP, device.UserAgent, now, rememberMe, string(sessionDataJSON),
	).Scan(
		&result.Id, &result.EmployeeId, &result.AccessTokenExpiry, &result.RefreshToken, &result.RefreshTokenExpiry,
		&result.Ip, &result.UserAgent, &result.LastActivity, &result.RememberMe, &result.SessionData, &result.Revoked, &result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (es *EmployeeService) Login(payload *dto.LoginDTO) (*ResponseLogin, error) {
	var responseLogin ResponseLogin

	query := `
	SELECT
		e.username,
		e.name,
		e.email,
		e.employee_code,
		e.status,
		b.branch_name
	FROM
		employees e
	JOIN
		branches b on b.id = e.branch_id
	`

	if err := es.db.QueryRow(query).Scan(
		responseLogin.Username, 
		responseLogin.Name, 
		responseLogin.Email, 
		responseLogin.EmployeeCode, 
		responseLogin.Status, 
		responseLogin.BranchName,
	); err != nil {
		return nil, err
	}

	return &responseLogin, nil
}

func (es *EmployeeService) GetAllEmployees(page, limit int, search string) ([]map[string]interface{}, int, error) {
	return nil, 0, nil
}

//* ===== INTERNAL FUNCTION =====
func (es *EmployeeService) isFieldExists(field string, value string) (bool, error) {
	allowedFields := map[string]bool{
		"email":         true,
		"username":      true,
		"employee_code": true,
	}

	if !allowedFields[field] {
		return false, exceptions.NewClientError(400, "Invalid field")
	}

	query := `
		SELECT EXISTS(
			SELECT 1 FROM employees WHERE ` + field + ` = $1
		)
	`

	var exists bool
	if err := es.db.QueryRow(query, value).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func parseExpiration(expiration string) time.Duration {
	units := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": 24 * time.Hour,
	}

	if len(expiration) < 2 {
		// fallback ke 7 hari
		return 7 * 24 * time.Hour
	}

	valueStr := expiration[:len(expiration)-1]
	unit := strings.ToLower(expiration[len(expiration)-1:])

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		// fallback ke 7 hari kalau konversi gagal
		return 7 * 24 * time.Hour
	}

	if u, ok := units[unit]; ok {
		return time.Duration(value) * u
	}

	// fallback ke 7 hari kalau unit tidak valid
	return 7 * 24 * time.Hour
}

func parseUserAgent(userAgent string) DeviceInfo2 {
	if strings.TrimSpace(userAgent) == "" {
		return DeviceInfo2{
			DeviceType: "other",
			Platform:   "other",
			DeviceName: "Unknown Device",
		}
	}

	ua := strings.ToLower(userAgent)
	deviceType := "desktop"
	platform := "other"
	deviceName := "Unknown Device"

	// device type detection
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		deviceType = "mobile"
	} else if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		deviceType = "tablet"
	} else if strings.Contains(ua, "tv") {
		deviceType = "tv"
	}

	// platform & device name detection (order mirrors JS logic)
	if strings.Contains(ua, "android") {
		platform = "android"
		deviceName = "Android Device"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ios") {
		platform = "ios"
		if strings.Contains(ua, "ipad") {
			deviceName = "iPad"
		} else {
			deviceName = "iPhone"
		}
	} else if strings.Contains(ua, "windows") {
		platform = "windows"
		deviceName = "Windows Device"
	} else if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os") {
		platform = "macos"
		deviceName = "Mac Device"
	} else if strings.Contains(ua, "linux") {
		platform = "linux"
		deviceName = "Linux Device"
	} else if strings.Contains(ua, "chrome") || strings.Contains(ua, "firefox") || strings.Contains(ua, "safari") || strings.Contains(ua, "edge") {
		platform = "web"
		deviceName = "Web Browser"
	}

	return DeviceInfo2{
		DeviceType: deviceType,
		Platform:   platform,
		DeviceName: deviceName,
	}
}