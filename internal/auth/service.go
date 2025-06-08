package auth

import (
	"errors"
	"time"

	"github.com/jatinnsharma/internal/config"
	"github.com/jatinnsharma/internal/database"
	"github.com/jatinnsharma/internal/token"
	"github.com/jatinnsharma/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db           *gorm.DB
	tokenService *token.Service
	cfg          *config.Config
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         *database.User `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresAt    time.Time      `json:"expires_at"`
}

func NewService(db *gorm.DB, tokenService *token.Service, cfg *config.Config) *Service {
	return &Service{
		db:           db,
		tokenService: tokenService,
		cfg:          cfg,
	}
}

func (s *Service) Signup(req *SignupRequest, userAgent, ipAddress string) (*database.User, error) {
	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password strength
	if !utils.IsValidPassword(req.Password) {
		return nil, errors.New("password must be at least 8 characters with uppercase, lowercase, and digit")
	}

	// Sanitize email
	email := utils.SanitizeEmail(req.Email)

	// Check if user already exists
	var existingUser database.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BCryptCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := database.User{
		Email:      email,
		Password:   string(hashedPassword),
		IsVerified: false,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Login(req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	email := utils.SanitizeEmail(req.Email)

	// Find user
	var user database.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Create session
	session := database.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(s.cfg.RefreshExpiry),
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.cfg.JWTExpiry),
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// Find session
	var session database.Session
	if err := s.db.Preload("User").Where("refresh_token = ? AND expires_at > ?", refreshToken, time.Now()).First(&session).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new tokens
	accessToken, err := s.tokenService.GenerateAccessToken(session.User.ID.String())
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Update session with new refresh token
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(s.cfg.RefreshExpiry)
	
	if err := s.db.Save(&session).Error; err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         &session.User,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(s.cfg.JWTExpiry),
	}, nil
}

func (s *Service) Logout(refreshToken string) error {
	// Delete session
	return s.db.Where("refresh_token = ?", refreshToken).Delete(&database.Session{}).Error
}