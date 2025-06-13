package auth

import (
	"context"
	"net/http"

	"github.com/indeedhat/barista/internal/server"
)

type RouteType uint8

const (
	UI RouteType = 1 << iota
	API
)

// IsLoggedInMiddleware will only accept requests from users with a valid login JWT
func IsLoggedInMiddleware(rt RouteType, repo Repository) server.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return UserHasPermissionMiddleware(rt, LevelAny, repo)(next)
	}
}

// UserHasPermissionMiddleware checks if the logged in user has a specific permission level
func UserHasPermissionMiddleware(rt RouteType, level Level, repo Repository) server.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			user := parseJwt(r, repo)
			if user == nil {
				redirectOrHeader(rw, r, http.StatusUnauthorized, rt, "/login")
				return
			}

			if user.Level&level != user.Level {
				redirectOrHeader(rw, r, http.StatusForbidden, rt, "/")
				return
			}

			r = r.WithContext(r.Context().(server.Context).WithValue("user", user))
			next(rw, r)
		}
	}
}

// AdminOrSelfMiddleware checks if the logged in user is eith an admin user or is performing the
// operation on itself
//
// This is a fairly sketch check that is designed to be used for user operations on self but the
// check only verifies that the {id} in the path matches that of the logged in user so if applied
// to an enpoint that is not a user op then it could allow unintended permissions
func AdminOrSelfMiddleware(rt RouteType, repo Repository) server.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			user := parseJwt(r, repo)
			if user == nil {
				redirectOrHeader(rw, r, http.StatusUnauthorized, rt, "/login")
				return
			}

			pathId, err := server.PathID(r)
			if err != nil {
				redirectOrHeader(rw, r, http.StatusUnauthorized, rt, "/login")
				return
			}

			if user.Level&LevelAdmin != LevelAdmin && user.ID != pathId {
				redirectOrHeader(rw, r, http.StatusForbidden, rt, "/login")
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user", user))
			next(rw, r)
		}
	}
}

// IsGuestMiddleware will only accept requests from users withot a valid login JWT
func IsGuestMiddleware(rt RouteType, repo Repository) server.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			user := parseJwt(r, repo)
			if user != nil {
				redirectOrHeader(rw, r, http.StatusForbidden, rt, "/")
				return
			}

			next(rw, r)
		}
	}
}

func parseJwt(r *http.Request, repo Repository) *User {
	jwt := extractJwtFromCookie(r)
	if jwt == "" {
		return nil
	}

	claims, err := verifyJwt(jwt)
	if err != nil {
		return nil
	}

	user, err := repo.FindUser(claims.UserId)
	if err != nil {
		return nil
	}

	if user.JwtKillSwitch != claims.KillSwitch {
		return nil
	}

	return user
}

func redirectOrHeader(rw http.ResponseWriter, r *http.Request, code int, rt RouteType, url string) {
	switch rt {
	case API:
		server.WriteResponse(rw, code, nil)
	case UI:
		http.Redirect(rw, r, url, http.StatusSeeOther)
	default:
		server.WriteResponse(rw, http.StatusInternalServerError, nil)
	}
}
