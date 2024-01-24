package sls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type ResourceRecord struct {
	Id             string `json:"id"`
	Tag            string `json:"tag"`
	Value          string `json:"value"`
	CreateTime     int64  `json:"createTime"`
	LastModifyTime int64  `json:"lastModifyTime"`
}

func (c *Client) CreateResourceRecordString(resourceName, recordStr string) error {
	body := []byte(recordStr)

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/resources/%s/records", resourceName)
	r, err := c.request("", "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) CreateResourceRecord(resourceName string, record *ResourceRecord) error {
	body, err := json.Marshal(record)
	if err != nil {
		return NewClientError(err)
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}
	uri := fmt.Sprintf("/resources/%s/records", resourceName)
	r, err := c.request("", "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateResourceRecord(resourceName string, record *ResourceRecord) error {
	body, err := json.Marshal(record)
	if err != nil {
		return NewClientError(err)
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/resources/%s/records/%s", resourceName, record.Id)
	r, err := c.request("", "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateResourceRecordString(resourceName, recordStr string) error {
	body := []byte(recordStr)

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/resources/%s/records", resourceName)
	r, err := c.request("", "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) DeleteResourceRecord(resourceName, recordId string) error {

	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/resources/%s/records?ids=%s", resourceName, recordId)
	r, err := c.request("", "DELETE", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) GetResourceRecord(resourceName, recordId string) (record *ResourceRecord, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := fmt.Sprintf("/resources/%s/records/%s", resourceName, recordId)
	r, err := c.request("", "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	record = &ResourceRecord{}
	if err = json.Unmarshal(buf, record); err != nil {
		err = NewClientError(err)
	}
	return record, err
}

func (c *Client) GetResourceRecordString(resourceName, recordId string) (recordStr string, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := fmt.Sprintf("/resources/%s/records/%s", resourceName, recordId)
	r, err := c.request("", "GET", uri, h, nil)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	return string(buf), err
}

func (c *Client) ListResourceRecord(resourceName string, offset, size int) (recordList []*ResourceRecord, count, total int, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
		"offset":            strconv.Itoa(offset),
		"size":              strconv.Itoa(size),
	}
	uri := fmt.Sprintf("/resources/%s/records", resourceName)
	r, err := c.request("", "GET", uri, h, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer r.Body.Close()
	type ListResourceRecordResponse struct {
		ResourceRecordList []*ResourceRecord `json:"items"`
		Total              int               `json:"total"`
		Count              int               `json:"count"`
	}

	buf, _ := ioutil.ReadAll(r.Body)
	resources := &ListResourceRecordResponse{}
	if err = json.Unmarshal(buf, resources); err != nil {
		err = NewClientError(err)
	}
	return resources.ResourceRecordList, resources.Count, resources.Total, err
}
