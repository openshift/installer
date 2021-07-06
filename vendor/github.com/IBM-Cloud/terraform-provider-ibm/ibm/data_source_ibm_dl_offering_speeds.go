// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dlSpeeds               = "offering_speeds"
	dlLinkSpeed            = "link_speed"
	dlOfferingType         = "offering_type"
	dlMacSecEnabled        = "macsec_enabled"
	dlMeteringCapabilities = "capabilities"
)

func dataSourceIBMDLOfferingSpeeds() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLOfferingSpeedsRead,
		Schema: map[string]*schema.Schema{
			dlOfferingType: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Direct Link offering type",
				ValidateFunc: InvokeDataSourceValidator("ibm_dl_offering_speeds", dlOfferingType),
			},
			dlSpeeds: {
				Type:        schema.TypeList,
				Description: "Collection of direct link speeds",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlMeteringCapabilities: {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of capabilities for billing option",
						},
						dlLinkSpeed: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Direct Link offering speed for the specified offering type",
						},
						dlMacSecEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicate whether speed supports MACsec",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMDLOfferingSpeedsRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	dlType := d.Get(dlOfferingType).(string)
	listSpeedsOptionsModel := &directlinkv1.ListOfferingTypeSpeedsOptions{}
	listSpeedsOptionsModel.OfferingType = &dlType
	listSpeeds, detail, err := directLink.ListOfferingTypeSpeeds(listSpeedsOptionsModel)

	if err != nil {
		log.Printf("Error reading list of direct link offering speeds:%s\n%s", err, detail)
		return err
	}
	speeds := make([]map[string]interface{}, 0)
	for _, instance := range listSpeeds.Speeds {
		speed := map[string]interface{}{}
		if instance.Capabilities != nil {
			speed[dlMeteringCapabilities] = flattenStringList(instance.Capabilities)
		}
		if instance.LinkSpeed != nil {
			speed[dlLinkSpeed] = *instance.LinkSpeed
		}
		if instance.MacsecEnabled != nil {
			speed[dlMacSecEnabled] = *instance.MacsecEnabled
		}
		speeds = append(speeds, speed)
	}
	d.SetId(dataSourceIBMDLOfferingSpeedsID(d))
	d.Set(dlSpeeds, speeds)
	return nil
}

// dataSourceIBMDLOfferingSpeedsID returns a reasonable ID for a direct link speeds list.
func dataSourceIBMDLOfferingSpeedsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func datasourceIBMDLOfferingSpeedsValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 2)
	dlTypeAllowedValues := "dedicated, connect"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 dlOfferingType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              dlTypeAllowedValues})

	ibmDLOfferingSpeedsDatasourceValidator := ResourceValidator{ResourceName: "ibm_dl_offering_speeds", Schema: validateSchema}
	return &ibmDLOfferingSpeedsDatasourceValidator
}
