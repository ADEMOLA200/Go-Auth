package models

import "time"

type User struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	FirstName        string     `json:"firstname" gorm:"not null"`
	LastName         string     `json:"lastname" gorm:"not null"`
	Username         string     `json:"username" gorm:"unique;type:varchar(100);not null"`
	Email            string     `json:"email" gorm:"unique;type:varchar(191)"`
	Password         string     `json:"password" gorm:"not null"`
	ConfirmPassword  string     `json:"confirm-password" gorm:"-"`
	EmailVerified    bool       `gorm:"default:false"`
	VerificationCode string     `json:"verification_code"`
	OTP              string     `json:"otp"`
	OTPTime          *time.Time `json:"otp_time,omitempty"`
}