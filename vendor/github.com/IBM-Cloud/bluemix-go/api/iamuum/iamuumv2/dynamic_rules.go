package iamuumv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type CreateRuleRequest struct {
	Name       string      `json:"name"`
	Expiration int         `json:"expiration"`
	RealmName  string      `json:"realm_name,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	Claim    string `json:"claim"`
	Operator string `json:"operator"`
	Value    string `json:"value,omitempty"`
}

type CreateRuleResponse struct {
	CreateRuleRequest
	RuleID           string `json:"id"`
	AccessGroupID    string `json:"access_group_id"`
	AccountID        string `json:"account_id"`
	CreatedAt        string `json:"created_at"`
	CreatedByID      string `json:"created_by_id"`
	LastModifiedAt   string `json:"last_modified_at"`
	LastModifiedByID string `json:"last_modified_by_id"`
}

type GetResponse struct {
	Rules []CreateRuleResponse `json:"rules"`
}

type DynamicRuleRepository interface {
	Create(groupID string, request CreateRuleRequest) (CreateRuleResponse, error)
	List(groupID string) ([]CreateRuleResponse, error)
	Get(groupID, ruleID string) (CreateRuleResponse, string, error)
	Replace(groupID, ruleID string, request CreateRuleRequest, etag string) (CreateRuleResponse, error)
	Delete(groupID, ruleID string) error
}

type dynamicRuleRepository struct {
	client *client.Client
}

func NewDynamicRuleRepository(c *client.Client) DynamicRuleRepository {
	return &dynamicRuleRepository{
		client: c,
	}
}

func (r *dynamicRuleRepository) Create(groupID string, request CreateRuleRequest) (CreateRuleResponse, error) {
	res := CreateRuleResponse{}
	_, err := r.client.Post(fmt.Sprintf("/v2/groups/%s/rules", groupID), &request, &res)
	if err != nil {
		return CreateRuleResponse{}, err
	}
	return res, nil
}

func (r *dynamicRuleRepository) List(groupID string) ([]CreateRuleResponse, error) {
	res := GetResponse{}
	_, err := r.client.Get(fmt.Sprintf("/v2/groups/%s/rules", groupID), &res)
	if err != nil {
		return []CreateRuleResponse{}, err
	}
	return res.Rules, nil
}

func (r *dynamicRuleRepository) Get(groupID, ruleID string) (CreateRuleResponse, string, error) {
	res := CreateRuleResponse{}
	response, err := r.client.Get(fmt.Sprintf("/v2/groups/%s/rules/%s", groupID, ruleID), &res)
	if err != nil {
		return CreateRuleResponse{}, "", err
	}
	return res, response.Header.Get("Etag"), nil
}

func (r *dynamicRuleRepository) Replace(groupID, ruleID string, request CreateRuleRequest, etag string) (CreateRuleResponse, error) {
	res := CreateRuleResponse{}
	header := make(map[string]string)

	header["IF-Match"] = etag
	_, err := r.client.Put(fmt.Sprintf("/v2/groups/%s/rules/%s", groupID, ruleID), &request, &res, header)
	if err != nil {
		return CreateRuleResponse{}, err
	}
	return res, nil
}

//Delete Function
func (r *dynamicRuleRepository) Delete(groupID, ruleID string) error {
	_, err := r.client.Delete(fmt.Sprintf("/v2/groups/%s/rules/%s", groupID, ruleID))
	return err
}
