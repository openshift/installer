package image

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/common"
	"github.com/openshift/installer/pkg/asset/agent/interactive"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

func TestInteractiveDisconnectedIgnition_Generate(t *testing.T) {
	skipTestIfnmstatectlIsMissing(t)

	workingDirectory, err := os.Getwd()
	assert.NoError(t, err)
	err = os.Chdir(path.Join(workingDirectory, "../../../../data"))
	assert.NoError(t, err)

	cases := []struct {
		name             string
		deps             []asset.Asset
		expectedError    string
		expectedFiles    []string
		expectedServices map[string]bool
	}{
		{
			name: "unsupported workflow",
			deps: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&interactive.InstallConfig{},
				&manifests.ClusterImageSet{},
				&manifests.AgentPullSecret{},
				&manifests.InfraEnv{},
				&common.InfraEnvID{},
			},
			expectedError: "AgentWorkflowType value not supported: install",
		},
		{
			name: "default",
			deps: interactiveDefaultAssets(),

			expectedFiles:    interactiveExpectedFiles(),
			expectedServices: interactiveExpectedServices(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.deps...)

			asset := &InteractiveDisconnectedIgnition{}
			err := asset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assertExpectedFiles(t, asset.Config, tc.expectedFiles, nil)
				assertServiceEnabled(t, asset.Config, tc.expectedServices)
			}
		})
	}
}

func interactiveDefaultAssets() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstallInteractiveDisconnected},
		&interactive.InstallConfig{},
		&manifests.ClusterImageSet{
			Config: &v1.ClusterImageSet{
				Spec: v1.ClusterImageSetSpec{
					ReleaseImage: "registry.ci.openshift.org/ocp/release:4.99",
				},
			},
			File: &asset.File{
				Filename: "cluster-image-set.yaml",
				Data:     []byte("cluster-image-set.yaml"),
			},
		},
		&manifests.AgentPullSecret{
			File: &asset.File{
				Filename: "pull-secret.yaml",
				Data:     []byte("pull-secret.yaml"),
			},
		},
		&manifests.InfraEnv{
			File: &asset.File{
				Filename: "infraenv.yaml",
				Data:     []byte("infraenv.yaml"),
			},
		},
		&common.InfraEnvID{},
	}
}

func interactiveExpectedServices() map[string]bool {
	return map[string]bool{
		// enabled services
		"agent.service":                             true,
		"agent-interactive-console.service":         true,
		"agent-interactive-console-serial@.service": true,
		"agent-import-cluster.service":              true,
		"agent-register-cluster.service":            true,
		"agent-register-infraenv.service":           true,

		"assisted-service-db.service":  true,
		"assisted-service-pod.service": true,
		"assisted-service.service":     true,

		"install-status.service":     true,
		"iscsistart.service":         true,
		"iscsiadm.service":           true,
		"node-zero.service":          true,
		"oci-eval-user-data.service": true,
		"selinux.service":            true,
		"set-hostname.service":       true,

		// disabled services
		"agent-add-node.service":             false,
		"agent-auth-token-status.service":    false,
		"agent-check-config-image.service":   false,
		"apply-host-config.service":          false,
		"load-config-iso@.service":           false,
		"pre-network-manager-config.service": false,
		"start-cluster-installation.service": false,
	}
}

func interactiveExpectedFiles() []string {
	return []string{
		"/etc/issue",

		"/etc/assisted/agent-installer.env",
		"/etc/assisted/rendezvous-host.env",

		"/etc/assisted/manifests/cluster-image-set.yaml",
		"/etc/assisted/manifests/infraenv.yaml",
		"/etc/assisted/manifests/pull-secret.yaml",

		"/etc/containers/containers.conf",
		"/etc/motd.d/10-agent-installer",
		"/etc/multipath.conf",
		"/etc/NetworkManager/conf.d/clientid.conf",
		"/etc/systemd/system.conf.d/10-default-env.conf",
		"/etc/udev/rules.d/80-agent-config-image.rules",

		"/root/assisted.te",
		"/root/.docker/config.json",

		"/usr/local/bin/add-node.sh",
		"/usr/local/bin/agent-auth-token-status.sh",
		"/usr/local/bin/agent-config-image-wait.sh",
		"/usr/local/bin/agent-gather",
		"/usr/local/bin/bootstrap-service-record.sh",
		"/usr/local/bin/common.sh",
		"/usr/local/bin/extract-agent.sh",
		"/usr/local/bin/get-container-images.sh",
		"/usr/local/bin/install-status.sh",
		"/usr/local/bin/issue_status.sh",
		"/usr/local/bin/load-config-iso.sh",
		"/usr/local/bin/oci-eval-user-data.sh",
		"/usr/local/bin/release-image.sh",
		"/usr/local/bin/release-image-download.sh",
		"/usr/local/bin/set-hostname.sh",
		"/usr/local/bin/set-node-zero.sh",
		"/usr/local/bin/start-agent.sh",
		"/usr/local/bin/start-cluster-installation.sh",
		"/usr/local/bin/wait-for-assisted-service.sh",

		"/usr/local/share/assisted-service/assisted-db.env",
		"/usr/local/share/assisted-service/assisted-service.env",
		"/usr/local/share/assisted-service/images.env",
		"/usr/local/share/start-cluster/start-cluster.env",
	}
}
