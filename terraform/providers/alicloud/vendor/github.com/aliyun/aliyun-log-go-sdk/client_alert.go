package sls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
)

// SavedSearch ...
type SavedSearch struct {
	SavedSearchName string `json:"savedsearchName"`
	SearchQuery     string `json:"searchQuery"`
	Logstore        string `json:"logstore"`
	Topic           string `json:"topic"`
	DisplayName     string `json:"displayName"`
}

type ResponseSavedSearchItem struct {
	SavedSearchName string `json:"savedsearchName"`
	DisplayName     string `json:"displayName"`
}

const (
	NotificationTypeSMS           = "SMS"
	NotificationTypeWebhook       = "Webhook"
	NotificationTypeDingTalk      = "DingTalk"
	NotificationTypeEmail         = "Email"
	NotificationTypeMessageCenter = "MessageCenter"
)

const (
	CountConditionKey = "__count__"
)

type Severity int

const (
	Report   Severity = 2
	Low      Severity = 4
	Medium   Severity = 6
	High     Severity = 8
	Critical Severity = 10
)

const (
	JoinTypeCross        = "cross_join"
	JoinTypeInner        = "inner_join"
	JoinTypeLeft         = "left_join"
	JoinTypeRight        = "right_join"
	JoinTypeFull         = "full_join"
	JoinTypeLeftExclude  = "left_exclude"
	JoinTypeRightExclude = "right_exclude"
	JoinTypeConcat       = "concat"
	JoinTypeNo           = "no_join"
)

const (
	GroupTypeNoGroup    = "no_group"
	GroupTypeLabelsAuto = "labels_auto"
	GroupTypeCustom     = "custom"
)

const (
	ScheduleTypeFixedRate = "FixedRate"
	ScheduleTypeHourly    = "Hourly"
	ScheduleTypeDaily     = "Daily"
	ScheduleTypeWeekly    = "Weekly"
	ScheduleTypeCron      = "Cron"
	ScheduleTypeDayRun    = "DryRun"
	ScheduleTypeResident  = "Resident"
)

const (
	StoreTypeLog    = "log"
	StoreTypeMetric = "metric"
	StoreTypeMeta   = "meta"
)

// SeverityConfiguration severity config by group
type SeverityConfiguration struct {
	Severity      Severity               `json:"severity"`
	EvalCondition ConditionConfiguration `json:"evalCondition"`
}

type ConditionConfiguration struct {
	Condition      string `json:"condition"`
	CountCondition string `json:"countCondition"`
}

type JoinConfiguration struct {
	Type      string `json:"type"`
	Condition string `json:"condition"`
}

type GroupConfiguration struct {
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Token struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Hide        bool   `json:"hide"`
}

type TemplateConfiguration struct {
	Id          string            `json:"id"`
	Type        string            `json:"type"`
	Version     string            `json:"version"`
	Lang        string            `json:"lang"`
	Tokens      map[string]string `json:"tokens"`
	Annotations map[string]string `json:"annotations"`
}

type PolicyConfiguration struct {
	UseDefault     bool   `json:"useDefault"`
	RepeatInterval string `json:"repeatInterval"`
	AlertPolicyId  string `json:"alertPolicyId"`
	ActionPolicyId string `json:"actionPolicyId"`
}

type Alert struct {
	Name             string              `json:"name"`
	DisplayName      string              `json:"displayName"`
	Description      string              `json:"description"`
	State            string              `json:"state"`
	Status           string              `json:"status"`
	Configuration    *AlertConfiguration `json:"configuration"`
	Schedule         *Schedule           `json:"schedule"`
	CreateTime       int64               `json:"createTime,omitempty"`
	LastModifiedTime int64               `json:"lastModifiedTime,omitempty"`
}

func (alert *Alert) MarshalJSON() ([]byte, error) {
	body := map[string]interface{}{
		"name":          alert.Name,
		"displayName":   alert.DisplayName,
		"description":   alert.Description,
		"state":         alert.State,
		"status":        alert.Status,
		"configuration": alert.Configuration,
		"schedule":      alert.Schedule,
		"type":          "Alert",
	}
	return json.Marshal(body)
}

type AlertQuery struct {
	ChartTitle   string `json:"chartTitle"`
	LogStore     string `json:"logStore"`
	Query        string `json:"query"`
	TimeSpanType string `json:"timeSpanType"`
	Start        string `json:"start"`
	End          string `json:"end"`

	StoreType   string `json:"storeType"`
	Project     string `json:"project"`
	Store       string `json:"store"`
	Region      string `json:"region"`
	RoleArn     string `json:"roleArn"`
	DashboardId string `json:"dashboardId"`
}

type Notification struct {
	Type       string            `json:"type"`
	Content    string            `json:"content"`
	EmailList  []string          `json:"emailList,omitempty"`
	Method     string            `json:"method,omitempty"`
	MobileList []string          `json:"mobileList,omitempty"`
	ServiceUri string            `json:"serviceUri,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

type Schedule struct {
	Type           string `json:"type"`
	Interval       string `json:"interval"`
	CronExpression string `json:"cronExpression"`
	Delay          int32  `json:"delay"`
	DayOfWeek      int32  `json:"dayOfWeek"`
	Hour           int32  `json:"hour"`
}

type AlertConfiguration struct {
	Condition        string          `json:"condition"`
	MuteUntil        int64           `json:"muteUntil"`
	NotificationList []*Notification `json:"notificationList"`
	NotifyThreshold  int32           `json:"notifyThreshold"`
	Throttling       string          `json:"throttling"`

	Version               string                 `json:"version"`
	Type                  string                 `json:"type"`
	TemplateConfiguration *TemplateConfiguration `json:"templateConfiguration"`

	Dashboard              string                   `json:"dashboard"`
	Threshold              int                      `json:"threshold"`
	NoDataFire             bool                     `json:"noDataFire"`
	NoDataSeverity         Severity                 `json:"noDataSeverity"`
	SendResolved           bool                     `json:"sendResolved"`
	QueryList              []*AlertQuery            `json:"queryList"`
	Annotations            []*Tag                   `json:"annotations"`
	Labels                 []*Tag                   `json:"labels"`
	SeverityConfigurations []*SeverityConfiguration `json:"severityConfigurations"`

	JoinConfigurations []*JoinConfiguration `json:"joinConfigurations"`
	GroupConfiguration GroupConfiguration   `json:"groupConfiguration"`

	PolicyConfiguration PolicyConfiguration `json:"policyConfiguration"`
}

func (c *Client) CreateSavedSearch(project string, savedSearch *SavedSearch) error {
	body, err := json.Marshal(savedSearch)
	if err != nil {
		return NewClientError(err)
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/savedsearches"
	r, err := c.request(project, "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateSavedSearch(project string, savedSearch *SavedSearch) error {
	body, err := json.Marshal(savedSearch)
	if err != nil {
		return NewClientError(err)
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/savedsearches/" + savedSearch.SavedSearchName
	r, err := c.request(project, "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) DeleteSavedSearch(project string, savedSearchName string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := "/savedsearches/" + savedSearchName
	r, err := c.request(project, "DELETE", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) GetSavedSearch(project string, savedSearchName string) (*SavedSearch, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := "/savedsearches/" + savedSearchName
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	savedSearch := &SavedSearch{}
	if err = json.Unmarshal(buf, savedSearch); err != nil {
		err = NewClientError(err)
	}
	return savedSearch, err
}

func (c *Client) ListSavedSearch(project string, savedSearchName string, offset, size int) (savedSearches []string, total int, count int, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
		"savedsearchName":   savedSearchName,
		"offset":            strconv.Itoa(offset),
		"size":              strconv.Itoa(size),
	}

	uri := "/savedsearches"
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer r.Body.Close()

	type ListSavedSearch struct {
		Total         int      `json:"total"`
		Count         int      `json:"count"`
		Savedsearches []string `json:"savedsearches"`
	}

	buf, _ := ioutil.ReadAll(r.Body)
	listSavedSearch := &ListSavedSearch{}
	if err = json.Unmarshal(buf, listSavedSearch); err != nil {
		err = NewClientError(err)
	}
	return listSavedSearch.Savedsearches, listSavedSearch.Total, listSavedSearch.Count, err
}

func (c *Client) ListSavedSearchV2(project string, savedSearchName string, offset, size int) (savedSearches []string, savedsearchItems []ResponseSavedSearchItem, total int, count int, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
		"savedsearchName":   savedSearchName,
		"offset":            strconv.Itoa(offset),
		"size":              strconv.Itoa(size),
	}

	uri := "/savedsearches"
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, nil, 0, 0, err
	}
	defer r.Body.Close()

	type ListSavedSearch struct {
		Total            int                       `json:"total"`
		Count            int                       `json:"count"`
		Savedsearches    []string                  `json:"savedsearches"`
		SavedsearchItems []ResponseSavedSearchItem `json:"savedsearchItems"`
	}

	buf, _ := ioutil.ReadAll(r.Body)
	listSavedSearch := &ListSavedSearch{}
	if err = json.Unmarshal(buf, listSavedSearch); err != nil {
		err = NewClientError(err)
	}
	return listSavedSearch.Savedsearches, listSavedSearch.SavedsearchItems, listSavedSearch.Total, listSavedSearch.Count, err
}

func (c *Client) CreateAlert(project string, alert *Alert) error {
	body, err := json.Marshal(alert)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/jobs"
	r, err := c.request(project, "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) CreateAlertString(project string, alert string) error {
	body := []byte(alert)
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/jobs"
	r, err := c.request(project, "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateAlert(project string, alert *Alert) error {
	body, err := json.Marshal(alert)
	if err != nil {
		return NewClientError(err)
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/jobs/" + alert.Name
	r, err := c.request(project, "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateAlertString(project string, alertName, alert string) error {
	body := []byte(alert)

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/jobs/" + alertName
	r, err := c.request(project, "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) DeleteAlert(project string, alertName string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := "/jobs/" + alertName
	r, err := c.request(project, "DELETE", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) DisableAlert(project string, alertName string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := fmt.Sprintf("/jobs/%s?action=disable", alertName)
	r, err := c.request(project, "PUT", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) EnableAlert(project string, alertName string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := fmt.Sprintf("/jobs/%s?action=enable", alertName)
	r, err := c.request(project, "PUT", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) GetAlert(project string, alertName string) (*Alert, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + alertName
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	alert := &Alert{}
	if err = json.Unmarshal(buf, alert); err != nil {
		err = NewClientError(err)
	}
	return alert, err
}

func (c *Client) GetAlertString(project string, alertName string) (string, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + alertName
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	return string(buf), err
}

func (c *Client) ListAlert(project, alertName, dashboard string, offset, size int) (alerts []*Alert, total int, count int, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	v := url.Values{}
	v.Add("jobName", alertName)
	v.Add("jobType", "Alert")
	v.Add("offset", fmt.Sprintf("%d", offset))
	v.Add("size", fmt.Sprintf("%d", size))
	if dashboard != "" {
		v.Add("resourceProvider", dashboard)
	}
	uri := "/jobs?" + v.Encode()
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer r.Body.Close()

	type AlertList struct {
		Total   int      `json:"total"`
		Count   int      `json:"count"`
		Results []*Alert `json:"results"`
	}
	buf, _ := ioutil.ReadAll(r.Body)
	listAlert := &AlertList{}
	if err = json.Unmarshal(buf, listAlert); err != nil {
		err = NewClientError(err)
	}
	return listAlert.Results, listAlert.Total, listAlert.Count, err
}
