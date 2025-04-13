package gomail

import (
	"strconv"

	gomailv2 "gopkg.in/gomail.v2"
)

var (
	Dialer *gomailv2.Dialer
	From, Name string
)

func NewDialer(host, port, user, pass string) *gomailv2.Dialer {
	smtpPort, _ := strconv.Atoi(port)

	return gomailv2.NewDialer(host, smtpPort, user, pass)
}
