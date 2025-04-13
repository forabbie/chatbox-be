package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"chatbox/pkg/database"
	"chatbox/pkg/database/postgres"
	"chatbox/pkg/email/gomail"
	"chatbox/pkg/jwt"
	"chatbox/pkg/settings"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
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

	ctx, cancel := context.WithTimeout(context.Background(), settings.Timeout)

	defer cancel()

	pg, err := postgres.Open(
		os.Getenv("PSQL_USER"),
		os.Getenv("PSQL_PASS"),
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_NAME"),
		os.Getenv("PSQL_SLLMODE"),
		ctx,
	)

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

	gomail.Dialer = gomail.NewDialer(
		os.Getenv("GOMAIL_HOST"),
		os.Getenv("GOMAIL_PORT"),
		os.Getenv("GOMAIL_USER"),
		os.Getenv("GOMAIL_PASS"),
	)

	gomail.From, gomail.Name = os.Getenv("GOMAIL_FROM"), os.Getenv("GOMAIL_NAME")

	jwt.AccessTokenKey, _ = jwt.GenerateKey(32)

	if _, ok := os.LookupEnv("JWT_ACCESS_TOKEN_KEY"); ok {
		jwt.AccessTokenKey = os.Getenv("JWT_ACCESS_TOKEN_KEY")
	}

	jwt.RefreshTokenKey, _ = jwt.GenerateKey(32)

	if _, ok := os.LookupEnv("JWT_REFRESH_TOKEN_KEY"); ok {
		jwt.RefreshTokenKey = os.Getenv("JWT_REFRESH_TOKEN_KEY")
	}

	app := New()

	port := os.Getenv("HTTP_PORT")

	if _, ok := os.LookupEnv("PORT"); ok {
		port = os.Getenv("PORT")
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		_ = app.Shutdown()

		log.Fatal(err)
	}
}
