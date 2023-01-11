# Network Type / Container Network Interface (CNI) plugin selection

### Post OpenShift 4.12 releases 

The default network type (CNI plugin) is OVNKubernetes for OpenShift v4.12. Prior to OpenShift v4.12, the default network type was OpenShiftSDN. OpenShiftSDN supports IPv4 only. For IPv6 or dual stack, a compatible NetworkType like OVNKubernetes must be used.

If you are using the install-config and agent-config yamls to generate the ZTP manifests that feeds into the agent ISO, the Network Type is set in install-config.yaml under Networking.NetworkType [https://docs.openshift.com/container-platform/4.11/installing/installing_bare_metal/installing-bare-metal-network-customizations.html]

If you are using ZTP manifests directly, the NetworkType is specified in the AgentClusterInstall under Spec.Networking.NetworkType.

### Pre Openshift 4.12 releases

The NetworkType always defaulted to OpenShiftSDN because the NetworkType specified either in the AgentClusterInstall or the InstallConfig wasn't being set in the cluster parameters being created. This meant that IPv6 did not work with the agent-installer prior to 4.12. This issue has been fixed in 4.12. 
