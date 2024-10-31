package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/repository"
	"first-project/src/routes"
	"first-project/src/seed"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.Default()

	var di = bootstrap.Run()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		di.Env.PRIMARY_DB.DB_USER,
		di.Env.PRIMARY_DB.DB_PASS,
		di.Env.PRIMARY_DB.DB_HOST,
		di.Env.PRIMARY_DB.DB_PORT,
		di.Env.PRIMARY_DB.DB_NAME,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entities.User{}, &entities.Password{})

	dbNumber, _ := strconv.Atoi(di.Env.PrimaryRedis.DB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     di.Env.PrimaryRedis.Addr,
		Password: di.Env.PrimaryRedis.Password,
		DB:       dbNumber,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}

	userRepository := repository.NewUserRepository(db)
	roleSeeder := seed.NewRoleSeeder(userRepository, &di.Env.Admin, &di.Env.Moderator)
	roleSeeder.SeedRoles()

	backgroundEnabled, _ := strconv.ParseBool(di.Env.Applications.BACKGROUND_SERVICE_ENABLED)
	if backgroundEnabled {
		emailService := application_communication.NewEmailService(&di.Env.Email)
		cronJob := application.NewCronJob(userRepository, emailService)
		cronJob.RunCronJob()
	}

	APIServiceEnabled, _ := strconv.ParseBool(di.Env.Applications.API_SERVICE_ENABLED)
	if APIServiceEnabled {
		routes.Run(ginEngine, di, db, rdb)
	}

	ginEngine.Run(":8080")
}
