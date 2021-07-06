package accountv2

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type GenericPaginatedResourcesHandler struct {
	resourceType reflect.Type
}

func NewAccountPaginatedResources(resource interface{}) GenericPaginatedResourcesHandler {
	return GenericPaginatedResourcesHandler{
		resourceType: reflect.TypeOf(resource),
	}
}

func (pr GenericPaginatedResourcesHandler) Resources(data []byte, curURL string) ([]interface{}, string, error) {
	var paginatedResources = struct {
		NextUrl        string          `json:"next_url"`
		ResourcesBytes json.RawMessage `json:"resources"`
	}{}

	err := json.Unmarshal(data, &paginatedResources)

	slicePtr := reflect.New(reflect.SliceOf(pr.resourceType))
	dc := json.NewDecoder(bytes.NewBuffer(paginatedResources.ResourcesBytes))
	dc.UseNumber()
	err = dc.Decode(slicePtr.Interface())
	slice := reflect.Indirect(slicePtr)

	contents := make([]interface{}, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		contents = append(contents, slice.Index(i).Interface())
	}

	return contents, paginatedResources.NextUrl, err
}
