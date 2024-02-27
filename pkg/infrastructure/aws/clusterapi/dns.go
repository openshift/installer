package clusterapi

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
)

func createDNSRecords(installConfig *installconfig.InstallConfig, apiTarget, apiIntTarget, phzID string) error {
	logrus.Infof("Creating Route53 records for control plane load balancer")
	ssn, err := installConfig.AWS.Session(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	client := awsconfig.NewClient(ssn)
	r53cfg := awsconfig.GetR53ClientCfg(ssn, installConfig.Config.AWS.HostedZoneRole)
	err = client.CreateOrUpdateRecord(installConfig.Config, r53cfg, apiTarget, apiIntTarget, phzID)
	if err != nil {
		return fmt.Errorf("failed to create route53 records: %w", err)
	}
	logrus.Infof("Created Route53 records for control plane load balancer")
	return nil
}

func createHostedZone(ctx context.Context, client route53iface.Route53API, userTags map[string]string, infraID, name, vpcID, region string, isPrivate bool) (*route53.HostedZone, error) {
	var res *route53.CreateHostedZoneOutput

	callRef := fmt.Sprintf("%d", time.Now().Unix())
	res, err := client.CreateHostedZoneWithContext(ctx, &route53.CreateHostedZoneInput{
		CallerReference: aws.String(callRef),
		Name:            aws.String(name),
		HostedZoneConfig: &route53.HostedZoneConfig{
			PrivateZone: aws.Bool(isPrivate),
			Comment:     aws.String("Created by Openshift Installer"),
		},
		VPC: &route53.VPC{
			VPCId:     aws.String(vpcID),
			VPCRegion: aws.String(region),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating private hosted zone: %w", err)
	}

	if res == nil {
		return nil, fmt.Errorf("unexpected output from hosted zone creation")
	}
	// Tag the hosted zone
	tags := mergeTags(userTags, map[string]string{"Name": fmt.Sprintf("%s-int", infraID)})
	_, err = client.ChangeTagsForResourceWithContext(ctx, &route53.ChangeTagsForResourceInput{
		ResourceType: aws.String("hostedzone"),
		ResourceId:   res.HostedZone.Id,
		AddTags:      r53Tags(tags),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to tag private hosted zone: %w", err)
	}
	logrus.Infoln("Tagged private hosted zone")

	// Set SOA minimum TTL
	recordSet, err := existingRecordSet(ctx, client, res.HostedZone.Id, name, "SOA")
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
		HostedZoneId: res.HostedZone.Id,
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

	return res.HostedZone, nil
}

func mergeTags(lhsTags, rhsTags map[string]string) map[string]string {
	merged := make(map[string]string, len(lhsTags)+len(rhsTags))
	for k, v := range lhsTags {
		merged[k] = v
	}
	for k, v := range rhsTags {
		merged[k] = v
	}
	return merged
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
	return nil, fmt.Errorf("not found")
}

func fqdn(name string) string {
	n := len(name)
	if n == 0 || name[n-1] == '.' {
		return name
	}
	return name + "."
}

func cleanRecordName(name string) string {
	s, err := strconv.Unquote(`"` + name + `"`)
	if err != nil {
		return name
	}
	return s
}
