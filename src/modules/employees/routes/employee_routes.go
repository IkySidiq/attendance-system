package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"attandance-system/src/modules/employees/handlers"
	"attandance-system/src/modules/employees/model"
	"attandance-system/src/middleware"
)

func RegisterUserRoutes(r *gin.Engine, db *sql.DB) {
	service := model.NewEmployeeService(db)
	employeeHandler := handlers.NewEmployeeHandler(service)

	userGroup := r.Group("/employees")
	{
		userGroup.POST("/register", employeeHandler.RegisterEmployee)
		userGroup.POST("/login", employeeHandler.LoginEmployee)
		userGroup.GET("/", middleware.JWTAuthMiddleware(), employeeHandler.GetAllUsers)
	}
}
