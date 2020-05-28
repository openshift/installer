# Resource Requirements

A standard installation creates the following resources:

- 1 Folder
- 1 Tag Category
- 1 Tag
- Virtual machines:
  - 1 template
  - 1 temporary bootstrap node
  - 3 control-plane nodes
  - 3 compute machines

## Requirements

### Storage
With the above resources, a standard installation requires a minimum of 800 GB of storage.

### DHCP
Installation requires DHCP for the network. 

## Limits
Available resources vary between clusters. The number of possible clusters within a vCenter will be primarily limited by storage space, plus any limitations upon the number of the resources limited above. Day Zero resources not provisioned by the installer, such as IP addresses and networks, should also be considered when planning.