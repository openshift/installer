package mccpv2

import (
	"encoding/json"
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//OrgQuota ...
type OrgQuota struct {
	GUID                    string
	Name                    string
	NonBasicServicesAllowed bool
	ServicesLimit           int
	RoutesLimit             int
	MemoryLimitInMB         int64
	InstanceMemoryLimitInMB int64
	TrialDBAllowed          bool
	AppInstanceLimit        int
	PrivateDomainsLimit     int
	AppTasksLimit           int
	ServiceKeysLimit        int
	RoutePortsLimit         int
}

//OrgQuotaFields ...
type OrgQuotaFields struct {
	Metadata OrgQuotaMetadata
	Entity   OrgQuotaEntity
}

//OrgQuotaMetadata ...
type OrgQuotaMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ErrCodeOrgQuotaDoesnotExist ...
const ErrCodeOrgQuotaDoesnotExist = "OrgQuotaDoesnotExist"

//OrgQuotaResource ...
type OrgQuotaResource struct {
	Resource
	Entity OrgQuotaEntity
}

//OrgQuotaEntity ...
type OrgQuotaEntity struct {
	Name                    string      `json:"name"`
	NonBasicServicesAllowed bool        `json:"non_basic_services_allowed"`
	ServicesLimit           int         `json:"total_services"`
	RoutesLimit             int         `json:"total_routes"`
	MemoryLimitInMB         int64       `json:"memory_limit"`
	InstanceMemoryLimitInMB int64       `json:"instance_memory_limit"`
	TrialDBAllowed          bool        `json:"trial_db_allowed"`
	AppInstanceLimit        json.Number `json:"app_instance_limit"`
	PrivateDomainsLimit     json.Number `json:"total_private_domains"`
	AppTasksLimit           json.Number `json:"app_tasks_limit"`
	ServiceKeysLimit        json.Number `json:"total_service_keys"`
	RoutePortsLimit         int         `json:"total_reserved_route_ports"`
}

//ToFields ...
func (resource OrgQuotaResource) ToFields() OrgQuota {
	entity := resource.Entity

	return OrgQuota{
		GUID:                    resource.Metadata.GUID,
		Name:                    entity.Name,
		NonBasicServicesAllowed: entity.NonBasicServicesAllowed,
		ServicesLimit:           entity.ServicesLimit,
		RoutesLimit:             entity.RoutesLimit,
		MemoryLimitInMB:         entity.MemoryLimitInMB,
		InstanceMemoryLimitInMB: entity.InstanceMemoryLimitInMB,
		TrialDBAllowed:          entity.TrialDBAllowed,
		AppInstanceLimit:        NumberToInt(entity.AppInstanceLimit, -1),
		PrivateDomainsLimit:     NumberToInt(entity.PrivateDomainsLimit, -1),
		AppTasksLimit:           NumberToInt(entity.AppTasksLimit, -1),
		ServiceKeysLimit:        NumberToInt(entity.ServiceKeysLimit, -1),
		RoutePortsLimit:         entity.RoutePortsLimit,
	}
}

//OrgQuotas ...
type OrgQuotas interface {
	FindByName(name string) (*OrgQuota, error)
	Get(orgQuotaGUID string) (*OrgQuotaFields, error)
	List() ([]OrgQuota, error)
}

type orgQuota struct {
	client *client.Client
}

func newOrgQuotasAPI(c *client.Client) OrgQuotas {
	return &orgQuota{
		client: c,
	}
}

func (r *orgQuota) FindByName(name string) (*OrgQuota, error) {
	rawURL := fmt.Sprintf("/v2/quota_definitions")
	req := rest.GetRequest(rawURL).Query("q", "name:"+name)

	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	orgQuotas, err := r.listOrgQuotaWithPath(path)
	if err != nil {
		return nil, err
	}

	if len(orgQuotas) == 0 {
		return nil, bmxerror.New(ErrCodeAppDoesnotExist,
			fmt.Sprintf("Given quota definition: %q doesn't exist", name))

	}
	return &orgQuotas[0], nil
}

func (r *orgQuota) List() ([]OrgQuota, error) {
	rawURL := fmt.Sprintf("/v2/quota_definitions")
	req := rest.GetRequest(rawURL)

	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	orgQuotas, err := r.listOrgQuotaWithPath(path)
	if err != nil {
		return nil, err
	}

	return orgQuotas, nil
}

func (r *orgQuota) listOrgQuotaWithPath(path string) ([]OrgQuota, error) {
	var orgQuota []OrgQuota
	_, err := r.client.GetPaginated(path, NewCCPaginatedResources(OrgQuotaResource{}), func(resource interface{}) bool {
		if orgQuotaResource, ok := resource.(OrgQuotaResource); ok {
			orgQuota = append(orgQuota, orgQuotaResource.ToFields())
			return true
		}
		return false
	})
	return orgQuota, err
}

func (r *orgQuota) Get(quotaGUID string) (*OrgQuotaFields, error) {
	rawURL := fmt.Sprintf("/v2/quota_definitions/%s", quotaGUID)
	orgQuotaFields := OrgQuotaFields{}
	_, err := r.client.Get(rawURL, &orgQuotaFields)
	if err != nil {
		return nil, err
	}

	return &orgQuotaFields, err
}
