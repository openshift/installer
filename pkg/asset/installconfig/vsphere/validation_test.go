package vsphere

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/simulator"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"

	_ "github.com/vmware/govmomi/vapi/simulator"
)

var (
	validCIDR = "10.0.0.0/16"
)

func validIPIInstallConfig(vcenter string) *types.InstallConfig {
	datacenter := "DC0"
	portgroup := "DVPG0"
	cluster := "DC0_C0"
	datastore := "LocalDS_0"
	username := "my-username"
	password := "my-password"

	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				Cluster:          cluster,
				Datacenter:       datacenter,
				DefaultDatastore: datastore,
				Network:          portgroup,
				Password:         password,
				Username:         username,
				VCenter:          vcenter,
				APIVIP:           "192.168.111.0",
				IngressVIP:       "192.168.111.1",
			},
		},
	}
}

func validUPIInstallConfig(vcenter string) *types.InstallConfig {
	datacenter := "DC0"
	datastore := "LocalDS_0"
	username := "my-username"
	password := "my-password"
	cluster := "DC0_C0"
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				Datacenter:       datacenter,
				DefaultDatastore: datastore,
				Password:         password,
				Username:         username,
				VCenter:          vcenter,
				Cluster:          cluster,
			},
		},
	}
}

func TestValidate(t *testing.T) {
	model := simulator.VPX()

	// The number of vCenter objects we will want to have
	model.Cluster = 1
	model.Datacenter = 1
	model.Datastore = 1
	model.Portgroup = 1
	model.ClusterHost = 3

	// We don't need standalone hosts
	model.Host = 0

	defer model.Remove()
	err := model.Create()

	if err != nil {
		t.Error(err)
	}

	// Set the simulator's username and password
	// The simulator will always run on loopback and
	// port 22443. That way we can bypass the host validation
	// that would otherwise interfere with the simulator url.
	model.Service.Listen = &url.URL{
		User: url.UserPassword("my-username", "my-password"),
		Host: "127.0.0.1:22443",
	}

	ca, err := createCertificateAuthority()

	if err != nil {
		t.Error(err)
	}

	cert, keypair, err := createServerCertificate(ca)

	if err != nil {
		t.Error(err)
	}

	// Provide the server certificate to the simulator
	model.Service.TLS = &tls.Config{
		Certificates: []tls.Certificate{*keypair},
	}
	model.Service.RegisterEndpoints = true

	s := model.Service.NewServer()
	defer s.Close()

	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(*types.InstallConfig, ...x509.Certificate) error
		expectErr        string
	}{{
		name:             "valid UPI install config",
		installConfig:    validUPIInstallConfig(s.URL.Host),
		validationMethod: Validate,
	}, {
		name:             "valid IPI install config",
		installConfig:    validIPIInstallConfig(s.URL.Host),
		validationMethod: ValidateForProvisioning,
	}, {
		name: "invalid IPI - no network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(s.URL.Host)
			c.Platform.VSphere.Network = ""
			return c
		}(),
		validationMethod: ValidateForProvisioning,
		expectErr:        `^platform\.vsphere\.network: Required value: must specify the network$`,
	}, {
		name: "invalid IPI - no cluster",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(s.URL.Host)
			c.Platform.VSphere.Cluster = ""
			return c
		}(),
		validationMethod: ValidateForProvisioning,
		expectErr:        `^platform\.vsphere\.cluster: Required value: must specify the cluster$`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.validationMethod(test.installConfig, *cert)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err.Error())
			}
		})
	}
}

// This section of functions below are dervied from examples in these two
// sources:
// - https://medium.com/@shaneutt/create-sign-x509-certificates-in-golang-8ac4ae49f903
// - https://gist.github.com/Mattemagikern/328cdd650be33bc33105e26db88e487d
// The simulator is started with a self-signed certificate which is not created
// correctly.
// To resolve this we create a CA and server certificate. The server certificate
// is signed with the CA.
// The server certificate is provided to the simulator. The CA is
// provided to CreateVSphereClients
// This allows us to test using `insecure: false`.

func createCertificateAuthority() (*x509.Certificate, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),

		Subject: pkix.Name{
			Organization:  []string{"Self-Signed Red Hat, Inc."},
			Country:       []string{"US"},
			Province:      []string{"North Carolina"},
			Locality:      []string{"Raleigh"},
			StreetAddress: []string{"100 E. Davie Street"},
			PostalCode:    []string{"27601"},
		},
		NotBefore:             time.Now().AddDate(0, 0, -1),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	cert, _, err := createx509Certificate(ca, ca)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func createServerCertificate(parent *x509.Certificate) (*x509.Certificate, *tls.Certificate, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),

		Subject: pkix.Name{
			Organization: []string{"Self-Signed Red Hat, Inc."},
			Country:      []string{"US"},
			Province:     []string{"North Carolina"},
			Locality:     []string{"Raleigh"},
		},
		NotBefore:             time.Now().AddDate(0, 0, -1),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
	}
	cert, keypair, err := createx509Certificate(ca, parent)
	if err != nil {
		return nil, nil, err
	}

	return cert, keypair, nil
}

func createx509Certificate(template, parent *x509.Certificate) (*x509.Certificate, *tls.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, err
	}

	certPEMBlock := new(bytes.Buffer)
	pem.Encode(certPEMBlock, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	keyPEMBlock := new(bytes.Buffer)
	pem.Encode(keyPEMBlock, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	keypair, err := tls.X509KeyPair(certPEMBlock.Bytes(), keyPEMBlock.Bytes())

	if err != nil {
		return nil, nil, err
	}

	return cert, &keypair, nil
}
