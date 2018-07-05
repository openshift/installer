package tls

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
	"time"
)

func TestSelfSignedCACert(t *testing.T) {
	key, err := PrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate Private Key: %v", err)
	}
	cases := []struct {
		cfg *CertCfg
		err bool
	}{
		{
			cfg: &CertCfg{
				Validity:  time.Hour * 5,
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Subject: pkix.Name{
					CommonName:         "root_ca",
					OrganizationalUnit: []string{"openshift"},
				},
				IsCA: true,
			},
			err: false,
		},
		{
			cfg: &CertCfg{
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Subject: pkix.Name{
					CommonName: "root_ca",
				},
				IsCA: false,
			},
			err: true,
		},
		{
			cfg: &CertCfg{
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Subject: pkix.Name{
					OrganizationalUnit: []string{"openshift"},
				},
			},
			err: true,
		},
	}
	for i, c := range cases {
		if _, err := SelfSignedCACert(c.cfg, key); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestSignedCertificate(t *testing.T) {
	key, err := PrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	cases := []struct {
		Subject            pkix.Name
		SignatureAlgorithm x509.SignatureAlgorithm
		err                bool
	}{
		{
			Subject: pkix.Name{
				CommonName:         "csr",
				OrganizationalUnit: []string{"openshift"},
			},
			err: false,
		},
		{
			Subject: pkix.Name{},
			err:     false,
		},
		{
			Subject: pkix.Name{
				CommonName:         "csr-wrong-alg",
				OrganizationalUnit: []string{"openshift"},
			},
			SignatureAlgorithm: 123,
			err:                true,
		},
	}
	for i, c := range cases {
		csrTmpl := x509.CertificateRequest{
			Subject:            c.Subject,
			SignatureAlgorithm: c.SignatureAlgorithm,
		}
		if _, err := x509.CreateCertificateRequest(rand.Reader, &csrTmpl, key); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}
