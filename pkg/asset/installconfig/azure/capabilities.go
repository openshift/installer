package azure

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types/azure"
)

// GetHyperVGenerationVersion returns a HyperVGeneration version compatible with that of the image's. If imageHyperVGen is empty, it returns the highest supported version.
func GetHyperVGenerationVersion(capabilities map[string]string, imageHyperVGen string) (string, error) {
	generations, err := GetHyperVGenerationVersions(capabilities)
	if err != nil {
		return "", err
	}
	// If there is a version compatible with the VM image, return it
	if imageHyperVGen != "" && generations.Has(imageHyperVGen) {
		return imageHyperVGen, nil
	} else if generations.Len() > 0 { // otherwise, return the highest version available
		return sets.List(generations)[generations.Len()-1], nil
	}
	if generations.Has("V2") {
		return "V2", nil
	}
	return "V1", nil
}

// GetHyperVGenerationVersions returns all the HyperVGeneration versions supported by the instance type according to its capabilities as a string set V = {"V1", "V2", ...}
func GetHyperVGenerationVersions(capabilities map[string]string) (sets.Set[string], error) {
	if val, ok := capabilities["HyperVGenerations"]; ok {
		generations := sets.New[string]()
		for _, g := range strings.Split(val, ",") {
			g = strings.TrimSpace(g)
			g = strings.ToUpper(g)
			generations.Insert(g)
		}
		return generations, nil
	}
	return nil, fmt.Errorf("unable to determine HyperVGeneration version")
}

// GetVMNetworkingCapability returns true if Accelerated networking is supported by the instance type according to its capabilities or false, otherwise
func GetVMNetworkingCapability(capabilities map[string]string) bool {
	val, ok := capabilities[string(azure.AcceleratedNetworkingEnabled)]
	return ok && strings.EqualFold(val, "True")
}
