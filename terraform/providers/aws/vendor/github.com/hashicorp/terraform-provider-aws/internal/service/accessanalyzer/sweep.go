//go:build sweep
// +build sweep

package accessanalyzer

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_accessanalyzer_analyzer", &resource.Sweeper{
		Name: "aws_accessanalyzer_analyzer",
		F:    sweepAnalyzers,
	})
}

func sweepAnalyzers(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).AccessAnalyzerClient(ctx)
	input := &accessanalyzer.ListAnalyzersInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	pages := accessanalyzer.NewListAnalyzersPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping IAM Access Analyzer Analyzer sweep for %s: %s", region, err)
			return nil
		}

		if err != nil {
			return fmt.Errorf("listing IAM Access Analyzer Analyzers (%s): %w", region, err)
		}

		for _, v := range page.Analyzers {
			r := resourceAnalyzer()
			d := r.Data(nil)
			d.SetId(aws.ToString(v.Name))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("sweeping IAM Access Analyzer Analyzers (%s): %w", region, err)
	}

	return nil
}
