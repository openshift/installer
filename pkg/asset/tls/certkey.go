package tls

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

const (
	// KeyIndex is the index into a CertKey asset's contents at which the key
	// can be found.
	KeyIndex = 0

	// CertIndex is the index into a CertKey asset's contents at which the
	// certificate can be found.
	CertIndex = 1
)

// CertKey contains the private key and the cert that's
// signed by the parent CA.
type CertKey struct {
	installConfig asset.Asset

	// Common fields.
	Subject      pkix.Name
	KeyUsages    x509.KeyUsage
	ExtKeyUsages []x509.ExtKeyUsage
	Validity     time.Duration
	KeyFileName  string
	CertFileName string
	ParentCA     asset.Asset

	IsCA         bool
	AppendParent bool // Whether append the parent CA in the cert.

	// Some certs might need to set Subject, DNSNames and IPAddresses.
	GenDNSNames    func(*types.InstallConfig) ([]string, error)
	GenIPAddresses func(*types.InstallConfig) ([]net.IP, error)
	GenSubject     func(*types.InstallConfig) (pkix.Name, error)
}

var _ asset.Asset = (*CertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (c *CertKey) Dependencies() []asset.Asset {
	parents := []asset.Asset{c.ParentCA}

	// Require the InstallConfig if we need additional info from install config.
	if c.GenDNSNames != nil || c.GenIPAddresses != nil || c.GenSubject != nil {
		parents = append(parents, c.installConfig)
	}

	return parents
}

// Generate generates the cert/key pair based on its dependencies.
func (c *CertKey) Generate(parents map[asset.Asset]*asset.State) (*asset.State, error) {
	cfg := &CertCfg{
		Subject:      c.Subject,
		KeyUsages:    c.KeyUsages,
		ExtKeyUsages: c.ExtKeyUsages,
		Validity:     c.Validity,
		IsCA:         c.IsCA,
	}

	if c.GenSubject != nil || c.GenDNSNames != nil || c.GenIPAddresses != nil {
		state, ok := parents[c.installConfig]
		if !ok {
			return nil, fmt.Errorf("failed to get install config state in the parent asset states")
		}

		var installConfig types.InstallConfig
		if err := yaml.Unmarshal(state.Contents[0].Data, &installConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal install config: %v", err)
		}

		var err error
		if c.GenSubject != nil {
			cfg.Subject, err = c.GenSubject(&installConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to generate Subject: %v", err)
			}
		}
		if c.GenDNSNames != nil {
			cfg.DNSNames, err = c.GenDNSNames(&installConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to generate DNSNames: %v", err)
			}
		}
		if c.GenIPAddresses != nil {
			cfg.IPAddresses, err = c.GenIPAddresses(&installConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to generate IPAddresses: %v", err)
			}
		}
	}

	var key *rsa.PrivateKey
	var crt *x509.Certificate
	var err error

	state, ok := parents[c.ParentCA]
	if !ok {
		return nil, fmt.Errorf("failed to get parent CA %v in the parent asset states", c.ParentCA)
	}

	caKey, caCert, err := parseCAFromAssetState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA from asset: %v", err)
	}

	key, crt, err = GenerateCert(caKey, caCert, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cert/key pair: %v", err)
	}

	keyData := []byte(PrivateKeyToPem(key))
	certData := []byte(CertToPem(crt))
	if c.AppendParent {
		certData = append(certData, '\n')
		certData = append(certData, []byte(CertToPem(caCert))...)
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: assetFilePath(c.KeyFileName),
				Data: keyData,
			},
			{
				Name: assetFilePath(c.CertFileName),
				Data: certData,
			},
		},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (c *CertKey) Name() string {
	return fmt.Sprintf("Certificate (%s)", c.Subject.CommonName)
}

func parseCAFromAssetState(ca *asset.State) (*rsa.PrivateKey, *x509.Certificate, error) {
	var key *rsa.PrivateKey
	var cert *x509.Certificate
	var err error

	if len(ca.Contents) != 2 {
		return nil, nil, fmt.Errorf("expect key and cert in the contents of CA, got: %v", ca)
	}

	for _, c := range ca.Contents {
		switch filepath.Ext(c.Name) {
		case ".key":
			key, err = PemToPrivateKey(c.Data)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse rsa private key: %v", err)
			}
		case ".crt":
			cert, err = PemToCertificate(c.Data)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse x509 certificate: %v", err)
			}
		default:
			return nil, nil, fmt.Errorf("unexpected content name: %v", c.Name)
		}
	}

	return key, cert, nil
}
