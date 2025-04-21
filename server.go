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

func setupLogFile(filename string) (*os.File, error) {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(filename, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Open the log file
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func main() {
	// Set up error log file
	errorLogFile, err := setupLogFile(settings.ErrorLogFilename)
	if err != nil {
		log.Fatal("failed to create or open error log file:", err)
	}
	defer errorLogFile.Close()

	// Set up access log file
	accessLogFile, err := setupLogFile(settings.AccessLogFilename)
	if err != nil {
		log.Fatal("failed to create or open access log file:", err)
	}
	defer accessLogFile.Close()

	// Set logging flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)

	// Set output to error log file
	log.SetOutput(errorLogFile)

	// Optional: You could also set up logging to access log file if needed
	// log.SetOutput(io.MultiWriter(errorLogFile, accessLogFile))

	// Configure port
	port := os.Getenv("HTTP_PORT")
	if _, ok := os.LookupEnv("PORT"); ok {
		port = os.Getenv("PORT")
	}

	// Set up PostgreSQL connection
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

	// Set up gomail configuration
	gomailConfig := gomail.Config{
		Host: os.Getenv("GOMAIL_HOST"),
		Port: os.Getenv("GOMAIL_PORT"),
		User: os.Getenv("GOMAIL_USER"),
		Pass: os.Getenv("GOMAIL_PASS"),
	}

	dialer := gomail.NewDialer(gomailConfig)
	email.GomailV2Dialer = dialer
	email.GomailV2From, email.GomailV2Name = os.Getenv("GOMAIL_FROM"), os.Getenv("GOMAIL_NAME")

	// Initialize and run the app
	app := New()

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		_ = app.Shutdown()
		log.Fatal(err)
	}
}
