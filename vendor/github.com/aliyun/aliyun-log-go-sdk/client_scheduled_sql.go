package sls

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

type SqlType string
type ResourcePool string
type DataFormat string
type JobType string
type Status string
type ScheduledSQLState string

const (
	STANDARD     SqlType = "standard"
	SEARCH_QUERY SqlType = "searchQuery"
)
const (
	DEFAULT  ResourcePool = "default"
	ENHANCED ResourcePool = "enhanced"
)
const (
	LOG_TO_LOG       DataFormat = "log2log"
	LOG_TO_METRIC    DataFormat = "log2metric"
	METRIC_TO_metric DataFormat = "metric2metric"
)
const (
	ALERT_JOB         JobType = "Alert"
	REPORT_JOB        JobType = "Report"
	ETL_JOB           JobType = "ETL"
	INGESTION_JOB     JobType = "Ingestion"
	REBUILD_INDEX_JOB JobType = "RebuildIndex"
	AUDIT_JOB_JOB     JobType = "AuditJob"
	EXPORT_JOB        JobType = "Export"
	SCHEDULED_SQL_JOB JobType = "ScheduledSQL"
)

const (
	ENABLED  Status = "Enabled"
	DISABLED Status = "Disabled"
)

const (
	ScheduledSQL_RUNNING   ScheduledSQLState = "RUNNING"
	ScheduledSQL_FAILED    ScheduledSQLState = "FAILED"
	ScheduledSQL_SUCCEEDED ScheduledSQLState = "SUCCEEDED"
)

type ScheduledSQL struct {
	Name             string                     `json:"name"`
	DisplayName      string                     `json:"displayName"`
	Description      string                     `json:"description"`
	Status           Status                     `json:"status"`
	ScheduleId       string                     `json:"scheduleId"`
	Configuration    *ScheduledSQLConfiguration `json:"configuration"`
	Schedule         *Schedule                  `json:"schedule"`
	CreateTime       int64                      `json:"createTime,omitempty"`
	LastModifiedTime int64                      `json:"lastModifiedTime,omitempty"`
	Type             JobType                    `json:"type"`
}

type ScheduledSQLConfiguration struct {
	SourceLogStore      string                  `json:"sourceLogstore"`
	DestProject         string                  `json:"destProject"`
	DestEndpoint        string                  `json:"destEndpoint"`
	DestLogStore        string                  `json:"destLogstore"`
	Script              string                  `json:"script"`
	SqlType             SqlType                 `json:"sqlType"`
	ResourcePool        ResourcePool            `json:"resourcePool"`
	RoleArn             string                  `json:"roleArn"`
	DestRoleArn         string                  `json:"destRoleArn"`
	FromTimeExpr        string                  `json:"fromTimeExpr"`
	ToTimeExpr          string                  `json:"toTimeExpr"`
	MaxRunTimeInSeconds int32                   `json:"maxRunTimeInSeconds"`
	MaxRetries          int32                   `json:"maxRetries"`
	FromTime            int64                   `json:"fromTime"`
	ToTime              int64                   `json:"toTime"`
	DataFormat          DataFormat              `json:"dataFormat"`
	Parameters          *ScheduledSQLParameters `json:"parameters,omitempty"`
}

func NewScheduledSQLConfiguration() *ScheduledSQLConfiguration {
	return &ScheduledSQLConfiguration{
		SqlType:      STANDARD,
		ResourcePool: DEFAULT,
		FromTime:     0,
		ToTime:       0,
		DataFormat:   LOG_TO_LOG,
	}
}

type ScheduledSQLParameters struct {
	TimeKey    string `json:"timeKey,omitempty"`
	LabelKeys  string `json:"labelKeys,omitempty"`
	MetricKeys string `json:"metricKeys,omitempty"`
	MetricName string `json:"metricName,omitempty"`
	HashLabels string `json:"hashLabels,omitempty"`
	AddLabels  string `json:"addLabels,omitempty"`
}

func (c *Client) CreateScheduledSQL(project string, scheduledsql *ScheduledSQL) error {
	fromTime := scheduledsql.Configuration.FromTime
	toTime := scheduledsql.Configuration.ToTime
	timeRange := fromTime > 1451577600 && toTime > fromTime
	sustained := fromTime > 1451577600 && toTime == 0
	if !timeRange && !sustained {
		return fmt.Errorf("invalid fromTime: %d toTime: %d, please ensure fromTime more than 1451577600", fromTime, toTime)
	}
	body, err := json.Marshal(scheduledsql)
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

func (c *Client) DeleteScheduledSQL(project string, name string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := "/jobs/" + name
	r, err := c.request(project, "DELETE", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateScheduledSQL(project string, scheduledsql *ScheduledSQL) error {
	body, err := json.Marshal(scheduledsql)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}

	uri := "/jobs/" + scheduledsql.Name
	r, err := c.request(project, "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) GetScheduledSQL(project string, name string) (*ScheduledSQL, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + name
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	scheduledSQL := &ScheduledSQL{}
	if err = json.Unmarshal(buf, scheduledSQL); err != nil {
		err = NewClientError(err)
	}
	return scheduledSQL, err
}

func (c *Client) ListScheduledSQL(project, name, displayName string, offset, size int) (scheduledsqls []*ScheduledSQL, total, count int, error error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	v := url.Values{}
	v.Add("jobName", name)
	if displayName != "" {
		v.Add("displayName", displayName)
	}
	v.Add("jobType", "ScheduledSQL")
	v.Add("offset", fmt.Sprintf("%d", offset))
	v.Add("size", fmt.Sprintf("%d", size))

	uri := "/jobs?" + v.Encode()
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer r.Body.Close()

	type ScheduledSqlList struct {
		Total   int             `json:"total"`
		Count   int             `json:"count"`
		Results []*ScheduledSQL `json:"results"`
	}
	buf, _ := ioutil.ReadAll(r.Body)
	scheduledSqlList := &ScheduledSqlList{}
	if err = json.Unmarshal(buf, scheduledSqlList); err != nil {
		err = NewClientError(err)
	}
	return scheduledSqlList.Results, scheduledSqlList.Total, scheduledSqlList.Count, err
}

type ScheduledSQLJobInstance struct {
	InstanceId           string            `json:"instanceId"`
	JobName              string            `json:"jobName,omitempty"`
	DisplayName          string            `json:"displayName,omitempty"`
	Description          string            `json:"description,omitempty"`
	JobScheduleId        string            `json:"jobScheduleId,omitempty"`
	CreateTimeInMillis   int64             `json:"createTimeInMillis"`
	ScheduleTimeInMillis int64             `json:"scheduleTimeInMillis"`
	UpdateTimeInMillis   int64             `json:"updateTimeInMillis"`
	State                ScheduledSQLState `json:"state"`
	ErrorCode            string            `json:"errorCode"`
	ErrorMessage         string            `json:"errorMessage"`
	Summary              string            `json:"summary,omitempty"`
}

func (c *Client) GetScheduledSQLJobInstance(projectName, jobName, instanceId string, result bool) (*ScheduledSQLJobInstance, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := fmt.Sprintf("/jobs/%s/jobinstances/%s?result=%t", jobName, instanceId, result)
	r, err := c.request(projectName, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	instance := &ScheduledSQLJobInstance{}
	if err = json.Unmarshal(buf, instance); err != nil {
		err = NewClientError(err)
	}
	return instance, err
}

func (c *Client) ModifyScheduledSQLJobInstanceState(projectName, jobName, instanceId string, state ScheduledSQLState) error {
	if ScheduledSQL_RUNNING != state {
		return NewClientError(errors.New(fmt.Sprintf("Invalid state: %s, state must be RUNNING.", state)))
	}
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}

	uri := fmt.Sprintf("/jobs/%s/jobinstances/%s?state=%s", jobName, instanceId, state)
	r, err := c.request(projectName, "PUT", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

type InstanceStatus struct {
	FromTime int64
	ToTime   int64
	Offset   int64
	Size     int64
	State    ScheduledSQLState
}

func (c *Client) ListScheduledSQLJobInstances(projectName, jobName string, status *InstanceStatus) (instances []*ScheduledSQLJobInstance, total, count int64, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	v := url.Values{}
	v.Add("jobType", "ScheduledSQL")
	v.Add("start", fmt.Sprintf("%d", status.FromTime))
	v.Add("end", fmt.Sprintf("%d", status.ToTime))
	v.Add("offset", fmt.Sprintf("%d", status.Offset))
	v.Add("size", fmt.Sprintf("%d", status.Size))
	if status.State != "" {
		v.Add("state", string(status.State))
	}

	uri := fmt.Sprintf("/jobs/%s/jobinstances?%s", jobName, v.Encode())
	r, err := c.request(projectName, "GET", uri, h, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer r.Body.Close()

	type ScheduledSqlJobInstances struct {
		Total   int64                      `json:"total"`
		Count   int64                      `json:"count"`
		Results []*ScheduledSQLJobInstance `json:"results"`
	}
	buf, _ := ioutil.ReadAll(r.Body)
	jobInstances := &ScheduledSqlJobInstances{}
	if err = json.Unmarshal(buf, jobInstances); err != nil {
		err = NewClientError(err)
	}
	return jobInstances.Results, jobInstances.Total, jobInstances.Count, err
}
