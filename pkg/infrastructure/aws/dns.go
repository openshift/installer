package aws

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
)

const errSharedCredsLoad = "SharedCredsLoad"

var cnameRegions = sets.New[string]("us-gov-west-1", "us-gov-east-1")

type dnsInputOptions struct {
	infraID           string
	region            string
	baseDomain        string
	clusterDomain     string
	vpcID             string
	lbExternalZoneID  string
	lbExternalZoneDNS string
	lbInternalZoneID  string
	lbInternalZoneDNS string
	internalZone      string
	isPrivateCluster  bool
	tags              map[string]string
}

func createDNSResources(ctx context.Context, logger logrus.FieldLogger, route53Client route53iface.Route53API, assumedRoleClient route53iface.Route53API, input *dnsInputOptions) error {
	apiName := fmt.Sprintf("api.%s", input.clusterDomain)
	apiIntName := fmt.Sprintf("api-int.%s", input.clusterDomain)
	useCNAME := cnameRegions.Has(input.region)

	if !input.isPrivateCluster {
		publicZone, err := existingHostedZone(ctx, route53Client, input.baseDomain, false)
		if err != nil {
			return fmt.Errorf("failed to find public zone (%s): %w", input.baseDomain, err)
		}
		zoneID := cleanZoneID(publicZone.Id)
		logger.WithFields(logrus.Fields{
			"domain": input.baseDomain,
			"id":     zoneID,
		}).Infoln("Found existing public zone")

		// Create API record in public zone
		_, err = createRecord(ctx, route53Client, zoneID, apiName, input.lbExternalZoneDNS, input.lbExternalZoneID, useCNAME)
		if err != nil {
			return fmt.Errorf("failed to create api record (%s) in public zone: %w", apiName, err)
		}
		logger.Infoln("Created api DNS record for public zone")
	}

	privateZoneID := cleanZoneID(aws.String(input.internalZone))
	if len(privateZoneID) == 0 {
		privateZone, err := ensurePrivateZone(ctx, logger, route53Client, input)
		if err != nil {
			return err
		}
		privateZoneID = cleanZoneID(privateZone.Id)
	}

	// Note: assumedRoleClient used below is either the standard route53 client
	// or a client with a custom role if one was supplied for a phz belonging
	// to a different account than the rest of the cluster resources.

	// Create API record in private zone
	_, err := createRecord(ctx, assumedRoleClient, privateZoneID, apiName, input.lbInternalZoneDNS, input.lbInternalZoneID, useCNAME)
	if err != nil {
		return fmt.Errorf("failed to create api record (%s) in private zone: %w", apiName, err)
	}
	logger.Infoln("Created api DNS record for private zone")

	// Create API-int record in privat zone
	_, err = createRecord(ctx, assumedRoleClient, privateZoneID, apiIntName, input.lbInternalZoneDNS, input.lbInternalZoneID, useCNAME)
	if err != nil {
		return fmt.Errorf("failed to create api-int record (%s) in private zone: %w", apiIntName, err)
	}
	logger.Infoln("Created api-int DNS record for private zone")

	return nil
}

func ensurePrivateZone(ctx context.Context, logger logrus.FieldLogger, client route53iface.Route53API, input *dnsInputOptions) (*route53.HostedZone, error) {
	createdOrFoundMsg := "Found existing private hosted zone"
	privateZone, err := existingHostedZone(ctx, client, input.clusterDomain, true)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		createdOrFoundMsg = "Created private hosted zone"
		privateZone, err = createHostedZone(ctx, client, input.clusterDomain, input.vpcID, input.region, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create private hosted zone: %w", err)
		}
	}
	l := logger.WithField("id", cleanZoneID(privateZone.Id))
	l.Infoln(createdOrFoundMsg)

	// Tag the hosted zone
	tags := mergeTags(input.tags, map[string]string{"Name": fmt.Sprintf("%s-int", input.infraID)})
	_, err = client.ChangeTagsForResourceWithContext(ctx, &route53.ChangeTagsForResourceInput{
		ResourceType: aws.String("hostedzone"),
		ResourceId:   privateZone.Id,
		AddTags:      r53Tags(tags),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to tag private hosted zone: %w", err)
	}
	l.Infoln("Tagged private hosted zone")

	// Set SOA minimum TTL
	recordSet, err := existingRecordSet(ctx, client, privateZone.Id, input.clusterDomain, "SOA")
	if err != nil {
		return nil, fmt.Errorf("failed to find SOA record set for private zone: %w", err)
	}
	if len(recordSet.ResourceRecords) == 0 || recordSet.ResourceRecords[0] == nil || recordSet.ResourceRecords[0].Value == nil {
		return nil, fmt.Errorf("failed to find SOA record for private zone")
	}
	record := recordSet.ResourceRecords[0]
	fields := strings.Split(aws.StringValue(record.Value), " ")
	if len(fields) != 7 {
		return nil, fmt.Errorf("SOA record value has %d fields, expected 7", len(fields))
	}
	fields[6] = "60"
	record.Value = aws.String(strings.Join(fields, " "))
	_, err = client.ChangeResourceRecordSetsWithContext(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: privateZone.Id,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String("UPSERT"),
					ResourceRecordSet: recordSet,
				},
			},
		},
	},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set SOA TTL to minimum: %w", err)
	}

	return privateZone, nil
}

func existingHostedZone(ctx context.Context, client route53iface.Route53API, name string, isPrivate bool) (*route53.HostedZone, error) {
	var hz *route53.HostedZone
	if err := retryWithBackoff(
		ctx,
		func(ctx context.Context) error {
			return client.ListHostedZonesPagesWithContext(
				ctx,
				&route53.ListHostedZonesInput{},
				func(res *route53.ListHostedZonesOutput, lastPage bool) bool {
					for i, zone := range res.HostedZones {
						if zone.Config != nil && aws.BoolValue(zone.Config.PrivateZone) == isPrivate && strings.TrimSuffix(aws.StringValue(zone.Name), ".") == strings.TrimSuffix(name, ".") {
							hz = res.HostedZones[i]
							return false
						}
					}
					return !lastPage
				},
			)
		},
	); err != nil {
		return nil, fmt.Errorf("failed to list hosted zones: %w", err)
	}
	if hz == nil {
		return nil, errNotFound
	}

	return hz, nil
}

func createHostedZone(ctx context.Context, client route53iface.Route53API, name, vpcID, region string, isPrivate bool) (*route53.HostedZone, error) {
	var res *route53.CreateHostedZoneOutput
	err := retryWithBackoff(
		ctx,
		func(ctx context.Context) error {
			var err error
			callRef := fmt.Sprintf("%d", time.Now().Unix())
			res, err = client.CreateHostedZoneWithContext(ctx, &route53.CreateHostedZoneInput{
				CallerReference: aws.String(callRef),
				Name:            aws.String(name),
				HostedZoneConfig: &route53.HostedZoneConfig{
					PrivateZone: aws.Bool(isPrivate),
					Comment:     aws.String(defaultDescription),
				},
				VPC: &route53.VPC{
					VPCId:     aws.String(vpcID),
					VPCRegion: aws.String(region),
				},
			})
			return err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create hosted zone: %w", err)
	}
	if res == nil {
		return nil, fmt.Errorf("unexpected output from hosted zone creation")
	}

	return res.HostedZone, nil
}

func createRecord(ctx context.Context, client route53iface.Route53API, zoneID string, name string, dnsName string, aliasZoneID string, useCNAME bool) (*route53.ChangeInfo, error) {
	recordSet := &route53.ResourceRecordSet{
		Name: aws.String(cleanRecordName(name)),
	}
	if useCNAME {
		recordSet.SetType("CNAME")
		recordSet.SetTTL(10)
		recordSet.SetResourceRecords([]*route53.ResourceRecord{
			{Value: aws.String(dnsName)},
		})
	} else {
		recordSet.SetType("A")
		recordSet.SetAliasTarget(&route53.AliasTarget{
			DNSName:              aws.String(dnsName),
			HostedZoneId:         aws.String(aliasZoneID),
			EvaluateTargetHealth: aws.Bool(false),
		})
	}
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String("UPSERT"),
					ResourceRecordSet: recordSet,
				},
			},
		},
	}
	res, err := client.ChangeResourceRecordSetsWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return res.ChangeInfo, nil
}

func existingRecordSet(ctx context.Context, client route53iface.Route53API, zoneID *string, recordName string, recordType string) (*route53.ResourceRecordSet, error) {
	name := fqdn(strings.ToLower(recordName))
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    zoneID,
		StartRecordName: aws.String(name),
		StartRecordType: aws.String(recordType),
		MaxItems:        aws.String("1"),
	}
	res, err := client.ListResourceRecordSetsWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	for _, rs := range res.ResourceRecordSets {
		resName := strings.ToLower(cleanRecordName(aws.StringValue(rs.Name)))
		resType := strings.ToUpper(aws.StringValue(rs.Type))
		if resName == name && resType == recordType {
			return rs, nil
		}
	}
	return nil, errNotFound
}

func cleanZoneID(zoneID *string) string {
	return strings.TrimPrefix(aws.StringValue(zoneID), "/hostedzone/")
}

func cleanRecordName(name string) string {
	s, err := strconv.Unquote(`"` + name + `"`)
	if err != nil {
		return name
	}
	return s
}

func fqdn(name string) string {
	n := len(name)
	if n == 0 || name[n-1] == '.' {
		return name
	}
	return name + "."
}

func retryWithBackoff(ctx context.Context, fn func(ctx context.Context) error) error {
	return wait.ExponentialBackoffWithContext(
		ctx,
		defaultBackoff,
		func(ctx context.Context) (bool, error) {
			if err := fn(ctx); err != nil {
				var awsErr awserr.Error
				if errors.As(err, &awsErr) && awsErr.Code() == errSharedCredsLoad {
					return true, err
				}
				return false, nil
			}
			return true, nil
		},
	)
}

func r53Tags(tags map[string]string) []*route53.Tag {
	rtags := make([]*route53.Tag, 0, len(tags))
	for k, v := range tags {
		k, v := k, v
		rtags = append(rtags, &route53.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return rtags
}
