package packngo

import (
	"fmt"
)

const (
	apiKeyUserBasePath    = "/user/api-keys"
	apiKeyProjectBasePath = "/projects/%s/api-keys"
)

// APIKeyService interface defines available device methods
type APIKeyService interface {
	UserList(*ListOptions) ([]APIKey, *Response, error)
	ProjectList(string, *ListOptions) ([]APIKey, *Response, error)
	UserGet(string, *GetOptions) (*APIKey, error)
	ProjectGet(string, string, *GetOptions) (*APIKey, error)
	Create(*APIKeyCreateRequest) (*APIKey, *Response, error)
	Delete(string) (*Response, error)
}

type apiKeyRoot struct {
	APIKeys []APIKey `json:"api_keys"`
}

type APIKey struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Token       string   `json:"token"`
	ReadOnly    bool     `json:"read_only"`
	Created     string   `json:"created_at"`
	Updated     string   `json:"updated_at"`
	User        *User    `json:"user"`
	Project     *Project `json:"project"`
}

// APIKeyCreateRequest type used to create an api key
type APIKeyCreateRequest struct {
	Description string `json:"description"`
	ReadOnly    bool   `json:"read_only"`
	ProjectID   string `json:"-"`
}

func (s APIKeyCreateRequest) String() string {
	return Stringify(s)
}

// APIKeyServiceOp implements APIKeyService
type APIKeyServiceOp struct {
	client *Client
}

func (s *APIKeyServiceOp) list(url string, lopts *ListOptions) ([]APIKey, *Response, error) {
	root := new(apiKeyRoot)
	params := createListOptionsURL(lopts)
	paramURL := fmt.Sprintf("%s?%s", url, params)

	resp, err := s.client.DoRequest("GET", paramURL, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.APIKeys, resp, err
}

// ProjectList lists api keys of a project
func (s *APIKeyServiceOp) ProjectList(projectID string, lopts *ListOptions) ([]APIKey, *Response, error) {
	return s.list(fmt.Sprintf(apiKeyProjectBasePath, projectID), lopts)

}

// UserList returns a user's api keys
func (s *APIKeyServiceOp) UserList(lopts *ListOptions) ([]APIKey, *Response, error) {
	return s.list(apiKeyUserBasePath, lopts)
}

// ProjectGet returns an api key by id
func (s *APIKeyServiceOp) ProjectGet(projectID, apiKeyID string, getOpt *GetOptions) (*APIKey, error) {
	var lopts *ListOptions
	if getOpt != nil {
		lopts = &ListOptions{Includes: getOpt.Includes, Excludes: getOpt.Excludes}
	}
	pkeys, _, err := s.ProjectList(projectID, lopts)
	if err != nil {
		return nil, err
	}
	for _, k := range pkeys {
		if k.ID == apiKeyID {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("Project (%s) API key %s not found", projectID, apiKeyID)
}

// UserGet returns a project api key by id
func (s *APIKeyServiceOp) UserGet(apiKeyID string, getOpt *GetOptions) (*APIKey, error) {
	var lopts *ListOptions
	if getOpt != nil {
		lopts = &ListOptions{Includes: getOpt.Includes, Excludes: getOpt.Excludes}
	}
	ukeys, _, err := s.UserList(lopts)
	if err != nil {
		return nil, err
	}
	for _, k := range ukeys {
		if k.ID == apiKeyID {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("User API key %s not found", apiKeyID)
}

// Create creates a new api key
func (s *APIKeyServiceOp) Create(createRequest *APIKeyCreateRequest) (*APIKey, *Response, error) {
	path := apiKeyUserBasePath
	if createRequest.ProjectID != "" {
		path = fmt.Sprintf(apiKeyProjectBasePath, createRequest.ProjectID)
	}
	apiKey := new(APIKey)

	resp, err := s.client.DoRequest("POST", path, createRequest, apiKey)
	if err != nil {
		return nil, resp, err
	}

	return apiKey, resp, err
}

// Delete deletes an api key
func (s *APIKeyServiceOp) Delete(apiKeyID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", apiKeyUserBasePath, apiKeyID)
	return s.client.DoRequest("DELETE", path, nil, nil)
}
