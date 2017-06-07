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

## Assume Role
Smoke tests should be run with a role with limited privileges to ensure that the Tectonic Installer is as locked-down as possible.
The following steps create or update a role with restricted access and assume that role for testing purposes.
To begin, verify that the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables are set in the current shell.
Now, edit the trust policy located at `Documentation/files/aws-sts-trust-policy.json` to include the ARN for the AWS user to be used for testing.
Then, define a role name, e.g. `tectonic-installer`, and use the `smoke.sh` script to create a role or update an existing one:
```sh
export ROLE_NAME=tectonic-installer
./smoke.sh set-role $ROLE_NAME ../../../Documentation/files/aws-policy.json ../../../Documentation/files/aws-sts-trust-policy.json
```

*Note*: the same command can be used to update the policies associated with an existing role.
Now that the role exists, assume the role:
```sh
eval "$(./smoke.sh assume-role "$ROLE_NAME")"
```

## Create and Test Cluster
Once the role has been assumed, select a Tectonic configuration to test, e.g. `aws-exp.tfvars` and plan a cluster:
```sh
export TEST_VARS=vars/aws-exp.tfvars
./smoke.sh plan $TEST_VARS
```

Continue by actually creating the cluster:
```sh
./smoke.sh create $TEST_VARS
```

Finally, test the cluster:
```sh
./smoke.sh test $TEST_VARS
```

## Cleanup
Once all testing has concluded, clean up the AWS resources that were created:
```sh
./smoke.sh destroy $TEST_VARS
```
