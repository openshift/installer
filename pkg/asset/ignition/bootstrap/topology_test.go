package bootstrap

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

func topologyManifests(mastersSchedulable bool, controlPlaneTopology, infrastructureTopology string) []*asset.File {
	return []*asset.File{
		{
			Filename: "manifests/cluster-scheduler-02-config.yml",
			Data: []byte(fmt.Sprintf(`apiVersion: config.openshift.io/v1
kind: Scheduler
metadata:
  name: cluster
spec:
  mastersSchedulable: %t
`, mastersSchedulable)),
		},
		{
			Filename: "manifests/cluster-infrastructure-02-config.yml",
			Data: []byte(fmt.Sprintf(`apiVersion: config.openshift.io/v1
kind: Infrastructure
metadata:
  name: cluster
status:
  controlPlaneTopology: %q
  infrastructureTopology: %q
`, controlPlaneTopology, infrastructureTopology)),
		},
	}
}

func TestValidateTopology(t *testing.T) {
	const (
		single  = "SingleReplica"
		dual    = "DualReplica"
		ha      = "HighlyAvailable"
		arbiter = "HighlyAvailableArbiter"
	)
	cases := []struct {
		name                 string
		controlPlaneReplicas int64
		arbiterReplicas      int64
		workerReplicas       int64
		mastersSchedulable   bool
		controlPlaneTopology string
		validInfraTopologies []string
	}{
		// Zero workers can mean UPI with externally provisioned workers, even
		// when the control plane remains schedulable.
		{"1+0 unschedulable", 1, 0, 0, false, single, []string{single, ha}},
		{"1+0 schedulable", 1, 0, 0, true, single, []string{single, ha}},
		{"1+1 unschedulable", 1, 0, 1, false, single, []string{single}},
		{"1+1 schedulable", 1, 0, 1, true, single, []string{ha}},
		{"1+2 unschedulable", 1, 0, 2, false, single, []string{ha}},
		{"1+2 schedulable", 1, 0, 2, true, single, []string{ha}},
		{"2+0 unschedulable", 2, 0, 0, false, dual, []string{single, ha}},
		{"2+0 schedulable", 2, 0, 0, true, dual, []string{ha}},
		{"2+1 unschedulable", 2, 0, 1, false, dual, []string{single}},
		{"2+1 schedulable", 2, 0, 1, true, dual, []string{ha}},
		{"2+2 unschedulable", 2, 0, 2, false, dual, []string{ha}},
		{"2+2 schedulable", 2, 0, 2, true, dual, []string{ha}},
		{"2+arbiter+0 unschedulable", 2, 1, 0, false, arbiter, []string{single, ha}},
		{"2+arbiter+0 schedulable", 2, 1, 0, true, arbiter, []string{ha}},
		{"2+arbiter+1 unschedulable", 2, 1, 1, false, arbiter, []string{single}},
		{"2+arbiter+1 schedulable", 2, 1, 1, true, arbiter, []string{ha}},
		{"2+arbiter+2 unschedulable", 2, 1, 2, false, arbiter, []string{ha}},
		{"2+arbiter+2 schedulable", 2, 1, 2, true, arbiter, []string{ha}},
		{"3+0 unschedulable", 3, 0, 0, false, ha, []string{single, ha}},
		{"3+0 schedulable", 3, 0, 0, true, ha, []string{ha}},
		{"3+1 unschedulable", 3, 0, 1, false, ha, []string{single}},
		{"3+1 schedulable", 3, 0, 1, true, ha, []string{ha}},
		{"3+2 unschedulable", 3, 0, 2, false, ha, []string{ha}},
		{"3+2 schedulable", 3, 0, 2, true, ha, []string{ha}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.InstallConfig{
				ControlPlane: &types.MachinePool{Replicas: ptr.To(tc.controlPlaneReplicas)},
				Arbiter:      &types.MachinePool{Replicas: ptr.To(tc.arbiterReplicas)},
				Compute:      []types.MachinePool{{Replicas: ptr.To(tc.workerReplicas)}},
			}
			// Check every topology value, including values only valid for the
			// control plane, typos, and explicitly empty fields.
			topologies := []string{single, dual, ha, arbiter, "External", "invalid", ""}
			for _, controlPlane := range topologies {
				for _, infrastructure := range topologies {
					t.Run(controlPlane+"/"+infrastructure, func(t *testing.T) {
						files := topologyManifests(tc.mastersSchedulable, controlPlane, infrastructure)
						err := validateTopology(files, config)
						validControlPlane := controlPlane == tc.controlPlaneTopology
						validInfra := slices.Contains(tc.validInfraTopologies, infrastructure)
						if validControlPlane && validInfra {
							assert.NoError(t, err)
						} else {
							assert.ErrorContains(t, err, files[1].Filename)
							if !validControlPlane {
								assert.ErrorContains(t, err, "status.controlPlaneTopology")
							}
							if !validInfra {
								assert.ErrorContains(t, err, "status.infrastructureTopology")
							}
						}
					})
				}
			}
		})
	}
}

func TestValidateTopologyWorkerPools(t *testing.T) {
	config := &types.InstallConfig{
		ControlPlane: &types.MachinePool{}, // Defaults to three replicas.
		Compute: []types.MachinePool{
			{Replicas: ptr.To(int64(1))},
			{}, // An omitted worker replica count contributes zero.
			{Replicas: ptr.To(int64(1))},
		},
	}
	err := validateTopology(topologyManifests(false, "HighlyAvailable", "SingleReplica"), config)
	assert.ErrorContains(t, err, "must be HighlyAvailable for 2 configured workers")
	assert.NoError(t, validateTopology(topologyManifests(false, "HighlyAvailable", "HighlyAvailable"), config))
}

func TestValidateTopologyManifestParsing(t *testing.T) {
	cases := []struct {
		name     string
		edit     func([]*asset.File) []*asset.File
		expected string
	}{
		{
			name: "renamed manifests",
			edit: func(files []*asset.File) []*asset.File {
				files[0].Filename = "manifests/custom-scheduler.yaml"
				files[1].Filename = "openshift/custom-infrastructure.yaml"
				return files
			},
			expected: "incompatible topology",
		},
		{
			name: "JSON manifest",
			edit: func(files []*asset.File) []*asset.File {
				files[0].Filename = "manifests/custom-scheduler.json"
				files[0].Data = []byte(`{"apiVersion":"config.openshift.io/v1","kind":"Scheduler","metadata":{"name":"cluster"},"spec":{"mastersSchedulable":true}}`)
				return files
			},
			expected: "incompatible topology",
		},
		{
			name: "missing scheduler",
			edit: func(files []*asset.File) []*asset.File { return files[1:] },
		},
		{
			name:     "missing infrastructure",
			edit:     func(files []*asset.File) []*asset.File { return files[:1] },
			expected: "missing Infrastructure cluster manifest",
		},
		{
			name: "unrelated resource",
			edit: func(files []*asset.File) []*asset.File {
				files[0].Data = []byte("apiVersion: example.com/v1\nkind: Scheduler\nmetadata:\n  name: cluster\nspec:\n  mastersSchedulable: true\n")
				return files
			},
		},
		{
			name: "invalid scheduler",
			edit: func(files []*asset.File) []*asset.File {
				files[0].Data = []byte(strings.ReplaceAll(string(files[0].Data), "mastersSchedulable: true", "mastersSchedulable: []"))
				return files
			},
			expected: "failed to unmarshal manifests/cluster-scheduler-02-config.yml",
		},
		{
			name: "invalid infrastructure",
			edit: func(files []*asset.File) []*asset.File {
				files[1].Data = []byte(strings.ReplaceAll(string(files[1].Data), `infrastructureTopology: "SingleReplica"`, "infrastructureTopology: []"))
				return files
			},
			expected: "failed to unmarshal manifests/cluster-infrastructure-02-config.yml",
		},
		{
			name: "omitted topology fields default to HA",
			edit: func(files []*asset.File) []*asset.File {
				files[1].Data = []byte("apiVersion: config.openshift.io/v1\nkind: Infrastructure\nmetadata:\n  name: cluster\nstatus: {}\n")
				return files
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.InstallConfig{
				ControlPlane: &types.MachinePool{Replicas: ptr.To(int64(3))},
				Compute:      []types.MachinePool{{Replicas: ptr.To(int64(1))}},
			}
			err := validateTopology(tc.edit(topologyManifests(true, "HighlyAvailable", "SingleReplica")), config)
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expected)
			}
		})
	}
}

// Exercise the manifest load and ignition generation boundary: edits to the
// scheduler must be checked against the loaded infrastructure, not regenerated
// defaults from the original install-config.
func TestGenerateConfigRejectsEditedTopology(t *testing.T) {
	files := topologyManifests(true, "HighlyAvailable", "SingleReplica")
	files = append(files, &asset.File{
		Filename: "manifests/cluster-config.yaml",
		Data:     []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: cluster-config-v1\n"),
	})
	fetcher := mock.NewMockFileFetcher(gomock.NewController(t))
	fetcher.EXPECT().FetchByPattern("manifests/*.yaml").Return(files[2:], nil)
	fetcher.EXPECT().FetchByPattern("manifests/*.yml").Return(files[:2], nil)
	fetcher.EXPECT().FetchByPattern("manifests/*.json").Return(nil, nil)
	loadedManifests := &manifests.Manifests{}
	found, err := loadedManifests.Load(fetcher)
	if !assert.NoError(t, err) || !assert.True(t, found) {
		return
	}

	parents := asset.Parents{}
	parents.Add(
		installconfig.MakeAsset(&types.InstallConfig{
			ControlPlane: &types.MachinePool{Replicas: ptr.To(int64(3))},
			Compute:      []types.MachinePool{{Replicas: ptr.To(int64(1))}},
		}),
		&installconfig.ClusterID{},
		&tls.BootstrapSSHKeyPair{},
		loadedManifests,
		&manifests.Openshift{},
	)
	common := &Common{}
	err = common.generateConfig(parents, &bootstrapTemplateData{})
	assert.ErrorContains(t, err, "incompatible topology")
	assert.Nil(t, common.Config, "invalid topology must be rejected before generating ignition")
}
