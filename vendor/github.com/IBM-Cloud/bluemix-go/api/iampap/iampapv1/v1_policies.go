package iampapv1

import (
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type SearchParams struct {
	AccountID     string
	IAMID         string
	AccessGroupID string
	Type          string
	ServiceType   string
	Sort          string
}

func (p SearchParams) buildRequest(r *rest.Request) {
	if p.AccountID != "" {
		r.Query("account_id", p.AccountID)
	}
	if p.IAMID != "" {
		r.Query("iam_id", p.IAMID)
	}
	if p.AccessGroupID != "" {
		r.Query("access_group_id", p.AccessGroupID)
	}
	if p.Type != "" {
		r.Query("type", p.Type)
	}
	if p.ServiceType != "" {
		r.Query("service_type", p.ServiceType)
	}
	if p.Sort != "" {
		r.Query("sort", p.Sort)
	}
}

type V1PolicyRepository interface {
	List(params SearchParams) ([]Policy, error)
	Get(policyID string) (Policy, error)
	Create(policy Policy) (Policy, error)
	Update(policyID string, policy Policy, version string) (Policy, error)
	Delete(policyID string) error
}

type v1PolicyRepository struct {
	client *client.Client
}

func NewV1PolicyRepository(c *client.Client) V1PolicyRepository {
	return &v1PolicyRepository{
		client: c,
	}
}

func (r *v1PolicyRepository) List(params SearchParams) ([]Policy, error) {
	request := rest.GetRequest(*r.client.Config.Endpoint + "/v1/policies")
	params.buildRequest(request)

	response := struct {
		Policies []Policy `json:"policies"`
	}{}
	_, err := r.client.SendRequest(request, &response)
	if err != nil {
		return []Policy{}, err
	}
	return response.Policies, nil
}

func (r *v1PolicyRepository) Get(policyID string) (Policy, error) {
	var response Policy
	resp, err := r.client.Get("/v1/policies/"+policyID, &response)
	if err != nil {
		return Policy{}, err
	}
	response.Version = resp.Header.Get("ETag")
	return response, nil
}

func (r *v1PolicyRepository) Create(policy Policy) (Policy, error) {
	var response Policy
	resp, err := r.client.Post("/v1/policies", &policy, &response)
	if err != nil {
		return Policy{}, err
	}
	response.Version = resp.Header.Get("ETag")
	return response, nil
}

func (r *v1PolicyRepository) Update(policyID string, policy Policy, version string) (Policy, error) {
	var response Policy
	request := rest.PutRequest(*r.client.Config.Endpoint + "/v1/policies/" + policyID)
	request = request.Set("If-Match", version).Body(&policy)

	resp, err := r.client.SendRequest(request, &response)
	if err != nil {
		return Policy{}, err
	}
	response.Version = resp.Header.Get("Etag")
	return response, nil
}

func (r *v1PolicyRepository) Delete(policyID string) error {
	_, err := r.client.Delete("/v1/policies/" + policyID)
	if err != nil {
		return err
	}
	return nil
}
