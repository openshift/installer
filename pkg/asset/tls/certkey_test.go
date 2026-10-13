package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
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
			rootCA := &RootCA{}
			err := rootCA.Generate(context.Background(), nil)
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
