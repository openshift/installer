package clusterapi

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// Provider is the vSphere implementation of the clusterapi InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

var _ clusterapi.InfraReadyProvider = (*Provider)(nil)
var _ clusterapi.Provider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)

// Name returns the PowerVS provider name.
func (p Provider) Name() string {
	return powervstypes.Name
}

func leftInContext(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return math.MaxInt64
	}

	duration := time.Until(deadline)

	return duration
}

// InfraReady is called once cluster.Status.InfrastructureReady
// is true, typically after load balancers have been provisioned. It can be used
// to create DNS records.
func (p Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	var (
		client      *powervsconfig.Client
		vpcRegion   string
		instanceCRN string
		rules       *vpcv1.SecurityGroupRuleCollection
		rule        *vpcv1.SecurityGroupRulePrototype
		wantedPorts = sets.New[int64](22, 10258, 22623)
		foundPorts  = sets.Set[int64]{}
		err         error
	)

	logrus.Debugf("InfraReady: in = %+v", in)
	logrus.Debugf("InfraReady: in.InstallConfig.Config = %+v", in.InstallConfig.Config)
	logrus.Debugf("InfraReady: in.InstallConfig.PowerVS = %+v", in.InstallConfig.PowerVS)

	powerVSCluster := &capibm.IBMPowerVSCluster{}

	// Get the cluster from the provider
	key := crclient.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	logrus.Debugf("InfraReady: cluster key = %+v", key)
	if err = in.Client.Get(ctx, key, powerVSCluster); err != nil {
		return fmt.Errorf("failed to get PowerVS cluster in InfraReady: %w", err)
	}
	logrus.Debugf("InfraReady: powerVSCluster = %+v", powerVSCluster)
	logrus.Debugf("InfraReady: powerVSCluster.Status = %+v", powerVSCluster.Status)
	if powerVSCluster.Status.VPC == nil || powerVSCluster.Status.VPC.ID == nil {
		return fmt.Errorf("vpc is empty in InfraReady?")
	}
	logrus.Debugf("InfraReady: powerVSCluster.Status.VPC.ID = %s", *powerVSCluster.Status.VPC.ID)

	// Get the image from the provider
	key = crclient.ObjectKey{
		Name:      fmt.Sprintf("rhcos-%s", in.InfraID),
		Namespace: capiutils.Namespace,
	}
	logrus.Debugf("InfraReady: image key = %+v", key)
	powerVSImage := &capibm.IBMPowerVSImage{}
	if err = in.Client.Get(ctx, key, powerVSImage); err != nil {
		return fmt.Errorf("failed to get PowerVS image in InfraReady: %w", err)
	}
	logrus.Debugf("InfraReady: image = %+v", powerVSImage)

	// SAD: client in the Metadata struct is lowercase and therefore private
	// client = in.InstallConfig.PowerVS.client
	client, err = powervsconfig.NewClient()
	if err != nil {
		return fmt.Errorf("failed to get NewClient in InfraReady: %w", err)
	}
	logrus.Debugf("InfraReady: NewClient returns %+v", client)

	// We need to set the region we will eventually query inside
	vpcRegion = in.InstallConfig.Config.Platform.PowerVS.VPCRegion
	if vpcRegion == "" {
		vpcRegion, err = powervstypes.VPCRegionForPowerVSRegion(in.InstallConfig.Config.Platform.PowerVS.Region)
		if err != nil {
			return fmt.Errorf("failed to get VPC region (%s) in InfraReady: %w", vpcRegion, err)
		}
	}
	logrus.Debugf("InfraReady: vpcRegion = %s", vpcRegion)
	if err = client.SetVPCServiceURLForRegion(ctx, vpcRegion); err != nil {
		return fmt.Errorf("failed to set the VPC service region (%s) in InfraReady: %w", vpcRegion, err)
	}

	// Step 1.
	// Create DNS records for the two load balancers
	// map[string]VPCLoadBalancerStatus
	instanceCRN, err = client.GetInstanceCRNByName(ctx,
		in.InstallConfig.PowerVS.BaseDomain,
		in.InstallConfig.Config.Publish)
	if err != nil {
		return fmt.Errorf("failed to get InstanceCRN (%s) by name in InfraReady: %w",
			in.InstallConfig.Config.Publish,
			err)
	}
	logrus.Debugf("InfraReady: instanceCRN = %s", instanceCRN)

	lbExtExp := regexp.MustCompile(`\b-loadbalancer\b$`)
	lbIntExp := regexp.MustCompile(`\b-loadbalancer-int\b$`)

	for lbKey, loadBalancerStatus := range powerVSCluster.Status.LoadBalancers {
		var (
			idx      int
			substr   string
			infraID  string
			hostname string
			prefix   string
		)

		// The infra id is "rdr-hamzy-test-dal10-846vd" and we need "rdr-hamzy-test-dal10"
		logrus.Debugf("in.InfraID = %s", in.InfraID)
		idx = strings.LastIndex(in.InfraID, "-")
		logrus.Debugf("idx = %d", idx)
		substr = in.InfraID[idx:]
		logrus.Debugf("substr = %s", substr)
		infraID = strings.ReplaceAll(in.InfraID, substr, "")
		logrus.Debugf("infraID = %s", infraID)

		// Is it external (public) or internal (private)?
		logrus.Debugf("lbKey = %s", lbKey)
		switch {
		case lbExtExp.MatchString(lbKey):
			prefix = "api."
		case lbIntExp.MatchString(lbKey):
			prefix = "api-int."
		}
		logrus.Debugf("prefix = %s", prefix)

		hostname = fmt.Sprintf("%s%s", prefix, infraID)

		logrus.Debugf("InfraReady: crn = %s, base domain = %s, hostname = %s, cname = %s",
			instanceCRN,
			in.InstallConfig.PowerVS.BaseDomain,
			hostname,
			*loadBalancerStatus.Hostname)

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := client.CreateDNSRecord(ctx,
				instanceCRN,
				in.InstallConfig.PowerVS.BaseDomain,
				hostname,
				*loadBalancerStatus.Hostname)
			if err2 == nil {
				return true, nil
			}
			return false, err2
		})
		if err != nil {
			return fmt.Errorf("failed to create a DNS CNAME record (%s, %s): %w",
				hostname,
				*loadBalancerStatus.Hostname,
				err)
		}
	}

	// Step 2.
	// See which ports are already allowed.
	rules, err = client.ListSecurityGroupRules(ctx, *powerVSCluster.Status.VPC.ID)
	if err != nil {
		return fmt.Errorf("failed to list security group rules: %w", err)
	}

	for _, existingRule := range rules.Rules {
		switch reflect.TypeOf(existingRule).String() {
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
			securityGroupRule, ok := existingRule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
			if !ok {
				return fmt.Errorf("could not convert to ProtocolTcpudp")
			}
			logrus.Debugf("InfraReady: VPC has rule: direction = %s, proto = %s, min = %d, max = %d",
				*securityGroupRule.Direction,
				*securityGroupRule.Protocol,
				*securityGroupRule.PortMin,
				*securityGroupRule.PortMax)
			if *securityGroupRule.Direction == "inbound" &&
				*securityGroupRule.Protocol == "tcp" {
				foundPorts.Insert(*securityGroupRule.PortMin)
			}
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		}
	}
	logrus.Debugf("InfraReady: foundPorts = %+v", foundPorts)
	logrus.Debugf("InfraReady: wantedPorts = %+v", wantedPorts)
	logrus.Debugf("InfraReady: wantedPorts.Difference(foundPorts) = %+v", wantedPorts.Difference(foundPorts))

	// Step 3.
	// Add to security group rules
	for port := range wantedPorts.Difference(foundPorts) {
		rule = &vpcv1.SecurityGroupRulePrototype{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			PortMin:   ptr.To(port),
			PortMax:   ptr.To(port),
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			logrus.Debugf("InfraReady: Adding port %d to security group rule to %v",
				port,
				*powerVSCluster.Status.VPC.ID)
			err2 := client.AddSecurityGroupRule(ctx, *powerVSCluster.Status.VPC.ID, rule)
			if err2 == nil {
				return true, nil
			}
			return false, err2
		})
		if err != nil {
			return fmt.Errorf("failed to add security group rule for port %d: %w", port, err)
		}
	}

	// Allow ping so we can debug
	rule = &vpcv1.SecurityGroupRulePrototype{
		Direction: ptr.To("inbound"),
		Protocol:  ptr.To("icmp"),
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		err2 := client.AddSecurityGroupRule(ctx, *powerVSCluster.Status.VPC.ID, rule)
		if err2 == nil {
			return true, nil
		}
		return false, err2
	})
	if err != nil {
		return fmt.Errorf("failed to add security group rule for icmp: %w", err)
	}

	return nil
}

func findMachineAddress(ctx context.Context, in clusterapi.PostProvisionInput, key crclient.ObjectKey) (string, error) {
	powerVSMachine := &capibm.IBMPowerVSMachine{}

	// Get the machine address
	// Unfortunately https://pkg.go.dev/k8s.io/apimachinery/pkg/util/wait#PollUntilContextCancel
	// can only return a bool.  It would be nice if it could return a pointer.
	if err := wait.PollUntilContextCancel(ctx, time.Second*10,
		false,
		func(ctx context.Context) (bool, error) {
			if err := in.Client.Get(ctx, key, powerVSMachine); err != nil {
				return false, fmt.Errorf("failed to get PowerVS machine in PostProvision: %w", err)
			}

			for _, address := range powerVSMachine.Status.Addresses {
				if address.Type == corev1.NodeInternalIP {
					logrus.Debugf("PostProvision: found %s address %s", key.Name, address.Address)
					return true, nil
				}
			}

			logrus.Debugf("PostProvision: waiting for %s machine", key.Name)
			return false, nil
		}); err != nil {
		return "", err
	}

	if err := in.Client.Get(ctx, key, powerVSMachine); err != nil {
		return "", fmt.Errorf("failed to get PowerVS machine in PostProvision: %w", err)
	}

	for _, address := range powerVSMachine.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			return address.Address, nil
		}
	}

	return "", fmt.Errorf("failed to get machine %s IP address", key.Name)
}

// PostProvision should be called to add or update PowerVS resources after provisioning has completed.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	var (
		client             *powervsconfig.Client
		ipAddr             string
		refServiceInstance *capibm.IBMPowerVSResourceReference
		sshKeyName         string
		err                error
		instanceID         *string
		fieldType          string
	)

	// SAD: client in the Metadata struct is lowercase and therefore private
	// client = in.InstallConfig.PowerVS.client
	client, err = powervsconfig.NewClient()
	if err != nil {
		return fmt.Errorf("failed to get NewClient in PostProvision: %w", err)
	}
	logrus.Debugf("PostProvision: NewClient returns %+v", client)

	// Step 1.
	// Wait until bootstrap and master nodes have IP addresses.  This will verify
	// that the Transit Gateway and DHCP server work correctly before continuing on.

	// Get master IP addresses
	masterCount := int64(1)
	if reps := in.InstallConfig.Config.ControlPlane.Replicas; reps != nil {
		masterCount = *reps
	}
	logrus.Debugf("PostProvision: masterCount = %d", masterCount)
	for i := int64(0); i < masterCount; i++ {
		key := crclient.ObjectKey{
			Name:      fmt.Sprintf("%s-master-%d", in.InfraID, i),
			Namespace: capiutils.Namespace,
		}
		if ipAddr, err = findMachineAddress(ctx, in, key); err != nil {
			return err
		}
		logrus.Debugf("PostProvision: %s ipAddr = %v", key.Name, ipAddr)
	}

	// Get the bootstrap machine from the provider
	key := crclient.ObjectKey{
		Name:      fmt.Sprintf("%s-bootstrap", in.InfraID),
		Namespace: capiutils.Namespace,
	}
	logrus.Debugf("PostProvision: machine key = %+v", key)

	// Find its address
	if ipAddr, err = findMachineAddress(ctx, in, key); err != nil {
		return err
	}
	logrus.Debugf("PostProvision: ipAddr = %v", ipAddr)

	// Get information about it
	powerVSMachine := &capibm.IBMPowerVSMachine{}
	if err := in.Client.Get(ctx, key, powerVSMachine); err != nil {
		return fmt.Errorf("failed to get PowerVS bootstrap machine in PostProvision: %w", err)
	}
	logrus.Debugf("PostProvision: machine = %+v", powerVSMachine)

	// Specifically the Power Virtual Server (PVS)
	logrus.Debugf("PostProvision: machine.Spec.ServiceInstance = %+v", powerVSMachine.Spec.ServiceInstance)
	refServiceInstance = powerVSMachine.Spec.ServiceInstance

	// Step 2.
	// Create worker ssh key in the PVS
	if in.InstallConfig.Config.SSHKey == "" {
		return fmt.Errorf("install config's ssh key is empty?")
	}

	sshKeyName = fmt.Sprintf("%s-key", in.InfraID)

	switch {
	case refServiceInstance.ID != nil:
		logrus.Debugf("PostProvision: CreateSSHKey: si id = %s, key = %s",
			*refServiceInstance.ID,
			in.InstallConfig.Config.SSHKey)
		instanceID = refServiceInstance.ID
		fieldType = "ID"
	case refServiceInstance.Name != nil:
		logrus.Debugf("PostProvision: CreateSSHKey: si name = %s, key = %s",
			*refServiceInstance.Name,
			in.InstallConfig.Config.SSHKey)
		guid, err := client.ServiceInstanceNameToGUID(ctx, *refServiceInstance.Name)
		if err != nil {
			return fmt.Errorf("failed to find id for ServiceInstance name %s: %w",
				*refServiceInstance.Name,
				err)
		}
		logrus.Debugf("PostProvision: CreateSSHKey: guid = %s", guid)
		instanceID = ptr.To(guid)
		fieldType = "Name"
	default:
		return fmt.Errorf("could not handle powerVSMachine.Spec.ServiceInstance")
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32,
	}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		err2 := client.CreateSSHKey(ctx,
			*instanceID,
			*powerVSMachine.Status.Zone,
			sshKeyName,
			in.InstallConfig.Config.SSHKey)
		if err2 == nil {
			return true, nil
		}
		return false, err2
	})
	if err != nil {
		return fmt.Errorf("failed to add SSH key for the workers(%s): %w", fieldType, err)
	}

	return nil
}
