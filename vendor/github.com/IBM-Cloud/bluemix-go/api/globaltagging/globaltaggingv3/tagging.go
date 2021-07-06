package globaltaggingv3

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type TaggingResult struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name string `json:"name"`
}

type TagUpdateResult struct {
	Results []TagResult `json:"results"`
}

type TagResult struct {
	ResourceID  string `json:"resource_id"`
	IsError     string `json:"isError"`
	Response    string `json:"response"`
	Message     string `json:"message"`
	Code        string `json:"code"`
	Level       string `json:"level"`
	HttpCode    int    `json:"httpCode"`
	Description string `json:"description"`
	MoreInfo    string `json:"more_info"`
}

type TaggingBody struct {
	TagResources []TagResource `json:"resources"`
	TagName      string        `json:"tag_name,omitempty"`
	TagNames     []string      `json:"tag_names,omitempty"`
}

type TagResource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type,omitempty"`
}

type Tags interface {
	GetTags(resourceID string) (TaggingResult, error)
	AttachTags(resourceID string, taglist []string) (TagUpdateResult, error)
	DetachTags(resourceID string, taglist []string) (TagUpdateResult, error)
	DeleteTag(tag string) (TagUpdateResult, error)
}

type tags struct {
	client *client.Client
}

func newTaggingAPI(c *client.Client) Tags {
	return &tags{
		client: c,
	}
}

func (r *tags) GetTags(resourceID string) (TaggingResult, error) {
	taggingResult := TaggingResult{}
	query := fmt.Sprintf("?attached_to=%v", resourceID)
	rawURL := fmt.Sprintf("/v3/tags" + query)
	_, err := r.client.Get(rawURL, &taggingResult)
	if err != nil {
		return taggingResult, err
	}
	return taggingResult, nil
}

func (r *tags) AttachTags(resourceID string, taglist []string) (TagUpdateResult, error) {
	tagUpdateResult := TagUpdateResult{}
	taggingBody := TaggingBody{
		TagResources: []TagResource{
			{ResourceID: resourceID},
		},
		TagNames: taglist,
	}
	rawURL := fmt.Sprintf("/v3/tags/attach")
	_, err := r.client.Post(rawURL, &taggingBody, &tagUpdateResult)
	if err != nil {
		return tagUpdateResult, err
	}
	return tagUpdateResult, nil

}

func (r *tags) DetachTags(resourceID string, taglist []string) (TagUpdateResult, error) {
	tagUpdateResult := TagUpdateResult{}
	taggingBody := TaggingBody{
		TagResources: []TagResource{
			{ResourceID: resourceID},
		},
		TagNames: taglist,
	}
	rawURL := fmt.Sprintf("/v3/tags/detach")
	_, err := r.client.Post(rawURL, &taggingBody, &tagUpdateResult)
	if err != nil {
		return tagUpdateResult, err
	}
	return tagUpdateResult, nil

}

func (r *tags) DeleteTag(tagin string) (TagUpdateResult, error) {
	tagUpdateResult := TagUpdateResult{}
	rawURL := fmt.Sprintf("/v3/tags/%s", tagin)
	_, err := r.client.Delete(rawURL, &tagUpdateResult)
	if err != nil {
		return tagUpdateResult, err
	}
	return tagUpdateResult, nil

}
