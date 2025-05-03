package main

import (
	"github.com/flastors/jwt-auth-golang/internal/app"
)

func main() {
	c := app.NewContext()
	c.Run()
}
