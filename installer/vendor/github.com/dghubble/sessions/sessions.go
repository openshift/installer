package sessions

import (
	"net/http"
)

const (
	defaultMaxAge = 3600 * 24 * 7 // 1 week
)

// Config is the set of session cookie properties.
type Config struct {
	// cookie domain/path scope (leave zeroed for requested resource scope)
	Domain string
	Path   string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int
	// browser should prohibit non-HTTP (i.e. javascript) cookie access
	HTTPOnly bool
	// cookie may only be transferred over HTTPS
	Secure bool
}

// Session represents Values state which  a named bundle of maintained web state
// stores web session state
type Session struct {
	name   string  // session cookie name
	Config *Config // session cookie config
	store  Store   // session store
	Values map[string]interface{}
}

// NewSession returns a new Session.
func NewSession(store Store, name string) *Session {
	return &Session{
		store:  store,
		name:   name,
		Values: make(map[string]interface{}),
	}
}

// Name returns the name of the session.
func (s *Session) Name() string {
	return s.name
}

// Save adds or updates the session. Identical to calling
// store.Save(w, session).
func (s *Session) Save(w http.ResponseWriter) error {
	return s.store.Save(w, s)
}

// Destroy destroys the session. Identical to calling
// store.Destroy(w, session.name).
func (s *Session) Destroy(w http.ResponseWriter) {
	s.store.Destroy(w, s.name)
}
