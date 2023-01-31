package ovirtclient

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// TLSProvider creates a TLS configuration for use by the oVirt client.
type TLSProvider interface {
	// CreateTLSConfig returns a working TLS configuration for the client, or an error if the configuration cannot be
	// created.
	CreateTLSConfig() (*tls.Config, error)
}

// BuildableTLSProvider is a TLSProvider that allows adding configuration options.
type BuildableTLSProvider interface {
	TLSProvider

	// Insecure disables CA certificate verification. This cannot be chained as any further option doesn't make any
	// sense.
	Insecure() TLSProvider

	// CACertsFromMemory adds one or more CA certificate from an in-memory byte slice containing PEM-encoded
	// certificates. This function can be called multiple times to add multiple certificates. The function will fail if
	// no certificates have been added.
	CACertsFromMemory(caCert []byte) BuildableTLSProvider

	// CACertsFromFile adds certificates from a single file. The certificate must be in PEM format. This function can
	// be called multiple times to add multiple files.
	CACertsFromFile(file string) BuildableTLSProvider

	// CACertsFromDir adds all PEM-encoded certificates from a directory. If one or more patterns are passed, the
	// files will only be added if the files match at least one of the patterns. The certificate will fail if one or
	// more matching files don't contain a valid certificate.
	CACertsFromDir(dir string, patterns ...*regexp.Regexp) BuildableTLSProvider

	// CACertsFromSystem adds the system certificate store. This may fail because the certificate store is not available
	// or not supported on the platform.
	CACertsFromSystem() BuildableTLSProvider

	// CACertsFromCertPool sets a certificate pool to use as a source for certificates. This is incompatible with  the
	// CACertsFromSystem call as both create a certificate pool. This function must not be called twice.
	CACertsFromCertPool(*x509.CertPool) BuildableTLSProvider
}

// TLS creates a BuildableTLSProvider that can be used to easily add trusted CA certificates and generally follows best
// practices in terms of TLS settings.
func TLS() BuildableTLSProvider {
	return &standardTLSProvider{
		lock: &sync.Mutex{},
	}
}

type standardTLSProvider struct {
	lock        *sync.Mutex
	insecure    bool
	caCerts     [][]byte
	files       []string
	directories []standardTLSProviderDirectory
	certPool    *x509.CertPool
	system      bool
	configured  bool
}

type standardTLSProviderDirectory struct {
	dir      string
	patterns []*regexp.Regexp
}

func (s *standardTLSProvider) Insecure() TLSProvider {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.configured = true
	s.insecure = true
	return s
}

func (s *standardTLSProvider) CACertsFromCertPool(certPool *x509.CertPool) BuildableTLSProvider {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.certPool != nil {
		panic(newError(EConflict, "the CACertsFromCertPool function has been called twice"))
	}
	s.configured = true
	s.certPool = certPool
	return s
}

func (s *standardTLSProvider) CACertsFromMemory(caCert []byte) BuildableTLSProvider {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.configured = true
	s.caCerts = append(s.caCerts, caCert)
	return s
}

func (s *standardTLSProvider) CACertsFromFile(file string) BuildableTLSProvider {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.configured = true
	s.files = append(s.files, file)
	return s
}

func (s *standardTLSProvider) CACertsFromDir(dir string, patterns ...*regexp.Regexp) BuildableTLSProvider {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.configured = true
	s.directories = append(s.directories, standardTLSProviderDirectory{
		dir,
		patterns,
	})
	return s
}

func (s *standardTLSProvider) CACertsFromSystem() BuildableTLSProvider {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.configured = true
	s.system = true
	return s
}

func (s *standardTLSProvider) CreateTLSConfig() (*tls.Config, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.configured {
		return nil, newError(
			ETLSError,
			"TLS not configured (Did you forget to call certificate configuration options on the TLS provider?)",
		)
	}
	if s.insecure {
		return &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		}, nil
	}
	tlsConfig := &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   false,
		SessionTicketKey:         [32]byte{},
		ClientSessionCache:       nil,
		// Based on Mozilla intermediate compatibility:
		// https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28recommended.29
		MinVersion: tls.VersionTLS12,
		MaxVersion: 0,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256, tls.CurveP384,
		},
		DynamicRecordSizingDisabled: false,
		Renegotiation:               0,
		KeyLogWriter:                nil,
	}

	certPool := s.certPool
	if certPool == nil {
		var err error
		if certPool, err = s.createCertPool(); err != nil {
			return nil, err
		}
	}
	if err := s.addCertsFromMemory(certPool); err != nil {
		return nil, err
	}
	if err := s.addCertsFromFile(certPool); err != nil {
		return nil, err
	}
	if err := s.addCertsFromDir(certPool); err != nil {
		return nil, err
	}
	tlsConfig.RootCAs = certPool
	return tlsConfig, nil
}

func (s *standardTLSProvider) addCertsFromDir(certPool *x509.CertPool) error {
	for _, dir := range s.directories {
		files, err := os.ReadDir(dir.dir)
		if err != nil {
			return wrap(
				err,
				ELocalIO,
				"failed to list contents of %s directory",
				dir.dir,
			)
		}
		for _, info := range files {
			if info.IsDir() {
				continue
			}
			if len(dir.patterns) > 0 {
				matches := false
				for _, pattern := range dir.patterns {
					if pattern.MatchString(info.Name()) {
						matches = true
						break
					}
				}
				if !matches {
					continue
				}
			}
			fullPath := filepath.Join(dir.dir, info.Name())
			data, err := os.ReadFile(fullPath) //nolint:gosec
			if err != nil {
				return wrap(
					err,
					EFileReadFailed,
					"failed to read certificate file: %s (%w)",
					fullPath,
				)
			}
			if !certPool.AppendCertsFromPEM(data) {
				return newError(
					ETLSError,
					"failed to add certificate from file: %s (certificate not in PEM format?)",
					fullPath,
				)
			}
		}
	}
	return nil
}

func (s *standardTLSProvider) addCertsFromFile(certPool *x509.CertPool) error {
	for _, file := range s.files {
		pemData, err := os.ReadFile(file) //nolint:gosec
		if err != nil {
			return wrap(err, EFileReadFailed, "failed to read CA certificate from file %s", file)
		}
		if ok := certPool.AppendCertsFromPEM(pemData); !ok {
			return newError(
				ETLSError,
				"the provided CA certificate is not a valid certificate in PEM format in file %s",
				file,
			)
		}
	}
	return nil
}

func (s *standardTLSProvider) addCertsFromMemory(certPool *x509.CertPool) error {
	for i, caCert := range s.caCerts {
		if ok := certPool.AppendCertsFromPEM(caCert); !ok {
			return newError(
				EBadArgument,
				"the provided CA certificate number #%d is not a valid certificate in PEM format",
				i,
			)
		}
	}
	return nil
}

func (s *standardTLSProvider) createCertPool() (*x509.CertPool, error) {
	if s.system && s.certPool != nil {
		return nil, newError(ETLSError, "both system and cert pool are specified, these options are incompatible")
	}
	if s.certPool != nil {
		return s.certPool, nil
	}
	if !s.system {
		return x509.NewCertPool(), nil
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		// This is the case on Windows before go 1.18 where the system certificate pool is not available.
		// See https://github.com/golang/go/issues/16736
		return nil, wrap(err, ETLSError, "system cert pool not available")
	}
	return certPool, nil
}
