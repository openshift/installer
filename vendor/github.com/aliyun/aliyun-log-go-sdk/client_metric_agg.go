package sls

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	// MetricAggRulesSQL sql type
	MetricAggRulesSQL = "sql"
	// MetricAggRulesPromQL promql type
	MetricAggRulesPromQL = "promql"
)

type MetricAggRules struct {
	ID   string
	Name string
	Desc string

	SrcStore           string
	SrcAccessKeyID     string // ETL_STS_DEFAULT
	SrcAccessKeySecret string // acs:ram::${aliuid}:role/aliyunlogetlrole

	DestEndpoint        string // same region, inner endpoint; different region, public endpoint
	DestProject         string
	DestStore           string
	DestAccessKeyID     string // ETL_STS_DEFAULT
	DestAccessKeySecret string // acs:ram::${aliuid}:role/aliyunlogetlrole

	AggRules []MetricAggRuleItem
}

type MetricAggRuleItem struct {
	Name        string
	QueryType   string
	Query       string
	TimeName    string
	MetricNames []string
	LabelNames  map[string]string

	BeginUnixTime int64
	EndUnixTime   int64
	Interval      int64
	DelaySeconds  int64
}

func (c *Client) getScheduledSQLParams(aggRules []MetricAggRuleItem) (map[string]string, error) {
	params := make(map[string]string)
	params["sls.config.job_mode"] = `{"type":"ml","source":"ScheduledSQL"}`

	var aggRuleJsons []interface{}

	for _, aggRule := range aggRules {
		aggRuleMap := make(map[string]interface{})

		aggRuleMap["rule_name"] = aggRule.Name

		advancedQueryMap := make(map[string]interface{})
		advancedQueryMap["type"] = aggRule.QueryType
		advancedQueryMap["query"] = aggRule.Query
		advancedQueryMap["time_name"] = aggRule.TimeName
		advancedQueryMap["metric_names"] = aggRule.MetricNames
		advancedQueryMap["labels"] = aggRule.LabelNames
		aggRuleMap["advanced_query"] = advancedQueryMap

		scheduleControlMap := make(map[string]interface{})
		scheduleControlMap["from_unixtime"] = aggRule.BeginUnixTime
		scheduleControlMap["to_unixtime"] = aggRule.EndUnixTime
		scheduleControlMap["granularity"] = aggRule.Interval
		scheduleControlMap["delay"] = aggRule.DelaySeconds
		aggRuleMap["schedule_control"] = scheduleControlMap

		aggRuleJsons = append(aggRuleJsons, aggRuleMap)
	}

	scheduledSql := make(map[string]interface{})
	scheduledSql["agg_rules"] = aggRuleJsons
	scheduledSqlJson, err := json.Marshal(scheduledSql)
	if err != nil {
		return nil, err
	}
	params["config.ml.scheduled_sql"] = string(scheduledSqlJson)

	return params, nil
}

func (c *Client) createMetricAggRulesConfig(aggRules *MetricAggRules) (*ETL, error) {
	etl := new(ETL)

	etl.Name = aggRules.ID
	etl.DisplayName = aggRules.Name
	etl.Description = aggRules.Desc
	etl.Type = "ETL"

	etl.Configuration.AccessKeyId = aggRules.SrcAccessKeyID
	etl.Configuration.AccessKeySecret = aggRules.SrcAccessKeySecret
	etl.Configuration.Script = ""
	etl.Configuration.Logstore = aggRules.SrcStore
	parameters, err := c.getScheduledSQLParams(aggRules.AggRules)
	if err != nil {
		return nil, err
	}
	etl.Configuration.Parameters = parameters
	etl.Configuration.FromTime = time.Now().Unix()

	var sink ETLSink
	sink.Endpoint = aggRules.DestEndpoint
	sink.Name = "sls-convert-metric"
	sink.AccessKeyId = aggRules.DestAccessKeyID
	sink.AccessKeySecret = aggRules.DestAccessKeySecret
	sink.Project = aggRules.DestProject
	sink.Logstore = aggRules.DestStore
	etl.Configuration.ETLSinks = append(etl.Configuration.ETLSinks, sink)

	etl.Schedule.Type = ScheduleTypeResident

	return etl, nil
}

func (c *Client) CreateMetricAggRules(project string, aggRules *MetricAggRules) error {
	etl, err := c.createMetricAggRulesConfig(aggRules)
	if err != nil {
		return err
	}
	if err := c.CreateETL(project, *etl); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateMetricAggRules(project string, aggRules *MetricAggRules) error {
	etl, err := c.createMetricAggRulesConfig(aggRules)
	if err != nil {
		return err
	}
	if err := c.UpdateETL(project, *etl); err != nil {
		return err
	}
	return nil
}

func (c *Client) castEtlToMetricAggRules(etl *ETL) (*MetricAggRules, error) {
	aggRules := new(MetricAggRules)
	aggRules.ID = etl.Name
	aggRules.Name = etl.DisplayName
	aggRules.Desc = etl.Description
	aggRules.SrcAccessKeyID = etl.Configuration.AccessKeyId
	aggRules.SrcAccessKeySecret = etl.Configuration.AccessKeySecret
	aggRules.SrcStore = etl.Configuration.Logstore

	scheduledSqlJson := etl.Configuration.Parameters["config.ml.scheduled_sql"]
	aggRuleJson := make(map[string][]map[string]interface{})
	err := json.Unmarshal([]byte(scheduledSqlJson), &aggRuleJson)
	if err != nil {
		return nil, err
	}
	aggRuleMaps := aggRuleJson["agg_rules"]

	var aggRuleItems []MetricAggRuleItem
	for _, aggRuleMap := range aggRuleMaps {
		aggRuleItem := new(MetricAggRuleItem)

		aggRuleItem.Name, err = castInterfaceToString(aggRuleMap, "rule_name")
		if err != nil {
			return nil, err
		}
		advancedQuery, err := castInterfaceToMap(aggRuleMap, "advanced_query")
		if err != nil {
			return nil, err
		}
		aggRuleItem.QueryType, err = castInterfaceToString(advancedQuery, "type")
		if err != nil {
			return nil, err
		}
		aggRuleItem.Query, err = castInterfaceToString(advancedQuery, "query")
		if err != nil {
			return nil, err
		}
		aggRuleItem.TimeName, err = castInterfaceToString(advancedQuery, "time_name")
		if err != nil {
			return nil, err
		}
		aggRuleItem.MetricNames, err = castInterfaceArrayToStringArray(advancedQuery, "metric_names")
		if err != nil {
			return nil, err
		}
		aggRuleItem.LabelNames, err = castInterfaceMapToStringMap(advancedQuery, "labels")
		if err != nil {
			return nil, err
		}
		scheduleControl, err := castInterfaceToMap(aggRuleMap, "schedule_control")
		if err != nil {
			return nil, err
		}
		aggRuleItem.BeginUnixTime, err = castInterfaceToInt(scheduleControl, "from_unixtime")
		if err != nil {
			return nil, err
		}
		aggRuleItem.EndUnixTime, err = castInterfaceToInt(scheduleControl, "to_unixtime")
		if err != nil {
			return nil, err
		}
		aggRuleItem.Interval, err = castInterfaceToInt(scheduleControl, "granularity")
		if err != nil {
			return nil, err
		}
		aggRuleItem.DelaySeconds, err = castInterfaceToInt(scheduleControl, "delay")
		if err != nil {
			return nil, err
		}

		aggRuleItems = append(aggRuleItems, *aggRuleItem)
	}
	aggRules.AggRules = aggRuleItems
	for _, sink := range etl.Configuration.ETLSinks {
		aggRules.DestEndpoint = sink.Endpoint
		aggRules.DestAccessKeyID = sink.AccessKeyId
		aggRules.DestAccessKeySecret = sink.AccessKeySecret
		aggRules.DestProject = sink.Project
		aggRules.DestStore = sink.Logstore
	}
	return aggRules, nil
}

func (c *Client) ListMetricAggRules(project string, offset int, size int) ([]*MetricAggRules, error) {
	listEtl, err := c.ListETL(project, offset, size)
	if err != nil {
		return nil, err
	}
	etls := listEtl.Results
	var aggRules []*MetricAggRules
	for _, etl := range etls {
		if _, ok := etl.Configuration.Parameters["config.ml.scheduled_sql"]; ok {
			aggRule, err := c.castEtlToMetricAggRules(etl)
			if err != nil {
				return nil, err
			}
			aggRules = append(aggRules, aggRule)
		}
	}
	return aggRules, nil
}

func (c *Client) GetMetricAggRules(project string, ruleID string) (*MetricAggRules, error) {
	etl, err := c.GetETL(project, ruleID)
	if err != nil {
		return nil, err
	}
	aggRules, err := c.castEtlToMetricAggRules(etl)
	if err != nil {
		return nil, err
	}
	return aggRules, nil

}

func (c *Client) DeleteMetricAggRules(project string, ruleID string) error {
	if err := c.DeleteETL(project, ruleID); err != nil {
		return err
	}
	return nil
}

func castInterfaceArrayToStringArray(inter map[string]interface{}, key string) ([]string, error) {
	t, ok := inter[key].([]interface{})
	if !ok {
		return nil, fmt.Errorf("castInterfaceArrayToStringArray is not ok, key: %s, value: %v\n", key, inter[key])
	}
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = v.(string)
	}
	return s, nil
}

func castInterfaceMapToStringMap(inter map[string]interface{}, key string) (map[string]string, error) {
	t, ok := inter[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("castInterfaceMapToStringMap is not ok, key: %s, value: %v\n", key, inter[key])
	}
	s := make(map[string]string, len(t))
	for k, v := range t {
		s[k] = v.(string)
	}
	return s, nil
}

func castInterfaceToInt(inter map[string]interface{}, key string) (int64, error) {
	t, ok := inter[key].(float64)
	if !ok {
		return 0, fmt.Errorf("castInterfaceToInt is not ok, key: %s, value: %v\n", key, inter[key])
	}
	return int64(t), nil
}

func castInterfaceToString(inter map[string]interface{}, key string) (string, error) {
	t, ok := inter[key].(string)
	if !ok {
		return "", fmt.Errorf("castInterfaceToString is not ok, key: %s, value: %v\n", key, inter[key])
	}
	return t, nil
}

func castInterfaceToMap(inter map[string]interface{}, key string) (map[string]interface{}, error) {
	t, ok := inter[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("castInterfaceToMap is not ok, key: %s, value: %v\n", key, inter[key])
	}
	return t, nil
}
