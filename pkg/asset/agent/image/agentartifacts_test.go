package image

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	agentconfig "github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/types"
	agenttype "github.com/openshift/installer/pkg/types/agent"
)

// ociPlatformSpec is used to skip fetchAgentTuiFiles in tests so Generate
// fails fast at the ISO extraction step (no network calls required).
var ociPlatformSpec = &hiveext.ExternalPlatformSpec{
	PlatformName: "oci",
}

// buildParents constructs the minimal asset.Parents required by AgentArtifacts.Generate.
// cpuArch is taken from Ignition; platformType is set on AgentManifests.AgentClusterInstall.
func buildParents(cpuArch string, minimalISO bool, platformType hiveext.PlatformType) asset.Parents {
	parents := asset.Parents{}
	parents.Add(
		&Ignition{CPUArch: cpuArch},
		&Kargs{},
		// BaseIso with a non-existent filename so prepareAgentArtifacts fails fast.
		&BaseIso{File: &asset.File{Filename: "nonexistent-agent.iso"}},
		&manifests.AgentManifests{
			AgentClusterInstall: &hiveext.AgentClusterInstall{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: hiveext.AgentClusterInstallSpec{
					PlatformType: platformType,
				},
			},
			ClusterImageSet: &hivev1.ClusterImageSet{},
			PullSecret: &corev1.Secret{
				StringData: map[string]string{".dockerconfigjson": "{}"},
			},
		},
		&mirror.RegistriesConf{},
		// AgentClusterInstall with OCI ExternalPlatformSpec so that
		// fetchAgentTuiFiles is bypassed (see agentartifacts.go line 127).
		&manifests.AgentClusterInstall{
			Config: &hiveext.AgentClusterInstall{
				Spec: hiveext.AgentClusterInstallSpec{
					ExternalPlatformSpec: ociPlatformSpec,
					PlatformType:         platformType,
				},
			},
		},
		&agentconfig.AgentConfig{
			Config: &agenttype.Config{MinimalISO: minimalISO},
		},
		&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
	)
	return parents
}

// TestAgentArtifacts_MinimalISO_s390x verifies the minimal ISO validation
// that lives in AgentArtifacts.Generate():
//   - s390x + minimalISO:true  → hard error
//   - s390x + external platform → a.MinimalISO forced to false (no minimalISO error)
//   - non-s390x + minimalISO:true → a.MinimalISO set true (no minimalISO error)
func TestAgentArtifacts_MinimalISO_s390x(t *testing.T) {
	cases := []struct {
		name             string
		cpuArch          string
		minimalISO       bool
		platformType     hiveext.PlatformType
		expectMinimalISO bool
		expectErrContain string
	}{
		{
			name:             "s390x + minimalISO true → error",
			cpuArch:          string(types.ArchitectureS390X),
			minimalISO:       true,
			expectErrContain: "minimal ISO is not supported on the s390x architecture",
		},
		{
			name:             "s390x + external platform → MinimalISO forced false",
			cpuArch:          string(types.ArchitectureS390X),
			minimalISO:       false,
			platformType:     hiveext.ExternalPlatformType,
			expectMinimalISO: false,
		},
		{
			name:             "non-s390x + minimalISO true → MinimalISO set true",
			cpuArch:          "x86_64",
			minimalISO:       true,
			expectMinimalISO: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := buildParents(tc.cpuArch, tc.minimalISO, tc.platformType)

			a := &AgentArtifacts{}
			err := a.Generate(context.Background(), parents)

			if tc.expectErrContain != "" {
				assert.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.expectErrContain),
					"expected error containing %q, got: %v", tc.expectErrContain, err)
				return
			}

			// For the non-error-validation cases Generate will still fail (no ISO
			// on disk), but a.MinimalISO is assigned before that IO step.
			assert.Equal(t, tc.expectMinimalISO, a.MinimalISO,
				"MinimalISO mismatch; Generate error was: %v", err)
			if err != nil {
				assert.False(t,
					strings.Contains(err.Error(), "minimalISO") ||
						strings.Contains(err.Error(), "minimal ISO is not supported"),
					"unexpected minimalISO validation error: %v", err)
			}
		})
	}
}
