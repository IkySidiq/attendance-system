package dto

import (

)

type CreateEmployeeDTO struct{
	Name string `json:"name" validate:"required,min=3,max=50"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email,no_space"`
	Password string `json:"password" validate:"required,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	EmployeeCode  string `json:"employee_code" validate:"required"`
	BranchId string `json:"branch_id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=active inactive suspended"`
}