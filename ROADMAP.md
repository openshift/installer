## Roadmap

### Documentation

- Give examples of how to customize each of the platforms (DNS, loadbalancing, etc)

### Testing and Tooling

- Run [Kubernetes e2e tests](https://github.com/coreos-inc/tectonic-platform-sdk/issues/6) for each platform
- Build a tool to walk the Terraform graph and warn if cluster won't comply with [Generic Platform](https://github.com/coreos-inc/tectonic-platform-sdk/blob/master/Documentation/generic-platform.md)
- Create a spec for generic and platform specific Terraform Variable files
- Create a tool to verify Terraform Variable files

### Platforms
- Move Azure support to beta [see issues](https://github.com/coreos-inc/tectonic-platform-sdk/labels/platform%2Fazure)
- Move OpenStack support to beta [see issues](https://github.com/coreos-inc/tectonic-platform-sdk/labels/platform%2Fopenstack)
- Move VMware support to beta [see issues](https://github.com/coreos-inc/tectonic-platform-sdk/labels/platform%2Fvmware)
- Introduce alpha GCP support [see issues](https://github.com/coreos-inc/tectonic-platform-sdk/labels/platform%2Fgcp)
- Introduce alpha DO support [see issues](https://github.com/coreos-inc/tectonic-platform-sdk/labels/platform%2Fdigitalocean)
- Rearchitect bare-metal installer [see issues](https://github.com/coreos-inc/tectonic-platform-sdk/labels/platform%2Fbare-metal)
- GUI support for all of the above

## OS Support

- RHEL/CentOS
- Other Linux distributions

### Other Work

- Deploy with other self-hosted tools like kubeadm
