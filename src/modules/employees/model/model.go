package model

import (
	"attandance-system/src/modules/employees/dto"
	"database/sql"
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

type DeviceInfo struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

type EmployeeService struct{
	db *sql.DB
}

func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{
		db: db,
	}
}

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

func (es *EmployeeService) CreateEmployee(dto *dto.CreateEmployeeDTO, userId string) (*CreatedEmployee, error) {
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

	var (
		id, name, username, email, branchId, createdBy string
		createdAt time.Time
	)

	query := `
		INSERT INTO employeem
			(id, name, username, email, branch_id, employee_code, status, password, created_at, created_by)
		VALUem
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
		now,
	).Scan(
		&id, &name, &username, &email, &branchId, &createdAt, &createdBy,
	); err != nil {
		return nil, err
	}

	return &CreatedEmployee{
		ID:        id,
		Name:      name,
		Username:  username,
		Email:     email,
		BranchId:  branchId,
		CreatedAt: createdAt,
		CreatedBy: createdBy,
	}, nil
}

func (es *EmployeeService) CreateSession(userID string, device DeviceInfo) (map[string]interface{}, error) {
	return nil, nil
}

func (es *EmployeeService) GetAllEmployees(page, limit int, search string) ([]map[string]interface{}, int, error) {
	return nil, 0, nil
}