package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
	"time"
)

func TestSelfSignedCACert(t *testing.T) {
	key, err := GeneratePrivateKey()
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
			},
			err: false,
		},
		{
			cfg: &CertCfg{
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Subject: pkix.Name{
					CommonName: "root_ca",
				},
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
