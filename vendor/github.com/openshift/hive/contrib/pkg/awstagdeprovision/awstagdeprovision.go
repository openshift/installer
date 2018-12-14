/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awstagdeprovision

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	log "github.com/sirupsen/logrus"
)

const (
	secondsToSleep = 10
)

// AWSFilter holds the key/value pairs for the tags we will be matching against.
//
// A resource matches the filter if all of the key/value pairs are in its tags.
type AWSFilter map[string]string

// awsObjectWithTags is a generic way to represent an AWS object and its tags so that
// filtering objects client-side can be done in a generic way
type awsObjectWithTags struct {
	Name string
	Tags map[string]string
}

// deleteFunc type is the interface a function needs to implement to be called as a goroutine.
// The (bool, error) return type mimics wait.ExponentialBackoff where the bool indicates successful
// completion, and the error is for unrecoverable errors.
type deleteFunc func(awsClient *session.Session, filters AWSFilter, clusterName string, logger log.FieldLogger) (bool, error)

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
	Filters     []AWSFilter // filter(s) we will be searching for
	Logger      log.FieldLogger
	LogLevel    string
	Region      string
	ClusterName string
}

func (o *ClusterUninstaller) validate() error {
	if len(o.Filters) == 0 {
		return fmt.Errorf("you must specify at least one tag filter")
	}
	if len(o.ClusterName) == 0 {
		return fmt.Errorf("you must specify cluster-name")
	}
	return nil
}

// populateDeleteFuncs is the list of functions that will be launched as goroutines
func populateDeleteFuncs(funcs map[string]deleteFunc) {
	funcs["deleteVPCs"] = deleteVPCs
	funcs["deleteEIPs"] = deleteEIPs
	funcs["deleteNATGateways"] = deleteNATGateways
	funcs["deleteInstances"] = deleteInstances
	funcs["deleteIAMresources"] = deleteIAMresources
	funcs["deleteSecurityGroups"] = deleteSecurityGroups
	funcs["deleteInternetGateways"] = deleteInternetGateways
	funcs["deleteSubnets"] = deleteSubnets
	funcs["deleteS3Buckets"] = deleteS3Buckets
	funcs["deleteRoute53"] = deleteRoute53
	funcs["deletePVs"] = deletePVs
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() error {
	err := o.validate()
	if err != nil {
		return err
	}
	deleteFuncs := map[string]deleteFunc{}
	populateDeleteFuncs(deleteFuncs)
	returnChannel := make(chan string)

	awsSession, err := getAWSSession(o.Region)
	if err != nil {
		return err
	}

	// launch goroutines
	goroutines := 0
	for name, function := range deleteFuncs {
		for _, filter := range o.Filters {
			go deleteRunner(name, function, awsSession, filter, o.ClusterName, o.Logger, returnChannel)
			goroutines++
		}
	}

	// wait for them to finish
	for goroutines > 0 {
		select {
		case res := <-returnChannel:
			goroutines--
			o.Logger.Debugf("goroutine %v complete (%d left)", res, goroutines)
		}
	}

	return nil
}

func deleteRunner(deleteFuncName string, dFunction deleteFunc, awsSession *session.Session, filters AWSFilter, clusterName string, logger log.FieldLogger, channel chan string) {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 10,
		Factor:   1.3,
		Steps:    100,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return dFunction(awsSession, filters, clusterName, logger)
	})

	if err != nil {
		logger.Fatalf("Unrecoverable error/timed out: %v", err)
		os.Exit(1)
	}

	// record that the goroutine has run to completion
	channel <- deleteFuncName
	return
}

func getAWSSession(region string) (*session.Session, error) {
	awsConfig := &aws.Config{Region: aws.String(region)}

	// Relying on appropriate AWS ENV vars (eg AWS_PROFILE, AWS_ACCESS_KEY_ID, etc)
	s, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func createEC2Filters(filters AWSFilter) []*ec2.Filter {
	AWSFilter := []*ec2.Filter{}
	for key, val := range filters {
		AWSFilter = append(AWSFilter, &ec2.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", key)),
			Values: []*string{aws.String(val)},
		})
	}

	return AWSFilter
}

// tagsToMap takes various types of AWS-object tags and returns a map-representation
func tagsToMap(tags interface{}) (map[string]string, error) {
	x := map[string]string{}

	switch v := tags.(type) {
	case []*autoscaling.TagDescription:
		for _, tag := range v {
			x[*tag.Key] = *tag.Value
		}
	case *elb.TagDescription:
		for _, tag := range v.Tags {
			x[*tag.Key] = *tag.Value
		}
	case []*s3.Tag:
		for _, tag := range v {
			x[*tag.Key] = *tag.Value
		}
	case []*route53.Tag:
		for _, tag := range v {
			x[*tag.Key] = *tag.Value
		}
	default:
		return x, fmt.Errorf("unable to convert type: %v", v)
	}

	return x, nil
}

// filterLBsByVPC will find all the load balancers in the provided list that are under the provided VPC
func filterLBsByVPC(lbs []*elb.LoadBalancerDescription, vpc *ec2.Vpc, logger log.FieldLogger) []*elb.LoadBalancerDescription {
	filteredLBs := []*elb.LoadBalancerDescription{}

	for _, lb := range lbs {
		if *lb.VPCId == *vpc.VpcId {
			filteredLBs = append(filteredLBs, lb)
		}
	}

	return filteredLBs
}

// deleteLBs finds all load balancers under the provided VPC and attempts to delete them
// returns bool representing whether it has completed its work (ie no LBs left to delete)
func deleteLBs(vpc *ec2.Vpc, awsSession *session.Session, logger log.FieldLogger) bool {
	logger.Debugf("Deleting load balancers (%s)", *vpc.VpcId)
	defer logger.Debugf("Exiting deleting load balancers (%s)", *vpc.VpcId)
	elbClient := elb.New(awsSession)

	describeLoadBalancersInput := elb.DescribeLoadBalancersInput{}
	results, err := elbClient.DescribeLoadBalancers(&describeLoadBalancersInput)
	if err != nil {
		logger.Errorf("Error listing load balancers: %v", err)
		return false
	}

	filteredLBs := filterLBsByVPC(results.LoadBalancerDescriptions, vpc, logger)
	logger.Debugf("from %d total load balancers, %d scheduled for deletion", len(results.LoadBalancerDescriptions), len(filteredLBs))

	if len(filteredLBs) == 0 {
		// no items left to delete
		return true
	}

	for _, lb := range filteredLBs {
		logger.Debugf("Deleting load balancer: %v", *lb.LoadBalancerName)
		_, err := elbClient.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
			LoadBalancerName: lb.LoadBalancerName,
		})
		if err != nil {
			logger.Debugf("Error deleting load balancer %v: %v", *lb.LoadBalancerName, err)
		} else {
			logger.WithField("name", *lb.LoadBalancerName).Info("Deleted load balancer")
		}
	}
	return false
}

// filterV2LBsByVPC will find all the load balancers in the provided list that are under the provided VPC
func filterV2LBsByVPC(lbs []*elbv2.LoadBalancer, vpc *ec2.Vpc, logger log.FieldLogger) []*elbv2.LoadBalancer {
	filteredLBs := []*elbv2.LoadBalancer{}

	for _, lb := range lbs {
		if *lb.VpcId == *vpc.VpcId {
			filteredLBs = append(filteredLBs, lb)
		}
	}

	return filteredLBs
}

func deleteV2LBs(vpc *ec2.Vpc, awsSession *session.Session, logger log.FieldLogger) bool {
	logger.Debugf("Deleting V2 load balancers (%s)", *vpc.VpcId)
	defer logger.Debugf("Exiting deleting V2 load balancers (%s)", *vpc.VpcId)
	elbv2Client := elbv2.New(awsSession)

	total := 0
	filteredLBs := []*elbv2.LoadBalancer{}
	if err := elbv2Client.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{}, func(results *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
		total += len(results.LoadBalancers)
		filteredLBs = append(filteredLBs, filterV2LBsByVPC(results.LoadBalancers, vpc, logger)...)
		return lastPage
	}); err != nil {
		logger.Errorf("Error listing V2 load balancers: %v", err)
		return false
	}
	logger.Debugf("from %d total V2 load balancers, %d scheduled for deletion", total, len(filteredLBs))

	if len(filteredLBs) == 0 {
		// no items left to delete
		// see if we can delete target groups.
		return deleteTargetGroups(vpc, awsSession, logger)
	}

	for _, lb := range filteredLBs {
		logger.Debugf("Deleting V2 load balancer: %v", *lb.LoadBalancerName)
		_, err := elbv2Client.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})
		if err != nil {
			logger.Debugf("Error deleting V2 load balancer %v: %v", *lb.LoadBalancerName, err)
		} else {
			logger.WithField("name", *lb.LoadBalancerName).Info("Deleted load balancer")
		}
	}
	// cleanup target groups
	return deleteTargetGroups(vpc, awsSession, logger)
}

// filterTargetGroupsByVPC will find all the target groups in the provided list that are under the provided VPC
func filterTargetGroupsByVPC(tgs []*elbv2.TargetGroup, vpc *ec2.Vpc, logger log.FieldLogger) []*elbv2.TargetGroup {
	filteredTGs := []*elbv2.TargetGroup{}

	for _, tg := range tgs {
		if *tg.VpcId == *vpc.VpcId {
			filteredTGs = append(filteredTGs, tg)
		}
	}

	return filteredTGs
}

func deleteTargetGroups(vpc *ec2.Vpc, awsSession *session.Session, logger log.FieldLogger) bool {
	logger.Debugf("Deleting target groups (%s)", *vpc.VpcId)
	defer logger.Debugf("Exiting deleting target groups (%s)", *vpc.VpcId)
	elbv2Client := elbv2.New(awsSession)

	total := 0
	filteredTGs := []*elbv2.TargetGroup{}
	if err := elbv2Client.DescribeTargetGroupsPages(&elbv2.DescribeTargetGroupsInput{}, func(results *elbv2.DescribeTargetGroupsOutput, lastPage bool) bool {
		total += len(results.TargetGroups)
		filteredTGs = append(filteredTGs, filterTargetGroupsByVPC(results.TargetGroups, vpc, logger)...)
		return lastPage
	}); err != nil {
		logger.Errorf("Error listing target groups: %v", err)
		return false
	}
	logger.Debugf("from %d total target groups, %d scheduled for deletion", total, len(filteredTGs))

	if len(filteredTGs) == 0 {
		// no items left to delete
		return true
	}

	for _, tg := range filteredTGs {
		logger.Debugf("Deleting target groups: %v", *tg.TargetGroupName)
		_, err := elbv2Client.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
			TargetGroupArn: tg.TargetGroupArn,
		})
		if err != nil {
			logger.Debugf("Error deleting target groups %v: %v", *tg.TargetGroupName, err)
		} else {
			logger.WithField("name", *tg.TargetGroupName).Info("Deleted target group")
		}
	}
	return false
}

// rtHasMainAssociation will check whether a given route table has an association marked 'Main'
func rtHasMainAssociation(rt *ec2.RouteTable) bool {
	for _, association := range rt.Associations {
		if *association.Main == true {
			return true
		}
	}
	return false
}

// deleteVPCEndpoints will find all VPC endpoints associated with the passed in VPC and attempt to delete them
func deleteVPCEndpoints(vpc *ec2.Vpc, ec2Client *ec2.EC2, logger log.FieldLogger) error {
	describeEndpointsInput := ec2.DescribeVpcEndpointsInput{}
	describeEndpointsInput.Filters = []*ec2.Filter{
		{
			Name:   aws.String("vpc-id"),
			Values: []*string{vpc.VpcId},
		},
	}

	results, err := ec2Client.DescribeVpcEndpoints(&describeEndpointsInput)
	if err != nil {
		logger.Debugf("error describing VPC endpoints: %v", err)
		return err
	}
	for _, ep := range results.VpcEndpoints {
		_, err := ec2Client.DeleteVpcEndpoints(&ec2.DeleteVpcEndpointsInput{
			VpcEndpointIds: []*string{ep.VpcEndpointId},
		})
		if err != nil {
			logger.Debugf("error deleting VPC endpoint: %v", err)
			return err
		}
		logger.WithField("id", *ep.VpcEndpointId).Info("Deleted VPC endpoint")
	}
	return nil
}

// deleteRouteTablesWithVPC will attempt to delete all route tables associated with a given VPC
func deleteRouteTablesWithVPC(vpc *ec2.Vpc, ec2Client *ec2.EC2, logger log.FieldLogger) error {
	var anyError error
	describeRouteTablesInput := ec2.DescribeRouteTablesInput{}
	describeRouteTablesInput.Filters = []*ec2.Filter{
		{
			Name:   aws.String("vpc-id"),
			Values: []*string{vpc.VpcId},
		},
	}

	results, err := ec2Client.DescribeRouteTables(&describeRouteTablesInput)
	if err != nil {
		logger.Debugf("error describing route tables: %v", err)
		return err
	}
	for _, rt := range results.RouteTables {
		err := disassociateRouteTable(rt, ec2Client, logger)
		if err != nil {
			logger.Debugf("error disassociating from route table: %v", err)
			return err
		}

		if rtHasMainAssociation(rt) {
			// can't delete route table with the 'Main' association
			// it will get cleaned up as part of deleting the VPC
			continue
		}
		// there is a certain order that route tables need to be deleted, just try to delete
		// all of them and eventually they will all be deleted
		logger.Debugf("deleting route table: %v", *rt.RouteTableId)
		_, err = ec2Client.DeleteRouteTable(&ec2.DeleteRouteTableInput{
			RouteTableId: rt.RouteTableId,
		})
		if err != nil {
			logger.Debugf("error deleting route table: %v", err)
			anyError = err
		} else {
			logger.WithField("id", *rt.RouteTableId).Info("Deleted route table")
		}
	}

	return anyError
}

// deleteVPCs will delete any VPCs that match the provided filters/tags
func deleteVPCs(awsSession *session.Session, filters AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting VPCs (%s)", filters)
	defer logger.Debugf("Exiting deleting VPCs (%s)", filters)
	ec2Client := getEC2Client(awsSession)

	describeVpcsInput := ec2.DescribeVpcsInput{}
	describeVpcsInput.Filters = createEC2Filters(filters)
	for {
		results, err := ec2Client.DescribeVpcs(&describeVpcsInput)
		if err != nil {
			logger.Errorf("Error listing VPCs: %v", err)
			return false, nil
		}

		if len(results.Vpcs) == 0 {
			break
		}

		for _, vpc := range results.Vpcs {
			// first delete any Load Balancers under this VPC (not all of them are tagged)
			v1lbcomplete := deleteLBs(vpc, awsSession, logger)
			v2lbcomplete := deleteV2LBs(vpc, awsSession, logger)
			if !v1lbcomplete || !v2lbcomplete {
				logger.Debugf("not finished deleting load balancers, will need to retry")
				return false, nil
			}

			// next delete any VPC endpoints associated with the VPC (they are not taggable)
			err := deleteVPCEndpoints(vpc, ec2Client, logger)
			if err != nil {
				logger.Debugf("error deleting VPC endpoint: %v", err)
				return false, nil
			}

			// next delete route tables associated with the VPC (not all of them are tagged)
			err = deleteRouteTablesWithVPC(vpc, ec2Client, logger)
			if err != nil {
				logger.Debugf("error deleting route tables: %v", err)
				return false, nil
			}

			logger.Debugf("deleting VPC: %v", *vpc.VpcId)
			_, err = ec2Client.DeleteVpc(&ec2.DeleteVpcInput{
				VpcId: vpc.VpcId,
			})
			if err != nil {
				logger.Debugf("error deleting VPC %v: %v", *vpc.VpcId, err)
				return false, nil
			}

			logger.WithField("id", *vpc.VpcId).Info("Deleted VPC")
		}

		return false, nil
	}

	return true, nil
}

// getEC2Client is just a wrapper for creating an EC2 client
func getEC2Client(awsSession *session.Session) *ec2.EC2 {
	return ec2.New(awsSession)
}

// deleteNATGateways will attempt to delete all NAT Gateways that match the provided filters
func deleteNATGateways(awsSession *session.Session, filters AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {

	logger.Debugf("Deleting NAT Gateways (%s)", filters)
	defer logger.Debugf("Exiting deleting NAT Gateways (%s)", filters)

	ec2Client := getEC2Client(awsSession)
	describeNatGatewaysInput := ec2.DescribeNatGatewaysInput{}
	describeNatGatewaysInput.Filter = createEC2Filters(filters)

	// NAT Gateways take a while to really disappear so only find the ones not already being deleted
	describeNatGatewaysInput.Filter = append(describeNatGatewaysInput.Filter, &ec2.Filter{
		Name:   aws.String("state"),
		Values: []*string{aws.String("available")},
	})

	for {
		results, err := ec2Client.DescribeNatGateways(&describeNatGatewaysInput)
		if err != nil {
			logger.Debugf("error listing NAT gateways: %v", err)
			return false, nil
		}

		if len(results.NatGateways) == 0 {
			break
		}

		for _, nat := range results.NatGateways {
			logger.Debugf("deleting NAT Gateway: %v", *nat.NatGatewayId)
			_, err := ec2Client.DeleteNatGateway(&ec2.DeleteNatGatewayInput{
				NatGatewayId: nat.NatGatewayId,
			})
			if err != nil {
				logger.Debugf("error deleting NAT gateway: %v", err)
				continue
			} else {
				logger.WithField("id", *nat.NatGatewayId).Info("Deleted NAT Gateway")
			}
		}

		return false, nil
	}

	return true, nil
}

// deleteNetworkIface will attempt to delete a specific network interface
func deleteNetworkIface(iface *string, ec2Client *ec2.EC2, logger log.FieldLogger) error {

	result, err := ec2Client.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{iface},
	})
	if err != nil {
		logger.Debugf("error listing network interface: %v", err)
		return err
	}

	if len(result.NetworkInterfaces) == 0 {
		// must have already been deleted
		return nil
	}

	for _, i := range result.NetworkInterfaces {
		logger.Debugf("deleting network interface: %v", *i.NetworkInterfaceId)
		_, err := ec2Client.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: i.NetworkInterfaceId,
		})
		if err != nil {
			logger.Debugf("error deleting network iface: %v", err)
			return err
		}

		logger.WithField("id", *i.NetworkInterfaceId).Info("Deleted network interface")
	}

	return nil
}

// deleteEIPs will attempt to delete any elastic IPs matching the provided filters
func deleteEIPs(awsSession *session.Session, filters AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting EIPs (%s)", filters)
	defer logger.Debugf("Exiting deleting EIPs (%s)", filters)
	ec2Client := getEC2Client(awsSession)

	describeAddressesInput := ec2.DescribeAddressesInput{}
	describeAddressesInput.Filters = createEC2Filters(filters)

	for {
		results, err := ec2Client.DescribeAddresses(&describeAddressesInput)
		if err != nil {
			logger.Debugf("error querying elastic IPs: %v", err)
			return false, nil
		}

		if len(results.Addresses) == 0 {
			// nothing left to delete
			break
		}

		for _, eip := range results.Addresses {
			// delete any network interface associated with the EIP (they are untagged)
			if eip.NetworkInterfaceId != nil {
				logger.Debugf("deleting EIP: %v", *eip.NetworkInterfaceId)
				err := deleteNetworkIface(eip.NetworkInterfaceId, ec2Client, logger)
				if err != nil {
					logger.Debugf("error deleting network iface: %v", err)
					continue
				}
			}

			_, err := ec2Client.ReleaseAddress(&ec2.ReleaseAddressInput{
				AllocationId: eip.AllocationId,
			})
			if err != nil {
				logger.Debugf("error deleting EIP: %v", err)
				continue
			} else {
				logger.WithField("ip", *eip.PublicIp).Info("Deleted Elastic IP")
			}

		}

		return false, nil
	}

	return true, nil
}

// deletePoliciesFromRole will attempt to delete any role policies from a provided role
func deletePoliciesFromRole(role *string, iamClient *iam.IAM) error {
	results, err := iamClient.ListRolePolicies(&iam.ListRolePoliciesInput{
		RoleName: role,
	})
	if err != nil {
		return err
	}

	for _, policy := range results.PolicyNames {
		_, err := iamClient.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
			RoleName:   role,
			PolicyName: policy,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteRolesFromInstanceProfile will attempt to delete any roles associated with a given instance profile
func deleteRolesFromInstanceProfile(ip *iam.InstanceProfile, iamClient *iam.IAM, logger log.FieldLogger) error {
	for _, role := range ip.Roles {
		logger.Debugf("deleting role %v from instance profile %v", *role.RoleName, *ip.InstanceProfileName)

		// empty the role
		logger.Debugf("deleting policies from role: %v", *role.RoleName)
		err := deletePoliciesFromRole(role.RoleName, iamClient)
		if err != nil {
			logger.Debugf("error deleting policies from role: %v", err)
			return err
		}

		logger.Infof("Deleted all policies from role: %v", *role.RoleName)

		// detach role from instance profile
		_, err = iamClient.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: ip.InstanceProfileName,
			RoleName:            role.RoleName,
		})
		if err != nil {
			logger.Debugf("error removing role from instance profile: %v", err)
			return err
		}

		logger.Infof("Removed role %v from instance profile %v", *role.RoleName, *ip.InstanceProfileName)

		// now delete the role
		// need to loop because this is the only time we'll have the name of the role
		// now that it has been detached from the instance profile
		for {
			_, err = iamClient.DeleteRole(&iam.DeleteRoleInput{
				RoleName: role.RoleName,
			})
			if err != nil {
				logger.Debugf("error deleting role %v from instance profile %v: %v", *role.RoleName, ip.InstanceProfileName, err)
			} else {
				logger.WithField("name", *role.RoleName).Info("Deleted role")
				break
			}

			time.Sleep(time.Second * secondsToSleep)
		}
	}

	return nil
}

// deleteInstanceProfile will attempt to delete the provided instance profile
func deleteInstanceProfile(instanceProfileID *string, iamClient *iam.IAM, logger log.FieldLogger) error {
	ipList, err := iamClient.ListInstanceProfiles(&iam.ListInstanceProfilesInput{})
	if err != nil {
		logger.Debugf("error listing instance profiles: %v", err)
		return err
	}

	var matchedIP *iam.InstanceProfile
	for _, ip := range ipList.InstanceProfiles {
		if *ip.InstanceProfileId == *instanceProfileID {
			matchedIP = ip
		}
	}

	if matchedIP == nil {
		// nothing found, so already deleted?
		return nil
	}

	// first delete any roles out of the instance profile
	err = deleteRolesFromInstanceProfile(matchedIP, iamClient, logger)
	if err != nil {
		return fmt.Errorf("error deleting roles from instance profile: %v", err)
	}

	logger.Debugf("deleting instance profile: %v", *matchedIP.InstanceProfileName)
	_, err = iamClient.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: matchedIP.InstanceProfileName,
	})
	if err != nil {
		logger.Debugf("error deleting instance profile: %v", err)
		return err
	} else if err == nil {
		logger.WithField("name", *matchedIP.InstanceProfileName).Info("Deleted instance profile")
	}

	return nil
}

// tryDeleteRoleProfileByName attempts to delete roles and profiles with given name ($CLUSTER_NAME-bootstrap|master|worker-role|profile)
func tryDeleteRoleProfileByName(roleName string, profileName string, session *session.Session, logger log.FieldLogger) error {
	logger.Debugf("deleting role: %s", roleName)
	describeRoleInput := iam.GetRoleInput{}
	describeRoleInput.RoleName = &roleName
	iamClient := iam.New(session)
	if _, err := iamClient.GetRole(&describeRoleInput); err != nil && err.(awserr.Error).Code() != iam.ErrCodeNoSuchEntityException {
		return err
	}

	// empty the role
	logger.Debugf("deleting policies from role: %s", roleName)
	if err := deletePoliciesFromRole(&roleName, iamClient); err != nil && err.(awserr.Error).Code() != iam.ErrCodeNoSuchEntityException {
		logger.Debugf("error deleting policies from role: %v", err)
		return err
	}
	describeProfileInput := iam.GetInstanceProfileInput{}
	describeProfileInput.InstanceProfileName = &profileName
	if _, err := iamClient.GetInstanceProfile(&describeProfileInput); err != nil && err.(awserr.Error).Code() != iam.ErrCodeNoSuchEntityException {
		return err
	}

	// detach role from profile
	logger.Debugf("detaching role from profile: %s", profileName)
	_, err := iamClient.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: &profileName,
		RoleName:            &roleName,
	})
	if err != nil && err.(awserr.Error).Code() != iam.ErrCodeNoSuchEntityException {
		logger.Debugf("error removing role from instance profile: %v", err)
		return err
	}
	if err == nil {
		logger.Infof("Removed role %v from instance profile %v", roleName, profileName)
	}
	// delete profile
	logger.Debugf("deleting instance profile: %v", profileName)
	_, err = iamClient.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: &profileName,
	})
	if err != nil && err.(awserr.Error).Code() != iam.ErrCodeNoSuchEntityException {
		logger.Debugf("error deleting instance profile %s: %v", profileName, err)
		return err
	}
	if err == nil {
		logger.Infof("deleted profile %s", profileName)
	}
	// now we can delete role
	logger.Debugf("deleted policies from role %s", roleName)
	deleteRoleInput := iam.DeleteRoleInput{}
	deleteRoleInput.RoleName = &roleName
	if _, err := iamClient.DeleteRole(&deleteRoleInput); err != nil && err.(awserr.Error).Code() != iam.ErrCodeNoSuchEntityException {
		logger.Debugf("error deleting role %s: %v", roleName, err)
		return err
	}
	if err == nil {
		logger.Infof("deleted role %s", roleName)
	}
	return nil
}

// deleteIAMresources will delete any IAM resources created by the installer that are not associated with a running instance
// Currently openshift/installer creates 3 roles per cluster, 1 for master|worker|bootstrap and identified by the
// cluster name used to install the cluster.
func deleteIAMresources(session *session.Session, filter AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting IAM resources (%s)", filter)
	defer logger.Debugf("Exiting deleting IAM resources (%s)", filter)
	installerType := []string{"master", "worker", "bootstrap"}
	for _, t := range installerType {
		// Naming of IAM resources expected from https://github.com/openshift/installer as follows:
		// $CLUSTER_NAME-master-role     $CLUSTER_NAME-worker-role    $CLUSTER_NAME-bootstrap-role
		// $CLUSTER_NAME-master-profile  $CLUSTER_NAME-worker-profile $CLUSTER_NAME-bootstrap-profile
		roleName := fmt.Sprintf("%s-%s-role", clusterName, t)
		instanceProfileName := fmt.Sprintf("%s-%s-profile", clusterName, t)
		if err := tryDeleteRoleProfileByName(roleName, instanceProfileName, session, logger); err != nil {
			logger.Debugf("error deleting instance profile %s: %v", instanceProfileName, err)
			return false, nil
		}
	}
	return true, nil
}

// deleteInstances will find any running/pending instances that match the given filter and terminate them
// and any instance profiles attached to the instance(s)
func deleteInstances(session *session.Session, filter AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting instances (%s)", filter)
	defer logger.Debugf("Exiting deleting instances (%s)", filter)

	ec2Client := getEC2Client(session)
	iamClient := iam.New(session)

	describeInstancesInput := ec2.DescribeInstancesInput{}
	describeInstancesInput.Filters = createEC2Filters(filter)

	// only fetch instances in 'running|pending' state since 'terminated' ones take a while to really get cleaned up
	describeInstancesInput.Filters = append(describeInstancesInput.Filters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("running"), aws.String("pending")},
	})

	instancesFound := false
	err := ec2Client.DescribeInstancesPages(&describeInstancesInput, func(results *ec2.DescribeInstancesOutput, lastPage bool) bool {
		instancesFound = instancesFound || len(results.Reservations) > 0
		for _, reservation := range results.Reservations {
			for _, instance := range reservation.Instances {
				// first delete any instance profiles (they are not tagged)
				if instance.IamInstanceProfile != nil {
					err := deleteInstanceProfile(instance.IamInstanceProfile.Id, iamClient, logger)
					if err != nil {
						logger.Debugf("error deleting instance profile: %v", err)
						continue
					}
				}

				// now delete the instance
				logger.Debugf("deleting instance: %v", *instance.InstanceId)
				_, err := ec2Client.TerminateInstances(&ec2.TerminateInstancesInput{
					InstanceIds: []*string{instance.InstanceId},
				})
				if err != nil {
					logger.Debugf("error deleting instance: %v", err)
					continue
				} else {
					logger.WithField("id", *instance.InstanceId).Info("Deleted instance")
				}
			}
		}

		return lastPage
	})
	if err != nil {
		logger.Debugf("error describing instances: %v", err)
		return false, nil
	}

	return !instancesFound, nil
}

// deleteSecurityGroupRules will attempt to delete all the rules defined in the given security group
// since some security groups have self-referencing rules that complicate being able to delete the security group
func deleteSecurityGroupRules(sg *ec2.SecurityGroup, ec2Client *ec2.EC2, logger log.FieldLogger) error {

	if len(sg.IpPermissions) > 0 {
		_, err := ec2Client.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
			GroupId:       sg.GroupId,
			IpPermissions: sg.IpPermissions,
		})
		if err != nil {
			logger.Debugf("error removing ingress permissions: %v", err)
		}
	}

	if len(sg.IpPermissionsEgress) > 0 {
		_, err := ec2Client.RevokeSecurityGroupEgress(&ec2.RevokeSecurityGroupEgressInput{
			GroupId:       sg.GroupId,
			IpPermissions: sg.IpPermissionsEgress,
		})
		if err != nil {
			logger.Debugf("error removing egress permissions: %v", err)
		}
	}

	return nil
}

// deleteSecurityGroups will attempt to delete all security groups matching the given filter
func deleteSecurityGroups(session *session.Session, filter AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting security groups (%s)", filter)
	defer logger.Debugf("Exiting deleting security groups (%s)", filter)

	ec2Client := getEC2Client(session)
	describeSecurityGroupsInput := ec2.DescribeSecurityGroupsInput{}
	describeSecurityGroupsInput.Filters = createEC2Filters(filter)

	for {
		results, err := ec2Client.DescribeSecurityGroups(&describeSecurityGroupsInput)
		if err != nil {
			logger.Debugf("error listing security groups %v", err)
			return false, nil
		}

		if len(results.SecurityGroups) == 0 {
			break
		}

		for _, sg := range results.SecurityGroups {
			// first delete rules (can get circular dependencies otherwise)
			deleteSecurityGroupRules(sg, ec2Client, logger)
			_, err := ec2Client.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
				GroupId: sg.GroupId,
			})
			if err != nil {
				logger.Debugf("error deleting security group: %v", err)
				continue
			} else {
				logger.WithField("id", *sg.GroupId).Info("Deleted security group")
			}
		}

		return false, nil
	}

	return true, nil
}

// detachInternetGateways will attempt to detach an internet gateway from the associated VPC(s)
func detachInternetGateways(gw *ec2.InternetGateway, ec2Client *ec2.EC2, logger log.FieldLogger) error {
	for _, vpc := range gw.Attachments {
		logger.Debugf("detaching Internet GW %v from VPC %v", *gw.InternetGatewayId, *vpc.VpcId)
		_, err := ec2Client.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
			InternetGatewayId: gw.InternetGatewayId,
			VpcId:             vpc.VpcId,
		})

		if err != nil {
			return fmt.Errorf("error detaching internet gateway: %v", err)
		} else if err == nil {
			logger.Infof("Detached Internet GW %v from VPC %v", *gw.InternetGatewayId, *vpc.VpcId)
		}
	}

	return nil
}

// deleteInternetGateways will attemp to delete any Internet Gateways matching the given filter
func deleteInternetGateways(session *session.Session, filter AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting internet gateways (%s)", filter)
	defer logger.Debugf("Exiting deleting internet gateways (%s)", filter)

	ec2Client := getEC2Client(session)

	describeInternetGatewaysInput := ec2.DescribeInternetGatewaysInput{}
	describeInternetGatewaysInput.Filters = createEC2Filters(filter)

	for {
		results, err := ec2Client.DescribeInternetGateways(&describeInternetGatewaysInput)
		if err != nil {
			logger.Debugf("error listing internet gateways: %v", err)
			return false, nil
		}

		if len(results.InternetGateways) == 0 {
			break
		}

		for _, gw := range results.InternetGateways {
			logger.Debugf("deleting internet gateway: %v", *gw.InternetGatewayId)

			err := detachInternetGateways(gw, ec2Client, logger)
			if err != nil {
				logger.Debugf("error detaching igw: %v", err)
				continue
			}

			_, err = ec2Client.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
				InternetGatewayId: gw.InternetGatewayId,
			})
			if err != nil {
				logger.Debugf("error deleting internet gateway: %v", err)
			} else {
				logger.WithField("id", *gw.InternetGatewayId).Info("Deleted internet gateway")
			}
		}

		return false, nil
	}

	return true, nil
}

// disassociateRouteTable will attempt to disassociate all except the 'Main' associations defined
// for the given Route Table
func disassociateRouteTable(rt *ec2.RouteTable, ec2Client *ec2.EC2, logger log.FieldLogger) error {
	for _, association := range rt.Associations {
		if *association.Main {
			// can't remove the 'Main' association
			continue
		}
		logger.Debugf("disassociating route table association %v", *association.RouteTableAssociationId)
		_, err := ec2Client.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
			AssociationId: association.RouteTableAssociationId,
		})
		if err != nil {
			logger.Debugf("error disassociating from route table: %v", err)
			return err
		} else if err == nil {
			logger.WithField("id", *association.RouteTableAssociationId).Info("Disassociated route table association")
		}
	}

	return nil
}

// deleteSubnets will attempt to delete all Subnets matching the given filter
func deleteSubnets(session *session.Session, filter AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting subnets (%s)", filter)
	defer logger.Debugf("Exiting deleting subnets (%s)", filter)

	ec2Client := getEC2Client(session)

	describeSubnetsInput := ec2.DescribeSubnetsInput{}
	describeSubnetsInput.Filters = createEC2Filters(filter)

	for {
		results, err := ec2Client.DescribeSubnets(&describeSubnetsInput)
		if err != nil {
			logger.Debugf("error listing subnets: %v", err)
			return false, nil
		}

		if len(results.Subnets) == 0 {
			break
		}

		for _, subnet := range results.Subnets {
			_, err := ec2Client.DeleteSubnet(&ec2.DeleteSubnetInput{
				SubnetId: subnet.SubnetId,
			})
			if err != nil {
				logger.Debugf("error deleting subnet: %v", err)
			} else {
				logger.WithField("id", *subnet.SubnetId).Info("Deleted subnet")
			}
		}

		return false, nil
	}

	return true, nil
}

// bucketsToAWSObjects will convert a list of S3 Buckets to awsObjectsWithTags (for easier filtering)
func bucketsToAWSObjects(buckets []*s3.Bucket, s3Client *s3.S3, logger log.FieldLogger) ([]awsObjectWithTags, error) {
	bucketObjects := []awsObjectWithTags{}

	for _, bucket := range buckets {
		tags, err := s3Client.GetBucketTagging(&s3.GetBucketTaggingInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			logger.Debugf("error getting tags for bucket %s: %v, skipping...", *bucket.Name, err)
			continue
		}

		tagsAsMap, err := tagsToMap(tags.TagSet)
		if err != nil {
			return bucketObjects, err
		}
		bucketObjects = append(bucketObjects, awsObjectWithTags{
			Name: *bucket.Name,
			Tags: tagsAsMap,
		})
	}

	return bucketObjects, nil
}

// filterObjects will do client-side filtering given an appropriately filled out list of awsObjectWithTags
func filterObjects(awsObjects []awsObjectWithTags, filters AWSFilter) []awsObjectWithTags {
	objectsWithTags := []awsObjectWithTags{}
	filteredObjects := []awsObjectWithTags{}

	// first find the objects that have all the desired tags
	for _, object := range awsObjects {
		allTagsFound := true
		for key := range filters {
			if _, ok := object.Tags[key]; !ok {
				// doesn't have one of the tags we're looking for so skip it
				allTagsFound = false
				break
			}
		}
		if allTagsFound {
			objectsWithTags = append(objectsWithTags, object)
		}
	}

	// now check that the values match
	for _, object := range objectsWithTags {
		valuesMatch := true
		for key, val := range filters {
			if object.Tags[key] != val {
				valuesMatch = false
				break
			}
		}
		if valuesMatch {
			filteredObjects = append(filteredObjects, object)
		}
	}
	return filteredObjects
}

// deleteS3Buckets will attempt to delete (and empty) any S3 bucket matching the provided filter
func deleteS3Buckets(session *session.Session, filter AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting S3 buckets (%s)", filter)
	defer logger.Debugf("Exiting deleting buckets (%s)", filter)

	s3Client := s3.New(session)

	listBucketsInput := s3.ListBucketsInput{}

	for {
		results, err := s3Client.ListBuckets(&listBucketsInput)
		if err != nil {
			logger.Debugf("error listing s3 buckets: %v", err)
			return false, nil
		}

		awsObjects, err := bucketsToAWSObjects(results.Buckets, s3Client, logger)
		if err != nil {
			logger.Debugf("error converting s3 buckets to native AWS objects: %v", err)
			return false, nil
		}

		filteredObjects := filterObjects(awsObjects, filter)
		logger.Debugf("from %d total s3 buckets, %d match filters", len(awsObjects), len(filteredObjects))
		if len(filteredObjects) == 0 {
			break
		}

		for _, bucket := range filteredObjects {
			logger.Debugf("deleting bucket: %v", bucket.Name)

			// first empty the bucket
			iter := s3manager.NewDeleteListIterator(s3Client, &s3.ListObjectsInput{
				Bucket: aws.String(bucket.Name),
			})
			err := s3manager.NewBatchDeleteWithClient(s3Client).Delete(aws.BackgroundContext(), iter)
			if err != nil {
				logger.Debugf("error emptying bucket %v: %v", bucket.Name, err)
				continue
			} else {
				logger.WithField("name", bucket.Name).Info("Emptied bucket")
			}

			// now delete the bucket
			_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
				Bucket: aws.String(bucket.Name),
			})
			if err != nil {
				logger.Debugf("error deleting bucket %v: %v", bucket.Name, err)
				continue
			} else {
				logger.WithField("name", bucket.Name).Info("Deleted bucket")
			}
		}

		return false, nil
	}

	return true, nil
}

// r53ZonesToAWSObjects will create a list of awsObjectsWithTags for the provided list of route53.HostedZone s
func r53ZonesToAWSObjects(zones []*route53.HostedZone, r53Client *route53.Route53) ([]awsObjectWithTags, error) {
	zonesAsAWSObjects := []awsObjectWithTags{}

	for _, zone := range zones {
		result, err := r53Client.ListTagsForResource(&route53.ListTagsForResourceInput{
			ResourceType: aws.String("hostedzone"),
			ResourceId:   zone.Id,
		})
		if err != nil {
			return zonesAsAWSObjects, err
		}

		tagsToMap, err := tagsToMap(result.ResourceTagSet.Tags)
		if err != nil {
			return zonesAsAWSObjects, err
		}

		zonesAsAWSObjects = append(zonesAsAWSObjects, awsObjectWithTags{
			Name: *zone.Id,
			Tags: tagsToMap,
		})

	}

	return zonesAsAWSObjects, nil
}

// deleteEntriesFromSharedR53Zone will find route53 entries for the shared (ie non-terraform-managed) route53 zone
// and remove them.
// Provide the terraform-created private zone, and the manually created public/shared zone, and it will find any
// entries in the public/shared zone that match entries in the private zone, and delete them
func deleteEntriesFromSharedR53Zone(zoneID string, sharedZoneID string, r53Client *route53.Route53, logger log.FieldLogger) error {

	zoneEntries, err := r53Client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
	})
	if err != nil {
		return err
	}

	sharedZoneEntries, err := r53Client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(sharedZoneID),
	})
	if err != nil {
		return err
	}

	for _, entry := range zoneEntries.ResourceRecordSets {
		// only interested in deleting 'A' records
		if *entry.Type != "A" {
			continue
		}
		for _, sharedEntry := range sharedZoneEntries.ResourceRecordSets {
			if *sharedEntry.Name == *entry.Name && *sharedEntry.Type == *entry.Type {
				_, err := r53Client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
					HostedZoneId: aws.String(sharedZoneID),
					ChangeBatch: &route53.ChangeBatch{
						Changes: []*route53.Change{
							{
								Action: aws.String("DELETE"),
								ResourceRecordSet: &route53.ResourceRecordSet{
									Name:        sharedEntry.Name,
									Type:        sharedEntry.Type,
									AliasTarget: sharedEntry.AliasTarget,
								},
							},
						},
					},
				})
				if err != nil {
					return err
				}

				logger.Infof("Deleted record %v from r53 zone %v", *sharedEntry.Name, sharedZoneID)
			}
		}
	}

	return nil
}

// getSharedHostedZone will find the zoneID of the non-terraform-managed public route53 zone given the
// terraform-managed private zoneID
func getSharedHostedZone(zoneID string, allZones []*route53.HostedZone) (string, error) {
	// given the ID, get the name of the zone
	zoneName := ""
	for _, zone := range allZones {
		if *zone.Id == zoneID {
			zoneName = *zone.Name
			break
		}
	}

	// now find the shared zone that matches by name
	for _, zone := range allZones {
		// skip the actual terraform-managed zone (we're looking for the shared zone)
		if *zone.Id == zoneID {
			continue
		}

		if *zone.Name == zoneName {
			return *zone.Id, nil
		}
	}

	// else we didn't find it
	return "", fmt.Errorf("could not find shared zone with name: %v", zoneName)
}

// emptyAndDeleteRoute53Zone will delete all the entries in the given route53 zone and delete the zone itself
func emptyAndDeleteRoute53Zone(zoneID string, r53Client *route53.Route53, logger log.FieldLogger) error {

	// first need to delete all non SOA and NS records
	results, err := r53Client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
	})
	if err != nil {
		return err
	}

	for _, entry := range results.ResourceRecordSets {
		if *entry.Type == "SOA" || *entry.Type == "NS" {
			// can't delete SOA and NS types
			continue
		}
		_, err := r53Client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(zoneID),
			ChangeBatch: &route53.ChangeBatch{
				Changes: []*route53.Change{
					{
						Action: aws.String("DELETE"),
						ResourceRecordSet: &route53.ResourceRecordSet{
							Name:            entry.Name,
							Type:            entry.Type,
							TTL:             entry.TTL,
							ResourceRecords: entry.ResourceRecords,
							AliasTarget:     entry.AliasTarget,
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		logger.Infof("Deleted record %v from r53 zone %v", *entry.Name, zoneID)
	}

	// now delete zone
	_, err = r53Client.DeleteHostedZone(&route53.DeleteHostedZoneInput{
		Id: aws.String(zoneID),
	})
	if err != nil {
		return err
	}

	logger.WithField("id", zoneID).Info("Deleted route53 zone")

	return nil
}

// deleteRoute53 will attempt to delete any route53 zone matching the given filter.
// it will also attempt to delete any entries in the shared/public route53 zone
func deleteRoute53(session *session.Session, filters AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {
	logger.Debugf("Deleting Route53 zones (%s)", filters)
	defer logger.Debugf("Exiting deleting Route53 zones (%s)", filters)

	r53Client := route53.New(session)

	listHostedZonesInput := route53.ListHostedZonesInput{}

	for {
		allZones, err := r53Client.ListHostedZones(&listHostedZonesInput)
		if err != nil {
			logger.Debugf("error listing route53 zones: %v", err)
			return false, nil
		}

		awsZones, err := r53ZonesToAWSObjects(allZones.HostedZones, r53Client)
		if err != nil {
			logger.Debugf("error converting r53Zones to native AWS objects: %v", err)
			return false, nil
		}

		filteredZones := filterObjects(awsZones, filters)
		logger.Debugf("from %d total r53 zones, %d match filters", len(awsZones), len(filteredZones))
		if len(filteredZones) == 0 {
			break
		}

		for _, zone := range filteredZones {
			// first find the shared hostedzone (will have same name as the tagged zone)
			sharedZoneID, err := getSharedHostedZone(zone.Name, allZones.HostedZones)
			if err != nil {
				logger.Debugf("%v", err)
				return false, nil
			}

			// first need to delete any 'A' entries from the shared-non-private Route53 zone
			// (eg. newcluster.subdomain.domain.com newcluster-api.subdomain.domain.com and
			// *.newcluster.subdomain.domain.com)
			err = deleteEntriesFromSharedR53Zone(zone.Name, sharedZoneID, r53Client, logger)
			if err != nil {
				logger.Debugf("error deleting entries from shared r53 zone: %v", err)
				return false, nil
			}

			// finally can delete the tagged hosted zone
			err = emptyAndDeleteRoute53Zone(zone.Name, r53Client, logger)
			if err != nil {
				logger.Debugf("error deleting zone %v: %v", zone.Name, err)
				return false, nil
			}
		}
		return false, nil
	}
	// all done deleting r53 entries/zones
	return true, nil
}

// deletePVs will find PVs based on provided filters and delete them
func deletePVs(session *session.Session, filters AWSFilter, clusterName string, logger log.FieldLogger) (bool, error) {

	logger.Debugf("Deleting PVs (%s)", filters)
	defer logger.Debugf("Exiting deleting PVs (%s)", filters)

	ec2Client := getEC2Client(session)
	describeVolumesInput := ec2.DescribeVolumesInput{}
	describeVolumesInput.Filters = createEC2Filters(filters)

	results, err := ec2Client.DescribeVolumes(&describeVolumesInput)
	if err != nil {
		logger.Debugf("error listing volumes: %v", err)
		return false, nil
	}

	if len(results.Volumes) == 0 {
		// nothing to delete, we must be done
		return true, nil
	}

	for _, vol := range results.Volumes {
		logger.Debugf("deleting volume: %v", *vol.VolumeId)
		_, err := ec2Client.DeleteVolume(&ec2.DeleteVolumeInput{
			VolumeId: vol.VolumeId,
		})
		if err != nil {
			logger.Debugf("error deleting volume: %v", err)
		} else {
			logger.WithField("id", *vol.VolumeId).Info("Deleted Volume")
		}
	}

	return false, nil
}
