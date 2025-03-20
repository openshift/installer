// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMKmsKMIPClientCertificate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKmsKMIPClientCertCreate,
		Read:     resourceIBMKmsKMIPClientCertRead,
		Delete:   resourceIBMKmsKMIPClientCertDelete,
		Exists:   resourceIBMKmsKMIPClientCertExists,
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name or UUID of the KMIP adapter that contains the cert",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The name of the KMIP client certificate",
			},
			"cert_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the KMIP cert",
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
		},
	}
}

func resourceIBMKmsKMIPClientCertCreate(d *schema.ResourceData, meta interface{}) error {
	certToCreate, adapterID, instanceID, err := ExtractAndValidateKMIPClientCertDataFromSchema(d)
	if err != nil {
		return err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	if strings.Contains((kpAPI.URL).String(), "private") || strings.Contains(kpAPI.Config.BaseURL, "private") {
		d.Set("endpoint_type", "private")
	} else {
		d.Set("endpoint_type", "public")
	}
	ctx := context.Background()
	cert, err := kpAPI.CreateKMIPClientCertificate(ctx,
		adapterID,
		certToCreate.Certificate,
		kp.WithKMIPClientCertName(certToCreate.Name),
	)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while creating KMIP client certificate: %s", err)
	}
	return populateKMIPClientCertSchemaDataFromStruct(d, *cert, adapterID, instanceID)
}

func resourceIBMKmsKMIPClientCertRead(d *schema.ResourceData, meta interface{}) error {
	// use instanceID and adapterID here to support terraform import case
	instanceID, adapterID, certID, err := splitCertID(d.Id())
	if err != nil {
		return err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	if strings.Contains((kpAPI.URL).String(), "private") || strings.Contains(kpAPI.Config.BaseURL, "private") {
		d.Set("endpoint_type", "private")
	} else {
		d.Set("endpoint_type", "public")
	}
	ctx := context.Background()
	cert, err := kpAPI.GetKMIPClientCertificate(ctx, adapterID, certID)
	if err != nil {
		return err
	}
	err = populateKMIPClientCertSchemaDataFromStruct(d, *cert, adapterID, instanceID)
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMKmsKMIPClientCertDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Get("instance_id").(string)
	adapterID := d.Get("adapter_id").(string)
	_, _, certID, err := splitCertID(d.Id())
	if err != nil {
		return err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	err = kpAPI.DeleteKMIPClientCertificate(context.Background(), adapterID, certID)
	return err
}

func resourceIBMKmsKMIPClientCertExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	// use instanceID and adapterID here to support terraform import case
	instanceID, adapterID, certID, err := splitCertID(d.Id())
	if err != nil {
		return false, err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	_, err = kpAPI.GetKMIPClientCertificate(ctx, adapterID, certID)
	if err != nil {
		if kpError, ok := err.(*kp.Error); ok {
			if kpError.StatusCode == 404 {
				return false, nil
			}
		}
		return false, wrapError(err, "Error checking KMIP Client Certificate existence")
	}
	return true, nil
}

func ExtractAndValidateKMIPClientCertDataFromSchema(d *schema.ResourceData) (cert kp.KMIPClientCertificate, adapterIDStr string, instanceID string, err error) {
	err = nil
	instanceID = getInstanceIDFromResourceData(d, "instance_id")

	cert = kp.KMIPClientCertificate{}
	if name, ok := d.GetOk("name"); ok {
		nameStr, ok2 := name.(string)
		if !ok2 {
			err = flex.FmtErrorf("[ERROR] Error converting name to string")
			return
		}
		cert.Name = nameStr
	}
	if certPayload, ok := d.GetOk("certificate"); ok {
		certStr, ok2 := certPayload.(string)
		if !ok2 {
			err = flex.FmtErrorf("[ERROR] Error converting certificate to string")
			return
		}
		cert.Certificate = certStr
	}

	adapterID := d.Get("adapter_id")
	adapterIDStr = adapterID.(string)
	return
}

func populateKMIPClientCertSchemaDataFromStruct(d *schema.ResourceData, cert kp.KMIPClientCertificate, adapterID string, instanceID string) (err error) {
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, adapterID, cert.ID))

	if err = d.Set("name", cert.Name); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting name: %s", err)
	}
	if err = d.Set("adapter_id", adapterID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_id: %s", err)
	}
	if err = d.Set("instance_id", instanceID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting instance_id: %s", err)
	}
	if err = d.Set("cert_id", cert.ID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting cert_id: %s", err)
	}
	if err = d.Set("certificate", cert.Certificate); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting certificate: %s", err)
	}
	if cert.CreatedAt != nil {
		if err = d.Set("created_at", cert.CreatedAt.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_at: %s", err)
		}
		if err = d.Set("created_by", cert.CreatedBy); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_by: %s", err)
		}
	}
	return nil
}

func splitCertID(terraformId string) (instanceID, adapterID, certID string, err error) {
	split, err := flex.SepIdParts(terraformId, "/")
	if err != nil {
		return "", "", "", err
	}
	if len(split) != 3 {
		return "", "", "", flex.FmtErrorf("[ERROR] The given id %s does not contain all expected sections, should be of format instance_id/adapter_id/cert_id", terraformId)
	}
	for index, id := range split {
		if uuid.Validate(id) != nil {
			return "", "", "", flex.FmtErrorf("[ERROR] The given id %s at index %d of instance_id/adapter_id/cert_id is not a valid UUID", id, index)
		}
	}
	return split[0], split[1], split[2], nil
}
