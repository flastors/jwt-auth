package empty

import (
	"github.com/flastors/jwt-auth-golang/internal/core/user"
)

type EmptyMailer struct {
}

func NewUserEmail() user.UserEmail {
	return &EmptyMailer{}
}

func (m *EmptyMailer) Send(to, subject, body string) error {
	return nil
}
