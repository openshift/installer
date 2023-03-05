package aws

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	operatorv1 "github.com/openshift/api/operator/v1"
)

// ClusterCSIDriverConfig is the AWS config for the cluster CSI driver.
type ClusterCSIDriverConfig struct {
	KMSKeyARN string
}

// YAML generates the cluster CSI driver config for the AWS platform.
func (params ClusterCSIDriverConfig) YAML() ([]byte, error) {
	obj := &operatorv1.ClusterCSIDriver{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1.GroupVersion.String(),
			Kind:       "ClusterCSIDriver",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: string(operatorv1.AWSEBSCSIDriver),
		},
		Spec: operatorv1.ClusterCSIDriverSpec{
			DriverConfig: operatorv1.CSIDriverConfigSpec{
				DriverType: operatorv1.AWSDriverType,
				AWS: &operatorv1.AWSCSIDriverConfigSpec{
					KMSKeyARN: params.KMSKeyARN,
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
