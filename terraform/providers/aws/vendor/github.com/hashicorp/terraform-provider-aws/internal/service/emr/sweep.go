//go:build sweep
// +build sweep

package emr

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_emr_cluster", &resource.Sweeper{
		Name: "aws_emr_cluster",
		F:    sweepClusters,
	})

	resource.AddTestSweepers("aws_emr_studio", &resource.Sweeper{
		Name: "aws_emr_studio",
		F:    sweepStudios,
	})
}

func sweepClusters(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EMRConn(ctx)
	input := &emr.ListClustersInput{
		ClusterStates: aws.StringSlice([]string{emr.ClusterStateBootstrapping, emr.ClusterStateRunning, emr.ClusterStateStarting, emr.ClusterStateWaiting}),
	}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListClustersPagesWithContext(ctx, input, func(page *emr.ListClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Clusters {
			id := aws.StringValue(v.Id)

			_, err := conn.SetTerminationProtectionWithContext(ctx, &emr.SetTerminationProtectionInput{
				JobFlowIds:           aws.StringSlice([]string{id}),
				TerminationProtected: aws.Bool(false),
			})

			if err != nil {
				log.Printf("[ERROR] unsetting EMR Cluster (%s) termination protection: %s", id, err)
			}

			r := ResourceCluster()
			d := r.Data(nil)
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EMR Clusters sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EMR Clusters (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EMR Clusters (%s): %w", region, err)
	}

	return nil
}

func sweepStudios(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.(*conns.AWSClient).EMRConn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var sweeperErrs *multierror.Error
	input := &emr.ListStudiosInput{}

	err = conn.ListStudiosPagesWithContext(ctx, input, func(page *emr.ListStudiosOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, studio := range page.Studios {
			r := ResourceStudio()
			d := r.Data(nil)
			d.SetId(aws.StringValue(studio.StudioId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EMR Studios sweep for %s: %s", region, sweeperErrs)
		return nil
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing EMR Studios for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping EMR Studios for %s: %w", region, err))
	}

	return sweeperErrs.ErrorOrNil()
}
