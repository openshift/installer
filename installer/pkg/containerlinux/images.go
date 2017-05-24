package containerlinux

import (
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-semver/semver"
)

const containerLinuxAMIsList = "https://stable.release.core-os.net/amd64-usr/current/coreos_production_ami_all.json"

// AMI represents a Container Linux AMI image.
type AMI struct {
	Name string `json:"name"`
	PV   string `json:"pv"`
	HVM  string `json:"hvm"`
}

// ListAMIImages returns the list of Container Linux AMIs.
func ListAMIImages(timeout time.Duration) ([]AMI, error) {
	client := &http.Client{
		Timeout:   time.Duration(timeout),
		Transport: http.DefaultTransport,
	}

	resp, err := client.Get(containerLinuxAMIsList)
	if err != nil {
		return []AMI{}, err
	}
	defer resp.Body.Close()

	st := struct {
		AMIs []AMI `json:"amis"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return []AMI{}, err
	}

	return st.AMIs, nil
}

// ListMatchboxImages fetches and parses the Container Linux images
// available at an endpoint (e.g. http://matchbox.foo:8080/assets/coreos/).
//
// The caller must provide a suitable image endpoint, Container Linux image
// indexes don't have a consistent identifier, just a list of image anchors.
func ListMatchboxImages(endpoint string, minVersion semver.Version, timeout time.Duration) ([]*semver.Version, error) {
	client := &http.Client{
		Timeout:   time.Duration(timeout),
		Transport: http.DefaultTransport,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var images []*semver.Version
	for _, img := range parseMatchboxImagesIndex(resp.Body) {
		if minVersion.LessThan(*img) {
			images = append(images, img)
		}
	}
	return images, nil
}

func parseMatchboxImagesIndex(r io.Reader) (images []*semver.Version) {
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch {
		case tt == html.StartTagToken:
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						vs := strings.Replace(attr.Val, "/", "", -1)
						v, err := semver.NewVersion(vs)
						if err != nil || v == nil {
							log.Warningf("Could not parse semver of Container Linux version %q: %s", vs, err)
							continue
						}
						images = append(images, v)
					}
				}
			}
		case tt == html.ErrorToken:
			return
		}
	}
}
