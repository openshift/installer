// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCRLocation            = "ibm_dns_custom_resolver_location"
	pdnsResolverID           = "resolver_id"
	pdnsCRLocationID         = "location_id"
	pdnsCRLocationSubnetCRN  = "subnet_crn"
	pdnsCRLocationEnable     = "enabled"
	pdnsCRLocationServerIP   = "dns_server_ip"
	pdnsCustomReolverEnabled = "cr_enabled"
)

func ResourceIBMPrivateDNSCRLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(fmt.Errorf("Resource ibm_dns_custom_resolver_location is deprecated. Use the composite Custom Resolver resource[ibm_dns_custom_resolver], which can handle locations instead."))
		},
		ReadContext: func(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(fmt.Errorf("Resource ibm_dns_custom_resolver_location is deprecated. Use the composite Custom Resolver resource[ibm_dns_custom_resolver], which can handle locations instead."))
		},
		UpdateContext: func(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(fmt.Errorf("Resource ibm_dns_custom_resolver_location is deprecated. Use the composite Custom Resolver resource[ibm_dns_custom_resolver], which can handle locations instead."))
		},
		DeleteContext: func(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(fmt.Errorf("Resource ibm_dns_custom_resolver_location is deprecated. Use the composite Custom Resolver resource[ibm_dns_custom_resolver], which can handle locations instead."))
		},
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID",
			},

			pdnsResolverID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom Resolver ID",
			},
			pdnsCRLocationID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRLocation ID",
			},

			pdnsCRLocationSubnetCRN: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CRLocation Subnet CRN",
			},

			pdnsCRLocationEnable: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "CRLocation Enabled",
			},

			pdnsCRLocationHealthy: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "CRLocation Healthy",
			},

			pdnsCRLocationServerIP: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRLocation Server IP",
			},
			pdnsCustomReolverEnabled: {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: flex.ApplyOnce,
			},
		},
		DeprecationMessage: "Resource ibm_dns_custom_resolver_location is deprecated. Using the deprecated resource can cause an outage. If you have used the `ibm_dns_custom_resolver_location` resource, change it to the composite Custom Resolver [ibm_dns_custom_resolver] resource before running terraform apply.",
	}
}
