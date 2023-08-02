package organizations

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_organizations_delegated_administrator")
func ResourceDelegatedAdministrator() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceDelegatedAdministratorCreate,
		ReadWithoutTimeout:   resourceDelegatedAdministratorRead,
		DeleteWithoutTimeout: resourceDelegatedAdministratorDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"service_principal": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delegation_enabled_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDelegatedAdministratorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).OrganizationsConn(ctx)

	accountID := d.Get("account_id").(string)
	servicePrincipal := d.Get("service_principal").(string)
	input := &organizations.RegisterDelegatedAdministratorInput{
		AccountId:        aws.String(accountID),
		ServicePrincipal: aws.String(servicePrincipal),
	}

	_, err := conn.RegisterDelegatedAdministratorWithContext(ctx, input)
	if err != nil {
		return diag.Errorf("creating Organizations DelegatedAdministrator (%s): %s", accountID, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", accountID, servicePrincipal))

	return resourceDelegatedAdministratorRead(ctx, d, meta)
}

func resourceDelegatedAdministratorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).OrganizationsConn(ctx)

	accountID, servicePrincipal, err := DecodeOrganizationDelegatedAdministratorID(d.Id())
	if err != nil {
		return diag.Errorf("decoding ID AWS Organization (%s) DelegatedAdministrators: %s", d.Id(), err)
	}
	input := &organizations.ListDelegatedAdministratorsInput{
		ServicePrincipal: aws.String(servicePrincipal),
	}
	var delegatedAccount *organizations.DelegatedAdministrator
	err = conn.ListDelegatedAdministratorsPagesWithContext(ctx, input, func(page *organizations.ListDelegatedAdministratorsOutput, lastPage bool) bool {
		for _, delegated := range page.DelegatedAdministrators {
			if aws.StringValue(delegated.Id) == accountID {
				delegatedAccount = delegated
			}
		}

		return !lastPage
	})
	if err != nil {
		return diag.Errorf("listing AWS Organization (%s) DelegatedAdministrators: %s", d.Id(), err)
	}

	if delegatedAccount == nil {
		if !d.IsNewResource() {
			log.Printf("[WARN] AWS Organization DelegatedAdministrators not found (%s), removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(&retry.NotFoundError{})
	}

	d.Set("arn", delegatedAccount.Arn)
	d.Set("delegation_enabled_date", aws.TimeValue(delegatedAccount.DelegationEnabledDate).Format(time.RFC3339))
	d.Set("email", delegatedAccount.Email)
	d.Set("joined_method", delegatedAccount.JoinedMethod)
	d.Set("joined_timestamp", aws.TimeValue(delegatedAccount.JoinedTimestamp).Format(time.RFC3339))
	d.Set("name", delegatedAccount.Name)
	d.Set("status", delegatedAccount.Status)
	d.Set("account_id", accountID)
	d.Set("service_principal", servicePrincipal)

	return nil
}

func resourceDelegatedAdministratorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).OrganizationsConn(ctx)

	accountID, servicePrincipal, err := DecodeOrganizationDelegatedAdministratorID(d.Id())
	if err != nil {
		return diag.Errorf("decoding ID AWS Organization (%s) DelegatedAdministrators: %s", d.Id(), err)
	}
	input := &organizations.DeregisterDelegatedAdministratorInput{
		AccountId:        aws.String(accountID),
		ServicePrincipal: aws.String(servicePrincipal),
	}

	_, err = conn.DeregisterDelegatedAdministratorWithContext(ctx, input)
	if err != nil {
		return diag.Errorf("deleting Organizations DelegatedAdministrator (%s): %s", d.Id(), err)
	}
	return nil
}

func DecodeOrganizationDelegatedAdministratorID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return "", "", fmt.Errorf("expected ID in the form of account_id/service_principal, given: %q", id)
	}
	return idParts[0], idParts[1], nil
}
