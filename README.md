```markdown
# Go Authentication App

This is a Go authentication application designed to handle user sign-up, sign-in, and OTP verification.

## Features

- User sign-up with email verification
- User sign-in with OTP (One-Time Password)
- Resend OTP functionality
- Logout functionality

## Prerequisites

Before running the application, ensure you have the following installed:

- Go (v1.13 or higher)
- Fiber (v2 or higher)
- GORM (v2 or higher)
- AfterShip/email-verifier
- TLS/SSL certificates for secure email transmission
- SMTP credentials (username and password)
- Elastic Email SMTP host (or any other SMTP host you prefer)
- Database (MySQL, PostgreSQL, etc.) configured and running

## Installation

1. Clone the repository:

```
git clone https://github.com/ADEMOLA200/Go-Auth.git
```

2. Navigate to the project directory:

```
cd Go-Auth
```

3. Build and run the application:

```
go build
./Go-Auth
```

## Configuration

Before running the application, make sure to set the following environment variables:

- `SMTP_USER`: SMTP username for email transmission
- `SMTP_PASSWORD`: SMTP password for email transmission

Additionally, ensure that the SMTP host and port are correctly set:

- `smtpHost`: SMTP host of your email service provider (e.g., smtp.elasticemail.com)
- `smtpPort`: SMTP port number of your email service provider (e.g., 2525)

## Usage

### Sign Up

To sign up a new user, send a POST request to `/signup` with the following JSON payload:

```json
{
  "username": "example",
  "email": "example@example.com",
  "password": "password123",
  "confirmPassword": "password123"
}
```

### Sign In

To sign in a user and generate OTP, send a POST request to `/signin` with the following JSON payload:

```json
{
  "email": "example@example.com",
  "password": "password123"
}
```

### Verify OTP

To verify OTP after signing in, send a POST request to `/verifyotp` with the following JSON payload:

```json
{
  "otp": "123456"
}
```

### Resend OTP

To resend OTP to the user's email, send a POST request to `/resendotp` with the following JSON payload:

```json
{
  "email": "example@example.com"
}
```

### Logout

To logout a user, send a GET request to `/logout`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```
