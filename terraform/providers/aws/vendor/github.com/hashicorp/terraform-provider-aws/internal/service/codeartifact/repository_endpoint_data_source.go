package codeartifact

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codeartifact"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKDataSource("aws_codeartifact_repository_endpoint")
func DataSourceRepositoryEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceRepositoryEndpointRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(codeartifact.PackageFormat_Values(), false),
			},
			"domain_owner": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"repository_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRepositoryEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodeArtifactConn(ctx)
	domainOwner := meta.(*conns.AWSClient).AccountID
	domain := d.Get("domain").(string)
	repo := d.Get("repository").(string)
	format := d.Get("format").(string)
	params := &codeartifact.GetRepositoryEndpointInput{
		Domain:     aws.String(domain),
		Repository: aws.String(repo),
		Format:     aws.String(format),
	}

	if v, ok := d.GetOk("domain_owner"); ok {
		params.DomainOwner = aws.String(v.(string))
		domainOwner = v.(string)
	}

	log.Printf("[DEBUG] Getting CodeArtifact Repository Endpoint")
	out, err := conn.GetRepositoryEndpointWithContext(ctx, params)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "getting CodeArtifact Repository Endpoint: %s", err)
	}
	log.Printf("[DEBUG] CodeArtifact Repository Endpoint: %#v", out)

	d.SetId(fmt.Sprintf("%s:%s:%s:%s", domainOwner, domain, repo, format))
	d.Set("repository_endpoint", out.RepositoryEndpoint)
	d.Set("domain_owner", domainOwner)

	return diags
}
