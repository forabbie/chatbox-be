package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"chatbox/pkg/database"
	"chatbox/pkg/database/postgres"
	"chatbox/pkg/email"
	"chatbox/pkg/email/gomail"
	"chatbox/pkg/settings"
)

func init() {
	if err := godotenv.Load(); err != nil {
		// log.Fatal(err)
		log.Println("No .env file found, relying on system environment variables.")

	}
}

func main() {
	file, err := os.OpenFile(settings.ErrorLogFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)

	log.SetOutput(file)

	file, err = os.OpenFile(settings.AccessLogFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	settings.LoggerConfig.Output = file

	// scheme := os.Getenv("HTTP_SCHEME")

	// host := os.Getenv("HTTP_HOST")

	port := os.Getenv("HTTP_PORT")

	if _, ok := os.LookupEnv("PORT"); ok {
		port = os.Getenv("PORT")
	}

	// baseURL := fmt.Sprintf("%s://%s:%s/", scheme, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	pgConfig := postgres.Config{
		User:     os.Getenv("PSQL_USER"),
		Pass:     os.Getenv("PSQL_PASS"),
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		Name:     os.Getenv("PSQL_NAME"),
		SSLMode:  os.Getenv("PSQL_SLLMODE"),
		TimeZone: "+00:00",
	}

	pg, err := postgres.Open(ctx, pgConfig)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := pg.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := pg.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	database.PostgresMain = pg

	ctx, cancel = context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	gomailConfig := gomail.Config{
		Host: os.Getenv("GOMAIL_HOST"),
		Port: os.Getenv("GOMAIL_PORT"),
		User: os.Getenv("GOMAIL_USER"),
		Pass: os.Getenv("GOMAIL_PASS"),
	}

	dialer := gomail.NewDialer(gomailConfig)

	email.GomailV2Dialer = dialer

	email.GomailV2From, email.GomailV2Name = os.Getenv("GOMAIL_FROM"), os.Getenv("GOMAIL_NAME")

	app := New()

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		_ = app.Shutdown()

		log.Fatal(err)
	}
}
