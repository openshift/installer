//go:build sweep
// +build sweep

package servicediscovery

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_service_discovery_http_namespace", &resource.Sweeper{
		Name: "aws_service_discovery_http_namespace",
		F:    sweepHTTPNamespaces,
		Dependencies: []string{
			"aws_service_discovery_service",
		},
	})

	resource.AddTestSweepers("aws_service_discovery_private_dns_namespace", &resource.Sweeper{
		Name: "aws_service_discovery_private_dns_namespace",
		F:    sweepPrivateDNSNamespaces,
		Dependencies: []string{
			"aws_service_discovery_service",
		},
	})

	resource.AddTestSweepers("aws_service_discovery_public_dns_namespace", &resource.Sweeper{
		Name: "aws_service_discovery_public_dns_namespace",
		F:    sweepPublicDNSNamespaces,
		Dependencies: []string{
			"aws_service_discovery_service",
		},
	})

	resource.AddTestSweepers("aws_service_discovery_service", &resource.Sweeper{
		Name: "aws_service_discovery_service",
		F:    sweepServices,
	})
}

func sweepHTTPNamespaces(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).ServiceDiscoveryConn
	sweepResources := make([]sweep.Sweepable, 0)

	namespaces, err := findNamespacesByType(context.Background(), conn, servicediscovery.NamespaceTypeHttp)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Service Discovery HTTP Namespace sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Service Discovery HTTP Namespaces (%s): %w", region, err)
	}

	for _, v := range namespaces {
		r := ResourceHTTPNamespace()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.Id))

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Service Discovery HTTP Namespaces (%s): %w", region, err)
	}

	return nil
}

func sweepPrivateDNSNamespaces(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).ServiceDiscoveryConn
	sweepResources := make([]sweep.Sweepable, 0)

	namespaces, err := findNamespacesByType(context.Background(), conn, servicediscovery.NamespaceTypeDnsPrivate)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Service Discovery Private DNS Namespace sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Service Discovery Private DNS Namespaces (%s): %w", region, err)
	}

	for _, v := range namespaces {
		r := ResourcePrivateDNSNamespace()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.Id))

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Service Discovery Private DNS Namespaces (%s): %w", region, err)
	}

	return nil
}

func sweepPublicDNSNamespaces(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).ServiceDiscoveryConn
	sweepResources := make([]sweep.Sweepable, 0)

	namespaces, err := findNamespacesByType(context.Background(), conn, servicediscovery.NamespaceTypeDnsPublic)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Service Discovery Public DNS Namespace sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Service Discovery Public DNS Namespaces (%s): %w", region, err)
	}

	for _, v := range namespaces {
		r := ResourcePrivateDNSNamespace()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.Id))

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Service Discovery Public DNS Namespaces (%s): %w", region, err)
	}

	return nil
}

func sweepServices(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).ServiceDiscoveryConn
	input := &servicediscovery.ListServicesInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	services, err := findServices(context.Background(), conn, input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Service Discovery Service sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Service Discovery Services (%s): %w", region, err)
	}

	for _, v := range services {
		r := ResourceService()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.Id))
		d.Set("force_destroy", true)

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Service Discovery Services (%s): %w", region, err)
	}

	return nil
}
