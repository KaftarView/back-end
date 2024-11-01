package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PRIMARY_DB    Database
	PrimaryRedis  RedisDB
	PrimaryBucket Bucket
	Applications  AppInfo
	Email         EmailInfo
	Admin         UserInfo
	Moderator     UserInfo
}

type Database struct {
	DB_HOST string
	DB_NAME string
	DB_PORT string
	DB_USER string
	DB_PASS string
}

type RedisDB struct {
	Port     string
	Addr     string
	Password string
	DB       string
}

type Bucket struct {
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

type AppInfo struct {
	BACKGROUND_SERVICE_ENABLED string
	API_SERVICE_ENABLED        string
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
		PrimaryRedis: RedisDB{
			Port:     os.Getenv("REDIS_PORT"),
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       os.Getenv("REDIS_DB"),
		},
		PrimaryBucket: Bucket{
			Region:    os.Getenv("BUCKET_REGION"),
			AccessKey: os.Getenv("BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("BUCKET_ENDPOINT"),
		},
		Applications: AppInfo{
			BACKGROUND_SERVICE_ENABLED: os.Getenv("BACKGROUND_SERVICE_ENABLED"),
			API_SERVICE_ENABLED:        os.Getenv("API_SERVICE_ENABLED"),
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
