# Tectonic Smoke Tests on Azure

This README documents how to run Tectonic smoke tests on Azure. This directory has the following configuration:
```
Azure
├── smoke.sh           # Utility script for planning and creating clusters and running tests
└── vars               # Terraform tfvars files for Azure smoke tests
    ├── Azure-exp.tfvars # Default Tectonic configuration + experimental features enabled
    ├── Azure.tfvars     # Default Tectonic configuration
    └── ...
```

All Tectonic smoke tests on Azure are run using the `smoke.sh` script found in this directory. For help using this tool, run:
```sh
./smoke.sh help
```

## Credentials
For now, the smoke tests rely on Azure service principal type of credentials being injected into the environment.

Here is an example of how to set variables that carry Azure credentials:
```
export ARM_SUBSCRIPTION_ID="1901df3f-af3b-4a10-8874-58c8db1d5c8f"
export ARM_CLIENT_ID="66e0e96e-b83e-43f1-b1dc-148e7c9b8f36"
export ARM_CLIENT_SECRET="a5ec54fc-65a6-4a30-a3ff-80afc7213539"
export ARM_TENANT_ID="7fc8fb88-785b-4194-830d-acd348ea3b7d"
```

## Create and Test Cluster
Select a Tectonic configuration to test, e.g. `Azure-exp.tfvars` and plan a cluster:
```sh
export TEST_VARS=vars/Azure-exp.tfvars
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
Once all testing has concluded, clean up the Azure resources that were created:
```sh
./smoke.sh destroy $TEST_VARS
```
