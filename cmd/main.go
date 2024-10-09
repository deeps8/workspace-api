package main

import (
	"log"
	"os"
	"strings"
	"work-space-backend/database"
	"work-space-backend/handlers"
	cronjob "work-space-backend/handlers/cron-job"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 1. creating new instance of echo(server)
	app := echo.New()
	// 2. adding logging middleware
	app.Use(middleware.Logger())

	// 3. load .env
	envErr := godotenv.Load()
	if envErr != nil {
		log.Printf("Error loading env %+v", envErr)
	}
	or := os.Getenv("ORIGIN")
	origins := strings.Split(or, ",")
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowCredentials: true, AllowOrigins: origins}))
	// 4. connect the database
	database.Db = database.NewDBconn()

	// close the db connection before server ends
	defer database.Db.Close()

	// 5. handle all the routes
	api := app.Group("/api")
	handlers.InitHandler(api)

	// getting the port from env and starting the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS"))
	if err != nil {
		panic(err)
	}

	database.Rdb = redis.NewClient(opt)

	cronjob.SyncDb()
	app.Logger.Fatal(app.Start("0.0.0.0:" + port))

}
