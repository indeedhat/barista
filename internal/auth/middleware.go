package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/indeedhat/barista/internal/server"
)

// IsLoggedInMiddleware will only accept requests from users with a valid login JWT
func IsLoggedInMiddleware(repo Repository) server.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return UserHasPermissionMiddleware(LevelAny, repo)(next)
	}
}

// UserHasPermissionMiddleware checks if the logged in user has a specific permission level
func UserHasPermissionMiddleware(level Level, repo Repository) server.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			jwt := extractJwtFromCookie(r)
			if jwt == "" {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			claims, err := verifyJwt(jwt)
			if err != nil {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			user, err := repo.FindUser(claims.UserId)
			if err != nil {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			if user.JwtKillSwitch != claims.KillSwitch {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			if user.Level&level != level {
				server.WriteResponse(rw, http.StatusForbidden, errors.New("Forbidden"))
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user", user))
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
func AdminOrSelfMiddleware(repo Repository) server.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			jwt := extractJwtFromCookie(r)
			if jwt == "" {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			claims, err := verifyJwt(jwt)
			if err != nil {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			user, err := repo.FindUser(claims.UserId)
			if err != nil {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			if user.JwtKillSwitch != claims.KillSwitch {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			pathId, err := server.PathID(r)
			if err != nil {
				server.WriteResponse(rw, http.StatusUnauthorized, errors.New("Not Authorized"))
				return
			}

			if user.Level&LevelAdmin != LevelAdmin && user.ID != pathId {
				server.WriteResponse(rw, http.StatusForbidden, errors.New("Forbidden"))
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user", user))
			next(rw, r)
		}
	}
}

// IsGuestMiddleware will only accept requests from users withot a valid login JWT
func IsGuestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		jwt := extractJwtFromAuthHeader(r)
		if jwt == "" {
			next(rw, r)
			return
		}

		if _, err := verifyJwt(jwt); err == nil {
			next(rw, r)
			return
		}

		server.WriteResponse(rw, http.StatusForbidden, errors.New("Already logged in"))
	}
}
