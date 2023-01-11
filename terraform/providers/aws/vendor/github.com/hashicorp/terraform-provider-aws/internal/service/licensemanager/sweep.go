//go:build sweep
// +build sweep

package licensemanager

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/licensemanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_licensemanager_license_configuration", &resource.Sweeper{
		Name: "aws_licensemanager_license_configuration",
		F:    sweepLicenseConfigurations,
	})
}

func sweepLicenseConfigurations(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).LicenseManagerConn
	input := &licensemanager.ListLicenseConfigurationsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = listLicenseConfigurationsPages(conn, input, func(page *licensemanager.ListLicenseConfigurationsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.LicenseConfigurations {
			r := ResourceLicenseConfiguration()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.LicenseConfigurationArn))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping License Manager License Configuration sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing License Manager License Configurations (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping License Manager License Configurations (%s): %w", region, err)
	}

	return nil
}
