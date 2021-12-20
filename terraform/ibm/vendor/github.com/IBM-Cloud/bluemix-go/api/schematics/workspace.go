package schematics

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type WorkspaceConfig struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Type               []string           `json:"type"`
	Description        string             `json:"description"`
	ResourceGroup      string             `json:"resource_group"`
	Location           string             `json:"location"`
	Tags               []string           `json:"tags"`
	CreatedAt          string             `json:"created_at"`
	CreatedBy          string             `json:"created_by"`
	Status             string             `json:"status"`
	WorkspaceStatusMsg StatusMsgInfo      `json:"workspace_status_msg"`
	WorkspaceStatus    StatusInfo         `json:"workspace_status"`
	TemplateRepo       RepoInfo           `json:"template_repo"`
	TemplateData       []TemplateDataInfo `json:"template_data"`
	RuntimeData        []RuntimeDataInfo  `json:"runtime_data"`
	SharedData         SharedDataInfo     `json:"shared_data"`
	UpdatedAt          string             `json:"updated_at"`
	LastHealthCheckAt  string             `json:"last_health_check_at"`
	CatalogRef         CatalogInfo        `json:"catalog_ref"`
	CRN                string             `json:"crn"`
}

type StatusMsgInfo struct {
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
type StatusInfo struct {
	Frozen     bool   `json:"frozen"`
	FrozenAt   string `json:"status_msg"`
	LockedBy   string `json:"locked_by"`
	FrozenBy   string `json:"frozen_by"`
	Locked     bool   `json:"locked"`
	LockedTime string `json:"locked_time"`
}

type TemplateDataInfo struct {
	Env                 []EnvValues     `json:"env_values"`
	Folder              string          `json:"folder"`
	TemplateID          string          `json:"id"`
	Type                string          `json:"type"`
	Locked              bool            `json:"locked"`
	UninstallScriptName string          `json:"uninstall_script_name"`
	Values              string          `json:"values"`
	ValuesMetadata      interface{}     `json:"values_metadata"`
	ValuesURL           string          `json:"values_url"`
	Variablestore       []Variablestore `json:"variablestore"`
}

type RuntimeDataInfo struct {
	EngineCmd     string                `json:"engine_cmd"`
	EngineName    string                `json:"engine_name"`
	TemplateID    string                `json:"id"`
	EngineVersion string                `json:"engine_version"`
	LogStoreURL   string                `json:"log_store_url"`
	OutputValues  map[string]string     `json:"output_values"`
	StateStoreURL string                `json:"state_store_url"`
	Resources     [][]map[string]string `json:"resources"`
}
type RepoInfo struct {
	Branch  string `json:"branch"`
	Release string `json:"release"`
	RepoURL string `json:"repo_url"`
	URL     string `json:"url"`
}
type SharedDataInfo struct {
	ClusterID       string              `json:"cluster_id"`
	ClusterName     string              `json:"cluster_name"`
	EntitlementKeys []map[string]string `json:"entitlement_keys"`
	Namespace       string              `json:"namespace"`
	Region          string              `json:"region"`
	ResourceGroupID string              `json:"resource_group_id"`
}

type EnvValues struct {
	Hidden bool   `json:"hidden"`
	Name   string `json:"name"`
	// Secure bool   `json:"secure"`
	Value string `json:"value"`
}
type OutputResponse struct {
	Folder     string                    `json:"folder"`
	TemplateID string                    `json:"id"`
	Type       string                    `json:"type"`
	Output     []map[string]OutputValues `json:"output_values"`
}

type OutputValues struct {
	Sensitive bool        `json:"sensitive"`
	Value     interface{} `json:"value"`
	Type      interface{} `json:"type"`
}

type CreateWorkspaceConfig struct {
	Name            string             `json:"name"`
	Type            []string           `json:"type"`
	Description     string             `json:"description"`
	Tags            []string           `json:"tags"`
	WorkspaceStatus StatusInfo         `json:"workspace_status"`
	TemplateRepo    RepoInfo           `json:"template_repo"`
	TemplateRef     string             `json:"template_ref"`
	TemplateData    []TemplateDataInfo `json:"template_data"`
}

type Payload struct {
	Name            string          `json:"name"`
	Type            []string        `json:"type"`
	Description     string          `json:"description"`
	Tags            []string        `json:"tags"`
	TemplateRef     string          `json:"template_ref"`
	TemplateRepo    TemplateRepo    `json:"template_repo"`
	WorkspaceStatus WorkspaceStatus `json:"workspace_status"`
	TemplateData    []TemplateData  `json:"template_data"`
}
type TemplateRepo struct {
	URL string `json:"url"`
}
type WorkspaceStatus struct {
	Frozen bool `json:"frozen"`
}
type Variablestore struct {
	Name        string `json:"name"`
	Secure      bool   `json:"secure,omitempty"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}
type TemplateData struct {
	Folder        string          `json:"folder"`
	Type          string          `json:"type"`
	Variablestore []Variablestore `json:"variablestore"`
}

type CatalogInfo struct {
	ItemID          string `json:"item_id"`
	ItemName        string `json:"item_name"`
	ItemURL         string `json:"item_url"`
	ItemReadmeURL   string `json:"item_readme_url"`
	ItemIconURL     string `json:"item_icon_url"`
	OfferingVersion string `json:"offering_version"`
}
type workspace struct {
	client *client.Client
}

type Workspaces interface {
	GetWorkspaceByID(WorskpaceID string) (WorkspaceConfig, error)
	GetOutputValues(WorskpaceID string) ([]OutputResponse, error)
	GetStateStore(WorskpaceID, TemplateID string) (interface{}, error)
	CreateWorkspace(createReq Payload) (WorkspaceConfig, error)
}

func newWorkspaceAPI(c *client.Client) Workspaces {
	return &workspace{
		client: c,
	}
}

func (r *workspace) GetWorkspaceByID(WorskpaceID string) (WorkspaceConfig, error) {
	var successV WorkspaceConfig
	_, err := r.client.Get(fmt.Sprintf("/v1/workspaces/%s", WorskpaceID), &successV)
	return successV, err
}
func (r *workspace) GetStateStore(WorskpaceID, TemplateID string) (interface{}, error) {
	var successV interface{}
	_, err := r.client.Get(fmt.Sprintf("/v1/workspaces/%s/runtime_data/%s/state_store", WorskpaceID, TemplateID), &successV)
	return successV, err
}
func (r *workspace) GetOutputValues(WorskpaceID string) ([]OutputResponse, error) {
	outputs := []OutputResponse{}
	_, err := r.client.Get(fmt.Sprintf("/v1/workspaces/%s/output_values", WorskpaceID), &outputs)
	if err != nil {
		return nil, err
	}
	return outputs, err
}
func (r *workspace) CreateWorkspace(createReq Payload) (WorkspaceConfig, error) {
	var successV WorkspaceConfig
	_, err := r.client.Post("/v1/workspaces", createReq, &successV)
	return successV, err
}
