package ibmcloud

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	operatorv1 "github.com/openshift/api/operator/v1"
)

// ClusterCSIDriverConfig is the IBM Cloud config for the cluster CSI Driver.
type ClusterCSIDriverConfig struct {
	EncryptionKeyCRN string
}

// YAML generates the cluster CSI Driver config for the IBM Cloud platform.
func (params ClusterCSIDriverConfig) YAML() ([]byte, error) {
	obj := &operatorv1.ClusterCSIDriver{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1.GroupVersion.String(),
			Kind:       "ClusterCSIDriver",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: string(operatorv1.IBMVPCBlockCSIDriver),
		},
		Spec: operatorv1.ClusterCSIDriverSpec{
			DriverConfig: operatorv1.CSIDriverConfigSpec{
				DriverType: operatorv1.IBMCloudDriverType,
				IBMCloud: &operatorv1.IBMCloudCSIDriverConfigSpec{
					EncryptionKeyCRN: params.EncryptionKeyCRN,
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
