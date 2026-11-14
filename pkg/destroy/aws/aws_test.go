package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/stretchr/testify/assert"

	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func TestEndpointsForRegion(t *testing.T) {
	customEndpoints := []awstypes.ServiceEndpoint{
		{
			Name: "tagging",
			URL:  "https://tagging.eu-west-1.example.com",
		},
		{
			Name: "ec2",
			URL:  "https://ec2.eu-west-1.example.com",
		},
		{
			Name: "sts",
			URL:  "https://sts.eu-west-1.example.com",
		},
	}

	uninstaller := &ClusterUninstaller{
		Region:    "eu-west-1",
		endpoints: customEndpoints,
	}

	t.Run("returns all endpoints when region matches install region", func(t *testing.T) {
		result := uninstaller.endpointsForRegion("eu-west-1")
		assert.Equal(t, customEndpoints, result)
	})

	t.Run("filters tagging endpoint when region differs", func(t *testing.T) {
		result := uninstaller.endpointsForRegion("us-east-1")
		assert.Equal(t, []awstypes.ServiceEndpoint{
			{Name: "ec2", URL: "https://ec2.eu-west-1.example.com"},
			{Name: "sts", URL: "https://sts.eu-west-1.example.com"},
		}, result)
	})

	t.Run("preserves STS endpoint for cross-region assumed role", func(t *testing.T) {
		result := uninstaller.endpointsForRegion(awstypes.UsEast1RegionID)
		var hasSTS, hasTagging bool
		for _, ep := range result {
			if ep.Name == "sts" {
				hasSTS = true
			}
			if ep.Name == "tagging" {
				hasTagging = true
			}
		}
		assert.True(t, hasSTS, "STS endpoint should be preserved for credential resolution")
		assert.False(t, hasTagging, "tagging endpoint should be filtered for cross-region client")
	})

	t.Run("returns nil when only tagging endpoint exists and region differs", func(t *testing.T) {
		taggingOnly := &ClusterUninstaller{
			Region: "eu-west-1",
			endpoints: []awstypes.ServiceEndpoint{
				{Name: "tagging", URL: "https://tagging.eu-west-1.example.com"},
			},
		}
		result := taggingOnly.endpointsForRegion("us-east-1")
		assert.Nil(t, result)
	})

	t.Run("returns nil for cn-northwest-1 when install region is eu-west-1", func(t *testing.T) {
		taggingOnly := &ClusterUninstaller{
			Region: "eu-west-1",
			endpoints: []awstypes.ServiceEndpoint{
				{Name: "tagging", URL: "https://tagging.eu-west-1.example.com"},
			},
		}
		result := taggingOnly.endpointsForRegion(awstypes.CnNorthwest1RegionID)
		assert.Nil(t, result)
	})

	t.Run("returns nil for us-gov-west-1 when install region is eu-west-1", func(t *testing.T) {
		taggingOnly := &ClusterUninstaller{
			Region: "eu-west-1",
			endpoints: []awstypes.ServiceEndpoint{
				{Name: "tagging", URL: "https://tagging.eu-west-1.example.com"},
			},
		}
		result := taggingOnly.endpointsForRegion(awstypes.UsGovWest1RegionID)
		assert.Nil(t, result)
	})

	t.Run("returns nil when no custom endpoints configured", func(t *testing.T) {
		emptyUninstaller := &ClusterUninstaller{
			Region: "eu-west-1",
		}
		result := emptyUninstaller.endpointsForRegion("us-east-1")
		assert.Nil(t, result)
	})

	t.Run("returns nil when install region has no custom endpoints and region matches", func(t *testing.T) {
		emptyUninstaller := &ClusterUninstaller{
			Region: "eu-west-1",
		}
		result := emptyUninstaller.endpointsForRegion("eu-west-1")
		assert.Nil(t, result)
	})

	t.Run("filters SDK service ID name for tagging when region differs", func(t *testing.T) {
		sdkNameUninstaller := &ClusterUninstaller{
			Region: "eu-west-1",
			endpoints: []awstypes.ServiceEndpoint{
				{Name: resourcegroupstaggingapi.ServiceID, URL: "https://tagging.eu-west-1.example.com"},
				{Name: "sts", URL: "https://sts.eu-west-1.example.com"},
			},
		}
		result := sdkNameUninstaller.endpointsForRegion("us-east-1")
		assert.Equal(t, []awstypes.ServiceEndpoint{
			{Name: "sts", URL: "https://sts.eu-west-1.example.com"},
		}, result)
	})
}
