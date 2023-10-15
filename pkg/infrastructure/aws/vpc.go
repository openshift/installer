package aws

import (
	"fmt"
	"math"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/sirupsen/logrus"
)

func createVPCResources(logger *logrus.Logger, session *session.Session, vpcInput *CreateInfraOptions) error {
	ec2Client := ec2.New(session)

	vpcInput.additionalEC2Tags = ec2CreateTags(vpcInput.AdditionalTags)

	var err error
	zoneToPrivateSubnetIDMap := make(map[string]string)
	vpcID := vpcInput.vpcID
	if vpcID == "" {
		vpcID, err = vpcInput.createVPC(logger, ec2Client)
		if err != nil {
			return err
		}
		vpcInput.vpcID = vpcID

		if err = vpcInput.CreateDHCPOptions(logger, ec2Client, vpcID); err != nil {
			return err
		}

		igwID, err := vpcInput.CreateInternetGateway(logger, ec2Client, vpcID)
		if err != nil {
			return err
		}

		bootstrapSG, err := vpcInput.CreateBootstrapSecurityGroup(logger, ec2Client, vpcID)
		if err != nil {
			return err
		}
		var machineV4Cidrs []string
		if vpcInput.public {
			machineV4Cidrs = []string{"0.0.0.0/0"}
		} else {
			machineV4Cidrs = vpcInput.cidrV4Blocks
		}
		bootstrapIngressPermissions := DefaultBootstrapSGIngressRules(vpcInput.bootstrapSecurityGroupID, machineV4Cidrs)
		if err := vpcInput.AttachSecurityGroupIngressRules(logger, ec2Client, bootstrapSG, bootstrapIngressPermissions); err != nil {
			return err
		}

		// Per zone resources
		_, network, err := net.ParseCIDR(vpcInput.cidrV4Blocks[0])
		if err != nil {
			return err
		}
		privateNetwork, err := cidr.Subnet(network, 1, 0)
		if err != nil {
			return err
		}
		publicNetwork, err := cidr.Subnet(network, 1, 1)
		if err != nil {
			return err
		}

		// If a single-zone deployment, the available CIDR block will be split
		// into two to allow user expansion
		if len(vpcInput.Zones) == 1 {
			privateNetwork, err = cidr.Subnet(privateNetwork, 1, 0)
			if err != nil {
				return err
			}
			publicNetwork, err = cidr.Subnet(publicNetwork, 1, 0)
			if err != nil {
				return err
			}
		}

		var publicSubnetIDs []string
		var privateSubnetIDs []string
		var endpointRouteTableIds []*string
		newBits := int(math.Ceil(math.Log2(float64(len(vpcInput.Zones)))))
		for i, zone := range vpcInput.Zones {
			privateCIDR, err := cidr.Subnet(privateNetwork, newBits, i)
			if err != nil {
				return err
			}
			privateSubnetID, err := vpcInput.CreatePrivateSubnet(logger, ec2Client, vpcID, zone, privateCIDR.String())
			if err != nil {
				return err
			}
			zoneToPrivateSubnetIDMap[zone] = privateSubnetID

			publicCIDR, err := cidr.Subnet(publicNetwork, newBits, i)
			if err != nil {
				return err
			}
			publicSubnetID, err := vpcInput.CreatePublicSubnet(logger, ec2Client, vpcID, zone, publicCIDR.String())
			if err != nil {
				return err
			}

			var natGatewayID string
			publicSubnetIDs = append(publicSubnetIDs, publicSubnetID)
			privateSubnetIDs = append(privateSubnetIDs, privateSubnetID)

			if !vpcInput.EnableProxy {
				natGatewayID, err = vpcInput.CreateNATGateway(logger, ec2Client, publicSubnetID, zone)
				if err != nil {
					return err
				}
			}
			privateRouteTable, err := vpcInput.CreatePrivateRouteTable(logger, ec2Client, vpcID, natGatewayID, privateSubnetID, zone)
			if err != nil {
				return err
			}
			endpointRouteTableIds = append(endpointRouteTableIds, aws.String(privateRouteTable))
		}
		vpcInput.publicSubnetIDs = publicSubnetIDs
		vpcInput.privateSubnetIDs = privateSubnetIDs

		publicRouteTable, err := vpcInput.CreatePublicRouteTable(logger, ec2Client, vpcID, igwID, publicSubnetIDs)
		if err != nil {
			return err
		}

		endpointRouteTableIds = append(endpointRouteTableIds, aws.String(publicRouteTable))
		err = vpcInput.CreateVPCS3Endpoint(logger, ec2Client, vpcID, endpointRouteTableIds)
		if err != nil {
			return err
		}
	} else {
		logger.WithField("id", vpcID).Debugln("Using user-supplied VPC")

		subnets, err := ec2GetSubnets(ec2Client, vpcInput.privateSubnetIDs)
		if err != nil {
			return err
		}
		for _, subnet := range subnets {
			zoneToPrivateSubnetIDMap[aws.StringValue(subnet.AvailabilityZone)] = aws.StringValue(subnet.SubnetId)
		}
	}
	vpcInput.zoneToSubnetIDMap = zoneToPrivateSubnetIDMap

	masterSG, err := vpcInput.CreateMasterSecurityGroup(logger, ec2Client, vpcID)
	if err != nil {
		return err
	}
	//vpcInput.masterSecurityGroupID = aws.StringValue(masterSG.GroupId)

	workerSG, err := vpcInput.CreateWorkerSecurityGroup(logger, ec2Client, vpcID)
	if err != nil {
		return err
	}
	//vpcInput.workerSecurityGroupID = aws.StringValue(workerSG.GroupId)

	masterIngressPermissions := DefaultMasterSGIngressRules(vpcInput.masterSecurityGroupID, vpcInput.workerSecurityGroupID, vpcInput.cidrV4Blocks)
	// masterIngressPermissions := DefaultAllowAllSGIngressRules(vpcInput.masterSecurityGroupID, []string{})
	if err := vpcInput.AttachSecurityGroupIngressRules(logger, ec2Client, masterSG, masterIngressPermissions); err != nil {
		return err
	}

	workerIngressPermissions := DefaultWorkerSGIngressRules(vpcInput.workerSecurityGroupID, vpcInput.masterSecurityGroupID, vpcInput.cidrV4Blocks)
	// workerIngressPermissions := DefaultAllowAllSGIngressRules(vpcInput.workerSecurityGroupID, []string{})
	if err := vpcInput.AttachSecurityGroupIngressRules(logger, ec2Client, workerSG, workerIngressPermissions); err != nil {
		return err
	}

	if err := vpcInput.CreateLoadBalancers(logger, session, vpcID, vpcInput.privateSubnetIDs, vpcInput.publicSubnetIDs, vpcInput.public); err != nil {
		return err
	}

	return nil
}

func (o *CreateInfraOptions) createVPC(l *logrus.Logger, client ec2iface.EC2API) (string, error) {
	vpcName := fmt.Sprintf("%s-vpc", o.InfraID)

	vpcID, err := o.existingVPC(client, vpcName)
	if err != nil {
		return "", err
	}
	if vpcID == "" {
		vpcID, err = ec2CreateVPC(client, vpcName, o.cidrV4Blocks[0], o.additionalEC2Tags)
		if err != nil {
			return "", err
		}
		l.WithField("id", vpcID).Infoln("Created VPC")
	} else {
		l.WithField("id", vpcID).Infoln("Found existing VPC")
	}

	logger := l.WithField("id", vpcID)
	if err := ec2VPCEnableDNSSupport(client, vpcID); err != nil {
		return "", err
	}
	logger.Info("Enabled DNS support on VPC")
	if err := ec2VPCEnableDNSHostnames(client, vpcID); err != nil {
		return "", err
	}
	logger.Info("Enabled DNS hostnames on VPC")

	return vpcID, nil
}

func (o *CreateInfraOptions) existingVPC(client ec2iface.EC2API, vpcName string) (string, error) {
	return ec2GetVPC(client, o.ec2Filters(vpcName))
}

func (o *CreateInfraOptions) ec2Filters(name string) []*ec2.Filter {
	filters := []*ec2.Filter{
		ec2CreateFilter(fmt.Sprintf("tag:%s", clusterTag(o.InfraID)), clusterTagValue),
	}
	if len(name) > 0 {
		filters = append(filters, ec2CreateFilter("tag:Name", name))
	}
	return filters
}

func clusterTag(infraID string) string {
	return fmt.Sprintf("kubernetes.io/cluster/%s", infraID)
}

func (o *CreateInfraOptions) CreateDHCPOptions(l *logrus.Logger, client ec2iface.EC2API, vpcID string) error {
	optID, err := o.existingDHCPOptions(client)
	if err != nil {
		return err
	}
	if optID == "" {
		domainName := "ec2.internal"
		if o.Region != "us-east-1" {
			domainName = fmt.Sprintf("%s.compute.internal", o.Region)
		}
		optID, err = ec2CreateDHCPOptions(client, "", domainName, o.additionalEC2Tags)
		if err != nil {
			return err
		}
		l.WithField("id", optID).Info("Created DHCP options")
	} else {
		l.WithField("id", optID).Info("Found existing DHCP options")
	}

	err = ec2AssociateDHCPOptionsToVPC(client, optID, vpcID)
	if err != nil {
		return err
	}
	l.WithField("vpc", vpcID).WithField("dhcp option", optID).Infoln("Associated DHCP options with VPC")

	return nil
}

func (o *CreateInfraOptions) existingDHCPOptions(client ec2iface.EC2API) (string, error) {
	return ec2GetDHCPOptions(client, o.ec2Filters(""))
}

func (o *CreateInfraOptions) CreateInternetGateway(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (string, error) {
	gatewayName := fmt.Sprintf("%s-igw", o.InfraID)
	igw, err := o.existingInternetGateway(client, gatewayName)
	if err != nil {
		return "", err
	}
	if igw == nil {
		igw, err = ec2CreateInternetGateway(client, gatewayName, o.additionalEC2Tags)
		if err != nil {
			return "", err
		}
		l.WithField("id", aws.StringValue(igw.InternetGatewayId)).Infoln("Created internet gateway")
	} else {
		l.WithField("id", aws.StringValue(igw.InternetGatewayId)).Infoln("Found existing internet gateway")
	}
	igwID := aws.StringValue(igw.InternetGatewayId)
	attached := false
	for _, attachment := range igw.Attachments {
		if aws.StringValue(attachment.VpcId) == vpcID {
			attached = true
			break
		}
	}
	if !attached {
		err := ec2AttachInternetGatewayToVPC(client, igwID, vpcID)
		if err != nil {
			return "", err
		}
		l.WithField("internet gateway", igwID).WithField("vpc", vpcID).Infoln("Attached internet gateway to VPC")
	}
	return igwID, nil
}

func (o *CreateInfraOptions) existingInternetGateway(client ec2iface.EC2API, name string) (*ec2.InternetGateway, error) {
	return ec2GetInternetGateway(client, o.ec2Filters(name))
}

func (o *CreateInfraOptions) createSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string, groupName string) (*ec2.SecurityGroup, error) {
	securityGroup, err := o.existingSecurityGroup(client, groupName)
	if err != nil {
		return nil, err
	}
	logger := l.WithField("name", groupName)
	if securityGroup != nil {
		logger.WithField("id", aws.StringValue(securityGroup.GroupId)).Infoln("Found existing security group")
		return securityGroup, nil
	}
	securityGroup, err = ec2CreateSecurityGroup(client, groupName, vpcID, o.additionalEC2Tags)
	if err != nil {
		return nil, err
	}
	logger.WithField("id", aws.StringValue(securityGroup.GroupId)).Infoln("Created security group")
	return securityGroup, nil
}

func (o *CreateInfraOptions) CreateBootstrapSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (*ec2.SecurityGroup, error) {
	groupName := fmt.Sprintf("%s-bootstrap-sg", o.InfraID)
	securityGroup, err := o.createSecurityGroup(l, client, vpcID, groupName)
	if err != nil {
		return nil, err
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	//sgUserID := aws.StringValue(securityGroup.OwnerId)
	egressPermissions := DefaultSGEgressRules(securityGroupID)
	if err := o.AttachSecurityGroupEgressRules(l, client, securityGroup, egressPermissions); err != nil {
		return nil, err
	}
	o.bootstrapSecurityGroupID = securityGroupID
	return securityGroup, nil
}

func (o *CreateInfraOptions) CreateMasterSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (*ec2.SecurityGroup, error) {
	groupName := fmt.Sprintf("%s-master-sg", o.InfraID)
	securityGroup, err := o.createSecurityGroup(l, client, vpcID, groupName)
	if err != nil {
		return nil, err
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	//sgUserID := aws.StringValue(securityGroup.OwnerId)
	egressPermissions := DefaultSGEgressRules(securityGroupID)
	if err := o.AttachSecurityGroupEgressRules(l, client, securityGroup, egressPermissions); err != nil {
		return nil, err
	}
	o.masterSecurityGroupID = securityGroupID
	return securityGroup, nil
}

func (o *CreateInfraOptions) CreateWorkerSecurityGroup(l *logrus.Logger, client ec2iface.EC2API, vpcID string) (*ec2.SecurityGroup, error) {
	groupName := fmt.Sprintf("%s-worker-sg", o.InfraID)
	securityGroup, err := o.createSecurityGroup(l, client, vpcID, groupName)
	if err != nil {
		return nil, err
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	//sgUserID := aws.StringValue(securityGroup.OwnerId)
	egressPermissions := DefaultSGEgressRules(securityGroupID)
	if err := o.AttachSecurityGroupEgressRules(l, client, securityGroup, egressPermissions); err != nil {
		return nil, err
	}

	o.workerSecurityGroupID = securityGroupID
	return securityGroup, nil
}

func (o *CreateInfraOptions) AttachSecurityGroupEgressRules(l *logrus.Logger, client ec2iface.EC2API, securityGroup *ec2.SecurityGroup, egressPermissions []*ec2.IpPermission) error {
	var egressToAuthorize []*ec2.IpPermission
	for _, permission := range egressPermissions {
		permission := permission
		if !includesPermission(securityGroup.IpPermissionsEgress, permission) {
			egressToAuthorize = append(egressToAuthorize, permission)
		}
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	logger := l.WithField("id", securityGroupID)
	if len(egressToAuthorize) > 0 {
		logger.Infoln("Authorizing egress rules on security group")
		err := ec2AuthorizeEgressRules(client, securityGroupID, egressToAuthorize)
		if err != nil {
			return err
		}
		logger.Infoln("Authorized egress rules on security group")
	}
	return nil
}

func (o *CreateInfraOptions) AttachSecurityGroupIngressRules(l *logrus.Logger, client ec2iface.EC2API, securityGroup *ec2.SecurityGroup, ingressPermissions []*ec2.IpPermission) error {
	var ingressToAuthorize []*ec2.IpPermission
	for _, permission := range ingressPermissions {
		permission := permission
		if !includesPermission(securityGroup.IpPermissions, permission) {
			ingressToAuthorize = append(ingressToAuthorize, permission)
		}
	}
	securityGroupID := aws.StringValue(securityGroup.GroupId)
	logger := l.WithField("id", securityGroupID)
	if len(ingressToAuthorize) > 0 {
		logger.Infoln("Authorizing ingress rules on security group")
		err := ec2AuthorizeIngressRules(client, securityGroupID, ingressToAuthorize)
		if err != nil {
			return err
		}
		logger.Infoln("Authorized ingress rules on security group")
	}
	return nil
}

func (o *CreateInfraOptions) existingSecurityGroup(client ec2iface.EC2API, name string) (*ec2.SecurityGroup, error) {
	return ec2GetSecurityGroup(client, o.ec2Filters(name))
}

func includesPermission(list []*ec2.IpPermission, permission *ec2.IpPermission) bool {
	for _, p := range list {
		if samePermission(p, permission) {
			return true
		}
	}
	return false
}

func samePermission(a, b *ec2.IpPermission) bool {
	if a == nil || b == nil {
		return false
	}
	if a.String() == b.String() {
		return true
	}
	return false
}

func DefaultSGEgressRules(securityGroupID string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		ec2CreateSGRule(securityGroupID, sgRuleInput{
			protocol:   "-1",
			cidrBlocks: []string{"0.0.0.0/0"},
			fromPort:   0,
			toPort:     0,
		}),
	}
}

func DefaultMasterSGIngressRules(masterSGID string, workerSGID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// master mcs
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			cidrBlocks: cidrBlocks,
			fromPort:   22623,
			toPort:     22623,
		}),
		// master icmp
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "icmp",
			cidrBlocks: cidrBlocks,
			fromPort:   -1,
			toPort:     -1,
		}),
		// master ssh
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			cidrBlocks: cidrBlocks,
			fromPort:   22,
			toPort:     22,
		}),
		// master https
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			cidrBlocks: cidrBlocks,
			fromPort:   6443,
			toPort:     6443,
		}),
		// master vxlan
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4789,
			toPort:   4789,
			self:     true,
		}),
		// master geneve
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 6081,
			toPort:   6081,
			self:     true,
		}),
		// master ike
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 500,
			toPort:   500,
			self:     true,
		}),
		// master ike nat_t
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4500,
			toPort:   4500,
			self:     true,
		}),
		// master esp
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "50",
			fromPort: 0,
			toPort:   0,
			self:     true,
		}),
		// master ovndb
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 6641,
			toPort:   6642,
			self:     true,
		}),
		// master internal
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// master internal udp
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// master kube scheduler
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10259,
			toPort:   10259,
			self:     true,
		}),
		// master kube controller manager
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10257,
			toPort:   10257,
			self:     true,
		}),
		// master kubelet secure
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10250,
			toPort:   10250,
			self:     true,
		}),
		// master etcd
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 2379,
			toPort:   2380,
			self:     true,
		}),
		// master services tcp
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// master services udp
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// master vxlan from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4789,
			toPort:     4789,
			sourceSGID: workerSGID,
		}),
		// master geneve from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   6081,
			toPort:     6081,
			sourceSGID: workerSGID,
		}),
		// master ike from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   500,
			toPort:     500,
			sourceSGID: workerSGID,
		}),
		// master ike nat_t from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4500,
			toPort:     4500,
			sourceSGID: workerSGID,
		}),
		// master esp from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "50",
			fromPort:   0,
			toPort:     0,
			sourceSGID: workerSGID,
		}),
		// master ovndb from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   6641,
			toPort:     6642,
			sourceSGID: workerSGID,
		}),
		// master internal from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: workerSGID,
		}),
		// master internal udp from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: workerSGID,
		}),
		// master kube scheduler from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10259,
			toPort:     10259,
			sourceSGID: workerSGID,
		}),
		// master kube controler manager from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10257,
			toPort:     10257,
			sourceSGID: workerSGID,
		}),
		// master kubelet secure from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10250,
			toPort:     10250,
			sourceSGID: workerSGID,
		}),
		// master services tcp from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: workerSGID,
		}),
		// master services udp from worker
		ec2CreateSGRule(masterSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: workerSGID,
		}),
	}
}

// DefaultWorkerSGIngressRules
func DefaultWorkerSGIngressRules(workerSGID string, masterSGID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// worker icmp
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "icmp",
			cidrBlocks: cidrBlocks,
			fromPort:   -1,
			toPort:     -1,
		}),
		// worker vxlan
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4789,
			toPort:   4789,
			self:     true,
		}),
		// worker geneve
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 6081,
			toPort:   6081,
			self:     true,
		}),
		// worker ike
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 500,
			toPort:   500,
			self:     true,
		}),
		// worker ike nat_t
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 4500,
			toPort:   4500,
			self:     true,
		}),
		// worker esp
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "50",
			fromPort: 0,
			toPort:   0,
			self:     true,
		}),
		// worker internal
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// worker internal udp
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 9000,
			toPort:   9999,
			self:     true,
		}),
		// worker kubelet insecure
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 10250,
			toPort:   10250,
			self:     true,
		}),
		// worker services tcp
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "tcp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// worker services udp
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol: "udp",
			fromPort: 30000,
			toPort:   32767,
			self:     true,
		}),
		// worker vxlan from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4789,
			toPort:     4789,
			sourceSGID: masterSGID,
		}),
		// worker geneve from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   6081,
			toPort:     6081,
			sourceSGID: masterSGID,
		}),
		// worker ike from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   500,
			toPort:     500,
			sourceSGID: masterSGID,
		}),
		// worker ike nat_t from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   4500,
			toPort:     4500,
			sourceSGID: masterSGID,
		}),
		// worker esp from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "50",
			fromPort:   0,
			toPort:     0,
			sourceSGID: masterSGID,
		}),
		// worker internal from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: masterSGID,
		}),
		// master internal udp from worker
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   9000,
			toPort:     9999,
			sourceSGID: masterSGID,
		}),
		// worker kubelet insecure from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   10250,
			toPort:     10250,
			sourceSGID: masterSGID,
		}),
		// worker services tcp from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: masterSGID,
		}),
		// worker services udp from master
		ec2CreateSGRule(workerSGID, sgRuleInput{
			protocol:   "udp",
			fromPort:   30000,
			toPort:     32767,
			sourceSGID: masterSGID,
		}),
	}
}

func DefaultAllowAllSGIngressRules(securityGroupID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		ec2CreateSGRule(securityGroupID, sgRuleInput{
			protocol:   "-1",
			cidrBlocks: []string{"0.0.0.0/0"},
		}),
	}
}
func DefaultBootstrapSGIngressRules(securityGroupID string, cidrBlocks []string) []*ec2.IpPermission {
	return []*ec2.IpPermission{
		// bootstrap ssh
		ec2CreateSGRule(securityGroupID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   22,
			toPort:     22,
			cidrBlocks: cidrBlocks,
		}),
		// bootstrap journald gateway
		ec2CreateSGRule(securityGroupID, sgRuleInput{
			protocol:   "tcp",
			fromPort:   19531,
			toPort:     19531,
			cidrBlocks: cidrBlocks,
		}),
	}
}

const (
	// tagNameSubnetInternalELB is the tag name used on a subnet to designate that
	// it should be used for internal ELBs
	tagNameSubnetInternalELB = "kubernetes.io/role/internal-elb"

	// tagNameSubnetPublicELB is the tag name used on a subnet to designate that
	// it should be used for internet ELBs
	tagNameSubnetPublicELB = "kubernetes.io/role/elb"
)

func (o *CreateInfraOptions) CreatePrivateSubnet(l *logrus.Logger, client ec2iface.EC2API, vpcID string, zone string, cidr string) (string, error) {
	tags := append(o.additionalEC2Tags, &ec2.Tag{
		Key:   aws.String(tagNameSubnetInternalELB),
		Value: aws.String("true"),
	})
	return o.CreateSubnet(l, client, vpcID, zone, cidr, fmt.Sprintf("%s-private-%s", o.InfraID, zone), tags)
}

func (o *CreateInfraOptions) CreatePublicSubnet(l *logrus.Logger, client ec2iface.EC2API, vpcID string, zone string, cidr string) (string, error) {
	tags := append(o.additionalEC2Tags, &ec2.Tag{
		Key:   aws.String(tagNameSubnetPublicELB),
		Value: aws.String("true"),
	})
	return o.CreateSubnet(l, client, vpcID, zone, cidr, fmt.Sprintf("%s-public-%s", o.InfraID, zone), tags)
}

func (o *CreateInfraOptions) CreateSubnet(l *logrus.Logger, client ec2iface.EC2API, vpcID, zone, cidr, name string, ec2Tags []*ec2.Tag) (string, error) {
	logger := l.WithField("name", name)
	subnetID, err := o.existingSubnet(client, name)
	if err != nil {
		return "", err
	}
	if len(subnetID) > 0 {
		logger.WithField("id", subnetID).Infoln("Found existing subnet")
		return subnetID, nil
	}

	subnetID, err = ec2CreateSubnet(client, name, zone, vpcID, cidr, ec2Tags)
	if err != nil {
		return "", err
	}
	logger.WithField("id", subnetID).Infoln("Created subnet")
	return subnetID, nil
}

func (o *CreateInfraOptions) existingSubnet(client ec2iface.EC2API, name string) (string, error) {
	return ec2GetSubnet(client, o.ec2Filters(name))
}

func (o *CreateInfraOptions) CreateNATGateway(l *logrus.Logger, client ec2iface.EC2API, publicSubnetID, availabilityZone string) (string, error) {
	natGatewayName := fmt.Sprintf("%s-nat-%s", o.InfraID, availabilityZone)
	natGateway, err := o.existingNATGateway(client, natGatewayName)
	if err != nil {
		return "", err
	}

	if natGateway != nil {
		l.WithField("id", aws.StringValue(natGateway.NatGatewayId)).Infoln("Found existing NAT gateway")
		return aws.StringValue(natGateway.NatGatewayId), nil
	}

	allocationID, err := ec2AllocateEIPAddress(client, append(ec2Tags(fmt.Sprintf("%s-eip-%s", o.InfraID, availabilityZone)), o.additionalEC2Tags...))
	if err != nil {
		return "", err
	}
	l.WithField("id", allocationID).Infoln("Created elastic IP for NAT gateway")

	natGatewayID, err := ec2CreateNatGateway(client, natGatewayName, allocationID, publicSubnetID, o.additionalEC2Tags)
	if err != nil {
		return "", err
	}
	l.WithField("id", natGatewayID).Infoln("Created NAT gateway")
	return natGatewayID, nil
}

func (o *CreateInfraOptions) existingNATGateway(client ec2iface.EC2API, name string) (*ec2.NatGateway, error) {
	return ec2GetNatGateway(client, o.ec2Filters(name))
}

func (o *CreateInfraOptions) CreatePrivateRouteTable(l *logrus.Logger, client ec2iface.EC2API, vpcID, natGatewayID, subnetID, zone string) (string, error) {
	tableName := fmt.Sprintf("%s-private-%s", o.InfraID, zone)
	routeTable, err := o.createRouteTable(l, client, vpcID, tableName)
	if err != nil {
		return "", err
	}
	tableID := aws.StringValue(routeTable.RouteTableId)
	logger := l.WithField("route table", tableID)

	// Everything below this is only needed if direct internet access is used
	if o.EnableProxy {
		return tableID, nil
	}

	if !ec2HasNATGatewayRoute(routeTable, natGatewayID) {
		if err := ec2CreateNatGatewayRoute(client, tableID, natGatewayID); err != nil {
			return "", err
		}
		logger.WithField("nat gateway", natGatewayID).Infoln("Created route to NAT gateway")
	} else {
		logger.WithField("nat gateway", natGatewayID).Infoln("Found existing route to NAT gateway")
	}
	if !ec2HasAssociatedSubnet(routeTable, subnetID) {
		if err := ec2AssociateRouteTable(client, tableID, subnetID); err != nil {
			return "", err
		}
		logger.WithField("subnet", subnetID).Infoln("Associated subnet with route table")
	} else {
		logger.WithField("subnet", subnetID).Infoln("Subnet already associated with route table")
	}
	return tableID, nil
}

func (o *CreateInfraOptions) CreatePublicRouteTable(l *logrus.Logger, client ec2iface.EC2API, vpcID, igwID string, subnetIDs []string) (string, error) {
	tableName := fmt.Sprintf("%s-public", o.InfraID)
	routeTable, err := o.createRouteTable(l, client, vpcID, tableName)
	if err != nil {
		return "", err
	}
	tableID := aws.StringValue(routeTable.RouteTableId)
	logger := l.WithField("route table", tableID)

	// Replace the VPC's main route table
	routeTableInfo, err := ec2GetRouteTable(client, []*ec2.Filter{
		ec2CreateFilter("vpc-id", vpcID),
		ec2CreateFilter("association.main", "true"),
	})
	if err != nil {
		return "", err
	}
	if routeTableInfo == nil {
		return "", fmt.Errorf("no route tables associated with the vpc")
	}
	// Replace route table association only if it's not the associated route table already
	if tID := aws.StringValue(routeTableInfo.RouteTableId); tID != tableID {
		var associationID string
		for _, assoc := range routeTableInfo.Associations {
			if aws.BoolValue(assoc.Main) {
				associationID = aws.StringValue(assoc.RouteTableAssociationId)
				break
			}
		}
		if err := ec2ReplaceRouteTableAssociation(client, tableID, associationID); err != nil {
			return "", err
		}
		logger.WithField("vpc", vpcID).Infoln("Set main VPC route table")
	}

	// Create route to internet gateway
	if !ec2HasInternetGatewayRoute(routeTable, igwID) {
		if err := ec2CreateRoute(client, tableID, igwID); err != nil {
			return "", err
		}
		logger.WithField("internet gateway", igwID).Infoln("Created route to internet gateway")
	} else {
		logger.WithField("internet gateway", igwID).Infoln("Found existing route to internet gateway")
	}

	// Associate the route table with the public subnet ID
	for _, subnetID := range subnetIDs {
		if !ec2HasAssociatedSubnet(routeTable, subnetID) {
			if err := ec2AssociateRouteTable(client, tableID, subnetID); err != nil {
				return "", err
			}
			logger.WithField("subnet", subnetID).Infoln("Associated route table with subnet")
		} else {
			logger.WithField("subnet", subnetID).Infoln("Found existing association between route table and subnet")
		}
	}
	return tableID, nil
}

func (o *CreateInfraOptions) createRouteTable(l *logrus.Logger, client ec2iface.EC2API, vpcID, name string) (*ec2.RouteTable, error) {
	table, err := o.existingRouteTable(l, client, name)
	if err != nil {
		return nil, err
	}
	if table == nil {
		table, err = ec2CreateRouteTable(client, name, vpcID, o.additionalEC2Tags)
		if err != nil {
			return nil, err
		}
	}
	l.WithField("name", name).WithField("id", aws.StringValue(table.RouteTableId)).Infoln("Created route table")
	return table, nil
}

func (o *CreateInfraOptions) existingRouteTable(l *logrus.Logger, client ec2iface.EC2API, name string) (*ec2.RouteTable, error) {
	table, err := ec2GetRouteTable(client, o.ec2Filters(name))
	if err != nil {
		return nil, err
	}
	if table != nil {
		l.WithField("name", name).WithField("id", aws.StringValue(table.RouteTableId)).Infoln("Found existing route table")
	}
	return table, nil
}

func (o *CreateInfraOptions) addTargetGroup(arn string) {
	for _, tg := range o.targetGroupARNs {
		if tg == arn {
			return
		}
	}
	o.targetGroupARNs = append(o.targetGroupARNs, arn)
}

func (o *CreateInfraOptions) createLB(l *logrus.Logger, elbClient *elbv2.ELBV2, name string, subnets []string, isPublic bool, elbTags []*elbv2.Tag) (*elbv2.LoadBalancer, error) {
	// Check if the internal load balancer already exists.
	lb, err := elbGetLoadBalancer(elbClient, name)
	if err != nil {
		return nil, err
	}
	if lb == nil {
		lb, err = elbCreateLoadBalancer(elbClient, name, subnets, isPublic, elbTags)
		if err != nil {
			return nil, err
		}
		l.WithField("name", aws.StringValue(lb.DNSName)).Infoln("Internal Load Balancer created")
	} else {
		l.WithField("name", aws.StringValue(lb.DNSName)).Infoln("Internal Load Balancer already exists")
	}
	return lb, nil
}

func (o *CreateInfraOptions) createTargeGroup(l *logrus.Logger, elbClient *elbv2.ELBV2, name, vpcID, healthCheckPath string, port int64, elbTags []*elbv2.Tag) (*elbv2.TargetGroup, error) {
	logger := l.WithField("target group", name)
	// Check if the target group already exists
	tg, err := elbGetTargetGroup(elbClient, name)
	if err != nil {
		return nil, err
	}
	if tg == nil {
		tg, err = elbCreateTargetGroup(elbClient, name, vpcID, healthCheckPath, port, elbTags)
		if err != nil {
			return nil, err
		}
		logger.WithField("arn", aws.StringValue(tg.TargetGroupArn)).Infoln("Target Group created")
	} else {
		logger.WithField("arn", aws.StringValue(tg.TargetGroupArn)).Infoln("Target Group already exists")
	}
	return tg, nil
}

func (o *CreateInfraOptions) createInternalLB(l *logrus.Logger, elbClient *elbv2.ELBV2, vpcID string, privateSubnets []string, elbTags []*elbv2.Tag) error {
	// Create internal LB.
	internalLBName := fmt.Sprintf("%s-int", o.InfraID)

	lb, err := o.createLB(l, elbClient, internalLBName, privateSubnets, false, elbTags)
	if err != nil {
		return err
	}
	internalLBARN := aws.StringValue(lb.LoadBalancerArn)
	o.LoadBalancers.Internal.ZoneID = aws.StringValue(lb.CanonicalHostedZoneId)
	o.LoadBalancers.Internal.DNSName = aws.StringValue(lb.DNSName)

	// Create InternalA TargetGroup.
	internalATGName := fmt.Sprintf("%s-aint", o.InfraID)

	tg, err := o.createTargeGroup(l, elbClient, internalATGName, vpcID, "/readyz", 6443, elbTags)
	if err != nil {
		return err
	}
	internalATGARN := aws.StringValue(tg.TargetGroupArn)
	o.addTargetGroup(internalATGARN)

	// Create InternalA Listener.
	internalAListenerName := fmt.Sprintf("%s-aint", o.InfraID)
	listener, err := elbCreateListener(elbClient, internalAListenerName, internalLBARN, internalATGARN, 6443, elbTags)
	if err != nil {
		return err
	}
	l.WithField("arn", aws.StringValue(listener.ListenerArn)).Infoln("Internal Listener created")

	// Create InternalS TargetGroup.
	internalSTGName := fmt.Sprintf("%s-sint", o.InfraID)
	stg, err := o.createTargeGroup(l, elbClient, internalSTGName, vpcID, "/healthz", 22623, elbTags)
	if err != nil {
		return err
	}
	internalSTGARN := aws.StringValue(stg.TargetGroupArn)
	o.addTargetGroup(internalSTGARN)

	internalSListenerName := fmt.Sprintf("%s-sint", o.InfraID)
	// Create a listener and associate it with the target group
	listener, err = elbCreateListener(elbClient, internalSListenerName, internalLBARN, internalSTGARN, 22623, elbTags)
	if err != nil {
		return err
	}
	l.WithField("arn", aws.StringValue(listener.ListenerArn)).Infoln("Internal Service Listener created")

	return nil
}

func (o *CreateInfraOptions) createExternalLB(l *logrus.Logger, elbClient *elbv2.ELBV2, vpcID string, publicSubnets []string, elbTags []*elbv2.Tag) error {
	// Create external LB.
	externalLBName := fmt.Sprintf("%s-ext", o.InfraID)
	lb, err := o.createLB(l, elbClient, externalLBName, publicSubnets, true, elbTags)
	if err != nil {
		return err
	}
	externalLBARN := aws.StringValue(lb.LoadBalancerArn)
	o.LoadBalancers.External.ZoneID = aws.StringValue(lb.CanonicalHostedZoneId)
	o.LoadBalancers.External.DNSName = aws.StringValue(lb.DNSName)

	// Create TargetGroup.
	externalTGName := fmt.Sprintf("%s-aext", o.InfraID)
	tg, err := o.createTargeGroup(l, elbClient, externalTGName, vpcID, "/readyz", 6443, elbTags)
	if err != nil {
		return err
	}
	externalTGARN := aws.StringValue(tg.TargetGroupArn)
	o.addTargetGroup(externalTGARN)

	externalListenerName := fmt.Sprintf("%s-aext", o.InfraID)
	elistener, err := elbCreateListener(elbClient, externalListenerName, externalLBARN, externalTGARN, 6443, elbTags)
	if err != nil {
		return err
	}
	l.WithField("arn", aws.StringValue(elistener.ListenerArn)).Infoln("External Listener created")

	return nil
}

func (o *CreateInfraOptions) CreateLoadBalancers(l *logrus.Logger, session *session.Session, vpcID string, privateSubnets []string, publicSubnets []string, external bool) error {
	elbClient := elbv2.New(session)
	elbTags := elbCreateTags(o.AdditionalTags)
	if err := o.createInternalLB(l, elbClient, vpcID, privateSubnets, elbTags); err != nil {
		return err
	}

	if !external {
		l.Debugln("Skipping creation of a public LB because of private cluster")
		return nil
	}

	if err := o.createExternalLB(l, elbClient, vpcID, publicSubnets, elbTags); err != nil {
		return err
	}

	return nil
}

func (o *CreateInfraOptions) existingVPCS3Endpoint(client ec2iface.EC2API) (string, error) {
	return ec2GetVPCS3Endpoint(client, o.ec2Filters(""))
}

func (o *CreateInfraOptions) CreateVPCS3Endpoint(l *logrus.Logger, client ec2iface.EC2API, vpcID string, routeTableIds []*string) error {
	existingEndpoint, err := o.existingVPCS3Endpoint(client)
	if err != nil {
		return err
	}
	if len(existingEndpoint) > 0 {
		l.WithField("id", existingEndpoint).Infoln("Found existing s3 VPC endpoint")
		return nil
	}
	serviceName := fmt.Sprintf("com.amazonaws.%s.s3", o.Region)
	endpoint, err := ec2CreateVPCS3Endpoint(client, "", vpcID, serviceName, routeTableIds, o.additionalEC2Tags)
	if err != nil {
		return err
	}
	if endpoint != nil {
		l.WithField("id", aws.StringValue(endpoint.VpcEndpointId)).Infoln("Created s3 VPC endpoint")
	}
	return nil
}

// CreateInfraOptions
// TODO(alberto): this is brought form hypershift just to satisfy current functions.
// Feel free to model this as you best see fit.
type CreateInfraOptions struct {
	Region             string
	InfraID            string
	AWSCredentialsFile string
	AWSKey             string
	AWSSecretKey       string
	Name               string
	BaseDomain         string
	BaseDomainPrefix   string
	Zones              []string
	OutputFile         string
	AdditionalTags     map[string]string
	EnableProxy        bool
	SSHKeyFile         string
	additionalEC2Tags  []*ec2.Tag
	cidrV4Blocks       []string
	cidrV6Blocks       []string

	// Additional output for consumers.
	vpcID            string
	public           bool
	publicSubnetIDs  []string
	privateSubnetIDs []string
	LoadBalancers    struct {
		Internal struct {
			DNSName string
			ZoneID  string
		}
		External struct {
			DNSName string
			ZoneID  string
		}
	}
	targetGroupARNs          []string
	bootstrapSecurityGroupID string
	masterSecurityGroupID    string
	workerSecurityGroupID    string
	zoneToSubnetIDMap        map[string]string
}

// CreateInfraOutput
// TODO(alberto): this is brought form hypershift just to satisfy current functions.
// Feel free to model this as you best see fit. E.g. see dnsInput or boostrapInput.
//type CreateInfraOutput struct {
//	Region           string                   `json:"region"`
//	Zone             string                   `json:"zone"`
//	InfraID          string                   `json:"infraID"`
//	MachineCIDR      string                   `json:"machineCIDR"`
//	VPCID            string                   `json:"vpcID"`
//	Zones            []*CreateInfraOutputZone `json:"zones"`
//	SecurityGroupID  string                   `json:"securityGroupID"`
//	Name             string                   `json:"Name"`
//	BaseDomain       string                   `json:"baseDomain"`
//	BaseDomainPrefix string                   `json:"baseDomainPrefix"`
//	PublicZoneID     string                   `json:"publicZoneID"`
//	PrivateZoneID    string                   `json:"privateZoneID"`
//	LocalZoneID      string                   `json:"localZoneID"`
//	ProxyAddr        string                   `json:"proxyAddr"`
//}
//
//type CreateInfraOutputZone struct {
//	Name     string `json:"name"`
//	SubnetID string `json:"subnetID"`
//}
