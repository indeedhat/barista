package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/indeedhat/barista/internal/server"
)

const sessionCookie = "bs"

type Controller struct {
	repo Repository
}

func NewController(repo Repository) Controller {
	return Controller{repo}
}

type createSuccessResponse struct {
	ID uint `json:"id"`
}

type createUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Level    uint8  `json:"level" validate:"required"` // TODO: validate level
}

func (c Controller) CreateUser(rw http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	user := User{
		Name:          req.Name,
		Password:      string(hash),
		Level:         Level(req.Level),
		JwtKillSwitch: time.Now().Unix(),
	}

	if err := c.repo.SaveUser(&user); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusCreated, createSuccessResponse{user.ID})
}

type updateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Level uint8  `json:"level" validate:"required"` // TODO: validate level
}

func (c Controller) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	var req updateUserRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := c.repo.FindUser(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("User not found"))
		return
	}

	user.Name = req.Name
	user.Level = Level(req.Level)

	if err := c.repo.SaveUser(user); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}

type changeUserPasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

func (c Controller) ChangePassword(rw http.ResponseWriter, r *http.Request) {
	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	var req changeUserPasswordRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := c.repo.FindUser(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("User not found"))
		return
	}

	if err := c.repo.UpdateUserPassword(user, req.Password); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}

type loginRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login handles user login attempts
func (c Controller) Login(rw http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := server.UnmarshalBody(r, &req); err != nil {
		server.WriteResponse(rw, http.StatusBadRequest, nil)
		return
	}

	if err := server.ValidateRequest(req); err != nil {
		server.WriteResponse(rw, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := c.repo.FindUserByLogin(req.Name, req.Password)
	if user == nil {
		server.WriteResponse(rw, http.StatusUnauthorized, "login failed")
	}

	jwt, err := GenerateUserJwt(user.ID, user.Name, uint8(user.Level), user.JwtKillSwitch)
	if err != nil {
		server.WriteResponse(
			rw,
			http.StatusInternalServerError,
			errors.New("Failed to process login"),
		)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     sessionCookie,
		Value:    jwt,
		HttpOnly: true,
		Domain:   r.URL.Host,
		Path:     "/",
		MaxAge:   86400 * 30,
	})
	server.WriteResponse(rw, http.StatusNoContent, nil)
}

func (c Controller) Logout(rw http.ResponseWriter, r *http.Request) {
	http.SetCookie(rw, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		HttpOnly: true,
		Domain:   r.URL.Host,
		Path:     "/",
		MaxAge:   0,
	})
	server.WriteResponse(rw, http.StatusNoContent, nil)
}

// GetLoggedInUser returns the json representation of the logged in user
func (c Controller) GetLoggedInUser(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	server.WriteResponse(rw, http.StatusOK, user)
}

// ForceLogoutUser resets the users JwtKillSwitch field invalidating all existing logins
func (c Controller) ForceLogoutUser(rw http.ResponseWriter, r *http.Request) {
	id, err := server.PathID(r)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, nil)
		return
	}

	user, err := c.repo.FindUser(id)
	if err != nil {
		server.WriteResponse(rw, http.StatusNotFound, errors.New("User not found"))
		return
	}

	user.JwtKillSwitch = time.Now().Unix()

	if err := c.repo.SaveUser(user); err != nil {
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
		return
	}

	server.WriteResponse(rw, http.StatusNoContent, nil)
}
