//go:build sweep
// +build sweep

package apigateway

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_api_gateway_rest_api", &resource.Sweeper{
		Name: "aws_api_gateway_rest_api",
		F:    sweepRestAPIs,
	})

	resource.AddTestSweepers("aws_api_gateway_vpc_link", &resource.Sweeper{
		Name: "aws_api_gateway_vpc_link",
		F:    sweepVPCLinks,
	})

	resource.AddTestSweepers("aws_api_gateway_client_certificate", &resource.Sweeper{
		Name: "aws_api_gateway_client_certificate",
		F:    sweepClientCertificates,
	})

	resource.AddTestSweepers("aws_api_gateway_usage_plan", &resource.Sweeper{
		Name: "aws_api_gateway_usage_plan",
		F:    sweepUsagePlans,
	})

	resource.AddTestSweepers("aws_api_gateway_api_key", &resource.Sweeper{
		Name: "aws_api_gateway_api_key",
		F:    sweepAPIKeys,
		Dependencies: []string{
			"aws_api_gateway_usage_plan",
		},
	})

	resource.AddTestSweepers("aws_api_gateway_domain_name", &resource.Sweeper{
		Name: "aws_api_gateway_domain_name",
		F:    sweepDomainNames,
	})
}

func sweepRestAPIs(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).APIGatewayConn

	err = conn.GetRestApisPages(&apigateway.GetRestApisInput{}, func(page *apigateway.GetRestApisOutput, lastPage bool) bool {
		for _, item := range page.Items {
			input := &apigateway.DeleteRestApiInput{
				RestApiId: item.Id,
			}
			log.Printf("[INFO] Deleting API Gateway REST API: %s", input)
			// TooManyRequestsException: Too Many Requests can take over a minute to resolve itself
			err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				_, err := conn.DeleteRestApi(input)
				if err != nil {
					if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeTooManyRequestsException) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete API Gateway REST API %s: %s", *item.Name, err)
				continue
			}
		}
		return !lastPage
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping API Gateway REST API sweep for %s: %s", region, err)
			return nil
		}
		return fmt.Errorf("Error retrieving API Gateway REST APIs: %s", err)
	}

	return nil
}

func sweepVPCLinks(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).APIGatewayConn

	sweepResources := make([]sweep.Sweepable, 0)
	var sweeperErrs *multierror.Error

	err = conn.GetVpcLinksPages(&apigateway.GetVpcLinksInput{}, func(page *apigateway.GetVpcLinksOutput, lastPage bool) bool {
		for _, item := range page.Items {
			id := aws.StringValue(item.Id)

			log.Printf("[INFO] Deleting API Gateway VPC Link (%s)", id)
			r := ResourceVPCLink()
			d := r.Data(nil)
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}
		return !lastPage
	})
	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping API Gateway VPC Link sweep for %s: %s", region, err)
		return nil
	}
	if err != nil {
		return fmt.Errorf("error retrieving API Gateway VPC Links: %w", err)
	}

	if err := sweep.SweepOrchestrator(sweepResources); err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping API Gateway VPC Links: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepClientCertificates(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).APIGatewayConn
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.GetClientCertificatesPages(&apigateway.GetClientCertificatesInput{}, func(page *apigateway.GetClientCertificatesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, clientCertificate := range page.Items {
			r := ResourceClientCertificate()
			d := r.Data(nil)
			d.SetId(aws.StringValue(clientCertificate.ClientCertificateId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing API Gateway Client Certificates for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping API Gateway Client Certificates for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping API Gateway Client Certificate sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}

func sweepUsagePlans(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	log.Printf("[INFO] Sweeping API Gateway Usage Plans for %s", region)

	conn := client.(*conns.AWSClient).APIGatewayConn
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.GetUsagePlansPages(&apigateway.GetUsagePlansInput{}, func(page *apigateway.GetUsagePlansOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		log.Printf("[INFO] API Gateway Usage Plans: %d", len(page.Items))

		for _, up := range page.Items {
			r := ResourceUsagePlan()
			d := r.Data(nil)
			d.SetId(aws.StringValue(up.Id))
			d.Set("api_stages", flattenAPIStages(up.ApiStages))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing API Gateway Usage Plans for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping API Gateway Usage Plans for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping API Gateway Usage Plan sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}

func sweepAPIKeys(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	log.Printf("[INFO] Sweeping API Gateway API Keys for %s", region)

	conn := client.(*conns.AWSClient).APIGatewayConn
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.GetApiKeysPages(&apigateway.GetApiKeysInput{}, func(page *apigateway.GetApiKeysOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		log.Printf("[INFO] API Gateway API Keys: %d", len(page.Items))

		for _, ak := range page.Items {
			r := ResourceAPIKey()
			d := r.Data(nil)
			d.SetId(aws.StringValue(ak.Id))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing API Gateway API Keys for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping API Gateway API Keys for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping API Gateway API Key sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}

func sweepDomainNames(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	log.Printf("[INFO] Sweeping API Gateway Domain Names for %s", region)

	conn := client.(*conns.AWSClient).APIGatewayConn
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.GetDomainNamesPages(&apigateway.GetDomainNamesInput{}, func(page *apigateway.GetDomainNamesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		log.Printf("[INFO] API Gateway Domain Names: %d", len(page.Items))

		for _, dn := range page.Items {
			r := ResourceDomainName()
			d := r.Data(nil)
			d.SetId(aws.StringValue(dn.DomainName))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("describing API Gateway Domain Names for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("sweeping API Gateway Domain Names for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping API Gateway Domain Name sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}
