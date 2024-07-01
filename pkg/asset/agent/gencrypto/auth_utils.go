package gencrypto

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// UserAuthHeaderWriter sets the JWT authorization token.
func UserAuthHeaderWriter(token string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		return r.SetHeaderParam("Authorization", token)
	})
}
