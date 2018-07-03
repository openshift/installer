package configgenerator

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/coreos/tectonic-installer/installer/pkg/tls"
)

func writeFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	w.Flush()

	return nil
}

func copyFile(fromFilePath, toFilePath string) error {
	from, err := os.Open(fromFilePath)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.OpenFile(toFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	return err
}

// privateKeyToPem gets the content of the private key and returns a pem string
func privateKeyToPem(key *rsa.PrivateKey) string {
	keyInBytes := x509.MarshalPKCS1PrivateKey(key)
	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyInBytes,
		},
	)
	return string(keyinPem)
}

func certToPem(cert *x509.Certificate) string {
	certInPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
	return string(certInPem)
}

// generatePrivateKey generates and returns an *rsa.Privatekey object
func generatePrivateKey(clusterDir string, path string) (*rsa.PrivateKey, error) {
	fileTargetPath := filepath.Join(clusterDir, path)
	key, err := tls.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("error generating private key: %v", err)
	}
	if err := writeFile(fileTargetPath, privateKeyToPem(key)); err != nil {
		return nil, err
	}
	return key, nil
}
