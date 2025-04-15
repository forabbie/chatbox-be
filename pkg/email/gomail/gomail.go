package gomail

import (
	"strconv"

	gomailv2 "gopkg.in/gomail.v2"
)

type Config struct {
	Host, Port, User, Pass string
}

func NewDialer(config Config) *gomailv2.Dialer {
	smtpPort, _ := strconv.Atoi(config.Port)

	return gomailv2.NewDialer(config.Host, smtpPort, config.User, config.Pass)
}
