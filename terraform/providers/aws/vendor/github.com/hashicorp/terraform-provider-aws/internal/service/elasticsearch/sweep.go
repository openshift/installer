//go:build sweep
// +build sweep

package elasticsearch

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_elasticsearch_domain", &resource.Sweeper{
		Name: "aws_elasticsearch_domain",
		F:    sweepDomains,
	})
}

func sweepDomains(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.(*conns.AWSClient).ElasticsearchConn
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	input := &elasticsearchservice.ListDomainNamesInput{}

	// ListDomainNames has no pagination support whatsoever
	output, err := conn.ListDomainNames(input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Elasticsearch Domain sweep for %s: %s", region, err)
		return errs.ErrorOrNil()
	}

	if err != nil {
		sweeperErr := fmt.Errorf("error listing Elasticsearch Domains: %w", err)
		log.Printf("[ERROR] %s", sweeperErr)
		errs = multierror.Append(errs, sweeperErr)
		return errs.ErrorOrNil()
	}

	if output == nil {
		log.Printf("[WARN] Skipping Elasticsearch Domain sweep for %s: empty response", region)
		return errs.ErrorOrNil()
	}

	for _, domainInfo := range output.DomainNames {
		if domainInfo == nil {
			continue
		}

		name := aws.StringValue(domainInfo.DomainName)

		// Elasticsearch Domains have regularly gotten stuck in a "being deleted" state
		// e.g. Deleted and Processing are both true for days in the API
		// Filter out domains that are Deleted already.

		output, err := FindDomainByName(conn, name)
		if err != nil {
			sweeperErr := fmt.Errorf("error describing Elasticsearch Domain (%s): %w", name, err)
			log.Printf("[ERROR] %s", sweeperErr)
			errs = multierror.Append(errs, sweeperErr)
			continue
		}

		if output != nil && aws.BoolValue(output.Deleted) {
			log.Printf("[INFO] Skipping Elasticsearch Domain (%s) with deleted status", name)
			continue
		}

		r := ResourceDomain()
		d := r.Data(nil)
		d.SetId(name)
		d.Set("domain_name", name)

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	if err = sweep.SweepOrchestrator(sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping Elasticsearch Domains for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping Elasticsearch Domain sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}
