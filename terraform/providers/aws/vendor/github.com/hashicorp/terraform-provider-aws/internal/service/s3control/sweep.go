//go:build sweep
// +build sweep

package s3control

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_s3_access_point", &resource.Sweeper{
		Name: "aws_s3_access_point",
		F:    sweepAccessPoints,
		Dependencies: []string{
			"aws_s3control_object_lambda_access_point",
		},
	})

	resource.AddTestSweepers("aws_s3control_multi_region_access_point", &resource.Sweeper{
		Name: "aws_s3control_multi_region_access_point",
		F:    sweepMultiRegionAccessPoints,
	})

	resource.AddTestSweepers("aws_s3control_object_lambda_access_point", &resource.Sweeper{
		Name: "aws_s3control_object_lambda_access_point",
		F:    sweepObjectLambdaAccessPoints,
	})
}

func sweepAccessPoints(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).S3ControlConn
	accountID := client.(*conns.AWSClient).AccountID
	input := &s3control.ListAccessPointsInput{
		AccountId: aws.String(accountID),
	}
	sweepResources := make([]sweep.Sweepable, 0)
	var sweeperErrs *multierror.Error

	err = conn.ListAccessPointsPages(input, func(page *s3control.ListAccessPointsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, accessPoint := range page.AccessPointList {
			r := ResourceAccessPoint()
			d := r.Data(nil)
			id, err := AccessPointCreateResourceID(aws.StringValue(accessPoint.AccessPointArn))
			if err != nil {
				sweeperErr := fmt.Errorf("error composing S3 Access Point ID (%s): %w", aws.StringValue(accessPoint.AccessPointArn), err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
			}
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping S3 Access Point sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil()
	}

	if err != nil {
		sweeperErr := fmt.Errorf("error listing S3 Access Points (%s): %w", region, err)
		if sweeperErrs.Len() > 0 {
			return multierror.Append(sweeperErr, sweeperErrs)
		}
		return sweeperErr
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		sweeperErr := fmt.Errorf("error sweeping S3 Access Points (%s): %w", region, err)
		if sweeperErrs.Len() > 0 {
			return multierror.Append(sweeperErr, sweeperErrs)
		}
		return sweeperErr
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepMultiRegionAccessPoints(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	if region != endpoints.UsWest2RegionID {
		log.Printf("[WARN] Skipping S3 Multi-Region Access Point sweep for region: %s", region)
		return nil
	}
	conn := client.(*conns.AWSClient).S3ControlConn
	accountID := client.(*conns.AWSClient).AccountID
	input := &s3control.ListMultiRegionAccessPointsInput{
		AccountId: aws.String(accountID),
	}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListMultiRegionAccessPointsPages(input, func(page *s3control.ListMultiRegionAccessPointsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, accessPoint := range page.AccessPoints {
			r := ResourceMultiRegionAccessPoint()
			d := r.Data(nil)
			d.SetId(MultiRegionAccessPointCreateResourceID(accountID, aws.StringValue(accessPoint.Name)))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping S3 Multi-Region Access Point sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing S3 Multi-Region Access Points (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping S3 Multi-Region Access Points (%s): %w", region, err)
	}

	return nil
}

func sweepObjectLambdaAccessPoints(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).S3ControlConn
	accountID := client.(*conns.AWSClient).AccountID
	input := &s3control.ListAccessPointsForObjectLambdaInput{
		AccountId: aws.String(accountID),
	}
	sweepResources := make([]sweep.Sweepable, 0)

	conn.ListAccessPointsForObjectLambdaPages(input, func(page *s3control.ListAccessPointsForObjectLambdaOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, accessPoint := range page.ObjectLambdaAccessPointList {
			r := ResourceObjectLambdaAccessPoint()
			d := r.Data(nil)
			d.SetId(ObjectLambdaAccessPointCreateResourceID(accountID, aws.StringValue(accessPoint.Name)))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping S3 Object Lambda Access Point sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing S3 Object Lambda Access Points (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping S3 Object Lambda Access Points (%s): %w", region, err)
	}

	return nil
}
