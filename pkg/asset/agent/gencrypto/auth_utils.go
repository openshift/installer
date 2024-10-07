package gencrypto

import (
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// WatcherAuthHeaderWriter sets the JWT authorization token.
func WatcherAuthHeaderWriter(token string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		return r.SetHeaderParam("Watcher-Authorization", token)
	})
}

// ParseToken checks if the token string is valid or not and returns JWT token claim.
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Errorf("malformed token claims in url")
	}
	return claims, nil
}

// ParseExpirationFromToken checks if the token is expired or not.
func ParseExpirationFromToken(tokenString string) (time.Time, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return time.Time{}, errors.Errorf("token missing 'exp' claim")
	}
	expTime := time.Unix(int64(exp), 0)
	expiresAt := strfmt.DateTime(expTime)
	expiryTime := time.Time(expiresAt)

	return expiryTime, nil
}
