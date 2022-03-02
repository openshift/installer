package icdv4

import (
	"encoding/json"
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

// AutoscalingSetGroup ...
type AutoscalingSetGroup struct {
	Autoscaling AutoscalingGroup `json:"autoscaling,omitempty"`
}

// AutoscalingGroup ...
type AutoscalingGroup struct {
	Memory *ASGBody `json:"memory,omitempty"`
	CPU    *ASGBody `json:"cpu,omitempty"`
	Disk   *ASGBody `json:"disk,omitempty"`
}

// ASGBody ...
type ASGBody struct {
	Scalers ScalersBody `json:"scalers,omitempty"`
	Rate    RateBody    `json:"rate,omitempty"`
}

// RateBody ...
type RateBody struct {
	IncreasePercent     int    `json:"increase_percent,omitempty"`
	PeriodSeconds       int    `json:"period_seconds,omitempty"`
	LimitCountPerMember int    `json:"limit_count_per_member,omitempty"`
	LimitMBPerMember    int    `json:"limit_mb_per_member,omitempty"`
	Units               string `json:"units,omitempty"`
}

// ScalersBody ...
type ScalersBody struct {
	Capacity *CapacityBody `json:"capacity,omitempty"`
	IO       *IOBody       `json:"io_utilization,omitempty"`
}

// CapacityBody ...
type CapacityBody struct {
	Enabled                   bool `json:"enabled"`
	FreeSpaceRemainingPercent int  `json:"free_space_remaining_percent,omitempty"`
	FreeSpaceLessThanPercent  int  `json:"free_space_less_than_percent,omitempty"`
}

// IOBody ...
type IOBody struct {
	Enabled      bool   `json:"enabled"`
	AbovePercent int    `json:"above_percent,omitempty"`
	OverPeriod   string `json:"over_period,omitempty"`
}

// AutoscalingGetGroup ...
type AutoscalingGetGroup struct {
	Autoscaling AutoscalingGet `json:"autoscaling,omitempty"`
}

// AutoscalingGet ...
type AutoscalingGet struct {
	Memory ASGGet `json:"memory,omitempty"`
	CPU    ASGGet `json:"cpu,omitempty"`
	Disk   ASGGet `json:"disk,omitempty"`
}

// ASGGet ...
type ASGGet struct {
	Scalers ScalersBody `json:"scalers,omitempty"`
	Rate    Rate        `json:"rate,omitempty"`
}

// Rate ...
type Rate struct {
	IncreasePercent     json.Number `json:"increase_percent,omitempty"`
	PeriodSeconds       int         `json:"period_seconds,omitempty"`
	LimitCountPerMember int         `json:"limit_count_per_member,omitempty"`
	LimitMBPerMember    json.Number `json:"limit_mb_per_member,omitempty"`
	Units               string      `json:"units,omitempty"`
}

type autoScaling struct {
	client *client.Client
}

// AutoScaling ...
type AutoScaling interface {
	GetAutoScaling(icdID string, groupID string) (AutoscalingGetGroup, error)
	SetAutoScaling(icdID string, groupID string, AutoScaleReq AutoscalingSetGroup) (Task, error)
}

func newAutoScalingAPI(c *client.Client) AutoScaling {
	return &autoScaling{
		client: c,
	}
}
func (r *autoScaling) GetAutoScaling(icdID string, groupID string) (AutoscalingGetGroup, error) {
	autoscalingGroupResult := AutoscalingGetGroup{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/groups/%s/autoscaling", utils.EscapeUrlParm(icdID), groupID)
	_, err := r.client.Get(rawURL, &autoscalingGroupResult)
	if err != nil {
		return autoscalingGroupResult, err
	}
	return autoscalingGroupResult, nil
}
func (r *autoScaling) SetAutoScaling(icdID string, groupID string, AutoScaleReq AutoscalingSetGroup) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/groups/%s/autoscaling", utils.EscapeUrlParm(icdID), groupID)
	_, err := r.client.Patch(rawURL, &AutoScaleReq, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}
