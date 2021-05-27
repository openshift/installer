package iamuumv2

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
	Groups []models.AccessGroupV2 `json:"groups"`
}

func (g *Groups) Resources() []interface{} {
	r := make([]interface{}, len(g.Groups))
	for i := range g.Groups {
		r[i] = g.Groups[i]
	}
	return r
}

type AccessGroupRepository interface {
	List(accountID string, queryParams ...string) ([]models.AccessGroupV2, error)
	Create(group models.AccessGroupV2, accountID string) (*models.AccessGroupV2, error)
	FindByName(name string, accountID string) ([]models.AccessGroupV2, error)
	Delete(accessGroupID string, recursive bool) error
	Update(accessGroupID string, group AccessGroupUpdateRequest, revision string) (models.AccessGroupV2, error)
	Get(accessGroupID string) (group *models.AccessGroupV2, revision string, err error)
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

func (r *accessGroupRepository) List(accountID string, queryParams ...string) ([]models.AccessGroupV2, error) {
	var groups []models.AccessGroupV2
	var err error
	if len(queryParams) != 0 {
		_, err = r.client.GetPaginated(fmt.Sprintf("/v2/groups?account_id=%s&iam_id=%s", url.QueryEscape(accountID), queryParams[0]), NewPaginatedResourcesHandler(&Groups{}), func(v interface{}) bool {
			groups = append(groups, v.(models.AccessGroupV2))
			return true
		})
	} else {
		_, err = r.client.GetPaginated(fmt.Sprintf("/v2/groups?account_id=%s", url.QueryEscape(accountID)), NewPaginatedResourcesHandler(&Groups{}), func(v interface{}) bool {
			groups = append(groups, v.(models.AccessGroupV2))
			return true
		})
	}
	if err != nil {
		return []models.AccessGroupV2{}, err
	}
	return groups, err
}

func (r *accessGroupRepository) Create(accessGroup models.AccessGroupV2, accountID string) (*models.AccessGroupV2, error) {
	req := rest.PostRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v2/groups")).Query("account_id", accountID).Body(accessGroup)

	newAccessGroup := models.AccessGroupV2{}
	_, err := r.client.SendRequest(req, &newAccessGroup)
	if err != nil {
		return nil, err
	}
	return &newAccessGroup, nil
}

func (r *accessGroupRepository) FindByName(name string, accountID string) ([]models.AccessGroupV2, error) {
	var groups []models.AccessGroupV2
	_, err := r.client.GetPaginated(fmt.Sprintf("/v2/groups?account=%s", url.QueryEscape(accountID)), NewPaginatedResourcesHandler(&Groups{}), func(v interface{}) bool {
		if v.(models.AccessGroupV2).AccessGroup.Name == name {
			groups = append(groups, v.(models.AccessGroupV2))
		}
		return true
	})
	if err != nil {
		return []models.AccessGroupV2{}, err
	}
	return groups, err
}

func (r *accessGroupRepository) Delete(accessGroupID string, recursive bool) error {
	req := rest.DeleteRequest((helpers.GetFullURL(*r.client.Config.Endpoint, "/v2/groups/"+accessGroupID)))

	if recursive {
		req = req.Query("force", "true")
	}
	_, err := r.client.SendRequest(req, nil)
	return err
}

func (r *accessGroupRepository) Update(accessGroupID string, group AccessGroupUpdateRequest, revision string) (models.AccessGroupV2, error) {
	req := rest.PatchRequest((helpers.GetFullURL(*r.client.Config.Endpoint, "/v2/groups/"+accessGroupID))).Body(group).Add("If-Match", revision)
	resp := models.AccessGroupV2{}
	_, err := r.client.SendRequest(req, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (r *accessGroupRepository) Get(accessGroupID string) (*models.AccessGroupV2, string, error) {
	group := models.AccessGroupV2{}
	response, err := r.client.Get("/v2/groups/"+url.PathEscape(accessGroupID), &group)
	if err != nil {
		return &group, "", err
	}
	return &group, response.Header.Get("Etag"), nil
}
