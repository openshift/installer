// Package machineconfigpool generates MachineConfigPool manifests for custom
// compute pools. It mirrors the machineconfig/ helper package: pure functions,
// no Asset interface.
package machineconfigpool

import (
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset"
)

const (
	machineConfigPoolFileName = "99_openshift-machineconfig_%s-mcp.yaml"
)

// ForCustomPool builds a MachineConfigPool for a user-defined compute pool.
// The pool inherits base worker MachineConfigs via machineConfigSelector and
// MCO uses nodeSelector to assign nodes to this pool on day 2 (driven by the
// initial MCS request to /config/<poolName>).
func ForCustomPool(poolName string) *mcfgv1.MachineConfigPool {
	return &mcfgv1.MachineConfigPool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfigPool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: poolName,
		},
		Spec: mcfgv1.MachineConfigPoolSpec{
			MachineConfigSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "machineconfiguration.openshift.io/role",
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{"worker", poolName},
					},
				},
			},
			NodeSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"node-role.kubernetes.io/" + poolName: "",
				},
			},
		},
	}
}

// Manifest serialises a MachineConfigPool for the given custom pool into an
// asset file ready to be written to the openshift/ manifest directory.
func Manifest(poolName, directory string) (*asset.File, error) {
	data, err := yaml.Marshal(ForCustomPool(poolName))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MachineConfigPool for pool %q: %w", poolName, err)
	}
	return &asset.File{
		Filename: filepath.Join(directory, fmt.Sprintf(machineConfigPoolFileName, poolName)),
		Data:     data,
	}, nil
}
