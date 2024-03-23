package middlewares

import "github.com/gofiber/fiber/v2"

// CorsMiddleware handles CORS headers for requests.
func CorsMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Set("Access-Control-Allow-Origin", "*")
        c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Set("Access-Control-Allow-Headers", "Content-Type, Origin, Authorization")

        // Handle preflight requests (OPTIONS)
        if c.Method() == fiber.MethodOptions {
            c.Status(fiber.StatusOK)
            return nil
        }

        return c.Next()
    }
}