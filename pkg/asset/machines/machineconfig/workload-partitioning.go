package machineconfig

import (
	"bytes"
	"encoding/json"
	"fmt"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/types"
)

// CrioWorkloadDropinContents generates the content expected by Cri-O for workload partitioning
//
// Example output:
// [crio.runtime.workloads.management]
// label             = "management.workload.openshift.io/cores"
// annotation_prefix = "io.openshift.workload.management"
// resources         = { "cpu" = "", "cpuset" = "0-1", }
func CrioWorkloadDropinContents(workloads []types.Workload) (string, error) {
	crioIni := ini.Empty()
	for _, w := range workloads {
		section := crioIni.Section(fmt.Sprintf("crio.runtime.workloads.%s", w.Name))
		err := section.ReflectFrom(&crioWorkloadCfg{
			Label:            fmt.Sprintf(`"workload.openshift.io/%s"`, w.Name),
			AnnotationPrefix: fmt.Sprintf(`"io.openshift.workload.%s"`, w.Name),
			Resources:        fmt.Sprintf(`{ "cpu" = "", "cpuset" = "%s", }`, w.CPUIDs),
		})
		if err != nil {
			return "", errors.Wrapf(err, "could not reflect %q structure to INI", w.Name)
		}
	}
	crioBuf := new(bytes.Buffer)
	if _, err := crioIni.WriteTo(crioBuf); err != nil {
		return "", errors.Wrap(err, "could not write INI to buffer")
	}
	return crioBuf.String(), nil
}

type crioWorkloadCfg struct {
	Label            string `ini:"label"`
	AnnotationPrefix string `ini:"annotation_prefix"`
	Resources        string `ini:"resources"`
}

// KubeletWorkloadDropinContents generates the content expected by Kubelet for workload partitioning
//
// Example output:
// {
//   "management": {
//     "cpuset": "0-1"
//   }
// }
func KubeletWorkloadDropinContents(workloads []types.Workload) (string, error) {
	kubeletWorkload := map[string]kubeletWorkloadEntry{}
	for _, w := range workloads {
		kubeletWorkload[string(w.Name)] = kubeletWorkloadEntry{Cpuset: w.CPUIDs}
	}
	kubeletCfg, err := json.MarshalIndent(kubeletWorkload, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "could not marshall JSON")
	}
	return string(kubeletCfg), nil
}

type kubeletWorkloadEntry struct {
	Cpuset string `json:"cpuset"`
}

// ForWorkloads creates the MachineConfig that configures Cri-O and
// Kubelet with the workload partition settings from install-config.yaml
func ForWorkloads(workloads []types.Workload, role string) (*mcfgv1.MachineConfig, error) {
	crioCfg, err := CrioWorkloadDropinContents(workloads)
	if err != nil {
		return nil, err
	}
	kubeletCfg, err := KubeletWorkloadDropinContents(workloads)
	if err != nil {
		return nil, err
	}
	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				ignition.FileFromString(
					"/etc/crio/crio.conf.d/01-workload-partitioning",
					"root", 0644, crioCfg),
				ignition.FileFromString(
					"/etc/kubernetes/workload-pinning",
					"root", 0644, kubeletCfg),
			},
		},
	}
	rawExt, err := ignition.ConvertToRawExtension(ignConfig)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert to raw extension")
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("02-%s-workload-partitioning", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}
