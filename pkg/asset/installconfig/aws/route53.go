package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awss "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

//go:generate mockgen -source=./route53.go -destination=mock/awsroute53_generated.go -package=mock

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
		return nil, errors.Wrapf(err, "could not get hosted zone: %s", hostedZone)
	}
	return hostedZoneOutput, nil
}

// ValidateZoneRecords Attempts to validate each of the candidate HostedZones against the Config
func (c *Client) ValidateZoneRecords(zone *route53.HostedZone, zoneName string, zonePath *field.Path, ic *types.InstallConfig, cfg *aws.Config) field.ErrorList {
	allErrs := field.ErrorList{}

	problematicRecords, err := c.GetSubDomainDNSRecords(zone, ic, cfg)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(zonePath,
			errors.Wrapf(err, "could not list record sets for domain %q", zoneName)))
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
		return nil, errors.Errorf("hosted zone domain %q is not a parent of the cluster domain %q", *hostedZone.Name, dottedClusterDomain)
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
		return nil, errors.Wrapf(err, "could not find public zone: %s", baseDomainName)
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
