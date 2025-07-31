package auth_controllers

import (
	"github.com/indeedhat/barista/internal/auth"
)

type Controller struct {
	repo auth.Repository
}

func New(repo auth.Repository) Controller {
	return Controller{repo}
}
