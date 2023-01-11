package licensemanager

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/licensemanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceAssociation() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAssociationCreate,
		ReadWithoutTimeout:   resourceAssociationRead,
		DeleteWithoutTimeout: resourceAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"license_configuration_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"resource_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func resourceAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LicenseManagerConn

	licenseConfigurationARN := d.Get("license_configuration_arn").(string)
	resourceARN := d.Get("resource_arn").(string)

	input := &licensemanager.UpdateLicenseSpecificationsForResourceInput{
		AddLicenseSpecifications: []*licensemanager.LicenseSpecification{{
			LicenseConfigurationArn: aws.String(licenseConfigurationARN),
		}},
		ResourceArn: aws.String(resourceARN),
	}

	log.Printf("[DEBUG] Creating License Manager Association: %s", input)
	_, err := conn.UpdateLicenseSpecificationsForResourceWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating License Manager Association: %s", err)
	}

	d.SetId(AssociationCreateResourceID(resourceARN, licenseConfigurationARN))

	return resourceAssociationRead(ctx, d, meta)
}

func resourceAssociationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LicenseManagerConn

	resourceARN, licenseConfigurationARN, err := AssociationParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	err = FindAssociation(ctx, conn, resourceARN, licenseConfigurationARN)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] License Manager Association %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading License Manager Association (%s): %s", d.Id(), err)
	}

	d.Set("license_configuration_arn", licenseConfigurationARN)
	d.Set("resource_arn", resourceARN)

	return nil
}

func resourceAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LicenseManagerConn

	resourceARN, licenseConfigurationARN, err := AssociationParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	input := &licensemanager.UpdateLicenseSpecificationsForResourceInput{
		RemoveLicenseSpecifications: []*licensemanager.LicenseSpecification{{
			LicenseConfigurationArn: aws.String(licenseConfigurationARN),
		}},
		ResourceArn: aws.String(resourceARN),
	}

	_, err = conn.UpdateLicenseSpecificationsForResourceWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("deleting License Manager Association (%s): %s", d.Id(), err)
	}

	return nil
}

func FindAssociation(ctx context.Context, conn *licensemanager.LicenseManager, resourceARN, licenseConfigurationARN string) error {
	input := &licensemanager.ListLicenseSpecificationsForResourceInput{
		ResourceArn: aws.String(resourceARN),
	}
	var output []*licensemanager.LicenseSpecification

	err := listLicenseSpecificationsForResourcePagesWithContext(ctx, conn, input, func(page *licensemanager.ListLicenseSpecificationsForResourceOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.LicenseSpecifications {
			if v != nil {
				output = append(output, v)
			}
		}

		return !lastPage
	})

	if err != nil {
		return err
	}

	for _, v := range output {
		if aws.StringValue(v.LicenseConfigurationArn) == licenseConfigurationARN {
			return nil
		}
	}

	return &resource.NotFoundError{}
}

const associationResourceIDSeparator = ","

func AssociationCreateResourceID(resourceARN, licenseConfigurationARN string) string {
	parts := []string{resourceARN, licenseConfigurationARN}
	id := strings.Join(parts, associationResourceIDSeparator)

	return id
}

func AssociationParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, associationResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected ResourceARN%[2]sLicenseConfigurationARN", id, associationResourceIDSeparator)
}
