package auth

import (
	"github.com/indeedhat/dotenv"
)

const (
	// Auth
	envJwtSecret dotenv.String = "JWT_SECRET"
	// Time since jwt generation that will cause the jwt to be refreshed
	envJwtRefreshAge dotenv.Int = "JWT_REFRESH_AGE"
	envJwtTTl        dotenv.Int = "JWT_TTL"
)

const (
	envRootUsername dotenv.String = "ROOT_USERNAME"
	envRootPassword dotenv.String = "ROOT_PASSWORD"

	defaultRootUsername = "admin"
	defaultRootPassword = "admin"
)

const EnvEnableRegister dotenv.Bool = "ENABLE_REGISTER"
