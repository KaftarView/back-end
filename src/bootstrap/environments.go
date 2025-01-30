package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PRIMARY_DB   Database
	PrimaryRedis RedisDB
	Storage      S3
	Applications AppInfo
	Email        EmailInfo
	SuperAdmin   AdminCredentials
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

type S3 struct {
	Buckets   BucketName
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

type BucketName struct {
	EventsBucket   string
	PodcastsBucket string
	JournalsBucket string
	ProfilesBucket string
	NewsBucket     string
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

type AdminCredentials struct {
	Name         string
	EmailAddress string
	Password     string
}

func NewEnvironments() *Env {
	godotenv.Load(".env")

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
		Storage: S3{
			Buckets: BucketName{
				EventsBucket:   os.Getenv("EVENTS_BUCKET_NAME"),
				PodcastsBucket: os.Getenv("PODCASTS_BUCKET_NAME"),
				JournalsBucket: os.Getenv("JOURNALS_BUCKET_NAME"),
				ProfilesBucket: os.Getenv("PROFILES_BUCKET_NAME"),
				NewsBucket:     os.Getenv("NEWS_BUCKET_NAME"),
			},
			Region:    os.Getenv("EVENTS_BUCKET_REGION"),
			AccessKey: os.Getenv("EVENTS_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("EVENTS_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("EVENTS_BUCKET_ENDPOINT"),
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
		SuperAdmin: AdminCredentials{
			Name:         "Admin",
			EmailAddress: os.Getenv("ADMIN_EMAIL"),
			Password:     os.Getenv("ADMIN_PASSWORD"),
		},
	}
}
