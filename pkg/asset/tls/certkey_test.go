package tls

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	libcrypto "github.com/openshift/library-go/pkg/crypto"
)

func TestSignedCertKeyGenerate(t *testing.T) {
	tests := []struct {
		name         string
		certCfg      *CertCfg
		filenameBase string
		certFileName string
		appendParent AppendParentChoice
		errString    string
	}{
		{
			name: "simple ca",
			certCfg: &CertCfg{
				Subject:   pkix.Name{CommonName: "test0-ca", OrganizationalUnit: []string{"openshift"}},
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Validity:  ValidityTenYears(),
				DNSNames:  []string{"test.openshift.io"},
			},
			filenameBase: "test0-ca",
			appendParent: DoNotAppendParent,
		},
		{
			name: "more complicated ca",
			certCfg: &CertCfg{
				Subject:     pkix.Name{CommonName: "test1-ca", OrganizationalUnit: []string{"openshift"}},
				KeyUsages:   x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Validity:    ValidityTenYears(),
				DNSNames:    []string{"test.openshift.io"},
				IPAddresses: []net.IP{net.ParseIP("10.0.0.1")},
			},
			filenameBase: "test1-ca",
			appendParent: AppendParent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCA := &SelfSignedCertKey{}
			rootCACfg := &CertCfg{
				Subject:   pkix.Name{CommonName: "test-root-ca", OrganizationalUnit: []string{"openshift"}},
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Validity:  ValidityTenYears(),
				IsCA:      true,
			}
			err := rootCA.Generate(context.Background(), rootCACfg, "test-root-ca", nil)
			assert.NoError(t, err, "failed to generate root CA")

			certKey := &SignedCertKey{}
			err = certKey.Generate(context.Background(), tt.certCfg, rootCA, tt.filenameBase, tt.appendParent)
			if err != nil {
				assert.EqualErrorf(t, err, tt.errString, tt.name)
				return
			} else if tt.errString != "" {
				t.Errorf("expect error %v, saw nil", err)
			}

			actualFiles := certKey.Files()

			assert.Equal(t, 2, len(actualFiles), "unexpected number of files")
			assert.Equal(t, assetFilePath(tt.filenameBase+".key"), actualFiles[0].Filename, "unexpected key file name")
			assert.Equal(t, assetFilePath(tt.filenameBase+".crt"), actualFiles[1].Filename, "unexpected cert file name")

			assert.Equal(t, certKey.Key(), actualFiles[0].Data, "key file data does not match key")
			assert.Equal(t, certKey.Cert(), actualFiles[1].Data, "cert file does not match cert")

			// Briefly check the certs.
			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(certKey.Cert()) {
				t.Error("failed to append certs from PEM")
			}

			opts := x509.VerifyOptions{
				Roots:   certPool,
				DNSName: tt.certCfg.Subject.CommonName,
			}
			if tt.certCfg.DNSNames != nil {
				opts.DNSName = "test.openshift.io"
			}

			cert, err := PemToCertificate(certKey.Cert())
			assert.NoError(t, err, tt.name)

			_, err = cert.Verify(opts)
			assert.NoError(t, err, tt.name)
		})
	}
}

func TestSelfSignedCertKeyGenerateWithKeyGen(t *testing.T) {
	cases := []struct {
		name            string
		keyGen          libcrypto.KeyPairGenerator
		expectKeyType   interface{}
		expectPubKeyAlg x509.PublicKeyAlgorithm
	}{
		{
			name:            "RSA 4096",
			keyGen:          libcrypto.RSAKeyPairGenerator{Bits: 4096},
			expectKeyType:   &rsa.PrivateKey{},
			expectPubKeyAlg: x509.RSA,
		},
		{
			name:            "ECDSA P384",
			keyGen:          libcrypto.ECDSAKeyPairGenerator{Curve: libcrypto.P384},
			expectKeyType:   &ecdsa.PrivateKey{},
			expectPubKeyAlg: x509.ECDSA,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &CertCfg{
				Subject:  pkix.Name{CommonName: "test-pki-ca", OrganizationalUnit: []string{"openshift"}},
				Validity: ValidityTenYears(),
				IsCA:     true,
			}

			ca := &SelfSignedCertKey{}
			err := ca.Generate(context.Background(), cfg, "test-pki-ca", tc.keyGen)
			assert.NoError(t, err)

			key, err := PemToPrivateKey(ca.Key())
			assert.NoError(t, err)
			assert.IsType(t, tc.expectKeyType, key)

			cert, err := PemToCertificate(ca.Cert())
			assert.NoError(t, err)
			assert.Equal(t, tc.expectPubKeyAlg, cert.PublicKeyAlgorithm)
			assert.True(t, cert.IsCA)
		})
	}
}

func TestCrossAlgorithmCertificateSigning(t *testing.T) {
	// Generate ECDSA P384 CA
	rootCA := &SelfSignedCertKey{}
	rootCACfg := &CertCfg{
		Subject:  pkix.Name{CommonName: "ecdsa-ca", OrganizationalUnit: []string{"openshift"}},
		Validity: ValidityTenYears(),
		IsCA:     true,
	}
	err := rootCA.Generate(context.Background(), rootCACfg, "ecdsa-ca", libcrypto.ECDSAKeyPairGenerator{Curve: libcrypto.P384})
	assert.NoError(t, err)

	// Verify CA key is ECDSA
	caKey, err := PemToPrivateKey(rootCA.Key())
	assert.NoError(t, err)
	assert.IsType(t, &ecdsa.PrivateKey{}, caKey)

	// Generate RSA leaf signed by ECDSA CA
	leafCfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "leaf-cert", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		Validity:  ValidityTenYears(),
		DNSNames:  []string{"test.openshift.io"},
	}
	certKey := &SignedCertKey{}
	err = certKey.Generate(context.Background(), leafCfg, rootCA, "cross-algo-leaf", DoNotAppendParent)
	assert.NoError(t, err)

	// Verify leaf key is RSA (SignedCertKey always generates RSA leaf keys)
	leafKey, err := PemToPrivateKey(certKey.Key())
	assert.NoError(t, err)
	assert.IsType(t, &rsa.PrivateKey{}, leafKey)

	// Verify the leaf cert was signed by the ECDSA CA
	leafCert, err := PemToCertificate(certKey.Cert())
	assert.NoError(t, err)
	assert.Equal(t, x509.ECDSAWithSHA384, leafCert.SignatureAlgorithm)

	// Verify cert chain: leaf validates against CA
	caCert, err := PemToCertificate(rootCA.Cert())
	assert.NoError(t, err)
	certPool := x509.NewCertPool()
	certPool.AddCert(caCert)
	_, err = leafCert.Verify(x509.VerifyOptions{
		Roots:     certPool,
		DNSName:   "test.openshift.io",
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	})
	assert.NoError(t, err, "leaf cert should validate against ECDSA CA")
}
