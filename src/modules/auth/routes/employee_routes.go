package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"attandance-system/src/modules/auth/handlers"
	"attandance-system/src/modules/auth/model"
	"attandance-system/src/middleware"
)

func RegisterEmployeeRoutes(r *gin.Engine, db *sql.DB) {
	service := model.NewEmployeeService(db)
	employeeHandler := handlers.NewEmployeeHandler(service)

	employeeGroup := r.Group("/employees")
	{
		employeeGroup.POST("/register", employeeHandler.CreateEmployee)
		employeeGroup.POST("/login", employeeHandler.Login)
		employeeGroup.GET("/", middleware.JWTAuthMiddleware(), employeeHandler.GetAllEmployees)
	}
}
