
package database

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email      string    `gorm:"uniqueIndex;not null" json:"email"`
	Password   string    `gorm:"not null" json:"-"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Session model
type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RefreshToken string    `gorm:"not null" json:"-"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// EmailOTP model
type EmailOTP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"not null" json:"email"`
	OTP       string    `gorm:"not null" json:"-"`
	Purpose   string    `gorm:"not null" json:"purpose"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func ConnectRedis(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	return client
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Session{}, &EmailOTP{})
}