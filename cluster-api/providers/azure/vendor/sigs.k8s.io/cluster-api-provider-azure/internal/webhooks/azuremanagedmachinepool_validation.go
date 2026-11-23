/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhooks

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	clusterctlv1alpha3 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
)

// validateLastSystemNodePool is used to check if the existing system node pool is the last system node pool.
// If it is a last system node pool it cannot be deleted or mutated to user node pool as AKS expects min 1 system node pool.
func validateLastSystemNodePool(cli client.Client, labels map[string]string, namespace string, annotations map[string]string) error {
	ctx := context.Background()

	// Fetch the Cluster.
	clusterName, ok := labels[clusterv1.ClusterNameLabel]
	if !ok {
		return nil
	}

	ownerCluster := &clusterv1.Cluster{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      clusterName,
	}

	if err := cli.Get(ctx, key, ownerCluster); err != nil {
		return err
	}

	if !ownerCluster.DeletionTimestamp.IsZero() {
		return nil
	}

	// checking if this AzureManagedMachinePool is going to be deleted for clusterctl move operation
	if _, ok := annotations[clusterctlv1alpha3.DeleteForMoveAnnotation]; ok {
		return nil
	}

	opt1 := client.InNamespace(namespace)
	opt2 := client.MatchingLabels(map[string]string{
		clusterv1.ClusterNameLabel: clusterName,
		infrav1.LabelAgentPoolMode: string(infrav1.NodePoolModeSystem),
	})

	ammpList := &infrav1.AzureManagedMachinePoolList{}
	if err := cli.List(ctx, ammpList, opt1, opt2); err != nil {
		return err
	}

	if len(ammpList.Items) <= 1 {
		return errors.New("AKS Cluster must have at least one system pool")
	}
	return nil
}

func validateMaxPods(maxPods *int, fldPath *field.Path) error {
	if maxPods != nil {
		if ptr.Deref(maxPods, 0) < 10 || ptr.Deref(maxPods, 0) > 250 {
			return field.Invalid(
				fldPath,
				maxPods,
				"MaxPods must be between 10 and 250")
		}
	}

	return nil
}

func validateOSType(mode string, osType *string, fldPath *field.Path) error {
	if mode == string(infrav1.NodePoolModeSystem) {
		if osType != nil && *osType != infrav1.LinuxOS {
			return field.Forbidden(
				fldPath,
				"System node pooll must have OSType 'Linux'")
		}
	}

	return nil
}

func validateMPName(mpName string, specName *string, osType *string, fldPath *field.Path) error {
	var name *string
	var fieldNameMessage string
	if specName == nil || *specName == "" {
		name = &mpName
		fieldNameMessage = "when spec.name is empty, metadata.name"
	} else {
		name = specName
		fieldNameMessage = "spec.name"
	}

	if err := validateNameLength(osType, name, fieldNameMessage, fldPath); err != nil {
		return err
	}
	return validateNamePattern(name, fieldNameMessage, fldPath)
}

func validateNameLength(osType *string, name *string, fieldNameMessage string, fldPath *field.Path) error {
	if osType != nil && *osType == infrav1.WindowsOS &&
		name != nil && len(*name) > 6 {
		return field.Invalid(
			fldPath,
			name,
			fmt.Sprintf("For OSType Windows, %s can not be longer than 6 characters.", fieldNameMessage))
	} else if (osType == nil || *osType == infrav1.LinuxOS) &&
		(name != nil && len(*name) > 12) {
		return field.Invalid(
			fldPath,
			osType,
			fmt.Sprintf("For OSType Linux, %s can not be longer than 12 characters.", fieldNameMessage))
	}
	return nil
}

func validateNamePattern(name *string, fieldNameMessage string, fldPath *field.Path) error {
	if name == nil || *name == "" {
		return nil
	}

	if !unicode.IsLower(rune((*name)[0])) {
		return field.Invalid(
			fldPath,
			name,
			fmt.Sprintf("%s must begin with a lowercase letter.", fieldNameMessage))
	}

	for _, char := range *name {
		if !unicode.IsLower(char) && !unicode.IsNumber(char) {
			return field.Invalid(
				fldPath,
				name,
				fmt.Sprintf("%s may only contain lowercase alphanumeric characters.", fieldNameMessage))
		}
	}
	return nil
}

func validateNodeLabels(nodeLabels map[string]string, fldPath *field.Path) error {
	for key := range nodeLabels {
		if azureutil.IsAzureSystemNodeLabelKey(key) {
			return field.Invalid(
				fldPath,
				key,
				fmt.Sprintf("Node pool label key must not start with %s", azureutil.AzureSystemNodeLabelPrefix))
		}
	}

	return nil
}

func validateNodePublicIPPrefixID(nodePublicIPPrefixID *string, fldPath *field.Path) error {
	if nodePublicIPPrefixID != nil && !validNodePublicPrefixID.MatchString(*nodePublicIPPrefixID) {
		return field.Invalid(
			fldPath,
			nodePublicIPPrefixID,
			fmt.Sprintf("resource ID must match %q", validNodePublicPrefixID.String()))
	}
	return nil
}

func validateEnableNodePublicIP(enableNodePublicIP *bool, nodePublicIPPrefixID *string, fldPath *field.Path) error {
	if (enableNodePublicIP == nil || !*enableNodePublicIP) &&
		nodePublicIPPrefixID != nil {
		return field.Invalid(
			fldPath,
			enableNodePublicIP,
			"must be set to true when NodePublicIPPrefixID is set")
	}
	return nil
}

func validateMPSubnetName(subnetName *string, fldPath *field.Path) error {
	if subnetName != nil {
		subnetRegex := "^[a-zA-Z0-9][a-zA-Z0-9._-]{0,78}[a-zA-Z0-9]$"
		regex := regexp.MustCompile(subnetRegex)
		if success := regex.MatchString(ptr.Deref(subnetName, "")); !success {
			return field.Invalid(fldPath, subnetName,
				fmt.Sprintf("name of subnet doesn't match regex %s", subnetRegex))
		}
	}
	return nil
}

// validateKubeletConfig enforces the AKS API configuration for KubeletConfig.
// See:  https://learn.microsoft.com/en-us/azure/aks/custom-node-configuration.
func validateKubeletConfig(kubeletConfig *infrav1.KubeletConfig, fldPath *field.Path) error {
	var allowedUnsafeSysctlsPatterns = []string{
		`^kernel\.shm.+$`,
		`^kernel\.msg.+$`,
		`^kernel\.sem$`,
		`^fs\.mqueue\..+$`,
		`^net\..+$`,
	}
	if kubeletConfig != nil {
		if kubeletConfig.CPUCfsQuotaPeriod != nil {
			if !strings.HasSuffix(ptr.Deref(kubeletConfig.CPUCfsQuotaPeriod, ""), "ms") {
				return field.Invalid(
					fldPath.Child("CPUfsQuotaPeriod"),
					kubeletConfig.CPUCfsQuotaPeriod,
					"must be a string value in milliseconds with a 'ms' suffix, e.g., '100ms'")
			}
		}
		if kubeletConfig.ImageGcHighThreshold != nil && kubeletConfig.ImageGcLowThreshold != nil {
			if ptr.Deref(kubeletConfig.ImageGcLowThreshold, 0) > ptr.Deref(kubeletConfig.ImageGcHighThreshold, 0) {
				return field.Invalid(
					fldPath.Child("ImageGcLowThreshold"),
					kubeletConfig.ImageGcLowThreshold,
					fmt.Sprintf("must not be greater than ImageGcHighThreshold, ImageGcLowThreshold=%d, ImageGcHighThreshold=%d",
						ptr.Deref(kubeletConfig.ImageGcLowThreshold, 0), ptr.Deref(kubeletConfig.ImageGcHighThreshold, 0)))
			}
		}
		for _, val := range kubeletConfig.AllowedUnsafeSysctls {
			var hasMatch bool
			for _, p := range allowedUnsafeSysctlsPatterns {
				if m, _ := regexp.MatchString(p, val); m {
					hasMatch = true
					break
				}
			}
			if !hasMatch {
				return field.Invalid(
					fldPath.Child("AllowedUnsafeSysctls"),
					kubeletConfig.AllowedUnsafeSysctls,
					fmt.Sprintf("%s is not a supported AllowedUnsafeSysctls configuration", val))
			}
		}
	}
	return nil
}

// validateLinuxOSConfig enforces AKS API configuration for Linux OS custom configuration
// See: https://learn.microsoft.com/en-us/azure/aks/custom-node-configuration#linux-os-custom-configuration for detailed information.
func validateLinuxOSConfig(linuxOSConfig *infrav1.LinuxOSConfig, kubeletConfig *infrav1.KubeletConfig, fldPath *field.Path) error {
	var errs []error
	if linuxOSConfig == nil {
		return nil
	}

	if linuxOSConfig.SwapFileSizeMB != nil {
		if kubeletConfig == nil || ptr.Deref(kubeletConfig.FailSwapOn, true) {
			errs = append(errs, field.Invalid(
				fldPath.Child("SwapFileSizeMB"),
				linuxOSConfig.SwapFileSizeMB,
				"KubeletConfig.FailSwapOn must be set to false to enable swap file on nodes"))
		}
	}

	if linuxOSConfig.Sysctls != nil && linuxOSConfig.Sysctls.NetIpv4IPLocalPortRange != nil {
		// match numbers separated by a space
		portRangeRegex := `^[0-9]+ [0-9]+$`
		portRange := *linuxOSConfig.Sysctls.NetIpv4IPLocalPortRange

		match, matchErr := regexp.MatchString(portRangeRegex, portRange)
		if matchErr != nil {
			errs = append(errs, matchErr)
		}
		if !match {
			errs = append(errs, field.Invalid(
				fldPath.Child("NetIpv4IpLocalPortRange"),
				linuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
				"LinuxOSConfig.Sysctls.NetIpv4IpLocalPortRange must be of the format \"<int> <int>\""))
		} else {
			ports := strings.Split(portRange, " ")
			firstPort, _ := strconv.Atoi(ports[0])
			lastPort, _ := strconv.Atoi(ports[1])

			if firstPort < 1024 || firstPort > 60999 {
				errs = append(errs, field.Invalid(
					fldPath.Child("NetIpv4IpLocalPortRange", "First"),
					linuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
					fmt.Sprintf("first port of NetIpv4IpLocalPortRange=%d must be in between [1024 - 60999]", firstPort)))
			}

			if lastPort < 32768 || lastPort > 65000 {
				errs = append(errs, field.Invalid(
					fldPath.Child("NetIpv4IpLocalPortRange", "Last"),
					linuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
					fmt.Sprintf("last port of NetIpv4IpLocalPortRange=%d must be in between [32768 -65000]", lastPort)))
			}

			if firstPort > lastPort {
				errs = append(errs, field.Invalid(
					fldPath.Child("NetIpv4IpLocalPortRange", "First"),
					linuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
					fmt.Sprintf("first port of NetIpv4IpLocalPortRange=%d cannot be greater than last port of NetIpv4IpLocalPortRange=%d", firstPort, lastPort)))
			}
		}
	}
	return kerrors.NewAggregate(errs)
}
