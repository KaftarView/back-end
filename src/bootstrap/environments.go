package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PRIMARY_DB   Database
	PrimaryRedis RedisDB
	Email        EmailInfo
	Admin        UserInfo
	Moderator    UserInfo
}

type Database struct {
	DB_HOST string
	DB_NAME string
	DB_PORT string
	DB_USER string
	DB_PASS string
}

type RedisDB struct {
	Addr     string
	Password string
	DB       int
}

type EmailInfo struct {
	EmailFrom     string
	EmailPassword string
	SMTPHost      string
	SMTPPort      string
}

type UserInfo struct {
	EmailAddress string
	Password     string
}

func NewEnvironments() *Env {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Env{
		PRIMARY_DB: Database{
			DB_HOST: os.Getenv("DB_HOST"),
			DB_NAME: os.Getenv("DB_NAME"),
			DB_PORT: os.Getenv("DB_PORT"),
			DB_USER: os.Getenv("DB_USER"),
			DB_PASS: os.Getenv("DB_PASS"),
		},
		Email: EmailInfo{
			EmailFrom:     os.Getenv("EMAIL_FROM"),
			EmailPassword: os.Getenv("EMAIL_PASSWORD"),
			SMTPHost:      os.Getenv("SMTP_HOST"),
			SMTPPort:      os.Getenv("SMTP_PORT"),
		},
		Admin: UserInfo{
			EmailAddress: os.Getenv("ADMIN_EMAIL"),
			Password:     os.Getenv("ADMIN_PASSWORD"),
		},
		Moderator: UserInfo{
			EmailAddress: os.Getenv("MODERATOR_EMAIL"),
			Password:     os.Getenv("MODERATOR_PASSWORD"),
		},
	}
}
