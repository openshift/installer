package libvirt

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// cachedImage leaves non-file:// image URIs unalterered.
// Other URIs are retrieved with a local cache at
// $XDG_CACHE_HOME/openshift-install/libvirt [1].  This allows you to
// use the same remote image URI multiple times without needing to
// worry about redundant downloads, although you will want to
// periodically blow away your cache.
//
// [1]: https://standards.freedesktop.org/basedir-spec/basedir-spec-0.7.html
func cachedImage(uri string) (string, error) {
	if strings.HasPrefix(uri, "file://") {
		return uri, nil
	}

	logrus.Infof("Fetching OS image: %s", filepath.Base(uri))

	// FIXME: Use os.UserCacheDir() once we bump to Go 1.11
	// baseCacheDir, err := os.UserCacheDir()
	// if err != nil {
	// 	return uri, err
	// }
	baseCacheDir := filepath.Join(os.Getenv("HOME"), ".cache")

	cacheDir := filepath.Join(baseCacheDir, "openshift-install", "libvirt")
	httpCacheDir := filepath.Join(cacheDir, "http")
	// We used to use httpCacheDir, warn if it still exists since the user may
	// want to delete it
	_, err := os.Stat(httpCacheDir)
	if err == nil || !os.IsNotExist(err) {
		logrus.Infof("The installer no longer uses %q, it can be deleted", httpCacheDir)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return uri, err
	}
	if resp.StatusCode != 200 {
		return uri, fmt.Errorf("%s while getting %s", resp.Status, uri)
	}
	defer resp.Body.Close()

	key, err := cacheKey(resp.Header.Get("ETag"))
	if err != nil {
		return uri, fmt.Errorf("invalid ETag for %s: %v", uri, err)
	}

	imageCacheDir := filepath.Join(cacheDir, "image")
	err = os.MkdirAll(imageCacheDir, 0777)
	if err != nil {
		return uri, err
	}

	imagePath := filepath.Join(imageCacheDir, key)
	_, err = os.Stat(imagePath)
	if err == nil {
		logrus.Debugf("Using cached OS image %q", imagePath)
	} else {
		if !os.IsNotExist(err) {
			return uri, err
		}

		err = cacheImage(resp.Body, imagePath)
		if err != nil {
			return uri, err
		}
	}

	return fmt.Sprintf("file://%s", filepath.ToSlash(imagePath)), nil
}

func cacheKey(etag string) (key string, err error) {
	if etag == "" {
		return "", fmt.Errorf("caching is not supported when ETag is unset")
	}
	etagSections := strings.SplitN(etag, "\"", 3)
	if len(etagSections) != 3 {
		return "", fmt.Errorf("broken quoting: %s", etag)
	}
	if etagSections[0] == "W/" {
		return "", fmt.Errorf("caching is not supported for weak ETags: %s", etag)
	}
	opaque := etagSections[1]
	if opaque == "" {
		return "", fmt.Errorf("caching is not supported when the opaque tag is unset: %s", etag)
	}
	hashed := md5.Sum([]byte(opaque))
	return hex.EncodeToString(hashed[:]), nil
}

type progressWatcher struct {
	total    uint64
	oldTotal uint64
}

func (w *progressWatcher) Write(p []byte) (int, error) {
	n := len(p)
	w.total += uint64(n)
	// Log progress every 3 MiB
	if w.total-w.oldTotal > 3000000 {
		logrus.Debugf("Read %d MiB of image", w.total/1000000)
		w.oldTotal = w.total
	}
	return n, nil
}

func (w *progressWatcher) finish() {
	logrus.Debugf("Read %d MiB of image", w.total/1000000)
}

func cacheImage(reader io.Reader, imagePath string) (err error) {
	logrus.Debugf("Unpacking OS image into %q...", imagePath)

	flockPath := fmt.Sprintf("%s.lock", imagePath)
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

	_, err = os.Stat(imagePath)
	if err != nil && !os.IsNotExist(err) {
		return nil // another cacheImage beat us to it
	}

	tempPath := fmt.Sprintf("%s.tmp", imagePath)
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

	watcher := &progressWatcher{}
	_, err = io.Copy(file, io.TeeReader(reader, watcher))
	if err != nil {
		return err
	}
	watcher.finish()

	err = file.Close()
	if err != nil {
		return err
	}
	closed = true

	return os.Rename(tempPath, imagePath)
}
