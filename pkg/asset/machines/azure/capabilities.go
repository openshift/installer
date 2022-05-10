package azure

import (
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types/azure"
	"k8s.io/apimachinery/pkg/util/sets"
)

func getHyperVGenerationVersion(capabilities map[string]string) (string, error) {
	if val, ok := capabilities["HyperVGenerations"]; ok {
		generations := sets.NewString()
		for _, g := range strings.Split(val, ",") {
			g = strings.TrimSpace(g)
			g = strings.ToUpper(g)
			generations.Insert(g)
		}
		if generations.Has("V2") {
			return "V2", nil
		}
		return "V1", nil
	}
	return "", fmt.Errorf("unable to determine HyperVGeneration version")
}

func getVMNetworkingCapability(capabilities map[string]string) bool {
	val, ok := capabilities[string(azure.AcceleratedNetworkingEnabled)]
	if !ok {
		return false
	}
	return strings.EqualFold(val, "True")
}

func getUltraSSDCapability(capabilities map[string]string) (bool, error) {
	val, ok := capabilities["UltraSSDAvailable"]
	if !ok {
		return false, fmt.Errorf("unable to determine ultra ssd capability")
	}
	return strings.EqualFold(val, "True"), nil
}
