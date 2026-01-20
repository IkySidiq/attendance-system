package model

import (
	"attandance-system/src/modules/employees/dto"
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"attandance-system/src/exceptions"

	"github.com/golang-jwt/jwt/v5"
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
	UserId string
	AccessTokenExpiry string
	RefreshToken string
	RefreshTokenExpiry time.Time
	Ip string
	UserAgent string
	LastActivity string
	RememberMe bool
	SessionData string
	Revoked string
	CreatedAt string
}

type DeviceInfo struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
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

func (es *EmployeeService) _isFieldExists(field string, value string) (bool, error) {
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

func _parseExpiration(expiration string) time.Duration {
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

func (es *EmployeeService) CreateEmployee(dto *dto.CreateEmployeeDTO, userId string) (*CreatedEmployee, error) {
	normalizedName := strings.ToLower(strings.TrimSpace(dto.Name))
	normalizedUsername := strings.TrimSpace(dto.Username)
	normalizedEmail := strings.TrimSpace(dto.Email)
	normalizedPassword := strings.TrimSpace(dto.Password)
	normalizedEmployeeCode := strings.ToLower(strings.TrimSpace(dto.EmployeeCode))

	// EMAIL
	if exists, err := es._isFieldExists("email", normalizedEmail); err != nil {
		return nil, err
	} else if exists {
		return nil, exceptions.NewClientError(409, "Email already taken")
	}

	// USERNAME
	if exists, err := es._isFieldExists("username", normalizedUsername); err != nil {
		return nil, err
	} else if exists {
		return nil, exceptions.NewClientError(409, "Username already taken")
	}

	// EMPLOYEE CODE
	if exists, err := es._isFieldExists("employee_code", normalizedEmployeeCode); err != nil {
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

func (es *EmployeeService) CreateSession(userID string, device DeviceInfo) (*CreatedSession, error) {
	tx, err := es.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback() // rollback jika error
	}()

	sessionId, _ := uuid.NewV7()
	now := time.Now()

	accessTokenExpiry := now.Add(_parseExpiration(os.Getenv("JWT_EXP")))
	refreshTokenExpiry := now.Add(_parseExpiration(os.Getenv("JWT_REFRESH_EXP")))

		// Generate JWT refresh token
	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     refreshTokenExpiry.Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	
	secret := []byte(os.Getenv("jill_secret"))
	refreshTokenString, err := refreshToken.SignedString(secret)
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
			id, user_id, expires_at, refresh_token, refresh_expires_at,
			ip, user_agent, last_activity, remember_me, session_data, revoked, created_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,false,CURRENT_TIMESTAMP
		)
		RETURNING id, user_id, expires_at, refresh_token, refresh_expires_at,
		          ip, user_agent, last_activity, remember_me, session_data, revoked, created_at
	`

	var result CreatedSession

	err = tx.QueryRow(
		query,
		sessionId, userID, accessTokenExpiry, refreshTokenString, refreshTokenExpiry,
		device.IP, device.UserAgent, now, false, string(sessionDataJSON),
	).Scan(
		&result.Id, &result.UserId, &result.AccessTokenExpiry, &result.RefreshToken, &result.RefreshTokenExpiry,
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

func (es *EmployeeService) GetAllEmployees(page, limit int, search string) ([]map[string]interface{}, int, error) {
	return nil, 0, nil
}