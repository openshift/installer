# Tectonic Smoke Tests on Azure

This README documents how to run Tectonic smoke tests on Azure. This directory has the following configuration:
```
Azure
└── vars               # Terraform tfvars files for Azure smoke tests
    ├── Azure-exp.tfvars # Default Tectonic configuration + experimental features enabled
    ├── Azure.tfvars     # Default Tectonic configuration
    └── ...
```

All Tectonic smoke tests on Azure are written using the `RSpec`.
For help using this tool, go to: [http://www.rubydoc.info/gems/rspec-core/](http://www.rubydoc.info/gems/rspec-core/)

## Credentials
For now, the smoke tests rely on Azure service principal type of credentials being injected into the environment.
They also require a running SSH agent and a valid SSH key exported as `TF_VAR_tectonic_azure_ssh_key`.
We will try to make this more automated in the future.

Here is an example of how to set variables that carry Azure credentials:
```sh
export ARM_SUBSCRIPTION_ID="1901df3f-af3b-4a10-8874-58c8db1d5c8f"
export ARM_CLIENT_ID="66e0e96e-b83e-43f1-b1dc-148e7c9b8f36"
export ARM_CLIENT_SECRET="a5ec54fc-65a6-4a30-a3ff-80afc7213539"
export ARM_TENANT_ID="7fc8fb88-785b-4194-830d-acd348ea3b7d"
export TF_VAR_tectonic_azure_ssh_key="${HOME}/.ssh/id_rsa.pub"
```

## Running a test
For convenience, we provide a make target that can run either the entire test suite or just one / some tests.
To run the entire test suite for Azure, use this command:
```sh
export TEST_VARS=vars/Azure-exp.tfvars
make tests/smoke TEST='azure_*'
```
