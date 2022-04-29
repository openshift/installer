package isosource

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

const (
	outputFile = "output/coreos.iso"

	isoURL    = "https://rhcos-redirector.apps.art.xq1c.p1.openshiftapps.com/art/storage/releases/rhcos-4.10/410.84.202201251210-0/x86_64/rhcos-410.84.202201251210-0-live.x86_64.iso"
	isoSha256 = "2905c1f0d85739e8600e8816c0d32711fb4002be4f845e0b20eeab35314e5b58"
)

func downloadIso(dest string) error {
	resp, err := http.Get(isoURL)
	if err != nil {
		return err
	}

	dir, _ := path.Split(dest)
	if err = os.MkdirAll(dir, 0775); err != nil {
		return err
	}

	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, resp.Body)
	return err
}

func haveValidIso(location string) bool {
	iso, err := os.OpenFile(location, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer iso.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, iso); err != nil {
		return false
	}

	expectedChecksum, err := hex.DecodeString(isoSha256)
	if err != nil {
		panic(err)
	}

	return bytes.Equal(hash.Sum(nil), expectedChecksum)
}

// EnsureIso downloads the ISO if it is not already present
func EnsureIso() (string, error) {
	if !haveValidIso(outputFile) {
		if err := downloadIso(outputFile); err != nil {
			return "", err
		}
		if !haveValidIso(outputFile) {
			return "", fmt.Errorf("Downloaded ISO is not valid")
		}
	}
	return outputFile, nil
}
