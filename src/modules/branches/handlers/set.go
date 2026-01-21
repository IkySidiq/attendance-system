package handlers

import (
	"attandance-system/src/exceptions"
	"attandance-system/src/modules/branches/dto"
	"attandance-system/src/modules/branches/model"
	response "attandance-system/src/utils/response_helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BranchesHandler struct{
	service *model.BranchesService
}

func NewBranchesHandler(service *model.BranchesService) *BranchesHandler {
	return &BranchesHandler{service: service}
}

func (bh *BranchesHandler) CreateBranch(ctx *gin.Context) {
	var payload dto.BranchDTO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	branch, err := bh.service.CreateBranch(&payload); 
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

		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	type BranchResponse struct {
			Id         string    `json:"id"`
			Name       string    `json:"name"`
			BranchCode string    `json:"branch_code"`
			Location   string    `json:"location"`
			IsActive   bool      `json:"is_active"`
			CreatedAt  time.Time `json:"created_at"`
	}

	var branchResponse = BranchResponse{
		Id: branch.ID,
		Name: branch.Name,
		BranchCode: branch.BranchCode,
		Location: branch.Location,
		IsActive: branch.IsActive,
		CreatedAt: branch.CreatedAt,
	}

	// Kirim response
	response.Created(ctx, branchResponse, "User registered successfully", nil)
}