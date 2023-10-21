//go:build sweep
// +build sweep

package appstream

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_appstream_directory_config", &resource.Sweeper{
		Name: "aws_appstream_directory_config",
		F:    sweepDirectoryConfigs,
	})

	resource.AddTestSweepers("aws_appstream_fleet", &resource.Sweeper{
		Name: "aws_appstream_fleet",
		F:    sweepFleets,
	})

	resource.AddTestSweepers("aws_appstream_image_builder", &resource.Sweeper{
		Name: "aws_appstream_image_builder",
		F:    sweepImageBuilders,
	})

	resource.AddTestSweepers("aws_appstream_stack", &resource.Sweeper{
		Name: "aws_appstream_stack",
		F:    sweepStacks,
	})
}

func sweepDirectoryConfigs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).AppStreamConn(ctx)
	input := &appstream.DescribeDirectoryConfigsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = describeDirectoryConfigsPages(ctx, conn, input, func(page *appstream.DescribeDirectoryConfigsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.DirectoryConfigs {
			r := ResourceDirectoryConfig()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.DirectoryName))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping AppStream Directory Config sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing AppStream Directory Configs (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping AppStream Directory Configs (%s): %w", region, err)
	}

	return nil
}

func sweepFleets(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).AppStreamConn(ctx)
	input := &appstream.DescribeFleetsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = describeFleetsPages(ctx, conn, input, func(page *appstream.DescribeFleetsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Fleets {
			r := ResourceFleet()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.Name))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping AppStream Fleet sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing AppStream Fleets (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping AppStream Fleets (%s): %w", region, err)
	}

	return nil
}

func sweepImageBuilders(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).AppStreamConn(ctx)
	input := &appstream.DescribeImageBuildersInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = describeImageBuildersPages(ctx, conn, input, func(page *appstream.DescribeImageBuildersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.ImageBuilders {
			r := ResourceImageBuilder()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.Name))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping AppStream Image Builder sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing AppStream Image Builders (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping AppStream Image Builders (%s): %w", region, err)
	}

	return nil
}

func sweepStacks(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).AppStreamConn(ctx)
	input := &appstream.DescribeStacksInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = describeStacksPages(ctx, conn, input, func(page *appstream.DescribeStacksOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Stacks {
			r := ResourceStack()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.Name))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping AppStream Stack sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing AppStream Stacks (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping AppStream Stacks (%s): %w", region, err)
	}

	return nil
}
