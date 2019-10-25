package middleware

import (
	"github.com/jamesburns-rts/base-go-server/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	isEnabled = true

	privateJwt echo.MiddlewareFunc = func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
)

// IsAuthed should add the JWT check to the echo
func IsAuthed(next echo.HandlerFunc) echo.HandlerFunc {
	return privateJwt(next)
}

func ConfigureAuth(props config.Application) error {
	isEnabled = props.AuthEnabled

	// configure uaa server or something
	return nil
}

func setPublicKeyFromUaa(props config.Application) error {

	//publicKey, err := uaaClient.GetTokenURL()
	publicKey := "my key"

	if key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey)); err == nil {
		privateJwt = middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey:    key,
			SigningMethod: "RS256",
			AuthScheme:    "bearer",
			ErrorHandler: func(e error) error {
				return echo.ErrUnauthorized
			},
		})
	}
	return nil
}
