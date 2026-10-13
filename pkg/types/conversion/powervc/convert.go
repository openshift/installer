package powervc

import (
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervc"
)

// ConvertPowerVCInstallConfig takes a PowerVC install config and creates an underlying
// OpenStack confing.
func ConvertPowerVCInstallConfig(c *types.InstallConfig) error {
	if c.Platform.Name() != powervc.Name {
		return fmt.Errorf("convert expected powervc, got %s", c.Platform.Name())
	}

	// Marshal the PowerVC platform structure into yaml
	powerData, err := yaml.Marshal(c.Platform)
	if err != nil {
		return err
	}

	// Cut off the powervc: prefix
	powerString := string(powerData)
	idx := strings.Index(powerString, ":")
	if idx == -1 {
		return fmt.Errorf("convert could not find a colon in powervc platform")
	}
	openString := powerString[idx+1:]

	// Unmarshal the PowerVC yaml data into the OpenStack platform structure
	err = yaml.UnmarshalStrict([]byte(openString), &c.Platform.OpenStack, yaml.DisallowUnknownFields)
	if err != nil {
		return err
	}

	// Marshal the PowerVC control plane structure into yaml
	powerData, err = yaml.Marshal(c.ControlPlane.Platform.PowerVC)
	if err != nil {
		return err
	}

	// Unmarshal the PowerVC yaml data into the OpenStack control plane structure
	err = yaml.UnmarshalStrict(powerData, &c.ControlPlane.Platform.OpenStack, yaml.DisallowUnknownFields)
	if err != nil {
		return err
	}

	for idx := range c.Compute {
		// Marshal the PowerVC compute structure into yaml
		powerData, err = yaml.Marshal(c.Compute[idx].Platform.PowerVC)
		if err != nil {
			return err
		}

		// Unmarshal the PowerVC yaml data into the OpenStack compute structure
		err = yaml.UnmarshalStrict(powerData, &c.Compute[idx].Platform.OpenStack, yaml.DisallowUnknownFields)
		if err != nil {
			return err
		}
	}

	return nil
}
