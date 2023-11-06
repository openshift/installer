package aws

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/sets"
)

type mockEC2Client struct {
	ec2iface.EC2API

	// swappable functions
	createVpc          func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error)
	describeVpcs       func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	modifyVpcAttribute func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error)

	createIgw    func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error)
	describeIgws func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error)
	attachIgw    func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error)

	createDHCPOpts   func(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error)
	describeDHCPOpts func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error)
	assocOpts        func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error)

	createRouteTable       func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error)
	describeRouteTables    func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error)
	replaceRouteTableAssoc func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error)
	assocRouteTable        func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error)
	createRoute            func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error)

	createSubnet    func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error)
	describeSubnets func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)

	allocAddr    func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error)
	createTags   func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)
	disassocAddr func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error)
	releaseAddr  func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error)

	createNat    func(*ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error)
	describeNats func(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error)

	createVpcEndpoint    func(*ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error)
	describeVpcEndpoints func(*ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error)

	describeSGs      func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	createSG         func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
	authorizeEgress  func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error)
	authorizeIngress func(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)

	describeInstances func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	runInstances      func(*ec2.RunInstancesInput) (*ec2.Reservation, error)
	getDefaultKmsKey  func(*ec2.GetEbsDefaultKmsKeyIdInput) (*ec2.GetEbsDefaultKmsKeyIdOutput, error)
}

func (m *mockEC2Client) DescribeVpcsWithContext(_ context.Context, in *ec2.DescribeVpcsInput, _ ...request.Option) (*ec2.DescribeVpcsOutput, error) {
	return m.describeVpcs(in)
}

func (m *mockEC2Client) CreateVpcWithContext(_ context.Context, in *ec2.CreateVpcInput, _ ...request.Option) (*ec2.CreateVpcOutput, error) {
	return m.createVpc(in)
}

func (m *mockEC2Client) ModifyVpcAttributeWithContext(_ context.Context, in *ec2.ModifyVpcAttributeInput, _ ...request.Option) (*ec2.ModifyVpcAttributeOutput, error) {
	return m.modifyVpcAttribute(in)
}

func (m *mockEC2Client) DescribeInternetGatewaysWithContext(_ context.Context, in *ec2.DescribeInternetGatewaysInput, _ ...request.Option) (*ec2.DescribeInternetGatewaysOutput, error) {
	return m.describeIgws(in)
}

func (m *mockEC2Client) CreateInternetGatewayWithContext(_ context.Context, in *ec2.CreateInternetGatewayInput, _ ...request.Option) (*ec2.CreateInternetGatewayOutput, error) {
	return m.createIgw(in)
}

func (m *mockEC2Client) AttachInternetGatewayWithContext(_ context.Context, in *ec2.AttachInternetGatewayInput, _ ...request.Option) (*ec2.AttachInternetGatewayOutput, error) {
	return m.attachIgw(in)
}

func (m *mockEC2Client) CreateDhcpOptionsWithContext(_ context.Context, in *ec2.CreateDhcpOptionsInput, _ ...request.Option) (*ec2.CreateDhcpOptionsOutput, error) {
	return m.createDHCPOpts(in)
}

func (m *mockEC2Client) DescribeDhcpOptionsWithContext(_ context.Context, in *ec2.DescribeDhcpOptionsInput, _ ...request.Option) (*ec2.DescribeDhcpOptionsOutput, error) {
	return m.describeDHCPOpts(in)
}

func (m *mockEC2Client) AssociateDhcpOptionsWithContext(_ context.Context, in *ec2.AssociateDhcpOptionsInput, _ ...request.Option) (*ec2.AssociateDhcpOptionsOutput, error) {
	return m.assocOpts(in)
}

func (m *mockEC2Client) CreateRouteTableWithContext(_ context.Context, in *ec2.CreateRouteTableInput, _ ...request.Option) (*ec2.CreateRouteTableOutput, error) {
	return m.createRouteTable(in)
}

func (m *mockEC2Client) DescribeRouteTablesWithContext(_ context.Context, in *ec2.DescribeRouteTablesInput, _ ...request.Option) (*ec2.DescribeRouteTablesOutput, error) {
	return m.describeRouteTables(in)
}

func (m *mockEC2Client) ReplaceRouteTableAssociationWithContext(_ context.Context, in *ec2.ReplaceRouteTableAssociationInput, _ ...request.Option) (*ec2.ReplaceRouteTableAssociationOutput, error) {
	return m.replaceRouteTableAssoc(in)
}

func (m *mockEC2Client) AssociateRouteTableWithContext(_ context.Context, in *ec2.AssociateRouteTableInput, _ ...request.Option) (*ec2.AssociateRouteTableOutput, error) {
	return m.assocRouteTable(in)
}

func (m *mockEC2Client) CreateRouteWithContext(_ context.Context, in *ec2.CreateRouteInput, _ ...request.Option) (*ec2.CreateRouteOutput, error) {
	return m.createRoute(in)
}

func (m *mockEC2Client) CreateSubnetWithContext(_ context.Context, in *ec2.CreateSubnetInput, _ ...request.Option) (*ec2.CreateSubnetOutput, error) {
	return m.createSubnet(in)
}

func (m *mockEC2Client) DescribeSubnetsWithContext(_ context.Context, in *ec2.DescribeSubnetsInput, _ ...request.Option) (*ec2.DescribeSubnetsOutput, error) {
	return m.describeSubnets(in)
}

func (m *mockEC2Client) AllocateAddressWithContext(_ context.Context, in *ec2.AllocateAddressInput, _ ...request.Option) (*ec2.AllocateAddressOutput, error) {
	return m.allocAddr(in)
}

func (m *mockEC2Client) DisassociateAddressWithContext(_ context.Context, in *ec2.DisassociateAddressInput, _ ...request.Option) (*ec2.DisassociateAddressOutput, error) {
	return m.disassocAddr(in)
}

func (m *mockEC2Client) ReleaseAddressWithContext(_ context.Context, in *ec2.ReleaseAddressInput, _ ...request.Option) (*ec2.ReleaseAddressOutput, error) {
	return m.releaseAddr(in)
}

func (m *mockEC2Client) CreateTagsWithContext(_ context.Context, in *ec2.CreateTagsInput, _ ...request.Option) (*ec2.CreateTagsOutput, error) {
	return m.createTags(in)
}

func (m *mockEC2Client) CreateNatGatewayWithContext(_ context.Context, in *ec2.CreateNatGatewayInput, _ ...request.Option) (*ec2.CreateNatGatewayOutput, error) {
	return m.createNat(in)
}

func (m *mockEC2Client) DescribeNatGatewaysWithContext(_ context.Context, in *ec2.DescribeNatGatewaysInput, _ ...request.Option) (*ec2.DescribeNatGatewaysOutput, error) {
	return m.describeNats(in)
}

func (m *mockEC2Client) CreateVpcEndpointWithContext(_ context.Context, in *ec2.CreateVpcEndpointInput, _ ...request.Option) (*ec2.CreateVpcEndpointOutput, error) {
	return m.createVpcEndpoint(in)
}

func (m *mockEC2Client) DescribeVpcEndpointsWithContext(_ context.Context, in *ec2.DescribeVpcEndpointsInput, _ ...request.Option) (*ec2.DescribeVpcEndpointsOutput, error) {
	return m.describeVpcEndpoints(in)
}

func (m *mockEC2Client) DescribeSecurityGroupsWithContext(_ context.Context, in *ec2.DescribeSecurityGroupsInput, _ ...request.Option) (*ec2.DescribeSecurityGroupsOutput, error) {
	return m.describeSGs(in)
}

func (m *mockEC2Client) CreateSecurityGroupWithContext(_ context.Context, in *ec2.CreateSecurityGroupInput, _ ...request.Option) (*ec2.CreateSecurityGroupOutput, error) {
	return m.createSG(in)
}

func (m *mockEC2Client) AuthorizeSecurityGroupEgressWithContext(_ context.Context, in *ec2.AuthorizeSecurityGroupEgressInput, _ ...request.Option) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	return m.authorizeEgress(in)
}

func (m *mockEC2Client) AuthorizeSecurityGroupIngressWithContext(_ context.Context, in *ec2.AuthorizeSecurityGroupIngressInput, _ ...request.Option) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return m.authorizeIngress(in)
}

func (m *mockEC2Client) DescribeInstancesWithContext(_ context.Context, in *ec2.DescribeInstancesInput, _ ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	return m.describeInstances(in)
}

func (m *mockEC2Client) RunInstancesWithContext(_ context.Context, in *ec2.RunInstancesInput, _ ...request.Option) (*ec2.Reservation, error) {
	return m.runInstances(in)
}

func (m *mockEC2Client) GetEbsDefaultKmsKeyIdWithContext(_ context.Context, in *ec2.GetEbsDefaultKmsKeyIdInput, _ ...request.Option) (*ec2.GetEbsDefaultKmsKeyIdOutput, error) { //nolint:revive,stylecheck //This is a mocked function so we cannot rename it
	return m.getDefaultKmsKey(in)
}

var errAwsSdk = errors.New("some AWS SDK error")

func TestEnsureVPC(t *testing.T) {
	vpcID := aws.String("vpc-1")
	expectedVpc := &ec2.Vpc{VpcId: vpcID}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.Vpc
		expectedErr string
	}{
		{
			name: "SDK error fetching VPCs",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return nil, errAwsSdk
				},
				createVpc: func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					panic("should not be called")
				},
				modifyVpcAttribute: func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list VPCs: some AWS SDK error$`,
		},
		{
			name: "VPC found but enabling DNS support fails",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return &ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{
							{VpcId: vpcID},
							{VpcId: aws.String("vpc-2")},
						},
					}, nil
				},
				createVpc: func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					panic("should not be called")
				},
				modifyVpcAttribute: func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to enable DNS support on VPC: some AWS SDK error$`,
		},
		{
			name: "VPC found but enabling DNS hostnames fails",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return &ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{{VpcId: vpcID}},
					}, nil
				},
				createVpc: func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					panic("should not be called")
				},
				modifyVpcAttribute: func(in *ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					if in.EnableDnsHostnames != nil {
						return nil, errAwsSdk
					}
					return &ec2.ModifyVpcAttributeOutput{}, nil
				},
			},
			expectedErr: `^failed to enable DNS hostnames on VPC: some AWS SDK error$`,
		},
		{
			name: "VPC found and attributes set",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return &ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{{VpcId: vpcID}},
					}, nil
				},
				createVpc: func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					panic("should not be called")
				},
				modifyVpcAttribute: func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					return &ec2.ModifyVpcAttributeOutput{}, nil
				},
			},
			expectedOut: expectedVpc,
		},
		{
			name: "VPC created",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return &ec2.DescribeVpcsOutput{Vpcs: []*ec2.Vpc{}}, nil
				},
				createVpc: func(in *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("vpc not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "vpc" {
						panic("vpc tagged with wrong resource type")
					}
					return &ec2.CreateVpcOutput{Vpc: &ec2.Vpc{VpcId: vpcID}}, nil
				},
				modifyVpcAttribute: func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					return &ec2.ModifyVpcAttributeOutput{}, nil
				},
			},
			expectedOut: expectedVpc,
		},
		{
			name: "VPC creation fails",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return &ec2.DescribeVpcsOutput{Vpcs: []*ec2.Vpc{}}, nil
				},
				createVpc: func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					return nil, errAwsSdk
				},
				modifyVpcAttribute: func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create VPC: some AWS SDK error$`,
		},
		{
			name: "VPC created but enabling DNS support fails",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return &ec2.DescribeVpcsOutput{Vpcs: []*ec2.Vpc{}}, nil
				},
				createVpc: func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					return &ec2.CreateVpcOutput{Vpc: &ec2.Vpc{VpcId: vpcID}}, nil
				},
				modifyVpcAttribute: func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `failed to enable DNS support on VPC: some AWS SDK error$`,
		},
		{
			name: "VPC created but enabling DNS hostnames fails",
			mockSvc: mockEC2Client{
				describeVpcs: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
					return &ec2.DescribeVpcsOutput{Vpcs: []*ec2.Vpc{}}, nil
				},
				createVpc: func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
					return &ec2.CreateVpcOutput{Vpc: &ec2.Vpc{VpcId: vpcID}}, nil
				},
				modifyVpcAttribute: func(in *ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
					if in.EnableDnsHostnames != nil {
						return nil, errAwsSdk
					}
					return &ec2.ModifyVpcAttributeOutput{}, nil
				},
			},
			expectedErr: `failed to enable DNS hostnames on VPC: some AWS SDK error$`,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
			}
			res, err := state.ensureVPC(context.TODO(), logger, &test.mockSvc)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureIntergetGateway(t *testing.T) {
	vpcID := aws.String("vpc-1")
	igwID := aws.String("igw-1")
	attachedIgw := &ec2.InternetGateway{
		InternetGatewayId: igwID,
		Attachments:       []*ec2.InternetGatewayAttachment{{VpcId: vpcID}},
	}
	unattachedIgw := &ec2.InternetGateway{InternetGatewayId: igwID}
	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.InternetGateway
		expectedErr string
	}{
		{
			name: "SDK error fetching Internet Gateways",
			mockSvc: mockEC2Client{
				describeIgws: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
					return nil, errAwsSdk
				},
				createIgw: func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					panic("should not be called")
				},
				attachIgw: func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list internet gateways: some AWS SDK error$`,
		},
		{
			name: "Internet Gateway found but attaching to VPC fails",
			mockSvc: mockEC2Client{
				describeIgws: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
					return &ec2.DescribeInternetGatewaysOutput{
						InternetGateways: []*ec2.InternetGateway{unattachedIgw},
					}, nil
				},
				createIgw: func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					panic("should not be called")
				},
				attachIgw: func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to attach internet gateway to VPC: some AWS SDK error$`,
		},
		{
			name: "Internet Gateway found and already attached to VPC",
			mockSvc: mockEC2Client{
				describeIgws: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
					return &ec2.DescribeInternetGatewaysOutput{
						InternetGateways: []*ec2.InternetGateway{attachedIgw},
					}, nil
				},
				createIgw: func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					panic("should not be called")
				},
				attachIgw: func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					panic("should not be called")
				},
			},
			expectedOut: attachedIgw,
		},
		{
			name: "Internet Gateway found and attaching to VPC succeeds",
			mockSvc: mockEC2Client{
				describeIgws: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
					return &ec2.DescribeInternetGatewaysOutput{
						InternetGateways: []*ec2.InternetGateway{unattachedIgw},
					}, nil
				},
				createIgw: func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					panic("should not be called")
				},
				attachIgw: func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					return &ec2.AttachInternetGatewayOutput{}, nil
				},
			},
			expectedOut: unattachedIgw,
		},
		{
			name: "Internet Gateway creation fails",
			mockSvc: mockEC2Client{
				describeIgws: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
					return &ec2.DescribeInternetGatewaysOutput{}, nil
				},
				createIgw: func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					return nil, errAwsSdk
				},
				attachIgw: func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create internet gateway: some AWS SDK error$`,
		},
		{
			name: "Internet Gateway created but attaching to VPC fails",
			mockSvc: mockEC2Client{
				describeIgws: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
					return &ec2.DescribeInternetGatewaysOutput{}, nil
				},
				createIgw: func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					return &ec2.CreateInternetGatewayOutput{InternetGateway: unattachedIgw}, nil
				},
				attachIgw: func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to attach internet gateway to VPC: some AWS SDK error$`,
		},
		{
			name: "Internet Gateway created and attaching to VPC succeeds",
			mockSvc: mockEC2Client{
				describeIgws: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
					return &ec2.DescribeInternetGatewaysOutput{}, nil
				},
				createIgw: func(in *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("internet gateway not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "internet-gateway" {
						panic("internet gateway tagged with wrong resource type")
					}
					return &ec2.CreateInternetGatewayOutput{InternetGateway: unattachedIgw}, nil
				},
				attachIgw: func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					return &ec2.AttachInternetGatewayOutput{}, nil
				},
			},
			expectedOut: unattachedIgw,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
			}
			res, err := state.ensureInternetGateway(context.TODO(), logger, &test.mockSvc, vpcID)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureDHCPOptions(t *testing.T) {
	vpcID := aws.String("vpc-1")
	expectedOpt := &ec2.DhcpOptions{
		DhcpOptionsId: aws.String("dhcp-opt-1"),
	}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.DhcpOptions
		expectedErr string
	}{
		{
			name: "SDK error fetching DHCP options",
			mockSvc: mockEC2Client{
				describeDHCPOpts: func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
					return nil, errAwsSdk
				},
				createDHCPOpts: func(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
					panic("should not be called")
				},
				assocOpts: func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list DHCP options: some AWS SDK error$`,
		},
		{
			name: "DHCP options found but association with VPC fails",
			mockSvc: mockEC2Client{
				describeDHCPOpts: func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
					return &ec2.DescribeDhcpOptionsOutput{
						DhcpOptions: []*ec2.DhcpOptions{expectedOpt},
					}, nil
				},
				createDHCPOpts: func(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
					panic("should not be called")
				},
				assocOpts: func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to associate DHCP options to VPC: some AWS SDK error$`,
		},
		{
			name: "DHCP options found and association with VPC succeeds",
			mockSvc: mockEC2Client{
				describeDHCPOpts: func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
					return &ec2.DescribeDhcpOptionsOutput{
						DhcpOptions: []*ec2.DhcpOptions{expectedOpt},
					}, nil
				},
				createDHCPOpts: func(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
					panic("should not be called")
				},
				assocOpts: func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
					return &ec2.AssociateDhcpOptionsOutput{}, nil
				},
			},
			expectedOut: expectedOpt,
		},
		{
			name: "DHCP creation fails",
			mockSvc: mockEC2Client{
				describeDHCPOpts: func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
					return &ec2.DescribeDhcpOptionsOutput{
						DhcpOptions: []*ec2.DhcpOptions{},
					}, nil
				},
				createDHCPOpts: func(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
					return nil, errAwsSdk
				},
				assocOpts: func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create DHCP options: some AWS SDK error$`,
		},
		{
			name: "DHCP created but association with VPC fails",
			mockSvc: mockEC2Client{
				describeDHCPOpts: func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
					return &ec2.DescribeDhcpOptionsOutput{
						DhcpOptions: []*ec2.DhcpOptions{},
					}, nil
				},
				createDHCPOpts: func(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
					return &ec2.CreateDhcpOptionsOutput{DhcpOptions: expectedOpt}, nil
				},
				assocOpts: func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to associate DHCP options to VPC: some AWS SDK error$`,
		},
		{
			name: "DHCP created and association with VPC succeeds",
			mockSvc: mockEC2Client{
				describeDHCPOpts: func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
					return &ec2.DescribeDhcpOptionsOutput{
						DhcpOptions: []*ec2.DhcpOptions{},
					}, nil
				},
				createDHCPOpts: func(in *ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("DHCP options not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "dhcp-options" {
						panic("DHCP options tagged with wrong resource type")
					}
					return &ec2.CreateDhcpOptionsOutput{DhcpOptions: expectedOpt}, nil
				},
				assocOpts: func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
					return &ec2.AssociateDhcpOptionsOutput{}, nil
				},
			},
			expectedOut: expectedOpt,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
			}
			res, err := state.ensureDHCPOptions(context.TODO(), logger, &test.mockSvc, vpcID)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureRouteTable(t *testing.T) {
	vpcID := aws.String("vpc-1")
	tableID := aws.String("route-table-1")
	tableName := "route-table"

	expectedTable := &ec2.RouteTable{
		RouteTableId: tableID,
	}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.RouteTable
		expectedErr string
	}{
		{
			name: "SDK error listing route tables",
			mockSvc: mockEC2Client{
				describeRouteTables: func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return nil, errAwsSdk
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `failed to list route tables: some AWS SDK error$`,
		},
		{
			name: "Route table found",
			mockSvc: mockEC2Client{
				describeRouteTables: func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{RouteTables: []*ec2.RouteTable{expectedTable}}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
			},
			expectedOut: expectedTable,
		},
		{
			name: "Route table creation fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{RouteTables: []*ec2.RouteTable{}}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^some AWS SDK error$`,
		},
		{
			name: "Route table created but not found",
			mockSvc: mockEC2Client{
				describeRouteTables: func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{RouteTables: []*ec2.RouteTable{}}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					return &ec2.CreateRouteTableOutput{RouteTable: expectedTable}, nil
				},
			},
			expectedErr: `^failed to find route table \(route\-table\-1\) that was just created: .*$`,
		},
		{
			name: "Route table creation succeeds",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeRouteTablesOutput{RouteTables: []*ec2.RouteTable{}}, nil
					}
					return &ec2.DescribeRouteTablesOutput{RouteTables: []*ec2.RouteTable{expectedTable}}, nil
				},
				createRouteTable: func(in *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("route table not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "route-table" {
						panic("route table tagged with wrong resource type")
					}
					return &ec2.CreateRouteTableOutput{RouteTable: expectedTable}, nil
				},
			},
			expectedOut: expectedTable,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
			}
			res, err := state.ensureRouteTable(context.TODO(), logger, &test.mockSvc, tableName)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsurePublicRouteTable(t *testing.T) {
	vpcID := aws.String("vpc-1")
	igwID := aws.String("igw-1")
	tableID := aws.String("route-table-1")
	mainTableID := aws.String("route-table-main")

	const vpcIDFilter = "vpc-id"

	nonMainTable := &ec2.RouteTable{
		RouteTableId: tableID,
	}
	mainTable := &ec2.RouteTable{
		RouteTableId: tableID,
		Associations: []*ec2.RouteTableAssociation{
			{
				RouteTableAssociationId: aws.String("assoc-1"),
				RouteTableId:            aws.String("route-table-2"),
			},
		},
	}
	mainTableNoAssoc := &ec2.RouteTable{
		RouteTableId: mainTableID,
		Associations: []*ec2.RouteTableAssociation{},
	}
	mainTableNotVPC := &ec2.RouteTable{
		RouteTableId: mainTableID,
		Associations: []*ec2.RouteTableAssociation{
			{
				RouteTableAssociationId: aws.String("assoc-1"),
				RouteTableId:            aws.String("route-table-2"),
			},
		},
	}
	mainTableIgw := &ec2.RouteTable{
		RouteTableId: tableID,
		Associations: []*ec2.RouteTableAssociation{
			{
				RouteTableAssociationId: aws.String("assoc-1"),
				RouteTableId:            aws.String("route-table-2"),
			},
		},
		Routes: []*ec2.Route{
			{
				GatewayId:            igwID,
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
			},
		},
	}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.RouteTable
		expectedErr string
	}{
		{
			name: "Route table creation fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return nil, errAwsSdk
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create public route table: failed to list route tables: some AWS SDK error$`,
		},
		{
			name: "Route table created but VPC main route table fetching fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					if aws.StringValue(in.Filters[0].Name) == vpcIDFilter {
						return nil, errAwsSdk
					}
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{nonMainTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to get main route table: failed to list route tables: some AWS SDK error$`,
		},
		{
			name: "Route table created but VPC main route table not found",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					// main route table filter
					if aws.StringValue(in.Filters[0].Name) == vpcIDFilter {
						return &ec2.DescribeRouteTablesOutput{
							RouteTables: []*ec2.RouteTable{},
						}, nil
					}
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{nonMainTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^no main route table associated with the VPC$`,
		},
		{
			name: "Route table created but VPC main route table has no associations",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					// main route table filter
					if aws.StringValue(in.Filters[0].Name) == vpcIDFilter {
						return &ec2.DescribeRouteTablesOutput{
							RouteTables: []*ec2.RouteTable{mainTableNoAssoc},
						}, nil
					}
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{nonMainTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^no associations found for main route table$`,
		},
		{
			name: "Route table created but replacing main route table association fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					// main route table filter
					if aws.StringValue(in.Filters[0].Name) == vpcIDFilter {
						return &ec2.DescribeRouteTablesOutput{
							RouteTables: []*ec2.RouteTable{mainTableNotVPC},
						}, nil
					}
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{nonMainTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					return nil, errAwsSdk
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to replace vpc main route table: some AWS SDK error$`,
		},
		{
			name: "Route table created and already main route table but Internet Gateway route fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{mainTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create route to internet gateway: some AWS SDK error$`,
		},
		{
			name: "Route table created and already main route table and Internet Gateway route already exists",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{mainTableIgw},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
			},
			expectedOut: mainTableIgw,
		},
		{
			name: "Route table created and already main route table and Internet Gateway route created",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{mainTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				replaceRouteTableAssoc: func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					return &ec2.CreateRouteOutput{}, nil
				},
			},
			expectedOut: mainTable,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
				igwID: igwID,
			}
			res, err := state.ensurePublicRouteTable(context.TODO(), logger, &test.mockSvc)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureSubnet(t *testing.T) {
	vpcID := aws.String("vpc-1")
	subnetName := "subnet-zone1"

	expectedSubnet := &ec2.Subnet{SubnetId: aws.String("subnet-1")}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.Subnet
		expectedErr string
	}{
		{
			name: "SDK error listing subnets",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return nil, errAwsSdk
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list subnets: some AWS SDK error$`,
		},
		{
			name: "Subnet creation fails",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return &ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{}}, nil
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create subnet: some AWS SDK error$`,
		},
		{
			name: "Subnet created but not found",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{}}, nil
					}
					return &ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{}}, nil
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					return &ec2.CreateSubnetOutput{Subnet: expectedSubnet}, nil
				},
			},
			expectedErr: `^failed to find subnet \(subnet-1\) that was just created: .*$`,
		},
		{
			name: "Subnet created and found",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{}}, nil
					}
					return &ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{expectedSubnet}}, nil
				},
				createSubnet: func(in *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("subnet not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "subnet" {
						panic("subnet tagged with wrong resource type")
					}
					return &ec2.CreateSubnetOutput{Subnet: expectedSubnet}, nil
				},
			},
			expectedOut: expectedSubnet,
		},
		{
			name: "Existing subnet found",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return &ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{expectedSubnet}}, nil
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					panic("should not be called")
				},
			},
			expectedOut: expectedSubnet,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
			}
			res, err := state.ensureSubnet(context.TODO(), logger, &test.mockSvc, "zone1", "10.0.1.0/17", subnetName, map[string]string{})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsurePublicSubnets(t *testing.T) {
	vpcID := aws.String("vpc-1")
	_, network, err := net.ParseCIDR("10.0.0.0/17")
	assert.NoError(t, err, "parsing CIDR should not fail")

	subnet1 := aws.String("subnet-zone1")
	subnet2 := aws.String("subnet-zone2")
	subnet3 := aws.String("subnet-zone3")
	expectSubnets := []*string{subnet1, subnet2, subnet3}
	expectedMap := map[string]*string{
		"zone1": subnet1,
		"zone2": subnet2,
		"zone3": subnet3,
	}

	tests := []struct {
		name            string
		mockSvc         mockEC2Client
		expectedSubnets []*string
		expectedMap     map[string]*string
		expectedErr     string
	}{
		{
			name: "AWS SDK error listing subnets",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create public subnets: failed to create public subnet \(infraID-public-zone1\): failed to list subnets: some AWS SDK error$`,
		},
		{
			name: "AWS SDK error creating a subnet",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{}, nil
					}
					panic("should not be reached")
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create public subnets: failed to create public subnet \(infraID-public-zone1\): failed to create subnet: some AWS SDK error$`,
		},
		{
			name: "Subnet created but failed to be fetched",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{}, nil
					}
					return &ec2.DescribeSubnetsOutput{}, nil
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					return &ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{SubnetId: subnet1},
					}, nil
				},
			},
			expectedErr: `^failed to create public subnets: failed to create public subnet \(infraID-public-zone1\): failed to find subnet \(subnet-zone1\) that was just created: .*$`,
		},
		{
			name: "Public subnets created but associating with public route table fails",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{}, nil
					}
					return &ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{SubnetId: in.SubnetIds[0]}},
					}, nil
				},
				createSubnet: func(in *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					publicELBTagFound := false
					for _, ec2Tag := range in.TagSpecifications[0].Tags {
						if aws.StringValue(ec2Tag.Key) == "kubernetes.io/role/elb" && aws.StringValue(ec2Tag.Value) == "true" {
							publicELBTagFound = true
							break
						}
					}
					if !publicELBTagFound {
						panic("public subnet not tagged with public ELB")
					}
					return &ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{SubnetId: expectedMap[*in.AvailabilityZone]},
					}, nil
				},
				assocRouteTable: func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to associate public subnet \(subnet-zone1\) with public route table \(table-1\): some AWS SDK error$`,
		},
		{
			name: "Public subnets created",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{}, nil
					}
					return &ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{SubnetId: in.SubnetIds[0]}},
					}, nil
				},
				createSubnet: func(in *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					publicELBTagFound := false
					for _, ec2Tag := range in.TagSpecifications[0].Tags {
						if aws.StringValue(ec2Tag.Key) == "kubernetes.io/role/elb" && aws.StringValue(ec2Tag.Value) == "true" {
							publicELBTagFound = true
							break
						}
					}
					if !publicELBTagFound {
						panic("public subnet not tagged with public ELB")
					}
					return &ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{SubnetId: expectedMap[*in.AvailabilityZone]},
					}, nil
				},
				assocRouteTable: func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
					return &ec2.AssociateRouteTableOutput{}, nil
				},
			},
			expectedSubnets: expectSubnets,
			expectedMap:     expectedMap,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					zones:       []string{"zone1", "zone2", "zone3"},
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
			}
			publicRouteTable := &ec2.RouteTable{RouteTableId: aws.String("table-1")}
			subnets, subnetMap, err := state.ensurePublicSubnets(context.TODO(), logger, &test.mockSvc, network, publicRouteTable)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedSubnets, subnets)
				assert.Equal(t, test.expectedMap, subnetMap)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsurePrivateSubnets(t *testing.T) {
	vpcID := aws.String("vpc-1")
	_, network, err := net.ParseCIDR("10.0.0.0/17")
	assert.NoError(t, err, "parsing CIDR should not fail")

	subnet1 := aws.String("subnet-zone1")
	subnet2 := aws.String("subnet-zone2")
	subnet3 := aws.String("subnet-zone3")
	expectedSubnets := []*string{subnet1, subnet2, subnet3}
	expectedMap := map[string]*string{
		"zone1": subnet1,
		"zone2": subnet2,
		"zone3": subnet3,
	}

	tests := []struct {
		name            string
		mockSvc         mockEC2Client
		expectedSubnets []*string
		expectedMap     map[string]*string
		expectedErr     string
	}{
		{
			name: "AWS SDK error listing subnets",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create private subnets: failed to create private subnet \(infraID-private-zone1\): failed to list subnets: some AWS SDK error$`,
		},
		{
			name: "AWS SDK error creating a subnet",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{}, nil
					}
					panic("should not be reached")
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create private subnets: failed to create private subnet \(infraID-private-zone1\): failed to create subnet: some AWS SDK error$`,
		},
		{
			name: "Subnet created but failed to be fetched",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{}, nil
					}
					return &ec2.DescribeSubnetsOutput{}, nil
				},
				createSubnet: func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					return &ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{SubnetId: subnet1},
					}, nil
				},
			},
			expectedErr: `^failed to create private subnets: failed to create private subnet \(infraID-private-zone1\): failed to find subnet \(subnet-zone1\) that was just created: .*$`,
		},
		{
			name: "Private subnets created",
			mockSvc: mockEC2Client{
				describeSubnets: func(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSubnetsOutput{}, nil
					}
					return &ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{SubnetId: in.SubnetIds[0]}},
					}, nil
				},
				createSubnet: func(in *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
					privateELBTagFound := false
					for _, ec2Tag := range in.TagSpecifications[0].Tags {
						if aws.StringValue(ec2Tag.Key) == "kubernetes.io/role/internal-elb" && aws.StringValue(ec2Tag.Value) == "true" {
							privateELBTagFound = true
							break
						}
					}
					if !privateELBTagFound {
						panic("private subnet not tagged with private ELB")
					}
					return &ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{SubnetId: expectedMap[*in.AvailabilityZone]},
					}, nil
				},
			},
			expectedSubnets: expectedSubnets,
			expectedMap:     expectedMap,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					zones:       []string{"zone1", "zone2", "zone3"},
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
			}
			subnets, subnetMap, err := state.ensurePrivateSubnets(context.TODO(), logger, &test.mockSvc, network)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedSubnets, subnets)
				assert.Equal(t, test.expectedMap, subnetMap)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureEIP(t *testing.T) {
	expectedAlloc := aws.String("eip-1")
	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *string
		expectedErr string
	}{
		{
			name: "AWS SDK error allocating EIP address",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return nil, errAwsSdk
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to allocate EIP: some AWS SDK error$`,
		},
		{
			name: "EIP created but tagging fails with non-retriable error",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: expectedAlloc}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return nil, awserr.New(errInvalidEIPNotFound, "", errAwsSdk)
				},
				disassocAddr: func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
					return &ec2.DisassociateAddressOutput{}, nil
				},
				releaseAddr: func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
					return &ec2.ReleaseAddressOutput{}, nil
				},
			},
			expectedErr: `^failed to tag EIP \(eip-1\): .*$`,
		},
		{
			name: "EIP created but tagging fails with retriable error",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: expectedAlloc}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return nil, errAwsSdk
				},
				disassocAddr: func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
					return &ec2.DisassociateAddressOutput{}, nil
				},
				releaseAddr: func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
					return &ec2.ReleaseAddressOutput{}, nil
				},
			},
			expectedErr: `^failed to tag EIP \(eip-1\): some AWS SDK error$`,
		},
		{
			name: "EIP created, tagging fails and disassociation fails",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: expectedAlloc}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return nil, errAwsSdk
				},
				disassocAddr: func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
					return nil, errAwsSdk
				},
				releaseAddr: func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to release untagged EIP \(eip-1\): .*$`,
		},
		{
			name: "EIP created, tagging fails, address disassociated but releasing fails",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: expectedAlloc}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return nil, errAwsSdk
				},
				disassocAddr: func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
					return &ec2.DisassociateAddressOutput{}, nil
				},
				releaseAddr: func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to release untagged EIP \(eip-1\): .*$`,
		},
		{
			name: "EIP created, tagging fails, address disassociated and already released",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: expectedAlloc}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return nil, errAwsSdk
				},
				disassocAddr: func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
					return &ec2.DisassociateAddressOutput{}, nil
				},
				releaseAddr: func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
					return nil, awserr.New(errAuthFailure, "", errAwsSdk)
				},
			},
			expectedErr: `^failed to tag EIP \(eip-1\): some AWS SDK error$`,
		},
		{
			name: "EIP created and tagged",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: expectedAlloc}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return &ec2.CreateTagsOutput{}, nil
				},
				disassocAddr: func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
					panic("should not be called")
				},
				releaseAddr: func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
					panic("should not be called")
				},
			},
			expectedOut: expectedAlloc,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ensureEIP(context.TODO(), &test.mockSvc, map[string]string{})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureNatGateway(t *testing.T) {
	vpcID := aws.String("vpc-1")
	allocID := aws.String("eip-1")
	subnet1 := aws.String("subnet-zone1")
	expectedNat := &ec2.NatGateway{}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.NatGateway
		expectedErr string
	}{
		{
			name: "AWS SDK error listing Nat gateways",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					panic("should not be called")
				},
				describeNats: func(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
					return nil, errAwsSdk
				},
				createNat: func(*ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list Nat gateways: some AWS SDK error$`,
		},
		{
			name: "AWS SDK error creating EIP for Nat gateway",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return nil, errAwsSdk
				},
				describeNats: func(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
					return &ec2.DescribeNatGatewaysOutput{}, nil
				},
				createNat: func(*ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to allocate EIP: some AWS SDK error$`,
		},
		{
			name: "EIP created but non-retriable error creating Nat gateway",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: allocID}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return &ec2.CreateTagsOutput{}, nil
				},
				describeNats: func(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
					return &ec2.DescribeNatGatewaysOutput{}, nil
				},
				createNat: func(*ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create nat gateway: some AWS SDK error$`,
		},
		{
			name: "EIP created but retriable error creating Nat gateway",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: allocID}, nil
				},
				createTags: func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					return &ec2.CreateTagsOutput{}, nil
				},
				describeNats: func(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
					return &ec2.DescribeNatGatewaysOutput{}, nil
				},
				createNat: func(*ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
					return nil, awserr.New(errInvalidSubnet, "", errAwsSdk)
				},
			},
			expectedErr: `^failed to create nat gateway: .*$`,
		},
		{
			name: "EIP and Nat gateway created",
			mockSvc: mockEC2Client{
				allocAddr: func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: allocID}, nil
				},
				createTags: func(in *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
					if len(in.Tags) == 0 {
						panic("EIP not tagged")
					}
					return &ec2.CreateTagsOutput{}, nil
				},
				describeNats: func(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
					return &ec2.DescribeNatGatewaysOutput{}, nil
				},
				createNat: func(in *ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("nat gateway not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "natgateway" {
						panic("nat gateway tagged with wrong resource type")
					}
					return &ec2.CreateNatGatewayOutput{NatGateway: expectedNat}, nil
				},
			},
			expectedOut: expectedNat,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					zones:       []string{"zone1", "zone2", "zone3"},
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
			}
			res, err := state.ensureNatGateway(context.TODO(), logger, &test.mockSvc, subnet1, "zone1")
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsurePrivateRouteTable(t *testing.T) {
	vpcID := aws.String("vpc-1")
	tableID := aws.String("route-table-1")
	natGwID := aws.String("natgw-1")
	subnetID := aws.String("subnet-1")

	emptyTable := &ec2.RouteTable{
		RouteTableId: tableID,
	}
	tableNat := &ec2.RouteTable{
		RouteTableId: tableID,
		Routes: []*ec2.Route{
			{
				NatGatewayId:         natGwID,
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
			},
		},
	}
	expectedTable := &ec2.RouteTable{
		RouteTableId: tableID,
		Associations: []*ec2.RouteTableAssociation{
			{SubnetId: subnetID},
		},
		Routes: []*ec2.Route{
			{
				NatGatewayId:         natGwID,
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
			},
		},
	}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.RouteTable
		expectedErr string
	}{
		{
			name: "Route table creation fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return nil, errAwsSdk
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
				assocRouteTable: func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list route tables: some AWS SDK error$`,
		},
		{
			name: "Route table created but Nat Gateway route fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{emptyTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create route to nat gateway \(natgw-1\): some AWS SDK error$`,
		},
		{
			name: "Route table created and Nat Gateway route created but subnet association to route table fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{emptyTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					return &ec2.CreateRouteOutput{}, nil
				},
				assocRouteTable: func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to associate subnet \(subnet-1\) to route table \(route-table-1\): some AWS SDK error$`,
		},
		{
			name: "Route table created and Nat Gateway route already exists but subnet association to route table fails",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{tableNat},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					panic("should not be called")
				},
				assocRouteTable: func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to associate subnet \(subnet-1\) to route table \(route-table-1\): some AWS SDK error$`,
		},
		{
			name: "Route table created, Nat Gateway route created and subnet associated",
			mockSvc: mockEC2Client{
				describeRouteTables: func(in *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{expectedTable},
					}, nil
				},
				createRouteTable: func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					panic("should not be called")
				},
				createRoute: func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					return &ec2.CreateRouteOutput{}, nil
				},
				assocRouteTable: func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
					return &ec2.AssociateRouteTableOutput{}, nil
				},
			},
			expectedOut: expectedTable,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
			}
			res, err := state.ensurePrivateRouteTable(context.TODO(), logger, &test.mockSvc, natGwID, subnetID, "zone1")
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureS3VPCEndpoint(t *testing.T) {
	vpcID := aws.String("vpc-1")
	endpointID := aws.String("endpoint-1")

	expectedEndpoint := &ec2.VpcEndpoint{
		VpcEndpointId: endpointID,
	}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.VpcEndpoint
		expectedErr string
	}{
		{
			name: "AWS SDK error listing VPC endpoints",
			mockSvc: mockEC2Client{
				describeVpcEndpoints: func(*ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error) {
					return nil, errAwsSdk
				},
				createVpcEndpoint: func(*ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list VPC endpoints: some AWS SDK error$`,
		},
		{
			name: "VPC S3 endpoint creation fails with non-retriable error",
			mockSvc: mockEC2Client{
				describeVpcEndpoints: func(*ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error) {
					return &ec2.DescribeVpcEndpointsOutput{}, nil
				},
				createVpcEndpoint: func(*ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create VPC endpoint: some AWS SDK error$`,
		},
		{
			name: "VPC S3 endpoint creation fails with retriable error",
			mockSvc: mockEC2Client{
				describeVpcEndpoints: func(*ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error) {
					return &ec2.DescribeVpcEndpointsOutput{}, nil
				},
				createVpcEndpoint: func(*ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error) {
					return nil, awserr.New(errRouteTableIDNotFound, "", errAwsSdk)
				},
			},
			expectedErr: `^failed to create VPC endpoint: .*$`,
		},
		{
			name: "VPC S3 endpoint created",
			mockSvc: mockEC2Client{
				describeVpcEndpoints: func(*ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error) {
					return &ec2.DescribeVpcEndpointsOutput{}, nil
				},
				createVpcEndpoint: func(in *ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error) {
					if aws.StringValue(in.ServiceName) != "com.amazonaws.region1.s3" {
						panic("S3 VPC endpoint has wrong service name")
					}
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("S3 VPC endpoint not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "vpc-endpoint" {
						panic("S3 VPC endpoint tagged with wrong resource type")
					}
					return &ec2.CreateVpcEndpointOutput{VpcEndpoint: expectedEndpoint}, nil
				},
			},
			expectedOut: expectedEndpoint,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := vpcState{
				input: &vpcInputOptions{
					infraID:     "infraID",
					region:      "region1",
					cidrV4Block: "10.0.0.0/16",
					tags:        map[string]string{"custom-tag": "custom-value"},
				},
				vpcID: vpcID,
			}
			res, err := state.ensureVPCS3Endpoint(context.TODO(), logger, &test.mockSvc, []*string{})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureUserVpc(t *testing.T) {
	const vpcID = "vpc-1"
	publicIDs := []string{"public-1", "public-2"}
	privateIDs := []string{"private-1", "private-2"}
	privateSubnets := []*ec2.Subnet{
		{
			SubnetId:         aws.String("private-1"),
			AvailabilityZone: aws.String("zone-1"),
		},
		{
			SubnetId:         aws.String("private-2"),
			AvailabilityZone: aws.String("zone-2"),
		},
	}
	privateSubnetZoneMap := map[string]string{
		"zone-1": "private-1",
		"zone-2": "private-2",
	}
	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		publicIDs   []string
		expectedOut *vpcOutput
		expectedErr string
	}{
		{
			name: "AWS SDK error describing subnets",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to retrieve user-supplied subnets: some AWS SDK error$`,
		},
		{
			name: "Subnet zones mapped and no resources created in private-only VPC",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return &ec2.DescribeSubnetsOutput{
						Subnets: privateSubnets,
					}, nil
				},
			},
			expectedOut: &vpcOutput{
				vpcID:            vpcID,
				privateSubnetIDs: privateIDs,
				zoneToSubnetMap:  privateSubnetZoneMap,
				publicSubnetIDs:  nil,
			},
		},
		{
			name: "Subnet zones mapped and no resources created",
			mockSvc: mockEC2Client{
				describeSubnets: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
					return &ec2.DescribeSubnetsOutput{
						Subnets: privateSubnets,
					}, nil
				},
			},
			publicIDs: publicIDs,
			expectedOut: &vpcOutput{
				vpcID:            vpcID,
				privateSubnetIDs: privateIDs,
				zoneToSubnetMap:  privateSubnetZoneMap,
				publicSubnetIDs:  publicIDs,
			},
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := &vpcInputOptions{
				infraID:          "infraID",
				region:           "region",
				vpcID:            vpcID,
				privateSubnetIDs: privateIDs,
				tags:             map[string]string{},
			}
			input.publicSubnetIDs = test.publicIDs
			res, err := createVPCResources(context.TODO(), logger, &test.mockSvc, input)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureSecurityGroup(t *testing.T) {
	emptySG := &ec2.SecurityGroup{
		GroupId: aws.String("sg-1"),
	}
	expectedSG := &ec2.SecurityGroup{
		GroupId: aws.String("sg-1"),
		IpPermissionsEgress: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					{CidrIp: aws.String("0.0.0.0/0")},
				},
			},
		},
	}

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		expectedOut *ec2.SecurityGroup
		expectedErr string
	}{
		{
			name: "AWS SDK error listing security groups",
			mockSvc: mockEC2Client{
				describeSGs: func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					return nil, errAwsSdk
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					panic("should not be called")
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list security groups: some AWS SDK error$`,
		},
		{
			name: "Failed to create security group",
			mockSvc: mockEC2Client{
				describeSGs: func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					return nil, errNotFound
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return nil, errAwsSdk
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create security group \(node-sg\): some AWS SDK error$`,
		},
		{
			name: "Security group created but not found",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					return nil, errNotFound
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return &ec2.CreateSecurityGroupOutput{
						GroupId: aws.String("sg-1"),
					}, nil
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to find security group \(node-sg\) that was just created: .*$`,
		},
		{
			name: "Security group created but authorizing egress fails",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 0 {
						return nil, errNotFound
					}
					return &ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{emptySG},
					}, nil
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return &ec2.CreateSecurityGroupOutput{
						GroupId: aws.String("sg-1"),
					}, nil
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to authorize egress rules for security group \(node-sg\): .*$`,
		},
		{
			name: "Security group exists but authorizing egress fails",
			mockSvc: mockEC2Client{
				describeSGs: func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					return &ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{emptySG},
					}, nil
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					panic("should not be called")
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to authorize egress rules for security group \(node-sg\): .*$`,
		},
		{
			name: "Security Group created and egress rules authorized",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 0 {
						return nil, errNotFound
					}
					return &ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{expectedSG},
					}, nil
				},
				createSG: func(in *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("security group not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "security-group" {
						panic("security group tagged with wrong resource type")
					}
					return &ec2.CreateSecurityGroupOutput{
						GroupId: aws.String("sg-1"),
					}, nil
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					return &ec2.AuthorizeSecurityGroupEgressOutput{
						SecurityGroupRules: []*ec2.SecurityGroupRule{{}},
					}, nil
				},
			},
			expectedOut: expectedSG,
		},
		{
			name: "Security Group already created and egress rules already authorized",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 0 {
						return nil, errNotFound
					}
					return &ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{expectedSG},
					}, nil
				},
				createSG: func(in *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					if len(in.TagSpecifications) == 0 || len(in.TagSpecifications[0].Tags) == 0 {
						panic("security group not tagged")
					}
					if aws.StringValue(in.TagSpecifications[0].ResourceType) != "security-group" {
						panic("security group tagged with wrong resource type")
					}
					return &ec2.CreateSecurityGroupOutput{
						GroupId: aws.String("sg-1"),
					}, nil
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					return nil, awserr.New(errDuplicatePermission, "", errAwsSdk)
				},
			},
			expectedOut: expectedSG,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ensureSecurityGroup(context.TODO(), logger, &test.mockSvc, "infraID", "vpc-1", "node-sg", map[string]string{"custom-tag": "custom-value"})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

//nolint:gocyclo // we don't care about the cyclomatic complexity
func TestCreateSecurityGroups(t *testing.T) {
	bootstrapSG := &ec2.SecurityGroup{
		GroupId: aws.String("sg-1"),
	}
	controlplaneSG := &ec2.SecurityGroup{
		GroupId: aws.String("sg-2"),
	}
	computeSG := &ec2.SecurityGroup{
		GroupId: aws.String("sg-3"),
	}
	cidrBlocks := []string{"10.0.0.0/16"}
	defaultEgress := defaultEgressRules(aws.String("sg-1"))
	bootstrapIngress := defaultBootstrapSGIngressRules(aws.String("sg-1"), cidrBlocks)
	masterIngress := defaultMasterSGIngressRules(aws.String("sg-2"), aws.String("sg-3"), cidrBlocks)
	workerIngress := defaultWorkerSGIngressRules(aws.String("sg-3"), aws.String("sg-2"), cidrBlocks)

	tests := []struct {
		name        string
		mockSvc     mockEC2Client
		private     bool
		expectedOut *sgOutput
		expectedErr string
	}{
		{
			name: "Creating bootstrap security group fails when listing security groups",
			mockSvc: mockEC2Client{
				describeSGs: func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					return nil, errAwsSdk
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create bootstrap security group: failed to list security groups: some AWS SDK error$`,
		},
		{
			name: "Creating bootstrap security group fails",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: []*ec2.SecurityGroup{}}, nil
					}
					panic("should not be called")
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create bootstrap security group: failed to create security group \(infraID-bootstrap-sg\): some AWS SDK error$`,
		},
		{
			name: "Bootstrap security group created but authorizing ingress fails",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 0 {
						return &ec2.DescribeSecurityGroupsOutput{
							SecurityGroups: []*ec2.SecurityGroup{bootstrapSG},
						}, nil
					}
					panic("should not be called")
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					panic("should not be called")
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					return &ec2.AuthorizeSecurityGroupEgressOutput{
						SecurityGroupRules: make([]*ec2.SecurityGroupRule, len(defaultEgress)),
					}, nil
				},
				authorizeIngress: func(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to attach ingress rules to bootstrap security group: .*$`,
		},
		{
			name: "Bootstrap security group created and authorizing ingress for public cluster uses all addresses but creating master security group fails",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 1 && aws.StringValue(in.Filters[1].Values[0]) == "infraID-bootstrap-sg" {
						return &ec2.DescribeSecurityGroupsOutput{
							SecurityGroups: []*ec2.SecurityGroup{bootstrapSG},
						}, nil
					}
					return &ec2.DescribeSecurityGroupsOutput{}, nil
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return nil, errAwsSdk
				},
				authorizeEgress: func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					return &ec2.AuthorizeSecurityGroupEgressOutput{
						SecurityGroupRules: []*ec2.SecurityGroupRule{{}},
					}, nil
				},
				authorizeIngress: func(in *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					if aws.StringValue(in.IpPermissions[0].IpRanges[0].CidrIp) != "0.0.0.0/0" {
						panic("should have used 0.0.0.0/0 for public clusters")
					}
					rules := make([]*ec2.SecurityGroupRule, len(bootstrapIngress))
					return &ec2.AuthorizeSecurityGroupIngressOutput{
						SecurityGroupRules: rules,
					}, nil
				},
			},
			expectedErr: `^failed to create control plane security group: failed to create security group \(infraID-master-sg\): some AWS SDK error$`,
		},
		{
			name: "Bootstrap security group created, master created but creating worker fails",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 1 {
						switch aws.StringValue(in.Filters[1].Values[0]) {
						case "infraID-bootstrap-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{bootstrapSG},
							}, nil
						case "infraID-master-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{controlplaneSG},
							}, nil
						default:
							return &ec2.DescribeSecurityGroupsOutput{}, nil
						}
					}
					return &ec2.DescribeSecurityGroupsOutput{}, nil
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return nil, errAwsSdk
				},
				authorizeEgress: func(in *ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1", "sg-2":
						return &ec2.AuthorizeSecurityGroupEgressOutput{
							SecurityGroupRules: make([]*ec2.SecurityGroupRule, len(defaultEgress)),
						}, nil
					}
					panic("should not be called")
				},
				authorizeIngress: func(in *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1":
						return &ec2.AuthorizeSecurityGroupIngressOutput{
							SecurityGroupRules: make([]*ec2.SecurityGroupRule, len(bootstrapIngress)),
						}, nil
					default:
						panic("should not be called")
					}
				},
			},
			expectedErr: `failed to create compute security group: failed to create security group \(infraID-worker-sg\): some AWS SDK error$`,
		},
		{
			name: "Bootstrap, master and worker security groups created but authorizing master ingress fails",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 1 {
						switch aws.StringValue(in.Filters[1].Values[0]) {
						case "infraID-bootstrap-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{bootstrapSG},
							}, nil
						case "infraID-master-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{controlplaneSG},
							}, nil
						case "infraID-worker-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{computeSG},
							}, nil
						default:
							panic("should not be called")
						}
					}
					return &ec2.DescribeSecurityGroupsOutput{}, nil
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return nil, errAwsSdk
				},
				authorizeEgress: func(in *ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1", "sg-2", "sg-3":
						return &ec2.AuthorizeSecurityGroupEgressOutput{
							SecurityGroupRules: make([]*ec2.SecurityGroupRule, len(defaultEgress)),
						}, nil
					default:
						panic("should not be called")
					}
				},
				authorizeIngress: func(in *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1":
						rules := make([]*ec2.SecurityGroupRule, len(bootstrapIngress))
						return &ec2.AuthorizeSecurityGroupIngressOutput{
							SecurityGroupRules: rules,
						}, nil
					case "sg-2":
						return nil, errAwsSdk
					default:
						panic("should not be called")
					}
				},
			},
			expectedErr: `failed to attach ingress rules to master security group: .*$`,
		},
		{
			name: "Bootstrap, master and worker security groups created but authorizing worker ingress fails",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 1 {
						switch aws.StringValue(in.Filters[1].Values[0]) {
						case "infraID-bootstrap-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{bootstrapSG},
							}, nil
						case "infraID-master-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{controlplaneSG},
							}, nil
						case "infraID-worker-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{computeSG},
							}, nil
						default:
							panic("should not be called")
						}
					}
					return &ec2.DescribeSecurityGroupsOutput{}, nil
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return nil, errAwsSdk
				},
				authorizeEgress: func(in *ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1", "sg-2", "sg-3":
						return &ec2.AuthorizeSecurityGroupEgressOutput{
							SecurityGroupRules: make([]*ec2.SecurityGroupRule, len(defaultEgress)),
						}, nil
					default:
						panic("should not be called")
					}
				},
				authorizeIngress: func(in *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1":
						rules := make([]*ec2.SecurityGroupRule, len(bootstrapIngress))
						return &ec2.AuthorizeSecurityGroupIngressOutput{
							SecurityGroupRules: rules,
						}, nil
					case "sg-2":
						rules := make([]*ec2.SecurityGroupRule, len(masterIngress))
						return &ec2.AuthorizeSecurityGroupIngressOutput{
							SecurityGroupRules: rules,
						}, nil
					case "sg-3":
						return nil, errAwsSdk
					default:
						panic("should not be called")
					}
				},
			},
			expectedErr: `failed to attach ingress rules to worker security group: .*$`,
		},
		{
			name: "All security groups created and rules authorized",
			mockSvc: mockEC2Client{
				describeSGs: func(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
					if len(in.Filters) > 1 {
						switch aws.StringValue(in.Filters[1].Values[0]) {
						case "infraID-bootstrap-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{bootstrapSG},
							}, nil
						case "infraID-master-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{controlplaneSG},
							}, nil
						case "infraID-worker-sg":
							return &ec2.DescribeSecurityGroupsOutput{
								SecurityGroups: []*ec2.SecurityGroup{computeSG},
							}, nil
						default:
							panic("should not be called")
						}
					}
					return &ec2.DescribeSecurityGroupsOutput{}, nil
				},
				createSG: func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
					return nil, errAwsSdk
				},
				authorizeEgress: func(in *ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1", "sg-2", "sg-3":
						return &ec2.AuthorizeSecurityGroupEgressOutput{
							SecurityGroupRules: make([]*ec2.SecurityGroupRule, len(defaultEgress)),
						}, nil
					default:
						panic("should not be called")
					}
				},
				authorizeIngress: func(in *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					switch aws.StringValue(in.GroupId) {
					case "sg-1":
						rules := make([]*ec2.SecurityGroupRule, len(bootstrapIngress))
						return &ec2.AuthorizeSecurityGroupIngressOutput{
							SecurityGroupRules: rules,
						}, nil
					case "sg-2":
						rules := make([]*ec2.SecurityGroupRule, len(masterIngress))
						return &ec2.AuthorizeSecurityGroupIngressOutput{
							SecurityGroupRules: rules,
						}, nil
					case "sg-3":
						rules := make([]*ec2.SecurityGroupRule, len(workerIngress))
						return &ec2.AuthorizeSecurityGroupIngressOutput{
							SecurityGroupRules: rules,
						}, nil
					default:
						panic("should not be called")
					}
				},
			},
			expectedOut: &sgOutput{
				bootstrap:    "sg-1",
				controlPlane: "sg-2",
				compute:      "sg-3",
			},
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := sgInputOptions{
				infraID:          "infraID",
				vpcID:            "vpc-1",
				cidrV4Blocks:     []string{"10.0.0.0/16"},
				isPrivateCluster: test.private,
				tags:             map[string]string{"custom-tag": "custom-value"},
			}
			res, err := createSecurityGroups(context.TODO(), logger, &test.mockSvc, &input)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureInstance(t *testing.T) {
	expectedInstance := &ec2.Instance{
		InstanceId:       aws.String("instance-1"),
		PrivateIpAddress: aws.String("ip-1"),
	}
	tests := []struct {
		name        string
		mockEC2     mockEC2Client
		mockELB     mockELBClient
		expectedOut *ec2.Instance
		expectedErr string
	}{
		{
			name: "AWS SDK error listing instances",
			mockEC2: mockEC2Client{
				describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					return nil, errAwsSdk
				},
				runInstances: func(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
					panic("should not be called")
				},
			},
			mockELB: mockELBClient{
				registerTargets: func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to find instance: some AWS SDK error$`,
		},
		{
			name: "Creating instance fails getting default KMS key",
			mockEC2: mockEC2Client{
				describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					return &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{},
					}, nil
				},
				getDefaultKmsKey: func(*ec2.GetEbsDefaultKmsKeyIdInput) (*ec2.GetEbsDefaultKmsKeyIdOutput, error) {
					return nil, errAwsSdk
				},
				runInstances: func(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
					panic("should not be called")
				},
			},
			mockELB: mockELBClient{
				registerTargets: func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to get default KMS key: some AWS SDK error$`,
		},
		{
			name: "Creating instance fails",
			mockEC2: mockEC2Client{
				describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					return &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{},
					}, nil
				},
				getDefaultKmsKey: func(*ec2.GetEbsDefaultKmsKeyIdInput) (*ec2.GetEbsDefaultKmsKeyIdOutput, error) {
					return &ec2.GetEbsDefaultKmsKeyIdOutput{
						KmsKeyId: aws.String("kms-1"),
					}, nil
				},
				runInstances: func(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
					return nil, errAwsSdk
				},
			},
			mockELB: mockELBClient{
				registerTargets: func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^some AWS SDK error$`,
		},
		{
			name: "Instance created but no reservations found",
			mockEC2: mockEC2Client{
				describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					return &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{},
					}, nil
				},
				getDefaultKmsKey: func(*ec2.GetEbsDefaultKmsKeyIdInput) (*ec2.GetEbsDefaultKmsKeyIdOutput, error) {
					return &ec2.GetEbsDefaultKmsKeyIdOutput{
						KmsKeyId: aws.String("kms-1"),
					}, nil
				},
				runInstances: func(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
					return &ec2.Reservation{
						Instances: []*ec2.Instance{},
					}, nil
				},
			},
			mockELB: mockELBClient{
				registerTargets: func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^instance was not created$`,
		},
		{
			name: "Instance created but fails to register target groups",
			mockEC2: mockEC2Client{
				describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					return &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{},
					}, nil
				},
				getDefaultKmsKey: func(*ec2.GetEbsDefaultKmsKeyIdInput) (*ec2.GetEbsDefaultKmsKeyIdOutput, error) {
					return &ec2.GetEbsDefaultKmsKeyIdOutput{
						KmsKeyId: aws.String("kms-1"),
					}, nil
				},
				runInstances: func(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
					return &ec2.Reservation{
						Instances: []*ec2.Instance{expectedInstance},
					}, nil
				},
			},
			mockELB: mockELBClient{
				registerTargets: func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to register target group \(tg-1\): some AWS SDK error$`,
		},
		{
			name: "Instance created and target groups registered",
			mockEC2: mockEC2Client{
				describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					return &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{},
					}, nil
				},
				getDefaultKmsKey: func(*ec2.GetEbsDefaultKmsKeyIdInput) (*ec2.GetEbsDefaultKmsKeyIdOutput, error) {
					return &ec2.GetEbsDefaultKmsKeyIdOutput{
						KmsKeyId: aws.String("kms-1"),
					}, nil
				},
				runInstances: func(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
					return &ec2.Reservation{
						Instances: []*ec2.Instance{expectedInstance},
					}, nil
				},
			},
			mockELB: mockELBClient{
				registerTargets: func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
					return &elbv2.RegisterTargetsOutput{}, nil
				},
			},
			expectedOut: expectedInstance,
		},
		{
			name: "Instance found and target groups registered",
			mockEC2: mockEC2Client{
				describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					return &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{Instances: []*ec2.Instance{expectedInstance}},
						},
					}, nil
				},
				getDefaultKmsKey: func(*ec2.GetEbsDefaultKmsKeyIdInput) (*ec2.GetEbsDefaultKmsKeyIdOutput, error) {
					return &ec2.GetEbsDefaultKmsKeyIdOutput{
						KmsKeyId: aws.String("kms-1"),
					}, nil
				},
				runInstances: func(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
					panic("should not be called")
				},
			},
			mockELB: mockELBClient{
				registerTargets: func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
					return &elbv2.RegisterTargetsOutput{}, nil
				},
			},
			expectedOut: expectedInstance,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := instanceInputOptions{
				infraID:         "infraID",
				amiID:           "ami-1",
				name:            "instance-1",
				instanceType:    "type-1",
				subnetID:        "subnet-1",
				targetGroupARNs: []string{"tg-1"},
				tags:            map[string]string{"custom-tag": "custom-value"},
			}
			res, err := ensureInstance(context.TODO(), logger, &test.mockEC2, &test.mockELB, &input)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

type mockIAMClient struct {
	iamiface.IAMAPI

	// swappable functions
	getRole          func(*iam.GetRoleInput) (*iam.GetRoleOutput, error)
	createRole       func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error)
	getProfile       func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error)
	createProfile    func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error)
	addRoleToProfile func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error)
	getRolePolicy    func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error)
	putRolePolicy    func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error)
}

func (m *mockIAMClient) GetRoleWithContext(_ context.Context, in *iam.GetRoleInput, _ ...request.Option) (*iam.GetRoleOutput, error) {
	return m.getRole(in)
}

func (m *mockIAMClient) CreateRoleWithContext(_ context.Context, in *iam.CreateRoleInput, _ ...request.Option) (*iam.CreateRoleOutput, error) {
	return m.createRole(in)
}

func (m *mockIAMClient) GetInstanceProfileWithContext(_ context.Context, in *iam.GetInstanceProfileInput, _ ...request.Option) (*iam.GetInstanceProfileOutput, error) {
	return m.getProfile(in)
}

func (m *mockIAMClient) CreateInstanceProfileWithContext(_ context.Context, in *iam.CreateInstanceProfileInput, _ ...request.Option) (*iam.CreateInstanceProfileOutput, error) {
	return m.createProfile(in)
}

func (m *mockIAMClient) AddRoleToInstanceProfileWithContext(_ context.Context, in *iam.AddRoleToInstanceProfileInput, _ ...request.Option) (*iam.AddRoleToInstanceProfileOutput, error) {
	return m.addRoleToProfile(in)
}

func (m *mockIAMClient) GetRolePolicyWithContext(_ context.Context, in *iam.GetRolePolicyInput, _ ...request.Option) (*iam.GetRolePolicyOutput, error) {
	return m.getRolePolicy(in)
}

func (m *mockIAMClient) PutRolePolicyWithContext(_ context.Context, in *iam.PutRolePolicyInput, _ ...request.Option) (*iam.PutRolePolicyOutput, error) {
	return m.putRolePolicy(in)
}

func TestEnsureInstanceProfile(t *testing.T) {
	expectedRole := &iam.Role{RoleName: aws.String("name-role")}
	expectedProfile := &iam.InstanceProfile{
		InstanceProfileName: aws.String("name-profile"),
		Roles:               []*iam.Role{expectedRole},
	}
	tests := []struct {
		name         string
		existingRole string
		mockIAM      mockIAMClient
		expectedOut  *iam.InstanceProfile
		expectedErr  string
	}{
		{
			name: "AWS SDK error when retrieving role",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return nil, errAwsSdk
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					panic("should not be called")
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					panic("should not be called")
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					panic("should not be called")
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					panic("should not be called")
				},
				getRolePolicy: func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error) {
					panic("should not be called")
				},
				putRolePolicy: func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to get existing role: some AWS SDK error$`,
		},
		{
			name: "Creating role fails",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					return nil, errAwsSdk
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					panic("should not be called")
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					panic("should not be called")
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					panic("should not be called")
				},
				getRolePolicy: func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error) {
					panic("should not be called")
				},
				putRolePolicy: func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create role \(name-role\): some AWS SDK error$`,
		},
		{
			name: "Role created but AWS SDK error when retrieving profile",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					return &iam.CreateRoleOutput{
						Role: expectedRole,
					}, nil
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to get instance profile: some AWS SDK error$`,
		},
		{
			name: "Role created but creating profile fails",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					return &iam.CreateRoleOutput{
						Role: expectedRole,
					}, nil
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create instance profile: some AWS SDK error$`,
		},
		{
			name: "Role created and profile created but times out waiting for profile to exist",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				createRole: func(in *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					if len(in.Tags) == 0 {
						panic("role not tagged")
					}
					return &iam.CreateRoleOutput{
						Role: expectedRole,
					}, nil
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				createProfile: func(in *iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					if len(in.Tags) == 0 {
						panic("profile not tagged")
					}
					return &iam.CreateInstanceProfileOutput{
						InstanceProfile: &iam.InstanceProfile{},
					}, nil
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create instance profile: timed out waiting for instance profile to exist: .*$`,
		},
		{
			name: "Role created, profile created but adding role to profile fails",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return &iam.GetRoleOutput{
						Role: expectedRole,
					}, nil
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					panic("should not be called")
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return &iam.GetInstanceProfileOutput{
						InstanceProfile: &iam.InstanceProfile{},
					}, nil
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					panic("should not be called")
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to add role \(name-role\) to instance profile \(name-profile\): some AWS SDK error$`,
		},
		{
			name: "Role created, profile created, role added to profile but getting policy fails",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return &iam.GetRoleOutput{
						Role: expectedRole,
					}, nil
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					panic("should not be called")
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return &iam.GetInstanceProfileOutput{
						InstanceProfile: &iam.InstanceProfile{},
					}, nil
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					panic("should not be called")
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					return &iam.AddRoleToInstanceProfileOutput{}, nil
				},
				getRolePolicy: func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error) {
					return nil, errAwsSdk
				},
				putRolePolicy: func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to get role policy: some AWS SDK error$`,
		},
		{
			name: "Role created, profile created, role added to profile but adding policy fails",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return &iam.GetRoleOutput{
						Role: expectedRole,
					}, nil
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					panic("should not be called")
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return &iam.GetInstanceProfileOutput{
						InstanceProfile: &iam.InstanceProfile{},
					}, nil
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					panic("should not be called")
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					return &iam.AddRoleToInstanceProfileOutput{}, nil
				},
				getRolePolicy: func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				putRolePolicy: func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create role policy: some AWS SDK error$`,
		},
		{
			name: "Role created, profile created, role added to profile and policy created",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					return &iam.GetRoleOutput{
						Role: expectedRole,
					}, nil
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					panic("should not be called")
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return &iam.GetInstanceProfileOutput{
						InstanceProfile: expectedProfile,
					}, nil
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					panic("should not be called")
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					panic("should not be called")
				},
				getRolePolicy: func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error) {
					return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "", errAwsSdk)
				},
				putRolePolicy: func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
					return &iam.PutRolePolicyOutput{}, nil
				},
			},
			expectedOut: expectedProfile,
		},
		{
			name:         "Use user-supplied role",
			existingRole: "existing-role",
			mockIAM: mockIAMClient{
				getRole: func(*iam.GetRoleInput) (*iam.GetRoleOutput, error) {
					panic("should not be called")
				},
				createRole: func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
					panic("should not be called")
				},
				getProfile: func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
					return &iam.GetInstanceProfileOutput{
						InstanceProfile: expectedProfile,
					}, nil
				},
				createProfile: func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
					panic("should not be called")
				},
				addRoleToProfile: func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
					return &iam.AddRoleToInstanceProfileOutput{}, nil
				},
				getRolePolicy: func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error) {
					panic("should not be called")
				},
				putRolePolicy: func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
					panic("should not be called")
				},
			},
			expectedOut: expectedProfile,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := instanceProfileOptions{
				namePrefix:       "name",
				assumeRolePolicy: "policy-role",
				policyDocument:   "policy",
				tags:             map[string]string{"custom-tag": "custom-value"},
			}
			if len(test.existingRole) > 0 {
				input.roleName = test.existingRole
			}
			res, err := createInstanceProfile(context.TODO(), logger, &test.mockIAM, &input)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

type mockELBClient struct {
	elbv2iface.ELBV2API

	// swappable functions
	createListener func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error)

	createTargetGroup    func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error)
	describeTargetGroups func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error)

	createLoadBalancer    func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error)
	describeLoadBalancers func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error)
	modifyLBAttr          func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error)

	registerTargets func(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error)
}

func (m *mockELBClient) CreateListenerWithContext(_ context.Context, in *elbv2.CreateListenerInput, _ ...request.Option) (*elbv2.CreateListenerOutput, error) {
	return m.createListener(in)
}

func (m *mockELBClient) CreateTargetGroupWithContext(_ context.Context, in *elbv2.CreateTargetGroupInput, _ ...request.Option) (*elbv2.CreateTargetGroupOutput, error) {
	return m.createTargetGroup(in)
}

func (m *mockELBClient) DescribeTargetGroupsWithContext(_ context.Context, in *elbv2.DescribeTargetGroupsInput, _ ...request.Option) (*elbv2.DescribeTargetGroupsOutput, error) {
	return m.describeTargetGroups(in)
}

func (m *mockELBClient) CreateLoadBalancerWithContext(_ context.Context, in *elbv2.CreateLoadBalancerInput, _ ...request.Option) (*elbv2.CreateLoadBalancerOutput, error) {
	return m.createLoadBalancer(in)
}

func (m *mockELBClient) DescribeLoadBalancersWithContext(_ context.Context, in *elbv2.DescribeLoadBalancersInput, _ ...request.Option) (*elbv2.DescribeLoadBalancersOutput, error) {
	return m.describeLoadBalancers(in)
}

func (m *mockELBClient) ModifyLoadBalancerAttributesWithContext(_ context.Context, in *elbv2.ModifyLoadBalancerAttributesInput, _ ...request.Option) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
	return m.modifyLBAttr(in)
}

func (m *mockELBClient) RegisterTargetsWithContext(_ context.Context, in *elbv2.RegisterTargetsInput, _ ...request.Option) (*elbv2.RegisterTargetsOutput, error) {
	return m.registerTargets(in)
}

func TestEnsureTargetGroup(t *testing.T) {
	expectedTargetGroup := &elbv2.TargetGroup{}

	tests := []struct {
		name        string
		mockSvc     mockELBClient
		expectedOut *elbv2.TargetGroup
		expectedErr string
	}{
		{
			name: "AWS SDK error listing target groups",
			mockSvc: mockELBClient{
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, errAwsSdk
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list target groups: some AWS SDK error$`,
		},
		{
			name: "Target group already created",
			mockSvc: mockELBClient{
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return &elbv2.DescribeTargetGroupsOutput{TargetGroups: []*elbv2.TargetGroup{expectedTargetGroup}}, nil
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
			},
			expectedOut: expectedTargetGroup,
		},
		{
			name: "Creating target group fails",
			mockSvc: mockELBClient{
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return &elbv2.DescribeTargetGroupsOutput{}, nil
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^some AWS SDK error$`,
		},
		{
			name: "Creating target group succeeds",
			mockSvc: mockELBClient{
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, awserr.New(elbv2.ErrCodeTargetGroupNotFoundException, "", errAwsSdk)
				},
				createTargetGroup: func(in *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return &elbv2.CreateTargetGroupOutput{TargetGroups: []*elbv2.TargetGroup{expectedTargetGroup}}, nil
				},
			},
			expectedOut: expectedTargetGroup,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ensureTargetGroup(context.TODO(), logger, &test.mockSvc, "tgName", "vpc-1", readyzPath, apiPort, map[string]string{})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureLoadBalancer(t *testing.T) {
	expectedLoadBalancer := &elbv2.LoadBalancer{}

	tests := []struct {
		name        string
		mockSvc     mockELBClient
		expectedOut *elbv2.LoadBalancer
		expectedErr string
	}{
		{
			name: "AWS SDK error listing load balancers",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errAwsSdk
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list load balancers: some AWS SDK error$`,
		},
		{
			name: "Load balancer already exists but enabling cross zone fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return &elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to enable cross_zone attribute: some AWS SDK error$`,
		},
		{
			name: "Creating load balancer fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, awserr.New(elbv2.ErrCodeLoadBalancerNotFoundException, "", errAwsSdk)
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return nil, errAwsSdk
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^some AWS SDK error$`,
		},
		{
			name: "Load balancer created but enabling cross zone fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to enable cross_zone attribute: some AWS SDK error$`,
		},
		{
			name: "Load balancer created and cross zone attribute set",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
			},
			expectedOut: expectedLoadBalancer,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ensureLoadBalancer(context.TODO(), logger, &test.mockSvc, "lbName", []string{}, true, map[string]string{})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureInternalLoadBalancer(t *testing.T) {
	expectedLoadBalancer := &elbv2.LoadBalancer{}
	const tgAName = "infraID-aint"
	targetA := &elbv2.TargetGroup{
		TargetGroupName: aws.String(tgAName),
		TargetGroupArn:  aws.String("aint-arn"),
	}
	listenerA := &elbv2.Listener{
		ListenerArn: aws.String("aint-arn"),
		Port:        aws.Int64(apiPort),
	}
	const tgSName = "infraID-sint"
	targetS := &elbv2.TargetGroup{
		TargetGroupName: aws.String(tgSName),
		TargetGroupArn:  aws.String("sint-arn"),
	}
	listenerS := &elbv2.Listener{
		ListenerArn: aws.String("sint-arn"),
		Port:        aws.Int64(servicePort),
	}

	tests := []struct {
		name        string
		mockSvc     mockELBClient
		expectedOut *elbv2.LoadBalancer
		expectedErr string
	}{
		{
			name: "AWS SDK error listing load balancers",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errAwsSdk
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					panic("should not be called")
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					panic("should not be called")
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list load balancers: some AWS SDK error$`,
		},
		{
			name: "Load balancer already exists but creating target group A fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return &elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, errNotFound
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return nil, errAwsSdk
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create internalA target group: some AWS SDK error$`,
		},
		{
			name: "Creating load balancer fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, awserr.New(elbv2.ErrCodeLoadBalancerNotFoundException, "", errAwsSdk)
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return nil, errAwsSdk
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					panic("should not be called")
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^some AWS SDK error$`,
		},
		{
			name: "Load balancer created but listing target group fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, errAwsSdk
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create internalA target group: failed to list target groups: some AWS SDK error$`,
		},
		{
			name: "Load balancer created but creating target group A fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, awserr.New(elbv2.ErrCodeTargetGroupNotFoundException, "", errAwsSdk)
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return nil, errAwsSdk
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create internalA target group: some AWS SDK error$`,
		},
		{
			name: "Load balancer created but creating listener A fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, errAwsSdk
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return &elbv2.CreateTargetGroupOutput{}, nil
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create internalA target group: failed to list target groups: some AWS SDK error$`,
		},
		{
			name: "Load balancer created, target A and listener created but listing target groups fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(in *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					if aws.StringValue(in.Names[0]) == tgAName {
						return &elbv2.DescribeTargetGroupsOutput{
							TargetGroups: []*elbv2.TargetGroup{targetA},
						}, nil
					}
					return nil, errAwsSdk
				},
				createTargetGroup: func(in *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(in *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					if aws.Int64Value(in.Port) == apiPort {
						return &elbv2.CreateListenerOutput{
							Listeners: []*elbv2.Listener{listenerA},
						}, nil
					}
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create internalS target group: failed to list target groups: some AWS SDK error$`,
		},
		{
			name: "Load balancer created, target A and listener created but creating target S fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return &elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(in *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					switch aws.StringValue(in.Names[0]) {
					case tgAName:
						return &elbv2.DescribeTargetGroupsOutput{
							TargetGroups: []*elbv2.TargetGroup{targetA},
						}, nil
					case tgSName:
						return nil, errNotFound
					default:
						panic("should not be called")
					}
				},
				createTargetGroup: func(in *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					switch aws.StringValue(in.Name) {
					case tgAName:
						panic("should not be called")
					case tgSName:
						return nil, errAwsSdk
					default:
						panic("should not be called")
					}
				},
				createListener: func(in *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					if aws.Int64Value(in.Port) == apiPort {
						return &elbv2.CreateListenerOutput{
							Listeners: []*elbv2.Listener{listenerA},
						}, nil
					}
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create internalS target group: some AWS SDK error$`,
		},
		{
			name: "Load balancer created, target A and listener created, target S created but listener fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return &elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(in *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					switch aws.StringValue(in.Names[0]) {
					case tgAName:
						return &elbv2.DescribeTargetGroupsOutput{
							TargetGroups: []*elbv2.TargetGroup{targetA},
						}, nil
					case tgSName:
						return nil, errNotFound
					default:
						panic("should not be called")
					}
				},
				createTargetGroup: func(in *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					switch aws.StringValue(in.Name) {
					case tgAName:
						panic("should not be called")
					case tgSName:
						return &elbv2.CreateTargetGroupOutput{
							TargetGroups: []*elbv2.TargetGroup{targetS},
						}, nil
					default:
						panic("should not be called")
					}
				},
				createListener: func(in *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					switch aws.Int64Value(in.Port) {
					case apiPort:
						return &elbv2.CreateListenerOutput{
							Listeners: []*elbv2.Listener{listenerA},
						}, nil
					case servicePort:
						return nil, errAwsSdk
					default:
						panic("should not be called")
					}
				},
			},
			expectedErr: `^failed to create internalS listener: some AWS SDK error$`,
		},
		{
			name: "Load balancer, target groups and listeners created",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(in *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					switch aws.StringValue(in.Names[0]) {
					case tgAName:
						return &elbv2.DescribeTargetGroupsOutput{
							TargetGroups: []*elbv2.TargetGroup{targetA},
						}, nil
					case tgSName:
						return &elbv2.DescribeTargetGroupsOutput{
							TargetGroups: []*elbv2.TargetGroup{targetS},
						}, nil
					default:
						return nil, errAwsSdk
					}
				},
				createTargetGroup: func(in *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(in *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					switch aws.Int64Value(in.Port) {
					case apiPort:
						return &elbv2.CreateListenerOutput{
							Listeners: []*elbv2.Listener{listenerA},
						}, nil
					case servicePort:
						return &elbv2.CreateListenerOutput{
							Listeners: []*elbv2.Listener{listenerS},
						}, nil
					default:
						panic("should not be called")
					}
				},
			},
			expectedOut: expectedLoadBalancer,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := lbState{
				input: &lbInputOptions{
					infraID: "infraID",
					vpcID:   "vpc-1",
				},
				targetGroupArns: sets.New[string](),
			}
			res, err := state.ensureInternalLoadBalancer(context.TODO(), logger, &test.mockSvc, []string{}, map[string]string{})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestEnsureExternalLoadBalancer(t *testing.T) {
	expectedLoadBalancer := &elbv2.LoadBalancer{}
	const tgAName = "infraID-aext"
	targetA := &elbv2.TargetGroup{
		TargetGroupName: aws.String(tgAName),
		TargetGroupArn:  aws.String("aext-arn"),
	}
	listenerA := &elbv2.Listener{
		ListenerArn: aws.String("aext-arn"),
		Port:        aws.Int64(apiPort),
	}

	tests := []struct {
		name        string
		mockSvc     mockELBClient
		expectedOut *elbv2.LoadBalancer
		expectedErr string
	}{
		{
			name: "AWS SDK error listing load balancers",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errAwsSdk
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					panic("should not be called")
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					panic("should not be called")
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list load balancers: some AWS SDK error$`,
		},
		{
			name: "Load balancer already exists but creating target group A fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return &elbv2.DescribeLoadBalancersOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					panic("should not be called")
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, errNotFound
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return nil, errAwsSdk
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create external target group: some AWS SDK error$`,
		},
		{
			name: "Creating load balancer fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, awserr.New(elbv2.ErrCodeLoadBalancerNotFoundException, "", errAwsSdk)
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return nil, errAwsSdk
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					panic("should not be called")
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^some AWS SDK error$`,
		},
		{
			name: "Load balancer created but listing target group fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, errAwsSdk
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create external target group: failed to list target groups: some AWS SDK error$`,
		},
		{
			name: "Load balancer created but creating target group A fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, awserr.New(elbv2.ErrCodeTargetGroupNotFoundException, "", errAwsSdk)
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return nil, errAwsSdk
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create external target group: some AWS SDK error$`,
		},
		{
			name: "Load balancer created but creating listener A fails",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					return nil, errAwsSdk
				},
				createTargetGroup: func(*elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					return &elbv2.CreateTargetGroupOutput{}, nil
				},
				createListener: func(*elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to create external target group: failed to list target groups: some AWS SDK error$`,
		},
		{
			name: "Load balancer, target groups and listeners created",
			mockSvc: mockELBClient{
				describeLoadBalancers: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
					return nil, errNotFound
				},
				createLoadBalancer: func(*elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
					return &elbv2.CreateLoadBalancerOutput{
						LoadBalancers: []*elbv2.LoadBalancer{expectedLoadBalancer},
					}, nil
				},
				modifyLBAttr: func(*elbv2.ModifyLoadBalancerAttributesInput) (*elbv2.ModifyLoadBalancerAttributesOutput, error) {
					return &elbv2.ModifyLoadBalancerAttributesOutput{}, nil
				},
				describeTargetGroups: func(in *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
					switch aws.StringValue(in.Names[0]) {
					case tgAName:
						return &elbv2.DescribeTargetGroupsOutput{
							TargetGroups: []*elbv2.TargetGroup{targetA},
						}, nil
					default:
						return nil, errAwsSdk
					}
				},
				createTargetGroup: func(in *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
					panic("should not be called")
				},
				createListener: func(in *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
					switch aws.Int64Value(in.Port) {
					case apiPort:
						return &elbv2.CreateListenerOutput{
							Listeners: []*elbv2.Listener{listenerA},
						}, nil
					default:
						panic("should not be called")
					}
				},
			},
			expectedOut: expectedLoadBalancer,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := lbState{
				input: &lbInputOptions{
					infraID: "infraID",
					vpcID:   "vpc-1",
				},
				targetGroupArns: sets.New[string](),
			}
			res, err := state.ensureExternalLoadBalancer(context.TODO(), logger, &test.mockSvc, []string{}, map[string]string{})
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

type mockR53Client struct {
	route53iface.Route53API

	createHostedZone func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error)
	listHostedZones  func(*route53.ListHostedZonesInput, func(*route53.ListHostedZonesOutput, bool) bool) error
	changeTags       func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error)

	listRecordSets  func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
	changeRecordSet func(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
}

func (m *mockR53Client) CreateHostedZoneWithContext(_ context.Context, in *route53.CreateHostedZoneInput, _ ...request.Option) (*route53.CreateHostedZoneOutput, error) {
	return m.createHostedZone(in)
}

func (m *mockR53Client) ListHostedZonesPagesWithContext(_ context.Context, in *route53.ListHostedZonesInput, fn func(*route53.ListHostedZonesOutput, bool) bool, _ ...request.Option) error {
	return m.listHostedZones(in, fn)
}

func (m *mockR53Client) ChangeTagsForResourceWithContext(_ context.Context, in *route53.ChangeTagsForResourceInput, _ ...request.Option) (*route53.ChangeTagsForResourceOutput, error) {
	return m.changeTags(in)
}

func (m *mockR53Client) ListResourceRecordSetsWithContext(_ context.Context, in *route53.ListResourceRecordSetsInput, _ ...request.Option) (*route53.ListResourceRecordSetsOutput, error) {
	return m.listRecordSets(in)
}

func (m *mockR53Client) ChangeResourceRecordSetsWithContext(_ context.Context, in *route53.ChangeResourceRecordSetsInput, _ ...request.Option) (*route53.ChangeResourceRecordSetsOutput, error) {
	return m.changeRecordSet(in)
}

func TestEnsurePrivateHostedZone(t *testing.T) {
	privateHostedZone := &route53.HostedZone{
		Id:     aws.String("hostedzone-1"),
		Name:   aws.String("domain.com."),
		Config: &route53.HostedZoneConfig{PrivateZone: aws.Bool(true)},
	}
	tests := []struct {
		name        string
		mockSvc     mockR53Client
		isPrivate   bool
		expectedOut *route53.HostedZone
		expectedErr string
	}{
		{
			name: "AWS SDK error listing hosted zones",
			mockSvc: mockR53Client{
				listHostedZones: func(*route53.ListHostedZonesInput, func(*route53.ListHostedZonesOutput, bool) bool) error {
					return errAwsSdk
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					panic("should not be called")
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to list hosted zones: .*$`,
		},
		{
			name: "Create hosted zone fails",
			mockSvc: mockR53Client{
				listHostedZones: func(*route53.ListHostedZonesInput, func(*route53.ListHostedZonesOutput, bool) bool) error {
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					return nil, errAwsSdk
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `^failed to create private hosted zone: failed to create hosted zone: .*$`,
		},
		{
			name: "Hosted zone created but tagging fails",
			mockSvc: mockR53Client{
				listHostedZones: func(*route53.ListHostedZonesInput, func(*route53.ListHostedZonesOutput, bool) bool) error {
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					return &route53.CreateHostedZoneOutput{
						HostedZone: privateHostedZone,
					}, nil
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `^failed to tag private hosted zone: some AWS SDK error$`,
		},
		{
			name: "Existing zone and tagging succeeds but listing records fails",
			mockSvc: mockR53Client{
				listHostedZones: func(_ *route53.ListHostedZonesInput, fn func(*route53.ListHostedZonesOutput, bool) bool) error {
					fn(&route53.ListHostedZonesOutput{
						HostedZones: []*route53.HostedZone{privateHostedZone},
					}, true)
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					panic("should not be called")
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					return &route53.ChangeTagsForResourceOutput{}, nil
				},
				listRecordSets: func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
					return nil, errAwsSdk
				},
				changeRecordSet: func(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
					panic("change records should not be called")
				},
			},
			expectedErr: `failed to find SOA record set for private zone: some AWS SDK error$`,
		},
		{
			name: "Created and tagged zone but SOA record set not found",
			mockSvc: mockR53Client{
				listHostedZones: func(_ *route53.ListHostedZonesInput, fn func(*route53.ListHostedZonesOutput, bool) bool) error {
					fn(&route53.ListHostedZonesOutput{
						HostedZones: []*route53.HostedZone{privateHostedZone},
					}, true)
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					panic("should not be called")
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					return &route53.ChangeTagsForResourceOutput{}, nil
				},
				listRecordSets: func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
					return &route53.ListResourceRecordSetsOutput{}, nil
				},
				changeRecordSet: func(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
					panic("should not be called")
				},
			},
			expectedErr: `failed to find SOA record set for private zone: not found$`,
		},
		{
			name: "Created and tagged zone but SOA not found",
			mockSvc: mockR53Client{
				listHostedZones: func(_ *route53.ListHostedZonesInput, fn func(*route53.ListHostedZonesOutput, bool) bool) error {
					fn(&route53.ListHostedZonesOutput{
						HostedZones: []*route53.HostedZone{privateHostedZone},
					}, true)
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					panic("create should not be called")
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					return &route53.ChangeTagsForResourceOutput{}, nil
				},
				listRecordSets: func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
					return &route53.ListResourceRecordSetsOutput{
						ResourceRecordSets: []*route53.ResourceRecordSet{
							{
								Name:            aws.String("domain.com."),
								Type:            aws.String("SOA"),
								ResourceRecords: []*route53.ResourceRecord{},
							},
						},
					}, nil
				},
				changeRecordSet: func(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
					panic("change records should not be called")
				},
			},
			expectedErr: `failed to find SOA record for private zone$`,
		},
		{
			name: "Created and tagged zone but SOA has wrong number of fields",
			mockSvc: mockR53Client{
				listHostedZones: func(_ *route53.ListHostedZonesInput, fn func(*route53.ListHostedZonesOutput, bool) bool) error {
					fn(&route53.ListHostedZonesOutput{
						HostedZones: []*route53.HostedZone{privateHostedZone},
					}, true)
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					panic("create should not be called")
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					return &route53.ChangeTagsForResourceOutput{}, nil
				},
				listRecordSets: func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
					return &route53.ListResourceRecordSetsOutput{
						ResourceRecordSets: []*route53.ResourceRecordSet{
							{
								Name: aws.String("domain.com."),
								Type: aws.String("SOA"),
								ResourceRecords: []*route53.ResourceRecord{
									{Value: aws.String("domain email 1")},
								},
							},
						},
					}, nil
				},
				changeRecordSet: func(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
					panic("change records should not be called")
				},
			},
			expectedErr: `SOA record value has [^7] fields, expected 7$`,
		},
		{
			name: "Created and tagged zone but setting SOA minimum TTL fails",
			mockSvc: mockR53Client{
				listHostedZones: func(_ *route53.ListHostedZonesInput, fn func(*route53.ListHostedZonesOutput, bool) bool) error {
					fn(&route53.ListHostedZonesOutput{
						HostedZones: []*route53.HostedZone{privateHostedZone},
					}, true)
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					panic("create should not be called")
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					return &route53.ChangeTagsForResourceOutput{}, nil
				},
				listRecordSets: func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
					return &route53.ListResourceRecordSetsOutput{
						ResourceRecordSets: []*route53.ResourceRecordSet{
							{
								Name: aws.String("domain.com."),
								Type: aws.String("SOA"),
								ResourceRecords: []*route53.ResourceRecord{
									{Value: aws.String("domain email 1 7200 900 1209600 86400")},
								},
							},
						},
					}, nil
				},
				changeRecordSet: func(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
					return nil, errAwsSdk
				},
			},
			expectedErr: `failed to set SOA TTL to minimum: some AWS SDK error$`,
		},
		{
			name: "Zone created, tagged zone and SOA set to minimum TTL",
			mockSvc: mockR53Client{
				listHostedZones: func(_ *route53.ListHostedZonesInput, fn func(*route53.ListHostedZonesOutput, bool) bool) error {
					fn(&route53.ListHostedZonesOutput{
						HostedZones: []*route53.HostedZone{privateHostedZone},
					}, true)
					return nil
				},
				createHostedZone: func(*route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
					panic("create should not be called")
				},
				changeTags: func(*route53.ChangeTagsForResourceInput) (*route53.ChangeTagsForResourceOutput, error) {
					return &route53.ChangeTagsForResourceOutput{}, nil
				},
				listRecordSets: func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
					return &route53.ListResourceRecordSetsOutput{
						ResourceRecordSets: []*route53.ResourceRecordSet{
							{
								Name: aws.String("domain.com."),
								Type: aws.String("SOA"),
								ResourceRecords: []*route53.ResourceRecord{
									{Value: aws.String("domain email 1 7200 900 1209600 86400")},
								},
							},
						},
					}, nil
				},
				changeRecordSet: func(in *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
					updated := aws.StringValue(in.ChangeBatch.Changes[0].ResourceRecordSet.ResourceRecords[0].Value)
					if updated != "domain email 1 7200 900 1209600 60" {
						panic("SOA minimum TTL not set")
					}
					return &route53.ChangeResourceRecordSetsOutput{}, nil
				},
			},
			expectedOut: privateHostedZone,
		},
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	input := dnsInputOptions{
		clusterDomain: "domain.com",
		vpcID:         "vpc-1",
		region:        "region-1",
		infraID:       "infraID",
		tags:          map[string]string{"custom-tag": "custom-value"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ensurePrivateZone(context.TODO(), logger, &test.mockSvc, &input)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, res)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestFQDN(t *testing.T) {
	assert.Equal(t, "domain.", fqdn("domain"))
	assert.Equal(t, "domain.", fqdn("domain."))
	assert.Equal(t, "domain.com.", fqdn("domain.com"))
	assert.Equal(t, "domain.com.", fqdn("domain.com."))
	assert.Equal(t, "", fqdn(""))
	assert.Equal(t, ".", fqdn("."))
}
