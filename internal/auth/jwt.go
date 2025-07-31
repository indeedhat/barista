package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/indeedhat/barista/internal/cookie"
)

var ErrInvalidJWT = errors.New("Invalid jwt")

type Claims struct {
	jwt.RegisteredClaims

	Name       string `json:"nme"`
	UserId     uint   `json:"uid"`
	Level      uint8  `json:"lvl"`
	KillSwitch int64  `json:"kil"`
}

// GenerateJWT will generate a new JWT for the given account model
func GenerateJWT(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(envJwtSecret.Get()))
}

// GenerateUserJwt genertes a new JWT specifically for a user login session
func GenerateUserJwt(id uint, name string, level uint8, killSwitch int64) (string, error) {
	return GenerateJWT(Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID: strconv.Itoa(int(time.Now().Unix())),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(envJwtTTl.Get()) * time.Second,
			)),
		},
		Name:       name,
		UserId:     id,
		Level:      level,
		KillSwitch: killSwitch,
	})
}

// extractJwtFromAuthHeader will verify that the Authorization header both exists and is in the
// Bearer format, if so it will extract the token (hopefully this should be a valid JWT)
func extractJwtFromAuthHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

// extractJwtFromAuthHeader will verify that the Authorization header both exists and is in the
// Bearer format, if so it will extract the token (hopefully this should be a valid JWT)
func extractJwtFromCookie(r *http.Request) string {
	c, err := r.Cookie(cookie.SessionKey)
	if err != nil {
		return ""
	}

	return c.Value
}

// VerifyJwt will check that the JWT is both a jwt and valid
func verifyJwt(jwtString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(jwtString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != "HS256" {
			return nil, ErrInvalidJWT
		}

		return []byte(envJwtSecret.Get()), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}
