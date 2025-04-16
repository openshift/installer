// Package dns provides access to the Akamai DNS V2 APIs
//
// See: https://techdocs.akamai.com/edge-dns/reference/edge-dns-api
package dns

import (
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// DNS is the dns api interface
	DNS interface {
		Zones
		TSIGKeys
		Authorities
		Records
		RecordSets
	}

	dns struct {
		session.Session
	}

	// Option defines a DNS option
	Option func(*dns)

	// ClientFunc is a dns client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) DNS
)

// Client returns a new dns Client instance with the specified controller
func Client(sess session.Session, opts ...Option) DNS {
	p := &dns{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Exec overrides the session.Exec to add dns options
func (p *dns) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {

	return p.Session.Exec(r, out, in...)
}
