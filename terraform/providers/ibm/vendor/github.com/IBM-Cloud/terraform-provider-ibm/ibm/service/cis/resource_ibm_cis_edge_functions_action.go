// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisEdgeFunctionsActionActionName = "action_name"
	cisEdgeFunctionsActionScript     = "script"
)

func ResourceIBMCISEdgeFunctionsAction() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISEdgeFunctionsActionCreate,
		Read:     ResourceIBMCISEdgeFunctionsActionRead,
		Update:   ResourceIBMCISEdgeFunctionsActionUpdate,
		Delete:   ResourceIBMCISEdgeFunctionsActionDelete,
		Exists:   ResourceIBMCISEdgeFunctionsActionExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeValidator("ibm_cis_edge_functions_action",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisEdgeFunctionsActionActionName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge function action script name",
			},
			cisEdgeFunctionsActionScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge function action script",
			},
		},
	}
}

func ResourceIBMCISEdgeFunctionsActionValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISEdgeFunctionsActionValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_edge_functions_action",
		Schema:       validateSchema}
	return &ibmCISEdgeFunctionsActionValidator
}

func ResourceIBMCISEdgeFunctionsActionCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	scriptName := d.Get(cisEdgeFunctionsActionActionName).(string)
	script := d.Get(cisEdgeFunctionsActionScript).(string)
	r := ioutil.NopCloser(strings.NewReader(script))
	opt := cisClient.NewUpdateEdgeFunctionsActionOptions(scriptName)
	opt.SetEdgeFunctionsAction(r)

	_, _, err = cisClient.UpdateEdgeFunctionsAction(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error: %v", err)
	}
	d.SetId(flex.ConvertCisToTfThreeVar(scriptName, zoneID, crn))
	return ResourceIBMCISEdgeFunctionsActionRead(d, meta)
}

func ResourceIBMCISEdgeFunctionsActionUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange(cisEdgeFunctionsActionScript) {
		return ResourceIBMCISEdgeFunctionsActionCreate(d, meta)
	}

	return ResourceIBMCISEdgeFunctionsActionRead(d, meta)
}

func ResourceIBMCISEdgeFunctionsActionRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	scriptName, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsActionOptions(scriptName)
	result, resp, err := cisClient.GetEdgeFunctionsAction(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error: %v", resp)
	}

	// read script content
	content := []byte{}
	p := make([]byte, 8)
	for {
		n, err := result.Read(p)
		content = append(content, p[:n]...)
		if err == io.EOF || n < 1 {
			break
		}
	}
	err = result.Close()
	if err != nil {
		return fmt.Errorf("[ERROR] Error in closing reader")
	}

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisEdgeFunctionsActionActionName, scriptName)
	d.Set(cisEdgeFunctionsActionScript, string(content))
	return nil
}

func ResourceIBMCISEdgeFunctionsActionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error in creating CIS object")
	}

	scriptName, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsActionOptions(scriptName)
	_, response, err := cisClient.GetEdgeFunctionsAction(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Edge functions action script is not found")
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error: %v", response)
	}
	return true, nil
}

func ResourceIBMCISEdgeFunctionsActionDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error in creating CIS object")
	}

	scriptName, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewDeleteEdgeFunctionsActionOptions(scriptName)
	_, response, err := cisClient.DeleteEdgeFunctionsAction(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in edge function action script deletion: %v", response)
	}
	return nil
}
