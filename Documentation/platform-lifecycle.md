## Platform Stability Lifecycle

Each platform is marked as either Pre-Alpha, Alpha, Beta, or Stable. This document outlines the rough criteria for each.
Each lifecycle phase is cumulative and assumes all previous phase criteria are included.


### Pre-Alpha

Initial platform assets are added to this repository and are undergoing active development.

*Feature Requirements*

- None

*Intended Usage*

- Developer workflows

*Suitable Environments*

- Developer sandbox

*Testing*

- Manual developer testing

*Packaging*

- None

*Support*

- None

*Updates Supported*

- No

*Docs*

- None


### Alpha

Can reliably produce minimally functioning clusters with manual testing.

*Feature Requirements*

- Kubernetes API works
- Authenticated Tectonic Console works
- All Kubernetes and Tectonic components function consistently

*Intended Usage*

- Developer workflows

*Suitable Environments*

- Developer sandbox
- Testing in select environments

*Testing*

- Smoke tests pass for one common configuration for all supported releases
- Smoke tests are integrated into testing framework and _can be_ run on pull requests

*Packaging*

- None

*Support*

- Informal

*Updates Supported*

- No

*Docs*

- Basic developer usage documentation
  - README
  - General usage
  - Platform caveats
- Variable documentation is auto-generated
- Example tfvars file is auto-generated


### Beta

*Feature Requirements*

- Best practices implemented for platform:
  - Network Security
  - DNS
  - Load Balancing
  - Generates HA / Multi-AZ infrastructure
- Cloud Provider enabled for the platform (if applicable)

*Intended Usage*

- Developer workflows
- CLI-based install from official release package

*Suitable Environments*

- Developer sandboxes
- Development environments
- Pre-production

*Testing*

- Smoke, Kubernetes conformance, and Tectonic integration tests are automated to run nightly on master for 3 most common configurations (using appropriate cloud-provider if applicable)
- Smoke, Kubernetes conformance, and Tectonic integration tests pass for 3 most common configurations for all supported releases (using appropriate cloud-provider if applicable)

*Packaging*

- Assets are packaged into each official Tectonic Installer release

*Support*

- Formal for paying customers
- Informal for non-paying customers

*Updates Supported*

- Yes

*Docs*

- User-facing documentation is committed, and covers all topics:
  - installation requirements
  - installation
  - troubleshooting
  - un-installation
  - all available customizations


### Stable


*Requirements*

- Adheres 100% to the [generic platform specification](generic-platform.md)
- All manifests are vetted and certified to not significantly diverge from other stable platform manifests, so much that cluster updates are not compromised
- All code follows style and testing guidelines
- (Optional) Tectonic Installer UI built for platform

*Intended Usage*

- Developer workflows
- CLI-based install from official release package
- Integration into CI systems
- (Optional) GUI-based install from official release package

*Suitable Environments*

- Developer sandboxes
- Development environments
- Pre-production
- Production

*Testing*

- Cluster upgrade tests pass for all supported releases

*Packaging*

- Assets are packaged into each official Tectonic Installer release
- (Optional) GUI Installer performs complete installation flow
- (Optional) GUI Installer is certified by CoreOS UX Team

*Support*

- Formal for paying customers
- Informal for non-paying customers

*Updates Supported*

- Yes

*Docs*

- Documentation is certified by CoreOS Documentation team
- Documentation is published on coreos.com

