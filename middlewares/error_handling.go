package middlewares

import "github.com/gofiber/fiber/v2"

// ErrorHandlingMiddleware handles errors and sends appropriate responses
func ErrorHandlingMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Proceed to next middleware or route handler
        err := c.Next()

        // If there's an error, send an appropriate response
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }

        return nil
    }
}