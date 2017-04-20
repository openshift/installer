package cloudforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/coreos/ipnets"

	"github.com/coreos/tectonic-installer/installer/server/defaults"
)

const (
	vpcLogicalName = "VPC"
	defaultVPCCIDR = "10.0.0.0/16"
)

// SetDefaults sets Config default values (idempotent).
func (c *Config) SetDefaults() {
	defInt := func(field *int, def int) {
		if *field == 0 {
			*field = def
		}
	}
	defString := func(field *string, def string) {
		if *field == "" {
			*field = def
		}
	}

	// Kubernetes
	defString(&c.PodCIDR, defaults.PodCIDR)
	defString(&c.ServiceCIDR, defaults.ServiceCIDR)

	// AWS Cloud Formation
	defString(&c.VPCCIDR, defaultVPCCIDR)

	defInt(&c.ETCDCount, 3)
	defString(&c.ETCDInstanceType, "m3.medium")
	defString(&c.ETCDRootVolumeType, "gp2")
	defInt(&c.ETCDRootVolumeIOPS, 0)
	defInt(&c.ETCDRootVolumeSize, 30)

	defInt(&c.ControllerCount, 1)
	defString(&c.ControllerInstanceType, "t2.medium")
	defString(&c.ControllerRootVolumeType, "gp2")
	defInt(&c.ControllerRootVolumeIOPS, 0)
	defInt(&c.ControllerRootVolumeSize, 30)

	defInt(&c.WorkerCount, 2)
	defString(&c.WorkerInstanceType, "t2.medium")
	defString(&c.WorkerRootVolumeType, "gp2")
	defInt(&c.WorkerRootVolumeIOPS, 0)
	defInt(&c.WorkerRootVolumeSize, 30)
}

// GetDefaultSubnets partitions a CIDR into subnets
func GetDefaultSubnets(sess *session.Session, vpcCIDR string) ([]VPCSubnet, []VPCSubnet, error) {
	zones, err := getAvailabilityZones(sess)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting availability zone %v", err)
	}

	_, vpcNet, err := net.ParseCIDR(vpcCIDR)
	if vpcNet == nil || err != nil {
		return nil, nil, fmt.Errorf("failed parsing VPC CIDR %v", err)
	}

	// Calculate subnetMultipler many times as many subnets as needed to
	// intentionally leave unused IPs for unspecified use. A multipler
	// of 1 divides the VPC among AZs, 2 leaves 50% of the VPC unallocated,
	// 4 leaves 75% unallocated, etc.
	cidrs, err := ipnets.SubnetInto(vpcNet, 2*2*len(zones))
	if err != nil {
		return nil, nil, fmt.Errorf("failed dividing VPC into subnets %v", err)
	}

	controllerSubnets := make([]VPCSubnet, len(zones))
	workerSubnets := make([]VPCSubnet, len(zones))

	// add generated multi-AZ subnets for controllers
	for i, zone := range zones {
		controllerSubnets[i] = VPCSubnet{
			AvailabilityZone: zone,
			InstanceCIDR:     cidrs[i].String(),
		}
	}

	// add generated multi-AZ subnets for workers
	for i, zone := range zones {
		workerSubnets[i] = VPCSubnet{
			AvailabilityZone: zone,
			InstanceCIDR:     cidrs[i+len(zones)].String(),
		}
	}

	return controllerSubnets, workerSubnets, nil
}

// SetComputed populates computed fields and may make calls to AWS endpoints.
func (c *Config) SetComputed(sess *session.Session) error {
	c.SecureAPIServers = fmt.Sprintf("https://%s:443", c.ControllerDomain)

	// Self-hosted Kubernetes IPs
	var err error
	c.APIServiceIP, err = defaults.APIServiceIP(c.ServiceCIDR)
	if err != nil {
		return err
	}
	c.DNSServiceIP, err = defaults.DNSServiceIP(c.ServiceCIDR)
	if err != nil {
		return err
	}

	// compute AMI from channel
	if !(c.Channel == "stable" || c.Channel == "beta" || c.Channel == "alpha") {
		return fmt.Errorf("channel must be stable, beta, or alpha, got %s", c.Channel)
	}
	if c.AMI, err = getAMI(c.Region, c.Channel); err != nil {
		return fmt.Errorf("failed getting AMI for config: %v - %v: %v", c.Region, c.Channel, err)
	}

	controllerSubnets, workerSubnets, err := GetDefaultSubnets(sess, c.VPCCIDR)
	if err != nil {
		return err
	}
	if len(c.ControllerSubnets) == 0 {
		c.ControllerSubnets = controllerSubnets
	}
	if len(c.WorkerSubnets) == 0 {
		c.WorkerSubnets = workerSubnets
	}

	// determine whether to create new subnets in the stack formation and to give
	// them indexed names
	c.CreateControllerSubnets = false
	for i := 0; i < len(c.ControllerSubnets); i++ {
		if c.ControllerSubnets[i].ID == "" {
			c.CreateControllerSubnets = true

			if c.ControllerSubnets[i].Name == "" {
				c.ControllerSubnets[i].Name = fmt.Sprintf("ControllerSubnet%d", i)
			}
		}
	}
	c.CreateWorkerSubnets = false
	for i := 0; i < len(c.WorkerSubnets); i++ {
		if c.WorkerSubnets[i].ID == "" {
			c.CreateWorkerSubnets = true

			if c.WorkerSubnets[i].Name == "" {
				c.WorkerSubnets[i].Name = fmt.Sprintf("WorkerSubnet%d", i)
			}
		}
	}

	// retrieve the hosted zone name
	c.HostedZoneName, err = getHostedZoneName(sess, c.HostedZoneID)
	if err != nil {
		return err
	}

	if c.ExternalETCDClient == "" {
		c.ETCDInstances, c.ETCDInitialCluster, c.ETCDEndpoints = PopulateETCDInstances(
			c.ClusterName,
			c.HostedZoneName,
			c.WorkerSubnets,
			c.ETCDCount,
		)
	} else {
		c.ETCDEndpoints = c.ExternalETCDClient
	}

	err = PopulateCIDRs(sess, c.VPCID, c.ControllerSubnets, c.WorkerSubnets)
	if err != nil {
		return err
	}

	// set logical name constants
	c.VPCLogicalName = vpcLogicalName
	c.InternetGatewayLogicalName = "VPCInternetGateway"

	// default to creating VPC / gateway with logical name
	c.VPCRef = fmt.Sprintf(`{ "Ref" : %q }`, c.VPCLogicalName)
	c.InternetGatewayRef = fmt.Sprintf(`{ "Ref" : %q}`, c.InternetGatewayLogicalName)

	// for existing VPC, set the VPC and internet gateway references to use
	if c.VPCID != "" {
		gateway, err := getInternetGateway(sess, c.VPCID)
		if err != nil {
			return err
		}

		c.VPCRef = fmt.Sprintf("%q", c.VPCID)
		c.InternetGatewayRef = fmt.Sprintf("%q", aws.StringValue(gateway.InternetGatewayId))
	}

	return nil
}

// PopulateETCDInstances initializes a slice of ETCDInstance, with the
// corresponding initial-cluster and endpoint variables.
func PopulateETCDInstances(clusterName, hostedZoneName string, subnets []VPCSubnet, count int) (instances []ETCDInstance, initialCluster, endpoints string) {
	var endpointsS, initialClusterS []string

	for i := 0; i < count; i++ {
		domainName := fmt.Sprintf("%s-etcd-%d.%s", clusterName, i, hostedZoneName)
		instances = append(instances, ETCDInstance{
			Name:       fmt.Sprintf("etcd%d", i),
			DomainName: domainName,
			Subnet:     subnets[i%len(subnets)],
		})
		endpointsS = append(endpointsS, domainName+":2379")
		initialClusterS = append(initialClusterS, fmt.Sprintf("etcd%d=http://%s:2380", i, domainName))
	}

	endpoints = strings.Join(endpointsS, ",")
	initialCluster = strings.Join(initialClusterS, ",")
	return
}

// PopulateCIDRs shoves some CIDRs into subnets when we know the IDs
func PopulateCIDRs(sess *session.Session, existingVPCID string, publicSubnets, privateSubnets []VPCSubnet) error {
	existingPublicSubnets, existingPrivateSubnets, err := GetVPCSubnets(sess, existingVPCID)
	if err != nil {
		return err
	}

	existingSubnets := append(existingPublicSubnets, existingPrivateSubnets...)
	for i, subnet := range publicSubnets {
		if subnet.ID == "" || subnet.InstanceCIDR != "" {
			continue
		}
		for _, existing := range existingSubnets {
			if subnet.ID == existing.ID {
				publicSubnets[i].InstanceCIDR = existing.InstanceCIDR
				break
			}
		}
	}
	for i, subnet := range privateSubnets {
		if subnet.ID == "" || subnet.InstanceCIDR != "" {
			continue
		}
		for _, existing := range existingSubnets {
			if subnet.ID == existing.ID {
				privateSubnets[i].InstanceCIDR = existing.InstanceCIDR
				break
			}
		}
	}
	return nil
}

// Valid returns true if the cloudform Config is valid.
func (c *Config) Valid() error {
	if c.ELBScheme != "internet-facing" && c.ELBScheme != "internal" {
		return errors.New("elbScheme must be either 'internet-facing' or 'internal'")
	}
	if c.ControllerDomain == "" {
		return errors.New("controllerDomain must be set")
	}
	if c.TectonicDomain == "" {
		return errors.New("tectonicDomain must be set")
	}
	if c.HostedZoneID == "" {
		return errors.New("hostedZoneID must be specified")
	}
	if c.KeyName == "" {
		return errors.New("keyName must be set")
	}
	if c.Region == "" {
		return errors.New("region must be set")
	}
	if len(c.ClusterName) == 0 || len(c.ClusterName) > 28 {
		return errors.New("clusterName must be between 1 and 28 characters")
	}

	if c.KMSKeyARN == "" {
		return errors.New("kmsKeyArn must be set")
	}

	if c.VPCID == "" && c.RouteTableID != "" {
		return errors.New("vpcId must be specified if routeTableId is specified")
	}

	_, vpcNet, err := net.ParseCIDR(c.VPCCIDR)
	if err != nil {
		return fmt.Errorf("invalid vpcCIDR: %v", err)
	}

	// pod network / vpc overlap
	_, podNet, err := net.ParseCIDR(c.PodCIDR)
	if err != nil {
		return fmt.Errorf("invalid podCIDR: %v", err)
	}
	if cidrOverlap(podNet, vpcNet) {
		return fmt.Errorf("vpcCIDR (%s) overlaps with podCIDR (%s)", c.VPCCIDR, c.PodCIDR)
	}

	// service / vpc network overlap
	_, serviceNet, err := net.ParseCIDR(c.ServiceCIDR)
	if err != nil {
		return fmt.Errorf("invalid serviceCIDR: %v", err)
	}
	if cidrOverlap(serviceNet, vpcNet) {
		return fmt.Errorf("vpcCIDR (%s) overlaps with serviceCIDR (%s)", c.VPCCIDR, c.ServiceCIDR)
	}

	// service / pod network overlap
	if cidrOverlap(serviceNet, podNet) {
		return fmt.Errorf("serviceCIDR (%s) overlaps with podCIDR (%s)", c.ServiceCIDR, c.PodCIDR)
	}

	if c.ControllerRootVolumeType == "io1" {
		if c.ControllerRootVolumeIOPS < 100 || c.ControllerRootVolumeIOPS > 2000 {
			return fmt.Errorf("invalid controllerRootVolumeIOPS: %d", c.ControllerRootVolumeIOPS)
		}
	} else {
		if c.ControllerRootVolumeIOPS != 0 {
			return fmt.Errorf("invalid controllerRootVolumeIOPS for volume type '%s': %d", c.ControllerRootVolumeType, c.ControllerRootVolumeIOPS)
		}

		if c.ControllerRootVolumeType != "standard" && c.ControllerRootVolumeType != "gp2" {
			return fmt.Errorf("invalid controllerRootVolumeType: %s", c.ControllerRootVolumeType)
		}
	}

	if c.WorkerRootVolumeType == "io1" {
		if c.WorkerRootVolumeIOPS < 100 || c.WorkerRootVolumeIOPS > 2000 {
			return fmt.Errorf("invalid workerRootVolumeIOPS: %d", c.WorkerRootVolumeIOPS)
		}
	} else {
		if c.WorkerRootVolumeIOPS != 0 {
			return fmt.Errorf("invalid workerRootVolumeIOPS for volume type '%s': %d", c.WorkerRootVolumeType, c.WorkerRootVolumeIOPS)
		}

		if c.WorkerRootVolumeType != "standard" && c.WorkerRootVolumeType != "gp2" {
			return fmt.Errorf("invalid workerRootVolumeType: %s", c.WorkerRootVolumeType)
		}
	}

	return nil
}

/*
Returns the availability zones referenced by the cluster configuration
*/
func (c *Config) availabilityZones() []string {

	azs := make([]string, len(c.ControllerSubnets))
	for i := range azs {
		azs[i] = c.ControllerSubnets[i].AvailabilityZone
	}

	return azs
}

// ValidateJSON returns detailed JSON validation errors.
func validateJSON(data []byte) error {
	// use unmarshal function to do syntax validation
	var jsonHolder map[string]interface{}
	if err := json.Unmarshal(data, &jsonHolder); err != nil {
		// attempt to provide more detail
		syntaxError, ok := err.(*json.SyntaxError)
		if ok {
			contextString := getContextString(data, int(syntaxError.Offset), 3)
			return fmt.Errorf("%v:\njson syntax error (offset=%d), in this region:\n-------\n%s\n-------\n", err, syntaxError.Offset, contextString)
		}
		return err
	}
	return nil
}

func getContextString(buf []byte, offset, lineCount int) string {
	linesSeen := 0
	var leftLimit int
	for leftLimit = offset; leftLimit > 0 && linesSeen <= lineCount; leftLimit-- {
		if buf[leftLimit] == '\n' {
			linesSeen++
		}
	}
	linesSeen = 0
	var rightLimit int
	for rightLimit = offset + 1; rightLimit < len(buf) && linesSeen <= lineCount; rightLimit++ {
		if buf[rightLimit] == '\n' {
			linesSeen++
		}
	}
	return string(buf[leftLimit:rightLimit])
}

// Does the address space of these networks "a" and "b" overlap?
func cidrOverlap(a, b *net.IPNet) bool {
	return a.Contains(b.IP) || b.Contains(a.IP)
}

func withTrailingDot(s string) string {
	if s == "" {
		return s
	}
	lastRune, _ := utf8.DecodeLastRuneInString(s)
	if lastRune != rune('.') {
		return s + "."
	}
	return s
}
