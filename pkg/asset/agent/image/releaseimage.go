package image

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/openshift/installer/pkg/version"
)

type releaseImage struct {
	ReleaseVersion string   `json:"openshift_version"`
	Arch           string   `json:"cpu_architecture"`
	Archs          []string `json:"cpu_architectures"`
	PullSpec       string   `json:"url"`
	Tag            string   `json:"version"`
}

func isDigest(pullspec string) bool {
	return regexp.MustCompile(`.*sha256:[a-fA-F0-9]{64}$`).MatchString(pullspec)
}

func releaseImageFromPullSpec(pullSpec string, arch string, archs []string, version string) (releaseImage, error) {
	// When the pullspec it's a digest let's use the current version
	// stored in the installer
	if isDigest(pullSpec) {
		return releaseImage{
			ReleaseVersion: version,
			Arch:           arch,
			Archs:          archs,
			PullSpec:       pullSpec,
			Tag:            version,
		}, nil
	}

	components := strings.Split(pullSpec, ":")
	if len(components) < 2 {
		return releaseImage{}, fmt.Errorf("invalid release image \"%s\"", pullSpec)
	}
	lastIndex := len(components) - 1
	tag := strings.TrimSuffix(components[lastIndex], fmt.Sprintf("-%s", arch))

	versionComponents := strings.Split(tag, ".")
	if len(versionComponents) < 2 {
		return releaseImage{}, fmt.Errorf("invalid release image version \"%s\"", tag)
	}
	relVersion := strings.Join(versionComponents[:2], ".")

	return releaseImage{
		ReleaseVersion: relVersion,
		Arch:           arch,
		Archs:          archs,
		PullSpec:       pullSpec,
		Tag:            tag,
	}, nil
}

func releaseImageList(pullSpec, arch string, archs []string) (string, error) {
	versionString, err := version.Version()
	if err != nil {
		return "", err
	}

	relImage, err := releaseImageFromPullSpec(pullSpec, arch, archs, versionString)
	if err != nil {
		return "", err
	}

	imageList := []interface{}{relImage}
	text, err := json.Marshal(imageList)
	return string(text), err
}

func releaseImageListWithVersion(pullSpec, arch string, archs []string, openshiftVersion string) (string, error) {
	relImage, err := releaseImageFromPullSpec(pullSpec, arch, archs, openshiftVersion)
	if err != nil {
		return "", err
	}

	imageList := []interface{}{relImage}
	text, err := json.Marshal(imageList)
	return string(text), err
}
