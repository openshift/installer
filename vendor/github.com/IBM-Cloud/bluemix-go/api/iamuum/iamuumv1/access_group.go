package iamuumv1

import (
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type Groups struct {
	PaginationFields
	Groups []models.AccessGroup `json:"groups"`
}

func (g *Groups) Resources() []interface{} {
	r := make([]interface{}, len(g.Groups))
	for i := range g.Groups {
		r[i] = g.Groups[i]
	}
	return r
}

type AccessGroupRepository interface {
	List(accountID string) ([]models.AccessGroup, error)
	Create(group models.AccessGroup, accountID string) (*models.AccessGroup, error)
	FindByName(name string, accountID string) ([]models.AccessGroup, error)
	Delete(accessGroupID string, recursive bool) error
	Update(accessGroupID string, group AccessGroupUpdateRequest, revision string) (models.AccessGroup, error)
	Get(accessGroupID string) (group *models.AccessGroup, revision string, err error)
}

type accessGroupRepository struct {
	client *client.Client
}

type AccessGroupUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func NewAccessGroupRepository(c *client.Client) AccessGroupRepository {
	return &accessGroupRepository{
		client: c,
	}
}

func (r *accessGroupRepository) List(accountID string) ([]models.AccessGroup, error) {
	var groups []models.AccessGroup
	_, err := r.client.GetPaginated(fmt.Sprintf("/v1/groups?account=%s", url.QueryEscape(accountID)), NewPaginatedResourcesHandler(&Groups{}), func(v interface{}) bool {
		groups = append(groups, v.(models.AccessGroup))
		return true
	})
	if err != nil {
		return []models.AccessGroup{}, err
	}
	return groups, err
}

func (r *accessGroupRepository) Create(accessGroup models.AccessGroup, accountID string) (*models.AccessGroup, error) {
	req := rest.PostRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/groups")).Query("account", accountID).Body(accessGroup)

	newAccessGroup := models.AccessGroup{}
	_, err := r.client.SendRequest(req, &newAccessGroup)
	if err != nil {
		return nil, err
	}
	return &newAccessGroup, nil
}

func (r *accessGroupRepository) FindByName(name string, accountID string) ([]models.AccessGroup, error) {
	var groups []models.AccessGroup
	_, err := r.client.GetPaginated(fmt.Sprintf("/v1/groups?account=%s", url.QueryEscape(accountID)), NewPaginatedResourcesHandler(&Groups{}), func(v interface{}) bool {
		if v.(models.AccessGroup).Name == name {
			groups = append(groups, v.(models.AccessGroup))
		}
		return true
	})
	if err != nil {
		return []models.AccessGroup{}, err
	}
	return groups, err
}

func (r *accessGroupRepository) Delete(accessGroupID string, recursive bool) error {
	req := rest.DeleteRequest((helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/groups/"+accessGroupID)))

	if recursive {
		req = req.Query("force", "true")
	}
	_, err := r.client.SendRequest(req, nil)
	return err
}

func (r *accessGroupRepository) Update(accessGroupID string, group AccessGroupUpdateRequest, revision string) (models.AccessGroup, error) {
	req := rest.PatchRequest((helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/groups/"+accessGroupID))).Body(group).Add("If-Match", revision)
	resp := models.AccessGroup{}
	_, err := r.client.SendRequest(req, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (r *accessGroupRepository) Get(accessGroupID string) (*models.AccessGroup, string, error) {
	group := models.AccessGroup{}
	response, err := r.client.Get("/v1/groups/"+url.PathEscape(accessGroupID), &group)
	if err != nil {
		return &group, "", err
	}
	return &group, response.Header.Get("Etag"), nil
}
