package routes

import (
   "github.com/ADEMOLA200/Go-Auth/controller"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

    app.Post("/sign-up", controller.SignUp)
    app.Get("/sign-in", controller.SignIn)
    app.Post("/logout", controller.Logout)

    // Route for OTP verification
    app.Post("/verify-otp", controller.VerifyOTP)

    // Define a new route for resending OTP
    app.Post("/resend-otp", controller.ResendOTPHandler)
    
}