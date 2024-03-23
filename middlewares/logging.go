package middlewares

import "github.com/gofiber/fiber/v2"

// LoggingMiddleware logs incoming requests
func LoggingMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Log the incoming request
        println("Incoming request:", c.Path())

        return c.Next()
    }
}