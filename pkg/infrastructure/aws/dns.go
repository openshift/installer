package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/sirupsen/logrus"
)

type dnsInput struct {
	clusterID                   string
	region                      string
	baseDomain                  string
	clusterDomain               string
	vpcID                       string
	loadBalancerExternalZoneID  string
	loadBalancerExternalZoneDNS string
	loadBalancerInternalZoneID  string
	loadBalancerInternalZoneDNS string
	additionalTags              map[string]string
	additionalR53Tags           []*route53.Tag
}

func createDNSResources(ctx context.Context, logger *logrus.Logger, session *session.Session, dnsInput *dnsInput) error {
	dnsInput.additionalR53Tags = r53CreateTags(dnsInput.additionalTags)
	route53Client := route53.New(session)
	publicZoneID, err := dnsInput.LookupPublicZone(ctx, logger, route53Client)
	if err != nil {
		return err
	}

	privateZoneID, err := dnsInput.CreatePrivateZone(ctx, logger, route53Client)
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
	logger := l.WithField("name", o.baseDomain)
	id, err := r53FindHostedZone(ctx, client, o.baseDomain, false)
	if err != nil {
		return "", err
	}
	if id == "" {
		return "", fmt.Errorf("public zone %s not found", o.baseDomain)
	}
	logger.WithField("id", id).Infoln("Found existing public zone")
	return id, nil
}

func (o *dnsInput) CreatePrivateZone(ctx context.Context, l *logrus.Logger, client route53iface.Route53API) (string, error) {
	logger := l.WithField("name", o.clusterDomain)
	id, err := r53FindHostedZone(ctx, client, o.clusterDomain, true)
	if err != nil {
		return "", err
	}
	if id == "" {
		id, err = r53CreateHostedZone(ctx, client, o.clusterDomain, o.vpcID, o.region, true)
		if err != nil {
			return "", err
		}
		logger.WithField("id", id).Infoln("Created private zone")
	} else {
		logger.WithField("id", id).Infoln("Found existing private zone")
	}

	// Tag the hosted zone
	tags := append(o.additionalR53Tags, r53Tags(fmt.Sprintf("%s-int", o.clusterID))...)
	if err := r53TagHostedZone(ctx, client, id, tags); err != nil {
		return "", err
	}
	logger.WithField("id", id).Infoln("Tagged private hosted zone")

	if err := setSOAMinimum(ctx, client, id, o.clusterDomain); err != nil {
		return "", err
	}
	logger.WithField("id", id).Infoln("Set private hosted zone SOA mininum")

	return id, nil
}

func createRecord(l *logrus.Logger, client route53iface.Route53API, zoneID, domain, aliasLBDNSName, aliasLBZoneID string) error {
	result, err := r53CreateRecord(client, zoneID, domain, aliasLBDNSName, aliasLBZoneID)
	if err != nil {
		return err
	}
	l.Infof("Created DNS record: %s", result.String())
	return nil
}

func setSOAMinimum(ctx context.Context, client route53iface.Route53API, id, name string) error {
	recordSet, err := r53GetRecord(ctx, client, id, name, "SOA")
	if recordSet == nil || recordSet.ResourceRecords[0] == nil || recordSet.ResourceRecords[0].Value == nil {
		return fmt.Errorf("SOA record for private zone %s not found: %w", name, err)
	}
	record := recordSet.ResourceRecords[0]
	fields := strings.Split(aws.StringValue(record.Value), " ")
	if len(fields) != 7 {
		return fmt.Errorf("SOA record value has %d fields, expected 7", len(fields))
	}
	fields[6] = "60"
	record.Value = aws.String(strings.Join(fields, " "))

	return r53HostedZoneChangeRecord(ctx, client, id, recordSet)
}
