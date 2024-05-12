package ibi

import (
	"github.com/openshift-kni/lifecycle-agent/api/ibiconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ImageBasedInstallConfigVersion is the version supported by this package.
const ImageBasedInstallConfigVersion = "v1beta1"

// Config or aka ImageBasedInstallConfig is the API for specifying configuration
// for the image-based installer.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	ibiconfig.IBIPrepareConfig `json:",inline"`
}
