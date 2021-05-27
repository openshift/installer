package management

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

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
	List(*ResourceGroupQuery) ([]models.ResourceGroup, error)
	// Get resource group by ID
	Get(id string) (*models.ResourceGroup, error)
	// Find resource groups having the specific name
	FindByName(*ResourceGroupQuery, string) ([]models.ResourceGroup, error)
	// Create a new resource group
	Create(models.ResourceGroup) (*models.ResourceGroup, error)
	// Delete an existing resource group
	Delete(id string) error
	// Update an existing resource group
	Update(id string, request *ResourceGroupUpdateRequest) (*models.ResourceGroup, error)
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

func (r *resourceGroup) List(query *ResourceGroupQuery) ([]models.ResourceGroup, error) {
	listRequest := rest.GetRequest("/v1/resource_groups")
	if query != nil {
		query.MakeRequest(listRequest)
	}
	req, err := listRequest.Build()
	if err != nil {
		return []models.ResourceGroup{}, err
	}

	var groups []models.ResourceGroup
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewRCPaginatedResources(models.ResourceGroup{}),
		func(resource interface{}) bool {
			if group, ok := resource.(models.ResourceGroup); ok {
				groups = append(groups, group)
				return true
			}
			return false
		},
	)

	if err != nil {
		return []models.ResourceGroup{}, err
	}

	return groups, nil
}

func (r *resourceGroup) FindByName(query *ResourceGroupQuery, name string) ([]models.ResourceGroup, error) {
	groups, err := r.List(query)
	if err != nil {
		return []models.ResourceGroup{}, err
	}

	filteredGroups := []models.ResourceGroup{}
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

func (r *resourceGroup) Create(group models.ResourceGroup) (*models.ResourceGroup, error) {
	newGroup := models.ResourceGroup{}
	_, err := r.client.Post("/v1/resource_groups", &group, &newGroup)
	if err != nil {
		return nil, err
	}
	return &newGroup, nil
}

func (r *resourceGroup) Delete(id string) error {
	_, err := r.client.Delete("/v1/resource_groups/" + id)
	return err
}

func (r *resourceGroup) Update(id string, request *ResourceGroupUpdateRequest) (*models.ResourceGroup, error) {
	updatedGroup := models.ResourceGroup{}
	_, err := r.client.Patch("/v1/resource_groups/"+id, request, &updatedGroup)
	if err != nil {
		return nil, err
	}
	return &updatedGroup, nil
}

func (r *resourceGroup) Get(id string) (*models.ResourceGroup, error) {
	group := models.ResourceGroup{}
	_, err := r.client.Get("/v1/resource_groups/"+id, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}
