package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseImageList(t *testing.T) {
	cases := []struct {
		name     string
		pullSpec string
		arch     string
		result   string
	}{
		{
			name:     "4.10rc",
			pullSpec: "quay.io/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64",
			arch:     "x86_64",
			result:   "[{\"openshift_version\":\"4.10\",\"cpu_architecture\":\"x86_64\",\"cpu_architectures\":[\"x86_64\"],\"url\":\"quay.io/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64\",\"version\":\"4.10.0-rc.1\"}]",
		},
		{
			name:     "pull-spec-includes-port-number",
			pullSpec: "quay.io:433/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64",
			arch:     "x86_64",
			result:   "[{\"openshift_version\":\"4.10\",\"cpu_architecture\":\"x86_64\",\"cpu_architectures\":[\"x86_64\"],\"url\":\"quay.io:433/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64\",\"version\":\"4.10.0-rc.1\"}]",
		},
		{
			name:     "arm",
			pullSpec: "quay.io/openshift-release-dev/ocp-release:4.10.0-rc.1-aarch64",
			arch:     "aarch64",
			result:   "[{\"openshift_version\":\"4.10\",\"cpu_architecture\":\"aarch64\",\"cpu_architectures\":[\"aarch64\"],\"url\":\"quay.io/openshift-release-dev/ocp-release:4.10.0-rc.1-aarch64\",\"version\":\"4.10.0-rc.1\"}]",
		},
		{
			name:     "4.11ci",
			pullSpec: "registry.ci.openshift.org/ocp/release:4.11.0-0.ci-2022-05-16-202609",
			arch:     "x86_64",
			result:   "[{\"openshift_version\":\"4.11\",\"cpu_architecture\":\"x86_64\",\"cpu_architectures\":[\"x86_64\"],\"url\":\"registry.ci.openshift.org/ocp/release:4.11.0-0.ci-2022-05-16-202609\",\"version\":\"4.11.0-0.ci-2022-05-16-202609\"}]",
		},
		{
			name:     "CI-ephemeral",
			pullSpec: "registry.build04.ci.openshift.org/ci-op-m7rfgytz/release@sha256:ebb203f24ee060d61bdb466696a9c20b3841f9929badf9b81fc99cbedc2a679e",
			arch:     "x86_64",
			result:   "[{\"openshift_version\":\"was not built correctly\",\"cpu_architecture\":\"x86_64\",\"cpu_architectures\":[\"x86_64\"],\"url\":\"registry.build04.ci.openshift.org/ci-op-m7rfgytz/release@sha256:ebb203f24ee060d61bdb466696a9c20b3841f9929badf9b81fc99cbedc2a679e\",\"version\":\"was not built correctly\"}]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := releaseImageList(tc.pullSpec, tc.arch, []string{tc.arch})
			assert.NoError(t, err)
			if err == nil {
				assert.Equal(t, tc.result, output)
			}
		})
	}
}

func TestReleaseImageListErrors(t *testing.T) {
	cases := []string{
		"",
		"quay.io/openshift-release-dev/ocp-release-4.10",
		"quay.io/openshift-release-dev/ocp-release:4",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			_, err := releaseImageList(tc, "x86_64", []string{"x86_64"})
			assert.Error(t, err)
		})
	}
}
