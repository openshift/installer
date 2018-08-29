package tls

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
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

func TestCertKeyGenerate(t *testing.T) {
	testDir, err := ioutil.TempDir(os.TempDir(), "certkey_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	root := &RootCA{rootDir: testDir}
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
		certKey *CertKey
		err     bool
		parents map[asset.Asset]*asset.State
	}{
		{
			&CertKey{
				rootDir:       testDir,
				installConfig: installConfig,
				Subject:       pkix.Name{CommonName: "test0-ca", OrganizationalUnit: []string{"openshift"}},
				KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Validity:      ValidityTenYears,
				KeyFileName:   "test0-ca.key",
				CertFileName:  "test0-ca.crt",
				IsCA:          true,
				ParentCA:      root,
			},
			false,
			map[asset.Asset]*asset.State{
				root: rootState,
			},
		},
		{
			&CertKey{
				rootDir:        testDir,
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
			false,
			map[asset.Asset]*asset.State{
				root:          rootState,
				installConfig: installConfigState,
			},
		},
		{
			&CertKey{
				rootDir:        testDir,
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
			true,
			nil,
		},
	}

	for i, tt := range tests {
		st, err := tt.certKey.Generate(tt.parents)
		if tt.err != (err != nil) {
			t.Errorf("test #%d error is not expected, expect %v, saw %v", i, tt.err, err != nil)
		}

		if err != nil {
			continue

		}

		keyFileName := assetFilePath(testDir, tt.certKey.KeyFileName)
		crtFileName := assetFilePath(testDir, tt.certKey.CertFileName)

		keyData, err := ioutil.ReadFile(keyFileName)
		if err != nil {
			t.Errorf("test #%d failed to read key file %q: %v", i, keyFileName, err)
		}

		crtData, err := ioutil.ReadFile(crtFileName)
		if err != nil {
			t.Errorf("test #%d failed to read key file %q: %v", i, crtFileName, err)
		}

		if !bytes.Equal(st.Contents[0].Data, keyData) {
			t.Errorf("test #%d expect key data: %q, saw %q", i, string(st.Contents[0].Data), string(keyData))
		}

		if !bytes.Equal(st.Contents[1].Data, crtData) {
			t.Errorf("test #%d expect crt data: %q, saw %q", i, string(st.Contents[1].Data), string(crtData))
		}

		// Briefly check the certs.
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(crtData) {
			t.Errorf("test #%d failed to append certs from PEM", i)
		}

		opts := x509.VerifyOptions{
			Roots:   certPool,
			DNSName: tt.certKey.Subject.CommonName,
		}
		if tt.certKey.GenDNSNames != nil {
			opts.DNSName = "test.openshift.io"
		}

		cert, err := PemToCertificate(crtData)
		if err != nil {
			t.Errorf("test #%d failed to parse certificate: %v", i, err)
		}

		if _, err := cert.Verify(opts); err != nil {
			t.Errorf("test #%d failed to verify cert: %v", i, err)
		}
	}
}
