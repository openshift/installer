package sls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

// ResourceTag define
type ResourceTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ResourceFilterTag define
type ResourceFilterTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// ResourceTags tag for tag sls resource, only support project
type ResourceTags struct {
	ResourceType string        `json:"resourceType"`
	ResourceID   []string      `json:"resourceId"`
	Tags         []ResourceTag `json:"tags"`
}

// ResourceUnTags tag for untag sls resouce
type ResourceUnTags struct {
	ResourceType string   `json:"resourceType"`
	ResourceID   []string `json:"resourceId"`
	Tags         []string `json:"tags"`
}

// ResourceTagResponse used for ListTagResources
type ResourceTagResponse struct {
	ResourceType string `json:"resourceType"`
	ResourceID   string `json:"resourceId"`
	TagKey       string `json:"tagKey"`
	TagValue     string `json:"tagValue"`
}

// NewProjectTags create a project tags
func NewProjectTags(project string, tags []ResourceTag) *ResourceTags {
	return &ResourceTags{
		ResourceType: "project",
		ResourceID:   []string{project},
		Tags:         tags,
	}
}

// NewProjectUnTags create a project tags
func NewProjectUnTags(project string, tags []string) *ResourceUnTags {
	return &ResourceUnTags{
		ResourceType: "project",
		ResourceID:   []string{project},
		Tags:         tags,
	}
}

// TagResources tag specific resource
func (c *Client) TagResources(project string, tags *ResourceTags) error {
	body, err := json.Marshal(tags)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/tag"
	r, err := c.request(project, "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

// UnTagResources untag specific resource
func (c *Client) UnTagResources(project string, tags *ResourceUnTags) error {
	body, err := json.Marshal(tags)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/untag"
	r, err := c.request(project, "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

// ListTagResources list rag resources
func (c *Client) ListTagResources(project string,
	resourceType string,
	resourceIDs []string,
	tags []ResourceFilterTag,
	nextToken string) (respTags []*ResourceTagResponse, respNextToken string, err error) {
	tagsBuf, err := json.Marshal(tags)
	if err != nil {
		return nil, "", NewClientError(err)
	}
	resourceIDBuf, err := json.Marshal(resourceIDs)
	if err != nil {
		return nil, "", NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	v := url.Values{}
	v.Add("tags", string(tagsBuf))
	v.Add("resourceType", resourceType)
	v.Add("resourceId", string(resourceIDBuf))
	if nextToken != "" {
		v.Add("nextToken", nextToken)
	}
	uri := "/tags?" + v.Encode()
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, "", err
	}
	defer r.Body.Close()

	type ListTagResp struct {
		NextToken   string                 `json:"nextToken"`
		TagResource []*ResourceTagResponse `json:"tagResources"`
	}
	buf, _ := ioutil.ReadAll(r.Body)
	listTagResp := &ListTagResp{}
	if err = json.Unmarshal(buf, listTagResp); err != nil {
		err = NewClientError(err)
	}
	return listTagResp.TagResource, listTagResp.NextToken, err
}
