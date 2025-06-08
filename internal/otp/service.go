package otp

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/jatinnsharma/internal/config"
	"github.com/jatinnsharma/internal/database"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

const (
	PurposeEmailVerification = "email_verification"
	PurposePasswordReset     = "password_reset"
)

func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

func (s *Service) GenerateOTP(email, purpose string) (string, error) {
	// Generate 6-digit OTP
	otp := fmt.Sprintf("%06d", generateRandomNumber(100000, 999999))

	// Delete existing OTPs for this email and purpose
	s.db.Where("email = ? AND purpose = ?", email, purpose).Delete(&database.EmailOTP{})

	// Create new OTP record
	emailOTP := database.EmailOTP{
		Email:     email,
		OTP:       otp,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(s.cfg.OTPExpiry),
	}

	if err := s.db.Create(&emailOTP).Error; err != nil {
		return "", err
	}

	return otp, nil
}

func (s *Service) VerifyOTP(email, otp, purpose string) error {
	var emailOTP database.EmailOTP
	if err := s.db.Where("email = ? AND otp = ? AND purpose = ? AND expires_at > ?", 
		email, otp, purpose, time.Now()).First(&emailOTP).Error; err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Delete used OTP
	s.db.Delete(&emailOTP)

	return nil
}

func generateRandomNumber(min, max int) int {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	
	num := int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	if num < 0 {
		num = -num
	}
	
	return min + (num % (max - min + 1))
}
