package gcp

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types"
	dnstypes "github.com/openshift/installer/pkg/types/dns"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// PermissionGroup is the group of permissions needed by cluster creation, operation, or teardown.
type PermissionGroup string

const (
	// PermissionCreateBase is a base set of permissions required in all installs where the installer creates resources.
	PermissionCreateBase PermissionGroup = "create-base"

	// PermissionDeleteBase is a base set of permissions required in all installs where the installer destroys resources.
	PermissionDeleteBase PermissionGroup = "delete-base"

	// PermissionDeleteNetwork is the set of permissions required to delete network components.
	PermissionDeleteNetwork PermissionGroup = "delete-network"

	// PermissionDeleteManagedFirewalls is the set of permissions required to delete firewall components.
	PermissionDeleteManagedFirewalls PermissionGroup = "delete-managed-firewalls"

	// PermissionDeleteServiceAccounts is the set of permissions required to delete Service Account resources.
	PermissionDeleteServiceAccounts PermissionGroup = "delete-service-accounts"

	// PermissionCreateFirewall is the set of permissions required to create firewall resources.
	PermissionCreateFirewall PermissionGroup = "create-firewall"

	// PermissionCreateNetwork is the set of permissions required to create network resources.
	PermissionCreateNetwork PermissionGroup = "create-network"

	// PermissionCreateDNS is the set of permissions required to create DNS resources.
	PermissionCreateDNS PermissionGroup = "create-dns"

	// PermissionCreateServiceAccounts is the set of permissions required to create Service Account resources.
	PermissionCreateServiceAccounts PermissionGroup = "create-service-accounts"

	// PermissionCreateKMSSupport is the set of permissions required to validate a user-supplied KMS
	// key. The installer neither creates nor deletes KMS resources -- GCP does not permit deleting
	// a key ring or a crypto key at all -- so these are read-only.
	PermissionCreateKMSSupport PermissionGroup = "create-kms-support"
)

var permissions = map[PermissionGroup][]string{
	// Base set of permissions for cluster creation
	PermissionCreateBase: {
		// Compute Instances & Disks
		"compute.disks.create",
		"compute.disks.get",
		"compute.disks.list",
		"compute.disks.setLabels",
		"compute.diskTypes.get",
		"compute.instanceGroups.create",
		"compute.instanceGroups.delete",
		"compute.instanceGroups.get",
		"compute.instanceGroups.list",
		"compute.instanceGroups.update",
		"compute.instanceGroups.use",
		"compute.instances.create",
		"compute.instances.delete",
		"compute.instances.get",
		"compute.instances.list",
		"compute.instances.setLabels",
		"compute.instances.setMetadata",
		"compute.instances.setServiceAccount",
		"compute.instances.setTags",
		"compute.instances.use",
		"compute.machineTypes.get",
		"compute.machineTypes.list",
		"compute.images.list",
		// Load Balancers
		"compute.backendServices.create",
		"compute.backendServices.get",
		"compute.backendServices.list",
		"compute.backendServices.update",
		"compute.backendServices.use",
		"compute.regionBackendServices.create",
		"compute.regionBackendServices.get",
		"compute.regionBackendServices.list",
		"compute.regionBackendServices.update",
		"compute.regionBackendServices.use",
		"compute.targetPools.addInstance",
		"compute.targetPools.create",
		"compute.targetPools.get",
		"compute.targetPools.list",
		"compute.targetPools.removeInstance",
		"compute.targetPools.use",
		"compute.targetTcpProxies.create",
		"compute.targetTcpProxies.get",
		"compute.targetTcpProxies.use",
		// Health Checks
		"compute.healthChecks.create",
		"compute.healthChecks.get",
		"compute.healthChecks.list",
		"compute.healthChecks.update",
		"compute.healthChecks.useReadOnly",
		"compute.httpHealthChecks.create",
		"compute.httpHealthChecks.get",
		"compute.httpHealthChecks.list",
		"compute.httpHealthChecks.useReadOnly",
		"compute.regionHealthChecks.create",
		"compute.regionHealthChecks.get",
		"compute.regionHealthChecks.update",
		"compute.regionHealthChecks.useReadOnly",
		// Storage
		"storage.buckets.create",
		"storage.buckets.delete",
		"storage.buckets.get",
		"storage.buckets.list",
		"storage.objects.create",
		"storage.objects.delete",
		"storage.objects.get",
		"storage.objects.list",
		// Zones, Regions & Operations
		"compute.globalOperations.get",
		"compute.regionOperations.get",
		"compute.regions.get",
		"compute.regions.list",
		"compute.zoneOperations.get",
		"compute.zones.get",
		"compute.zones.list",
		// Monitoring & Quotas
		"monitoring.timeSeries.list",
		"serviceusage.quotas.get",
		"serviceusage.services.list",
		// IAM Roles
		"iam.roles.create",
		"iam.roles.get",
		"iam.roles.update",
		// The following permission is required for signing the bucket
		// for bootstrapping
		"iam.serviceAccounts.signBlob",
		// The following permission is required for bootstrap gather
		// and cannot be determined from an install-config
		"compute.instances.getSerialPortOutput",
	},
	PermissionCreateFirewall: {
		"compute.firewalls.create",
		"compute.firewalls.delete",
	},
	PermissionCreateNetwork: {
		"compute.addresses.create",
		"compute.addresses.createInternal",
		"compute.addresses.delete",
		"compute.addresses.get",
		"compute.addresses.list",
		"compute.addresses.use",
		"compute.addresses.useInternal",
		"compute.firewalls.get",
		"compute.firewalls.list",
		"compute.forwardingRules.create",
		"compute.forwardingRules.get",
		"compute.forwardingRules.list",
		"compute.forwardingRules.setLabels",
		"compute.globalAddresses.create",
		"compute.globalAddresses.get",
		"compute.globalAddresses.use",
		"compute.globalForwardingRules.create",
		"compute.globalForwardingRules.get",
		"compute.globalForwardingRules.setLabels",
		"compute.networks.create",
		"compute.networks.get",
		"compute.networks.list",
		"compute.networks.updatePolicy",
		"compute.networks.use",
		"compute.routers.create",
		"compute.routers.get",
		"compute.routers.list",
		"compute.routers.update",
		"compute.routes.list",
		"compute.subnetworks.create",
		"compute.subnetworks.get",
		"compute.subnetworks.list",
		"compute.subnetworks.use",
		"compute.subnetworks.useExternalIp",
	},
	PermissionCreateDNS: {
		"dns.managedZones.get",
		"dns.managedZones.list",
		"dns.changes.create",
		"dns.changes.get",
		"dns.networks.bindPrivateDNSZone",
		"dns.resourceRecordSets.create",
		"dns.resourceRecordSets.list",
		"dns.managedZones.create",
	},
	PermissionCreateServiceAccounts: {
		"iam.serviceAccountKeys.create",
		"iam.serviceAccountKeys.delete",
		"iam.serviceAccountKeys.get",
		"iam.serviceAccountKeys.list",
		"iam.serviceAccounts.actAs",
		"iam.serviceAccounts.create",
		"iam.serviceAccounts.delete",
		"iam.serviceAccounts.get",
		"iam.serviceAccounts.list",
		"resourcemanager.projects.get",
		"resourcemanager.projects.getIamPolicy",
		// The installer grants roles to the service accounts it creates
		// (AddServiceAccountRoles -> Projects.SetIamPolicy), and it does so in
		// every credentials mode.
		"resourcemanager.projects.setIamPolicy",
	},
	PermissionDeleteBase: {
		// Compute
		"compute.disks.delete",
		"compute.disks.list",
		"compute.instanceGroups.delete",
		"compute.instanceGroups.list",
		"compute.instances.delete",
		"compute.instances.list",
		"compute.instances.stop",
		"compute.machineTypes.list",
		// Load Balancers
		"compute.backendServices.delete",
		"compute.backendServices.list",
		"compute.regionBackendServices.delete",
		"compute.regionBackendServices.list",
		"compute.targetPools.delete",
		"compute.targetPools.list",
		"compute.targetTcpProxies.delete",
		"compute.targetTcpProxies.list",
		// Health Checks
		"compute.healthChecks.delete",
		"compute.healthChecks.list",
		"compute.httpHealthChecks.delete",
		"compute.httpHealthChecks.list",
		"compute.regionHealthChecks.delete",
		"compute.regionHealthChecks.list",
		// Storage
		"storage.buckets.delete",
		"storage.buckets.getIamPolicy",
		"storage.buckets.list",
		"storage.objects.delete",
		"storage.objects.list",
		// Image
		"compute.images.list",
		// DNS
		"dns.changes.create",
		"dns.managedZones.delete",
		"dns.managedZones.get",
		"dns.managedZones.list",
		"dns.resourceRecordSets.delete",
		"dns.resourceRecordSets.list",
	},
	PermissionDeleteManagedFirewalls: {
		"compute.firewalls.delete",
	},
	PermissionDeleteNetwork: {
		"compute.addresses.delete",
		"compute.addresses.deleteInternal",
		"compute.addresses.list",
		"compute.addresses.setLabels",
		"compute.firewalls.list",
		"compute.forwardingRules.delete",
		"compute.forwardingRules.list",
		"compute.globalAddresses.delete",
		"compute.globalAddresses.list",
		"compute.globalForwardingRules.delete",
		"compute.globalForwardingRules.list",
		"compute.networks.delete",
		"compute.networks.list",
		"compute.networks.updatePolicy",
		"compute.routers.delete",
		"compute.routers.list",
		"compute.routes.list",
		"compute.subnetworks.delete",
		"compute.subnetworks.list",
	},
	PermissionDeleteServiceAccounts: {
		"iam.serviceAccounts.delete",
		"iam.serviceAccounts.get",
		"iam.serviceAccounts.list",
		"resourcemanager.projects.getIamPolicy",
	},
	PermissionCreateKMSSupport: {
		"cloudkms.keyRings.list",
		"cloudkms.cryptoKeys.getIamPolicy",
	},
}

const (
	// CreateFirewallPermission is the permission to create firewall rules in the google cloud provider.
	CreateFirewallPermission = "compute.firewalls.create" // required to create GCP firewall rules

	// DeleteFirewallPermission is the permission to delete firewall rules in the google cloud provider.
	DeleteFirewallPermission = "compute.firewalls.delete"

	// UpdateNetworksPermission is the permission to update networks and the network resources in the google cloud provider.
	UpdateNetworksPermission = "compute.networks.updatePolicy"
)

// CustomIAMRole contains the data to create the custom IAM role file.
type CustomIAMRole struct {
	Title               string   `json:"title,omitempty"`
	Description         string   `json:"description,omitempty"`
	Stage               string   `json:"stage,omitempty"`
	IncludedPermissions []string `json:"includedPermissions"`
}

// HasPermission determines if the permission exists for the service account in the project.
func HasPermission(ctx context.Context, projectID string, permissions []string, endpoint *gcptypes.PSCEndpoint) (bool, error) {
	client, err := NewClient(ctx, endpoint)
	if err != nil {
		return false, fmt.Errorf("failed to create client permission check: %w", err)
	}

	foundPermissions, err := client.GetProjectPermissions(ctx, projectID, permissions)
	if err != nil {
		return false, fmt.Errorf("failed to find project permissions: %w", err)
	}

	permissionsValid := true
	for _, permission := range permissions {
		if hasPermission := foundPermissions.Has(permission); !hasPermission {
			logrus.Debugf("permission %s not found", permission)
			permissionsValid = false
		}
	}

	return permissionsValid, nil
}

// Permissions returns the list of permissions associated with `group`.
func Permissions(group PermissionGroup) ([]string, error) {
	groupPerms, ok := permissions[group]
	if !ok {
		return nil, fmt.Errorf("unable to access permissions group %s", group)
	}
	return groupPerms, nil
}

// hasKMSKey reports whether the machine pool encrypts its OS disk with a customer-supplied KMS key.
func hasKMSKey(mpool *gcptypes.MachinePool) bool {
	return mpool != nil && mpool.OSDisk.EncryptionKey != nil && mpool.OSDisk.EncryptionKey.KMSKey != nil
}

// RequiredPermissionGroups returns the permissions the installer itself exercises for a given
// cluster configuration.
//
// The credentials mode deliberately plays no part in the result. The installer makes the same GCP
// API calls in mint, passthrough and manual mode: the only mode-dependent behaviour in the install
// flow is whether the root credentials secret manifest is written to disk (see
// pkg/asset/manifests/openshift.go), which requires no permission. What does vary by mode is who
// consumes the credentials afterwards -- the cloud-credential-operator mints per-component service
// accounts in mint mode, and every operator receives the credentials verbatim in passthrough mode
// -- but those permissions are exercised in-cluster and so are out of scope here. See
// PermissionPassthroughMode.
func RequiredPermissionGroups(ic *types.InstallConfig) []PermissionGroup {
	permissionGroups := []PermissionGroup{
		PermissionCreateBase,
		PermissionDeleteBase,
	}

	if ic.GCP.UserProvisionedDNS != dnstypes.UserProvisionedDNSEnabled {
		permissionGroups = append(permissionGroups, PermissionCreateDNS)
	}

	usingExistingVPC := len(ic.GCP.Network) != 0 && (len(ic.GCP.ControlPlaneSubnet) != 0 && len(ic.GCP.ComputeSubnet) != 0)
	if !usingExistingVPC {
		permissionGroups = append(permissionGroups, PermissionCreateNetwork, PermissionDeleteNetwork)
		if ic.GCP.FirewallRulesManagement == gcptypes.ManagedFirewallRules {
			permissionGroups = append(permissionGroups, PermissionCreateFirewall, PermissionDeleteManagedFirewalls)
		}
	}

	serviceAccountPermissionsNeeded := ic.ControlPlane == nil || ic.ControlPlane.Platform.GCP == nil || ic.ControlPlane.Platform.GCP.ServiceAccount == ""

	// Only check compute if control plane turned up nothing.
	if !serviceAccountPermissionsNeeded {
		for _, compute := range ic.Compute {
			if compute.Platform.GCP == nil || compute.Platform.GCP.ServiceAccount == "" {
				serviceAccountPermissionsNeeded = true
				break
			}
		}
	}

	// Only check default machine platform if the compute and control plane turned up nothing.
	if !serviceAccountPermissionsNeeded {
		if ic.GCP.DefaultMachinePlatform == nil || ic.GCP.DefaultMachinePlatform.ServiceAccount == "" {
			serviceAccountPermissionsNeeded = true
		}
	}

	if serviceAccountPermissionsNeeded {
		permissionGroups = append(permissionGroups, PermissionCreateServiceAccounts, PermissionDeleteServiceAccounts)
	}

	// ValidateKMSKeys checks a key wherever it is configured -- on the control plane pool, on any
	// compute pool, or on the default machine platform -- so the permissions have to follow the
	// same rule. Checking only the default machine platform leaves a per-pool key unvalidatable.
	kmsKeyPermissionsNeeded := ic.GCP.DefaultMachinePlatform != nil && hasKMSKey(ic.GCP.DefaultMachinePlatform)

	if !kmsKeyPermissionsNeeded && ic.ControlPlane != nil {
		kmsKeyPermissionsNeeded = hasKMSKey(ic.ControlPlane.Platform.GCP)
	}

	if !kmsKeyPermissionsNeeded {
		for _, compute := range ic.Compute {
			if hasKMSKey(compute.Platform.GCP) {
				kmsKeyPermissionsNeeded = true
				break
			}
		}
	}

	if kmsKeyPermissionsNeeded {
		permissionGroups = append(permissionGroups, PermissionCreateKMSSupport)
	}

	if len(ic.GCP.NetworkProjectID) > 0 && serviceAccountPermissionsNeeded {
		// Shared VPC installs grant roles to the machine service accounts in the host project as
		// well, so resourcemanager.projects.setIamPolicy is needed there too. A project-scoped
		// custom role cannot express that, so surface it instead.
		logrus.Infof("Shared VPC install: the host project %s also requires resourcemanager.projects.setIamPolicy and resourcemanager.projects.getIamPolicy", ic.GCP.NetworkProjectID)
	}

	return permissionGroups
}
