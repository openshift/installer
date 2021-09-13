// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/networking-go-sdk/filtersv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISFilters        = "ibm_cis_filter"
	cisFilterExpression  = "expression"
	cisFilterPaused      = "paused"
	cisFilterDescription = "description"
	cisFilterID          = "filter_id"
)

func resourceIBMCISFilter() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISFilterCreate,
		Read:     resourceIBMCISFilterRead,
		Update:   resourceIBMCISFilterUpdate,
		Delete:   resourceIBMCISFilterDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisFilterPaused: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter Paused",
			},
			cisFilterID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filter ID",
			},
			cisFilterExpression: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Filter Expression",
			},
			cisFilterDescription: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter Description",
				ValidateFunc: InvokeValidator(ibmCISFilters, cisFilterDescription),
			},
		},
	}
}
func resourceIBMCISFilterCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return fmt.Errorf("Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(ClientSession).CisFiltersSession()
	if err != nil {
		return fmt.Errorf("Error while getting the CisFiltersSession %s", err)
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))

	var newfilter filtersv1.FilterInput

	if p, ok := d.GetOkExists(cisFilterPaused); ok {
		paused := p.(bool)
		newfilter.Paused = &paused
	}
	if des, ok := d.GetOk(cisFilterDescription); ok {
		description := des.(string)
		newfilter.Description = &description
	}
	if e, ok := d.GetOk(cisFilterExpression); ok {
		expression := e.(string)
		newfilter.Expression = &expression
	}

	opt := cisClient.NewCreateFilterOptions(xAuthtoken, crn, zoneID)

	opt.SetFilterInput([]filtersv1.FilterInput{newfilter})

	result, resp, err := cisClient.CreateFilter(opt)
	if err != nil || result == nil {
		return fmt.Errorf("Error creating Filter for zone %q: %s %s", zoneID, err, resp)
	}
	d.SetId(convertCisToTfThreeVar(*result.Result[0].ID, zoneID, crn))
	return resourceIBMCISFilterRead(d, meta)

}
func resourceIBMCISFilterRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return fmt.Errorf("Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(ClientSession).CisFiltersSession()
	if err != nil {
		return fmt.Errorf("Error while getting the CisFiltersSession %s", err)
	}
	filterid, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	opt := cisClient.NewGetFilterOptions(xAuthtoken, crn, zoneID, filterid)

	result, response, err := cisClient.GetFilter(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Error GetFilter not found ")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error finding GetFilter %q: %s %s", d.Id(), err, response)
	}
	if result.Result != nil {
		d.Set(cisID, crn)
		d.Set(cisDomainID, zoneID)
		d.Set(cisFilterID, result.Result.ID)
		d.Set(cisFilterPaused, result.Result.Paused)
		d.Set(cisFilterDescription, result.Result.Description)
		d.Set(cisFilterExpression, result.Result.Expression)
	}
	return nil
}
func resourceIBMCISFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return fmt.Errorf("Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(ClientSession).CisFiltersSession()
	if err != nil {
		return fmt.Errorf("Error while getting the CisFiltersSession %s", err)
	}

	filterid, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange(cisFilterExpression) ||
		d.HasChange(cisFilterPaused) ||
		d.HasChange(cisFilterDescription) {

		var updatefilter filtersv1.FilterUpdateInput
		updatefilter.ID = &filterid

		if p, ok := d.GetOkExists(cisFilterPaused); ok {
			paused := p.(bool)
			updatefilter.Paused = &paused
		}
		if des, ok := d.GetOk(cisFilterDescription); ok {
			description := des.(string)
			updatefilter.Description = &description
		}
		if e, ok := d.GetOk(cisFilterExpression); ok {
			expression := e.(string)
			updatefilter.Expression = &expression
		}

		opt := cisClient.NewUpdateFiltersOptions(xAuthtoken, crn, zoneID)

		opt.SetFilterUpdateInput([]filtersv1.FilterUpdateInput{updatefilter})

		result, resp, err := cisClient.UpdateFilters(opt)
		if err != nil {
			return fmt.Errorf("Error updating Filter for zone %q: %s %s", zoneID, err, resp)
		}

		if *result.Result[0].ID == "" {
			return fmt.Errorf("Error failed to find id in Update response; resource was empty")
		}
	}
	return resourceIBMCISFilterRead(d, meta)
}
func resourceIBMCISFilterDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	xAuthtoken := sess.Config.IAMAccessToken
	cisClient, err := meta.(ClientSession).CisFiltersSession()
	if err != nil {
		return err
	}
	filterid, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	opt := cisClient.NewDeleteFiltersOptions(xAuthtoken, crn, zoneID, filterid)
	_, _, err = cisClient.DeleteFilters(opt)
	if err != nil {
		return fmt.Errorf("Error deleting Filter: %s", err)
	}

	return nil
}
func resourceIBMCISFilterValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisFilterDescription,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "Filter-creation"})

	ibmCISFiltersResourceValidator := ResourceValidator{ResourceName: ibmCISFilters, Schema: validateSchema}
	return &ibmCISFiltersResourceValidator
}
