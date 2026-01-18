package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"attandance-system/cmd/config"
	employeesRoutes "attandance-system/src/modules/employees/routes"
	branchesRoutes "attandance-system/src/modules/branches/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := gin.Default()

	employeesRoutes.RegisterUserRoutes(r, db)
	branchesRoutes.RegisterBranchesRoutes(r, db)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("Server running on port", port)
	r.Run(":" + port)
}