output "vpc_id" {
  value = ibm_is_vpc.ocp_vpc.id
}

output "vpc_subnet_id" {
  value = ibm_is_subnet.ocp_vpc_subnet.id
}

output "vpc_crn" {
  value = ibm_is_vpc.ocp_vpc.crn
}
