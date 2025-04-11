// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package pag

import (
	"context"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPag() *schema.Resource {
	riSchema := resourcecontroller.ResourceIBMResourceInstance().Schema

	riSchema["parameters_json"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "PAG parameters to pass in Json string format",
		ForceNew:    true,
		Optional:    true,
	}

	return &schema.Resource{
		Create:   resourcecontroller.ResourceIBMResourceInstanceCreate,
		Read:     resourcecontroller.ResourceIBMResourceInstanceRead,
		Update:   resourcecontroller.ResourceIBMResourceInstanceUpdate,
		Delete:   resourcecontroller.ResourceIBMResourceInstanceDelete,
		Exists:   resourcecontroller.ResourceIBMResourceInstanceExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: riSchema,
	}
}
