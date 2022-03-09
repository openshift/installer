package manifests

import (
	"fmt"
	"os"

	"github.com/go-openapi/swag"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/thoas/go-funk"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func GetPullSecret() string {
	secretData, err := os.ReadFile("./data/manifests/pull-secret.yaml")
	if err != nil {
		fmt.Errorf("Error reading pull secret: %w", err)
	}
	var secret corev1.Secret
	if err := yaml.Unmarshal(secretData, &secret); err != nil {
		fmt.Errorf("Error unmarshalling pull secret: %w", err)
	}
	pullSecret := secret.StringData[".dockerconfigjson"]
	return pullSecret
}

func getClusterDeployment() hivev1.ClusterDeployment {
	cdData, err := os.ReadFile("./data/manifests/cluster-deployment.yaml")
	if err != nil {
		fmt.Errorf("Error reading cluster deployment CR: %w", err)
	}
	var cd hivev1.ClusterDeployment
	if err := yaml.Unmarshal(cdData, &cd); err != nil {
		fmt.Errorf("Error unmarshalling cluster deployment CR: %w", err)
	}
	return cd
}

func getAgentClusterInstall() hiveext.AgentClusterInstall {
	aciData, err := os.ReadFile("./data/manifests/agent-cluster-install.yaml")
	if err != nil {
		fmt.Errorf("Error reading AgentClusterInstall CR: %w", err)
	}
	var aci hiveext.AgentClusterInstall
	if err := yaml.Unmarshal(aciData, &aci); err != nil {
		fmt.Errorf("Error unmarshalling AgentClusterInstall CR: %w", err)
	}
	return aci
}

// createClusterParams and associated functions were copied from
// https://github.com/openshift/assisted-service/blob/c5eacda676475f5a6de123678c1af353a2368bd3/internal/controller/controllers/clusterdeployments_controller.go#L1088
// TODO: Refactor clusterdeployments_controller to have a CreateClusterParams function that can be used in controller and here.
//       After the refactoring most of the code below goes away, especially the helper functions that are being carried over here.
func CreateClusterParams() *models.ClusterCreateParams {
	cd := getClusterDeployment()
	aci := getAgentClusterInstall()
	clusterInstall := &aci
	// TODO: Have single source for image version and cpu arch
	releaseImageVersion := "4.10.0-rc.1"
	releaseImageCPUArch := "x86_64"
	pullSecret := GetPullSecret()

	clusterParams := &models.ClusterCreateParams{
		BaseDNSDomain:         cd.Spec.BaseDomain,
		Name:                  swag.String(cd.Spec.ClusterName),
		OpenshiftVersion:      &releaseImageVersion,
		OlmOperators:          nil, // TODO: handle operators
		PullSecret:            swag.String(pullSecret),
		VipDhcpAllocation:     swag.Bool(false),
		IngressVip:            clusterInstall.Spec.IngressVIP,
		SSHPublicKey:          clusterInstall.Spec.SSHPublicKey,
		CPUArchitecture:       releaseImageCPUArch,
		UserManagedNetworking: swag.Bool(isUserManagedNetwork(clusterInstall)),
	}

	if len(clusterInstall.Spec.Networking.ClusterNetwork) > 0 {
		for _, net := range clusterInstall.Spec.Networking.ClusterNetwork {
			clusterParams.ClusterNetworks = append(clusterParams.ClusterNetworks, &models.ClusterNetwork{
				Cidr:       models.Subnet(net.CIDR),
				HostPrefix: int64(net.HostPrefix)})
		}
	}

	if len(clusterInstall.Spec.Networking.ServiceNetwork) > 0 {
		for _, cidr := range clusterInstall.Spec.Networking.ServiceNetwork {
			clusterParams.ServiceNetworks = append(clusterParams.ServiceNetworks, &models.ServiceNetwork{
				Cidr: models.Subnet(cidr),
			})
		}
	}

	if clusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents == 1 &&
		clusterInstall.Spec.ProvisionRequirements.WorkerAgents == 0 {
		clusterParams.HighAvailabilityMode = swag.String("None")
	}

	if hyperthreadingInSpec(clusterInstall) {
		clusterParams.Hyperthreading = getHyperthreading(clusterInstall)
	}

	if isDiskEncryptionEnabled(clusterInstall) {
		clusterParams.DiskEncryption = &models.DiskEncryption{
			EnableOn:    clusterInstall.Spec.DiskEncryption.EnableOn,
			Mode:        clusterInstall.Spec.DiskEncryption.Mode,
			TangServers: clusterInstall.Spec.DiskEncryption.TangServers,
		}
	}

	return clusterParams
}

func isUserManagedNetwork(clusterInstall *hiveext.AgentClusterInstall) bool {
	return clusterInstall.Spec.Networking.UserManagedNetworking ||
		clusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents == 1 && clusterInstall.Spec.ProvisionRequirements.WorkerAgents == 0
}

//see https://docs.openshift.com/container-platform/4.7/installing/installing_platform_agnostic/installing-platform-agnostic.html#installation-bare-metal-config-yaml_installing-platform-agnostic
func hyperthreadingInSpec(clusterInstall *hiveext.AgentClusterInstall) bool {
	//check if either master or worker pool hyperthreading settings are explicitly specified
	return clusterInstall.Spec.ControlPlane != nil ||
		funk.Contains(clusterInstall.Spec.Compute, func(pool hiveext.AgentMachinePool) bool {
			return pool.Name == hiveext.WorkerAgentMachinePool
		})
}

func getHyperthreading(clusterInstall *hiveext.AgentClusterInstall) *string {
	const (
		None    = 0
		Workers = 1
		Masters = 2
		All     = 3
	)
	var config uint = 0

	//if there is no configuration of hyperthreading in the Spec then
	//we are opting of the default behavior which is all enabled
	if !hyperthreadingInSpec(clusterInstall) {
		config = All
	}

	//check if the Spec enables hyperthreading for workers
	for _, machinePool := range clusterInstall.Spec.Compute {
		if machinePool.Name == hiveext.WorkerAgentMachinePool && machinePool.Hyperthreading == hiveext.HyperthreadingEnabled {
			config = config | Workers
		}
	}

	//check if the Spec enables hyperthreading for masters
	if clusterInstall.Spec.ControlPlane != nil {
		if clusterInstall.Spec.ControlPlane.Hyperthreading == hiveext.HyperthreadingEnabled {
			config = config | Masters
		}
	}

	//map between CRD Spec and cluster API
	switch config {
	case None:
		return swag.String(models.ClusterHyperthreadingNone)
	case Workers:
		return swag.String(models.ClusterHyperthreadingWorkers)
	case Masters:
		return swag.String(models.ClusterHyperthreadingMasters)
	default:
		return swag.String(models.ClusterHyperthreadingAll)
	}
}

func isDiskEncryptionEnabled(clusterInstall *hiveext.AgentClusterInstall) bool {
	if clusterInstall.Spec.DiskEncryption == nil {
		return false
	}
	switch swag.StringValue(clusterInstall.Spec.DiskEncryption.EnableOn) {
	case models.DiskEncryptionEnableOnAll, models.DiskEncryptionEnableOnMasters, models.DiskEncryptionEnableOnWorkers:
		return true
	case models.DiskEncryptionEnableOnNone:
		return false
	default:
		return false
	}
}
