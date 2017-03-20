# AWS with Terraform

*Prerequsities*

1. Ensure all [common prerequisites](../../../README.md#common-prerequisites) are met
1. Configure AWS credentials via environment variables. 
[See docs](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-environment)
1. Configure a region by setting `AWS_REGION` environment variable

## Usage

1. Ensure all *prerequsities* are met.
1. From the root of the repo, run `make PLATFORM=aws-asg CLUSTER=<cluster-name>`

To clean up run `make destroy PLATFORM=aws-asg CLUSTER=<cluster-name>`

[platform-lifecycle]: Documentation/platform-lifecycle.md
