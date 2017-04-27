## Platform Stability

Each platform is marked as either pre-alpha, alpha, beta, or stable. This document outlines the rough criteria for each.

### Pre-Alpha

*Requirements*

- Initial static Terraform assets are added to this repo and underging active development
- No installer integration

*Packaging*

- None

*User Flow*

Development workflows only

### Alpha

- Manually tested and can reliably produce minimally functioning clusters
  - Kuberntes API works
  - Authenticated Tectonic Console works

*Packaging*

- Assets are packaged into the official Tectonic Installer tarball with every Tectonic release.

*User Flow*

1. Manually create `terraform.tfvars` file
1. Use static Terraform assets contained in installer tarball
1. Manually run Terraform


### Beta

*Requirements*

- Best practices for platform implemented
  - Network security
  - Automated testing results are published
  - DNS & Load Balancing
  - Generates HA / Multi-AZ infrastructure
- Automated testing results are published
  - Kubernetes e2e and conformance tests pass
  - Tectonic smoke tests work
- Cloud Provider enabled for the platform
- Tectonic automated updates work
- README documents all customizations
- Two documented platform users

*Packaging*

- Assets are packaged into the official Tectonic Installer tarball with every Tectonic release.
- (Optional) Installer UI is built to guide user through the configuration process.

*User Flow*

1. Use GUI installer to generate `terraform.tfvars` file, or manually create it.
1. Use static Terraform assets contained in installer tarball
1. Manually run terraform


### Stable

*Requirements*

- Automated tests pass for all supported releases
- Tectonic Installer UI for platform (optional)

*Packaging*

- (Optional) GUI Installer reads Terraform assets
- (Optional) GUI Installer backend runs equivalent of `terraform apply` on behalf of the user via Go client libraries
- (Optional) GUI shows status info and links to Console once complete

*User Flow*

GUI Only

1. Use GUI installer configure cluster
1. Click to provision cluster
1. GUI shows status info
1. Click to go directly to  Tectonic Console

Non-GUI

1. Manually create `terraform.tfvars` file
1. Use static Terraform assets contained in installer tarball
1. Manually run Terraform
