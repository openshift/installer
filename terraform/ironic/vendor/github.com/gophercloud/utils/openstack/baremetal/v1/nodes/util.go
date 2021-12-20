package nodes

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Gzips a file
func GzipFile(path string) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(contents)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	result := buf.Bytes()
	return result, nil
}

// Packs a directory into a gzipped ISO image
func PackDirectoryAsISO(path string) ([]byte, error) {
	iso, err := ioutil.TempFile("", "gophercloud-iso")
	if err != nil {
		return nil, err
	}
	iso.Close()
	defer os.Remove(iso.Name())
	cmd := exec.Command(
		"mkisofs",
		"-o", iso.Name(),
		"-ldots",
		"-allow-lowercase",
		"-allow-multidot", "-l",
		"-publisher", "gophercloud",
		"-quiet", "-J",
		"-r", "-V", "config-2",
		path,
	)
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("error creating configdrive iso: %s", err.Error())
	}

	return GzipFile(iso.Name())
}
