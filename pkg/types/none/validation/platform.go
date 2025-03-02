package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/common"
	"github.com/openshift/installer/pkg/types/none"
)

// ValidateFencingCredentials checks that the provided fencing credentials are valid.
func ValidateFencingCredentials(fencingCredentials []*none.FencingCredential, fldPath *field.Path) (errors field.ErrorList) {
	return common.ValidateUniqueAndRequiredFields(fencingCredentials, fldPath, func([]byte) bool { return false }, "fencingCredentials")
}
