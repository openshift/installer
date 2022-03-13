package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

// ZoneCreateOpts represents the attributes used when creating a new DNS zone.
type ZoneCreateOpts struct {
	zones.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToZoneCreateMap casts a CreateOpts struct to a map.
// It overrides zones.ToZoneCreateMap to add the ValueSpecs field.
func (opts ZoneCreateOpts) ToZoneCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]interface{}); ok {
		if opts.TTL > 0 {
			m["ttl"] = opts.TTL
		}

		return m, nil
	}

	return nil, fmt.Errorf("Expected map but got %T", b[""])
}

const headerAuthSudoTenantID string = "X-Auth-Sudo-Tenant-ID"
const headerAuthAllProjects string = "X-Auth-All-Projects"

// dnsClientSetAuthHeaders sets auth headers for interacting with different projects.
func dnsClientSetAuthHeader(resourceData *schema.ResourceData, dnsClient *gophercloud.ServiceClient) error {
	// Extracting project ID from token to compare with provided one
	project, err := getProjectFromToken(dnsClient)
	if err != nil {
		return fmt.Errorf("Error extracting project ID from token: %s", err)
	}
	headers := make(map[string]string)
	// If all projects need to be listed to lookup a zone, set AuthAllProjects header
	if v, ok := resourceData.GetOk("all_projects"); ok {
		if allProjects, ok := v.(string); ok {
			headers[headerAuthAllProjects] = allProjects
		} else {
			return fmt.Errorf("Expected all_projects as string, but got %T", v)
		}
	}

	// If project_id is different from auth one, set AuthSudo header
	if v, ok := resourceData.GetOk("project_id"); ok {
		if projectID, ok := v.(string); ok {
			if project != nil && project.ID != projectID {
				headers[headerAuthSudoTenantID] = projectID
			}
		} else {
			return fmt.Errorf("Expected project_id as string, but got %T", v)
		}
	}

	if len(headers) != 0 {
		dnsClient.MoreHeaders = headers
		log.Printf("[DEBUG] request headers set: %#v", headers)
	}

	return nil
}

func dnsZoneV2RefreshFunc(dnsClient *gophercloud.ServiceClient, zoneID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		zone, err := zones.Get(dnsClient, zoneID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return zone, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_dns_zone_v2 %s current status: %s", zone.ID, zone.Status)
		return zone, zone.Status, nil
	}
}

func getProjectFromToken(dnsClient *gophercloud.ServiceClient) (*tokens.Project, error) {
	var (
		project *tokens.Project
		err     error
	)
	r := dnsClient.ProviderClient.GetAuthResult()
	switch result := r.(type) {
	case tokens.CreateResult:
		project, err = result.ExtractProject()
		if err != nil {
			return nil, err
		}
	case tokens.GetResult:
		project, err = result.ExtractProject()
		if err != nil {
			return nil, err
		}
	default:
		res := tokens.Get(dnsClient, dnsClient.ProviderClient.TokenID)
		project, err = res.ExtractProject()
		if err != nil {
			return nil, err
		}
	}
	return project, nil
}
