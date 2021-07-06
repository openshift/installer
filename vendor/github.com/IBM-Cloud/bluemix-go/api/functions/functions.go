package functions

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// const ..
const (
	NamespaceTypeCFBased     = 1
	NamespaceTypeIamMigrated = 2
	NamespaceTypeIamBased    = 3

	DefaultServiceURL  = "https://gateway.watsonplatform.net/servicebroker/API/v1"
	DefaultServiceName = "ibm_cloud_functions_namespace_API"
)

// functions
type functions struct {
	client *client.Client
}

// Functions ..
type Functions interface {
	GetCloudFoundaryNamespaces() (NamespaceResponseList, error)
	DeleteNamespace(namespaceID string) (NamespaceResponse, error)
	CreateNamespace(CreateNamespaceOptions) (NamespaceResponse, error)
	GetNamespaces() (NamespaceResponseList, error)
	GetNamespace(payload GetNamespaceOptions) (NamespaceResponse, error)
	UpdateNamespace(payload UpdateNamespaceOptions) (NamespaceResponse, error)
}

func newFunctionsAPI(c *client.Client) Functions {
	return &functions{
		client: c,
	}
}

func (r *functions) GetCloudFoundaryNamespaces() (NamespaceResponseList, error) {
	var successV NamespaceResponseList
	formData := make(map[string]string)
	formData["accessToken"] = r.client.Config.UAAAccessToken[7:len(r.client.Config.UAAAccessToken)]
	formData["refreshToken"] = r.client.Config.UAARefreshToken
	_, err := r.client.PostWithForm("/bluemix/v2/authenticate", formData, &successV)
	return successV, err
}

func (r *functions) GetNamespaces() (NamespaceResponseList, error) {
	var successV NamespaceResponseList
	_, err := r.client.Get("/api/v1/namespaces", &successV)
	return successV, err
}

func (r *functions) CreateNamespace(payload CreateNamespaceOptions) (NamespaceResponse, error) {
	var successV NamespaceResponse
	_, err := r.client.Post("/api/v1/namespaces", payload, &successV)
	return successV, err
}

func (r *functions) GetNamespace(payload GetNamespaceOptions) (NamespaceResponse, error) {
	var successV NamespaceResponse
	_, err := r.client.Get(fmt.Sprintf("/api/v1/namespaces/%s", *payload.ID), &successV)
	return successV, err
}

func (r *functions) DeleteNamespace(namespaceID string) (NamespaceResponse, error) {
	var successV NamespaceResponse
	_, err := r.client.Delete(fmt.Sprintf("/api/v1/namespaces/%s", namespaceID))
	return successV, err
}

func (r *functions) UpdateNamespace(payload UpdateNamespaceOptions) (NamespaceResponse, error) {
	var successV NamespaceResponse
	_, err := r.client.Patch(fmt.Sprintf("/api/v1/namespaces/%s", *payload.ID), payload, &successV)
	return successV, err
}
