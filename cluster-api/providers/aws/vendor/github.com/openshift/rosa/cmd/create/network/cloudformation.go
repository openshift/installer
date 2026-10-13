package network

// This file is used for binary builds
// nolint:lll
const CloudFormationTemplateFile = `
AWSTemplateFormatVersion: '2010-09-09'
Description: CloudFormation template to create a ROSA Quickstart default VPC.
  This CloudFormation template may not work with rosa CLI versions later than 1.2.48.
  Please ensure that you are using the compatible CLI version before deploying this template.

Parameters:
  AvailabilityZoneCount:
    Type: Number
    Description: "Number of Availability Zones to use"
    Default: 1
    MinValue: 1
    MaxValue: 4
  AZ1:
    Type: String
    Description: "First availability zone to use"
    Default: ""
  AZ2:
    Type: String
    Description: "Second availability zone to use"
    Default: ""
  AZ3:
    Type: String
    Description: "Third availability zone to use"
    Default: ""
  AZ4:
    Type: String
    Description: "Fourth availability zone to use"
    Default: ""
  Region:
    Type: String
    Description: "AWS Region"
    Default: "us-west-2"
  Name:
    Type: String
    Description: "Name prefix for resources"
  VpcCidr:
    Type: String
    Description: CIDR block for the VPC
    Default: '10.0.0.0/16'

Conditions:
  AZ1Explicit: !Not [!Equals [!Ref AZ1, ""]]
  AZ2Explicit: !Not [!Equals [!Ref AZ2, ""]]
  AZ3Explicit: !Not [!Equals [!Ref AZ3, ""]]
  AZ4Explicit: !Not [!Equals [!Ref AZ4, ""]]

  ExplicitAZs:   !Or [!Condition AZ1Explicit, !Condition AZ2Explicit, !Condition AZ3Explicit, !Condition AZ4Explicit]
  NoExplicitAZs: !Not [!Condition ExplicitAZs]

  AZ4Implicit: !Equals [!Ref AvailabilityZoneCount, 4]
  AZ3Implicit: !Or [!Equals [!Ref AvailabilityZoneCount, 3], !Condition AZ4Implicit]
  AZ2Implicit: !Or [!Equals [!Ref AvailabilityZoneCount, 2], !Condition AZ3Implicit]
  AZ1Implicit: !Or [!Equals [!Ref AvailabilityZoneCount, 1], !Condition AZ2Implicit]

  One:   !Or [!And [!Condition ExplicitAZs, !Condition AZ1Explicit], !And [!Condition NoExplicitAZs, !Condition AZ1Implicit]]
  Two:   !Or [!And [!Condition ExplicitAZs, !Condition AZ2Explicit], !And [!Condition NoExplicitAZs, !Condition AZ2Implicit]]
  Three: !Or [!And [!Condition ExplicitAZs, !Condition AZ3Explicit], !And [!Condition NoExplicitAZs, !Condition AZ3Implicit]]
  Four:  !Or [!And [!Condition ExplicitAZs, !Condition AZ4Explicit], !And [!Condition NoExplicitAZs, !Condition AZ4Implicit]]

Resources:
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: !Ref VpcCidr
      EnableDnsSupport: true
      EnableDnsHostnames: true
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  S3VPCEndpoint:
    Type: AWS::EC2::VPCEndpoint
    Properties:
      VpcId: !Ref VPC
      ServiceName: !Sub "com.amazonaws.${Region}.s3"
      VpcEndpointType: Gateway
      RouteTableIds:
        - !Ref PublicRouteTable
        - !Ref PrivateRouteTable

  SubnetPublic1:
    Type: AWS::EC2::Subnet
    Condition: One
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [0, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ1, !Select [0, !GetAZs '']]
      MapPublicIpOnLaunch: true
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Public-Subnet-1"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/elb'
          Value: '1'

  SubnetPrivate1:
    Type: AWS::EC2::Subnet
    Condition: One
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [1, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ1, !Select [0, !GetAZs '']]
      MapPublicIpOnLaunch: false
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Private-Subnet-1"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/internal-elb'
          Value: '1'

  SubnetPublic2:
    Type: AWS::EC2::Subnet
    Condition: Two
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [2, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ2, !Select [1, !GetAZs '']]
      MapPublicIpOnLaunch: true
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Public-Subnet-2"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/elb'
          Value: '1'

  SubnetPrivate2:
    Type: AWS::EC2::Subnet
    Condition: Two
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [3, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ2, !Select [1, !GetAZs '']]
      MapPublicIpOnLaunch: false
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Private-Subnet-2"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/internal-elb'
          Value: '1'

  SubnetPublic3:
    Type: AWS::EC2::Subnet
    Condition: Three
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [4, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ3, !Select [2, !GetAZs '']]
      MapPublicIpOnLaunch: true
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Public-Subnet-3"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/elb'
          Value: '1'

  SubnetPrivate3:
    Type: AWS::EC2::Subnet
    Condition: Three
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [5, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ3, !Select [2, !GetAZs '']]
      MapPublicIpOnLaunch: false
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Private-Subnet-3"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/internal-elb'
          Value: '1'

  SubnetPublic4:
    Type: AWS::EC2::Subnet
    Condition: Four
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [6, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ4, !Select [3, !GetAZs '']]
      MapPublicIpOnLaunch: true
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Public-Subnet-4"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/elb'
          Value: '1'

  SubnetPrivate4:
    Type: AWS::EC2::Subnet
    Condition: Four
    Properties:
      VpcId: !Ref VPC
      CidrBlock: !Select [7, !Cidr [!Ref VpcCidr, 8, 8]]
      AvailabilityZone: !If [ExplicitAZs, !Ref AZ4, !Select [3, !GetAZs '']]
      MapPublicIpOnLaunch: false
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Private-Subnet-4"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'kubernetes.io/role/internal-elb'
          Value: '1'

  InternetGateway:
    Type: AWS::EC2::InternetGateway
    Properties:
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  AttachGateway:
    Type: AWS::EC2::VPCGatewayAttachment
    Properties:
      VpcId: !Ref VPC
      InternetGatewayId: !Ref InternetGateway

  ElasticIP1:
    Condition: One
    Type: AWS::EC2::EIP
    Properties:
      Domain: vpc
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  ElasticIP2:
    Condition: Two
    Type: AWS::EC2::EIP
    Properties:
      Domain: vpc
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  ElasticIP3:
    Condition: Three
    Type: AWS::EC2::EIP
    Properties:
      Domain: vpc
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  ElasticIP4:
    Condition: Four
    Type: AWS::EC2::EIP
    Properties:
      Domain: vpc
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  NATGateway1:
    Condition: One
    Type: 'AWS::EC2::NatGateway'
    Properties:
      AllocationId: !GetAtt ElasticIP1.AllocationId
      SubnetId: !Ref SubnetPublic1
      Tags:
        - Key: Name
          Value: !Sub "${Name}-NAT-1"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  NATGateway2:
    Condition: Two
    Type: 'AWS::EC2::NatGateway'
    Properties:
      AllocationId: !GetAtt ElasticIP2.AllocationId
      SubnetId: !Ref SubnetPublic2
      Tags:
        - Key: Name
          Value: !Sub "${Name}-NAT-2"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  NATGateway3:
    Condition: Three
    Type: 'AWS::EC2::NatGateway'
    Properties:
      AllocationId: !GetAtt ElasticIP3.AllocationId
      SubnetId: !Ref SubnetPublic3
      Tags:
        - Key: Name
          Value: !Sub "${Name}-NAT-3"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  NATGateway4:
    Condition: Four
    Type: 'AWS::EC2::NatGateway'
    Properties:
      AllocationId: !GetAtt ElasticIP4.AllocationId
      SubnetId: !Ref SubnetPublic4
      Tags:
        - Key: Name
          Value: !Sub "${Name}-NAT-4"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  PublicRouteTable:
    Type: AWS::EC2::RouteTable
    Properties:
      VpcId: !Ref VPC
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  PublicRoute:
    Type: AWS::EC2::Route
    DependsOn: AttachGateway
    Properties:
      RouteTableId: !Ref PublicRouteTable
      DestinationCidrBlock: 0.0.0.0/0
      GatewayId: !Ref InternetGateway

  PrivateRouteTable:
    Type: AWS::EC2::RouteTable
    Properties:
      VpcId: !Ref VPC
      Tags:
        - Key: Name
          Value: !Sub "${Name}-Private-Route-Table"
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'
        - Key: 'service'
          Value: 'ROSA'

  PrivateRoute:
    Type: AWS::EC2::Route
    Properties:
      RouteTableId: !Ref PrivateRouteTable
      DestinationCidrBlock: 0.0.0.0/0
      NatGatewayId: !If
        - One
        - !Ref NATGateway1
        - !If
          - Two
          - !Ref NATGateway2
          - !If
            - Three
            - !Ref NATGateway3
            - !If
              - Four
              - !Ref NATGateway4
              - !Ref "AWS::NoValue"

  PublicSubnetRouteTableAssociation1:
    Condition: One
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPublic1
      RouteTableId: !Ref PublicRouteTable

  PublicSubnetRouteTableAssociation2:
    Condition: Two
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPublic2
      RouteTableId: !Ref PublicRouteTable

  PublicSubnetRouteTableAssociation3:
    Condition: Three
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPublic3
      RouteTableId: !Ref PublicRouteTable

  PublicSubnetRouteTableAssociation4:
    Condition: Four
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPublic4
      RouteTableId: !Ref PublicRouteTable

  PrivateSubnetRouteTableAssociation1:
    Condition: One
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPrivate1
      RouteTableId: !Ref PrivateRouteTable

  PrivateSubnetRouteTableAssociation2:
    Condition: Two
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPrivate2
      RouteTableId: !Ref PrivateRouteTable

  PrivateSubnetRouteTableAssociation3:
    Condition: Three
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPrivate3
      RouteTableId: !Ref PrivateRouteTable

  PrivateSubnetRouteTableAssociation4:
    Condition: Four
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref SubnetPrivate4
      RouteTableId: !Ref PrivateRouteTable
  
  SecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: "Authorize inbound VPC traffic"
      VpcId: !Ref VPC
      SecurityGroupIngress:
        - IpProtocol: -1
          FromPort: 0
          ToPort: 0
          CidrIp: "10.0.0.0/16"
      SecurityGroupEgress:
        - IpProtocol: -1
          FromPort: 0
          ToPort: 0
          CidrIp: 0.0.0.0/0
      Tags:
        - Key: Name
          Value: !Ref Name
        - Key: 'service'
          Value: 'ROSA'
        - Key: 'rosa_managed_policies'
          Value: 'true'
        - Key: 'rosa_hcp_policies'
          Value: 'true'

  EC2VPCEndpoint:
    Type: AWS::EC2::VPCEndpoint
    Properties:
      VpcId: !Ref VPC
      ServiceName: !Sub "com.amazonaws.${Region}.ec2"
      PrivateDnsEnabled: true
      VpcEndpointType: Interface
      SubnetIds:
        - !If [One, !Ref SubnetPrivate1, !Ref "AWS::NoValue"]
        - !If [Two, !Ref SubnetPrivate2, !Ref "AWS::NoValue"]
        - !If [Three, !Ref SubnetPrivate3, !Ref "AWS::NoValue"]
        - !If [Four, !Ref SubnetPrivate4, !Ref "AWS::NoValue"]
      SecurityGroupIds:
        - !Ref SecurityGroup

  KMSVPCEndpoint:
    Type: AWS::EC2::VPCEndpoint
    Properties:
      VpcId: !Ref VPC
      ServiceName: !Sub "com.amazonaws.${Region}.kms"
      PrivateDnsEnabled: true
      VpcEndpointType: Interface
      SubnetIds:
        - !If [One, !Ref SubnetPrivate1, !Ref "AWS::NoValue"]
        - !If [Two, !Ref SubnetPrivate2, !Ref "AWS::NoValue"]
        - !If [Three, !Ref SubnetPrivate3, !Ref "AWS::NoValue"]
        - !If [Four, !Ref SubnetPrivate4, !Ref "AWS::NoValue"]
      SecurityGroupIds:
        - !Ref SecurityGroup

  STSVPCEndpoint:
    Type: AWS::EC2::VPCEndpoint
    Properties:
      VpcId: !Ref VPC
      ServiceName: !Sub "com.amazonaws.${Region}.sts"
      PrivateDnsEnabled: true
      VpcEndpointType: Interface
      SubnetIds:
        - !If [One, !Ref SubnetPrivate1, !Ref "AWS::NoValue"]
        - !If [Two, !Ref SubnetPrivate2, !Ref "AWS::NoValue"]
        - !If [Three, !Ref SubnetPrivate3, !Ref "AWS::NoValue"]
        - !If [Four, !Ref SubnetPrivate4, !Ref "AWS::NoValue"]
      SecurityGroupIds:
        - !Ref SecurityGroup

  EcrApiVPCEndpoint:
    Type: AWS::EC2::VPCEndpoint
    Properties:
      VpcId: !Ref VPC
      ServiceName: !Sub "com.amazonaws.${Region}.ecr.api"
      PrivateDnsEnabled: true
      VpcEndpointType: Interface
      SubnetIds:
        - !If [One, !Ref SubnetPrivate1, !Ref "AWS::NoValue"]
        - !If [Two, !Ref SubnetPrivate2, !Ref "AWS::NoValue"]
        - !If [Three, !Ref SubnetPrivate3, !Ref "AWS::NoValue"]
        - !If [Four, !Ref SubnetPrivate4, !Ref "AWS::NoValue"]
      SecurityGroupIds:
        - !Ref SecurityGroup

  EcrDkrVPCEndpoint:
    Type: AWS::EC2::VPCEndpoint
    Properties:
      VpcId: !Ref VPC
      ServiceName: !Sub "com.amazonaws.${Region}.ecr.dkr"
      PrivateDnsEnabled: true
      VpcEndpointType: Interface
      SubnetIds:
        - !If [One, !Ref SubnetPrivate1, !Ref "AWS::NoValue"]
        - !If [Two, !Ref SubnetPrivate2, !Ref "AWS::NoValue"]
        - !If [Three, !Ref SubnetPrivate3, !Ref "AWS::NoValue"]
        - !If [Four, !Ref SubnetPrivate4, !Ref "AWS::NoValue"]
      SecurityGroupIds: 
        - !Ref SecurityGroup

Outputs:
  VPCId:
    Description: "VPC Id"
    Value: !Ref VPC
    Export:
      Name: !Sub "${Name}-VPCId"

  VPCEndpointId:
    Description: The ID of the VPC Endpoint
    Value: !Ref S3VPCEndpoint
    Export:
      Name: !Sub "${Name}-VPCEndpointId"

  PublicSubnets:
    Description: "Public Subnet Ids"
    Value: !Join [",", [!If [One, !Ref SubnetPublic1, !Ref "AWS::NoValue"], !If [Two, !Ref SubnetPublic2, !Ref "AWS::NoValue"], !If [Three, !Ref SubnetPublic3, !Ref "AWS::NoValue"], !If [Four, !Ref SubnetPublic4, !Ref "AWS::NoValue"]]]
    Export:
      Name: !Sub "${Name}-PublicSubnets"

  PrivateSubnets:
    Description: "Private Subnet Ids"
    Value: !Join [",", [!If [One, !Ref SubnetPrivate1, !Ref "AWS::NoValue"], !If [Two, !Ref SubnetPrivate2, !Ref "AWS::NoValue"], !If [Three, !Ref SubnetPrivate3, !Ref "AWS::NoValue"], !If [Four, !Ref SubnetPublic4, !Ref "AWS::NoValue"]]]
    Export:
      Name: !Sub "${Name}-PrivateSubnets"

  EIP1AllocationId:
    Condition: One
    Description: Allocation ID for ElasticIP1
    Value: !GetAtt ElasticIP1.AllocationId
    Export:
      Name: !Sub "${Name}-EIP1-AllocationId"

  EIP2AllocationId:
    Condition: Two
    Description: Allocation ID for ElasticIP2
    Value: !GetAtt ElasticIP2.AllocationId
    Export:
      Name: !Sub "${Name}-EIP2-AllocationId"

  EIP3AllocationId:
    Condition: Three
    Description: Allocation ID for ElasticIP3
    Value: !GetAtt ElasticIP3.AllocationId
    Export:
      Name: !Sub "${Name}-EIP3-AllocationId"

  EIP4AllocationId:
    Condition: Four
    Description: Allocation ID for ElasticIP4
    Value: !GetAtt ElasticIP4.AllocationId
    Export:
      Name: !Sub "${Name}-EIP4-AllocationId"

  NatGatewayId:
    Description: The NAT Gateway IDs
    Value: !Join [",", [!If [One, !Ref NATGateway1, !Ref "AWS::NoValue"], !If [Two, !Ref NATGateway2, !Ref "AWS::NoValue"], !If [Three, !Ref NATGateway3, !Ref "AWS::NoValue"], !If [Four, !Ref NATGateway4, !Ref "AWS::NoValue"]]]
    Export:
      Name: !Sub "${Name}-NatGatewayId"

  InternetGatewayId:
    Description: The ID of the Internet Gateway
    Value: !Ref InternetGateway
    Export:
      Name: !Sub "${Name}-InternetGatewayId"

  PublicRouteTableId:
    Description: The ID of the public route table
    Value: !Ref PublicRouteTable
    Export:
      Name: !Sub "${Name}-PublicRouteTableId"

  PrivateRouteTableId:
    Description: The ID of the private route table
    Value: !Ref PrivateRouteTable
    Export:
      Name: !Sub "${Name}-PrivateRouteTableId"

  EC2VPCEndpointId:
    Description: The ID of the EC2 VPC Endpoint
    Value: !Ref EC2VPCEndpoint
    Export:
      Name: !Sub "${Name}-EC2VPCEndpointId"

  KMSVPCEndpointId:
    Description: The ID of the KMS VPC Endpoint
    Value: !Ref KMSVPCEndpoint
    Export:
      Name: !Sub "${Name}-KMSVPCEndpointId"

  STSVPCEndpointId:
    Description: The ID of the STS VPC Endpoint
    Value: !Ref STSVPCEndpoint
    Export:
      Name: !Sub "${Name}-STSVPCEndpointId"

  EcrApiVPCEndpointId:
    Description: The ID of the ECR API VPC Endpoint
    Value: !Ref EcrApiVPCEndpoint
    Export:
      Name: !Sub "${Name}-EcrApiVPCEndpointId"

  EcrDkrVPCEndpointId:
    Description: The ID of the ECR DKR VPC Endpoint
    Value: !Ref EcrDkrVPCEndpoint
    Export:
      Name: !Sub "${Name}-EcrDkrVPCEndpointId"`
