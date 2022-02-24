package models

import (
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
)

type MetadataType struct {
	ID        string     `json:"id"`
	Guid      string     `json:"guid"`
	Url       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type ServiceInstance struct {
	*MetadataType
	Name                string `json:"name"`
	RegionID            string `json:"region_id"`
	AccountID           string `json:"account_id"`
	ServicePlanID       string `json:"resource_plan_id"`
	ServicePlanName     string
	ResourceGroupID     string `json:"resource_group_id"`
	ResourceGroupName   string
	Crn                 crn.CRN                `json:"crn,omitempty"`
	Tags                []string               `json:"tags,omitempty"`
	Parameters          map[string]interface{} `json:"parameters,omitempty"`
	Extensions          map[string]interface{} `json:"extensions,omitempty"`
	CreateTime          int64                  `json:"create_time"`
	State               string                 `json:"state"`
	Type                string                 `json:"type"`
	ServiceID           string                 `json:"resource_id"`
	ServiceName         string
	DashboardUrl        *string            `json:"dashboard_url"`
	LastOperation       *LastOperationType `json:"last_operation"`
	AccountUrl          string             `json:"account_url"`
	ResourcePlanUrl     string             `json:"resource_plan_url"`
	ResourceBindingsUrl string             `json:"resource_bindings_url"`
	ResourceAliasesUrl  string             `json:"resource_aliases_url"`
	SiblingsUrl         string             `json:"siblings_url"`
	TargetCrn           crn.CRN            `json:"target_crn"`
}

type LastOperationType struct {
	Type        string     `json:"type"`
	State       string     `json:"state"`
	Description *string    `json:"description"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ServiceInstanceV2 struct {
	ServiceInstance
	ScheduledReclaimAt interface{}       `json:"scheduled_reclaim_at"`
	RestoredAt         interface{}       `json:"restored_at"`
	ScheduledReclaimBy string            `json:"scheduled_reclaim_by"`
	RestoredBy         string            `json:"restored_by"`
	ResourcePlanID     string            `json:"resource_plan_id"`
	ResourceGroupCrn   string            `json:"resource_group_crn"`
	AllowCleanup       bool              `json:"allow_cleanup"`
	ResourceKeysURL    string            `json:"resource_keys_url"`
	PlanHistory        []PlanHistoryData `json:"plan_history"`
}

type PlanHistoryData struct {
	ResourcePlanID string    `json:"resource_plan_id"`
	StartDate      time.Time `json:"start_date"`
	RequestorID    string    `json:"requestor_id"`
	Migrated       bool      `json:"migrated"`
	ControlledBy   string    `json:"controlled_by"`
	Locked         bool      `json:"locked"`
}
