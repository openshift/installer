package iampapv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type AccessPolicyRequest struct {
	Roles     []Roles     `json:"roles" binding:"required"`
	Resources []Resources `json:"resources" binding:"required"`
}

type AccessPolicyResponse struct {
	ID        string
	Roles     []Roles
	Resources []Resources
}

type AccessPolicyListResponse struct {
	Policies []AccessPolicyResponse
}

type Roles struct {
	ID          string `json:"id" binding:"required"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

type Resources struct {
	ServiceName     string `json:"serviceName,omitempty"`
	ServiceInstance string `json:"serviceInstance,omitempty"`
	Region          string `json:"region,omitempty"`
	ResourceType    string `json:"resourceType,omitempty"`
	Resource        string `json:"resource,omitempty"`
	SpaceId         string `json:"spaceId,omitempty"`
	AccountId       string `json:"accountId,omitempty"`
	OrganizationId  string `json:"organizationId,omitempty"`
}

type IAMPolicy interface {
	Create(scope, userId string, params AccessPolicyRequest) (AccessPolicyResponse, string, error)
	List(scope, userId string) (AccessPolicyListResponse, error)
	Delete(scope, userId, policyId string) error
	Get(scope, userId, policyId string) (AccessPolicyResponse, error)
	Update(scope, userId, policyId, etag string, params AccessPolicyRequest) (AccessPolicyResponse, string, error)
}

type iampolicy struct {
	client *client.Client
}

const IAM_ACCOUNT_ESCAPE = "a%2f"

func newIAMPolicyAPI(c *client.Client) IAMPolicy {
	return &iampolicy{
		client: c,
	}
}

//Create ...
func (r *iampolicy) Create(scope, userId string, params AccessPolicyRequest) (AccessPolicyResponse, string, error) {
	var accessPolicy AccessPolicyResponse
	rawURL := fmt.Sprintf("/acms/v1/scopes/%s/users/%s/policies", IAM_ACCOUNT_ESCAPE+scope, userId)
	resp, err := r.client.Post(rawURL, params, &accessPolicy)
	eTag := resp.Header.Get("etag")
	return accessPolicy, eTag, err
}

//List ...
func (r *iampolicy) List(scope, userId string) (AccessPolicyListResponse, error) {
	var accessPolicyListResponse AccessPolicyListResponse
	rawURL := fmt.Sprintf("/acms/v1/scopes/%s/users/%s/policies", IAM_ACCOUNT_ESCAPE+scope, userId)
	_, err := r.client.Get(rawURL, &accessPolicyListResponse)
	return accessPolicyListResponse, err
}

//Find ...
func (r *iampolicy) Get(scope, userId, policyId string) (AccessPolicyResponse, error) {
	var accessPolicyResponse AccessPolicyResponse
	rawURL := fmt.Sprintf("/acms/v1/scopes/%s/users/%s/policies/%s", IAM_ACCOUNT_ESCAPE+scope, userId, policyId)
	_, err := r.client.Get(rawURL, &accessPolicyResponse)
	return accessPolicyResponse, err
}

//Update ...
func (r *iampolicy) Update(scope, userId, policyId, etag string, params AccessPolicyRequest) (AccessPolicyResponse, string, error) {
	var accessPolicy AccessPolicyResponse
	rawURL := fmt.Sprintf("/acms/v1/scopes/%s/users/%s/policies/%s", IAM_ACCOUNT_ESCAPE+scope, userId, policyId)
	header := make(map[string]string)

	header["IF-Match"] = etag
	accessPolicyResp, err := r.client.Put(rawURL, params, &accessPolicy, header)
	eTag := accessPolicyResp.Header.Get("etag")
	return accessPolicy, eTag, err
}

//Delete ...
func (r *iampolicy) Delete(scope, userId, policyId string) error {
	rawURL := fmt.Sprintf("/acms/v1/scopes/%s/users/%s/policies/%s", IAM_ACCOUNT_ESCAPE+scope, userId, policyId)
	_, err := r.client.Delete(rawURL)
	return err
}
