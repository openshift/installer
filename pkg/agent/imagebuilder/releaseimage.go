package imagebuilder

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/openshift/installer/pkg/version"
)

type releaseImage struct {
	ReleaseVersion string `json:"openshift_version"`
	Arch           string `json:"cpu_architecture"`
	PullSpec       string `json:"url"`
	Tag            string `json:"version"`
}

func isDigest(pullspec string) bool {
	return regexp.MustCompile(`.*sha256:[a-fA-F0-9]{64}$`).MatchString(pullspec)
}

func releaseImageFromPullSpec(pullSpec, arch string) (releaseImage, error) {

	// When the pullspec it's a digest let's use the current version
	// stored in the installer
	if isDigest(pullSpec) {
		versionString, err := version.Version()
		if err != nil {
			return releaseImage{}, err
		}

		return releaseImage{
			ReleaseVersion: versionString,
			Arch:           arch,
			PullSpec:       pullSpec,
			Tag:            versionString,
		}, nil
	}

	components := strings.SplitN(pullSpec, ":", 2)
	if len(components) < 2 {
		return releaseImage{}, fmt.Errorf("invalid release image \"%s\"", pullSpec)
	}
	tag := strings.TrimSuffix(components[1], fmt.Sprintf("-%s", arch))

	versionComponents := strings.Split(tag, ".")
	if len(versionComponents) < 2 {
		return releaseImage{}, fmt.Errorf("invalid release image version \"%s\"", tag)
	}
	relVersion := strings.Join(versionComponents[:2], ".")

	return releaseImage{
		ReleaseVersion: relVersion,
		Arch:           arch,
		PullSpec:       pullSpec,
		Tag:            tag,
	}, nil
}

func releaseImageList(pullSpec, arch string) (string, error) {

	relImage, err := releaseImageFromPullSpec(pullSpec, arch)
	if err != nil {
		return "", err
	}

	imageList := []interface{}{relImage}
	text, err := json.Marshal(imageList)
	return string(text), err
}
