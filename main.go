package main

import (
	"github.com/ADEMOLA200/Go-Auth/database"
	"github.com/ADEMOLA200/Go-Auth/middlewares"
	"github.com/ADEMOLA200/Go-Auth/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connect to the database
	database.ConnectDB()
	
    // Create a new Fiber app
    app := fiber.New()

	// Register middleware
	app.Use(middlewares.LoggingMiddleware())
	app.Use(middlewares.ErrorHandlingMiddleware())
	app.Use(middlewares.CorsMiddleware())
	//app.Use(middleware.IsAuthenticated)

	// Initialize routes
	routes.SetupRoutes(app)

	// Start the server on port 3000
	app.Listen(":8000")
}