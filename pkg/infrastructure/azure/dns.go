package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/sirupsen/logrus"

	aztypes "github.com/openshift/installer/pkg/types/azure"
)

type CreateDNSEntriesInput struct {
	SubscriptionID              string
	ResourceGroupName           string
	BaseDomainResourceGroupName string
	BaseDomain                  string
	ClusterDomain               string
	ClusterName                 string
	Region                      string
	InfraID                     string
	Private                     bool
	UseIPv6                     bool
	Tags                        map[string]*string
	CloudName                   aztypes.CloudEnvironment
	TokenCredential             azcore.TokenCredential
	CloudConfiguration          cloud.Configuration
}

// CreateDNSEntries creates DNS entries in public and private DNS zones.
func CreateDNSEntries(ctx context.Context, in *CreateDNSEntriesInput) error {
	var (
		subscriptionID          = in.SubscriptionID
		resourceGroup           = in.ResourceGroupName
		clusterDomain           = in.ClusterDomain
		baseDomainResourceGroup = in.BaseDomainResourceGroupName
		baseDomain              = in.BaseDomain
		private                 = in.Private
		//region                  = in.Region
		tokenCreds         = in.TokenCredential
		useIPv6            = in.UseIPv6
		cloudConfiguration = in.CloudConfiguration
		tags               = in.Tags

		apiExternalName       = fmt.Sprintf("api.%s", in.ClusterName)
		apiExternalNameV6     = fmt.Sprintf("v6-api.%s", in.ClusterName)
		publicIPv4AddressName = fmt.Sprintf("%s-pip-v4", in.InfraID)
		publicIPv6AddressName = fmt.Sprintf("%s-pip-v6", in.InfraID)
		//virtualNetworkName    = fmt.Sprintf("%s-vnet", in.InfraID)
		//virtualNetworkLinkName = fmt.Sprintf("%s-network-link", in.InfraID)
		loadBalancerName = fmt.Sprintf("%s-internal", in.InfraID)
	)

	// Get the internal load balancer
	loadBalancersClient, err := armnetwork.NewLoadBalancersClient(subscriptionID, tokenCreds, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudConfiguration,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get load balancers client: %v", err)
	}
	loadBalancer, err := loadBalancersClient.Get(ctx, resourceGroup, loadBalancerName, nil)
	if err != nil {
		return fmt.Errorf("failed to get internal load balancer %s: %v", loadBalancerName, err)
	}

	// XXX: What do we do with multiple frontend IP
	// configurations/addresses ?
	lbIPv4Address := loadBalancer.Properties.FrontendIPConfigurations[0].Properties.PrivateIPAddress

	// List for private DNS records
	type privateRecordList struct {
		Name       string
		RecordType armprivatedns.RecordType
		RecordSet  armprivatedns.RecordSet
	}

	// Private A records for "api" and "api-int"
	privateRecords := []privateRecordList{
		privateRecordList{
			Name:       "api",
			RecordType: armprivatedns.RecordTypeA,
			RecordSet: armprivatedns.RecordSet{
				Properties: &armprivatedns.RecordSetProperties{
					ARecords: []*armprivatedns.ARecord{
						{IPv4Address: lbIPv4Address},
					},
					TTL:      to.Ptr(int64(300)),
					Metadata: tags,
				},
			},
		}, privateRecordList{
			Name:       "api-int",
			RecordType: armprivatedns.RecordTypeA,
			RecordSet: armprivatedns.RecordSet{
				Properties: &armprivatedns.RecordSetProperties{
					ARecords: []*armprivatedns.ARecord{
						{IPv4Address: lbIPv4Address},
					},
					TTL:      to.Ptr(int64(300)),
					Metadata: tags,
				},
			},
		},
	}
	// Private AAAA records for "api" and "api-int"
	if useIPv6 && len(loadBalancer.Properties.FrontendIPConfigurations) > 1 {
		lbIPv6Address := loadBalancer.Properties.FrontendIPConfigurations[1].Properties.PrivateIPAddress
		privateRecords = append(privateRecords, privateRecordList{
			Name:       "api",
			RecordType: armprivatedns.RecordTypeAAAA,
			RecordSet: armprivatedns.RecordSet{
				Properties: &armprivatedns.RecordSetProperties{
					AaaaRecords: []*armprivatedns.AaaaRecord{
						{IPv6Address: lbIPv6Address},
					},
					TTL:      to.Ptr(int64(300)),
					Metadata: tags,
				},
			},
		}, privateRecordList{
			Name:       "api-int",
			RecordType: armprivatedns.RecordTypeAAAA,
			RecordSet: armprivatedns.RecordSet{
				Properties: &armprivatedns.RecordSetProperties{
					AaaaRecords: []*armprivatedns.AaaaRecord{
						{IPv6Address: lbIPv6Address},
					},
					TTL:      to.Ptr(int64(300)),
					Metadata: tags,
				},
			},
		})
	}

	// Create private DNS zone
	/*
		privateZonesClient, err := armprivatedns.NewPrivateZonesClient(subscriptionID, tokenCreds, &arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloudConfiguration,
			},
		})
		if err != nil {
			return err
		}
		privateZonesPollerResponse, err := privateZonesClient.BeginCreateOrUpdate(ctx, resourceGroup, privateZoneName,
			armprivatedns.PrivateZone{
				Location: to.Ptr(region),
			}, nil)
		if err != nil {
			return err
		}
		_, err = privateZonesPollerResponse.PollUntilDone(ctx, nil)
		if err != nil {
			return nil
		}
	*/

	// Create records in private DNS zone
	privateRecordSetsClient, err := armprivatedns.NewRecordSetsClient(subscriptionID, tokenCreds, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudConfiguration,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get recordset client: %v", err)
	}
	for _, privateRecord := range privateRecords {
		logrus.Debugf("DNS: creating private record %s with value %s in zone %s",
			privateRecord.Name,
			privateRecord.RecordSet.Properties.ARecords[0].IPv4Address,
			clusterDomain,
		)
		_, err = privateRecordSetsClient.CreateOrUpdate(ctx, resourceGroup, clusterDomain, privateRecord.RecordType, privateRecord.Name, privateRecord.RecordSet, nil)
		if err != nil {
			return fmt.Errorf("failed to create private DNS record %s in zone %s: %v", privateRecord.Name, clusterDomain, err)
		}
	}

	// Get virtual network ID
	/*
		virtualNetworksClient, err := armnetwork.NewVirtualNetworksClient(subscriptionID, tokenCreds, &arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloudConfiguration,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get virtual networks client: %v", err)
		}
		virtualNetwork, err := virtualNetworksClient.Get(ctx, resourceGroup, virtualNetworkName, nil)
		if err != nil {
			return fmt.Errorf("failed to get virtual network %s: %v", virtualNetworkName, err)
		}
	*/

	// Create the private link
	/*
		virtualNetworkLinksClient, err := armprivatedns.NewVirtualNetworkLinksClient(subscriptionID, tokenCreds, &arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloudConfiguration,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get virtual network links client: %v", err)
		}
		virtualNetworkLinkPollerResponse, err := virtualNetworkLinksClient.BeginCreateOrUpdate(ctx, resourceGroup, clusterDomain,
			virtualNetworkLinkName, armprivatedns.VirtualNetworkLink{
				Location: to.Ptr(region),
				Properties: &armprivatedns.VirtualNetworkLinkProperties{
					VirtualNetwork: &armprivatedns.SubResource{
						ID: virtualNetwork.ID,
					},
				},
			}, nil)
		if err != nil {
			return fmt.Errorf("failed to create virtual network link %s: %v", virtualNetworkLinkName, err)
		}
		_, err = virtualNetworkLinkPollerResponse.PollUntilDone(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to finish creating virtual network link %s: %v", virtualNetworkLinkName, err)
		}
	*/

	// Return if we don't need to create any public records
	if private {
		return nil
	}

	// XXX: return for now - Need to handle waiting for the public IP
	// address to become available
	return nil

	// List for public DNS records
	type publicRecordList struct {
		Name       string
		RecordType armdns.RecordType
		RecordSet  armdns.RecordSet
	}

	// Get public IP addresses so we can set up the CNAME records
	publicIPAddressClient, err := armnetwork.NewPublicIPAddressesClient(subscriptionID, tokenCreds, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudConfiguration,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get public ip address client: %v", err)
	}
	publicIPv4Address, err := publicIPAddressClient.Get(ctx, resourceGroup, publicIPv4AddressName, nil)
	if err != nil {
		return fmt.Errorf("failed to get public ip address %s: %v", publicIPv4AddressName, err)
	}
	publicFqdnV4 := publicIPv4Address.Properties.DNSSettings.Fqdn

	// Public CNAME records for "api"
	publicRecords := []publicRecordList{
		publicRecordList{
			Name:       apiExternalName,
			RecordType: armdns.RecordTypeCNAME,
			RecordSet: armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					CnameRecord: &armdns.CnameRecord{
						Cname: publicFqdnV4,
					},
					TTL:      to.Ptr(int64(300)),
					Metadata: tags,
				},
			},
		},
	}
	if useIPv6 {
		publicIPv6Address, err := publicIPAddressClient.Get(ctx, resourceGroup, publicIPv6AddressName, nil)
		if err != nil {
			return nil
		}
		publicFqdnV6 := publicIPv6Address.Properties.DNSSettings.Fqdn

		publicRecords = append(publicRecords, publicRecordList{
			Name:       apiExternalNameV6,
			RecordType: armdns.RecordTypeCNAME,
			RecordSet: armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					CnameRecord: &armdns.CnameRecord{
						Cname: publicFqdnV6,
					},
					TTL:      to.Ptr(int64(300)),
					Metadata: tags,
				},
			},
		})
	}

	// Create records in public DNS zone
	recordSetClient, err := armdns.NewRecordSetsClient(subscriptionID, tokenCreds, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudConfiguration,
		},
	})
	if err != nil {
		return err
	}
	for _, publicRecord := range publicRecords {
		logrus.Debugf("DNS: creating public record %s with value %s in zone %s",
			publicRecord.Name,
			publicRecord.RecordSet.Properties.CnameRecord.Cname,
			clusterDomain,
		)
		_, err = recordSetClient.CreateOrUpdate(ctx, baseDomainResourceGroup, baseDomain, publicRecord.Name, publicRecord.RecordType, publicRecord.RecordSet, nil)
		if err != nil {
			return fmt.Errorf("failed to create public DNS record %s in zone %s: %v", publicRecord.Name, clusterDomain, err)
		}
	}

	return err
}
