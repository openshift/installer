package azure

import (
	"encoding/json"
	"testing"

	machineapi "github.com/openshift/api/machine/v1beta1"

	"github.com/openshift/installer/pkg/types/azure"
)

func FuzzTFVars(f *testing.F) {
	// No control-plane machine used to panic on MasterConfigs[0].
	f.Add(string(azure.PublicCloud), false, true, "", false, false, false)
	f.Add(string(azure.PublicCloud), true, false, "", false, false, false)
	f.Add(string(azure.PublicCloud), true, true, "", false, false, true)
	// Security type set without the matching settings object used to panic.
	f.Add(string(azure.PublicCloud), true, true, string(machineapi.SecurityTypesConfidentialVM), false, false, false)
	f.Add(string(azure.PublicCloud), true, true, string(machineapi.SecurityTypesTrustedLaunch), false, false, false)
	f.Add(string(azure.PublicCloud), true, true, string(machineapi.SecurityTypesConfidentialVM), true, false, false)
	f.Add(string(azure.PublicCloud), true, true, string(machineapi.SecurityTypesTrustedLaunch), false, true, false)
	f.Add(string(azure.StackCloud), true, true, "", false, false, false)
	f.Add("no-such-cloud", true, true, "", false, false, false)

	f.Fuzz(func(t *testing.T, cloud string, hasMaster, hasWorker bool, securityType string, withConfidential, withTrusted, nilMaster bool) {
		sources := tfVarsSources(cloud, hasMaster, hasWorker, securityType, withConfidential, withTrusted, nilMaster)
		out, err := TFVars(sources)
		missingMachine := !hasMaster || !hasWorker || nilMaster
		if missingMachine {
			if err == nil || out != nil {
				t.Fatalf("missing machine configuration returned %q, %v", out, err)
			}
			return
		}
		if !supportedCloud(azure.CloudEnvironment(cloud)) {
			if err == nil || out != nil {
				t.Fatalf("cloud %q returned %q, %v", cloud, out, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("TFVars(%q) error: %v", cloud, err)
		}
		var decoded map[string]any
		if jsonErr := json.Unmarshal(out, &decoded); jsonErr != nil {
			t.Fatalf("TFVars output is not JSON: %v\n%s", jsonErr, out)
		}
		if decoded["azure_environment"] == nil && cloud != string(azure.StackCloud) {
			t.Fatalf("TFVars output missing azure_environment: %s", out)
		}
	})
}

func tfVarsSources(cloud string, hasMaster, hasWorker bool, securityType string, withConfidential, withTrusted, nilMaster bool) TFVarsSources {
	sources := TFVarsSources{CloudName: azure.CloudEnvironment(cloud), InfrastructureName: "infra"}
	if nilMaster {
		sources.MasterConfigs = []*machineapi.AzureMachineProviderSpec{nil}
	} else if hasMaster {
		spec := &machineapi.AzureMachineProviderSpec{Location: "eastus", VMSize: "Standard_D8s_v3", Subnet: "master"}
		if securityType != "" {
			spec.SecurityProfile = &machineapi.SecurityProfile{}
			spec.SecurityProfile.Settings.SecurityType = machineapi.SecurityTypes(securityType)
			if withConfidential {
				spec.SecurityProfile.Settings.ConfidentialVM = &machineapi.ConfidentialVM{}
			}
			if withTrusted {
				spec.SecurityProfile.Settings.TrustedLaunch = &machineapi.TrustedLaunch{}
			}
		}
		sources.MasterConfigs = []*machineapi.AzureMachineProviderSpec{spec}
	}
	if hasWorker {
		sources.WorkerConfigs = []*machineapi.AzureMachineProviderSpec{{Subnet: "worker"}}
	}
	return sources
}

func supportedCloud(cloud azure.CloudEnvironment) bool {
	switch cloud {
	case azure.PublicCloud, azure.USGovernmentCloud, azure.ChinaCloud, azure.GermanCloud, azure.StackCloud:
		return true
	default:
		return false
	}
}
