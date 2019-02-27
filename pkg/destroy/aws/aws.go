package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/version"
)

var (
	exists = struct{}{}
)

// Filter holds the key/value pairs for the tags we will be matching against.
//
// A resource matches the filter if all of the key/value pairs are in its tags.
type Filter map[string]string

// ClusterUninstaller holds the various options for the cluster we want to delete
type ClusterUninstaller struct {

	// Filters is a slice of filters for matching resources.  A
	// resources matches the whole slice if it matches any of the
	// entries.  For example:
	//
	//   filter := []map[string]string{
	//     {
	//       "a": "b",
	//       "c": "d:,
	//     },
	//     {
	//       "d": "e",
	//     },
	//   }
	//
	// will match resources with (a:b and c:d) or d:e.
	Filters []Filter // filter(s) we will be searching for
	Logger  logrus.FieldLogger
	Region  string
}

func (o *ClusterUninstaller) validate() error {
	if len(o.Filters) == 0 {
		return errors.Errorf("you must specify at least one tag filter")
	}
	return nil
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() error {
	err := o.validate()
	if err != nil {
		return err
	}

	awsConfig := &aws.Config{Region: aws.String(o.Region)}

	// Relying on appropriate AWS ENV vars (eg AWS_PROFILE, AWS_ACCESS_KEY_ID, etc)
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return err
	}
	awsSession.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Destroyer", version.Raw),
	})

	tagClients := []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI{
		resourcegroupstaggingapi.New(awsSession),
	}
	tagClientNames := map[*resourcegroupstaggingapi.ResourceGroupsTaggingAPI]string{
		tagClients[0]: o.Region,
	}
	if o.Region != "us-east-1" {
		tagClient := resourcegroupstaggingapi.New(
			awsSession, aws.NewConfig().WithRegion("us-east-1"),
		)
		tagClients = append(tagClients, tagClient)
		tagClientNames[tagClient] = "us-east-1"
	}

	deleted := map[string]struct{}{}
	iamClient := iam.New(awsSession)
	iamRoleSearch := &iamRoleSearch{
		client:  iamClient,
		filters: o.Filters,
		logger:  o.Logger,
	}
	iamUserSearch := &iamUserSearch{
		client:  iamClient,
		filters: o.Filters,
		logger:  o.Logger,
	}

	var loopError error
	for len(tagClients) > 0 || loopError != nil {
		loopError = nil
		nextTagClients := tagClients[:0]
		for _, tagClient := range tagClients {
			matched := false
			for _, filter := range o.Filters {
				o.Logger.Debugf("search for and delete matching resources by tag in %s matching %#+v", tagClientNames[tagClient], filter)
				tagFilters := make([]*resourcegroupstaggingapi.TagFilter, 0, len(filter))
				for key, value := range filter {
					tagFilters = append(tagFilters, &resourcegroupstaggingapi.TagFilter{
						Key:    aws.String(key),
						Values: []*string{aws.String(value)},
					})
				}
				err = tagClient.GetResourcesPages(
					&resourcegroupstaggingapi.GetResourcesInput{TagFilters: tagFilters},
					func(results *resourcegroupstaggingapi.GetResourcesOutput, lastPage bool) bool {
						for _, resource := range results.ResourceTagMappingList {
							arn := *resource.ResourceARN
							if _, ok := deleted[arn]; !ok {
								matched = true
								err := deleteARN(awsSession, arn, o.Logger)
								if err != nil {
									err = errors.Wrapf(err, "deleting %s", arn)
									o.Logger.Debug(err)
									continue
								}
								deleted[arn] = exists
							}
						}

						return !lastPage
					},
				)
				if err != nil {
					err = errors.Wrapf(err, "get tagged resources")
					o.Logger.Info(err)
					matched = true
					loopError = err
				}
			}

			if matched {
				nextTagClients = append(nextTagClients, tagClient)
			} else {
				o.Logger.Debugf("no deletions from %s, removing client", tagClientNames[tagClient])
			}
		}
		tagClients = nextTagClients

		o.Logger.Debug("search for IAM roles")
		arns, err := iamRoleSearch.arns()
		if err != nil {
			o.Logger.Info(err)
			loopError = err
		}

		o.Logger.Debug("search for IAM users")
		userARNs, err := iamUserSearch.arns()
		if err != nil {
			o.Logger.Info(err)
			loopError = err
		}
		arns = append(arns, userARNs...)

		if len(arns) > 0 {
			o.Logger.Debug("delete IAM roles and users")
		}
		for _, arn := range arns {
			if _, ok := deleted[arn]; !ok {
				err = deleteARN(awsSession, arn, o.Logger)
				if err != nil {
					err = errors.Wrapf(err, "deleting %s", arn)
					o.Logger.Debug(err)
					loopError = err
					continue
				}
				deleted[arn] = exists
			}
		}
	}

	return nil
}

func splitSlash(name string, input string) (base string, suffix string, err error) {
	segments := strings.SplitN(input, "/", 2)
	if len(segments) != 2 {
		return "", "", errors.Errorf("%s %q does not contain the expected slash", name, input)
	}
	return segments[0], segments[1], nil
}

func tagMatch(filters []Filter, tags map[string]string) bool {
	for _, filter := range filters {
		match := true
		for filterKey, filterValue := range filter {
			tagValue, ok := tags[filterKey]
			if !ok {
				match = false
				break
			}
			if tagValue != filterValue {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return len(filters) == 0
}

type iamRoleSearch struct {
	client    *iam.IAM
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *iamRoleSearch) arns() ([]string, error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

	arns := []string{}
	var lastError error
	err := search.client.ListRolesPages(
		&iam.ListRolesInput{},
		func(results *iam.ListRolesOutput, lastPage bool) bool {
			for _, role := range results.Roles {
				if _, ok := search.unmatched[*role.Arn]; ok {
					continue
				}

				// Unfortunately role.Tags is empty from ListRoles, so we need to query each one
				var response *iam.GetRoleOutput
				response, lastError = search.client.GetRole(&iam.GetRoleInput{RoleName: role.RoleName})
				if lastError != nil {
					if lastError.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
						search.unmatched[*role.Arn] = exists
					} else {
						lastError = errors.Wrapf(lastError, "get tags for %s", *role.Arn)
						search.logger.Info(lastError)
					}
				} else {
					role = response.Role
					tags := make(map[string]string, len(role.Tags))
					for _, tag := range role.Tags {
						tags[*tag.Key] = *tag.Value
					}
					if tagMatch(search.filters, tags) {
						arns = append(arns, *role.Arn)
					} else {
						search.unmatched[*role.Arn] = exists
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return arns, lastError
	}
	return arns, err
}

type iamUserSearch struct {
	client    *iam.IAM
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *iamUserSearch) arns() ([]string, error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

	arns := []string{}
	var lastError error
	err := search.client.ListUsersPages(
		&iam.ListUsersInput{},
		func(results *iam.ListUsersOutput, lastPage bool) bool {
			for _, user := range results.Users {
				if _, ok := search.unmatched[*user.Arn]; ok {
					continue
				}

				// Unfortunately user.Tags is empty from ListUsers, so we need to query each one
				var response *iam.GetUserOutput
				response, lastError = search.client.GetUser(&iam.GetUserInput{UserName: aws.String(*user.UserName)})
				if lastError != nil {
					if lastError.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
						search.unmatched[*user.Arn] = exists
					} else {
						lastError = errors.Wrapf(lastError, "get tags for %s", *user.Arn)
						search.logger.Info(lastError)
					}
				} else {
					user = response.User
					tags := make(map[string]string, len(user.Tags))
					for _, tag := range user.Tags {
						tags[*tag.Key] = *tag.Value
					}
					if tagMatch(search.filters, tags) {
						arns = append(arns, *user.Arn)
					} else {
						search.unmatched[*user.Arn] = exists
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return arns, lastError
	}
	return arns, err
}

// getSharedHostedZone will find the ID of the non-Terraform-managed public route53 zone given the
// Terraform-managed zone's privateID.
func getSharedHostedZone(client *route53.Route53, privateID string, logger logrus.FieldLogger) (string, error) {
	response, err := client.GetHostedZone(&route53.GetHostedZoneInput{
		Id: aws.String(privateID),
	})
	if err != nil {
		return "", err
	}

	privateName := *response.HostedZone.Name

	if response.HostedZone.Config != nil && response.HostedZone.Config.PrivateZone != nil {
		if !*response.HostedZone.Config.PrivateZone {
			return "", errors.Errorf("getSharedHostedZone requires a private ID, but was passed the public %s", privateID)
		}
	} else {
		logger.WithField("hosted zone", privateName).Warn("could not determine whether hosted zone is private")
	}

	domain := privateName
	parents := []string{domain}
	for {
		idx := strings.Index(domain, ".")
		if idx == -1 {
			break
		}
		if len(domain[idx+1:]) > 0 {
			parents = append(parents, domain[idx+1:])
		}
		domain = domain[idx+1:]
	}

	for _, p := range parents {
		sZone, err := findPublicRoute53(client, p, logger)
		if err != nil {
			return "", err
		}
		if sZone != "" {
			return sZone, nil
		}
	}
	return "", nil
}

// findPublicRoute53 finds a public route53 zone matching the dnsName.
// It returns "", when no public route53 zone could be found.
func findPublicRoute53(client *route53.Route53, dnsName string, logger logrus.FieldLogger) (string, error) {
	request := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(dnsName),
	}
	for i := 0; true; i++ {
		logger.Debugf("listing AWS hosted zones %q (page %d)", dnsName, i)
		list, err := client.ListHostedZonesByName(request)
		if err != nil {
			return "", err
		}

		for _, zone := range list.HostedZones {
			if *zone.Name != dnsName {
				// No name after this can match dnsName
				return "", nil
			}
			if zone.Config == nil || zone.Config.PrivateZone == nil {
				logger.WithField("hosted zone", *zone.Name).Warn("could not determine whether hosted zone is private")
				continue
			}
			if !*zone.Config.PrivateZone {
				return *zone.Id, nil
			}
		}

		if *list.IsTruncated && *list.NextDNSName == *request.DNSName {
			request.HostedZoneId = list.NextHostedZoneId
			continue
		}

		break
	}
	return "", nil
}

func deleteARN(session *session.Session, arnString string, logger logrus.FieldLogger) error {
	logger = logger.WithField("arn", arnString)

	parsed, err := arn.Parse(arnString)
	if err != nil {
		return err
	}

	switch parsed.Service {
	case "ec2":
		return deleteEC2(session, parsed, logger)
	case "elasticloadbalancing":
		return deleteElasticLoadBalancing(session, parsed, logger)
	case "iam":
		return deleteIAM(session, parsed, logger)
	case "route53":
		return deleteRoute53(session, parsed, logger)
	case "s3":
		return deleteS3(session, parsed, logger)
	default:
		return errors.Errorf("unrecognized ARN service %s (%s)", parsed.Service, arnString)
	}
}

func deleteEC2(session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := ec2.New(session)

	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "dhcp-options":
		return deleteEC2DHCPOptions(client, id, logger)
	case "elastic-ip":
		return deleteEC2ElasticIP(client, id, logger)
	case "image":
		return deleteEC2Image(client, id, logger)
	case "instance":
		return deleteEC2Instance(client, iam.New(session), id, logger)
	case "internet-gateway":
		return deleteEC2InternetGateway(client, id, logger)
	case "natgateway":
		return deleteEC2NATGateway(client, id, logger)
	case "route-table":
		return deleteEC2RouteTable(client, id, logger)
	case "security-group":
		return deleteEC2SecurityGroup(client, id, logger)
	case "subnet":
		return deleteEC2Subnet(client, id, logger)
	case "volume":
		return deleteEC2Volume(client, id, logger)
	case "vpc":
		return deleteEC2VPC(client, elb.New(session), elbv2.New(session), id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteEC2DHCPOptions(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteDhcpOptions(&ec2.DeleteDhcpOptionsInput{
		DhcpOptionsId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidDhcpOptionsID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2Image(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeregisterImage(&ec2.DeregisterImageInput{
		ImageId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAMIID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2ElasticIP(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAllocationID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Released")
	return nil
}

func deleteEC2Instance(ec2Client *ec2.EC2, iamClient *iam.IAM, id string, logger logrus.FieldLogger) error {
	response, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(id)},

		// only fetch instances in 'running|pending' state since 'terminated' ones take a while to really get cleaned up
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidInstanceID.NotFound" {
			return nil
		}
		return err
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			if instance.IamInstanceProfile != nil {
				parsed, err := arn.Parse(*instance.IamInstanceProfile.Arn)
				if err != nil {
					return errors.Wrap(err, "parse ARN for IAM instance profile")
				}

				err = deleteIAMInstanceProfile(iamClient, parsed, logger.WithField("IAM instance profile", parsed.String()))
				if err != nil {
					return errors.Wrapf(err, "deleting %s", parsed.String())
				}
			}

			_, err := ec2Client.TerminateInstances(&ec2.TerminateInstancesInput{
				InstanceIds: []*string{instance.InstanceId},
			})
			if err != nil {
				return err
			}

			logger.Info("Deleted")
		}
	}

	return nil
}

func deleteEC2InternetGateway(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(id)},
	})
	if err != nil {
		return err
	}

	for _, gateway := range response.InternetGateways {
		for _, vpc := range gateway.Attachments {
			if vpc.VpcId == nil {
				logger.Warn("gateway does not have a VPC ID")
				continue
			}
			_, err := client.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
				InternetGatewayId: gateway.InternetGatewayId,
				VpcId:             vpc.VpcId,
			})
			if err == nil {
				logger.WithField("vpc", *vpc.VpcId).Debug("Detached")
			} else if err.(awserr.Error).Code() != "Gateway.NotAttached" {
				return errors.Wrapf(err, "detaching from %s", *vpc.VpcId)
			}
		}
	}

	_, err = client.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NATGateway(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteNatGateway(&ec2.DeleteNatGatewayInput{
		NatGatewayId: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTable(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidRouteTableID.NotFound" {
			return nil
		}
		return err
	}

	for _, table := range response.RouteTables {
		err = deleteEC2RouteTableObject(client, table, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEC2RouteTableObject(client *ec2.EC2, table *ec2.RouteTable, logger logrus.FieldLogger) error {
	hasMain := false
	for _, association := range table.Associations {
		if *association.Main {
			// can't remove the 'Main' association
			hasMain = true
			continue
		}
		_, err := client.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
			AssociationId: association.RouteTableAssociationId,
		})
		if err != nil {
			return errors.Wrapf(err, "dissociating %s", *association.RouteTableAssociationId)
		}
		logger.WithField("id", *association.RouteTableAssociationId).Info("Disassociated")
	}

	if hasMain {
		// can't delete route table with the 'Main' association
		// it will get cleaned up as part of deleting the VPC
		return nil
	}

	_, err := client.DeleteRouteTable(&ec2.DeleteRouteTableInput{
		RouteTableId: table.RouteTableId,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTablesByVPC(client *ec2.EC2, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeRouteTablesPages(
		&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeRouteTablesOutput, lastPage bool) bool {
			for _, table := range results.RouteTables {
				lastError := deleteEC2RouteTableObject(client, table, logger.WithField("table", *table.RouteTableId))
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting EC2 route table %s", *table.RouteTableId)
					logger.Info(lastError)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2SecurityGroup(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(id)},
	})
	if err != nil {
		return err
	}

	for _, group := range response.SecurityGroups {
		if len(group.IpPermissions) > 0 {
			_, err := client.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
				GroupId:       group.GroupId,
				IpPermissions: group.IpPermissions,
			})
			if err != nil {
				return errors.Wrap(err, "revoking ingress permissions")
			}
			logger.Debug("Revoked ingress permissions")
		}

		if len(group.IpPermissionsEgress) > 0 {
			_, err := client.RevokeSecurityGroupEgress(&ec2.RevokeSecurityGroupEgressInput{
				GroupId:       group.GroupId,
				IpPermissions: group.IpPermissionsEgress,
			})
			if err != nil {
				return errors.Wrap(err, "revoking egress permissions")
			}
			logger.Debug("Revoked egress permissions")
		}
	}

	_, err = client.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2Subnet(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidSubnetID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2Volume(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVolume(&ec2.DeleteVolumeInput{
		VolumeId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVolume.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPC(ec2Client *ec2.EC2, elbClient *elb.ELB, elbv2Client *elbv2.ELBV2, id string, logger logrus.FieldLogger) error {
	// first delete any Load Balancers under this VPC (not all of them are tagged)
	v1lbError := deleteElasticLoadBalancerClassicByVPC(elbClient, id, logger)
	v2lbError := deleteElasticLoadBalancerV2ByVPC(elbv2Client, id, logger)
	if v1lbError != nil {
		if v2lbError != nil {
			logger.Info(v2lbError)
		}
		return v1lbError
	} else if v2lbError != nil {
		return v2lbError
	}

	// next delete any VPC endpoints associated with the VPC (they are not taggable)
	err := deleteEC2VPCEndpointsByVPC(ec2Client, id, logger)
	if err != nil {
		return err
	}

	// next delete route tables associated with the VPC (not all of them are tagged)
	err = deleteEC2RouteTablesByVPC(ec2Client, id, logger)
	if err != nil {
		return err
	}

	_, err = ec2Client.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpoint(client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVpcEndpoints(&ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []*string{aws.String(id)},
	})
	if err != nil {
		return errors.Wrapf(err, "cannot delete VPC endpoint %s", id)
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpointsByVPC(client *ec2.EC2, vpc string, logger logrus.FieldLogger) error {
	response, err := client.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpc)},
			},
		},
	})

	if err != nil {
		return err
	}

	for _, endpoint := range response.VpcEndpoints {
		err := deleteEC2VPCEndpoint(client, *endpoint.VpcEndpointId, logger.WithField("VPC endpoint", *endpoint.VpcEndpointId))
		if err != nil {
			if err.(awserr.Error).Code() == "InvalidVpcID.NotFound" {
				return nil
			}
			return err
		}
	}

	return nil
}

func deleteElasticLoadBalancing(session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "loadbalancer":
		segments := strings.SplitN(id, "/", 2)
		if len(segments) == 1 {
			return deleteElasticLoadBalancerClassic(elb.New(session), id, logger)
		} else if len(segments) != 2 {
			return errors.Errorf("cannot parse subresource %q into {subtype}/{id}", id)
		}
		subtype := segments[0]
		id = segments[1]
		switch subtype {
		case "net":
			return deleteElasticLoadBalancerV2(elbv2.New(session), arn, logger)
		default:
			return errors.Errorf("unrecognized elastic load balancing resource subtype %s", subtype)
		}
	case "targetgroup":
		return deleteElasticLoadBalancerTargetGroup(elbv2.New(session), arn, logger)
	default:
		return errors.Errorf("unrecognized elastic load balancing resource type %s", resourceType)
	}
}

func deleteElasticLoadBalancerClassic(client *elb.ELB, name string, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerClassicByVPC(client *elb.ELB, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPages(
		&elb.DescribeLoadBalancersInput{},
		func(results *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, lb := range results.LoadBalancerDescriptions {
				lbLogger := logger.WithField("classic load balancer", *lb.LoadBalancerName)

				if lb.VPCId == nil {
					lbLogger.Warn("classic load balancer does not have a VPC ID so could not determine whether it should be deleted")
					continue
				}

				if *lb.VPCId != vpc {
					continue
				}

				lastError = deleteElasticLoadBalancerClassic(client, *lb.LoadBalancerName, lbLogger)
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting classic load balancer %s", *lb.LoadBalancerName)
					logger.Info(lastError)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteElasticLoadBalancerTargetGroup(client *elbv2.ELBV2, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn.String()),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerTargetGroupsByVPC(client *elbv2.ELBV2, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeTargetGroupsPages(
		&elbv2.DescribeTargetGroupsInput{},
		func(results *elbv2.DescribeTargetGroupsOutput, lastPage bool) bool {
			for _, group := range results.TargetGroups {
				if group.VpcId == nil {
					logger.WithField("target group", *group.TargetGroupArn).Warn("load balancer target group does not have a VPC ID so could not determine whether it should be deleted")
					continue
				}

				if *group.VpcId != vpc {
					continue
				}

				var parsed arn.ARN
				parsed, lastError = arn.Parse(*group.TargetGroupArn)
				if lastError != nil {
					lastError = errors.Wrap(lastError, "parse ARN for target group")
					logger.Info(lastError)
					continue
				}

				lastError = deleteElasticLoadBalancerTargetGroup(client, parsed, logger.WithField("target group", parsed.Resource))
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting %s", parsed.String())
					logger.Info(lastError)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteElasticLoadBalancerV2(client *elbv2.ELBV2, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn.String()),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerV2ByVPC(client *elbv2.ELBV2, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPages(
		&elbv2.DescribeLoadBalancersInput{},
		func(results *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, lb := range results.LoadBalancers {
				if lb.VpcId == nil {
					logger.WithField("load balancer", *lb.LoadBalancerArn).Warn("load balancer does not have a VPC ID so could not determine whether it should be deleted")
					continue
				}

				if *lb.VpcId != vpc {
					continue
				}

				var parsed arn.ARN
				parsed, lastError = arn.Parse(*lb.LoadBalancerArn)
				if lastError != nil {
					lastError = errors.Wrap(lastError, "parse ARN for load balancer")
					logger.Info(lastError)
					continue
				}

				lastError = deleteElasticLoadBalancerV2(client, parsed, logger.WithField("load balancer", parsed.Resource))
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting %s", parsed.String())
					logger.Info(lastError)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteIAM(session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := iam.New(session)

	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "instance-profile":
		return deleteIAMInstanceProfile(client, arn, logger)
	case "role":
		return deleteIAMRole(client, arn, logger)
	case "user":
		return deleteIAMUser(client, id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteIAMInstanceProfile(client *iam.IAM, profileARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", profileARN.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("name", name)

	if resourceType != "instance-profile" {
		return errors.Errorf("%s ARN passed to deleteIAMInstanceProfile: %s", resourceType, profileARN.String())
	}

	response, err := client.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: &name,
	})
	if err != nil {
		if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
			return nil
		}
		return err
	}
	profile := response.InstanceProfile

	for _, role := range profile.Roles {
		_, err = client.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: profile.InstanceProfileName,
			RoleName:            role.RoleName,
		})
		if err != nil {
			return errors.Wrapf(err, "dissociating %s", *role.RoleName)
		}
		logger.WithField("role", *role.RoleName).Info("Disassociated")
	}

	_, err = client.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: profile.InstanceProfileName,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteIAMRole(client *iam.IAM, roleARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", roleARN.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("name", name)

	if resourceType != "role" {
		return errors.Errorf("%s ARN passed to deleteIAMRole: %s", resourceType, roleARN.String())
	}

	var lastError error
	err = client.ListRolePoliciesPages(
		&iam.ListRolePoliciesInput{RoleName: &name},
		func(results *iam.ListRolePoliciesOutput, lastPage bool) bool {
			for _, policy := range results.PolicyNames {
				_, lastError = client.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
					RoleName:   &name,
					PolicyName: policy,
				})
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting IAM role policy %s", *policy)
					logger.Info(lastError)
				}
				logger.WithField("policy", *policy).Info("Deleted")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM role policies")
	}

	err = client.ListInstanceProfilesForRolePages(
		&iam.ListInstanceProfilesForRoleInput{RoleName: &name},
		func(results *iam.ListInstanceProfilesForRoleOutput, lastPage bool) bool {
			for _, profile := range results.InstanceProfiles {
				parsed, lastError := arn.Parse(*profile.Arn)
				if lastError != nil {
					lastError = errors.Wrap(lastError, "parse ARN for IAM instance profile")
					logger.Info(lastError)
					continue
				}

				lastError = deleteIAMInstanceProfile(client, parsed, logger.WithField("IAM instance profile", parsed.String()))
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting %s", parsed.String())
					logger.Info(lastError)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM instance profiles")
	}

	_, err = client.DeleteRole(&iam.DeleteRoleInput{RoleName: &name})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteIAMUser(client *iam.IAM, id string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.ListUserPoliciesPages(
		&iam.ListUserPoliciesInput{UserName: &id},
		func(results *iam.ListUserPoliciesOutput, lastPage bool) bool {
			for _, policy := range results.PolicyNames {
				_, lastError = client.DeleteUserPolicy(&iam.DeleteUserPolicyInput{
					UserName:   &id,
					PolicyName: policy,
				})
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting IAM user policy %s", *policy)
					logger.Info(lastError)
				}
				logger.WithField("policy", *policy).Info("Deleted")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM user policies")
	}

	err = client.ListAccessKeysPages(
		&iam.ListAccessKeysInput{UserName: &id},
		func(results *iam.ListAccessKeysOutput, lastPage bool) bool {
			for _, key := range results.AccessKeyMetadata {
				_, lastError := client.DeleteAccessKey(&iam.DeleteAccessKeyInput{
					UserName:    &id,
					AccessKeyId: key.AccessKeyId,
				})
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting IAM access key %s", *key.AccessKeyId)
					logger.Info(lastError)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM access keys")
	}

	_, err = client.DeleteUser(&iam.DeleteUserInput{
		UserName: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteRoute53(session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	if resourceType != "hostedzone" {
		return errors.Errorf("unrecognized Route 53 resource type %s", resourceType)
	}

	client := route53.New(session)

	sharedZoneID, err := getSharedHostedZone(client, id, logger)
	if err != nil {
		return err
	}

	recordSetKey := func(recordSet *route53.ResourceRecordSet) string {
		return fmt.Sprintf("%s %s", *recordSet.Type, *recordSet.Name)
	}

	sharedEntries := map[string]*route53.ResourceRecordSet{}
	if len(sharedZoneID) != 0 {
		err = client.ListResourceRecordSetsPages(
			&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(sharedZoneID)},
			func(results *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
				for _, recordSet := range results.ResourceRecordSets {
					key := recordSetKey(recordSet)
					sharedEntries[key] = recordSet
				}

				return !lastPage
			},
		)
		if err != nil {
			return err
		}
	} else {
		logger.Debug("shared public zone not found")
	}

	var lastError error
	err = client.ListResourceRecordSetsPages(
		&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(id)},
		func(results *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			for _, recordSet := range results.ResourceRecordSets {
				if *recordSet.Type == "SOA" || *recordSet.Type == "NS" {
					// can't delete SOA and NS types
					continue
				}
				key := recordSetKey(recordSet)
				if sharedEntry, ok := sharedEntries[key]; ok {
					lastError = deleteRoute53RecordSet(client, sharedZoneID, sharedEntry, logger.WithField("public zone", sharedZoneID))
					if lastError != nil {
						lastError = errors.Wrapf(lastError, "deleting public zone %s", sharedZoneID)
						logger.Info(lastError)
					}
				}

				lastError = deleteRoute53RecordSet(client, id, recordSet, logger)
				if lastError != nil {
					lastError = errors.Wrapf(lastError, "deleting record set %#+v from zone %s", recordSet, id)
					logger.Info(lastError)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return err
	}

	_, err = client.DeleteHostedZone(&route53.DeleteHostedZoneInput{
		Id: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteRoute53RecordSet(client *route53.Route53, zoneID string, recordSet *route53.ResourceRecordSet, logger logrus.FieldLogger) error {
	logger = logger.WithField("record set", fmt.Sprintf("%s %s", *recordSet.Type, *recordSet.Name))
	_, err := client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String("DELETE"),
					ResourceRecordSet: recordSet,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteS3(session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := s3.New(session)

	iter := s3manager.NewDeleteListIterator(client, &s3.ListObjectsInput{
		Bucket: aws.String(arn.Resource),
	})
	err := s3manager.NewBatchDeleteWithClient(client).Delete(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}
	logger.Debug("Emptied")

	_, err = client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(arn.Resource),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}
