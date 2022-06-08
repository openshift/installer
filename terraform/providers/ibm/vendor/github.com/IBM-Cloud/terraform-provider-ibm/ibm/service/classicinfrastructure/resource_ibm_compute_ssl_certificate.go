// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package classicinfrastructure

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func ResourceIBMComputeSSLCertificate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeSSLCertificateCreate,
		Read:     resourceIBMComputeSSLCertificateRead,
		Update:   resourceIBMComputeSSLCertificateUpdate,
		Delete:   resourceIBMComputeSSLCertificateDelete,
		Exists:   resourceIBMComputeSSLCertificateExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc:   normalizeCert,
				Description: "SSL Certifcate",
			},

			"intermediate_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				StateFunc:   normalizeCert,
				Description: "Intermediate certificate value",
			},

			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				StateFunc:   normalizeCert,
				Description: "SSL Private Key",
			},

			"common_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Common name",
			},

			"organization_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization name",
			},

			"validity_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Validity begins from",
			},

			"validity_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Validity days",
			},

			"validity_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Validity ends before",
			},

			"key_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "SSL key size",
			},

			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "certificate creation date",
			},

			"modify_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "certificate modificatiob date",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Tags set for resource",
			},
		},
	}
}

func resourceIBMComputeSSLCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetSecurityCertificateService(sess.SetRetries(0))

	template := datatypes.Security_Certificate{
		Certificate:             sl.String(d.Get("certificate").(string)),
		IntermediateCertificate: sl.String(d.Get("intermediate_certificate").(string)),
		PrivateKey:              sl.String(d.Get("private_key").(string)),
	}

	log.Printf("[INFO] Creating Security Certificate")

	cert, err := service.CreateObject(&template)

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating Security Certificate: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", *cert.Id))

	return resourceIBMComputeSSLCertificateRead(d, meta)
}

func resourceIBMComputeSSLCertificateRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetSecurityCertificateService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Not  a valid ID, must be an integer: %s", err)
	}

	cert, err := service.Id(id).GetObject()

	if err != nil {
		return fmt.Errorf("[ERROR] Unable  to get Security Certificate: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", *cert.Id))
	d.Set("certificate", *cert.Certificate)
	if cert.IntermediateCertificate != nil {
		d.Set("intermediate_certificate", *cert.IntermediateCertificate)
	}
	if cert.PrivateKey != nil {
		d.Set("private_key", *cert.PrivateKey)
	}
	d.Set("common_name", *cert.CommonName)
	d.Set("organization_name", *cert.OrganizationName)
	validityBegin := *cert.ValidityBegin
	d.Set("validity_begin", validityBegin.String())
	d.Set("validity_days", *cert.ValidityDays)
	validityEnd := *cert.ValidityEnd
	d.Set("validity_end", validityEnd.String())
	d.Set("key_size", *cert.KeySize)
	createDate := *cert.CreateDate
	d.Set("create_date", createDate.String())
	modifyDate := *cert.ModifyDate
	d.Set("modify_date", modifyDate.String())

	return nil
}

func resourceIBMComputeSSLCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	//Only tags are updated and that too locally hence nothing to validate and update in terms of real API at this point
	return nil
}

func resourceIBMComputeSSLCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetSecurityCertificateService(sess)
	id, err := strconv.Atoi(d.Id())
	_, err = service.Id(id).DeleteObject()
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting Security Certificate %d: %s", id, err)
	}

	return nil
}

func resourceIBMComputeSSLCertificateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetSecurityCertificateService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("[ERROR] Not  a valid ID, must be an integer: %s", err)
	}

	cert, err := service.Id(id).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting compute ssl certificate: %s", err)
	}
	return cert.Id != nil && *cert.Id == id, nil
}

func normalizeCert(cert interface{}) string {
	if cert == nil || cert == (*string)(nil) {
		return ""
	}

	switch cert.(type) {
	case string:
		return strings.TrimSpace(cert.(string))
	default:
		return ""
	}
}
