package controller

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/ADEMOLA200/Go-Auth/database"
	"github.com/ADEMOLA200/Go-Auth/models"
	"github.com/AfterShip/email-verifier"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	verifier = emailverifier.NewVerifier()
)

const (
	smtpUserEnv          = "SMTP_USER"     // SMTP username environment variable make sure to set SMTP_USER env variables
	smtpPasswordEnv      = "SMTP_PASSWORD" // SMTP password environment variable make sure to set SMTP_PASSWORD env variables
	smtpHost             = "smtp.elasticemail.com" // SMTP host of Elastic Email make sure to use your own SMTP host
	smtpPort             = 2525 // SMTP port number of Elastic Email make sure to use your own SMTP port
	authentication       = "plain"
	enable_starttls_auto = true
)

// extractDomain extracts the domain part from an email address
func extractDomain(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) == 2 {
        return parts[1] // return the domain part
    }
    return "" // return empty string if email format is invalid
}

// SignUp handles the user signup process
func SignUp(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Invalid JSON format",
        })
    }

    if user.Password != user.ConfirmPassword {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Passwords do not match",
        })
    }

    result, err := verifier.Verify(user.Email)
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Error verifying email",
        })
    }

    if !result.Syntax.Valid {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Invalid email syntax",
        })
    }

    // Extract domain from the email address
    domain := extractDomain(user.Email)

    // Check if the email domain is allowed
	allowedDomains := []string{"gmail.com", "email.com", "yahoo.com", "outlook.com", "hotmail.com", "aol.com", "mail.com"}
	domainAllowed := false
	for _, allowedDomain := range allowedDomains {
		if domain == allowedDomain {
			domainAllowed = true
			break
		}
	}
	if !domainAllowed {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Email address must be from " + strings.Join(allowedDomains, ", "),
		})
	}

    if result.Disposable {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Disposable email not allowed",
        })
    }

    if result.Reachable == "no" {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Email address not reachable",
        })
    }

    if !result.HasMxRecords {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Domain not properly set up to receive emails",
        })
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Error hashing password",
        })
    }
    user.Password = string(hashedPassword)

    r := database.DB.Create(&user)
    if r.Error != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Error creating user",
        })
    }

    c.Status(http.StatusOK)
    return c.JSON(fiber.Map{
        "message": "User created successfully",
    })
}


// Controller logic for sending OTP after successful sign-in attempt
func SignIn(c *fiber.Ctx) error {
    var loginRequest map[string]string

    if err := c.BodyParser(&loginRequest); err != nil {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Invalid JSON format",
        })
    }

    var user models.User

    r := database.DB.Where("email = ?", loginRequest["email"]).First(&user)
    if r.Error != nil {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "User not found",
        })
    }

    // Verify user's password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest["password"]))
    if err != nil {
        c.Status(http.StatusUnauthorized)
        return c.JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    // Generate OTP
    otp, err := generateOTP()
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Failed to generate OTP",
        })
    }

    // Log the generated OTP
    fmt.Println("Generated OTP:", otp)

    // Associate OTP with user
    user.OTP = otp
    currentTime := time.Now()
    user.OTPTime = &currentTime // Set OTPTime to current time
    if err := database.DB.Save(&user).Error; err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Failed to save OTP",
        })
    }

    // Send OTP to user's email
    err = sendOTPEmail(user.Email, otp)
    if err != nil {
        fmt.Printf("Error sending OTP email to %s: %v\n", user.Email, err)
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": fmt.Sprintf("Error sending OTP email to %s: %v", user.Email, err),
        })
    }

	log.Println("OTP has been sent to the user's email: ", user.Email)

    c.Status(http.StatusOK)
    return c.JSON(fiber.Map{
        "message": "OTP has been sent to your email login with the OTP provided. Please check your email.",
    })
}

// Controller logic for verifying OTP
func VerifyOTP(c *fiber.Ctx) error {
    var verifyRequest map[string]string

    if err := c.BodyParser(&verifyRequest); err != nil {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Invalid JSON format",
        })
    }

    var user models.User

    // Fetch the user based on the OTP
    r := database.DB.Where("otp = ?", verifyRequest["otp"]).First(&user)
    if r.Error != nil {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Invalid OTP",
        })
    }

   // Check if OTP has expired
	if user.OTPTime != nil {
		creationTime := *user.OTPTime
		expirationTime := creationTime.Add(1 * time.Minute)
		if time.Now().After(expirationTime) {
			c.Status(http.StatusBadRequest)
			log.Println("OTP has expired for user:", user.ID)
			return c.JSON(fiber.Map{
				"error": "OTP has expired, request a new OTP code.",
			})
		}
	}

	if user.ID != 0 {
		log.Println("User found with OTP:", user.ID) // Log user found with OTP
	} else {
		log.Println("User not found with OTP") // Log user not found with OTP
	}

    // Clear OTP after successful verification
    user.OTP = ""
    if err := database.DB.Save(&user).Error; err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Failed to clear OTP",
        })
    }

    // Log successful sign-in
    log.Println("User", user.Username, "has successfully signed in")

    // Proceed with sign-in
    // You can set up session or JWT token here

    c.Status(http.StatusOK)
    return c.JSON(fiber.Map{
        "message": "OTP verification successful. You can now sign in.",
    })
}

// ResendOTP resends OTP to the user's email
func ResendOTP(email string) error {
    var user models.User
    if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return err
    }

    otp, err := generateOTP()
    if err != nil {
        return err
    }

    // Log the generated OTP
    fmt.Println("ReGenerated OTP:", otp)

    user.OTP = otp
    // Update OTPTime to the current time
    currentTime := time.Now()
    user.OTPTime = &currentTime
    if err := database.DB.Save(&user).Error; err != nil {
        return err
    }

    if err := sendOTPEmail(email, otp); err != nil {
        return err
    }

	log.Println("OTP has been resent to user's email: ", user.Email)

    return nil
}

// ResendOTPHandler handles the request to resend OTP
func ResendOTPHandler(c *fiber.Ctx) error {
	var resendRequest map[string]string

	// Parse the request body
	if err := c.BodyParser(&resendRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Retrieve the email from the request
	email := resendRequest["email"]

	// Call the existing ResendOTP function with the email
	if err := ResendOTP(email); err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to resend OTP: %v", err),
		})
	}

	// Respond with success message
	return c.JSON(fiber.Map{
		"message": "OTP resent to your email",
	})
}

// generateOTP generates a random OTP
func generateOTP() (string, error) {
	// Generate a random 6-digit OTP
	otp := ""
	for i := 0; i < 6; i++ {
		otp += string(rune('0' + rand.Intn(10)))
	}
	if len(otp) != 6 {
		return "", fmt.Errorf("failed to generate OTP")
	}
	return otp, nil
}

// sendOTPEmail sends OTP to the user's email
func sendOTPEmail(email, otp string) error {
	// Retrieve SMTP credentials from environment variables
	smtpUser := os.Getenv(smtpUserEnv)
	smtpPassword := os.Getenv(smtpPasswordEnv)

	// Check if SMTP credentials are empty
	if smtpUser == "" || smtpPassword == "" {
		return errors.New("SMTP credentials not set")
	}

	// Set up TLS configuration
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	// Connect to the SMTP server with TLS
	client, err := smtp.Dial(smtpHost + ":" + strconv.Itoa(smtpPort))
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer client.Close()

	// Start TLS encryption
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	// Authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Load the HTML email template
	htmlTemplate := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>OTP Email</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				padding: 20px;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				background-color: #fff;
				padding: 20px;
				border-radius: 10px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			.header {
				background-color: #3498db;
				color: #fff;
				text-align: center;
				padding: 10px 0;
				border-top-left-radius: 10px;
				border-top-right-radius: 10px;
			}
			.content {
				padding: 20px;
			}
			.footer {
				text-align: center;
				padding: 10px 0;
				border-bottom-left-radius: 10px;
				border-bottom-right-radius: 10px;
			}
			.otp {
				font-size: 24px;
				text-align: center;
				margin-bottom: 20px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h2> Ready  Foods </h2>
			</div>
			<div class="content">
				<p>Dear User,</p>
				<p>Your OTP (One-Time Password) for sign-in is:</p>
				<div class="otp">{{.OTP}}</div>
				<p>Please use this OTP to proceed with your sign-in.</p>
			</div>
			<div class="footer">
				<p>This is an automated email. Please do not reply.</p>
				<p> If you didn't request this OTP, please contact us immediately.</p>
			</div>
		</div>
	</body>
	</html>
	`

	// Create a new template and parse the HTML
	t := template.New("emailTemplate")
	t, err = t.Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	// Prepare data to be passed into the template
	data := struct {
		OTP string
	}{
		OTP: otp,
	}

	// Execute the template to generate the HTML body
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	htmlBody := tpl.String()

	// Compose the email message
	fromAddress := "Email Verification  <no-reply@example.com>"
	toAddress := email
	subject := "Your OTP for sign-in"
	contentType := "text/html; charset=UTF-8"
	msg := []byte("From: " + fromAddress + "\r\n" +
		"To: " + toAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"\r\n" +
		htmlBody)

	// Send the email
	if err := client.Mail(smtpUser); err != nil {
		return fmt.Errorf("failed to send MAIL command: %v", err)
	}
	if err := client.Rcpt(toAddress); err != nil {
		return fmt.Errorf("failed to send RCPT command: %v", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %v", err)
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	return nil
}


// Logout handles the user logout process
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}