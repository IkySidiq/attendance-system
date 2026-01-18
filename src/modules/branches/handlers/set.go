package handlers

import (
	"attandance-system/src/exceptions"
	"attandance-system/src/modules/branches/dto"
	"attandance-system/src/modules/branches/model"
	response "attandance-system/src/utils/response_helper"
	"net/http"

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

		// Build response user
	branchResponse := map[string]interface{}{
		"id":       branch.ID,
		"name":    branch.Name,
		"branchCode": branch.BranchCode,
		"location":     branch.Location,
		"isActive":     branch.IsActive,
		"created_at":    branch.CreatedAt,
	}

	// Kirim response
	response.Created(ctx, map[string]interface{}{
		"user":    branchResponse,
	}, "User registered successfully", nil)
}