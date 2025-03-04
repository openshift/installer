package none

import "github.com/openshift/installer/pkg/types/common"

// Platform stores any global configuration used for generic
// platforms.
type Platform struct {
	// FencingCredentials stores the information about a baremetal host's management controller.
	// +optional
	FencingCredentials []*common.FencingCredential `json:"fencingCredentials,omitempty"`
}
