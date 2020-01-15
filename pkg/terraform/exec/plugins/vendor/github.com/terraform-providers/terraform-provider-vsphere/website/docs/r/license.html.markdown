---
subcategory: "Administration"
layout: "vsphere"
page_title: "VMware vSphere: vsphere_license"
sidebar_current: "docs-vsphere-resource-admin-license"
description: |-
  Provides a VMware vSphere license resource. This can be used to add and remove license keys.
---

# vsphere\_license

Provides a VMware vSphere license resource. This can be used to add and remove license keys.

## Example Usage

```hcl
resource "vsphere_license" "licenseKey" {
  license_key = "452CQ-2EK54-K8742-00000-00000"

  labels {
    VpxClientLicenseLabel = "Hello World"
    Workflow = "Hello World"
  }
  
}
```

## Argument Reference

The following arguments are supported:

* `license_key` - (Required) The license key to add.
* `labels` - (Optional) A map of key/value pairs to be attached as labels (tags) to the license key.


## Attributes Reference

The following attributes are exported:

* `edition_key` - The product edition of the license key.
* `total` - Total number of units (example: CPUs) contained in the license.
* `used` - The number of units (example: CPUs) assigned to this license.
* `name` - The display name for the license.
