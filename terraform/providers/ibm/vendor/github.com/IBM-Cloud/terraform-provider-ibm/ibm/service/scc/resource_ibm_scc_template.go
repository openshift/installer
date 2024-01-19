// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMSccTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Errorf("resource ibm_scc_template has been deprecated")
		},
		ReadContext: func(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Errorf("resource ibm_scc_template has been deprecated")
		},
		DeleteContext: func(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Errorf("resource ibm_scc_template has been deprecated")
		},
	}
}
