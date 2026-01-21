package handlers

import (
    "math"
    "strconv"

    response "attandance-system/src/utils/response_helper"
    "github.com/gin-gonic/gin"
)

func (h *EmployeeHandler) GetAllUsers(ctx *gin.Context) {
    // Query params dari URL
    page := ctx.DefaultQuery("page", "1")
    limit := ctx.DefaultQuery("limit", "10")
    search := ctx.DefaultQuery("search", "")

    // Konversi page & limit dari string ke int
    pageInt, err := strconv.Atoi(page)
    if err != nil || pageInt < 1 {
        pageInt = 1
    }
    limitInt, err := strconv.Atoi(limit)
    if err != nil || limitInt < 1 {
        limitInt = 10
    }

    users, total, err := h.service.GetAllEmployees(pageInt, limitInt, search)
    if err != nil {
        response.InternalServerError(ctx, "Failed to fetch users", nil)
        return
    }

    // Buat object pagination
    pagination := map[string]interface{}{
        "page":       pageInt,
        "limit":      limitInt,
        "total":      total,
        "totalPages": int(math.Ceil(float64(total) / float64(limitInt))),
    }

    message := "Users retrieved successfully"
    response.Paginated(ctx, users, pagination, message)
}