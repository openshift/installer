# Tectonic Smoke Tests on AWS

This README documents how to run Tectonic smoke tests on AWS. This directory has the following configuration:
```
aws
├── smoke.sh           # Utility script for planning and creating clusters and running tests
└── vars               # Terraform tfvars files for AWS smoke tests
    ├── aws-exp.tfvars # Default Tectonic configuration + experimental features enabled
    ├── aws.tfvars     # Default Tectonic configuration
    └── ...
```

All Tectonic smoke tests on AWS are run using the `smoke.sh` script found in this directory. For help using this tool, run:
```sh
./smoke.sh help
```

## Environment
To begin, verify that the following environment variables are set:

- `AWS_PROFILE` or alternatively `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`: These credentials are used by Terraform to spawn clusters on AWS.
- `TF_VAR_tectonic_pull_secret_path` and `TF_VAR_tectonic_license_path`: The local path to the pull secret and Tectonic license file.
- `TF_VAR_tectonic_aws_ssh_key`: The AWS ssh key pair which enables ssh'ing into the created machines using the `core` user.
  It must be present in AWS under "EC2 -> Network & Security -> Key Pairs".
- (optional) `BUILD_ID`: Any number >= 1. Based on this number the region will be selected of the deployed cluster.
  See the `REGIONS` variable under `smoke.sh` for details.
- (optional) `BRANCH_NAME`: The local branch name used as an infix for cluster names.
  A sensible value is `git rev-parse --abbrev-ref HEAD`.

Example:
```
$ export AWS_ACCESS_KEY_ID=AKIAIQ5TVFGQ7CKWD6IA
$ export AWS_SECRET_ACCESS_KEY_ID=rtp62V7H/JDY3cNBAs5vA0coaTou/OQbqJk96Hws
$ export TF_VAR_tectonic_license_path="/home/user/tectonic-license"
$ export TF_VAR_tectonic_pull_secret_path="/home/user/coreos-inc/pull-secret"
$ export TF_VAR_tectonic_aws_ssh_key="user"
```

## Assume Role
Smoke tests should be run with a role with limited privileges to ensure that the Tectonic Installer is as locked-down as possible.
The following steps create or update a role with restricted access and assume that role for testing purposes.
Edit the trust policy located at `Documentation/files/aws-sts-trust-policy.json` to include the ARN for the AWS user to be used for testing.
Then, define a role name, e.g. `tectonic-installer`, and use the `smoke.sh` script to create a role or update an existing one:
```sh
export ROLE_NAME=tectonic-installer
./smoke.sh set-role $ROLE_NAME ../../../Documentation/files/aws-policy.json ../../../Documentation/files/aws-sts-trust-policy.json
```

*Note*: the same command can be used to update the policies associated with an existing role.
Now that the role exists, assume the role:
```sh
source ./smoke.sh assume-role "$ROLE_NAME"
```

Note that this step can be skipped for local testing.

## Create and Test Cluster
Once the role has been assumed, select a Tectonic configuration to test, e.g. `vars/aws-exp.tfvars` and plan a cluster:
```sh
export TEST_VARS=vars/aws-exp.tfvars
./smoke.sh plan $TEST_VARS
```

This will create a terraform state directory in the project's top-level build directory, i.e. `/build/aws-exp-master-1012345678901`.

Continue by actually creating the cluster:
```sh
./smoke.sh create $TEST_VARS
```

Finally, test the cluster:
```sh
./smoke.sh test $TEST_VARS
=== RUN   TestCluster
=== RUN   TestCluster/APIAvailable
=== RUN   TestCluster/AllNodesRunning
=== RUN   TestCluster/GetLogs
=== RUN   TestCluster/AllPodsRunning
=== RUN   TestCluster/KillAPIServer
...
PASS
```

## Cleanup
Once all testing has concluded, clean up the AWS resources that were created:
```sh
./smoke.sh destroy $TEST_VARS
```

## Sanity test cheatsheet
To be able to ssh into the created machines, determine the generated cluster name and use the [AWS client](http://docs.aws.amazon.com/cli/latest/userguide/installing.html) to retrieve the public IP address or search for nodes having the cluster name via the AWS Web UI in "EC2 -> Instances":

```sh
$ ls build
aws-exp-master-1012345678901

$ export CLUSTER_NAME=aws-exp-master-1012345678901

$ aws autoscaling describe-auto-scaling-groups \
    | jq -r '.AutoScalingGroups[] | select(.AutoScalingGroupName | contains("'${CLUSTER_NAME}'")) | .Instances[].InstanceId' \
    | xargs aws ec2 describe-instances --instance-ids \
    | jq '.Reservations[].Instances[] | select(.PublicIpAddress != null) | .PublicIpAddress'
"52.15.184.15"

$ ssh -A core@52.15.184.15
```

Once connected to the master node, follow the [troubleshooting guide](../../../Documentation/troubleshooting/troubleshooting.md) for master, worker, and etcd nodes to investigate the following checklist:

- SSH connectivity to the master/worker/etcd nodes
- Successful start of all relevant installation service units on the corresponding nodes
- Successful login to the Tectonic Console
