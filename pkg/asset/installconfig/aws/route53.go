package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

//go:generate mockgen -source=./route53.go -destination=mock/awsroute53_generated.go -package=mock

const (
	// RecordChangeMaxWaitTime defines the max duration to wait for a record set
	// to apply changes.
	RecordChangeMaxWaitTime = 5 * time.Minute
)

// regions for which ALIAS records are not available.
// https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/govcloud-r53.html
var cnameRegions = sets.New[string]("us-gov-west-1", "us-gov-east-1")

// API represents the calls made to the AWS Route 53 API.
type API interface {
	GetHostedZone(ctx context.Context, hostedZone string) (*route53.GetHostedZoneOutput, error)
	ValidateZoneRecords(ctx context.Context, zone *route53types.HostedZone, zoneName string, fldPath *field.Path, ic *types.InstallConfig) field.ErrorList
	GetBaseDomain(ctx context.Context, baseDomainName string) (*route53types.HostedZone, error)
	GetSubDomainDNSRecords(ctx context.Context, hostedZone *route53types.HostedZone, ic *types.InstallConfig) ([]string, error)
}

// Client wraps Route53 Client to define custom API calls
// to the AWS Route53 API.
type Client struct {
	route53Client *route53.Client
}

// Make sure Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient initializes a client with a session.
// To configure the client to assume an IAM role, define a non-empty roleArn.
// Note: If defined, the client will use the assumed identity for all its operations.
func NewClient(ctx context.Context, opts EndpointOptions, roleArn string) (*Client, error) {
	client, err := NewRoute53Client(ctx, opts, roleArn)
	if err != nil {
		return nil, err
	}
	return &Client{route53Client: client}, nil
}

// GetHostedZone attempts to get the hosted zone from the AWS Route53 instance.
func (c *Client) GetHostedZone(ctx context.Context, hostedZone string) (*route53.GetHostedZoneOutput, error) {
	hostedZoneOutput, err := c.route53Client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: aws.String(hostedZone)})
	if err != nil {
		return nil, fmt.Errorf("failed to get hosted zone: %s: %w", hostedZone, err)
	}
	return hostedZoneOutput, nil
}

// ValidateZoneRecords attempts to validate each of the candidate HostedZones against the Config.
func (c *Client) ValidateZoneRecords(ctx context.Context, zone *route53types.HostedZone, zoneName string, fldPath *field.Path, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	problematicRecords, err := c.GetSubDomainDNSRecords(ctx, zone, ic)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath,
			fmt.Errorf("could not list record sets for domain %q: %w", zoneName, err)))
	}

	// validate that the hosted zone does not already have any record sets for the cluster domain
	if len(problematicRecords) > 0 {
		detail := fmt.Sprintf(
			"the zone already has record sets for the domain of the cluster: [%s]",
			strings.Join(problematicRecords, ", "),
		)
		allErrs = append(allErrs, field.Invalid(fldPath, zoneName, detail))
	}

	return allErrs
}

// GetSubDomainDNSRecords validates the hostedZone against the cluster domain, and ensures that the
// cluster domain does not have a current record set for the hostedZone.
func (c *Client) GetSubDomainDNSRecords(ctx context.Context, hostedZone *route53types.HostedZone, ic *types.InstallConfig) ([]string, error) {
	dottedClusterDomain := ic.ClusterDomain() + "."

	// validate that the domain of the hosted zone is the cluster domain or a parent of the cluster domain
	if !isHostedZoneDomainParentOfClusterDomain(hostedZone, dottedClusterDomain) {
		return nil, fmt.Errorf("hosted zone domain %q is not a parent of the cluster domain %q", *hostedZone.Name, dottedClusterDomain)
	}

	var problematicRecords []string

	paginator := route53.NewListResourceRecordSetsPaginator(c.route53Client, &route53.ListResourceRecordSetsInput{HostedZoneId: hostedZone.Id})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list record sets for zone %s: %w", aws.ToString(hostedZone.Id), err)
		}

		for _, recordSet := range page.ResourceRecordSets {
			name := aws.ToString(recordSet.Name)
			if skipRecord(name, dottedClusterDomain) {
				continue
			}
			problematicRecords = append(problematicRecords, fmt.Sprintf("%s (%s)", name, recordSet.Type))
		}
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

func isHostedZoneDomainParentOfClusterDomain(hostedZone *route53types.HostedZone, dottedClusterDomain string) bool {
	if aws.ToString(hostedZone.Name) == dottedClusterDomain {
		return true
	}
	return strings.HasSuffix(dottedClusterDomain, "."+aws.ToString(hostedZone.Name))
}

// GetBaseDomain gets the Domain Zone with the matching domain name.
func (c *Client) GetBaseDomain(ctx context.Context, baseDomainName string) (*route53types.HostedZone, error) {
	baseDomainZone, err := GetPublicZone(ctx, c.route53Client, baseDomainName)
	if err != nil {
		return nil, fmt.Errorf("failed to find public zone: %s: %w", baseDomainName, err)
	}
	return baseDomainZone, nil
}

// CreateRecordInput collects information for creating a record.
type CreateRecordInput struct {
	// Fully qualified record domain name.
	Name string
	// Cluster Region.
	Region string
	// Where to route the DNS queries to.
	DNSTarget string
	// ID of the Hosted Zone.
	ZoneID string
	// ID of the Hosted Zone for Alias record.
	AliasZoneID string
}

// CreateOrUpdateRecord Creates or Updates the Route53 Record for the cluster endpoint.
func (c *Client) CreateOrUpdateRecord(ctx context.Context, in *CreateRecordInput) error {
	recordSet := &route53types.ResourceRecordSet{
		Name: aws.String(in.Name),
	}
	if cnameRegions.Has(in.Region) {
		recordSet.Type = route53types.RRTypeCname
		recordSet.TTL = aws.Int64(10)
		recordSet.ResourceRecords = []route53types.ResourceRecord{
			{Value: aws.String(in.DNSTarget)},
		}
	} else {
		recordSet.Type = route53types.RRTypeA
		recordSet.AliasTarget = &route53types.AliasTarget{
			DNSName:              aws.String(in.DNSTarget),
			HostedZoneId:         aws.String(in.AliasZoneID),
			EvaluateTargetHealth: false,
		}
	}

	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(in.ZoneID),
		ChangeBatch: &route53types.ChangeBatch{
			Comment: aws.String(fmt.Sprintf("Creating record %s", in.Name)),
			Changes: []route53types.Change{
				{
					Action:            route53types.ChangeActionUpsert,
					ResourceRecordSet: recordSet,
				},
			},
		},
	}

	_, err := c.route53Client.ChangeResourceRecordSets(ctx, input)
	return err
}

// HostedZoneInput defines the input parameters for hosted zone creation.
type HostedZoneInput struct {
	Name     string
	InfraID  string
	VpcID    string
	Region   string
	UserTags map[string]string
}

// CreateHostedZone creates a private hosted zone.
func (c *Client) CreateHostedZone(ctx context.Context, input *HostedZoneInput) (*route53types.HostedZone, error) {
	// CallerReference needs to be a unique string. We include the infra id,
	// which is unique, in case that is helpful for human debugging. A random
	// string of an arbitrary length is appended in case the infra id is reused
	// which is generally not supposed to happen but does in some edge cases.
	callerRef := aws.String(fmt.Sprintf("%s-%s", input.InfraID, rand.String(5)))

	res, err := c.route53Client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
		CallerReference: callerRef,
		Name:            aws.String(input.Name),
		HostedZoneConfig: &route53types.HostedZoneConfig{
			PrivateZone: true,
			Comment:     aws.String("Created by Openshift Installer"),
		},
		VPC: &route53types.VPC{
			VPCId:     aws.String(input.VpcID),
			VPCRegion: route53types.VPCRegion(input.Region),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create private hosted zone: %w", err)
	}

	// This case should not occur, but we safeguard it against nil panic.
	if res == nil {
		return nil, fmt.Errorf("failed to create private hosted zone")
	}

	// Tag the hosted zone
	tags := mergeTags(input.UserTags, map[string]string{
		"Name": fmt.Sprintf("%s-int", input.InfraID),
		fmt.Sprintf("kubernetes.io/cluster/%s", input.InfraID): "owned",
	})

	// Of all the route53 actions being used in the installer,
	// the AWS SDK does sanitize the hosted zone ID to remove any resource prefix
	// except ChangeTagsForResource.
	hostedZoneID := sanitizeHostedZoneID(aws.ToString(res.HostedZone.Id))
	_, err = c.route53Client.ChangeTagsForResource(ctx, &route53.ChangeTagsForResourceInput{
		ResourceType: route53types.TagResourceTypeHostedzone,
		ResourceId:   aws.String(hostedZoneID),
		AddTags:      r53Tags(tags),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to tag private hosted zone: %w", err)
	}

	// Set SOA minimum TTL
	recordSet, err := existingRecordSet(ctx, c.route53Client, res.HostedZone.Id, input.Name, route53types.RRTypeSoa)
	if err != nil {
		return nil, fmt.Errorf("failed to find SOA record set for private zone: %w", err)
	}
	if len(recordSet.ResourceRecords) == 0 || recordSet.ResourceRecords[0].Value == nil {
		return nil, fmt.Errorf("failed to find SOA record for private zone")
	}
	record := recordSet.ResourceRecords[0]
	fields := strings.Split(aws.ToString(record.Value), " ")
	if len(fields) != 7 {
		return nil, fmt.Errorf("SOA record value has %d fields, expected 7", len(fields))
	}
	fields[0] = "60"
	record.Value = aws.String(strings.Join(fields, " "))
	req, err := c.route53Client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: res.HostedZone.Id,
		ChangeBatch: &route53types.ChangeBatch{
			Changes: []route53types.Change{
				{
					Action:            route53types.ChangeActionUpsert,
					ResourceRecordSet: recordSet,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set SOA TTL to minimum: %w", err)
	}

	waiter := route53.NewResourceRecordSetsChangedWaiter(c.route53Client)
	if err := waiter.Wait(ctx, &route53.GetChangeInput{Id: req.ChangeInfo.Id}, RecordChangeMaxWaitTime); err != nil {
		return nil, fmt.Errorf("failed to wait for SOA TTL change: %w", err)
	}

	return res.HostedZone, nil
}

func existingRecordSet(ctx context.Context, client *route53.Client, zoneID *string, recordName string, recordType route53types.RRType) (*route53types.ResourceRecordSet, error) {
	name := fqdn(strings.ToLower(recordName))
	res, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:    zoneID,
		StartRecordName: aws.String(name),
		StartRecordType: recordType,
		MaxItems:        aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list record sets: %w", err)
	}

	for _, rs := range res.ResourceRecordSets {
		resName := strings.ToLower(cleanRecordName(aws.ToString(rs.Name)))
		if resName == name && rs.Type == recordType {
			return &rs, nil
		}
	}

	return nil, fmt.Errorf("record set not found for type %s, name %s", recordType, recordName)
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

func r53Tags(tags map[string]string) []route53types.Tag {
	rtags := make([]route53types.Tag, 0, len(tags))
	for k, v := range tags {
		rtags = append(rtags, route53types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return rtags
}

// sanitizeHostedZoneID splits apart the Route53 hostedZone ID and returns the last piece.
// For example, an input of /hostedzone/Z12345678ABCDEFGHIJK will return Z12345678ABCDEFGHIJK.
func sanitizeHostedZoneID(hostedZoneID string) string {
	idx := strings.LastIndex(hostedZoneID, "/")
	if idx > 0 {
		return hostedZoneID[idx+1:]
	}
	return hostedZoneID
}

// HostedZoneIDPerRegionNLBMap maps HostedZoneIDs from known regions.
// See https://docs.aws.amazon.com/general/latest/gr/elb.html#elb_region
var HostedZoneIDPerRegionNLBMap = map[string]string{
	AfSouth1RegionID:     "Z203XCE67M25HM",
	ApEast1RegionID:      "Z12Y7K3UBGUAD1",
	ApNortheast1RegionID: "Z31USIVHYNEOWT",
	ApNortheast2RegionID: "ZIBE1TIR4HY56",
	ApNortheast3RegionID: "Z1GWIQ4HH19I5X",
	ApSouth1RegionID:     "ZVDDRBQ08TROA",
	ApSouth2RegionID:     "Z0711778386UTO08407HT",
	ApSoutheast1RegionID: "ZKVM4W9LS7TM",
	ApSoutheast2RegionID: "ZCT6FZBF4DROD",
	ApSoutheast3RegionID: "Z01971771FYVNCOVWJU1G",
	ApSoutheast4RegionID: "Z01156963G8MIIL7X90IV",
	CaCentral1RegionID:   "Z2EPGBW3API2WT",
	CnNorth1RegionID:     "Z3QFB96KMJ7ED6",
	CnNorthwest1RegionID: "ZQEIKTCZ8352D",
	EuCentral1RegionID:   "Z3F0SRJ5LGBH90",
	EuCentral2RegionID:   "Z02239872DOALSIDCX66S",
	EuNorth1RegionID:     "Z1UDT6IFJ4EJM",
	EuSouth1RegionID:     "Z23146JA1KNAFP",
	EuSouth2RegionID:     "Z1011216NVTVYADP1SSV",
	EuWest1RegionID:      "Z2IFOLAFXWLO4F",
	EuWest2RegionID:      "ZD4D7Y8KGAS4G",
	EuWest3RegionID:      "Z1CMS0P5QUZ6D5",
	MeCentral1RegionID:   "Z00282643NTTLPANJJG2P",
	MeSouth1RegionID:     "Z3QSRYVP46NYYV",
	SaEast1RegionID:      "ZTK26PT1VY4CU",
	UsEast1RegionID:      "Z26RNL4JYFTOTI",
	UsEast2RegionID:      "ZLMOA37VPKANP",
	UsGovEast1RegionID:   "Z1ZSMQQ6Q24QQ8",
	UsGovWest1RegionID:   "ZMG1MZ2THAWF1",
	UsWest1RegionID:      "Z24FKFUX50B4VW",
	UsWest2RegionID:      "Z18D5FSROUN65G",
}
