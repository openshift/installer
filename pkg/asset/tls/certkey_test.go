package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	"github.com/stretchr/testify/assert"
)

type fakeInstallConfig bool

var _ asset.Asset = fakeInstallConfig(false)

func (f fakeInstallConfig) Dependencies() []asset.Asset {
	return nil
}

func (f fakeInstallConfig) Generate(map[asset.Asset]*asset.State) (*asset.State, error) {
	return &asset.State{
		Contents: []asset.Content{
			{
				Name: "fakeInstallConfig",
				Data: []byte{},
			},
		},
	}, nil
}

func (f fakeInstallConfig) Name() string {
	return "Fake Install Config"
}

func TestCertKeyGenerate(t *testing.T) {
	root := &RootCA{}
	rootState, err := root.Generate(nil)
	if err != nil {
		t.Fatal(err)
	}

	var installConfig fakeInstallConfig
	installConfigState, err := installConfig.Generate(nil)
	if err != nil {
		t.Fatal(err)
	}

	testGenSubject := func(*types.InstallConfig) (pkix.Name, error) {
		return pkix.Name{CommonName: "test", OrganizationalUnit: []string{"openshift"}}, nil
	}

	testGenDNSNames := func(*types.InstallConfig) ([]string, error) {
		return []string{"test.openshift.io"}, nil
	}

	testGenIPAddresses := func(*types.InstallConfig) ([]net.IP, error) {
		return []net.IP{net.ParseIP("10.0.0.1")}, nil
	}

	tests := []struct {
		name      string
		certKey   *CertKey
		errString string
		parents   map[asset.Asset]*asset.State
	}{
		{
			name: "simple ca",
			certKey: &CertKey{
				installConfig: installConfig,
				Subject:       pkix.Name{CommonName: "test0-ca", OrganizationalUnit: []string{"openshift"}},
				KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Validity:      ValidityTenYears,
				KeyFileName:   "test0-ca.key",
				CertFileName:  "test0-ca.crt",
				IsCA:          true,
				ParentCA:      root,
			},
			parents: map[asset.Asset]*asset.State{
				root: rootState,
			},
		},
		{
			name: "more complicated ca",
			certKey: &CertKey{
				installConfig:  installConfig,
				Subject:        pkix.Name{CommonName: "test1-ca", OrganizationalUnit: []string{"openshift"}},
				KeyUsages:      x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Validity:       ValidityTenYears,
				KeyFileName:    "test1-ca.key",
				CertFileName:   "test1-ca.crt",
				IsCA:           true,
				ParentCA:       root,
				AppendParent:   true,
				GenSubject:     testGenSubject,
				GenDNSNames:    testGenDNSNames,
				GenIPAddresses: testGenIPAddresses,
			},
			parents: map[asset.Asset]*asset.State{
				root:          rootState,
				installConfig: installConfigState,
			},
		},
		{
			name: "can't find parents",
			certKey: &CertKey{
				installConfig:  installConfig,
				Subject:        pkix.Name{CommonName: "test1-ca", OrganizationalUnit: []string{"openshift"}},
				KeyUsages:      x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Validity:       ValidityTenYears,
				KeyFileName:    "test2-ca.key",
				CertFileName:   "test2-ca.crt",
				IsCA:           true,
				ParentCA:       root,
				AppendParent:   true,
				GenSubject:     testGenSubject,
				GenDNSNames:    testGenDNSNames,
				GenIPAddresses: testGenIPAddresses,
			},
			errString: "failed to get install config state in the parent asset states",
			parents:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := tt.certKey.Generate(tt.parents)
			if err != nil {
				assert.EqualErrorf(t, err, tt.errString, tt.name)
				return
			} else if tt.errString != "" {
				t.Errorf("expect error %v, saw nil", err)
			}

			assert.Equal(t, assetFilePath(tt.certKey.KeyFileName), st.Contents[0].Name, "unexpected key file name")
			assert.Equal(t, assetFilePath(tt.certKey.CertFileName), st.Contents[1].Name, "unexpected cert file name")

			// Briefly check the certs.
			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(st.Contents[1].Data) {
				t.Error("failed to append certs from PEM")
			}

			opts := x509.VerifyOptions{
				Roots:   certPool,
				DNSName: tt.certKey.Subject.CommonName,
			}
			if tt.certKey.GenDNSNames != nil {
				opts.DNSName = "test.openshift.io"
			}

			cert, err := PemToCertificate(st.Contents[1].Data)
			assert.NoError(t, err, tt.name)

			_, err = cert.Verify(opts)
			assert.NoError(t, err, tt.name)
		})
	}
}
