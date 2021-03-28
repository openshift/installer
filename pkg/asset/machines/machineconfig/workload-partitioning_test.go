package machineconfig

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
)

func expectedCrioCfg(workloads []types.Workload) string {
	parts := []string{}
	for _, w := range workloads {
		parts = append(parts, fmt.Sprintf(`[crio.runtime.workloads.%[1]s]
label             = "workload.openshift.io/%[1]s"
annotation_prefix = "io.openshift.workload.%[1]s"
resources         = { "cpu" = "", "cpuset" = "%s", }
`, w.Name, w.CPUIDs))
	}
	return strings.Join(parts, "\n") + "\n"
}

func expectedKubeletCfg(workloads []types.Workload) string {
	parts := []string{}
	for _, w := range workloads {
		parts = append(parts, fmt.Sprintf(`  %q: {
    "cpuset": %q
  }`, w.Name, w.CPUIDs))
	}
	return "{\n" + strings.Join(parts, ",\n") + "\n}"
}

func TestWorkloading(t *testing.T) {
	cases := []struct {
		workloads []types.Workload
		role      string
	}{
		{
			workloads: []types.Workload{
				{
					Name:   types.ManagementWorkload,
					CPUIDs: "0-1",
				},
			},
			role: "master",
		},
		{
			workloads: []types.Workload{
				{
					Name:   types.ManagementWorkload,
					CPUIDs: "0-1",
				},
				{
					Name:   "secondary",
					CPUIDs: "50-51",
				},
			},
			role: "master",
		},
	}

	t.Run("test", func(t *testing.T) {
		for _, tc := range cases {
			expectedCrioCfg := expectedCrioCfg(tc.workloads)
			crioCfg, err := CrioWorkloadDropinContents(tc.workloads)
			assert.Equal(t, err, nil, "No err")
			assert.Equal(t, expectedCrioCfg, crioCfg)

			expectedKubeletCfg := expectedKubeletCfg(tc.workloads)
			kubeletCfg, err := KubeletWorkloadDropinContents(tc.workloads)
			assert.Equal(t, err, nil, "No err")
			assert.Equal(t, expectedKubeletCfg, kubeletCfg)

			result, err := ForWorkloads(tc.workloads, tc.role)
			assert.Equal(t, err, nil, "No err")
			assert.Equal(t, tc.role, result.ObjectMeta.Labels["machineconfiguration.openshift.io/role"])

			cfg := igntypes.Config{}
			err = json.Unmarshal(result.Spec.Config.Raw, &cfg)
			assert.Equal(t, err, nil, "No err")
			files := cfg.Storage.Files
			assert.Len(t, files, 2, "Two files in the machineconfig")
			assert.Equal(t, "/etc/crio/crio.conf.d/01-workload-partitioning", files[0].Path, "Cri-O drop-in is present")
			assert.Equal(t, dataurl.EncodeBytes([]byte(expectedCrioCfg)), *files[0].Contents.Source)
			assert.Equal(t, "/etc/kubernetes/workload-pinning", files[1].Path, "Kubelet drop-in is present")
			assert.Equal(t, dataurl.EncodeBytes([]byte(expectedKubeletCfg)), *files[1].Contents.Source)
		}
	})
}
