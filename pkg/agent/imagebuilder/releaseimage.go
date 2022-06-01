package imagebuilder

import (
	"encoding/json"
	"fmt"
	"strings"
)

type releaseImage struct {
	ReleaseVersion string `json:"openshift_version"`
	Arch           string `json:"cpu_architecture"`
	PullSpec       string `json:"url"`
	Tag            string `json:"version"`
}

func releaseImageFromPullSpec(pullSpec, arch string) (releaseImage, error) {
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
	// TODO: re-enable this code instead of hard-coding
	/*
		relImage, err := releaseImageFromPullSpec(pullSpec, arch)
		if err != nil {
			return "", err
		}
	*/
	relImage := releaseImage{
		ReleaseVersion: "4.10",
		Arch:           arch,
		PullSpec:       "quay.io/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64",
		Tag:            "4.10.0-rc.1",
	}

	imageList := []interface{}{relImage}
	text, err := json.Marshal(imageList)
	return string(text), err
}
