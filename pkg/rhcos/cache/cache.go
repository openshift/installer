package cache

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/filetype/matchers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
	"github.com/ulikunitz/xz"
)

const (
	// InstallerApplicationName is to use as application name by installer.
	InstallerApplicationName = "openshift-installer"
	// AgentApplicationName is to use as application name used by agent.
	AgentApplicationName = "agent"
	// ImageBasedApplicationName is to use as application name used by image-based.
	ImageBasedApplicationName = "imagebased"
	// ImageDataType is used by installer.
	ImageDataType = "image"
	// FilesDataType is used by agent.
	FilesDataType = "files"
)

// GetFileFromCache returns path of the cached file if found, otherwise returns an empty string
// or error.
func GetFileFromCache(fileName string, cacheDir string) (string, error) {
	filePath := filepath.Join(cacheDir, fileName)

	// If the file has already been cached, return its path
	_, err := os.Stat(filePath)
	if err == nil {
		logrus.Debugf("The file was found in cache: %v. Reusing...", filePath)
		return filePath, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}

	return "", nil
}

// GetCacheDir returns a local path of the cache, where the installer should put the data:
// <user_cache_dir>/agent/<dataType>_cache
// If the directory doesn't exist, it will be automatically created.
func GetCacheDir(dataType, applicationName string) (string, error) {
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

// cacheFile puts data in the cache.
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

	err = flockFile(flock, true)
	if err != nil {
		return err
	}
	defer func() {
		err2 := flockFile(flock, false)
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
		uncompressor, err := xz.NewReader(reader)
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

	written, err := io.Copy(file, reader)
	if err != nil {
		return err
	}

	// Let's find out how much data was written
	// for future troubleshooting
	logrus.Debugf("writing the RHCOS image was %d bytes", written)

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

// urlWithIntegrity pairs a URL with an optional expected sha256 checksum (after decompression, if any)
// If the query string contains sha256 parameter (i.e. https://example.com/data.bin?sha256=098a5a...),
// then the downloaded data checksum will be compared with the provided value.
type urlWithIntegrity struct {
	location           url.URL
	uncompressedSHA256 string
}

func (u *urlWithIntegrity) uncompressedName() string {
	n := filepath.Base(u.location.Path)
	return strings.TrimSuffix(strings.TrimSuffix(n, ".gz"), ".xz")
}

// download obtains a file from a given URL, puts it in the cache folder, defined by dataType parameter,
// and returns the local file path.
func (u *urlWithIntegrity) download(dataType, applicationName string) (string, error) {
	fileName := u.uncompressedName()

	cacheDir, err := GetCacheDir(dataType, applicationName)
	if err != nil {
		return "", err
	}

	filePath, err := GetFileFromCache(fileName, cacheDir)
	if err != nil {
		return "", err
	}
	if filePath != "" {
		// Found cached file
		return filePath, nil
	}

	// Send a request to get the file
	err = retry.DoFunc(3, 5*time.Second, func() error {
		resp, err := http.Get(u.location.String())
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Let's find the content length for future debugging
		logrus.Debugf("image download content length: %d", resp.ContentLength)

		// Check server response
		if resp.StatusCode != http.StatusOK {
			return errors.Errorf("bad status: %s", resp.Status)
		}

		filePath = filepath.Join(cacheDir, fileName)
		return cacheFile(resp.Body, filePath, u.uncompressedSHA256)
	})
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// DownloadImageFile is a helper function that obtains an image file from a given URL,
// puts it in the cache and returns the local file path.  If the file is compressed
// by a known compressor, the file is uncompressed prior to being returned.
func DownloadImageFile(baseURL string, applicationName string) (string, error) {
	return DownloadImageFileWithSha(baseURL, applicationName, "")
}

// DownloadImageFileWithSha sets the sha256Checksum which is checked on download.
func DownloadImageFileWithSha(baseURL string, applicationName string, sha256Checksum string) (string, error) {
	logrus.Debugf("Obtaining RHCOS image file from '%v'", baseURL)

	var u urlWithIntegrity
	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}
	q := parsedURL.Query()
	if sha256Checksum != "" {
		u.uncompressedSHA256 = sha256Checksum
	}
	if uncompressedSHA256, ok := q["sha256"]; ok {
		if sha256Checksum != "" && uncompressedSHA256[0] != sha256Checksum {
			return "", errors.Errorf("supplied sha256Checksum does not match URL")
		}
		u.uncompressedSHA256 = uncompressedSHA256[0]
		q.Del("sha256")
		parsedURL.RawQuery = q.Encode()
	}
	u.location = *parsedURL

	return u.download(ImageDataType, applicationName)
}
