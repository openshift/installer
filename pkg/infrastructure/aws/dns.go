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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

type dnsInput struct {
	region                      string
	baseDomain                  string
	clusterDomain               string
	vpcID                       string
	loadBalancerExternalZoneID  string
	loadBalancerExternalZoneDNS string
	loadBalancerInternalZoneID  string
	loadBalancerInternalZoneDNS string
}

func createDNSResources(ctx context.Context, logger *logrus.Logger, session *session.Session, dnsInput *dnsInput) error {
	route53Client := route53.New(session)
	publicZoneID, err := dnsInput.LookupPublicZone(ctx, logger, route53Client)
	if err != nil {
		return err
	}

	privateZoneID, err := dnsInput.CreatePrivateZone(ctx, logger, route53Client, dnsInput.clusterDomain, dnsInput.vpcID)
	if err != nil {
		return err
	}

	// Create API record in public zone.
	if err := createRecord(logger, route53Client, publicZoneID, fmt.Sprintf("api.%s", dnsInput.clusterDomain),
		dnsInput.loadBalancerExternalZoneDNS, dnsInput.loadBalancerExternalZoneID); err != nil {
		return err
	}

	// Create API records in private zone.
	if err := createRecord(logger, route53Client, privateZoneID, fmt.Sprintf("api.%s", dnsInput.clusterDomain),
		dnsInput.loadBalancerInternalZoneDNS, dnsInput.loadBalancerInternalZoneID); err != nil {
		return err
	}

	if err := createRecord(logger, route53Client, privateZoneID, fmt.Sprintf("api-int.%s", dnsInput.clusterDomain),
		dnsInput.loadBalancerInternalZoneDNS, dnsInput.loadBalancerInternalZoneID); err != nil {
		return err
	}

	return nil
}
func (o *dnsInput) LookupPublicZone(ctx context.Context, l *logrus.Logger, client route53iface.Route53API) (string, error) {
	name := o.baseDomain
	id, err := lookupZone(ctx, client, name, false)
	if err != nil {
		l.Error(err, "Public zone not found", "name", name)
		return "", err
	}
	l.Info("Found existing public zone", "name", name, "id", id)
	return id, nil
}

func lookupZone(ctx context.Context, client route53iface.Route53API, name string, isPrivateZone bool) (string, error) {
	var res *route53.HostedZone
	f := func(resp *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
		for idx, zone := range resp.HostedZones {
			if zone.Config != nil && isPrivateZone == aws.BoolValue(zone.Config.PrivateZone) && strings.TrimSuffix(aws.StringValue(zone.Name), ".") == strings.TrimSuffix(name, ".") {
				res = resp.HostedZones[idx]
				return false
			}
		}
		return !lastPage
	}
	if err := retryRoute53WithBackoff(ctx, func() error {
		if err := client.ListHostedZonesPagesWithContext(ctx, &route53.ListHostedZonesInput{}, f); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to list hosted zones: %w", err)
	}
	if res == nil {
		return "", fmt.Errorf("hosted zone %s not found", name)
	}
	return cleanZoneID(*res.Id), nil
}

func (o *dnsInput) CreatePrivateZone(ctx context.Context, l *logrus.Logger, client route53iface.Route53API, name, vpcID string) (string, error) {
	id, err := lookupZone(ctx, client, name, true)
	if err == nil {
		l.Info("Found existing private zone", "name", name, "id", id)
		err := setSOAMinimum(ctx, client, id, name)
		if err != nil {
			return "", err
		}
		return id, err
	}

	var res *route53.CreateHostedZoneOutput
	if err := retryRoute53WithBackoff(ctx, func() error {
		callRef := fmt.Sprintf("%d", time.Now().Unix())
		if output, err := client.CreateHostedZoneWithContext(ctx, &route53.CreateHostedZoneInput{
			CallerReference: aws.String(callRef),
			Name:            aws.String(name),
			HostedZoneConfig: &route53.HostedZoneConfig{
				PrivateZone: aws.Bool(true),
			},
			VPC: &route53.VPC{
				VPCId:     aws.String(vpcID),
				VPCRegion: aws.String(o.region),
			},
		}); err != nil {
			return err
		} else {
			res = output
			return nil
		}
	}); err != nil {
		return "", fmt.Errorf("failed to create hosted zone: %w", err)
	}
	if res == nil {
		return "", fmt.Errorf("unexpected output from hosted zone creation")
	}
	id = cleanZoneID(*res.HostedZone.Id)
	l.Info("Created private zone", "name", name, "id", id)

	err = setSOAMinimum(ctx, client, id, name)
	if err != nil {
		return "", err
	}

	return id, nil
}

func setSOAMinimum(ctx context.Context, client route53iface.Route53API, id, name string) error {
	recordSet, err := findRecord(ctx, client, id, name, "SOA")
	if err != nil {
		return err
	}
	if recordSet == nil || recordSet.ResourceRecords[0] == nil || recordSet.ResourceRecords[0].Value == nil {
		return fmt.Errorf("SOA record for private zone %s not found: %w", name, err)
	}
	record := recordSet.ResourceRecords[0]
	fields := strings.Split(*record.Value, " ")
	if len(fields) != 7 {
		return fmt.Errorf("SOA record value has %d fields, expected 7", len(fields))
	}
	fields[6] = "60"
	record.Value = aws.String(strings.Join(fields, " "))
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(id),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String("UPSERT"),
					ResourceRecordSet: recordSet,
				},
			},
		},
	}
	_, err = client.ChangeResourceRecordSetsWithContext(ctx, input)
	return err
}

func findRecord(ctx context.Context, client route53iface.Route53API, id, name string, recordType string) (*route53.ResourceRecordSet, error) {
	recordName := fqdn(strings.ToLower(name))
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(id),
		StartRecordName: aws.String(recordName),
		StartRecordType: aws.String(recordType),
		MaxItems:        aws.String("1"),
	}

	var record *route53.ResourceRecordSet
	err := client.ListResourceRecordSetsPagesWithContext(ctx, input, func(resp *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
		if len(resp.ResourceRecordSets) == 0 {
			return false
		}

		recordSet := resp.ResourceRecordSets[0]
		responseName := strings.ToLower(cleanRecordName(*recordSet.Name))
		responseType := strings.ToUpper(*recordSet.Type)

		if recordName != responseName {
			return false
		}
		if recordType != responseType {
			return false
		}

		record = recordSet
		return false
	})

	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, fmt.Errorf("record not found")
	}
	return record, nil
}

func cleanRecordName(name string) string {
	str := name
	s, err := strconv.Unquote(`"` + str + `"`)
	if err != nil {
		return str
	}
	return s
}

func fqdn(name string) string {
	n := len(name)
	if n == 0 || name[n-1] == '.' {
		return name
	} else {
		return name + "."
	}
}

func retryRoute53WithBackoff(ctx context.Context, fn func() error) error {
	backoff := wait.Backoff{
		Duration: 1 * time.Second,
		Steps:    10,
		Factor:   1.5,
	}
	retriable := func(e error) bool {
		if !IsErrorRetryable(e) {
			return false
		}
		select {
		case <-ctx.Done():
			return false
		default:
			return true
		}
	}
	// TODO: inspect the error for throttling details?
	return retry.OnError(backoff, retriable, fn)
}

func cleanZoneID(ID string) string {
	return strings.TrimPrefix(ID, "/hostedzone/")
}

func IsErrorRetryable(err error) bool {
	if aggregate, isAggregate := err.(utilerrors.Aggregate); isAggregate {
		if len(aggregate.Errors()) == 1 {
			err = aggregate.Errors()[0]
		} else {
			// We aggregate all errors, utilerrors.Aggregate does for safety reasons not support
			// errors.As (As it can't know what to do when there are multiple matches), so we
			// iterate and bail out if there are only credential load errors
			hasOnlyCredentialLoadErrors := true
			for _, err := range aggregate.Errors() {
				if !isCredentialLoadError(err) {
					hasOnlyCredentialLoadErrors = false
					break
				}
			}
			if hasOnlyCredentialLoadErrors {
				return false
			}

		}
	}

	if isCredentialLoadError(err) {
		return false
	}
	return true
}

func isCredentialLoadError(err error) bool {
	if awsErr := awserr.Error(nil); errors.As(err, &awsErr) && awsErr.Code() == "SharedCredsLoad" {
		return true
	}

	return false
}

func createRecord(l *logrus.Logger, client route53iface.Route53API, zoneID, domain, aliasLBDNSName, aliasLBZoneID string) error {
	// Create alias records for public endpoints
	createAliasRecordInput := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(cleanRecordName(domain)),
						Type: aws.String("A"),
						AliasTarget: &route53.AliasTarget{
							DNSName:              aws.String(aliasLBDNSName),
							HostedZoneId:         aws.String(aliasLBZoneID),
							EvaluateTargetHealth: aws.Bool(false),
						},
					},
				},
			},
		},
	}

	result, err := client.ChangeResourceRecordSets(createAliasRecordInput)
	if err != nil {
		return fmt.Errorf("error creating alias record: %w", err)
	}
	l.Infof("Created DNS record: %v", result.ChangeInfo)

	return nil
}
