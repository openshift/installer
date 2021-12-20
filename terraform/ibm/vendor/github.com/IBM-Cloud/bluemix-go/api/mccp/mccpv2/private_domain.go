package mccpv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ErrCodePrivateDomainDoesnotExist ...
var ErrCodePrivateDomainDoesnotExist = "PrivateDomainDoesnotExist"

//PrivateDomainRequest ...
type PrivateDomainRequest struct {
	Name    string `json:"name,omitempty"`
	OrgGUID string `json:"owning_organization_guid,omitempty"`
}

//PrivateDomaineMetadata ...
type PrivateDomainMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//PrivateDomainEntity ...
type PrivateDomainEntity struct {
	Name                   string `json:"name"`
	OwningOrganizationGUID string `json:"owning_organization_guid"`
	OwningOrganizationURL  string `json:"owning_organization_url"`
	SharedOrganizationURL  string `json:"shared_organizations_url"`
}

//PrivateDomainResource ...
type PrivateDomainResource struct {
	Resource
	Entity PrivateDomainEntity
}

//PrivateDomainFields ...
type PrivateDomainFields struct {
	Metadata PrivateDomainMetadata
	Entity   PrivateDomainEntity
}

//ToFields ..
func (resource PrivateDomainResource) ToFields() PrivateDomain {
	entity := resource.Entity

	return PrivateDomain{
		GUID:                   resource.Metadata.GUID,
		Name:                   entity.Name,
		OwningOrganizationGUID: entity.OwningOrganizationGUID,
		OwningOrganizationURL:  entity.OwningOrganizationURL,
		SharedOrganizationURL:  entity.OwningOrganizationURL,
	}
}

//PrivateDomain model
type PrivateDomain struct {
	GUID                   string
	Name                   string
	OwningOrganizationGUID string
	OwningOrganizationURL  string
	SharedOrganizationURL  string
}

//PrivateDomains ...
type PrivateDomains interface {
	FindByNameInOrg(orgGUID, domainName string) (*PrivateDomain, error)
	FindByName(domainName string) (*PrivateDomain, error)
	Create(req PrivateDomainRequest, opts ...bool) (*PrivateDomainFields, error)
	Get(privateDomainGUID string) (*PrivateDomainFields, error)
	Delete(privateDomainGUID string, opts ...bool) error
}

type privateDomain struct {
	client *client.Client
}

func newPrivateDomainAPI(c *client.Client) PrivateDomains {
	return &privateDomain{
		client: c,
	}
}

func (d *privateDomain) FindByNameInOrg(orgGUID, domainName string) (*PrivateDomain, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/private_domains", orgGUID)
	req := rest.GetRequest(rawURL).Query("q", "name:"+domainName)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	domain, err := listPrivateDomainWithPath(d.client, path)
	if err != nil {
		return nil, err
	}
	if len(domain) == 0 {
		return nil, bmxerror.New(ErrCodePrivateDomainDoesnotExist, fmt.Sprintf("Private Domain: %q doesn't exist", domainName))
	}
	return &domain[0], nil
}

func (d *privateDomain) FindByName(domainName string) (*PrivateDomain, error) {
	rawURL := fmt.Sprintf("/v2/private_domains")
	req := rest.GetRequest(rawURL).Query("q", "name:"+domainName)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	domain, err := listPrivateDomainWithPath(d.client, path)
	if err != nil {
		return nil, err
	}
	if len(domain) == 0 {
		return nil, bmxerror.New(ErrCodePrivateDomainDoesnotExist, fmt.Sprintf("Private Domain: %q doesn't exist", domainName))
	}
	return &domain[0], nil
}

func listPrivateDomainWithPath(c *client.Client, path string) ([]PrivateDomain, error) {
	var privateDomain []PrivateDomain
	_, err := c.GetPaginated(path, NewCCPaginatedResources(PrivateDomainResource{}), func(resource interface{}) bool {
		if privateDomainResource, ok := resource.(PrivateDomainResource); ok {
			privateDomain = append(privateDomain, privateDomainResource.ToFields())
			return true
		}
		return false
	})
	return privateDomain, err
}

/* opts is list of boolean parametes
opts[0] - async - Will run the create request in a background job. Recommended: 'true'. Default to 'true'.
*/
func (d *privateDomain) Create(req PrivateDomainRequest, opts ...bool) (*PrivateDomainFields, error) {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/private_domains?async=%t", async)
	privateDomainFields := PrivateDomainFields{}
	_, err := d.client.Post(rawURL, req, &privateDomainFields)
	if err != nil {
		return nil, err
	}
	return &privateDomainFields, nil
}

func (d *privateDomain) Get(privateDomainGUID string) (*PrivateDomainFields, error) {
	rawURL := fmt.Sprintf("/v2/private_domains/%s", privateDomainGUID)
	privateDomainFields := PrivateDomainFields{}
	_, err := d.client.Get(rawURL, &privateDomainFields, nil)
	if err != nil {
		return nil, err
	}
	return &privateDomainFields, nil
}

// opts is list of boolean parametes
// opts[0] - async - Will run the delete request in a background job. Recommended: 'true'. Default to 'true'.

func (d *privateDomain) Delete(privateDomainGUID string, opts ...bool) error {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/private_domains/%s?async=%t", privateDomainGUID, async)
	_, err := d.client.Delete(rawURL)
	return err
}
