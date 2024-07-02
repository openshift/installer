package image

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

func TestInstallationService_Generate(t *testing.T) {
	cases := []struct {
		name         string
		dependencies []asset.Asset

		expectedContent string
	}{
		{
			name: "default",
			dependencies: []asset.Asset{
				&ImageBasedInstallationConfig{
					Config: ibiConfig().build(),
				},
			},

			expectedContent: "[Unit]\nWants=network-online.target\nAfter=network-online.target\nDescription=SNO Image-based Installation\n[Service]\nEnvironment=SEED_IMAGE=quay.io/openshift-kni/seed-image:4.16.0\nEnvironment=HTTP_PROXY=\nEnvironment=http_proxy=\nEnvironment=HTTPS_PROXY=\nEnvironment=https_proxy=\nEnvironment=NO_PROXY=\nEnvironment=no_proxy=\nEnvironment=IBI_CONFIGURATION_FILE=/var/tmp/ibi-configuration.json\nEnvironment=PULL_SECRET_FILE=/var/tmp/pull-secret.json\nType=oneshot\nRemainAfterExit=yes\nExecStartPre=/usr/bin/chcon -t install_exec_t /usr/local/bin/install-rhcos-and-restore-seed.sh\nExecStart=/usr/local/bin/install-rhcos-and-restore-seed.sh\n[Install]\nWantedBy=multi-user.target\n",
		},
		{
			name: "with proxy configuration",
			dependencies: []asset.Asset{
				&ImageBasedInstallationConfig{
					Config: ibiConfig().
						proxy(&types.Proxy{
							HTTPProxy:  "a-http-proxy",
							HTTPSProxy: "an-https-proxy",
							NoProxy:    "localhost",
						}).
						build(),
				},
			},

			expectedContent: "[Unit]\nWants=network-online.target\nAfter=network-online.target\nDescription=SNO Image-based Installation\n[Service]\nEnvironment=SEED_IMAGE=quay.io/openshift-kni/seed-image:4.16.0\nEnvironment=HTTP_PROXY=a-http-proxy\nEnvironment=http_proxy=a-http-proxy\nEnvironment=HTTPS_PROXY=an-https-proxy\nEnvironment=https_proxy=an-https-proxy\nEnvironment=NO_PROXY=localhost\nEnvironment=no_proxy=localhost\nEnvironment=IBI_CONFIGURATION_FILE=/var/tmp/ibi-configuration.json\nEnvironment=PULL_SECRET_FILE=/var/tmp/pull-secret.json\nType=oneshot\nRemainAfterExit=yes\nExecStartPre=/usr/bin/chcon -t install_exec_t /usr/local/bin/install-rhcos-and-restore-seed.sh\nExecStart=/usr/local/bin/install-rhcos-and-restore-seed.sh\n[Install]\nWantedBy=multi-user.target\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			installationService := &InstallationService{}
			err := installationService.Generate(parents)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedContent, installationService.Content)
		})
	}
}
