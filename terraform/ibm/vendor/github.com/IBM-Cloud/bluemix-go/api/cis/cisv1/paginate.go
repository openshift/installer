package cisv1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type GenericPaginatedResourcesHandler struct {
	resourceType reflect.Type
}

func NewDNSPaginatedResources(resource interface{}) GenericPaginatedResourcesHandler {
	return GenericPaginatedResourcesHandler{
		resourceType: reflect.TypeOf(resource),
	}
}

func (pr GenericPaginatedResourcesHandler) Resources(bytes []byte, curURL string) ([]interface{}, string, error) {
	var paginatedResources = struct {
		ResultInfo struct {
			Page       int `json:"page"`
			TotalPages int `json:"total_pages"`
		} `json:"result_info"`
		Result json.RawMessage `json:"result"`
	}{}

	if err := json.Unmarshal(bytes, &paginatedResources); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal paginated response as json: %s", err)
	}

	slicePtr := reflect.New(reflect.SliceOf(pr.resourceType))
	dc := json.NewDecoder(strings.NewReader(string(paginatedResources.Result)))
	dc.UseNumber()
	if err := dc.Decode(slicePtr.Interface()); err != nil {
		return nil, "", fmt.Errorf("failed to decode paginated objects as %T: %s", pr.resourceType, err)
	}
	slice := reflect.Indirect(slicePtr)

	contents := make([]interface{}, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		contents = append(contents, slice.Index(i).Interface())
	}

	if paginatedResources.ResultInfo.Page >= paginatedResources.ResultInfo.TotalPages {
		return contents, "", nil
	}

	return contents, strings.Replace(curURL, fmt.Sprintf("page=%d", paginatedResources.ResultInfo.Page), fmt.Sprintf("page=%d", paginatedResources.ResultInfo.Page+1), 1), nil
}
