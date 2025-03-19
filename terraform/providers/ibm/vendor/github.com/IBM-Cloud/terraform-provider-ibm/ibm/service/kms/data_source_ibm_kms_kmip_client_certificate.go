// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMKmsKMIPClientCertificateBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cert_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The id of the KMIP Client Certificate to be fetched",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the KMIP Client Certificate to be fetched",
		},
		"certificate": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The PEM-encoded contents of the certificate",
		},
		"created_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique identifier that is associated with the entity that created the adapter.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The date when a resource was created. The date format follows RFC 3339.",
		},
	}
}

func DataSourceIBMKmsKMIPClientCertificate() *schema.Resource {
	baseMap := dataSourceIBMKmsKMIPClientCertificateBaseSchema()

	baseMap["endpoint_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
		Description:  "public or private",
	}

	baseMap["instance_id"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Key protect Instance GUID",
		ForceNew:         true,
		DiffSuppressFunc: suppressKMSInstanceIDDiff,
	}
	baseMap["adapter_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The id of the KMIP adapter that contains the cert",
		ForceNew:     true,
		ExactlyOneOf: []string{"adapter_id", "adapter_name"},
	}
	baseMap["adapter_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The name of the KMIP adapter that contains the cert",
		ForceNew:     true,
		ExactlyOneOf: []string{"adapter_id", "adapter_name"},
	}

	baseMap["cert_id"].Optional = true
	baseMap["cert_id"].ExactlyOneOf = []string{"cert_id", "name"}

	baseMap["name"].Optional = true
	baseMap["name"].ExactlyOneOf = []string{"cert_id", "name"}

	return &schema.Resource{
		Read:     dataSourceIBMKmsKMIPClientCertRead,
		Importer: &schema.ResourceImporter{},
		Schema:   baseMap,
	}
}

func dataSourceIBMKmsKMIPClientCertRead(d *schema.ResourceData, meta interface{}) error {
	// initialize API
	instanceID := getInstanceIDFromResourceData(d, "instance_id")
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}

	// get adapterID and certID
	nameOrID, hasID := d.GetOk("adapter_id")
	if !hasID {
		nameOrID = d.Get("adapter_name")
	}
	adapterNameOrID := nameOrID.(string)

	nameOrID, hasID = d.GetOk("cert_id")
	if !hasID {
		nameOrID = d.Get("name")
	}
	certNameOrID := nameOrID.(string)

	ctx := context.Background()
	adapter, err := kpAPI.GetKMIPAdapter(ctx, adapterNameOrID)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while retriving KMIP adapter to get certificate: %s", err)
	}
	if err = d.Set("adapter_id", adapter.ID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_id: %s", err)
	}
	if err = d.Set("adapter_name", adapter.Name); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_name: %s", err)
	}

	cert, err := kpAPI.GetKMIPClientCertificate(ctx, adapterNameOrID, certNameOrID)
	if err != nil {
		return err
	}
	populateKMIPClientCertSchemaDataFromStruct(d, *cert, adapter.ID, instanceID)
	return nil
}
