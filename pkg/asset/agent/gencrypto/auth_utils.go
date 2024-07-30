package gencrypto

import (
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// UserAuthHeaderWriter sets the JWT authorization token.
func UserAuthHeaderWriter(token string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		return r.SetHeaderParam("Authorization", token)
	})
}

// ParseExpirationFromToken checks if the token is expired or not.
func ParseExpirationFromToken(tokenString string) (time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return time.Time{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, errors.Errorf("malformed token claims in url")
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
