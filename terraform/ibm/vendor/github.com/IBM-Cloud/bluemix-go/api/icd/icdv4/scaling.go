package icdv4

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

type GroupList struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Id      string  `json:"id"`
	Count   int     `json:"count"`
	Members Members `json:"members"`
	Memory  Memory `json:"memory"`
	Cpu     Cpu    `json:"cpu"`
	Disk    Disk   `json:"disk"`
}

type Members struct {
	Units           string `json:"units"`
	AllocationCount int    `json:"allocation_count"`
	MinimumCount    int    `json:"minimum_count"`
	MaximumCount    int    `json:"maximum_count"`
	StepSizeCount   int    `json:"step_size_count"`
	IsAdjustable    bool   `json:"is_adjustable"`
	CanScaleDown    bool   `json:"can_scale_down"`
}

type Memory struct {
	Units        string `json:"units"`
	AllocationMb int    `json:"allocation_mb"`
	MinimumMb    int    `json:"minimum_mb"`
	MaximumMb    int    `json:"maximum_mb"`
	StepSizeMb   int    `json:"step_size_mb"`
	IsAdjustable bool   `json:"is_adjustable"`
	CanScaleDown bool   `json:"can_scale_down"`
}

type Cpu struct {
	Units           string `json:"units"`
	AllocationCount int    `json:"allocation_count"`
	MinimumCount    int    `json:"minimum_count"`
	MaximumCount    int    `json:"maximum_count"`
	StepSizeCount   int    `json:"step_size_count"`
	IsAdjustable    bool   `json:"is_adjustable"`
	CanScaleDown    bool   `json:"can_scale_down"`
}

type Disk struct {
	Units        string `json:"units"`
	AllocationMb int    `json:"allocation_mb"`
	MinimumMb    int    `json:"minimum_mb"`
	MaximumMb    int    `json:"maximum_mb"`
	StepSizeMb   int    `json:"step_size_mb"`
	IsAdjustable bool   `json:"is_adjustable"`
	CanScaleDown bool   `json:"can_scale_down"`
}

type GroupReq struct {
	GroupBdy GroupBdy `json:"group"`
}

type GroupBdy struct {
	Members *MembersReq `json:"members,omitempty"`
	Memory  *MemoryReq  `json:"memory,omitempty"`
	Cpu     *CpuReq     `json:"cpu,omitempty"`
	Disk    *DiskReq    `json:"disk,omitempty"`
}

type MembersReq struct {
	AllocationCount int `json:"allocation_count,omitempty"`
}
type MemoryReq struct {
	AllocationMb int `json:"allocation_mb,omitempty"`
}
type CpuReq struct {
	AllocationCount int `json:"allocation_count,omitempty"`
}
type DiskReq struct {
	AllocationMb int `json:"allocation_mb,omitempty"`
}

type Groups interface {
	GetDefaultGroups(groupType string) (GroupList, error)
	GetGroups(icdId string) (GroupList, error)
	UpdateGroup(icdId string, groupId string, groupReq GroupReq) (Task, error)
}

type groups struct {
	client *client.Client
}

func newGroupAPI(c *client.Client) Groups {
	return &groups{
		client: c,
	}
}

func (r *groups) GetDefaultGroups(groupType string) (GroupList, error) {
	groupList := GroupList{}
	rawURL := fmt.Sprintf("/v4/ibm/deployables/%s/groups", groupType)
	_, err := r.client.Get(rawURL, &groupList)
	if err != nil {
		return groupList, err
	}
	return groupList, nil
}

func (r *groups) GetGroups(icdId string) (GroupList, error) {
	groupList := GroupList{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/groups", utils.EscapeUrlParm(icdId))
	_, err := r.client.Get(rawURL, &groupList)
	if err != nil {
		return groupList, err
	}
	return groupList, nil
}

func (r *groups) UpdateGroup(icdId string, groupId string, groupReq GroupReq) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/groups/%s", utils.EscapeUrlParm(icdId), groupId)
	_, err := r.client.Patch(rawURL, &groupReq, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}


