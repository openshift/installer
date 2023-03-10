package gcp

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	operatorv1 "github.com/openshift/api/operator/v1"
)

// ClusterCSIDriverConfig is the GCP config for the cluster CSI Driver.
type ClusterCSIDriverConfig struct {
	Name      string
	KeyRing   string
	ProjectID string
	Location  string
}

// YAML generates the cluster CSI driver config for the GCP platform.
func (params ClusterCSIDriverConfig) YAML() ([]byte, error) {
	obj := &operatorv1.ClusterCSIDriver{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1.GroupVersion.String(),
			Kind:       "ClusterCSIDriver",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: string(operatorv1.GCPPDCSIDriver),
		},
		Spec: operatorv1.ClusterCSIDriverSpec{
			DriverConfig: operatorv1.CSIDriverConfigSpec{
				DriverType: operatorv1.GCPDriverType,
				GCP: &operatorv1.GCPCSIDriverConfigSpec{
					KMSKey: &operatorv1.GCPKMSKeyReference{
						Name:      params.Name,
						KeyRing:   params.KeyRing,
						ProjectID: params.ProjectID,
						Location:  params.Location,
					},
				},
			},
			OperatorSpec: operatorv1.OperatorSpec{
				ManagementState: operatorv1.Managed,
			},
		},
	}

	configData, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return configData, nil
}
