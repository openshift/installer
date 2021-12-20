package certificatemanager

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type GenericPaginatedResourcesHandler struct {
	resourceType reflect.Type
}

func NewCMSPaginatedResources(resource interface{}) GenericPaginatedResourcesHandler {
	return GenericPaginatedResourcesHandler{
		resourceType: reflect.TypeOf(resource),
	}
}

func (pr GenericPaginatedResourcesHandler) Resources(bytes []byte, curURL string) ([]interface{}, string, error) {
	var paginatedResources = struct {
		NextPageInfo struct {
			StartDocId        string `json:"startWithDocId"`
			StartOrderByValue string `json:"startWithOrderByValue"`
		} `json:"nextPageInfo"`
		Certificates json.RawMessage `json:"certificates"`
		TotalDocs    int             `json:"totalScannedDocs"`
	}{}

	if err := json.Unmarshal(bytes, &paginatedResources); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal paginated response as json: %s", err)
	}
	slicePtr := reflect.New(reflect.SliceOf(pr.resourceType))
	dc := json.NewDecoder(strings.NewReader(string(paginatedResources.Certificates)))
	dc.UseNumber()
	if err := dc.Decode(slicePtr.Interface()); err != nil {
		return nil, "", fmt.Errorf("failed to decode paginated objects as %T: %s", pr.resourceType, err)
	}
	slice := reflect.Indirect(slicePtr)

	contents := make([]interface{}, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		contents = append(contents, slice.Index(i).Interface())
	}
	if paginatedResources.NextPageInfo.StartDocId == "" && paginatedResources.NextPageInfo.StartOrderByValue == "" {
		return contents, "", nil
	}
	urlprefix := strings.Split(curURL, "?")[0]
	nextURL := fmt.Sprintf("%s?page_number=1&&page_size=200&&start_from_document_id=%s&&start_from_orderby_value=%s", urlprefix, strings.Replace(paginatedResources.NextPageInfo.StartDocId, "/", "%2F", -1), paginatedResources.NextPageInfo.StartOrderByValue)
	return contents, nextURL, nil
}
