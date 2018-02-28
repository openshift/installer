# PLATFORMS is a dict matching {"name": "src", ...}.
PLATFORMS = {p: "//:platforms/" + p for p in ["azure", "gcp", "govcloud", "metal", "vmware"]}
PLATFORMS["openstack-neutron"] = "//:platforms/openstack/neutron"
PLATFORMS["aws"] = "//:steps"
