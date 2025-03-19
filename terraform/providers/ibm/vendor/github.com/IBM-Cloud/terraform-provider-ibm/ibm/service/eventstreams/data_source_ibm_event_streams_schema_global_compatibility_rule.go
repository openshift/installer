// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventstreams

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/eventstreams-go-sdk/pkg/schemaregistryv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	schemaGlobalCompatibilityRuleResourceType = "schema-global-compatibility-rule"
	schemaCompatibilityRuleType               = "COMPATIBILITY"
	schemaCompatibilityDefaultValue           = "NONE"
)

var (
	schemaCompatiblityRuleValidConfigValues = []string{schemaCompatibilityDefaultValue, "FULL", "FULL_TRANSITIVE", "FORWARD", "FORWARD_TRANSITIVE", "BACKWARD", "BACKWARD_TRANSITIVE"}
)

// The global compatibility rule in an Event Streams service instance.
// The ID is the CRN with the last two components "schema-global-compatibility:".
func DataSourceIBMEventStreamsSchemaGlobalCompatibilityRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEventStreamsSchemaGlobalCompatibilityRuleRead,

		Schema: map[string]*schema.Schema{
			"resource_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID or CRN of the Event Streams service instance",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the global schema compatibility rule",
			},
		},
	}
}

// read global compatibility rule properties using the schema registry API
func dataSourceIBMEventStreamsSchemaGlobalCompatibilityRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaregistryClient, err := meta.(conns.ClientSession).ESschemaRegistrySession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "Error getting Event Streams schema registry session", "ibm_event_streams_schema_global_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	adminURL, instanceCRN, err := getSchemaRuleInstanceURL(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "Error getting Event Streams schema registry URL", "ibm_event_streams_schema_global_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	schemaregistryClient.SetServiceURL(adminURL)

	getOpts := &schemaregistryv1.GetGlobalRuleOptions{}
	getOpts.SetRule(schemaCompatibilityRuleType)
	rule, _, err := schemaregistryClient.GetGlobalRuleWithContext(context, getOpts)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "GetGlobalRule returned error", "ibm_event_streams_schema_global_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if rule.Config == nil {
		tfErr := flex.TerraformErrorf(err, "Unexpected nil config when getting global compatibility rule", "ibm_event_streams_schema_global_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId(getSchemaGlobalCompatibilityRuleID(instanceCRN))
	d.Set("resource_instance_id", instanceCRN)
	d.Set("config", *rule.Config)
	return nil
}

func getSchemaRuleInstanceURL(d *schema.ResourceData, meta interface{}) (string, string, error) {
	instanceCRN := d.Get("resource_instance_id").(string)
	if instanceCRN == "" { // importing
		id := d.Id()
		crnSegments := strings.Split(id, ":")
		if len(crnSegments) != 10 || crnSegments[8] != schemaGlobalCompatibilityRuleResourceType {
			return "", "", fmt.Errorf("ID '%s' is not a schema global compatibility resource", id)
		}
		crnSegments[8] = ""
		crnSegments[9] = ""
		instanceCRN = strings.Join(crnSegments, ":")
		d.Set("resource_instance_id", instanceCRN)
	}

	instance, err := getInstanceDetails(instanceCRN, meta)
	if err != nil {
		return "", "", err
	}
	adminURL := instance.Extensions["kafka_http_url"].(string)
	planID := *instance.ResourcePlanID
	valid := strings.Contains(planID, "enterprise")
	if !valid {
		return "", "", fmt.Errorf("schema registry is not supported by the Event Streams %s plan, enterprise plan is expected",
			planID)
	}
	return adminURL, instanceCRN, nil
}

func getSchemaGlobalCompatibilityRuleID(instanceCRN string) string {
	crnSegments := strings.Split(instanceCRN, ":")
	crnSegments[8] = schemaGlobalCompatibilityRuleResourceType
	crnSegments[9] = ""
	return strings.Join(crnSegments, ":")
}
