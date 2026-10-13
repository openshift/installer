package image

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

func TestKargs_Generate(t *testing.T) {
	cases := []struct {
		name                string
		workflow            workflow.AgentWorkflowType
		agentClusterInstall *manifests.AgentClusterInstall
		clusterInfo         *joiner.ClusterInfo
		expectedArgs        string
		expectedErr         string
	}{
		{
			name:         "install workflow - default",
			workflow:     workflow.AgentWorkflowTypeInstall,
			expectedArgs: "",
		},
		{
			name:     "install workflow - fips enabled",
			workflow: workflow.AgentWorkflowTypeInstall,
			agentClusterInstall: &manifests.AgentClusterInstall{
				Config: &v1beta1.AgentClusterInstall{
					ObjectMeta: v1.ObjectMeta{
						Annotations: map[string]string{
							"agent-install.openshift.io/install-config-overrides": `{"fips": true}`,
						},
					},
				},
			},
			expectedArgs: " fips=1",
		},
		{
			name:     "install workflow - oci with fips enabled",
			workflow: workflow.AgentWorkflowTypeInstall,
			agentClusterInstall: &manifests.AgentClusterInstall{
				Config: &v1beta1.AgentClusterInstall{
					ObjectMeta: v1.ObjectMeta{
						Annotations: map[string]string{
							"agent-install.openshift.io/install-config-overrides": `{"fips": true}`,
						},
					},
					Spec: v1beta1.AgentClusterInstallSpec{
						ExternalPlatformSpec: &v1beta1.ExternalPlatformSpec{
							PlatformName: agent.ExternalPlatformNameOci,
						},
					},
				},
			},
			expectedArgs: " console=ttyS0 fips=1",
		},
		{
			name:         "add-nodes workflow - default",
			workflow:     workflow.AgentWorkflowTypeAddNodes,
			expectedArgs: "",
		},
		{
			name:     "add-nodes workflow - fips enabled",
			workflow: workflow.AgentWorkflowTypeAddNodes,
			clusterInfo: &joiner.ClusterInfo{
				FIPS: true,
			},
			expectedArgs: " fips=1",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dependencies := []asset.Asset{
				&workflow.AgentWorkflow{Workflow: tc.workflow},
			}
			aci := &manifests.AgentClusterInstall{
				Config: &v1beta1.AgentClusterInstall{},
			}
			if tc.agentClusterInstall != nil {
				aci = tc.agentClusterInstall
			}
			ci := &joiner.ClusterInfo{}
			if tc.clusterInfo != nil {
				ci = tc.clusterInfo
			}

			dependencies = append(dependencies, ci)
			dependencies = append(dependencies, aci)
			parents := asset.Parents{}
			parents.Add(dependencies...)

			kargs := &Kargs{}
			err := kargs.Generate(context.Background(), parents)

			if tc.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArgs, kargs.KernelCmdLine())
			} else {
				assert.Regexp(t, tc.expectedErr, err.Error())
			}
		})
	}
}
