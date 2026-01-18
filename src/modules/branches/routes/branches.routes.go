package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"attandance-system/src/modules/branches/handlers"
	"attandance-system/src/modules/branches/model"
)

func RegisterBranchesRoutes(r *gin.Engine, db *sql.DB) {
	service := model.NewBranchesService(db)
	branchesHandler := handlers.NewBranchesHandler(service)

	branchesGroup := r.Group("/branches")
	{
		branchesGroup.POST("/create", branchesHandler.CreateBranch)
	}
}
