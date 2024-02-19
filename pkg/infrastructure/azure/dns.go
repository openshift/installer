package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
)

// Create DNS entries for azure.
func createAzureDNSEntries(ctx context.Context, in clusterapi.InfraReadyInput) error {
	clusterName := in.InfraID
	private := in.InstallConfig.Config.Publish == types.InternalPublishingStrategy
	resourceGroup := in.InstallConfig.Config.Azure.ResourceGroupName
	zoneName := in.InstallConfig.Config.Azure.DefaultMachinePlatform.Zones
	subscriptionId := in.InstallConfig.Azure.Credentials.SubscriptionID
	apiExternalName := fmt.Sprintf("api.%s", clusterName)
	apiExternalNameV6 := fmt.Sprintf("v6-api.%s", clusterName)
	// TODO: set value from either manifests or getting from client.
	var azureTags map[string]*string
	var ipv4Elb string
	var ipv4Ilb string
	var ipv6Elb string
	var ipv6Ilb string

	useIPv6 := false
	for _, network := range in.InstallConfig.Config.Networking.ServiceNetwork {
		if network.IP.To4() == nil {
			useIPv6 = true
		}
	}

	type recordList struct {
		Name       string
		RecordType armdns.RecordType
		RecordSet  armdns.RecordSet
	}
	records := []recordList{}
	ttl := int64(300)
	if !useIPv6 {
		records = append(records, recordList{
			Name:       "api-int",
			RecordType: armdns.RecordTypeA,
			RecordSet: armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					ARecords: []*armdns.ARecord{
						{
							IPv4Address: &ipv4Ilb,
						},
					},
					TTL:      &ttl,
					Metadata: azureTags,
				},
			},
		}, recordList{
			Name:       "api",
			RecordType: armdns.RecordTypeA,
			RecordSet: armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					ARecords: []*armdns.ARecord{
						{
							IPv4Address: &ipv4Elb,
						},
					},
					TTL:      &ttl,
					Metadata: azureTags,
				},
			},
		})
	} else {
		records = append(records, recordList{
			Name:       "api-int",
			RecordType: armdns.RecordTypeAAAA,
			RecordSet: armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					AaaaRecords: []*armdns.AaaaRecord{
						{
							IPv6Address: &ipv6Ilb,
						},
					},
					TTL:      &ttl,
					Metadata: azureTags,
				},
			},
		}, recordList{
			Name:       "api",
			RecordType: armdns.RecordTypeAAAA,
			RecordSet: armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					AaaaRecords: []*armdns.AaaaRecord{
						{
							IPv6Address: &ipv6Elb,
						},
					},
					TTL:      &ttl,
					Metadata: azureTags,
				},
			},
		})
	}

	if !private {
		cnameRecordName := apiExternalName
		if useIPv6 {
			cnameRecordName = apiExternalNameV6
		}
		records = append(records, recordList{
			Name:       "api-int",
			RecordType: armdns.RecordTypeCNAME,
			RecordSet: armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					CnameRecord: &armdns.CnameRecord{
						Cname: &cnameRecordName,
					},
					TTL:      &ttl,
					Metadata: azureTags,
				},
			},
		})
	}

	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return err
	}
	tokenCreds, err := azidentity.NewClientSecretCredential(session.Credentials.TenantID, session.Credentials.ClientID, session.Credentials.ClientSecret, nil)
	if err != nil {
		return err
	}
	recordSetClient, err := armdns.NewRecordSetsClient(subscriptionId, tokenCreds, nil)
	if err != nil {
		return err
	}

	for _, zone := range zoneName {
		for _, record := range records {
			_, err = recordSetClient.CreateOrUpdate(ctx, resourceGroup, zone, record.Name, record.RecordType, record.RecordSet, nil)
			if err != nil {
				return err
			}
		}
	}
	return err
}
