package vsphere

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/cns"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/pbm"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// Session manages a vCenter connection lifecycle including SOAP, REST, and
// auxiliary clients. It provides a single entry point for creating and
// reusing vSphere client sessions with proper caching and cleanup.
type Session struct {
	server   string
	govmomi  *govmomi.Client
	vim25    *vim25.Client
	rest     *rest.Client
	cns      *cns.Client
	pbm      *pbm.Client
	finder   *find.Finder
	logger   logrus.FieldLogger
	once     sync.Once
	closeErr error
}

// SessionOptions configures a vSphere session.
type SessionOptions struct {
	// Timeout for the initial connection. Defaults to 60s if zero.
	Timeout time.Duration
	// Insecure skips TLS certificate verification. Defaults to false.
	Insecure bool
	// Logger for structured logging. If nil, uses logrus.StandardLogger().
	Logger logrus.FieldLogger
}

// SessionOption configures a session via functional options.
type SessionOption func(*SessionOptions)

// WithTimeout sets the connection timeout.
func WithTimeout(timeout time.Duration) SessionOption {
	return func(o *SessionOptions) { o.Timeout = timeout }
}

// WithInsecure skips TLS certificate verification.
func WithInsecure(insecure bool) SessionOption {
	return func(o *SessionOptions) { o.Insecure = insecure }
}

// WithLogger sets the logger for structured logging.
func WithLogger(logger logrus.FieldLogger) SessionOption {
	return func(o *SessionOptions) { o.Logger = logger }
}

// applyDefaults sets reasonable defaults for unset options.
func (o *SessionOptions) applyDefaults() {
	if o.Timeout == 0 {
		o.Timeout = 60 * time.Second
	}
	if o.Logger == nil {
		o.Logger = logrus.StandardLogger()
	}
}

// SessionCacheKey returns a deterministic key for session caching based on
// server, username, and password.
func SessionCacheKey(server, username, password string) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", server, username, password)))
	return fmt.Sprintf("%x", hash[:16])
}

var (
	sessionsMu sync.Mutex
	sessions   = make(map[string]*Session)
)

// redactHostname returns a short deterministic hash of the hostname, suitable
// for logging. The plain hostname is never exposed in logs.
func redactHostname(hostname string) string {
	h := sha256.Sum256([]byte(hostname))
	return fmt.Sprintf("%x", h[:16])
}

// NewSession creates or retrieves a cached vSphere session for the given
// server, username, and password. Returns the session and a cleanup function.
//
// The cleanup function must be called when the session is no longer needed.
// It is idempotent and safe to call multiple times.
//
// Options allow configuring timeout, TLS, and other settings.
func NewSession(ctx context.Context, server, username, password string, opts ...SessionOption) (*Session, func(), error) {
	options := makeSessionOptions(opts...)
	options.applyDefaults()

	if server == "" || username == "" || password == "" {
		return nil, nil, errors.New("vSphere server, username, and password are required")
	}

	key := SessionCacheKey(server, username, password)

	sessionsMu.Lock()
	if sess, ok := sessions[key]; ok {
		if sess.Error() != nil {
			// Cached session is closed — evict it so a fresh one is created.
			delete(sessions, key)
			sessionsMu.Unlock()
			// Fall through to create a new session.
		} else {
			sessionsMu.Unlock()
			return sess, func() { sess.Close() }, nil
		}
	}
	sessionsMu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	sess, err := newSessionFromConn(ctx, server, username, password, options)
	if err != nil {
		return nil, nil, err
	}

	sessionsMu.Lock()
	// Double check after acquiring the lock (another goroutine may have created
	// the session in the meantime).
	if s, ok := sessions[key]; ok {
		// A cached session arrived in the meantime — close the one we just
		// created to avoid leaking an authenticated client.
		sess.Close()
		sessionsMu.Unlock()
		return s, func() { s.Close() }, nil
	}
	sessions[key] = sess
	sessionsMu.Unlock()

	return sess, func() { sess.Close() }, nil
}

func makeSessionOptions(opts ...SessionOption) *SessionOptions {
	o := &SessionOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Close releases all resources held by the session. It is safe to call
// multiple times. If any cleanup operation fails, the first error is
// stored and returned via Error().
func (s *Session) Close() {
	s.once.Do(func() {
		var firstErr error

		s.logger.WithField("vcenter", redactHostname(s.server)).Debug("Closing vSphere session")

		if s.rest != nil {
			if err := s.rest.Logout(context.Background()); err != nil {
				s.logger.WithError(err).Warn("Failed to logout REST client")
				if firstErr == nil {
					firstErr = err
				}
			}
		}

		if s.govmomi != nil {
			if err := s.govmomi.Logout(context.Background()); err != nil {
				s.logger.WithError(err).Warn("Failed to logout SOAP client")
				if firstErr == nil {
					firstErr = err
				}
			}
		}

		s.closeErr = firstErr
		s.logger.WithField("vcenter", redactHostname(s.server)).Debug("vSphere session closed")
	})
}

// Error returns any error that occurred during Close.
func (s *Session) Error() error {
	return s.closeErr
}

func newSessionFromConn(ctx context.Context, server, username, password string, opts *SessionOptions) (*Session, error) {
	s := &Session{
		server: server,
		logger: opts.Logger,
	}

	u, err := soap.ParseURL(server)
	if err != nil {
		return nil, errors.Wrap(err, "parse vCenter URL")
	}
	u.User = url.UserPassword(username, password)

	c, err := govmomi.NewClient(ctx, u, opts.Insecure)
	if err != nil {
		return nil, errors.Wrap(err, "create SOAP client")
	}
	s.govmomi = c
	s.vim25 = c.Client
	s.finder = find.NewFinder(s.vim25, true)

	restClient := rest.NewClient(c.Client)
	if err := restClient.Login(ctx, u.User); err != nil {
		// Log out SOAP client before returning error — this replaces the
		// buggy pattern in CreateVSphereClients where the logout error
		// was swallowed and could mask the original error.
		if logoutErr := c.Logout(context.Background()); logoutErr != nil {
			opts.Logger.WithError(logoutErr).Warn("Failed to logout SOAP client after REST login failure")
		}
		return nil, errors.Wrap(err, "REST client login")
	}
	s.rest = restClient

	if cnsClient, err := cns.NewClient(context.Background(), s.vim25); err == nil {
		s.cns = cnsClient
	} else {
		opts.Logger.WithError(err).Debug("CNS client creation failed (optional)")
	}

	if pbmClient, err := pbm.NewClient(context.Background(), s.vim25); err == nil {
		s.pbm = pbmClient
	} else {
		opts.Logger.WithError(err).Debug("PBM client creation failed (optional)")
	}

	s.logger.WithField("vcenter", redactHostname(server)).Info("vSphere session created")
	return s, nil
}

// Vim25Client returns the underlying SOAP client for general vSphere API calls.
func (s *Session) Vim25Client() *vim25.Client {
	return s.vim25
}

// RestClient returns the REST client for vSphere tag and policy operations.
func (s *Session) RestClient() *rest.Client {
	return s.rest
}

// CNSClient returns the CNS (vSAN) client for storage volume operations.
func (s *Session) CNSClient() *cns.Client {
	return s.cns
}

// PBMClient returns the PBM (storage policy) client.
func (s *Session) PBMClient() *pbm.Client {
	return s.pbm
}

// Finder returns a finder for locating vSphere inventory objects.
func (s *Session) Finder() *find.Finder {
	return s.finder
}

// Server returns the vCenter server URL.
func (s *Session) Server() string {
	return s.server
}
