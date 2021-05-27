package models

import (
	"encoding/json"

	"github.com/IBM-Cloud/bluemix-go/crn"
)

type Service struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CatalogCRN string `json:"catalog_crn"`
	URL        string `json:"url"`
	Kind       string `json:"kind"`

	Metadata ServiceMetadata `json:"-"`
	Children []Service       `json:"children"`
	Active   bool            `json:"active"`
}

type ServicePlan struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CatalogCRN string `json:"catalog_crn"`
	URL        string `json:"url"`
	Kind       string `json:"kind"`
}

type ServiceDeployment struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	CatalogCRN string             `json:"catalog_crn"`
	Metadata   DeploymentMetaData `json:"metadata,omitempty"`
}

type ServiceDeploymentAlias struct {
	Metadata DeploymentMetaData `json:"metadata,omitempty"`
}

type DeploymentMetaData struct {
	RCCompatible  bool                       `json:"rc_compatible"`
	IAMCompatible bool                       `json:"iam_compatible"`
	Deployment    MetadataDeploymentFragment `json:"deployment,omitempty"`
	Service       MetadataServiceFragment    `json:"service,omitempty"`
}

type MetadataDeploymentFragment struct {
	DeploymentID string  `json:"deployment_id,omitempty"`
	TargetCrn    crn.CRN `json:"target_crn"`
	Location     string  `json:"location"`
}

type ServiceMetadata interface{}

type ServiceResourceMetadata struct {
	Service MetadataServiceFragment `json:"service"`
}

type MetadataServiceFragment struct {
	Bindable            bool   `json:"bindable"`
	IAMCompatible       bool   `json:"iam_compatible"`
	RCProvisionable     bool   `json:"rc_provisionable"`
	PlanUpdateable      bool   `json:"plan_updateable"`
	ServiceCheckEnabled bool   `json:"service_check_enabled"`
	ServiceKeySupported bool   `json:"service_key_supported"`
	State               string `json:"state"`
	TestCheckInterval   int    `json:"test_check_interval"`
	UniqueAPIKey        bool   `json:"unique_api_key"`

	// CF properties
	ServiceBrokerGUID string `json:"service_broker_guid"`
}

type PlatformServiceResourceMetadata struct {
}

type TemplateResourceMetadata struct {
}

type RuntimeResourceMetadata struct {
}

// UnmarshalJSON provide custom JSON unmarshal behavior to support multiple types
// of `metadata`
func (s *Service) UnmarshalJSON(data []byte) error {
	type Copy Service

	trial := &struct {
		*Copy
		Metadata json.RawMessage `json:"metadata"`
	}{
		Copy: (*Copy)(s),
	}

	if err := json.Unmarshal(data, trial); err != nil {
		return err
	}

	if len(trial.Metadata) == 0 {
		s.Metadata = nil
		return nil
	}

	switch s.Kind {
	case "runtime":
		s.Metadata = &RuntimeResourceMetadata{}
	case "service", "iaas":
		s.Metadata = &ServiceResourceMetadata{}
	case "platform_service":
		s.Metadata = &PlatformServiceResourceMetadata{}
	case "template":
		s.Metadata = &TemplateResourceMetadata{}
	default:
		s.Metadata = nil
		return nil
	}

	if err := json.Unmarshal(trial.Metadata, s.Metadata); err != nil {
		return err
	}
	return nil
}