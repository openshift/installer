//go:build sweep
// +build sweep

package storagegateway

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/storagegateway"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_storagegateway_gateway", &resource.Sweeper{
		Name: "aws_storagegateway_gateway",
		F:    sweepGateways,
		Dependencies: []string{
			"aws_storagegateway_file_system_association",
		},
	})

	resource.AddTestSweepers("aws_storagegateway_tape_pool", &resource.Sweeper{
		Name: "aws_storagegateway_tape_pool",
		F:    sweepTapePools,
	})

	resource.AddTestSweepers("aws_storagegateway_file_system_association", &resource.Sweeper{
		Name: "aws_storagegateway_file_system_association",
		F:    sweepFileSystemAssociations,
	})
}

func sweepGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).StorageGatewayConn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListGatewaysPagesWithContext(ctx, &storagegateway.ListGatewaysInput{}, func(page *storagegateway.ListGatewaysOutput, lastPage bool) bool {
		if len(page.Gateways) == 0 {
			log.Print("[DEBUG] No Storage Gateway Gateways to sweep")
			return true
		}

		for _, gateway := range page.Gateways {
			r := ResourceGateway()
			d := r.Data(nil)
			d.SetId(aws.StringValue(gateway.GatewayARN))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Storage Gateway Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Storage Gateway Gateways (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Storage Gateway Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepTapePools(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).StorageGatewayConn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListTapePoolsPagesWithContext(ctx, &storagegateway.ListTapePoolsInput{}, func(page *storagegateway.ListTapePoolsOutput, lastPage bool) bool {
		if len(page.PoolInfos) == 0 {
			log.Print("[DEBUG] No Storage Gateway Tape Pools to sweep")
			return true
		}

		for _, pool := range page.PoolInfos {
			r := ResourceTapePool()
			d := r.Data(nil)
			d.SetId(aws.StringValue(pool.PoolARN))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Storage Gateway Tape Pool sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Storage Gateway Tape Pools (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Storage Gateway Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepFileSystemAssociations(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).StorageGatewayConn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListFileSystemAssociationsPagesWithContext(ctx, &storagegateway.ListFileSystemAssociationsInput{}, func(page *storagegateway.ListFileSystemAssociationsOutput, lastPage bool) bool {
		if len(page.FileSystemAssociationSummaryList) == 0 {
			log.Print("[DEBUG] No Storage Gateway File System Associations to sweep")
			return true
		}

		for _, assoc := range page.FileSystemAssociationSummaryList {
			r := ResourceFileSystemAssociation()
			d := r.Data(nil)
			d.SetId(aws.StringValue(assoc.FileSystemAssociationARN))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Storage Gateway File System Association sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Storage Gateway File System Associations (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Storage Gateway File System Associations (%s): %w", region, err)
	}

	return nil
}
