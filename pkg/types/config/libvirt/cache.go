package libvirt

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

// UseCachedImage leaves non-file:// image URIs unalterered.
// Other URIs are retrieved with a local cache at
// $XDG_CACHE_HOME/openshift-install/libvirt [1].  This allows you to
// use the same remote image URI multiple times without needing to
// worry about redundant downloads, although you will want to
// periodically blow away your cache.
//
// [1]: https://standards.freedesktop.org/basedir-spec/basedir-spec-0.7.html
func (libvirt *Libvirt) UseCachedImage() (err error) {
	// FIXME: set the default URI here?  Leave it elsewhere?

	if strings.HasPrefix(libvirt.Image, "file://") {
		return nil
	}

	// FIXME: Use os.UserCacheDir() once we bump to Go 1.11
	// baseCacheDir, err := os.UserCacheDir()
	// if err != nil {
	// 	return err
	// }
	baseCacheDir := filepath.Join(os.Getenv("HOME"), ".cache")

	cacheDir := filepath.Join(baseCacheDir, "openshift-install", "libvirt")
	err = os.MkdirAll(cacheDir, 0777)
	if err != nil {
		return err
	}

	cache := diskcache.New(cacheDir)
	transport := httpcache.NewTransport(cache)
	resp, err := transport.Client().Get(libvirt.Image)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s while getting %s", resp.Status, libvirt.Image)
	}
	defer resp.Body.Close()

	var reader io.Reader
	if strings.HasSuffix(libvirt.Image, ".gz") {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
	} else {
		reader = resp.Body
	}

	// FIXME: diskcache's diskv backend doesn't expose direct file access.
	// We can write our own cache implementation to get around this,
	// but for now just dump this into /tmp without cleanup.
	file, err := ioutil.TempFile("", "openshift-install-")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	libvirt.Image = fmt.Sprintf("file://%s", filepath.ToSlash(file.Name()))
	return nil
}
