// Package google provides a cluster-destroyer for google clusters.
package google

import (
	"github.com/openshift/installer/pkg/destroy"
)

func init() {
	destroy.Registry["google"] = New
}
