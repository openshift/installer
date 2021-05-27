package iampapv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type AuthorizationPolicy struct {
	ID        string                  `json:"id,omitempty"`
	Roles     []models.PolicyRole     `json:"roles"`
	Resources []models.PolicyResource `json:"resources"`
	Subjects  []models.PolicyResource `json:"subjects"`
	Type      string                  `json:"type,omitempty"`
	Version   string                  `json:"-"`
}

type AuthorizationPolicySearchQuery struct {
	SubjectID     string
	Type          string
	AccessGroupID string
}

func (q *AuthorizationPolicySearchQuery) setQuery(r *rest.Request) {
	if q.SubjectID != "" {
		r.Query("subjectId", q.SubjectID)
	}
	if q.Type != "" {
		r.Query("type", q.Type)
	}
	if q.AccessGroupID != "" {
		r.Query("accessGroupId", q.AccessGroupID)
	}
}

const (
	AuthorizationPolicyType = "authorization"
	AccessPolicyType        = "access"
)

type AuthorizationPolicyRepository interface {
	List(accountID string, query *AuthorizationPolicySearchQuery) ([]AuthorizationPolicy, error)
	Get(accountID string, policyID string) (AuthorizationPolicy, error)
	Create(accountID string, policy AuthorizationPolicy) (AuthorizationPolicy, error)
	Update(accountID string, policyID string, policy AuthorizationPolicy, version string) (AuthorizationPolicy, error)
	Delete(accountID string, policyID string) error
	// Purge(accountID string, request DeleteAuthorizationPolicyRequest) (error)
}

type authorizationPolicyRepository struct {
	client *client.Client
}

func NewAuthorizationPolicyRepository(c *client.Client) AuthorizationPolicyRepository {
	return &authorizationPolicyRepository{
		client: c,
	}
}

type listAuthorizationPolicyResponse struct {
	Policies []AuthorizationPolicy `json:"policies"`
}

func (r *authorizationPolicyRepository) List(accountID string, query *AuthorizationPolicySearchQuery) ([]AuthorizationPolicy, error) {
	request := rest.GetRequest(*r.client.Config.Endpoint + fmt.Sprintf("/acms/v2/accounts/%s/policies", accountID))

	if query != nil {
		query.setQuery(request)
	}

	var response listAuthorizationPolicyResponse
	_, err := r.client.SendRequest(request, &response)
	if err != nil {
		return []AuthorizationPolicy{}, err
	}
	return response.Policies, nil
}

func (r *authorizationPolicyRepository) Get(accountID string, policyID string) (AuthorizationPolicy, error) {
	var policy AuthorizationPolicy

	resp, err := r.client.Get(fmt.Sprintf("/acms/v2/accounts/%s/policies/%s", accountID, policyID), &policy)
	if err != nil {
		return AuthorizationPolicy{}, err
	}
	policy.Version = resp.Header.Get("Etag")
	return policy, nil
}

func (r *authorizationPolicyRepository) Create(accountID string, policy AuthorizationPolicy) (AuthorizationPolicy, error) {
	var policyCreated AuthorizationPolicy

	resp, err := r.client.Post(fmt.Sprintf("/acms/v2/accounts/%s/policies", accountID), &policy, &policyCreated)
	if err != nil {
		return AuthorizationPolicy{}, err
	}
	policyCreated.Version = resp.Header.Get("Etag")
	return policyCreated, nil
}

func (r *authorizationPolicyRepository) Update(accountID string, policyID string, policy AuthorizationPolicy, version string) (AuthorizationPolicy, error) {
	var policyUpdated AuthorizationPolicy
	request := rest.PutRequest(*r.client.Config.Endpoint + fmt.Sprintf("/acms/v2/accounts/%s/policies/%s", accountID, policyID)).Body(policy)
	if version != "" {
		request = request.Set("If-Match", version)
	}

	resp, err := r.client.SendRequest(request, &policyUpdated)
	if err != nil {
		return AuthorizationPolicy{}, err
	}
	policyUpdated.Version = resp.Header.Get("Etag")

	return policyUpdated, nil
}

func (r *authorizationPolicyRepository) Delete(accountID string, policyID string) error {
	_, err := r.client.Delete(fmt.Sprintf("/acms/v1/policies/%s?scope=%s", policyID, "a/"+accountID))
	if err != nil {
		return err
	}
	return nil
}
