//go:build sweep
// +build sweep

package applicationinsights

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationinsights"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_applicationinsights_application", &resource.Sweeper{
		Name: "aws_applicationinsights_application",
		F:    sweepApplications,
	})
}

func sweepApplications(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).ApplicationInsightsConn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.ListApplicationsPagesWithContext(ctx, &applicationinsights.ListApplicationsInput{}, func(resp *applicationinsights.ListApplicationsOutput, lastPage bool) bool {
		if len(resp.ApplicationInfoList) == 0 {
			log.Print("[DEBUG] No ApplicationInsights Applications to sweep")
			return !lastPage
		}

		for _, c := range resp.ApplicationInfoList {
			r := ResourceApplication()
			d := r.Data(nil)
			d.SetId(aws.StringValue(c.ResourceGroupName))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing ApplicationInsights Applications: %w", err))
		// in case work can be done, don't jump out yet
	}

	if err = sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping ApplicationInsights Applications for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping ApplicationInsights Application sweep for %s: %s", region, err)
		return nil
	}

	return errs.ErrorOrNil()
}
