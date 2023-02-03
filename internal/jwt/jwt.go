package jwt

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"time"
)

type AccessToken struct {
	UserID int `json:"userId"`

	// Issued At as unix timestamp
	Iat int64 `json:"iat"`

	// Expire Time as unix timestamp
	Exp int64 `json:"exp"`
}

// implementation of jwt.Claims

// Valid implementation of jwt.Claims
func (o *AccessToken) Valid() error {
	if o.UserID == 0 {
		return errors.New("invalid wex user id")
	}
	if o.Iat == 0 {
		return errors.New("invalid issued at")
	}
	if o.Exp == 0 {
		return errors.New("invalid expiration")
	}
	return nil
}

type AccessTokenParams struct {
	UserID int `json:"userId"`
}

func NewTokenGenerator(props config.JWT) (func(params AccessTokenParams) (string, error), error) {
	// parse private key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(props.RSAPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	// parse lifespan
	tokenLifespan, err := time.ParseDuration(props.Lifespan)
	if err != nil {
		return nil, fmt.Errorf("invalid token lifespan: %w", err)
	}

	return func(params AccessTokenParams) (string, error) {
		claims := AccessToken{
			UserID: params.UserID,
			Iat:    time.Now().Unix(),
			Exp:    time.Now().Add(tokenLifespan).Unix(),
		}
		return jwt.NewWithClaims(jwt.SigningMethodRS256, &claims).SignedString(privateKey)
	}, nil
}

func ParseAccessToken(tokenString string, publicKey *rsa.PublicKey, leeway time.Duration) (*AccessToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessToken{}, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims := token.Claims.(*AccessToken)

	if time.Until(time.Unix(claims.Exp, 0).Add(leeway)) < 0 {
		return claims, errors.New("token is expired")
	}

	return claims, nil
}

// ParseRSAPublicKeyFromPEM parse bytes into rsa public key
func ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	key, err := base64DecodeIfNecessary(key)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(key)
}

// ParseRSAPrivateKeyFromPEM parse bytes into rsa private key
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	key, err := base64DecodeIfNecessary(key)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(key)
}

func base64DecodeIfNecessary(key []byte) ([]byte, error) {
	key = bytes.TrimSpace(key)

	if len(key) == 0 {
		return nil, errors.New("empty key")
	}
	if key[0] == '-' {
		return key, nil
	}

	// assume base64 encoded
	encoding := base64.StdEncoding
	dest := make([]byte, encoding.DecodedLen(len(key)))
	_, err := encoding.Decode(dest, key)
	if err != nil {
		return key, err
	}
	return dest, nil
}
