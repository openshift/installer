#!/bin/sh

die() {
	# shellcheck disable=SC2059
	printf "${@}" >&2
	exit 1
}

REGION="$(aws configure get region)"
if test -z "${REGION}"
then
	die 'no AWS region selected\n'
fi

printf 'count\tlimit\tregion\tcode\tname\n' || die 'failed to write header'
aws --region us-east-1 support describe-trusted-advisor-checks --language en --query "checks[? category == 'service_limits'].{id: @.id, name: @.name}" --output text | while read -r CHECK
do
	CHECK_ID="$(echo "${CHECK}" | cut -d '	' -f 1)" || die 'failed to extract ID from %s\n' "${CHECK}"
	CHECK_NAME="$(echo "${CHECK}" | cut -d '	' -f 2)" || die 'failed to extract name from %s\n' "${CHECK}"
	RESULT="$(aws --region us-east-1 support describe-trusted-advisor-check-result --check-id "${CHECK_ID}" --query "join(\`\\n\`, result.flaggedResources[].join(\`\\t\`, [@.metadata[4] || '0', @.metadata[3], @.region || '-', '${CHECK_ID}', @.metadata[2]]))" --output text)" || die 'failed to check %s (%s)\n' "${CHECK_ID}" "${CHECK_NAME}"
	if test -n "${RESULT}"
	then
		echo "${RESULT}" || die 'failed to write result for %s (%s)\n' "${CHECK_ID}" "${CHECK_NAME}"
	fi
done

BUCKETS="$(aws --region "${REGION}" s3api list-buckets --query "join(\`\\n\`, @.Buckets[].Name)" --output text)" || die 'failed to list S3 buckets\n'
printf '%d\t?\t%s\t-\tS3 buckets\n' "$(echo "${BUCKETS}" | wc -l)" "${REGION}" || die 'failed to write result for S3 buckets\n'

GATEWAY_VPC_ENDPOINTS="$(aws --region "${REGION}" ec2 describe-vpc-endpoints --query "join(\`\\n\`, @.VpcEndpoints[? @.VpcEndpointType == \`Gateway\`].VpcEndpointId)" --output text)" || die 'failed to list gateway VPC endpoints\n'
# https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html#limits_vpc
# Defaults to 20,  You cannot have more than 255 gateway endpoints per VPC.
GATEWAY_VPC_LIMIT='20?'  # per region
printf '%d\t%s\t%s\t-\tGateway VPC endpoints\n' "$(echo "${GATEWAY_VPC_ENDPOINTS}" | wc -l)" "${GATEWAY_VPC_LIMIT}" "${REGION}" || die 'failed to write result for gateway VPC endpoints\n'

NETWORK_INTERFACES="$(aws --region "${REGION}" ec2 describe-network-interfaces --query "join(\`\\n\`, @.NetworkInterfaces[].NetworkInterfaceId)" --output text)" || die 'failed to list EC2 network interfaces\n'
# https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html#limits_vpc
# This limit is the greater of either the default limit (350) or your
# On-Demand Instance limit multiplied by 5.  The default limit for
# On-Demand Instances is 20.  If your On-Demand Instance limit is
# below 70, the default limit of 350 applies.  To increase this limit,
# submit a request or increase your On-Demand Instance limit.
NETWORK_INTERFACE_LIMIT='350?'  # per region
printf '%d\t%s\t%s\t-\tEC2 network interfaces\n' "$(echo "${NETWORK_INTERFACES}" | wc -l)" "${NETWORK_INTERFACE_LIMIT}" "${REGION}" || die 'failed to write result for EC2 network interfaces\n'

NAT_GATEWAYS="$(aws --region "${REGION}" ec2 describe-nat-gateways --query "join(\`\\n\`, @.NatGateways[].NatGatewayId)" --output text)" || die 'failed to list EC2 NAT gateways\n'
# https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html#limits_vpc
NAT_GATEWAYS_LIMIT='5-per-zone?'  # per availability zone
printf '%d\t%s\t%s\t-\tEC2 NAT gateways\n' "$(echo "${NAT_GATEWAYS}" | wc -l)" "${NAT_GATEWAYS_LIMIT}" "${REGION}" || die 'failed to write result for EC2 NAT gateways\n'

SECURITY_GROUPS="$(aws --region "${REGION}" ec2 describe-security-groups --query "join(\`\\n\`, @.SecurityGroups[].GroupName)" --output text)" || die 'failed to list EC2 security groups\n'
# https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html#limits_vpc
SECURITY_GROUPS_LIMIT='2500?'  # per region
printf '%d\t%s\t%s\t-\tEC2 security groups\n' "$(echo "${SECURITY_GROUPS}" | wc -l)" "${SECURITY_GROUPS_LIMIT}" "${REGION}" || die 'failed to write result for EC2 security groups\n'

NETWORK_LOAD_BALANCERS="$(aws elbv2 describe-load-balancers --query "join(\`\\n\`, @.LoadBalancers[].LoadBalancerArn)" --output text)" || die 'failed to list network load balancers\n'
NETWORK_LOAD_BALANCER_LIMIT="$(aws elbv2 describe-account-limits --query "Limits[? @.Name == 'network-load-balancers'].Max" --output text)" || die 'failed to get network load balancer limit\n'
printf '%d\t%d\t%s\t-\tNetwork load balancers\n' "$(echo "${NETWORK_LOAD_BALANCERS}" | wc -l)" "${NETWORK_LOAD_BALANCER_LIMIT}" "${REGION}" || die 'failed to write result for network load balancers\n'
