package manifests

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"sigs.k8s.io/yaml"
)

func getInfraEnv() aiv1beta1.InfraEnv {
	infraEnvData, err := os.ReadFile("/data/manifests/infraenv.yaml")
	if err != nil {
		fmt.Errorf("Error reading pull secret: %w", err)
	}
	var infraEnv aiv1beta1.InfraEnv
	if err := yaml.Unmarshal(infraEnvData, &infraEnv); err != nil {
		fmt.Errorf("Error unmarshalling pull secret: %w", err)
	}
	return infraEnv
}

// createInfraEnvParams body was copied from
// https://github.com/openshift/assisted-service/blob/5d4d836747862f43fa2ec882e5871648bd12c780/internal/controller/controllers/infraenv_controller.go#L339
// TODO: Refactor infraenv_controller to have a CreateInfraEnvParams function that can be used in controller and here.
func CreateInfraEnvParams() *models.InfraEnvCreateParams {
	infraEnv := getInfraEnv()
	// TODO: Have single source for image version and cpu arch
	releaseImageVersion := "4.10.0-rc.1"
	releaseImageCPUArch := "x86_64"
	pullSecret := GetPullSecret()

	createParams := &models.InfraEnvCreateParams{
		Name:                   &infraEnv.Name,
		ImageType:              "full-iso",
		IgnitionConfigOverride: infraEnv.Spec.IgnitionConfigOverride,
		PullSecret:             &pullSecret,
		SSHAuthorizedKey:       &infraEnv.Spec.SSHAuthorizedKey,
		CPUArchitecture:        releaseImageCPUArch,
	}
	if infraEnv.Spec.Proxy != nil {
		proxy := &models.Proxy{
			HTTPProxy:  &infraEnv.Spec.Proxy.HTTPProxy,
			HTTPSProxy: &infraEnv.Spec.Proxy.HTTPSProxy,
			NoProxy:    &infraEnv.Spec.Proxy.NoProxy,
		}
		createParams.Proxy = proxy
	}

	if len(infraEnv.Spec.AdditionalNTPSources) > 0 {
		createParams.AdditionalNtpSources = swag.String(strings.Join(infraEnv.Spec.AdditionalNTPSources[:], ","))
	}

	// cluster-id is set in shell script
	var tempClusterID strfmt.UUID
	tempClusterID = "replace-cluster-id"
	createParams.ClusterID = &tempClusterID
	createParams.OpenshiftVersion = &releaseImageVersion

	return createParams
}
