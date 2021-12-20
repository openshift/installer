package catalog

import (
	"encoding/json"
	"reflect"
	"strings"
)

type ResourceCatalogPaginatedResourcesHandler struct {
	resourceType reflect.Type
	baseURL      string
}

func NewResourceCatalogPaginatedResources(resource interface{}, baseURL string) ResourceCatalogPaginatedResourcesHandler {
	return ResourceCatalogPaginatedResourcesHandler{
		resourceType: reflect.TypeOf(resource),
		baseURL:      baseURL,
	}
}

func (pr ResourceCatalogPaginatedResourcesHandler) Resources(bytes []byte, curPath string) ([]interface{}, string, error) {
	var paginatedResources = struct {
		NextUrl        string          `json:"next"`
		ResourcesBytes json.RawMessage `json:"resources"`
	}{}

	err := json.Unmarshal(bytes, &paginatedResources)

	slicePtr := reflect.New(reflect.SliceOf(pr.resourceType))
	dc := json.NewDecoder(strings.NewReader(string(paginatedResources.ResourcesBytes)))
	dc.UseNumber()
	err = dc.Decode(slicePtr.Interface())
	slice := reflect.Indirect(slicePtr)

	contents := make([]interface{}, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		contents = append(contents, slice.Index(i).Interface())
	}
	//The next URL in response is a full qualified URL like https://resource-catalog.stage1.ng.bluemix.net/api/v1?_offset=50&languages=en_US%2Cen
	//So need to cut the baseURL from it.
	index := strings.Index(paginatedResources.NextUrl, pr.baseURL)
	//NextUrl contains baseURL, means need to cut
	if index != -1 {
		url := paginatedResources.NextUrl[index+len(pr.baseURL):]
		return contents, url, err
	}
	return contents, paginatedResources.NextUrl, err
}
