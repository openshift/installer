// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSchematicsOutput() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSchematicsOutputRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace for which you want to retrieve output values. To find the workspace ID, use the `GET /workspaces` API.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Region of the workspace.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of template",
			},
			"output_values": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"output_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The json output in string",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this Workspace",
			},
		},
	}
}

func dataSourceIBMSchematicsOutputRead(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	workspaceID := d.Get("workspace_id").(string)
	templateID := d.Get("template_id").(string)

	if r, ok := d.GetOk("location"); ok {
		region := r.(string)
		schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
		if updatedURL {
			schematicsClient.Service.Options.URL = schematicsURL
		}
	}

	getWorkspaceOutputsOptions := &schematicsv1.GetWorkspaceOutputsOptions{}

	getWorkspaceOutputsOptions.SetWID(d.Get("workspace_id").(string))

	outputValuesList, response, err := schematicsClient.GetWorkspaceOutputs(getWorkspaceOutputsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetWorkspaceOutputs failed %s\n%s", err, response)
		return err
	}

	var outputJSON string
	items := make(map[string]interface{})
	found := false
	for _, fields := range outputValuesList {
		if *fields.ID == templateID {
			output := fields.OutputValues
			found = true
			outputByte, err := json.MarshalIndent(output, "", "")
			if err != nil {
				return err
			}
			outputJSON = string(outputByte[:])
			// items := map[string]interface{}
			for _, value := range output {
				for key, val := range value {
					val2 := val.(map[string]interface{})["value"]
					items[key] = val2
				}
			}
		}
	}

	if !(found) {
		return fmt.Errorf("[ERROR] Error while fetching template id in workspace: %s", workspaceID)
	}
	d.Set("output_json", outputJSON)
	d.SetId(fmt.Sprintf("%s/%s", workspaceID, templateID))
	d.Set("output_values", flex.Flatten(items))

	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}

	d.Set(flex.ResourceControllerURL, controller+"/schematics")

	return nil
}

// dataSourceIBMSchematicsOutputID returns a reasonable ID for the list.
func dataSourceIBMSchematicsOutputID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceOutputValuesListFlattenOutputValues(result []schematicsv1.OutputValuesInner) (outputValues interface{}) {
	for _, outputValuesItem := range result {
		outputValues = dataSourceOutputValuesListOutputValuesToMap(outputValuesItem)
	}

	return outputValues
}

func dataSourceOutputValuesListOutputValuesToMap(outputValuesItem schematicsv1.OutputValuesInner) (outputValuesMap map[string]interface{}) {
	outputValuesMap = map[string]interface{}{}

	if outputValuesItem.Folder != nil {
		outputValuesMap["folder"] = outputValuesItem.Folder
	}
	if outputValuesItem.ID != nil {
		outputValuesMap["id"] = outputValuesItem.ID
	}

	m := []flex.Map{}

	for _, outputValues := range outputValuesItem.OutputValues {
		m = append(m, flex.Flatten(outputValues))
	}

	if outputValuesItem.OutputValues != nil {
		outputValuesMap["output_values"] = m
	}
	if outputValuesItem.ValueType != nil {
		outputValuesMap["value_type"] = outputValuesItem.ValueType
	}

	return outputValuesMap
}
