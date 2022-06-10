package tls

import (
	"bytes"
	"encoding/pem"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CertBundle contains a multiple certificates in a bundle.
type CertBundle struct {
	asset.DefaultFileListWriter
	BundleRaw []byte
}

// Cert returns the certificate bundle.
func (b *CertBundle) Cert() []byte {
	return b.BundleRaw
}

// Generate generates the cert bundle from certs.
func (b *CertBundle) Generate(filename string, certs ...CertInterface) error {
	if len(certs) < 1 {
		return errors.New("atleast one certificate required for a bundle")
	}

	buf := bytes.Buffer{}
	for _, c := range certs {
		cert, err := PemToCertificate(c.Cert())
		if err != nil {
			logrus.Debugf("Failed to decode bundle certificate: %s", err)
			return errors.Wrap(err, "decoding certificate from PEM")
		}
		if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
			logrus.Debugf("Failed to encode bundle certificates: %s", err)
			return errors.Wrap(err, "encoding certificate to PEM")
		}
	}
	b.BundleRaw = buf.Bytes()
	b.FileList = []*asset.File{
		{
			Filename: assetFilePath(filename + ".crt"),
			Data:     b.BundleRaw,
		},
	}
	return nil
}

// Load is a no-op because TLS assets are not written to disk.
func (b *CertBundle) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}
