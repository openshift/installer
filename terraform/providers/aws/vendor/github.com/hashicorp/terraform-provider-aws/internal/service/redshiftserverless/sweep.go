//go:build sweep
// +build sweep

package redshiftserverless

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_redshiftserverless_namespace", &resource.Sweeper{
		Name: "aws_redshiftserverless_namespace",
		F:    sweepNamespaces,
		Dependencies: []string{
			"aws_redshiftserverless_workgroup",
		},
	})

	resource.AddTestSweepers("aws_redshiftserverless_workgroup", &resource.Sweeper{
		Name: "aws_redshiftserverless_workgroup",
		F:    sweepWorkgroups,
	})
}

func sweepNamespaces(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).RedshiftServerlessConn
	input := &redshiftserverless.ListNamespacesInput{}
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.ListNamespacesPages(input, func(page *redshiftserverless.ListNamespacesOutput, lastPage bool) bool {
		if len(page.Namespaces) == 0 {
			log.Print("[DEBUG] No Redshift Serverless Namespaces to sweep")
			return !lastPage
		}

		for _, namespace := range page.Namespaces {
			r := ResourceNamespace()
			d := r.Data(nil)
			d.SetId(aws.StringValue(namespace.NamespaceName))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing Redshift Serverless Namespaces: %w", err))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping Redshift Serverless Namespaces for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping Redshift Serverless Namespaces sweep for %s: %s", region, errs)
		return nil
	}

	return nil
}

func sweepWorkgroups(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).RedshiftServerlessConn
	input := &redshiftserverless.ListWorkgroupsInput{}
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.ListWorkgroupsPages(input, func(page *redshiftserverless.ListWorkgroupsOutput, lastPage bool) bool {
		if len(page.Workgroups) == 0 {
			log.Print("[DEBUG] No Redshift Serverless Workgroups to sweep")
			return !lastPage
		}

		for _, workgroup := range page.Workgroups {
			r := ResourceWorkgroup()
			d := r.Data(nil)
			d.SetId(aws.StringValue(workgroup.WorkgroupName))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing Redshift Serverless Workgroups: %w", err))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping Redshift Serverless Workgroups for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping Redshift Serverless Workgroups sweep for %s: %s", region, errs)
		return nil
	}

	return nil
}
