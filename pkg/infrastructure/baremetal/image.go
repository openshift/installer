// This file is largely based on existing code from terraform-provider-libvirt 0.6.12.
// https://github.com/dmacvicar/terraform-provider-libvirt
// Original code distributed under the terms of Apache License 2.0.
package baremetal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"libvirt.org/go/libvirtxml"
)

type image interface {
	Size() (uint64, error)
	Import(func(io.Reader) error, libvirtxml.StorageVolume) error
	String() string
	IsQCOW2() (bool, error)
}

type localImage struct {
	path string
}

func (i *localImage) String() string {
	return i.path
}

func isQCOW2Header(buf []byte) (bool, error) {
	if len(buf) < 8 {
		return false, fmt.Errorf("expected header of 8 bytes. Got %d", len(buf))
	}
	if buf[0] == 'Q' && buf[1] == 'F' && buf[2] == 'I' && buf[3] == 0xfb && buf[4] == 0x00 && buf[5] == 0x00 && buf[6] == 0x00 && buf[7] == 0x03 {
		return true, nil
	}
	return false, nil
}

func (i *localImage) Size() (uint64, error) {
	fi, err := os.Stat(i.path)
	if err != nil {
		return 0, err
	}
	return uint64(fi.Size()), nil
}

func (i *localImage) IsQCOW2() (bool, error) {
	file, err := os.Open(i.path)
	if err != nil {
		return false, fmt.Errorf("error while opening %s: %w", i.path, err)
	}
	defer file.Close()
	buf := make([]byte, 8)
	_, err = io.ReadAtLeast(file, buf, 8)
	if err != nil {
		return false, err
	}
	return isQCOW2Header(buf)
}

func (i *localImage) Import(copier func(io.Reader) error, vol libvirtxml.StorageVolume) error {
	file, err := os.Open(i.path)
	if err != nil {
		return fmt.Errorf("error while opening %s: %w", i.path, err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return err
	}
	// we can skip the upload if the modification times are the same
	if vol.Target.Timestamps != nil && vol.Target.Timestamps.Mtime != "" {
		if fi.ModTime().Equal(timeFromEpoch(vol.Target.Timestamps.Mtime)) {
			logrus.Info("Modification time is the same: skipping image copy")
			return nil
		}
	}

	return copier(file)
}

type httpImage struct {
	url *url.URL
}

func (i *httpImage) String() string {
	return i.url.String()
}

func (i *httpImage) Size() (uint64, error) {
	response, err := http.Head(i.url.String())
	if err != nil {
		return 0, err
	}
	if response.StatusCode == 403 {
		// possibly only the HEAD method is forbidden, try a Body-less GET instead
		response, err = http.Get(i.url.String())
		if err != nil {
			return 0, err
		}

		response.Body.Close()
	}
	if response.StatusCode != 200 {
		return 0,
			fmt.Errorf(
				"error accessing remote resource: %s - %s",
				i.url.String(),
				response.Status)
	}

	length, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		err = fmt.Errorf(
			"error while getting Content-Length of \"%s\": %w - got %s",
			i.url.String(),
			err,
			response.Header.Get("Content-Length"))
		return 0, err
	}
	return uint64(length), nil
}

func (i *httpImage) IsQCOW2() (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", i.url.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Range", "bytes=0-7")
	response, err := client.Do(req)

	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode != 206 {
		return false, fmt.Errorf(
			"can't retrieve partial header of resource to determine file type: %s - %s",
			i.url.String(),
			response.Status)
	}

	header, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	if len(header) < 8 {
		return false, fmt.Errorf(
			"can't retrieve read header of resource to determine file type: %s - %d bytes read",
			i.url.String(),
			len(header))
	}

	return isQCOW2Header(header)
}

func (i *httpImage) Import(copier func(io.Reader) error, vol libvirtxml.StorageVolume) error {
	// number of download retries on non client errors (eg. 5xx)
	const maxHTTPRetries int = 3
	// wait time between retries
	const retryWait time.Duration = 2 * time.Second

	client := &http.Client{}
	req, err := http.NewRequest("GET", i.url.String(), nil)

	if err != nil {
		return fmt.Errorf("error while downloading %s: %w", i.url.String(), err)
	}

	if vol.Target.Timestamps != nil && vol.Target.Timestamps.Mtime != "" {
		req.Header.Set("If-Modified-Since", timeFromEpoch(vol.Target.Timestamps.Mtime).UTC().Format(http.TimeFormat))
	}

	var response *http.Response
	for retryCount := 0; retryCount < maxHTTPRetries; retryCount++ {
		response, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("error while downloading %s: %w", i.url.String(), err)
		}
		defer response.Body.Close()

		logrus.Debugf("url resp status code %s (retry #%d)\n", response.Status, retryCount)

		switch response.StatusCode {
		case http.StatusNotModified:
			return nil
		case http.StatusOK:
			return copier(response.Body)
		default:
			if response.StatusCode < 500 {
				break
			}
			// The problem is not client but server side
			// retry a few times after a small wait
			if retryCount < maxHTTPRetries {
				time.Sleep(retryWait)
			}
		}
	}
	return fmt.Errorf("error while downloading %s: %v", i.url.String(), response)
}

func newImage(source string) (image, error) {
	url, err := url.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("can't parse source '%s' as url: %w", source, err)
	}

	if strings.HasPrefix(url.Scheme, "http") {
		return &httpImage{url: url}, nil
	}

	if url.Scheme == "file" || url.Scheme == "" {
		return &localImage{path: url.Path}, nil
	}

	return nil, fmt.Errorf("don't know how to read from '%s': %w", url.String(), err)
}

func timeFromEpoch(str string) time.Time {
	var s, ns int
	var err error

	ts := strings.Split(str, ".")
	if len(ts) == 2 {
		ns, err = strconv.Atoi(ts[1])
		if err != nil {
			ns = 0
		}
	}
	s, err = strconv.Atoi(ts[0])
	if err != nil {
		s = 0
	}

	return time.Unix(int64(s), int64(ns))
}
