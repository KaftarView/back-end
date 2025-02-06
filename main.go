package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/routes"
	"first-project/src/websocket"
	"first-project/src/wire"
)

func main() {
	gin.DisableConsoleColor()
	ginEngine := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5174", "http://localhost:5173", "https://cesaiust.ir", "http://cesaiust.ir"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	ginEngine.Use(cors.New(config))

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
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&entities.Category{},
		&entities.ChatMessage{},
		&entities.ChatRoom{},
		&entities.Comment{},
		&entities.Discount{},
		&entities.Episode{},
		&entities.Event{},
		&entities.Media{},
		&entities.Commentable{},
		&entities.Councilor{},
		&entities.Journal{},
		&entities.Organizer{},
		&entities.Permission{},
		&entities.Podcast{},
		&entities.News{},
		&entities.Order{},
		&entities.OrderItem{},
		// &entities.Purchasable{},
		&entities.Reservation{},
		&entities.ReservationItem{},
		&entities.Role{},
		&entities.Ticket{},
		&entities.Transaction{},
		&entities.User{},
	)

	dbNumber, _ := strconv.Atoi(di.Env.PrimaryRedis.DB)
	addr := fmt.Sprintf("%s:%s", di.Env.PrimaryRedis.Addr, di.Env.PrimaryRedis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: di.Env.PrimaryRedis.Password,
		DB:       dbNumber,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}

	hub := websocket.NewHub()
	go hub.Run()

	app, err := wire.InitializeApplication(di, db, rdb, hub)
	if err != nil {
		panic(err)
	}

	app.Seeders.RoleSeeder.SeedRoles()

	backgroundEnabled, err := strconv.ParseBool(di.Env.Applications.BACKGROUND_SERVICE_ENABLED)
	if err != nil {
		log.Fatal("Error during checking background service enable")
	}
	if backgroundEnabled {
		app.CronJobs.CronJob.RunCronJob()
	}

	APIServiceEnabled, err := strconv.ParseBool(di.Env.Applications.API_SERVICE_ENABLED)
	if err != nil {
		log.Fatal("Error during checking API service enable")
	}
	if APIServiceEnabled {
		routes.Run(ginEngine, app)
	}

	ginEngine.Run(":8080")
}
