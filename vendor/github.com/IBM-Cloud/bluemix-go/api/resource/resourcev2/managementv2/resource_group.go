package managementv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type ResourceOrigin string

func (o ResourceOrigin) String() string {
	return string(o)
}

//ResourceGroupUpdateRequest ...
type ResourceGroupUpdateRequest struct {
	Name    string `json:"name,omitempty"`
	QuotaID string `json:"quota_id,omitempty"`
	Default *bool  `json:"default,omitempty"`
}

// ResourceGroupQuery is the query filters to get resource groups
type ResourceGroupQuery struct {
	AccountID      string
	Default        bool
	ResourceID     string
	ResourceOrigin models.ResourceOrigin
}

//ErrCodeResourceGroupDoesnotExist ...
const ErrCodeResourceGroupDoesnotExist = "ResourceGroupDoesnotExist"

//ResourceGroupQuery ...
type ResourceGroupRepository interface {
	// List all available resource groups
	List(*ResourceGroupQuery) ([]models.ResourceGroupv2, error)
	// Get resource group by ID
	Get(id string) (*models.ResourceGroupv2, error)
	// Find resource groups having the specific name
	FindByName(*ResourceGroupQuery, string) ([]models.ResourceGroupv2, error)
	// Create a new resource group
	Create(models.ResourceGroupv2) (*models.ResourceGroupv2, error)
	// Delete an existing resource group
	Delete(id string) error
	// Update an existing resource group
	Update(id string, request *ResourceGroupUpdateRequest) (*models.ResourceGroupv2, error)
}

type resourceGroup struct {
	client *client.Client
}

func newResourceGroupAPI(c *client.Client) ResourceGroupRepository {
	return &resourceGroup{
		client: c,
	}
}

// populate query part of HTTP requests
func (q ResourceGroupQuery) MakeRequest(r *rest.Request) *rest.Request {
	if q.AccountID != "" {
		r.Query("account_id", q.AccountID)
	}
	if q.Default {
		r.Query("default", "true")
	}
	if q.ResourceID != "" {
		r.Query("resource_id", q.ResourceID)
	}
	if q.ResourceOrigin != "" {
		r.Query("resource_origin", q.ResourceOrigin.String())
	}
	return r
}

func (r *resourceGroup) List(query *ResourceGroupQuery) ([]models.ResourceGroupv2, error) {
	listRequest := rest.GetRequest("/v2/resource_groups")
	if query != nil {
		query.MakeRequest(listRequest)
	}
	req, err := listRequest.Build()
	if err != nil {
		return []models.ResourceGroupv2{}, err
	}

	var groups []models.ResourceGroupv2
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewRCPaginatedResources(models.ResourceGroupv2{}),
		func(resource interface{}) bool {
			if group, ok := resource.(models.ResourceGroupv2); ok {
				groups = append(groups, group)
				return true
			}
			return false
		},
	)

	if err != nil {
		return []models.ResourceGroupv2{}, err
	}

	return groups, nil
}

func (r *resourceGroup) FindByName(query *ResourceGroupQuery, name string) ([]models.ResourceGroupv2, error) {
	groups, err := r.List(query)
	if err != nil {
		return []models.ResourceGroupv2{}, err
	}

	filteredGroups := []models.ResourceGroupv2{}
	for _, group := range groups {
		if group.Name == name {
			filteredGroups = append(filteredGroups, group)
		}
	}

	if len(filteredGroups) == 0 {
		return filteredGroups, bmxerror.New(ErrCodeResourceGroupDoesnotExist,
			fmt.Sprintf("Given resource Group : %q doesn't exist", name))

	}
	return filteredGroups, nil
}

func (r *resourceGroup) Create(group models.ResourceGroupv2) (*models.ResourceGroupv2, error) {
	newGroup := models.ResourceGroupv2{}
	_, err := r.client.Post("/v2/resource_groups", &group, &newGroup)
	if err != nil {
		return nil, err
	}
	return &newGroup, nil
}

func (r *resourceGroup) Delete(id string) error {
	_, err := r.client.Delete("/v2/resource_groups/" + id)
	return err
}

func (r *resourceGroup) Update(id string, request *ResourceGroupUpdateRequest) (*models.ResourceGroupv2, error) {
	updatedGroup := models.ResourceGroupv2{}
	_, err := r.client.Patch("/v2/resource_groups/"+id, request, &updatedGroup)
	if err != nil {
		return nil, err
	}
	return &updatedGroup, nil
}

func (r *resourceGroup) Get(id string) (*models.ResourceGroupv2, error) {
	group := models.ResourceGroupv2{}
	_, err := r.client.Get("/v2/resource_groups/"+id, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}
