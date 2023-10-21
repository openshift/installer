package acm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_acm_certificate_validation")
func resourceCertificateValidation() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCertificateValidationCreate,
		ReadWithoutTimeout:   resourceCertificateValidationRead,
		DeleteWithoutTimeout: schema.NoopContext,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(75 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"certificate_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"validation_record_fqdns": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCertificateValidationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ACMClient(ctx)

	arn := d.Get("certificate_arn").(string)
	certificate, err := findCertificateByARN(ctx, conn, arn)

	if err != nil {
		return diag.Errorf("reading ACM Certificate (%s): %s", arn, err)
	}

	if v := certificate.Type; v != types.CertificateTypeAmazonIssued {
		return diag.Errorf("ACM Certificate (%s) has type %s, no validation necessary", arn, v)
	}

	if v, ok := d.GetOk("validation_record_fqdns"); ok && v.(*schema.Set).Len() > 0 {
		fqdns := make(map[string]types.DomainValidation)

		for _, domainValidation := range certificate.DomainValidationOptions {
			if v := domainValidation.ValidationMethod; v != types.ValidationMethodDns {
				return diag.Errorf("validation_record_fqdns is not valid for %s validation", v)
			}

			if v := domainValidation.ResourceRecord; v != nil {
				if v := aws.ToString(v.Name); v != "" {
					fqdns[strings.TrimSuffix(v, ".")] = domainValidation
				}
			}
		}

		for _, v := range v.(*schema.Set).List() {
			delete(fqdns, strings.TrimSuffix(v.(string), "."))
		}

		if len(fqdns) > 0 {
			var errs *multierror.Error

			for fqdn, domainValidation := range fqdns {
				errs = multierror.Append(errs, fmt.Errorf("missing %s DNS validation record: %s", aws.ToString(domainValidation.DomainName), fqdn))
			}

			return diag.FromErr(errs)
		}
	}

	if _, err := waitCertificateIssued(ctx, conn, arn, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("waiting for ACM Certificate (%s) to be issued: %s", arn, err)
	}

	d.SetId(aws.ToTime(certificate.IssuedAt).String())

	return resourceCertificateValidationRead(ctx, d, meta)
}

func resourceCertificateValidationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ACMClient(ctx)

	arn := d.Get("certificate_arn").(string)
	certificate, err := findCertificateValidationByARN(ctx, conn, arn)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] ACM Certificate %s not found, removing from state", arn)
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading ACM Certificate (%s): %s", arn, err)
	}

	d.Set("certificate_arn", certificate.CertificateArn)

	return nil
}

func findCertificateValidationByARN(ctx context.Context, conn *acm.Client, arn string) (*types.CertificateDetail, error) {
	output, err := findCertificateByARN(ctx, conn, arn)

	if err != nil {
		return nil, err
	}

	if status := output.Status; status != types.CertificateStatusIssued {
		return nil, &retry.NotFoundError{
			Message:     string(status),
			LastRequest: arn,
		}
	}

	return output, nil
}

func statusCertificate(ctx context.Context, conn *acm.Client, arn string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Don't call findCertificateByARN as it maps useful status codes to NotFoundError.
		input := &acm.DescribeCertificateInput{
			CertificateArn: aws.String(arn),
		}

		output, err := findCertificate(ctx, conn, input)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, string(output.Status), nil
	}
}

func waitCertificateIssued(ctx context.Context, conn *acm.Client, arn string, timeout time.Duration) (*types.CertificateDetail, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.CertificateStatusPendingValidation),
		Target:  enum.Slice(types.CertificateStatusIssued),
		Refresh: statusCertificate(ctx, conn, arn),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*types.CertificateDetail); ok {
		switch output.Status {
		case types.CertificateStatusFailed:
			tfresource.SetLastError(err, errors.New(string(output.FailureReason)))
		case types.CertificateStatusRevoked:
			tfresource.SetLastError(err, errors.New(string(output.RevocationReason)))
		}

		return output, err
	}

	return nil, err
}
