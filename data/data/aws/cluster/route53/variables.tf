variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = string
}

variable "base_domain" {
  description = "The base domain used for public records."
  type        = string
}

variable "vpc_id" {
  description = "The VPC used to create the private route53 zone."
}

variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "internal_zone" {
  type        = string
  description = "An existing hosted zone (zone ID) to use for the internal API."
}

variable "internal_zone_role" {
  type        = string
  default     = null
  description = "(optional) A role to assume when using an existing hosted zone from another account."
}

variable "api_external_lb_dns_name" {
  description = "External API's LB DNS name"
  type        = string
}

variable "api_external_lb_zone_id" {
  description = "External API's LB Zone ID"
  type        = string
}

variable "api_internal_lb_dns_name" {
  description = "Internal API's LB DNS name"
  type        = string
}

variable "api_internal_lb_zone_id" {
  description = "Internal API's LB Zone ID"
  type        = string
}

variable "publish_strategy" {
  type        = string
  description = <<EOF
The publishing strategy for endpoints like load balancers

Because of the issue https://github.com/hashicorp/terraform/issues/12570, the consumers cannot count 0/1
based on if api_external_lb_dns_name for example, which will be null when there is no external lb for API.
So publish_strategy serves an coordinated proxy for that decision.
EOF
}

variable "region" {
  type = string
  description = "The target AWS region for the cluster."
}

variable "custom_endpoints" {
  type = map(string)

  description = <<EOF
(optional) Custom AWS endpoints to override existing services.
Check - https://www.terraform.io/docs/providers/aws/guides/custom-service-endpoints.html

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}
