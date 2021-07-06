// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisEdgeFunctionsActionActionName = "action_name"
	cisEdgeFunctionsActionScript     = "script"
)

func resourceIBMCISEdgeFunctionsAction() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISEdgeFunctionsActionCreate,
		Read:     resourceIBMCISEdgeFunctionsActionRead,
		Update:   resourceIBMCISEdgeFunctionsActionUpdate,
		Delete:   resourceIBMCISEdgeFunctionsActionDelete,
		Exists:   resourceIBMCISEdgeFunctionsActionExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
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

func resourceIBMCISEdgeFunctionsActionCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	scriptName := d.Get(cisEdgeFunctionsActionActionName).(string)
	script := d.Get(cisEdgeFunctionsActionScript).(string)
	r := ioutil.NopCloser(strings.NewReader(script))
	opt := cisClient.NewUpdateEdgeFunctionsActionOptions(scriptName)
	opt.SetEdgeFunctionsAction(r)

	_, _, err = cisClient.UpdateEdgeFunctionsAction(opt)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	d.SetId(convertCisToTfThreeVar(scriptName, zoneID, crn))
	return resourceIBMCISEdgeFunctionsActionRead(d, meta)
}

func resourceIBMCISEdgeFunctionsActionUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange(cisEdgeFunctionsActionScript) {
		return resourceIBMCISEdgeFunctionsActionCreate(d, meta)
	}

	return resourceIBMCISEdgeFunctionsActionRead(d, meta)
}

func resourceIBMCISEdgeFunctionsActionRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	scriptName, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsActionOptions(scriptName)
	result, resp, err := cisClient.GetEdgeFunctionsAction(opt)
	if err != nil {
		return fmt.Errorf("Error: %v", resp)
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
		return fmt.Errorf("Error in closing reader")
	}

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisEdgeFunctionsActionActionName, scriptName)
	d.Set(cisEdgeFunctionsActionScript, string(content))
	return nil
}

func resourceIBMCISEdgeFunctionsActionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return false, fmt.Errorf("Error in creating CIS object")
	}

	scriptName, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsActionOptions(scriptName)
	_, response, err := cisClient.GetEdgeFunctionsAction(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Edge functions action script is not found")
			return false, nil
		}
		return false, fmt.Errorf("Error: %v", response)
	}
	return true, nil
}

func resourceIBMCISEdgeFunctionsActionDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return fmt.Errorf("Error in creating CIS object")
	}

	scriptName, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewDeleteEdgeFunctionsActionOptions(scriptName)
	_, response, err := cisClient.DeleteEdgeFunctionsAction(opt)
	if err != nil {
		return fmt.Errorf("Error in edge function action script deletion: %v", response)
	}
	return nil
}
