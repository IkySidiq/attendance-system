package dto

import (

)

type BranchDTO struct{
	Name string `json:"name" validate:"required,min=3,max=50"`
	BranchCode string `json:"branch_code" validate:"required,min=3,max=15"`
	Location string `json:"location" validate:"required"`
	IsActive bool `json:"is_active" validate:"required"`
}