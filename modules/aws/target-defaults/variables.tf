variable "profile" {
  type    = "string"
  default = ""

  description = <<EOF
(Optional) This is the AWS profile name as set in the shared credentials file.
This is passed through to the Terraform aws provider: https://www.terraform.io/docs/providers/aws/#profile
EOF
}

variable "region" {
  type = "string"

  description = <<EOF
This is the AWS region.
It is passed through to the Terraform aws provider: https://www.terraform.io/docs/providers/aws/#region
EOF
}

variable "role_arn" {
  type    = "string"
  default = ""

  description = <<EOF
(Optional) Name (full ARN) of IAM role to use to access AWS.
It is passed through to the Terraform aws provider: https://www.terraform.io/docs/providers/aws/#role_arn
EOF
}

variable "etcd_count" {
  type    = "string"
  default = "0"

  description = <<EOF
The number of etcd nodes to be created.
If set to zero, the count of etcd nodes will be determined automatically.
EOF
}
