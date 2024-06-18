package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awss "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

//go:generate mockgen -source=./route53.go -destination=mock/awsroute53_generated.go -package=mock

// regions for which ALIAS records are not available
// https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/govcloud-r53.html
var cnameRegions = sets.New[string]("us-gov-west-1", "us-gov-east-1")

// API represents the calls made to the API.
type API interface {
	GetHostedZone(hostedZone string, cfg *aws.Config) (*route53.GetHostedZoneOutput, error)
	ValidateZoneRecords(zone *route53.HostedZone, zoneName string, zonePath *field.Path, ic *types.InstallConfig, cfg *aws.Config) field.ErrorList
	GetBaseDomain(baseDomainName string) (*route53.HostedZone, error)
	GetSubDomainDNSRecords(hostedZone *route53.HostedZone, ic *types.InstallConfig, cfg *aws.Config) ([]string, error)
}

// Client makes calls to the AWS Route53 API.
type Client struct {
	ssn *awss.Session
}

// NewClient initializes a client with a session.
func NewClient(ssn *awss.Session) *Client {
	client := &Client{
		ssn: ssn,
	}
	return client
}

// GetHostedZone attempts to get the hosted zone from the AWS Route53 instance
func (c *Client) GetHostedZone(hostedZone string, cfg *aws.Config) (*route53.GetHostedZoneOutput, error) {
	// build a new Route53 instance from the same session that made it here
	r53 := route53.New(c.ssn, cfg)

	// validate that the hosted zone exists
	hostedZoneOutput, err := r53.GetHostedZone(&route53.GetHostedZoneInput{Id: aws.String(hostedZone)})
	if err != nil {
		return nil, fmt.Errorf("could not get hosted zone: %s: %w", hostedZone, err)
	}
	return hostedZoneOutput, nil
}

// ValidateZoneRecords Attempts to validate each of the candidate HostedZones against the Config
func (c *Client) ValidateZoneRecords(zone *route53.HostedZone, zoneName string, zonePath *field.Path, ic *types.InstallConfig, cfg *aws.Config) field.ErrorList {
	allErrs := field.ErrorList{}

	problematicRecords, err := c.GetSubDomainDNSRecords(zone, ic, cfg)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(zonePath,
			fmt.Errorf("could not list record sets for domain %q: %w", zoneName, err)))
	}

	if len(problematicRecords) > 0 {
		detail := fmt.Sprintf(
			"the zone already has record sets for the domain of the cluster: [%s]",
			strings.Join(problematicRecords, ", "),
		)
		allErrs = append(allErrs, field.Invalid(zonePath, zoneName, detail))
	}

	return allErrs
}

// GetSubDomainDNSRecords Validates the hostedZone against the cluster domain, and ensures that the
// cluster domain does not have a current record set for the hostedZone
func (c *Client) GetSubDomainDNSRecords(hostedZone *route53.HostedZone, ic *types.InstallConfig, cfg *aws.Config) ([]string, error) {
	dottedClusterDomain := ic.ClusterDomain() + "."

	// validate that the domain of the hosted zone is the cluster domain or a parent of the cluster domain
	if !isHostedZoneDomainParentOfClusterDomain(hostedZone, dottedClusterDomain) {
		return nil, fmt.Errorf("hosted zone domain %q is not a parent of the cluster domain %q", *hostedZone.Name, dottedClusterDomain)
	}

	r53 := route53.New(c.ssn, cfg)

	var problematicRecords []string
	// validate that the hosted zone does not already have any record sets for the cluster domain
	if err := r53.ListResourceRecordSetsPages(
		&route53.ListResourceRecordSetsInput{HostedZoneId: hostedZone.Id},
		func(out *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			for _, recordSet := range out.ResourceRecordSets {
				name := aws.StringValue(recordSet.Name)
				if skipRecord(name, dottedClusterDomain) {
					continue
				}
				problematicRecords = append(problematicRecords, fmt.Sprintf("%s (%s)", name, aws.StringValue(recordSet.Type)))
			}
			return !lastPage
		},
	); err != nil {
		return nil, err
	}

	return problematicRecords, nil
}

func skipRecord(recordName string, dottedClusterDomain string) bool {
	// skip record sets that are not sub-domains of the cluster domain. Such record sets may exist for
	// hosted zones that are used for other clusters or other purposes.
	if !strings.HasSuffix(recordName, "."+dottedClusterDomain) {
		return true
	}
	// skip record sets that are the cluster domain. Record sets for the cluster domain are fine. If the
	// hosted zone has the name of the cluster domain, then there will be NS and SOA record sets for the
	// cluster domain.
	if len(recordName) == len(dottedClusterDomain) {
		return true
	}

	return false
}

func isHostedZoneDomainParentOfClusterDomain(hostedZone *route53.HostedZone, dottedClusterDomain string) bool {
	if *hostedZone.Name == dottedClusterDomain {
		return true
	}
	return strings.HasSuffix(dottedClusterDomain, "."+*hostedZone.Name)
}

// GetBaseDomain Gets the Domain Zone with the matching domain name from the session
func (c *Client) GetBaseDomain(baseDomainName string) (*route53.HostedZone, error) {
	baseDomainZone, err := GetPublicZone(c.ssn, baseDomainName)
	if err != nil {
		return nil, fmt.Errorf("could not find public zone: %s: %w", baseDomainName, err)
	}
	return baseDomainZone, nil
}

// GetR53ClientCfg creates a config for the route53 client by determining
// whether it is needed to obtain STS assume role credentials.
func GetR53ClientCfg(sess *awss.Session, roleARN string) *aws.Config {
	if roleARN == "" {
		return nil
	}

	creds := stscreds.NewCredentials(sess, roleARN)
	return &aws.Config{Credentials: creds}
}

// CreateOrUpdateRecord Creates or Updates the Route53 Record for the cluster endpoint.
func (c *Client) CreateOrUpdateRecord(ctx context.Context, ic *types.InstallConfig, target string, intTarget string, phzID string, aliasZoneID string) error {
	useCNAME := cnameRegions.Has(ic.AWS.Region)

	apiName := fmt.Sprintf("api.%s.", ic.ClusterDomain())
	apiIntName := fmt.Sprintf("api-int.%s.", ic.ClusterDomain())

	// Create api record in public zone
	if ic.Publish == types.ExternalPublishingStrategy {
		zone, err := c.GetBaseDomain(ic.BaseDomain)
		if err != nil {
			return err
		}

		svc := route53.New(c.ssn) // we dont want to assume role here
		if _, err := createRecord(ctx, svc, aws.StringValue(zone.Id), apiName, target, aliasZoneID, useCNAME); err != nil {
			return fmt.Errorf("failed to create records for api: %w", err)
		}
		logrus.Debugln("Created public API record in public zone")
	}

	// Create service with assumed role for PHZ
	svc := route53.New(c.ssn, GetR53ClientCfg(c.ssn, ic.AWS.HostedZoneRole))

	// Create api record in private zone
	if _, err := createRecord(ctx, svc, phzID, apiName, intTarget, aliasZoneID, useCNAME); err != nil {
		return fmt.Errorf("failed to create records for api: %w", err)
	}
	logrus.Debugln("Created public API record in private zone")

	// Create api-int record in private zone
	if _, err := createRecord(ctx, svc, phzID, apiIntName, intTarget, aliasZoneID, useCNAME); err != nil {
		return fmt.Errorf("failed to create records for api-int: %w", err)
	}
	logrus.Debugln("Created private API record in private zone")

	return nil
}

func createRecord(ctx context.Context, client *route53.Route53, zoneID, name, dnsName, aliasZoneID string, useCNAME bool) (*route53.ChangeInfo, error) {
	recordSet := &route53.ResourceRecordSet{
		Name: aws.String(name),
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
			Comment: aws.String(fmt.Sprintf("Creating record %s", name)),
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

// HostedZoneInput defines the input parameters for hosted zone creation.
type HostedZoneInput struct {
	Name     string
	InfraID  string
	VpcID    string
	Region   string
	Role     string
	UserTags map[string]string
}

// CreateHostedZone creates a private hosted zone.
func (c *Client) CreateHostedZone(ctx context.Context, input *HostedZoneInput) (*route53.HostedZone, error) {
	cfg := GetR53ClientCfg(c.ssn, input.Role)
	svc := route53.New(c.ssn, cfg)

	// CallerReference needs to be a unique string. We include the infra id,
	// which is unique, in case that is helpful for human debugging. A random
	// string of an arbitrary length is appended in case the infra id is reused
	// which is generally not supposed to happen but does in some edge cases.
	callerRef := aws.String(fmt.Sprintf("%s-%s", input.InfraID, rand.String(5)))

	res, err := svc.CreateHostedZoneWithContext(ctx, &route53.CreateHostedZoneInput{
		CallerReference: callerRef,
		Name:            aws.String(input.Name),
		HostedZoneConfig: &route53.HostedZoneConfig{
			PrivateZone: aws.Bool(true),
			Comment:     aws.String("Created by Openshift Installer"),
		},
		VPC: &route53.VPC{
			VPCId:     aws.String(input.VpcID),
			VPCRegion: aws.String(input.Region),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating private hosted zone: %w", err)
	}

	if res == nil {
		return nil, fmt.Errorf("error creating private hosted zone: %w", err)
	}

	// Tag the hosted zone
	tags := mergeTags(input.UserTags, map[string]string{
		"Name": fmt.Sprintf("%s-int", input.InfraID),
	})
	_, err = svc.ChangeTagsForResourceWithContext(ctx, &route53.ChangeTagsForResourceInput{
		ResourceType: aws.String("hostedzone"),
		ResourceId:   res.HostedZone.Id,
		AddTags:      r53Tags(tags),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to tag private hosted zone: %w", err)
	}

	// Set SOA minimum TTL
	recordSet, err := existingRecordSet(ctx, svc, res.HostedZone.Id, input.Name, "SOA")
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
	fields[0] = "60"
	record.Value = aws.String(strings.Join(fields, " "))
	req, err := svc.ChangeResourceRecordSetsWithContext(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: res.HostedZone.Id,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String("UPSERT"),
					ResourceRecordSet: recordSet,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set SOA TTL to minimum: %w", err)
	}

	if err = svc.WaitUntilResourceRecordSetsChangedWithContext(ctx, &route53.GetChangeInput{Id: req.ChangeInfo.Id}); err != nil {
		return nil, fmt.Errorf("failed to wait for SOA TTL change: %w", err)
	}

	return res.HostedZone, nil
}

func existingRecordSet(ctx context.Context, client *route53.Route53, zoneID *string, recordName string, recordType string) (*route53.ResourceRecordSet, error) {
	name := fqdn(strings.ToLower(recordName))
	res, err := client.ListResourceRecordSetsWithContext(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:    zoneID,
		StartRecordName: aws.String(name),
		StartRecordType: aws.String(recordType),
		MaxItems:        aws.String("1"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list record sets: %w", err)
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
		rtags = append(rtags, &route53.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return rtags
}

// See https://docs.aws.amazon.com/general/latest/gr/elb.html#elb_region

// HostedZoneIDPerRegionNLBMap maps HostedZoneIDs from known regions.
var HostedZoneIDPerRegionNLBMap = map[string]string{
	endpoints.AfSouth1RegionID:     "Z203XCE67M25HM",
	endpoints.ApEast1RegionID:      "Z12Y7K3UBGUAD1",
	endpoints.ApNortheast1RegionID: "Z31USIVHYNEOWT",
	endpoints.ApNortheast2RegionID: "ZIBE1TIR4HY56",
	endpoints.ApNortheast3RegionID: "Z1GWIQ4HH19I5X",
	endpoints.ApSouth1RegionID:     "ZVDDRBQ08TROA",
	endpoints.ApSouth2RegionID:     "Z0711778386UTO08407HT",
	endpoints.ApSoutheast1RegionID: "ZKVM4W9LS7TM",
	endpoints.ApSoutheast2RegionID: "ZCT6FZBF4DROD",
	endpoints.ApSoutheast3RegionID: "Z01971771FYVNCOVWJU1G",
	endpoints.ApSoutheast4RegionID: "Z01156963G8MIIL7X90IV",
	endpoints.CaCentral1RegionID:   "Z2EPGBW3API2WT",
	endpoints.CnNorth1RegionID:     "Z3QFB96KMJ7ED6",
	endpoints.CnNorthwest1RegionID: "ZQEIKTCZ8352D",
	endpoints.EuCentral1RegionID:   "Z3F0SRJ5LGBH90",
	endpoints.EuCentral2RegionID:   "Z02239872DOALSIDCX66S",
	endpoints.EuNorth1RegionID:     "Z1UDT6IFJ4EJM",
	endpoints.EuSouth1RegionID:     "Z23146JA1KNAFP",
	endpoints.EuSouth2RegionID:     "Z1011216NVTVYADP1SSV",
	endpoints.EuWest1RegionID:      "Z2IFOLAFXWLO4F",
	endpoints.EuWest2RegionID:      "ZD4D7Y8KGAS4G",
	endpoints.EuWest3RegionID:      "Z1CMS0P5QUZ6D5",
	endpoints.MeCentral1RegionID:   "Z00282643NTTLPANJJG2P",
	endpoints.MeSouth1RegionID:     "Z3QSRYVP46NYYV",
	endpoints.SaEast1RegionID:      "ZTK26PT1VY4CU",
	endpoints.UsEast1RegionID:      "Z26RNL4JYFTOTI",
	endpoints.UsEast2RegionID:      "ZLMOA37VPKANP",
	endpoints.UsGovEast1RegionID:   "Z1ZSMQQ6Q24QQ8",
	endpoints.UsGovWest1RegionID:   "ZMG1MZ2THAWF1",
	endpoints.UsWest1RegionID:      "Z24FKFUX50B4VW",
	endpoints.UsWest2RegionID:      "Z18D5FSROUN65G",
}
