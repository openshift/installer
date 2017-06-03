#!/bin/bash

# TODO: Use jq instead of sed everywhere
# TODO: Make xargs ignore all errors (might be pre-empted): use xargs -I{} sh -c 'aws ec2 ... {} || true'
# TODO: Delete all other kind of resources

echo Deleting all private zones
aws route53 list-hosted-zones | jq ".HostedZones[] | select(.Config.PrivateZone == true) | .Id" | xargs -L 1 aws route53 delete-hosted-zone --id
aws route53 list-resource-record-sets --hosted-zone-id 


echo Deleting all ASGs
aws autoscaling describe-auto-scaling-groups | jq ".AutoScalingGroups[] .AutoScalingGroupName" | xargs -L 1 aws autoscaling delete-auto-scaling-group --force-delete --auto-scaling-group-name

echo Deleting all ELBs
aws elb describe-load-balancers | jq ".LoadBalancerDescriptions[] .LoadBalancerName" | xargs -L 1 aws elb delete-load-balancer --load-balancer-name

echo Deleting all NAT gateways
aws ec2 describe-nat-gateways | jq '.NatGateways[] | select(.State == "available") | .NatGatewayId' | xargs -L 1 -I {} aws ec2 delete-nat-gateway --nat-gateway-id {} 1> /dev/null

echo "Waiting for the Elastic Network Interfaces from the ELBs, and for the NAT gateways, to be detached/deleted"
sleep 30

echo Disassociating/Releasing all EIPs
aws ec2 describe-addresses | grep AssociationId | sed -E 's/^.*(eipassoc-[a-z0-9]+).*$/\1/' | xargs -L 1 aws ec2 disassociate-address --association-id
aws ec2 describe-addresses | grep AllocationId | sed -E 's/^.*(eipalloc-[a-z0-9]+).*$/\1/' | xargs -L 1 aws ec2 release-address --allocation-id

echo Detaching/Deleting all Internet Gateways
aws ec2 describe-internet-gateways | jq '.InternetGateways[] | "--internet-gateway-id=" + .InternetGatewayId + " --vpc-id=" + .Attachments?[].VpcId' -r | xargs -L 1 aws ec2 detach-internet-gateway
aws ec2 describe-internet-gateways | jq '.InternetGateways[] .InternetGatewayId' | xargs -L 1 aws ec2 delete-internet-gateway --internet-gateway-id

echo Deleting all subnets
aws ec2 describe-subnets | jq ".Subnets[] .SubnetId" | xargs -L 1 aws ec2 delete-subnet --subnet-id

echo "Deleting all route tables (except the main ones)"
#aws ec2 describe-route-tables | jq '.RouteTables[] | select(.Routes?[] | select(.DestinationCidrBlock == "0.0.0.0/0")) | .RouteTableId' | xargs -L 1 aws ec2 delete-route --destination-cidr-block 0.0.0.0/0 --route-table-id # First, delete the default route, as they might be pointed at a non-existing gateway.
# shellcheck disable=SC2016
aws ec2 describe-route-tables --query 'RouteTables[?Associations[0].Main != `true`]' | jq ".[] .RouteTableId" | xargs -L 1 aws ec2 delete-route-table --route-table-id

echo Deleting security groups
aws ec2 describe-security-groups | jq ".SecurityGroups[] | select(.GroupName != \"default\" and (.IpPermissions|length > 0)) | \"--group-id \" + .GroupId + \" --ip-permissions '\" + (.IpPermissions|tostring) + \"'\"" -r | xargs -L 1 aws ec2 revoke-security-group-ingress
aws ec2 describe-security-groups | jq ".SecurityGroups[] | select(.GroupName != \"default\" and (.IpPermissionsEgress|length > 0)) | \"--group-id \" + .GroupId + \" --ip-permissions '\" + (.IpPermissionsEgress|tostring) + \"'\"" -r | xargs -L 1 aws ec2 revoke-security-group-egress
aws ec2 describe-security-groups | jq '.SecurityGroups[] | select(.GroupName != "default") | .GroupId' | xargs -L 1 -I {} sh -c 'aws ec2 delete-security-group --group-id {} || true'

echo Deleting VPCs
aws ec2 describe-vpcs | jq ".Vpcs[] .VpcId" | xargs -L 1 -I {} sh -c 'aws ec2 delete-vpc --vpc-id {} || true'
