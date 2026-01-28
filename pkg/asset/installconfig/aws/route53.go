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

// https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/govcloud-r53.html
// cnameRegions defines regions for which ALIAS records are not available.
var cnameRegions = sets.New("us-gov-west-1", "us-gov-east-1")

// Route53Clientset is the Route 53 API client set.
type Route53Clientset struct {
	// clients is the map of clients indexed by IAM role if any.
	// An empty role ARN means using loaded credentials without assuming any IAM role.
	clients map[string]*Route53Client

	// endpointOpts is the configurations for AWS SDK client's endpoint resolver.
	endpointOpts EndpointOptions
}

// WithAssumedRole returns a Route 53 client using credentials obtained from assuming an IAM role.
func (c *Route53Clientset) WithAssumedRole(ctx context.Context, roleArn string) (*Route53Client, error) {
	if client, ok := c.clients[roleArn]; ok {
		return client, nil
	}

	client, err := NewRoute53Client(ctx, c.endpointOpts, roleArn)
	if err != nil {
		return nil, fmt.Errorf("failed to create route 53 client: %w", err)
	}

	c.clients[roleArn] = &Route53Client{client}
	return c.clients[roleArn], nil
}

// WithDefault returns a Route 53 client using credentials loaded directly from environment
// without assuming an IAM role.
func (c *Route53Clientset) WithDefault(ctx context.Context) (*Route53Client, error) {
	return c.WithAssumedRole(ctx, "")
}

// NewRoute53Clientset initializes a client with a session.
func NewRoute53Clientset(endpointOpts EndpointOptions) *Route53Clientset {
	return &Route53Clientset{
		endpointOpts: endpointOpts,
		clients:      make(map[string]*Route53Client),
	}
}

// Route53API defines the Route 53 API interface.
type Route53API interface {
	GetHostedZone(ctx context.Context, hostedZone string) (*route53.GetHostedZoneOutput, error)
	ValidateZoneRecords(ctx context.Context, zone *route53types.HostedZone, zoneName string, fldPath *field.Path, ic *types.InstallConfig) field.ErrorList
	GetBaseDomain(ctx context.Context, baseDomainName string) (*route53types.HostedZone, error)
	GetSubDomainDNSRecords(ctx context.Context, hostedZone *route53types.HostedZone, ic *types.InstallConfig) ([]string, error)
}

// Route53Client is the Route 53 API client.
type Route53Client struct {
	client *route53.Client
}

// Make sure Client implements the API interface.
var _ Route53API = (*Route53Client)(nil)

// GetHostedZone attempts to get the hosted zone from the AWS Route53 instance.
func (c *Route53Client) GetHostedZone(ctx context.Context, hostedZone string) (*route53.GetHostedZoneOutput, error) {
	// validate that the hosted zone exists
	hostedZoneOutput, err := c.client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: aws.String(hostedZone)})
	if err != nil {
		return nil, fmt.Errorf("could not get hosted zone: %s: %w", hostedZone, err)
	}
	return hostedZoneOutput, nil
}

// ValidateZoneRecords Attempts to validate each of the candidate HostedZones against the Config
func (c *Route53Client) ValidateZoneRecords(ctx context.Context, zone *route53types.HostedZone, zoneName string, zonePath *field.Path, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	problematicRecords, err := c.GetSubDomainDNSRecords(ctx, zone, ic)
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
func (c *Route53Client) GetSubDomainDNSRecords(ctx context.Context, hostedZone *route53types.HostedZone, ic *types.InstallConfig) ([]string, error) {
	dottedClusterDomain := ic.ClusterDomain() + "."

	// validate that the domain of the hosted zone is the cluster domain or a parent of the cluster domain
	if !isHostedZoneDomainParentOfClusterDomain(hostedZone, dottedClusterDomain) {
		return nil, fmt.Errorf("hosted zone domain %q is not a parent of the cluster domain %q", *hostedZone.Name, dottedClusterDomain)
	}

	var problematicRecords []string

	// validate that the hosted zone does not already have any record sets for the cluster domain
	paginator := route53.NewListResourceRecordSetsPaginator(c.client, &route53.ListResourceRecordSetsInput{HostedZoneId: hostedZone.Id})
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

// GetBaseDomain Gets the Domain Zone with the matching domain name from the session
func (c *Route53Client) GetBaseDomain(ctx context.Context, baseDomainName string) (*route53types.HostedZone, error) {
	baseDomainZone, err := GetPublicZone(ctx, c.client, baseDomainName)
	if err != nil {
		return nil, fmt.Errorf("could not find public zone: %s: %w", baseDomainName, err)
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
func (c *Route53Client) CreateOrUpdateRecord(ctx context.Context, in *CreateRecordInput) error {
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

	_, err := c.client.ChangeResourceRecordSets(ctx, input)

	return err
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
func (c *Route53Client) CreateHostedZone(ctx context.Context, input *HostedZoneInput) (*route53types.HostedZone, error) {
	// CallerReference needs to be a unique string. We include the infra id,
	// which is unique, in case that is helpful for human debugging. A random
	// string of an arbitrary length is appended in case the infra id is reused
	// which is generally not supposed to happen but does in some edge cases.
	callerRef := aws.String(fmt.Sprintf("%s-%s", input.InfraID, rand.String(5)))

	res, err := c.client.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
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

	if res == nil {
		return nil, fmt.Errorf("failed to create private hosted zone")
	}

	// Tag the hosted zone
	tags := mergeTags(input.UserTags, map[string]string{
		"Name": fmt.Sprintf("%s-int", input.InfraID),
		fmt.Sprintf("kubernetes.io/cluster/%s", input.InfraID): "owned",
	})

	// Of all the route53 actions being used in the installer, the AWS SDK v2 does sanitize the hosted zone ID to remove any resource prefix except ChangeTagsForResource.
	hostedZoneID := sanitizeHostedZoneID(aws.ToString(res.HostedZone.Id))
	_, err = c.client.ChangeTagsForResource(ctx, &route53.ChangeTagsForResourceInput{
		ResourceType: route53types.TagResourceTypeHostedzone,
		ResourceId:   aws.String(hostedZoneID),
		AddTags:      r53Tags(tags),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to tag private hosted zone: %w", err)
	}

	// Set SOA minimum TTL
	recordSet, err := existingRecordSet(ctx, c.client, res.HostedZone.Id, input.Name, route53types.RRTypeSoa)
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
	req, err := c.client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
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

	waiter := route53.NewResourceRecordSetsChangedWaiter(c.client)

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
		resType := rs.Type
		if resName == name && resType == recordType {
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

// r53Tags converts go map to SDK route53 tag list representation.
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
