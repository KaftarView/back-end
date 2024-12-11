package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PRIMARY_DB     Database
	PrimaryRedis   RedisDB
	BannersBucket  Bucket
	SessionsBucket Bucket
	PodcastsBucket Bucket
	ProfileBucket  Bucket
	Applications   AppInfo
	Email          EmailInfo
	Admin          UserInfo
	Moderator      UserInfo
	PayInfo        PayInfo
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

type PayInfo struct {
	ZarinMerchantID string
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
		BannersBucket: Bucket{
			Name:      os.Getenv("BANNERS_BUCKET_NAME"),
			Region:    os.Getenv("BANNERS_BUCKET_REGION"),
			AccessKey: os.Getenv("BANNERS_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("BANNERS_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("BANNERS_BUCKET_ENDPOINT"),
		},
		SessionsBucket: Bucket{
			Name:      os.Getenv("SESSIONS_BUCKET_NAME"),
			Region:    os.Getenv("SESSIONS_BUCKET_REGION"),
			AccessKey: os.Getenv("SESSIONS_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("SESSIONS_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("SESSIONS_BUCKET_ENDPOINT"),
		},
		PodcastsBucket: Bucket{
			Name:      os.Getenv("PODCAST_BUCKET_NAME"),
			Region:    os.Getenv("PODCAST_BUCKET_REGION"),
			AccessKey: os.Getenv("PODCAST_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("PODCAST_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("PODCAST_BUCKET_ENDPOINT"),
		},
		ProfileBucket: Bucket{
			Name:      os.Getenv("PROFILE_BUCKET_NAME"),
			Region:    os.Getenv("PROFILE_BUCKET_REGION"),
			AccessKey: os.Getenv("PROFILE_BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("PROFILE_BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("PROFILE_BUCKET_ENDPOINT"),
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
		PayInfo: PayInfo{
			ZarinMerchantID: os.Getenv("ZARIN_MERCHANT_ID"),
		},
	}
}
