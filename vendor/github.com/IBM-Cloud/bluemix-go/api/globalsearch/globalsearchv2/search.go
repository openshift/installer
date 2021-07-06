package globalsearchv2

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type SearchResult struct {
	Items       []Item `json:"items"`
	MoreData    bool   `json:"more_data"`
	Token       string `json:"token"`
	FilterError bool   `json:"filter_error"`
	PartialData int    `json:"partial_data"`
}

type Item struct {
	Name        string   `json:"name,omitempty"`
	CRN         string   `json:"crn,omitempty"`
	ServiceName string   `json:"service_name,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type SearchBody struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields,omitempty"`
	Token  string   `json:"token,omitempty"`
}

type Searches interface {
	PostQuery(searchBody SearchBody) (SearchResult, error)
}

type searches struct {
	client *client.Client
}

func newSearchAPI(c *client.Client) Searches {
	return &searches{
		client: c,
	}
}

func (r *searches) PostQuery(searchBody SearchBody) (SearchResult, error) {
	searchResult := SearchResult{}
	rawURL := fmt.Sprintf("/v2/resources/search")
	_, err := r.client.Post(rawURL, &searchBody, &searchResult)
	if err != nil {
		return searchResult, err
	}
	return searchResult, nil
}
