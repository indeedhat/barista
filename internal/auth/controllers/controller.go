package auth_controllers

import (
	"github.com/indeedhat/barista/internal/auth"
)

const sessionCookie = "bs"

type Controller struct {
	repo auth.Repository
}

func NewController(repo auth.Repository) Controller {
	return Controller{repo}
}
