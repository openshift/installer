package containerv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type ApiKeyInfo struct {
	ID    string
	Name  string
	Email string
}

// Apikeys ...
type Apikeys interface {
	GetApiKeyInfo(clusterID string, target ClusterTargetHeader) (ApiKeyInfo, error)
	ResetApiKey(target ClusterTargetHeader) error
}

type apikeys struct {
	client *client.Client
}

func newApiKeyAPI(c *client.Client) Apikeys {
	return &apikeys{
		client: c,
	}
}

//GetApiKeyInfo ...
func (r *apikeys) GetApiKeyInfo(cluster string, target ClusterTargetHeader) (ApiKeyInfo, error) {
	retVal := ApiKeyInfo{}
	req := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, fmt.Sprintf("/v1/logging/%s/clusterkeyowner", cluster)))

	for key, value := range target.ToMap() {
		req.Set(key, value)
	}

	_, err := r.client.SendRequest(req, &retVal)
	if err != nil {
		return retVal, err
	}
	return retVal, err
}

//ResetApiKey ...
func (r *apikeys) ResetApiKey(target ClusterTargetHeader) error {
	req := rest.PostRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/keys"))

	for key, value := range target.ToMap() {
		req.Set(key, value)
	}

	_, err := r.client.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return err
}
