package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PRIMARY_DB     Database
	PrimaryRedis   RedisDB
	EventsBucket   Bucket
	PodcastsBucket Bucket
	NewsBucket     Bucket
	JournalsBucket Bucket
	ProfilesBucket Bucket
	Applications   AppInfo
	Email          EmailInfo
	Admin          UserInfo
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
	Name      string
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
		EventsBucket: Bucket{
			Name:      os.Getenv("EVENTS_BUCKET_NAME"),
			Region:    os.Getenv("EVENTS_BUCKET_REGION"),
			AccessKey: os.Getenv("EVENTS_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("EVENTS_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("EVENTS_BUCKET_ENDPOINT"),
		},
		PodcastsBucket: Bucket{
			Name:      os.Getenv("PODCASTS_BUCKET_NAME`"),
			Region:    os.Getenv("PODCASTS_BUCKET_REGION"),
			AccessKey: os.Getenv("PODCASTS_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("PODCASTS_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("PODCASTS_BUCKET_ENDPOINT"),
		},
		NewsBucket: Bucket{
			Name:      os.Getenv("NEWS_BUCKET_NAME"),
			Region:    os.Getenv("NEWS_BUCKET_REGION"),
			AccessKey: os.Getenv("NEWS_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("NEWS_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("NEWS_BUCKET_ENDPOINT"),
		},
		JournalsBucket: Bucket{
			Name:      os.Getenv("JOURNALS_BUCKET_NAME"),
			Region:    os.Getenv("JOURNALS_BUCKET_REGION"),
			AccessKey: os.Getenv("JOURNALS_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("JOURNALS_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("JOURNALS_BUCKET_ENDPOINT"),
		},
		ProfilesBucket: Bucket{
			Name:      os.Getenv("PROFILES_BUCKET_NAME"),
			Region:    os.Getenv("PROFILES_BUCKET_REGION"),
			AccessKey: os.Getenv("PROFILES_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("PROFILES_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("PROFILES_BUCKET_ENDPOINT"),
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
	}
}
