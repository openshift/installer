package flex

import (
	"errors"
	"fmt"

	v "github.com/IBM-Cloud/terraform-provider-ibm/version"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// TerraformProblem provides a type that holds standardized information
// suitable to problems that occur in the Terraform Provider code.
type TerraformProblem struct {
	*core.IBMProblem

	Resource  string
	Operation string
}

// GetID returns a hash value computed from stable fields in the
// TerraformProblem instance, including Resource and Operation.
func (e *TerraformProblem) GetID() string {
	return core.CreateIDHash("terraform", e.GetBaseSignature(), e.Resource, e.Operation)
}

// GetConsoleMessage returns the fields of the problem that
// are relevant to a user, formatted as a YAML string.
func (e *TerraformProblem) GetConsoleMessage() string {
	return core.ComputeConsoleMessage(e)
}

// GetConsoleMessage returns the fields of the problem that
// are relevant to a developer, formatted as a YAML string.
func (e *TerraformProblem) GetDebugMessage() string {
	return core.ComputeDebugMessage(e)
}

func (e *TerraformProblem) GetConsoleOrderedMaps() *core.OrderedMaps {
	orderedMaps := core.NewOrderedMaps()

	orderedMaps.Add("id", e.GetID())
	orderedMaps.Add("summary", e.Summary)
	orderedMaps.Add("severity", e.Severity)
	orderedMaps.Add("resource", e.Resource)
	orderedMaps.Add("operation", e.Operation)
	orderedMaps.Add("component", e.Component)

	return orderedMaps
}

func (e *TerraformProblem) GetDebugOrderedMaps() *core.OrderedMaps {
	orderedMaps := e.GetConsoleOrderedMaps()

	var orderableCausedBy core.OrderableProblem
	if errors.As(e.GetCausedBy(), &orderableCausedBy) {
		orderedMaps.Add("caused_by", orderableCausedBy.GetDebugOrderedMaps().GetMaps())
	}

	return orderedMaps
}

// GetDiag returns a new Diagnostics object using the console
// message as the summary. It is used to create a Diagnostics
// object from a TerraformProblem in the resource/data source code.
func (e *TerraformProblem) GetDiag() diag.Diagnostics {
	return diag.Errorf("%s", e.GetConsoleMessage())
}

// TerraformErrorf creates and returns a new instance of `TerraformProblem`
// with "error" level severity and a blank discriminator - the "caused by"
// error is used to ensure uniqueness. This is a convenience function to
// use when creating a new TerraformProblem instance from an error that
// came from the SDK.
func TerraformErrorf(err error, summary, resource, operation string) *TerraformProblem {
	return DiscriminatedTerraformErrorf(err, summary, resource, operation, "")
}

// DiscriminatedTerraformErrorf creates and returns a new instance
// of `TerraformProblem` with "error" level severity that contains
// a discriminator used to make the instance unique relative to
// other problem scenarios in the same resource/operation.
func DiscriminatedTerraformErrorf(err error, summary, resource, operation, discriminator string) *TerraformProblem {
	return &TerraformProblem{
		IBMProblem: core.IBMErrorf(err, getComponentInfo(), summary, discriminator),
		Resource:   resource,
		Operation:  operation,
	}
}

func getComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent("github.com/IBM-Cloud/terraform-provider-ibm", v.Version)
}

// FmtErrorf wraps `fmt.Errorf(format string, a ...interface{}) error`
// and attempts to return a TerraformProblem instance instead of a
// plain error instance, if an error object is found among the arguments
func FmtErrorf(format string, a ...interface{}) error {
	intendedError := fmt.Errorf(format, a...)

	var err error
	for _, arg := range a {
		// Look for an error instance among the arguments.

		if errArg, ok := arg.(error); ok {
			err = errArg
		} else if ser, ok := arg.(*ServiceErrorResponse); ok {
			// Deal with the "ServiceErrorResponse" type, which
			// wraps errors in some of the handwritten code.
			err = ser.Error
		}

		if err != nil {
			var tfError *TerraformProblem
			if !errors.As(err, &tfError) {
				tfError = TerraformErrorf(err, err.Error(), "", "")
			}

			tfError.Summary = intendedError.Error()
			return tfError
		}
	}

	return intendedError
}
