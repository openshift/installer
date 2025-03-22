// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMKmsKMIPClientCertificates() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMKmsKMIPClientCertList,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
				Description:  "public or private",
			},
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Key protect Instance GUID",
				ForceNew:         true,
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"adapter_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The id of the KMIP adapter that contains the cert",
				ForceNew:     true,
				ExactlyOneOf: []string{"adapter_id", "adapter_name"},
			},
			"adapter_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the KMIP adapter that contains the cert",
				ForceNew:     true,
				ExactlyOneOf: []string{"adapter_id", "adapter_name"},
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Limit of how many adapters to be fetched",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Offset of adapters to be fetched",
			},
			"show_total_count": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to return the count of how many adapters there are in total",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "If show_total_count is true, this will contain the total number of certificates.",
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The certificates contained in the specified adapter",
				Elem: &schema.Resource{
					Schema: dataSourceIBMKmsKMIPClientCertificateBaseSchema(),
				},
			},
		},
	}
}

func dataSourceIBMKmsKMIPClientCertList(d *schema.ResourceData, meta interface{}) error {
	// initialize API
	instanceID := getInstanceIDFromResourceData(d, "instance_id")
	api, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	// get adapterID
	nameOrID, hasID := d.GetOk("adapter_id")
	if !hasID {
		nameOrID = d.Get("adapter_name")
	}
	adapterNameOrID := nameOrID.(string)

	opts := &kp.ListOptions{}
	if limit, ok := d.GetOk("limit"); ok {
		limitVal := uint32(limit.(int))
		opts.Limit = &limitVal
	}
	if offset, ok := d.GetOk("offset"); ok {
		offsetVal := uint32(offset.(int))
		opts.Offset = &offsetVal
	}
	showTotalCountEnabled := false
	if showTotalCount, ok := d.GetOk("show_total_count"); ok {
		showTotalCountEnabled = showTotalCount.(bool)
		opts.TotalCount = &showTotalCountEnabled
	}

	ctx := context.Background()
	adapter, err := api.GetKMIPAdapter(ctx, adapterNameOrID)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while retriving KMIP adapter to list certificates: %s", err)
	}
	if err = d.Set("adapter_id", adapter.ID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_id: %s", err)
	}
	if err = d.Set("adapter_name", adapter.Name); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_name: %s", err)
	}

	certs, err := api.GetKMIPClientCertificates(ctx, adapter.ID, opts)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while listing KMIP certs for adapter %s: %s", adapter.ID, err)
	}

	certsList := certs.Certificates

	// set computed values
	mySlice := make([]map[string]interface{}, 0, len(certsList))
	for _, cert := range certsList {
		certMap := dataSourceIBMKMSKmipClientCertToMap(cert)
		mySlice = append(mySlice, certMap)
	}
	d.Set("certificates", mySlice)
	d.SetId(adapter.ID)
	if showTotalCountEnabled {
		d.Set("total_count", certs.Metadata.TotalCount)
	}
	return nil
}

func dataSourceIBMKMSKmipClientCertToMap(model kp.KMIPClientCertificate) map[string]interface{} {
	modelMap := make(map[string]interface{})
	modelMap["cert_id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["created_by"] = model.CreatedBy
	return modelMap
}
