package sls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

const ResourceTypeUserDefine = "userdefine"

type Resource struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Schema         string `json:"schema"`
	Description    string `json:"description"`
	ExtInfo        string `json:"extInfo"`
	CreateTime     int64  `json:"createTime"`
	LastModifyTime int64  `json:"lastModifyTime"`
}

type ResourceSchema struct {
	Schema []*ResourceSchemaItem `json:"schema"`
}

type ResourceSchemaItem struct {
	Column   string      `json:"column"`
	Desc     string      `json:"desc"`
	ExtInfo  interface{} `json:"ext_info"`
	Required bool        `json:"required"`
	Type     string      `json:"type"`
}

func (rs *ResourceSchema) ToString() string {
	rsBytes, _ := json.Marshal(rs)
	return string(rsBytes)
}

func (rs *ResourceSchema) FromJsonString(schema string) error {
	return json.Unmarshal([]byte(schema), rs)
}

func (c *Client) CreateResourceString(resourceStr string) error {
	body := []byte(resourceStr)

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/resources"
	r, err := c.request("", "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) CreateResource(resource *Resource) error {
	body, err := json.Marshal(resource)
	if err != nil {
		return NewClientError(err)
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/resources"
	r, err := c.request("", "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateResource(resource *Resource) error {
	body, err := json.Marshal(resource)
	if err != nil {
		return NewClientError(err)
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/resources/" + resource.Name
	r, err := c.request("", "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateResourceString(resourceName, resourceStr string) error {
	body := []byte(resourceStr)

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/resources/" + resourceName
	r, err := c.request("", "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) DeleteResource(name string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := "/resources/" + name
	r, err := c.request("", "DELETE", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) GetResource(name string) (resource *Resource, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/resources/" + name
	r, err := c.request("", "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	resource = &Resource{}
	if err = json.Unmarshal(buf, resource); err != nil {
		err = NewClientError(err)
	}
	return resource, err
}

func (c *Client) GetResourceString(name string) (resource string, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/resources/" + name
	r, err := c.request("", "GET", uri, h, nil)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	return string(buf), err
}

func (c *Client) ListResource(resourceType string, resourceName string, offset, size int) (resourceList []*Resource, count, total int, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
		"offset":            strconv.Itoa(offset),
		"size":              strconv.Itoa(size),
	}
	uri := fmt.Sprintf("/resources?type=%s&names=%s", resourceType, resourceName)
	r, err := c.request("", "GET", uri, h, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer r.Body.Close()
	type ListResourceResponse struct {
		ResourceList []*Resource `json:"items"`
		Total        int         `json:"total"`
		Count        int         `json:"count"`
	}

	buf, _ := ioutil.ReadAll(r.Body)
	resources := &ListResourceResponse{}
	if err = json.Unmarshal(buf, resources); err != nil {
		err = NewClientError(err)
	}
	return resources.ResourceList, resources.Count, resources.Total, err
}
