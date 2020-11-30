package utils

import "fmt"

func BuildLabels(infraID string) map[string]string {
	return map[string]string{
		fmt.Sprintf("tenantcluster-%s-machine.openshift.io", infraID): "owned",
	}
}
