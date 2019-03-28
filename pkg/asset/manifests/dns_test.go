package manifests

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	azure "github.com/openshift/installer/pkg/types/azure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAzure(t *testing.T) {

	dnsAsset := DNS{}
	config := &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-name",
		},
		BaseDomain: "testdomain.com",
		Platform: types.Platform{
			Azure: &azure.Platform{
				Region: "eastus",
			},
		},
	}
	configAsset := &installconfig.InstallConfig{
		Config: config,
	}

	cid := &installconfig.ClusterID{
		InfraID: "infraId",
		UUID:    uuid.New().String(),
	}
	parents := asset.Parents{}
	parents.Add(configAsset, cid)

	dnsConfig := &configv1.DNS{}

	// config.Spec.PublicZone = &configv1.DNSZone{ID: strings.TrimPrefix(*zone.Id, "/hostedzone/")}
	// config.Spec.PrivateZone = &configv1.DNSZone{Tags: map[string]string{
	// 	  fmt.Sprintf("kubernetes.io/cluster/%s", clusterID.InfraID): "owned",
	// 	  "Name": fmt.Sprintf("%s-int", clusterID.InfraID),
	// }}
	err := dnsAsset.Generate(parents)
	assert.NoError(t, err)

	dnsyaml := dnsAsset.Files()[0].Data
	yaml.Unmarshal(dnsyaml, dnsConfig)

	assert.Equal(t, "", dnsConfig.Spec.PublicZone)
}
