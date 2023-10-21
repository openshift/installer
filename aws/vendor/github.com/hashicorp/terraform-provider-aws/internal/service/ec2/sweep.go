//go:build sweep
// +build sweep

package ec2

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_customer_gateway", &resource.Sweeper{
		Name: "aws_customer_gateway",
		F:    sweepCustomerGateways,
		Dependencies: []string{
			"aws_vpn_connection",
		},
	})

	resource.AddTestSweepers("aws_ec2_capacity_reservation", &resource.Sweeper{
		Name: "aws_ec2_capacity_reservation",
		F:    sweepCapacityReservations,
	})

	resource.AddTestSweepers("aws_ec2_carrier_gateway", &resource.Sweeper{
		Name: "aws_ec2_carrier_gateway",
		F:    sweepCarrierGateways,
	})

	resource.AddTestSweepers("aws_ec2_client_vpn_endpoint", &resource.Sweeper{
		Name: "aws_ec2_client_vpn_endpoint",
		F:    sweepClientVPNEndpoints,
		Dependencies: []string{
			"aws_ec2_client_vpn_network_association",
		},
	})

	resource.AddTestSweepers("aws_ec2_client_vpn_network_association", &resource.Sweeper{
		Name: "aws_ec2_client_vpn_network_association",
		F:    sweepClientVPNNetworkAssociations,
	})

	resource.AddTestSweepers("aws_ec2_fleet", &resource.Sweeper{
		Name: "aws_ec2_fleet",
		F:    sweepFleets,
	})

	resource.AddTestSweepers("aws_ebs_volume", &resource.Sweeper{
		Name: "aws_ebs_volume",
		Dependencies: []string{
			"aws_instance",
		},
		F: sweepEBSVolumes,
	})

	resource.AddTestSweepers("aws_ebs_snapshot", &resource.Sweeper{
		Name: "aws_ebs_snapshot",
		F:    sweepEBSSnapshots,
		Dependencies: []string{
			"aws_ami",
		},
	})

	resource.AddTestSweepers("aws_egress_only_internet_gateway", &resource.Sweeper{
		Name: "aws_egress_only_internet_gateway",
		F:    sweepEgressOnlyInternetGateways,
	})

	resource.AddTestSweepers("aws_eip", &resource.Sweeper{
		Name: "aws_eip",
		Dependencies: []string{
			"aws_vpc",
		},
		F: sweepEIPs,
	})

	resource.AddTestSweepers("aws_flow_log", &resource.Sweeper{
		Name: "aws_flow_log",
		F:    sweepFlowLogs,
	})

	resource.AddTestSweepers("aws_ec2_host", &resource.Sweeper{
		Name: "aws_ec2_host",
		F:    sweepHosts,
		Dependencies: []string{
			"aws_instance",
		},
	})

	resource.AddTestSweepers("aws_instance", &resource.Sweeper{
		Name: "aws_instance",
		F:    sweepInstances,
		Dependencies: []string{
			"aws_autoscaling_group",
			"aws_spot_fleet_request",
			"aws_spot_instance_request",
		},
	})

	resource.AddTestSweepers("aws_internet_gateway", &resource.Sweeper{
		Name: "aws_internet_gateway",
		Dependencies: []string{
			"aws_subnet",
		},
		F: sweepInternetGateways,
	})

	resource.AddTestSweepers("aws_key_pair", &resource.Sweeper{
		Name: "aws_key_pair",
		Dependencies: []string{
			"aws_elastic_beanstalk_environment",
			"aws_instance",
			"aws_spot_fleet_request",
			"aws_spot_instance_request",
		},
		F: sweepKeyPairs,
	})

	resource.AddTestSweepers("aws_launch_template", &resource.Sweeper{
		Name: "aws_launch_template",
		Dependencies: []string{
			"aws_autoscaling_group",
			"aws_batch_compute_environment",
		},
		F: sweepLaunchTemplates,
	})

	resource.AddTestSweepers("aws_nat_gateway", &resource.Sweeper{
		Name: "aws_nat_gateway",
		F:    sweepNATGateways,
	})

	resource.AddTestSweepers("aws_network_acl", &resource.Sweeper{
		Name: "aws_network_acl",
		F:    sweepNetworkACLs,
	})

	resource.AddTestSweepers("aws_network_interface", &resource.Sweeper{
		Name: "aws_network_interface",
		F:    sweepNetworkInterfaces,
		Dependencies: []string{
			"aws_db_proxy",
			"aws_directory_service_directory",
			"aws_ec2_client_vpn_endpoint",
			"aws_ec2_transit_gateway_vpc_attachment",
			"aws_eks_cluster",
			"aws_elb",
			"aws_instance",
			"aws_lb",
			"aws_nat_gateway",
			"aws_rds_cluster",
			"aws_rds_global_cluster",
		},
	})

	resource.AddTestSweepers("aws_ec2_network_insights_path", &resource.Sweeper{
		Name: "aws_ec2_network_insights_path",
		F:    sweepNetworkInsightsPaths,
	})

	resource.AddTestSweepers("aws_placement_group", &resource.Sweeper{
		Name: "aws_placement_group",
		F:    sweepPlacementGroups,
		Dependencies: []string{
			"aws_autoscaling_group",
			"aws_instance",
			"aws_launch_template",
			"aws_spot_fleet_request",
			"aws_spot_instance_request",
		},
	})

	resource.AddTestSweepers("aws_route_table", &resource.Sweeper{
		Name: "aws_route_table",
		F:    sweepRouteTables,
	})

	resource.AddTestSweepers("aws_security_group", &resource.Sweeper{
		Name: "aws_security_group",
		Dependencies: []string{
			"aws_subnet",
		},
		F: sweepSecurityGroups,
	})

	resource.AddTestSweepers("aws_spot_fleet_request", &resource.Sweeper{
		Name: "aws_spot_fleet_request",
		F:    sweepSpotFleetRequests,
	})

	resource.AddTestSweepers("aws_spot_instance_request", &resource.Sweeper{
		Name: "aws_spot_instance_request",
		F:    sweepSpotInstanceRequests,
	})

	resource.AddTestSweepers("aws_subnet", &resource.Sweeper{
		Name: "aws_subnet",
		F:    sweepSubnets,
		Dependencies: []string{
			"aws_appstream_fleet",
			"aws_appstream_image_builder",
			"aws_autoscaling_group",
			"aws_batch_compute_environment",
			"aws_elastic_beanstalk_environment",
			"aws_cloud9_environment_ec2",
			"aws_cloudhsm_v2_cluster",
			"aws_codestarconnections_host",
			"aws_db_subnet_group",
			"aws_directory_service_directory",
			"aws_dms_replication_instance",
			"aws_docdb_subnet_group",
			"aws_ec2_client_vpn_endpoint",
			"aws_ec2_transit_gateway_vpc_attachment",
			"aws_efs_file_system",
			"aws_eks_cluster",
			"aws_elasticache_cluster",
			"aws_elasticache_replication_group",
			"aws_elasticache_subnet_group",
			"aws_elasticsearch_domain",
			"aws_elb",
			"aws_emr_cluster",
			"aws_emr_studio",
			"aws_fsx_lustre_file_system",
			"aws_fsx_ontap_file_system",
			"aws_fsx_openzfs_file_system",
			"aws_fsx_windows_file_system",
			"aws_iot_topic_rule_destination",
			"aws_lambda_function",
			"aws_lb",
			"aws_memorydb_subnet_group",
			"aws_mq_broker",
			"aws_msk_cluster",
			"aws_network_interface",
			"aws_networkfirewall_firewall",
			"aws_opensearch_domain",
			"aws_redshift_cluster",
			"aws_redshift_subnet_group",
			"aws_route53_resolver_endpoint",
			"aws_sagemaker_notebook_instance",
			"aws_spot_fleet_request",
			"aws_spot_instance_request",
			"aws_vpc_endpoint",
			"aws_grafana_workspace",
		},
	})

	resource.AddTestSweepers("aws_ec2_transit_gateway_peering_attachment", &resource.Sweeper{
		Name: "aws_ec2_transit_gateway_peering_attachment",
		F:    sweepTransitGatewayPeeringAttachments,
	})

	resource.AddTestSweepers("aws_ec2_transit_gateway_multicast_domain", &resource.Sweeper{
		Name: "aws_ec2_transit_gateway_multicast_domain",
		F:    sweepTransitGatewayMulticastDomains,
	})

	resource.AddTestSweepers("aws_ec2_transit_gateway", &resource.Sweeper{
		Name: "aws_ec2_transit_gateway",
		F:    sweepTransitGateways,
		Dependencies: []string{
			"aws_dx_gateway_association",
			"aws_ec2_transit_gateway_vpc_attachment",
			"aws_ec2_transit_gateway_peering_attachment",
			"aws_vpn_connection",
		},
	})

	resource.AddTestSweepers("aws_ec2_transit_gateway_connect_peer", &resource.Sweeper{
		Name: "aws_ec2_transit_gateway_connect_peer",
		F:    sweepTransitGatewayConnectPeers,
	})

	resource.AddTestSweepers("aws_ec2_transit_gateway_connect", &resource.Sweeper{
		Name: "aws_ec2_transit_gateway_connect",
		F:    sweepTransitGatewayConnects,
		Dependencies: []string{
			"aws_ec2_transit_gateway_connect_peer",
		},
	})

	resource.AddTestSweepers("aws_ec2_transit_gateway_vpc_attachment", &resource.Sweeper{
		Name: "aws_ec2_transit_gateway_vpc_attachment",
		F:    sweepTransitGatewayVPCAttachments,
		Dependencies: []string{
			"aws_ec2_transit_gateway_connect",
			"aws_ec2_transit_gateway_multicast_domain",
		},
	})

	resource.AddTestSweepers("aws_vpc_dhcp_options", &resource.Sweeper{
		Name: "aws_vpc_dhcp_options",
		F:    sweepVPCDHCPOptions,
	})

	resource.AddTestSweepers("aws_vpc_endpoint_service", &resource.Sweeper{
		Name: "aws_vpc_endpoint_service",
		F:    sweepVPCEndpointServices,
		Dependencies: []string{
			"aws_vpc_endpoint",
		},
	})

	resource.AddTestSweepers("aws_vpc_endpoint", &resource.Sweeper{
		Name: "aws_vpc_endpoint",
		F:    sweepVPCEndpoints,
		Dependencies: []string{
			"aws_route_table",
			"aws_sagemaker_workforce",
		},
	})

	resource.AddTestSweepers("aws_vpc_peering_connection", &resource.Sweeper{
		Name: "aws_vpc_peering_connection",
		F:    sweepVPCPeeringConnections,
	})

	resource.AddTestSweepers("aws_vpc", &resource.Sweeper{
		Name: "aws_vpc",
		Dependencies: []string{
			"aws_ec2_carrier_gateway",
			"aws_egress_only_internet_gateway",
			"aws_internet_gateway",
			"aws_nat_gateway",
			"aws_network_acl",
			"aws_route_table",
			"aws_security_group",
			"aws_subnet",
			"aws_vpc_peering_connection",
			"aws_vpn_gateway",
		},
		F: sweepVPCs,
	})

	resource.AddTestSweepers("aws_vpn_connection", &resource.Sweeper{
		Name: "aws_vpn_connection",
		F:    sweepVPNConnections,
	})

	resource.AddTestSweepers("aws_vpn_gateway", &resource.Sweeper{
		Name: "aws_vpn_gateway",
		F:    sweepVPNGateways,
		Dependencies: []string{
			"aws_dx_gateway_association",
			"aws_vpn_connection",
		},
	})

	resource.AddTestSweepers("aws_vpc_ipam", &resource.Sweeper{
		Name: "aws_vpc_ipam",
		F:    sweepIPAMs,
	})

	resource.AddTestSweepers("aws_vpc_ipam_resource_discovery", &resource.Sweeper{
		Name: "aws_vpc_ipam_resource_discovery",
		F:    sweepIPAMResourceDiscoveries,
	})

	resource.AddTestSweepers("aws_ami", &resource.Sweeper{
		Name: "aws_ami",
		F:    sweepAMIs,
	})

	// aws_vpc_network_performance_metric_subscription
	resource.AddTestSweepers("aws_vpc_network_performance_metric_subscription", &resource.Sweeper{
		Name: "aws_vpc_network_performance_metric_subscription",
		F:    sweepNetworkPerformanceMetricSubscriptions,
	})
}

func sweepCapacityReservations(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)

	resp, err := conn.DescribeCapacityReservationsWithContext(ctx, &ec2.DescribeCapacityReservationsInput{})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Capacity Reservation sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error retrieving EC2 Capacity Reservations: %s", err)
	}

	if len(resp.CapacityReservations) == 0 {
		log.Print("[DEBUG] No EC2 Capacity Reservations to sweep")
		return nil
	}

	for _, r := range resp.CapacityReservations {
		if aws.StringValue(r.State) != ec2.CapacityReservationStateCancelled && aws.StringValue(r.State) != ec2.CapacityReservationStateExpired {
			id := aws.StringValue(r.CapacityReservationId)

			log.Printf("[INFO] Cancelling EC2 Capacity Reservation EC2 Instance: %s", id)

			opts := &ec2.CancelCapacityReservationInput{
				CapacityReservationId: aws.String(id),
			}

			_, err := conn.CancelCapacityReservationWithContext(ctx, opts)

			if err != nil {
				log.Printf("[ERROR] Error cancelling EC2 Capacity Reservation (%s): %s", id, err)
			}
		}
	}

	return nil
}

func sweepCarrierGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeCarrierGatewaysInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeCarrierGatewaysPagesWithContext(ctx, input, func(page *ec2.DescribeCarrierGatewaysOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.CarrierGateways {
			r := ResourceCarrierGateway()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.CarrierGatewayId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Carrier Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Carrier Gateways (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Carrier Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepClientVPNEndpoints(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeClientVpnEndpointsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeClientVpnEndpointsPagesWithContext(ctx, input, func(page *ec2.DescribeClientVpnEndpointsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.ClientVpnEndpoints {
			r := ResourceClientVPNEndpoint()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.ClientVpnEndpointId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Client VPN Endpoint sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Client VPN Endpoints (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Client VPN Endpoints (%s): %w", region, err)
	}

	return nil
}

func sweepClientVPNNetworkAssociations(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeClientVpnEndpointsInput{}
	var sweeperErrs *multierror.Error
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeClientVpnEndpointsPagesWithContext(ctx, input, func(page *ec2.DescribeClientVpnEndpointsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.ClientVpnEndpoints {
			input := &ec2.DescribeClientVpnTargetNetworksInput{
				ClientVpnEndpointId: v.ClientVpnEndpointId,
			}

			err := conn.DescribeClientVpnTargetNetworksPagesWithContext(ctx, input, func(page *ec2.DescribeClientVpnTargetNetworksOutput, lastPage bool) bool {
				if page == nil {
					return !lastPage
				}

				for _, v := range page.ClientVpnTargetNetworks {
					r := ResourceClientVPNNetworkAssociation()
					d := r.Data(nil)
					d.SetId(aws.StringValue(v.AssociationId))
					d.Set("client_vpn_endpoint_id", v.ClientVpnEndpointId)

					sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
				}

				return !lastPage
			})

			if sweep.SkipSweepError(err) {
				continue
			}

			if err != nil {
				sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing EC2 Client VPN Network Associations (%s): %w", region, err))
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Client VPN Network Association sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing EC2 Client VPN Endpoints (%s): %w", region, err))
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping EC2 Client VPN Network Associations (%s): %w", region, err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepFleets(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)

	var sweeperErrs *multierror.Error
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeFleetsPagesWithContext(ctx, &ec2.DescribeFleetsInput{}, func(page *ec2.DescribeFleetsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, fleet := range page.Fleets {
			if aws.StringValue(fleet.FleetState) == ec2.FleetStateCodeDeleted || aws.StringValue(fleet.FleetState) == ec2.FleetStateCodeDeletedTerminating {
				continue
			}

			r := ResourceFleet()
			d := r.Data(nil)
			d.SetId(aws.StringValue(fleet.FleetId))
			d.Set("terminate_instances", true)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Fleet sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil()
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("retrieving EC2 Fleets: %w", err))
	}

	if err := sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping EC2 Fleets for %s: %w", region, err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepEBSVolumes(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeVolumesInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeVolumesPagesWithContext(ctx, input, func(page *ec2.DescribeVolumesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Volumes {
			id := aws.StringValue(v.VolumeId)

			if state := aws.StringValue(v.State); state != ec2.VolumeStateAvailable {
				log.Printf("[INFO] Skipping EC2 EBS Volume (%s): %s", state, id)
				continue
			}

			r := ResourceEBSVolume()
			d := r.Data(nil)
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 EBS Volume sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 EBS Volumes (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 EBS Volumes (%s): %w", region, err)
	}

	return nil
}

func sweepEBSSnapshots(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	input := &ec2.DescribeSnapshotsInput{
		OwnerIds: aws.StringSlice([]string{"self"}),
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeSnapshotsPagesWithContext(ctx, input, func(page *ec2.DescribeSnapshotsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Snapshots {
			r := ResourceEBSSnapshot()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.SnapshotId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EBS Snapshot sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EBS Snapshots (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EBS Snapshots (%s): %w", region, err)
	}

	return nil
}

func sweepEgressOnlyInternetGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	input := &ec2.DescribeEgressOnlyInternetGatewaysInput{}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeEgressOnlyInternetGatewaysPagesWithContext(ctx, input, func(page *ec2.DescribeEgressOnlyInternetGatewaysOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.EgressOnlyInternetGateways {
			r := ResourceEgressOnlyInternetGateway()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.EgressOnlyInternetGatewayId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Egress-only Internet Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Egress-only Internet Gateways (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Egress-only Internet Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepEIPs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).EC2Conn(ctx)

	// There is currently no paginator or Marker/NextToken
	input := &ec2.DescribeAddressesInput{}

	output, err := conn.DescribeAddressesWithContext(ctx, input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 EIP sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing EC2 EIPs: %s", err)
	}

	if output == nil || len(output.Addresses) == 0 {
		log.Print("[DEBUG] No EC2 EIPs to sweep")
		return nil
	}

	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	for _, address := range output.Addresses {
		publicIP := aws.StringValue(address.PublicIp)

		if address.AssociationId != nil {
			log.Printf("[INFO] Skipping EC2 EIP (%s) with association: %s", publicIP, aws.StringValue(address.AssociationId))
			continue
		}

		r := ResourceEIP()
		d := r.Data(nil)
		if address.AllocationId != nil {
			d.SetId(aws.StringValue(address.AllocationId))
		} else {
			d.SetId(aws.StringValue(address.PublicIp))
		}

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	if err = sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping EC2 EIPs for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping EC2 EIP sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}

func sweepFlowLogs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeFlowLogsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeFlowLogsPagesWithContext(ctx, input, func(page *ec2.DescribeFlowLogsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, flowLog := range page.FlowLogs {
			r := ResourceFlowLog()
			d := r.Data(nil)
			d.SetId(aws.StringValue(flowLog.FlowLogId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Flow Log sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Flow Logs (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Flow Logs (%s): %w", region, err)
	}

	return nil
}

func sweepHosts(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeHostsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeHostsPagesWithContext(ctx, input, func(page *ec2.DescribeHostsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, host := range page.Hosts {
			r := ResourceHost()
			d := r.Data(nil)
			d.SetId(aws.StringValue(host.HostId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Host sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Hosts (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Hosts (%s): %w", region, err)
	}

	return nil
}

func sweepInstances(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.DescribeInstancesPagesWithContext(ctx, &ec2.DescribeInstancesInput{}, func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, reservation := range page.Reservations {
			if reservation == nil {
				continue
			}

			for _, instance := range reservation.Instances {
				id := aws.StringValue(instance.InstanceId)

				if instance.State != nil && aws.StringValue(instance.State.Name) == ec2.InstanceStateNameTerminated {
					log.Printf("[INFO] Skipping terminated EC2 Instance: %s", id)
					continue
				}

				r := ResourceInstance()
				d := r.Data(nil)
				d.SetId(id)
				d.Set("disable_api_stop", false)

				sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
			}
		}
		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing EC2 Instances for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping EC2 Instances for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping EC2 Instance sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}

func sweepInternetGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)

	defaultVPCID := ""
	describeVpcsInput := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("isDefault"),
				Values: aws.StringSlice([]string{"true"}),
			},
		},
	}

	describeVpcsOutput, err := conn.DescribeVpcsWithContext(ctx, describeVpcsInput)

	if err != nil {
		return fmt.Errorf("error describing VPCs: %w", err)
	}

	if describeVpcsOutput != nil && len(describeVpcsOutput.Vpcs) == 1 {
		defaultVPCID = aws.StringValue(describeVpcsOutput.Vpcs[0].VpcId)
	}

	input := &ec2.DescribeInternetGatewaysInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeInternetGatewaysPagesWithContext(ctx, input, func(page *ec2.DescribeInternetGatewaysOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, internetGateway := range page.InternetGateways {
			internetGatewayID := aws.StringValue(internetGateway.InternetGatewayId)
			isDefaultVPCInternetGateway := false

			for _, attachment := range internetGateway.Attachments {
				if aws.StringValue(attachment.VpcId) == defaultVPCID {
					isDefaultVPCInternetGateway = true
					break
				}
			}

			if isDefaultVPCInternetGateway {
				log.Printf("[DEBUG] Skipping Default VPC Internet Gateway: %s", internetGatewayID)
				continue
			}

			r := ResourceInternetGateway()
			d := r.Data(nil)
			d.SetId(internetGatewayID)
			if len(internetGateway.Attachments) > 0 {
				d.Set("vpc_id", internetGateway.Attachments[0].VpcId)
			}

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Internet Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Internet Gateways (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Internet Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepKeyPairs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeKeyPairsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	output, err := conn.DescribeKeyPairsWithContext(ctx, input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Key Pair sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Key Pairs (%s): %w", region, err)
	}

	for _, v := range output.KeyPairs {
		r := ResourceKeyPair()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.KeyName))

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Key Pairs (%s): %w", region, err)
	}

	return nil
}

func sweepLaunchTemplates(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeLaunchTemplatesInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeLaunchTemplatesPagesWithContext(ctx, input, func(page *ec2.DescribeLaunchTemplatesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.LaunchTemplates {
			r := ResourceLaunchTemplate()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.LaunchTemplateId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Launch Template sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Launch Templates (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Launch Templates (%s): %w", region, err)
	}

	return nil
}

func sweepNATGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	input := &ec2.DescribeNatGatewaysInput{}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeNatGatewaysPagesWithContext(ctx, input, func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.NatGateways {
			r := ResourceNATGateway()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.NatGatewayId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 NAT Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 NAT Gateways (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 NAT Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepNetworkACLs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	input := &ec2.DescribeNetworkAclsInput{}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeNetworkAclsPagesWithContext(ctx, input, func(page *ec2.DescribeNetworkAclsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.NetworkAcls {
			if aws.BoolValue(v.IsDefault) {
				continue
			}

			r := ResourceNetworkACL()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.NetworkAclId))

			var subnetIDs []string
			for _, v := range v.Associations {
				subnetIDs = append(subnetIDs, aws.StringValue(v.SubnetId))
			}
			d.Set("subnet_ids", subnetIDs)

			d.Set("vpc_id", v.VpcId)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Network ACL sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Network ACLs (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Network ACLs (%s): %w", region, err)
	}

	return nil
}

func sweepNetworkInterfaces(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeNetworkInterfacesInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeNetworkInterfacesPagesWithContext(ctx, input, func(page *ec2.DescribeNetworkInterfacesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.NetworkInterfaces {
			id := aws.StringValue(v.NetworkInterfaceId)

			if status := aws.StringValue(v.Status); status != ec2.NetworkInterfaceStatusAvailable {
				log.Printf("[INFO] Skipping EC2 Network Interface (%s): %s", status, id)
				continue
			}

			r := ResourceNetworkInterface()
			d := r.Data(nil)
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Network Interface sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Network Interfaces (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Network Interfaces (%s): %w", region, err)
	}

	return nil
}

func sweepNetworkInsightsPaths(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.DescribeNetworkInsightsPathsPagesWithContext(ctx, &ec2.DescribeNetworkInsightsPathsInput{}, func(page *ec2.DescribeNetworkInsightsPathsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, nip := range page.NetworkInsightsPaths {
			id := aws.StringValue(nip.NetworkInsightsPathId)

			r := ResourceNetworkInsightsPath()
			d := r.Data(nil)

			d.SetId(id)
			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error listing Network Insights Paths for %s: %w", region, err))
	}
	if err := sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping Network Insights Paths for %s: %w", region, err))
	}
	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Network Insights Path sweep for %s: %s", region, errs)
		return nil
	}
	return errs.ErrorOrNil()
}

func sweepPlacementGroups(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribePlacementGroupsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	output, err := conn.DescribePlacementGroupsWithContext(ctx, input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Placement Group sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Placement Groups (%s): %w", region, err)
	}

	for _, placementGroup := range output.PlacementGroups {
		r := ResourcePlacementGroup()
		d := r.Data(nil)
		d.SetId(aws.StringValue(placementGroup.GroupName))

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Placement Groups (%s): %w", region, err)
	}

	return nil
}

func sweepRouteTables(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.(*conns.AWSClient).EC2Conn(ctx)

	var sweeperErrs *multierror.Error

	input := &ec2.DescribeRouteTablesInput{}

	err = conn.DescribeRouteTablesPagesWithContext(ctx, input, func(page *ec2.DescribeRouteTablesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, routeTable := range page.RouteTables {
			if routeTable == nil {
				continue
			}

			id := aws.StringValue(routeTable.RouteTableId)
			isMainRouteTableAssociation := false

			for _, routeTableAssociation := range routeTable.Associations {
				if routeTableAssociation == nil {
					continue
				}

				if aws.BoolValue(routeTableAssociation.Main) {
					isMainRouteTableAssociation = true
					break
				}

				associationID := aws.StringValue(routeTableAssociation.RouteTableAssociationId)

				input := &ec2.DisassociateRouteTableInput{
					AssociationId: routeTableAssociation.RouteTableAssociationId,
				}

				log.Printf("[DEBUG] Deleting EC2 Route Table Association: %s", associationID)
				_, err := conn.DisassociateRouteTableWithContext(ctx, input)

				if err != nil {
					sweeperErr := fmt.Errorf("error deleting EC2 Route Table (%s) Association (%s): %w", id, associationID, err)
					log.Printf("[ERROR] %s", sweeperErr)
					sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
					continue
				}
			}

			if isMainRouteTableAssociation {
				for _, route := range routeTable.Routes {
					if route == nil {
						continue
					}

					if gatewayID := aws.StringValue(route.GatewayId); gatewayID == gatewayIDLocal || gatewayID == gatewayIDVPCLattice {
						continue
					}

					// Prevent deleting default VPC route for Internet Gateway
					// which some testing is still reliant on operating correctly
					if strings.HasPrefix(aws.StringValue(route.GatewayId), "igw-") && aws.StringValue(route.DestinationCidrBlock) == "0.0.0.0/0" {
						continue
					}

					input := &ec2.DeleteRouteInput{
						DestinationCidrBlock:     route.DestinationCidrBlock,
						DestinationIpv6CidrBlock: route.DestinationIpv6CidrBlock,
						RouteTableId:             routeTable.RouteTableId,
					}

					log.Printf("[DEBUG] Deleting EC2 Route Table (%s) Route", id)
					_, err := conn.DeleteRouteWithContext(ctx, input)

					if err != nil {
						sweeperErr := fmt.Errorf("error deleting EC2 Route Table (%s) Route: %w", id, err)
						log.Printf("[ERROR] %s", sweeperErr)
						sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
						continue
					}
				}

				continue
			}

			input := &ec2.DeleteRouteTableInput{
				RouteTableId: routeTable.RouteTableId,
			}

			log.Printf("[DEBUG] Deleting EC2 Route Table: %s", id)
			_, err := conn.DeleteRouteTableWithContext(ctx, input)

			if err != nil {
				sweeperErr := fmt.Errorf("error deleting EC2 Route Table (%s): %w", id, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Route Table sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil()
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing EC2 Route Tables: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepSecurityGroups(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)

	input := &ec2.DescribeSecurityGroupsInput{}

	// Delete all non-default EC2 Security Group Rules to prevent DependencyViolation errors
	err = conn.DescribeSecurityGroupsPagesWithContext(ctx, input, func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
		for _, sg := range page.SecurityGroups {
			if aws.StringValue(sg.GroupName) == "default" {
				log.Printf("[DEBUG] Skipping default EC2 Security Group: %s", aws.StringValue(sg.GroupId))
				continue
			}

			if sg.IpPermissions != nil {
				req := &ec2.RevokeSecurityGroupIngressInput{
					GroupId:       sg.GroupId,
					IpPermissions: sg.IpPermissions,
				}

				if _, err = conn.RevokeSecurityGroupIngressWithContext(ctx, req); err != nil {
					log.Printf("[ERROR] Error revoking ingress rule for Security Group (%s): %s", aws.StringValue(sg.GroupId), err)
				}
			}

			if sg.IpPermissionsEgress != nil {
				req := &ec2.RevokeSecurityGroupEgressInput{
					GroupId:       sg.GroupId,
					IpPermissions: sg.IpPermissionsEgress,
				}

				if _, err = conn.RevokeSecurityGroupEgressWithContext(ctx, req); err != nil {
					log.Printf("[ERROR] Error revoking egress rule for Security Group (%s): %s", aws.StringValue(sg.GroupId), err)
				}
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Security Group sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error retrieving EC2 Security Groups: %w", err)
	}

	err = conn.DescribeSecurityGroupsPagesWithContext(ctx, input, func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
		for _, sg := range page.SecurityGroups {
			if aws.StringValue(sg.GroupName) == "default" {
				log.Printf("[DEBUG] Skipping default EC2 Security Group: %s", aws.StringValue(sg.GroupId))
				continue
			}

			input := &ec2.DeleteSecurityGroupInput{
				GroupId: sg.GroupId,
			}

			// Handle EC2 eventual consistency
			err := retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
				_, err := conn.DeleteSecurityGroupWithContext(ctx, input)

				if tfawserr.ErrCodeEquals(err, "DependencyViolation") {
					return retry.RetryableError(err)
				}
				if err != nil {
					return retry.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				log.Printf("[ERROR] Error deleting Security Group (%s): %s", aws.StringValue(sg.GroupId), err)
			}
		}

		return !lastPage
	})

	if err != nil {
		return fmt.Errorf("Error retrieving EC2 Security Groups: %w", err)
	}

	return nil
}

func sweepSpotFleetRequests(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.DescribeSpotFleetRequestsPagesWithContext(ctx, &ec2.DescribeSpotFleetRequestsInput{}, func(page *ec2.DescribeSpotFleetRequestsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		if len(page.SpotFleetRequestConfigs) == 0 {
			log.Print("[DEBUG] No Spot Fleet Requests to sweep")
			return false
		}

		for _, config := range page.SpotFleetRequestConfigs {
			id := aws.StringValue(config.SpotFleetRequestId)

			r := ResourceSpotFleetRequest()
			d := r.Data(nil)
			d.SetId(id)
			d.Set("terminate_instances_with_expiration", true)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing EC2 Spot Fleet Requests for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping EC2 Spot Fleet Requests for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping EC2 Spot Fleet Requests sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}

func sweepSpotInstanceRequests(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.DescribeSpotInstanceRequestsPagesWithContext(ctx, &ec2.DescribeSpotInstanceRequestsInput{}, func(page *ec2.DescribeSpotInstanceRequestsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		if len(page.SpotInstanceRequests) == 0 {
			log.Print("[DEBUG] No Spot Instance Requests to sweep")
			return false
		}

		for _, config := range page.SpotInstanceRequests {
			id := aws.StringValue(config.SpotInstanceRequestId)

			r := ResourceSpotInstanceRequest()
			d := r.Data(nil)
			d.SetId(id)
			d.Set("spot_instance_id", config.InstanceId)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing EC2 Spot Instance Requests for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestratorWithContext(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping EC2 Spot Instance Requests for %s: %w", region, err))
	}

	if sweep.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping EC2 Spot Instance Requests sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}

func sweepSubnets(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeSubnetsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeSubnetsPagesWithContext(ctx, input, func(page *ec2.DescribeSubnetsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Subnets {
			// Skip default subnets.
			if aws.BoolValue(v.DefaultForAz) {
				continue
			}

			r := ResourceSubnet()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.SubnetId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Subnet sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Subnets (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Subnets (%s): %w", region, err)
	}

	return nil
}

func sweepTransitGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeTransitGatewaysInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeTransitGatewaysPagesWithContext(ctx, input, func(page *ec2.DescribeTransitGatewaysOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.TransitGateways {
			if aws.StringValue(v.State) == ec2.TransitGatewayStateDeleted {
				continue
			}

			r := ResourceTransitGateway()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.TransitGatewayId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Transit Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Transit Gateways (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Transit Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepTransitGatewayConnectPeers(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeTransitGatewayConnectPeersInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeTransitGatewayConnectPeersPagesWithContext(ctx, input, func(page *ec2.DescribeTransitGatewayConnectPeersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.TransitGatewayConnectPeers {
			if aws.StringValue(v.State) == ec2.TransitGatewayConnectPeerStateDeleted {
				continue
			}

			r := ResourceTransitGatewayConnectPeer()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.TransitGatewayConnectPeerId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Transit Gateway Connect Peer sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Transit Gateway Connect Peers (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Transit Gateway Connect Peers (%s): %w", region, err)
	}

	return nil
}

func sweepTransitGatewayConnects(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeTransitGatewayConnectsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeTransitGatewayConnectsPagesWithContext(ctx, input, func(page *ec2.DescribeTransitGatewayConnectsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.TransitGatewayConnects {
			if aws.StringValue(v.State) == ec2.TransitGatewayAttachmentStateDeleted {
				continue
			}

			r := ResourceTransitGatewayConnect()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.TransitGatewayAttachmentId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Transit Gateway Connect sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Transit Gateway Connects (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Transit Gateway Connects (%s): %w", region, err)
	}

	return nil
}

func sweepTransitGatewayMulticastDomains(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeTransitGatewayMulticastDomainsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeTransitGatewayMulticastDomainsPagesWithContext(ctx, input, func(page *ec2.DescribeTransitGatewayMulticastDomainsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.TransitGatewayMulticastDomains {
			if aws.StringValue(v.State) == ec2.TransitGatewayMulticastDomainStateDeleted {
				continue
			}

			r := ResourceTransitGatewayMulticastDomain()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.TransitGatewayMulticastDomainId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Transit Gateway Multicast Domain sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Transit Gateway Multicast Domains (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Transit Gateway Multicast Domains (%s): %w", region, err)
	}

	return nil
}

func sweepTransitGatewayPeeringAttachments(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeTransitGatewayPeeringAttachmentsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeTransitGatewayPeeringAttachmentsPagesWithContext(ctx, input, func(page *ec2.DescribeTransitGatewayPeeringAttachmentsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.TransitGatewayPeeringAttachments {
			if aws.StringValue(v.State) == ec2.TransitGatewayAttachmentStateDeleted {
				continue
			}

			r := ResourceTransitGatewayPeeringAttachment()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.TransitGatewayAttachmentId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Transit Gateway Peering Attachment sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Transit Gateway Peering Attachments (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Transit Gateway Peering Attachments (%s): %w", region, err)
	}

	return nil
}

func sweepTransitGatewayVPCAttachments(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeTransitGatewayVpcAttachmentsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeTransitGatewayVpcAttachmentsPagesWithContext(ctx, input, func(page *ec2.DescribeTransitGatewayVpcAttachmentsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.TransitGatewayVpcAttachments {
			if aws.StringValue(v.State) == ec2.TransitGatewayAttachmentStateDeleted {
				continue
			}

			r := ResourceTransitGatewayVPCAttachment()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.TransitGatewayAttachmentId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Transit Gateway VPC Attachment sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Transit Gateway VPC Attachments (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Transit Gateway VPC Attachments (%s): %w", region, err)
	}

	return nil
}

func sweepVPCDHCPOptions(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	input := &ec2.DescribeDhcpOptionsInput{}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeDhcpOptionsPagesWithContext(ctx, input, func(page *ec2.DescribeDhcpOptionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.DhcpOptions {
			// Skip the default DHCP Options.
			var defaultDomainNameFound, defaultDomainNameServersFound bool

			for _, v := range v.DhcpConfigurations {
				if aws.StringValue(v.Key) == "domain-name" {
					if len(v.Values) != 1 || v.Values[0] == nil {
						continue
					}

					if aws.StringValue(v.Values[0].Value) == RegionalPrivateDNSSuffix(region) {
						defaultDomainNameFound = true
					}
				} else if aws.StringValue(v.Key) == "domain-name-servers" {
					if len(v.Values) != 1 || v.Values[0] == nil {
						continue
					}

					if aws.StringValue(v.Values[0].Value) == "AmazonProvidedDNS" {
						defaultDomainNameServersFound = true
					}
				}
			}

			if defaultDomainNameFound && defaultDomainNameServersFound {
				continue
			}

			r := ResourceVPCDHCPOptions()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.DhcpOptionsId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 DHCP Options Set sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 DHCP Options Sets (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 DHCP Options Sets (%s): %w", region, err)
	}

	return nil
}

func sweepVPCEndpointServices(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeVpcEndpointServiceConfigurationsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeVpcEndpointServiceConfigurationsPagesWithContext(ctx, input, func(page *ec2.DescribeVpcEndpointServiceConfigurationsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.ServiceConfigurations {
			if aws.StringValue(v.ServiceState) == ec2.ServiceStateDeleted {
				continue
			}

			r := ResourceVPCEndpointService()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.ServiceId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 VPC Endpoint Service sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 VPC Endpoint Services (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 VPC Endpoint Services (%s): %w", region, err)
	}

	return nil
}

func sweepVPCEndpoints(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeVpcEndpointsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeVpcEndpointsPagesWithContext(ctx, input, func(page *ec2.DescribeVpcEndpointsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.VpcEndpoints {
			if aws.StringValue(v.State) == vpcEndpointStateDeleted {
				continue
			}

			r := ResourceVPCEndpoint()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.VpcEndpointId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 VPC Endpoint sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 VPC Endpoints (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 VPC Endpoints (%s): %w", region, err)
	}

	return nil
}

func sweepVPCPeeringConnections(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	input := &ec2.DescribeVpcPeeringConnectionsInput{}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeVpcPeeringConnectionsPagesWithContext(ctx, input, func(page *ec2.DescribeVpcPeeringConnectionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.VpcPeeringConnections {
			r := ResourceVPCPeeringConnection()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.VpcPeeringConnectionId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 VPC Peering Connection sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 VPC Peering Connections (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 VPC Peering Connections (%s): %w", region, err)
	}

	return nil
}

func sweepVPCs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeVpcsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeVpcsPagesWithContext(ctx, input, func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Vpcs {
			// Skip default VPCs.
			if aws.BoolValue(v.IsDefault) {
				continue
			}

			r := ResourceVPC()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.VpcId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 VPC sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 VPCs (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 VPCs (%s): %w", region, err)
	}

	return nil
}

func sweepVPNConnections(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeVpnConnectionsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	output, err := conn.DescribeVpnConnectionsWithContext(ctx, input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 VPN Connection sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 VPN Connections (%s): %w", region, err)
	}

	for _, v := range output.VpnConnections {
		if aws.StringValue(v.State) == ec2.VpnStateDeleted {
			continue
		}

		r := ResourceVPNConnection()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.VpnConnectionId))

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 VPN Connections (%s): %w", region, err)
	}

	return nil
}

func sweepVPNGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeVpnGatewaysInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	output, err := conn.DescribeVpnGatewaysWithContext(ctx, input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 VPN Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 VPN Gateways (%s): %w", region, err)
	}

	for _, v := range output.VpnGateways {
		if aws.StringValue(v.State) == ec2.VpnStateDeleted {
			continue
		}

		r := ResourceVPNGateway()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.VpnGatewayId))

		for _, v := range v.VpcAttachments {
			if aws.StringValue(v.State) != ec2.AttachmentStatusDetached {
				d.Set("vpc_id", v.VpcId)

				break
			}
		}

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 VPN Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepCustomerGateways(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeCustomerGatewaysInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	output, err := conn.DescribeCustomerGatewaysWithContext(ctx, input)

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Customer Gateway sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 Customer Gateways (%s): %w", region, err)
	}

	for _, v := range output.CustomerGateways {
		if aws.StringValue(v.State) == CustomerGatewayStateDeleted {
			continue
		}

		r := ResourceCustomerGateway()
		d := r.Data(nil)
		d.SetId(aws.StringValue(v.CustomerGatewayId))

		sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 Customer Gateways (%s): %w", region, err)
	}

	return nil
}

func sweepIPAMs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeIpamsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeIpamsPagesWithContext(ctx, input, func(page *ec2.DescribeIpamsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Ipams {
			r := ResourceIPAM()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.IpamId))
			d.Set("cascade", true)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping IPAM sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing IPAMs (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping IPAMs (%s): %w", region, err)
	}

	return nil
}

func sweepIPAMResourceDiscoveries(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeIpamResourceDiscoveriesInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeIpamResourceDiscoveriesPagesWithContext(ctx, input, func(page *ec2.DescribeIpamResourceDiscoveriesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.IpamResourceDiscoveries {
			// do not attempt to delete default resource created by each ipam
			if !aws.BoolValue(v.IsDefault) {
				r := ResourceIPAMResourceDiscovery()
				d := r.Data(nil)
				d.SetId(aws.StringValue(v.IpamResourceDiscoveryId))

				sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping IPAM Resource Discovery sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing IPAM Resource Discoveries (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Resource Discoveries (%s): %w", region, err)
	}

	return nil
}

func sweepAMIs(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	input := &ec2.DescribeImagesInput{
		Owners: aws.StringSlice([]string{"self"}),
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeImagesPagesWithContext(ctx, input, func(page *ec2.DescribeImagesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Images {
			r := ResourceAMI()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.ImageId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping AMI sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing AMIs (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping AMIs (%s): %w", region, err)
	}

	return nil
}

func sweepNetworkPerformanceMetricSubscriptions(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EC2Conn(ctx)
	input := &ec2.DescribeAwsNetworkPerformanceMetricSubscriptionsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeAwsNetworkPerformanceMetricSubscriptionsPagesWithContext(ctx, input, func(page *ec2.DescribeAwsNetworkPerformanceMetricSubscriptionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Subscriptions {
			r := ResourceNetworkPerformanceMetricSubscription()
			id := NetworkPerformanceMetricSubscriptionCreateResourceID(aws.StringValue(v.Source), aws.StringValue(v.Destination), aws.StringValue(v.Metric), aws.StringValue(v.Statistic))
			d := r.Data(nil)
			d.SetId(id)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 AWS Network Performance Metric Subscription sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EC2 AWS Network Performance Metric Subscriptions (%s): %w", region, err)
	}

	err = sweep.SweepOrchestratorWithContext(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EC2 AWS Network Performance Metric Subscriptions (%s): %w", region, err)
	}

	return nil
}
