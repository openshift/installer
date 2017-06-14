# Tectonic Tests

This directory contains all smoke and integration tests for Tectonic. Each subdirectory is further partitioned into directories for particular platforms as necessary. This directory should conform to the following directory hierarchy:

```
tests/
├── gui            # Integration tests for the installer user interface
├── integration    # Integration tests for Tectonic
└── smoke          # Smoke tests for all platforms
    ├── aws        # Smoke tests for AWS
    │   └── vars   # Terraform tfvars files for AWS smoke tests
    ├── azure      # Smoke tests for Azure
    │   └── vars   # Terraform tfvars files for Azure smoke tests
    ├── bare-metal # Smoke tests for bare-metal
    │   └── vars   # Terraform tfvars files for bare-metal smoke tests
    └── ...
```
