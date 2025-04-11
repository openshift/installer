package ocm

import (
	"fmt"
	"net"

	msv1 "github.com/openshift-online/ocm-sdk-go/servicemgmt/v1"

	"github.com/openshift/rosa/pkg/fedramp"
)

var fedrampError = fmt.Errorf("managed services are not supported for FedRAMP clusters")

type CreateManagedServiceArgs struct {
	ServiceType string
	ClusterName string

	Parameters map[string]string
	Properties map[string]string

	AwsAccountID           string
	AwsRoleARN             string
	AwsSupportRoleARN      string
	AwsControlPlaneRoleARN string
	AwsWorkerRoleARN       string
	AwsRegion              string

	AwsOperatorIamRoleList []OperatorIAMRole

	// Custom network configuration
	MultiAZ           bool
	Privatelink       bool
	AvailabilityZones []string
	SubnetIDs         []string
	MachineCIDR       net.IPNet
	PodCIDR           net.IPNet
	ServiceCIDR       net.IPNet
	HostPrefix        int

	// create a fake cluster with no aws resources
	FakeCluster bool
}

func (c *Client) CreateManagedService(args CreateManagedServiceArgs) (*msv1.ManagedService, error) {
	if fedramp.Enabled() {
		return nil, fedrampError
	}

	operatorIamRoles := []*msv1.OperatorIAMRoleBuilder{}
	for _, operatorIAMRole := range args.AwsOperatorIamRoleList {
		operatorIamRoles = append(operatorIamRoles,
			msv1.NewOperatorIAMRole().
				Name(operatorIAMRole.Name).
				Namespace(operatorIAMRole.Namespace).
				RoleARN(operatorIAMRole.RoleARN))
	}

	parameters := []*msv1.ServiceParameterBuilder{}
	for id, val := range args.Parameters {
		parameters = append(parameters,
			msv1.NewServiceParameter().ID(id).Value(val))
	}

	var listeningMethod msv1.ListeningMethod
	if args.Privatelink {
		listeningMethod = msv1.ListeningMethodInternal
	} else {
		listeningMethod = msv1.ListeningMethodExternal
	}

	var network *msv1.NetworkBuilder
	if !IsEmptyCIDR(args.MachineCIDR) ||
		!IsEmptyCIDR(args.PodCIDR) ||
		!IsEmptyCIDR(args.ServiceCIDR) ||
		args.HostPrefix > 0 {

		network = msv1.NewNetwork()

		if !IsEmptyCIDR(args.MachineCIDR) {
			network = network.MachineCIDR(args.MachineCIDR.String())
		}
		if !IsEmptyCIDR(args.ServiceCIDR) {
			network = network.ServiceCIDR(args.ServiceCIDR.String())
		}
		if !IsEmptyCIDR(args.PodCIDR) {
			network = network.PodCIDR(args.PodCIDR.String())
		}
		if args.HostPrefix > 0 {
			network = network.HostPrefix(args.HostPrefix)
		}
	}

	service, err := msv1.NewManagedService().
		Service(args.ServiceType).
		Parameters(parameters...).
		Cluster(
			msv1.NewCluster().
				Name(args.ClusterName).
				Properties(args.Properties).
				Region(
					msv1.NewCloudRegion().
						ID(args.AwsRegion)).
				MultiAZ(args.MultiAZ).
				Network(network).
				API(msv1.NewClusterAPI().
					Listening(listeningMethod)).
				AWS(
					msv1.NewAWS().
						STS(msv1.NewSTS().
							RoleARN(args.AwsRoleARN).
							SupportRoleARN(args.AwsSupportRoleARN).
							InstanceIAMRoles(msv1.NewInstanceIAMRoles().
								MasterRoleARN(args.AwsControlPlaneRoleARN).
								WorkerRoleARN(args.AwsWorkerRoleARN)).
							OperatorIAMRoles(operatorIamRoles...)).
						AccountID(args.AwsAccountID).
						SubnetIDs(args.SubnetIDs...).
						PrivateLink(args.Privatelink)).
				Nodes(msv1.NewClusterNodes().
					AvailabilityZones(args.AvailabilityZones...))).
		Build()

	if err != nil {
		return nil, fmt.Errorf("failed to create Managed Service call: %w", err)
	}

	serviceCall, err := c.ocm.ServiceMgmt().V1().Services().
		Add().
		Body(service).
		Send()
	if err != nil {
		return nil, handleErr(serviceCall.Error(), err)
	}

	return serviceCall.Body(), nil
}

func (c *Client) ListManagedServices(count int) (*msv1.ManagedServiceList, error) {
	if fedramp.Enabled() {
		return nil, fedrampError
	}

	if count < 0 {
		return nil, fmt.Errorf("invalid services count")
	}

	response, err := c.ocm.ServiceMgmt().V1().Services().List().Send()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve services list: %w", err)
	}
	return response.Items(), nil

}

type DescribeManagedServiceArgs struct {
	ID string
}

func (c *Client) GetManagedService(args DescribeManagedServiceArgs) (*msv1.ManagedService, error) {
	if fedramp.Enabled() {
		return nil, fedrampError
	}

	response, err := c.ocm.ServiceMgmt().V1().Services().Service(args.ID).Get().Send()
	if err != nil {
		return nil, fmt.Errorf("failed to get managed service with id %s: %w", args.ID, err)
	}
	return response.Body(), nil
}

type DeleteManagedServiceArgs struct {
	ID string
}

func (c *Client) DeleteManagedService(args DeleteManagedServiceArgs) (*msv1.ManagedServiceDeleteResponse, error) {
	if fedramp.Enabled() {
		return nil, fedrampError
	}

	deleteResponse, err := c.ocm.ServiceMgmt().V1().Services().
		Service(args.ID).
		Delete().
		Send()
	if err != nil {
		return nil, handleErr(deleteResponse.Error(), err)
	}

	return deleteResponse, nil
}

type UpdateManagedServiceArgs struct {
	ID         string
	Parameters map[string]string
}

func (c *Client) UpdateManagedService(args UpdateManagedServiceArgs) error {
	if fedramp.Enabled() {
		return fedrampError
	}

	parameters := []*msv1.ServiceParameterBuilder{}
	for id, val := range args.Parameters {
		parameters = append(parameters,
			msv1.NewServiceParameter().ID(id).Value(val))
	}

	serviceSpec, err := msv1.NewManagedService().
		Parameters(parameters...).
		Build()

	if err != nil {
		return fmt.Errorf("failed to create Managed Service call: %w", err)
	}

	serviceCall, err := c.ocm.ServiceMgmt().V1().Services().Service(args.ID).
		Update().
		Body(serviceSpec).
		Send()
	if err != nil {
		return handleErr(serviceCall.Error(), err)
	}

	return nil
}
