package conversion

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/agent"
)

// ConvertAgentConfig is modeled after the k8s conversion schemes, which is
// how deprecated values are upconverted.
// This updates the APIVersion to reflect the fact that we've internally
// upconverted.
func ConvertAgentConfig(config *agent.Config) error {
	// check that the version is convertible
	switch config.APIVersion {
	case agent.AgentConfigVersion, "v1alpha1":
		// works
	case "":
		return field.Required(field.NewPath("apiVersion"), "no version was provided")
	default:
		return field.Invalid(field.NewPath("apiVersion"), config.APIVersion, fmt.Sprintf("cannot upconvert from version %s", config.APIVersion))
	}

	config.APIVersion = agent.AgentConfigVersion
	return nil
}
