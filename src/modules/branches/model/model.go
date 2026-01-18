package model

import (
	"attandance-system/src/exceptions"
	"attandance-system/src/modules/branches/dto"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type BranchesService struct{
	db *sql.DB
}

type CreatedBranch struct{
	ID string
	Name string
	BranchCode string
	Location string
	IsActive bool
	CreatedAt time.Time
}

func NewBranchesService(db *sql.DB) *BranchesService  {
	return &BranchesService{db: db}
}

func (bs *BranchesService) isFieldExists(field string, value string) (bool, error) {
	allowedFields := map[string]bool{
		"branch_code":         true,
	}

	if !allowedFields[field] {
		return false, exceptions.NewClientError(400, "Invalid field")
	}

	query := `
		SELECT EXISTS(
			SELECT 1 FROM branches WHERE ` + field + ` = $1
		)
	`

	var exists bool
	if err := bs.db.QueryRow(query, value).Scan(&exists); err != nil {
		fmt.Println("hor", exists)
		return false, err
	}

	return exists, nil
}

func (bs *BranchesService) CreateBranch(dto *dto.BranchDTO) (*CreatedBranch, error) {
	// 1. Normalisasi data
	normalizedName := strings.TrimSpace(dto.Name)
	normalizedBranchCode := strings.TrimSpace(dto.BranchCode)
	normalizedLocation := strings.TrimSpace(dto.Location)

	if exists, err := bs.isFieldExists("branch_code", normalizedBranchCode); err != nil {
		return nil, err
	} else if exists {
		return nil, exceptions.NewClientError(409, "Branch code already taken")
	}

	branchId, _ := uuid.NewV7()
	now := time.Now()

	// 2. Tentukan Folder dan Path Gambar
	// Folder 'public/qrcodes' biasanya bisa diakses langsung lewat web server
	storageDir := "public/qrcodes"
	fileName := fmt.Sprintf("%s.png", branchId.String())
	filePath := filepath.Join(storageDir, fileName)

	// Buat folder jika belum ada
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// 3. Generate QR Code
	qrContent := fmt.Sprintf("branch_id:%s", branchId.String())
	png, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	// 4. Simpan byte gambar ke file fisik di storageDir
	if err := os.WriteFile(filePath, png, 0644); err != nil {
		return nil, fmt.Errorf("failed to save qr file: %v", err)
	}

	// 5. Simpan ke Database (yang disimpan adalah filePath-nya)
	var result CreatedBranch
	createdBy := "admin_system"

	query := `
		INSERT INTO branches 
			(id, name, branch_code, location, qr_code_attendance, is_active, created_at, created_by)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING 
			id, name, branch_code, location, is_active, created_at
	`

	// qr_code_attendance sekarang berisi link file: "public/qrcodes/xxx.png"
	err = bs.db.QueryRow(
		query,
		branchId,
		normalizedName,
		normalizedBranchCode,
		normalizedLocation,
		filePath, 
		dto.IsActive,
		now,
		createdBy,
	).Scan(
		&result.ID,
		&result.Name,
		&result.BranchCode,
		&result.Location,
		&result.IsActive,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}