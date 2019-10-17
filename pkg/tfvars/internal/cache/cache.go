package cache

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
func cacheFile(reader io.Reader, filePath string) (err error) {
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

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	closed = true

	return os.Rename(tempPath, filePath)
}

// DownloadFile obtains a file from a given URL, puts it in the cache and returns
// the local file path.
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
	} else {
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

		err = cacheFile(resp.Body, filePath)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

// DownloadImageFile is a helper function that obtains an image file from a given URL,
// puts it in the cache and returns the local file path.
func DownloadImageFile(baseURL string) (string, error) {
	logrus.Infof("Obtaining RHCOS image file from '%v'", baseURL)

	return DownloadFile(baseURL, imageDataType)
}
