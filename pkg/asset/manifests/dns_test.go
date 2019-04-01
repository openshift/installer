package manifests

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/dns"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	azure "github.com/openshift/installer/pkg/types/azure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
)

func setupParents() asset.Parents {
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
	return parents
}

func TestGenerateAzure(t *testing.T) {
	expectedPublicZoneID := "/subscriptions/<subid>/resourceGroups/<rg>/providers/Microsoft.Network/dnszones/my-dns-zone.com"
	dnsAsset := DNS{}
	dnsAsset.DNSConfig = &dns.MockConfigProvider{
		BaseDomain: "",
		PublicZone: expectedPublicZoneID,
	}

	err := dnsAsset.Generate(setupParents())
	assert.NoError(t, err)
	dnsyaml := dnsAsset.Files()[0].Data

	dnsConfig := configv1.DNS{}
	yaml.Unmarshal(dnsyaml, &dnsConfig)
	assert.Equal(t, expectedPublicZoneID, dnsConfig.Spec.PublicZone.ID)
}
