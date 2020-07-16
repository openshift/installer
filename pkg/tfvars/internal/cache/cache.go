package cache

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/h2non/filetype/matchers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xi2/xz"
	"golang.org/x/sys/unix"
)

const (
	applicationName = "openshift-installer"
	imageDataType   = "image"
)

// getCacheDir returns a local path of the cache, where the installer should put the data:
// <user_cache_dir>/openshift-installer/<dataType>_cache
// If the directory doesn't exist, it will be automatically created.
func getCacheDir(dataType string) (string, error) {
	if dataType == "" {
		return "", errors.Errorf("data type can't be an empty string")
	}

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(userCacheDir, applicationName, dataType+"_cache")

	_, err = os.Stat(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(cacheDir, 0755)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return cacheDir, nil
}

// cacheFile puts data in the cache
func cacheFile(reader io.Reader, filePath string, sha256Checksum string) (err error) {
	logrus.Debugf("Unpacking file into %q...", filePath)

	flockPath := fmt.Sprintf("%s.lock", filePath)
	flock, err := os.Create(flockPath)
	if err != nil {
		return err
	}
	defer flock.Close()
	defer func() {
		err2 := os.Remove(flockPath)
		if err == nil {
			err = err2
		}
	}()

	err = unix.Flock(int(flock.Fd()), unix.LOCK_EX)
	if err != nil {
		return err
	}
	defer func() {
		err2 := unix.Flock(int(flock.Fd()), unix.LOCK_UN)
		if err == nil {
			err = err2
		}
	}()

	_, err = os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil // another cacheFile beat us to it
	}

	tempPath := fmt.Sprintf("%s.tmp", filePath)

	// Delete the temporary file that may have been left over from previous launches.
	err = os.Remove(tempPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Errorf("failed to clean up %s: %v", tempPath, err)
		}
	} else {
		logrus.Debugf("Temporary file %v that remained after the previous launches was deleted", tempPath)
	}

	file, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0444)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			file.Close()
		}
	}()

	// Detect whether we know how to decompress the file
	// See http://golang.org/pkg/net/http/#DetectContentType for why we use 512
	buf := make([]byte, 512)
	_, err = reader.Read(buf)
	if err != nil {
		return err
	}

	reader = io.MultiReader(bytes.NewReader(buf), reader)
	switch {
	case matchers.Gz(buf):
		logrus.Debug("decompressing the image archive as gz")
		uncompressor, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer uncompressor.Close()
		reader = uncompressor
	case matchers.Xz(buf):
		logrus.Debug("decompressing the image archive as xz")
		uncompressor, err := xz.NewReader(reader, 0)
		if err != nil {
			return err
		}
		reader = uncompressor
	default:
		// No need for an interposer otherwise
		logrus.Debug("no known archive format detected for image, assuming no decompression necessary")
	}

	// Wrap the reader in TeeReader to calculate sha256 checksum on the fly
	hasher := sha256.New()
	if sha256Checksum != "" {
		reader = io.TeeReader(reader, hasher)
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	closed = true

	// Validate sha256 checksum
	if sha256Checksum != "" {
		foundChecksum := fmt.Sprintf("%x", hasher.Sum(nil))
		if sha256Checksum != foundChecksum {
			logrus.Error("File sha256 checksum is invalid.")
			return errors.Errorf("Checksum mismatch for %s; expected=%s found=%s", filePath, sha256Checksum, foundChecksum)
		}

		logrus.Debug("Checksum validation is complete...")
	}

	return os.Rename(tempPath, filePath)
}

// DownloadFile obtains a file from a given URL, puts it in the cache folder, defined by dataType parameter,
// and returns the local file path.
// If the query string contains sha256 parameter (i.e. https://example.com/data.bin?sha256=098a5a...),
// then the downloaded data checksum will be compared with the provided value.
func DownloadFile(baseURL string, dataType string) (string, error) {
	// Convert the given URL into a file name using md5 algorithm
	fileName := fmt.Sprintf("%x", md5.Sum([]byte(baseURL)))

	cacheDir, err := getCacheDir(dataType)
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(cacheDir, fileName)

	// If the file has already been cached, return its path
	_, err = os.Stat(filePath)
	if err == nil {
		logrus.Infof("The file was found in cache: %v. Reusing...", filePath)
		return filePath, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}

	// Send a request
	resp, err := http.Get(baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("bad status: %s", resp.Status)
	}

	// Get sha256 checksum if it was provided as a part of the URL
	var sha256Checksum string
	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}
	if sha256Checksums, ok := parsedURL.Query()["sha256"]; ok {
		sha256Checksum = sha256Checksums[0]
	}

	err = cacheFile(resp.Body, filePath, sha256Checksum)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// DownloadImageFile is a helper function that obtains an image file from a given URL,
// puts it in the cache and returns the local file path.  If the file is compressed
// by a known compressor, the file is uncompressed prior to being returned.
func DownloadImageFile(baseURL string) (string, error) {
	logrus.Infof("Obtaining RHCOS image file from '%v'", baseURL)

	return DownloadFile(baseURL, imageDataType)
}
