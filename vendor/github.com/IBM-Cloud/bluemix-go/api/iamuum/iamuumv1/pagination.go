package iamuumv1

import (
	"encoding/json"
	"net/url"
	"reflect"
)

type PaginatedResourcesHandler struct {
	resourcesType reflect.Type
}

func NewPaginatedResourcesHandler(resources PaginatedResources) PaginatedResourcesHandler {
	return PaginatedResourcesHandler{
		resourcesType: reflect.TypeOf(resources).Elem(),
	}
}

func (pr PaginatedResourcesHandler) Resources(bytes []byte, curPath string) ([]interface{}, string, error) {
	paginatedResources := reflect.New(pr.resourcesType).Interface().(PaginatedResources)
	err := json.Unmarshal(bytes, paginatedResources)

	if err != nil {
		return []interface{}{}, "", err
	}

	nextPath, err := paginatedResources.NextPath()

	if err != nil {
		return []interface{}{}, "", err
	}

	return paginatedResources.Resources(), nextPath, nil
}

type PaginationHref struct {
	Href string `json:"href"`
}

type PaginationFields struct {
	First    PaginationHref `json:"first"`
	Last     PaginationHref `json:"last"`
	Next     PaginationHref `json:"next"`
	Previous PaginationHref `json:"previous"`

	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
}

func (p *PaginationFields) NextPath() (string, error) {
	if p.Next.Href == "" {
		return "", nil
	}

	u, err := url.Parse(p.Next.Href)
	if err == nil {
		u.Scheme = ""
		u.Opaque = ""
		u.Host = ""
		u.User = nil
		return u.String(), nil
	}
	return "", err
}

type PaginatedResources interface {
	NextPath() (string, error)
	Resources() []interface{}
}
