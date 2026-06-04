# Azure SDK V2 Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate all legacy Azure SDK imports (profiles/services packages) to V2 SDK (resourcemanager packages) in the OpenShift installer.

**Architecture:** The migration converts legacy autorest-based clients to modern azcore-based clients. The `Session` type already provides V2 authentication (`TokenCreds`, `CloudConfig`). Each service client will be converted to use client factories and the new async polling pattern.

**Tech Stack:** Go, Azure SDK for Go v2 (`github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/*`), azcore, azidentity

---

## File Structure

### Files to Modify

| File | Current State | Target State |
|------|--------------|--------------|
| `pkg/asset/installconfig/azure/dns.go` | Legacy `profiles/2018-03-01/dns` | V2 `resourcemanager/dns/armdns` |
| `pkg/asset/installconfig/azure/client.go` | Mixed legacy profiles + V2 | Fully V2 SDK |
| `pkg/asset/installconfig/azure/validation.go` | Legacy profiles types | V2 types |
| `pkg/asset/installconfig/azure/mock/azureclient_generated.go` | Auto-generated | Re-generate after API changes |
| `pkg/destroy/azure/azure.go` | Mixed legacy + V2 | Fully V2 SDK |

### Reference Files (Already V2)

| File | Patterns to Copy |
|------|-----------------|
| `pkg/infrastructure/azure/dns.go` | DNS record creation with armdns |
| `pkg/infrastructure/azure/compute.go` | Client factory, async polling |
| `pkg/infrastructure/azure/network.go` | armnetwork/v2 usage |
| `pkg/asset/installconfig/azure/session.go` | TokenCreds, CloudConfig setup |

---

## Task 1: Migrate dns.go to V2 SDK

**Files:**
- Modify: `pkg/asset/installconfig/azure/dns.go:1-153`

### Step 1.1: Update imports

- [ ] **Step 1.1.1: Replace legacy imports with V2**

Replace:
```go
import (
	"context"
	"errors"
	"fmt"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	azdns "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/dns/mgmt/dns"
	"github.com/Azure/go-autorest/autorest/to"
)
```

With:
```go
import (
	"context"
	"errors"
	"fmt"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
)
```

### Step 1.2: Update ZonesClient struct

- [ ] **Step 1.2.1: Update ZonesClient to use V2 client**

Replace:
```go
// ZonesClient wraps the azure ZonesClient internal
type ZonesClient struct {
	azureClient azdns.ZonesClient
}
```

With:
```go
// ZonesClient wraps the azure ZonesClient internal
type ZonesClient struct {
	azureClient *armdns.ZonesClient
}
```

### Step 1.3: Update RecordSetsClient struct

- [ ] **Step 1.3.1: Update RecordSetsClient to use V2 client**

Replace:
```go
// RecordSetsClient wraps the azure RecordSetsClient internal
type RecordSetsClient struct {
	azureClient azdns.RecordSetsClient
}
```

With:
```go
// RecordSetsClient wraps the azure RecordSetsClient internal
type RecordSetsClient struct {
	azureClient *armdns.RecordSetsClient
}
```

### Step 1.4: Update GetDNSRecordSet method signature

- [ ] **Step 1.4.1: Update method to use V2 types**

Replace:
```go
// GetDNSRecordSet gets a record set for the zone identified by publicZoneID
func (config DNSConfig) GetDNSRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType azdns.RecordType) (*azdns.RecordSet, error) {
	recordsetsClient := newRecordSetsClient(config.session)
	return recordsetsClient.GetRecordSet(rgName, zoneName, relativeRecordSetName, recordType)
}
```

With:
```go
// GetDNSRecordSet gets a record set for the zone identified by publicZoneID
func (config DNSConfig) GetDNSRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType armdns.RecordType) (*armdns.RecordSet, error) {
	recordsetsClient, err := newRecordSetsClient(config.session)
	if err != nil {
		return nil, err
	}
	return recordsetsClient.GetRecordSet(rgName, zoneName, relativeRecordSetName, recordType)
}
```

### Step 1.5: Update newZonesClient function

- [ ] **Step 1.5.1: Replace legacy client creation with V2**

Replace:
```go
func newZonesClient(session *Session) ZonesGetter {
	azureClient := azdns.NewZonesClientWithBaseURI(session.Environment.ResourceManagerEndpoint, session.Credentials.SubscriptionID)
	azureClient.Authorizer = session.Authorizer
	return &ZonesClient{azureClient: azureClient}
}
```

With:
```go
func newZonesClient(session *Session) (ZonesGetter, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: session.CloudConfig,
		},
	}
	azureClient, err := armdns.NewZonesClient(session.Credentials.SubscriptionID, session.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS zones client: %w", err)
	}
	return &ZonesClient{azureClient: azureClient}, nil
}
```

### Step 1.6: Update newRecordSetsClient function

- [ ] **Step 1.6.1: Replace legacy client creation with V2**

Replace:
```go
func newRecordSetsClient(session *Session) *RecordSetsClient {
	azureClient := azdns.NewRecordSetsClientWithBaseURI(session.Environment.ResourceManagerEndpoint, session.Credentials.SubscriptionID)
	azureClient.Authorizer = session.Authorizer
	return &RecordSetsClient{azureClient: azureClient}
}
```

With:
```go
func newRecordSetsClient(session *Session) (*RecordSetsClient, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: session.CloudConfig,
		},
	}
	azureClient, err := armdns.NewRecordSetsClient(session.Credentials.SubscriptionID, session.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS record sets client: %w", err)
	}
	return &RecordSetsClient{azureClient: azureClient}, nil
}
```

### Step 1.7: Update GetAllPublicZones method

- [ ] **Step 1.7.1: Replace legacy pagination with V2 pager**

Replace:
```go
// GetAllPublicZones get all public zones from the current subscription
func (client *ZonesClient) GetAllPublicZones() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	allZones := map[string]string{}
	for zonesPage, err := client.azureClient.List(ctx, to.Int32Ptr(100)); zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}
		for _, zone := range zonesPage.Values() {
			allZones[to.String(zone.Name)] = to.String(zone.ID)
		}
	}
	return allZones, nil
}
```

With:
```go
// GetAllPublicZones get all public zones from the current subscription
func (client *ZonesClient) GetAllPublicZones() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	allZones := map[string]string{}
	pager := client.azureClient.NewListPager(&armdns.ZonesClientListOptions{Top: to.Ptr(int32(100))})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, zone := range page.Value {
			if zone.Name != nil && zone.ID != nil {
				allZones[*zone.Name] = *zone.ID
			}
		}
	}
	return allZones, nil
}
```

### Step 1.8: Update GetRecordSet method

- [ ] **Step 1.8.1: Replace legacy Get with V2 Get**

Replace:
```go
// GetRecordSet gets an Azure DNS recordset by zone, name and recordset type
func (client *RecordSetsClient) GetRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType azdns.RecordType) (*azdns.RecordSet, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	recordset, err := client.azureClient.Get(ctx, rgName, zoneName, relativeRecordSetName, recordType)
	if err != nil {
		return nil, err
	}

	return &recordset, nil
}
```

With:
```go
// GetRecordSet gets an Azure DNS recordset by zone, name and recordset type
func (client *RecordSetsClient) GetRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType armdns.RecordType) (*armdns.RecordSet, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	resp, err := client.azureClient.Get(ctx, rgName, zoneName, relativeRecordSetName, recordType, nil)
	if err != nil {
		return nil, err
	}

	return &resp.RecordSet, nil
}
```

### Step 1.9: Update GetDNSZone to handle errors

- [ ] **Step 1.9.1: Update GetDNSZone to handle new error return**

Replace:
```go
// GetDNSZone returns a DNS zone selected by survey
func (config DNSConfig) GetDNSZone() (*Zone, error) {
	//call azure api using the session to retrieve available base domain
	zonesClient := newZonesClient(config.session)
	allZones, _ := zonesClient.GetAllPublicZones()
```

With:
```go
// GetDNSZone returns a DNS zone selected by survey
func (config DNSConfig) GetDNSZone() (*Zone, error) {
	//call azure api using the session to retrieve available base domain
	zonesClient, err := newZonesClient(config.session)
	if err != nil {
		return nil, fmt.Errorf("failed to create zones client: %w", err)
	}
	allZones, _ := zonesClient.GetAllPublicZones()
```

### Step 1.10: Add missing import

- [ ] **Step 1.10.1: Add azcore import**

Add to imports:
```go
"github.com/Azure/azure-sdk-for-go/sdk/azcore"
```

- [ ] **Step 1.11: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 1.12: Run tests**

Run: `go test ./pkg/asset/installconfig/azure/... -v -short`
Expected: Tests pass (dns.go has no direct tests, validation_test.go may need updates)

- [ ] **Step 1.13: Commit**

```bash
git add pkg/asset/installconfig/azure/dns.go
git commit -m "$(cat <<'EOF'
azure: migrate dns.go to V2 SDK

Convert dns.go from legacy Azure SDK profiles to V2 SDK:
- Replace azdns (profiles/2018-03-01/dns) with armdns
- Replace autorest/to with azcore/to
- Update ZonesClient and RecordSetsClient to use V2 clients
- Update pagination to use V2 pager pattern
- Update error handling for client creation

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: Update validation.go DNS RecordType references

**Files:**
- Modify: `pkg/asset/installconfig/azure/validation.go:1-50` (imports and DNS types)

### Step 2.1: Update imports

- [ ] **Step 2.1.1: Replace DNS import**

Replace:
```go
azdns "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/dns/mgmt/dns"
```

With:
```go
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
```

### Step 2.2: Update ValidatePublicDNS function

- [ ] **Step 2.2.1: Update RecordType references**

Replace:
```go
// Look for an existing CNAME first
rs, err := azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.CNAME)
if err == nil && rs.CnameRecord != nil {
	return fmt.Errorf(fmtStr, zoneName, azdns.CNAME, clusterName)
}

// Look for an A record
rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.A)
if err == nil && rs.ARecords != nil && len(*rs.ARecords) > 0 {
	return fmt.Errorf(fmtStr, zoneName, azdns.A, clusterName)
}

// Look for an AAAA record
rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.AAAA)
if err == nil && rs.AaaaRecords != nil && len(*rs.AaaaRecords) > 0 {
	return fmt.Errorf(fmtStr, zoneName, azdns.AAAA, clusterName)
}
```

With:
```go
// Look for an existing CNAME first
rs, err := azureDNS.GetDNSRecordSet(rgName, zoneName, record, armdns.RecordTypeCNAME)
if err == nil && rs.Properties != nil && rs.Properties.CnameRecord != nil {
	return fmt.Errorf(fmtStr, zoneName, armdns.RecordTypeCNAME, clusterName)
}

// Look for an A record
rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, armdns.RecordTypeA)
if err == nil && rs.Properties != nil && rs.Properties.ARecords != nil && len(rs.Properties.ARecords) > 0 {
	return fmt.Errorf(fmtStr, zoneName, armdns.RecordTypeA, clusterName)
}

// Look for an AAAA record
rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, armdns.RecordTypeAAAA)
if err == nil && rs.Properties != nil && rs.Properties.AaaaRecords != nil && len(rs.Properties.AaaaRecords) > 0 {
	return fmt.Errorf(fmtStr, zoneName, armdns.RecordTypeAAAA, clusterName)
}
```

- [ ] **Step 2.3: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 2.4: Run tests**

Run: `go test ./pkg/asset/installconfig/azure/... -v -short`
Expected: Tests pass

- [ ] **Step 2.5: Commit**

```bash
git add pkg/asset/installconfig/azure/validation.go
git commit -m "$(cat <<'EOF'
azure: migrate validation.go DNS types to V2 SDK

Update ValidatePublicDNS to use V2 DNS types:
- Replace azdns.CNAME/A/AAAA with armdns.RecordType* constants
- Update RecordSet property access to use V2 struct layout

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: Migrate client.go Network clients to V2 SDK

**Files:**
- Modify: `pkg/asset/installconfig/azure/client.go:1-150` (network-related methods)

### Step 3.1: Update imports - remove network profiles

- [ ] **Step 3.1.1: Replace aznetwork profile import**

Replace:
```go
aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
```

With:
```go
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
```

Note: `armnetwork/v2` is already imported, so just remove the profile import.

### Step 3.2: Update API interface - network types

- [ ] **Step 3.2.1: Update API interface method signatures**

Replace:
```go
GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*aznetwork.VirtualNetwork, error)
CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*aznetwork.IPAddressAvailabilityResult, error)
GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*aznetwork.Subnet, error)
```

With:
```go
GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*armnetwork.VirtualNetwork, error)
CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*armnetwork.IPAddressAvailabilityResult, error)
GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*armnetwork.Subnet, error)
GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subnet string) (*armnetwork.Subnet, error)
```

### Step 3.3: Update GetVirtualNetwork method

- [ ] **Step 3.3.1: Replace legacy client with V2 client**

Replace:
```go
// GetVirtualNetwork gets an Azure virtual network by name
func (c *Client) GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*aznetwork.VirtualNetwork, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient(ctx)
	if err != nil {
		return nil, err
	}

	vnet, err := vnetClient.Get(ctx, resourceGroupName, virtualNetwork, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual network %s: %w", virtualNetwork, err)
	}

	return &vnet, nil
}
```

With:
```go
// GetVirtualNetwork gets an Azure virtual network by name
func (c *Client) GetVirtualNetwork(ctx context.Context, resourceGroupName, virtualNetwork string) (*armnetwork.VirtualNetwork, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient()
	if err != nil {
		return nil, err
	}

	resp, err := vnetClient.Get(ctx, resourceGroupName, virtualNetwork, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual network %s: %w", virtualNetwork, err)
	}

	return &resp.VirtualNetwork, nil
}
```

### Step 3.4: Update CheckIPAddressAvailability method

- [ ] **Step 3.4.1: Replace legacy client with V2 client**

Replace:
```go
// CheckIPAddressAvailability checks availability of an IP address in an Azure virtual network.
func (c *Client) CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*aznetwork.IPAddressAvailabilityResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient(ctx)
	if err != nil {
		return nil, err
	}

	availability, err := vnetClient.CheckIPAddressAvailability(ctx, resourceGroupName, virtualNetwork, ipAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure ip availability: %w", err)
	}

	return &availability, nil
}
```

With:
```go
// CheckIPAddressAvailability checks availability of an IP address in an Azure virtual network.
func (c *Client) CheckIPAddressAvailability(ctx context.Context, resourceGroupName, virtualNetwork, ipAddr string) (*armnetwork.IPAddressAvailabilityResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vnetClient, err := c.getVirtualNetworksClient()
	if err != nil {
		return nil, err
	}

	resp, err := vnetClient.CheckIPAddressAvailability(ctx, resourceGroupName, virtualNetwork, ipAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure ip availability: %w", err)
	}

	return &resp.IPAddressAvailabilityResult, nil
}
```

### Step 3.5: Update getSubnet method

- [ ] **Step 3.5.1: Replace legacy client with V2 client**

Replace:
```go
// getSubnet gets an Azure subnet by name
func (c *Client) getSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*aznetwork.Subnet, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	subnetsClient, err := c.getSubnetsClient(ctx)
	if err != nil {
		return nil, err
	}

	subnet, err := subnetsClient.Get(ctx, resourceGroupName, virtualNetwork, subNetwork, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get subnet %s: %w", subNetwork, err)
	}

	return &subnet, nil
}
```

With:
```go
// getSubnet gets an Azure subnet by name
func (c *Client) getSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*armnetwork.Subnet, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	subnetsClient, err := c.getSubnetsClient()
	if err != nil {
		return nil, err
	}

	resp, err := subnetsClient.Get(ctx, resourceGroupName, virtualNetwork, subNetwork, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subnet %s: %w", subNetwork, err)
	}

	return &resp.Subnet, nil
}
```

### Step 3.6: Update GetComputeSubnet and GetControlPlaneSubnet

- [ ] **Step 3.6.1: Update return types**

Replace:
```go
// GetComputeSubnet gets the Azure compute subnet
func (c *Client) GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*aznetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}

// GetControlPlaneSubnet gets the Azure control plane subnet
func (c *Client) GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*aznetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}
```

With:
```go
// GetComputeSubnet gets the Azure compute subnet
func (c *Client) GetComputeSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*armnetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}

// GetControlPlaneSubnet gets the Azure control plane subnet
func (c *Client) GetControlPlaneSubnet(ctx context.Context, resourceGroupName, virtualNetwork, subNetwork string) (*armnetwork.Subnet, error) {
	return c.getSubnet(ctx, resourceGroupName, virtualNetwork, subNetwork)
}
```

### Step 3.7: Update getVirtualNetworksClient

- [ ] **Step 3.7.1: Replace legacy client factory with V2**

Replace:
```go
// getVnetsClient sets up a new client to retrieve vnets
func (c *Client) getVirtualNetworksClient(ctx context.Context) (*aznetwork.VirtualNetworksClient, error) {
	vnetsClient := aznetwork.NewVirtualNetworksClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	vnetsClient.Authorizer = c.ssn.Authorizer
	return &vnetsClient, nil
}
```

With:
```go
// getVirtualNetworksClient sets up a new client to retrieve vnets
func (c *Client) getVirtualNetworksClient() (*armnetwork.VirtualNetworksClient, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	return armnetwork.NewVirtualNetworksClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
}
```

### Step 3.8: Update getSubnetsClient

- [ ] **Step 3.8.1: Replace legacy client factory with V2**

Replace:
```go
// getSubnetsClient sets up a new client to retrieve a subnet
func (c *Client) getSubnetsClient(ctx context.Context) (*aznetwork.SubnetsClient, error) {
	subnetClient := aznetwork.NewSubnetsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	subnetClient.Authorizer = c.ssn.Authorizer
	return &subnetClient, nil
}
```

With:
```go
// getSubnetsClient sets up a new client to retrieve a subnet
func (c *Client) getSubnetsClient() (*armnetwork.SubnetsClient, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	return armnetwork.NewSubnetsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
}
```

- [ ] **Step 3.9: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 3.10: Commit**

```bash
git add pkg/asset/installconfig/azure/client.go
git commit -m "$(cat <<'EOF'
azure: migrate client.go network clients to V2 SDK

Convert network-related methods from legacy SDK profiles to V2:
- Replace aznetwork (profiles/2020-09-01/network) with armnetwork/v2
- Update VirtualNetworksClient and SubnetsClient to V2
- Update API interface method signatures for network types

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: Migrate client.go Subscriptions client to V2 SDK

**Files:**
- Modify: `pkg/asset/installconfig/azure/client.go:150-180` (subscriptions-related methods)
- Add new import: `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions`

### Step 4.1: Add subscriptions import

- [ ] **Step 4.1.1: Add armsubscriptions import**

Add to imports:
```go
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
```

### Step 4.2: Remove legacy subscriptions import

- [ ] **Step 4.2.1: Remove azsubs import**

Remove:
```go
azsubs "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/subscriptions"
```

### Step 4.3: Update API interface - subscriptions types

- [ ] **Step 4.3.1: Update ListLocations signature**

Replace:
```go
ListLocations(ctx context.Context) (*[]azsubs.Location, error)
```

With:
```go
ListLocations(ctx context.Context) ([]*armsubscriptions.Location, error)
```

### Step 4.4: Update ListLocations method

- [ ] **Step 4.4.1: Replace legacy client with V2 client**

Replace:
```go
// ListLocations lists the Azure regions dir the given subscription
func (c *Client) ListLocations(ctx context.Context) (*[]azsubs.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	subsClient, err := c.getSubscriptionsClient(ctx)
	if err != nil {
		return nil, err
	}

	locations, err := subsClient.ListLocations(ctx, c.ssn.Credentials.SubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}

	return locations.Value, nil
}
```

With:
```go
// ListLocations lists the Azure regions for the given subscription
func (c *Client) ListLocations(ctx context.Context) ([]*armsubscriptions.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	subsClient, err := c.getSubscriptionsClient()
	if err != nil {
		return nil, err
	}

	var locations []*armsubscriptions.Location
	pager := subsClient.NewListLocationsPager(c.ssn.Credentials.SubscriptionID, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list locations: %w", err)
		}
		locations = append(locations, page.Value...)
	}

	return locations, nil
}
```

### Step 4.5: Update getSubscriptionsClient

- [ ] **Step 4.5.1: Replace legacy client factory with V2**

Replace:
```go
// getSubscriptionsClient sets up a new client to retrieve subscription data
func (c *Client) getSubscriptionsClient(ctx context.Context) (azsubs.Client, error) {
	client := azsubs.NewClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint)
	client.Authorizer = c.ssn.Authorizer
	return client, nil
}
```

With:
```go
// getSubscriptionsClient sets up a new client to retrieve subscription data
func (c *Client) getSubscriptionsClient() (*armsubscriptions.Client, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	return armsubscriptions.NewClient(c.ssn.TokenCreds, clientOpts)
}
```

- [ ] **Step 4.6: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 4.7: Commit**

```bash
git add pkg/asset/installconfig/azure/client.go
git commit -m "$(cat <<'EOF'
azure: migrate client.go subscriptions client to V2 SDK

Convert subscriptions-related methods from legacy SDK to V2:
- Replace azsubs (profiles/2018-03-01/resources/subscriptions) with armsubscriptions
- Update ListLocations to use V2 pager pattern
- Update API interface method signature

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: Migrate client.go Resources/Providers client to V2 SDK

**Files:**
- Modify: `pkg/asset/installconfig/azure/client.go:178-275` (resources-related methods)

### Step 5.1: Remove legacy resources import

- [ ] **Step 5.1.1: Remove azres import**

Remove:
```go
azres "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/resources"
```

### Step 5.2: Update API interface - resources types

- [ ] **Step 5.2.1: Update GetResourcesProvider and GetGroup signatures**

Replace:
```go
GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*azres.Provider, error)
GetGroup(ctx context.Context, groupName string) (*azres.Group, error)
```

With:
```go
GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*armresources.Provider, error)
GetGroup(ctx context.Context, groupName string) (*armresources.ResourceGroup, error)
```

### Step 5.3: Update GetResourcesProvider method

- [ ] **Step 5.3.1: Replace legacy client with V2 client**

Replace:
```go
// GetResourcesProvider gets the Azure resource provider
func (c *Client) GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*azres.Provider, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	providersClient, err := c.getProvidersClient(ctx)
	if err != nil {
		return nil, err
	}

	provider, err := providersClient.Get(ctx, resourceProviderNamespace, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get resource provider %s: %w", resourceProviderNamespace, err)
	}

	return &provider, nil
}
```

With:
```go
// GetResourcesProvider gets the Azure resource provider
func (c *Client) GetResourcesProvider(ctx context.Context, resourceProviderNamespace string) (*armresources.Provider, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	providersClient, err := c.getProvidersClient()
	if err != nil {
		return nil, err
	}

	resp, err := providersClient.Get(ctx, resourceProviderNamespace, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource provider %s: %w", resourceProviderNamespace, err)
	}

	return &resp.Provider, nil
}
```

### Step 5.4: Update getProvidersClient

- [ ] **Step 5.4.1: Replace legacy client factory with V2**

Replace:
```go
// getProvidersClient sets up a new client to retrieve providers data
func (c *Client) getProvidersClient(ctx context.Context) (azres.ProvidersClient, error) {
	client := azres.NewProvidersClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	return client, nil
}
```

With:
```go
// getProvidersClient sets up a new client to retrieve providers data
func (c *Client) getProvidersClient() (*armresources.ProvidersClient, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	return armresources.NewProvidersClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
}
```

### Step 5.5: Update GetGroup method

- [ ] **Step 5.5.1: Replace legacy client with V2 client**

Replace:
```go
// GetGroup returns resource group for the groupName.
func (c *Client) GetGroup(ctx context.Context, groupName string) (*azres.Group, error) {
	client := azres.NewGroupsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := client.Get(ctx, groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource group: %w", err)
	}
	return &res, nil
}
```

With:
```go
// GetGroup returns resource group for the groupName.
func (c *Client) GetGroup(ctx context.Context, groupName string) (*armresources.ResourceGroup, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armresources.NewResourceGroupsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource groups client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, groupName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource group: %w", err)
	}
	return &resp.ResourceGroup, nil
}
```

### Step 5.6: Update ListResourceIDsByGroup method

- [ ] **Step 5.6.1: Replace legacy client with V2 client**

Replace:
```go
// ListResourceIDsByGroup returns a list of resource IDs for resource group groupName.
func (c *Client) ListResourceIDsByGroup(ctx context.Context, groupName string) ([]string, error) {
	client := azres.NewClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var res []string
	resPage, err := client.ListByResourceGroup(ctx, groupName, "", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}
	for ; resPage.NotDone(); err = resPage.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("error fetching resource pages: %w", err)
		}
		for _, page := range resPage.Values() {
			res = append(res, to.String(page.ID))
		}
	}
	return res, nil
}
```

With:
```go
// ListResourceIDsByGroup returns a list of resource IDs for resource group groupName.
func (c *Client) ListResourceIDsByGroup(ctx context.Context, groupName string) ([]string, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armresources.NewClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var res []string
	pager := client.NewListByResourceGroupPager(groupName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching resource pages: %w", err)
		}
		for _, resource := range page.Value {
			if resource.ID != nil {
				res = append(res, *resource.ID)
			}
		}
	}
	return res, nil
}
```

- [ ] **Step 5.7: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 5.8: Commit**

```bash
git add pkg/asset/installconfig/azure/client.go
git commit -m "$(cat <<'EOF'
azure: migrate client.go resources clients to V2 SDK

Convert resources-related methods from legacy SDK to V2:
- Replace azres (profiles/2018-03-01/resources) with armresources
- Update ProvidersClient, ResourceGroupsClient, and Client to V2
- Update API interface method signatures
- Update pagination to use V2 pager pattern

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: Migrate client.go Compute clients to V2 SDK

**Files:**
- Modify: `pkg/asset/installconfig/azure/client.go:203-450` (compute-related methods)

### Step 6.1: Update imports - replace compute profile

- [ ] **Step 6.1.1: Replace azenc import**

Replace:
```go
azenc "github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
```

With:
```go
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
```

### Step 6.2: Update API interface - compute types

- [ ] **Step 6.2.1: Update compute method signatures**

Replace:
```go
GetVirtualMachineSku(ctx context.Context, name, region string) (*azenc.ResourceSku, error)
GetDiskSkus(ctx context.Context, region string) ([]azenc.ResourceSku, error)
GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName string, diskEncryptionSetName string) (*azenc.DiskEncryptionSet, error)
GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (azenc.VirtualMachineImage, error)
GetLocationInfo(ctx context.Context, region string, instanceType string) (*azenc.ResourceSkuLocationInfo, error)
```

With:
```go
GetVirtualMachineSku(ctx context.Context, name, region string) (*armcompute.ResourceSKU, error)
GetDiskSkus(ctx context.Context, region string) ([]*armcompute.ResourceSKU, error)
GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName string, diskEncryptionSetName string) (*armcompute.DiskEncryptionSet, error)
GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (*armcompute.VirtualMachineImage, error)
GetLocationInfo(ctx context.Context, region string, instanceType string) (*armcompute.ResourceSKULocationInfo, error)
```

### Step 6.3: Update GetDiskSkus method

- [ ] **Step 6.3.1: Replace legacy client with V2 client**

Replace:
```go
// GetDiskSkus returns all the disk SKU pages for a given region.
func (c *Client) GetDiskSkus(ctx context.Context, region string) ([]azenc.ResourceSku, error) {
	client := azenc.NewResourceSkusClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	// See https://issues.redhat.com/browse/OCPBUGS-29469 before changing this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var sku []azenc.ResourceSku
	filter := fmt.Sprintf("location eq '%s'", region)
	// This has to be initialized outside the `for` because we need access to
	// `err`. If initialized in the loop and the API call fails right away,
	// `page.NotDone()` will return `false` and we'll never check for the error
	skuPage, err := client.List(ctx, filter, "false")
	if err != nil {
		return nil, fmt.Errorf("failed to list SKUs: %w", err)
	}
	for ; skuPage.NotDone(); err = skuPage.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, page := range skuPage.Values() {
			for _, diskRegion := range to.StringSlice(page.Locations) {
				if strings.EqualFold(diskRegion, region) {
					sku = append(sku, page)
				}
			}
		}
	}

	if len(sku) != 0 {
		return sku, nil
	}

	return nil, fmt.Errorf("no disks for specified subscription in region %s", region)
}
```

With:
```go
// GetDiskSkus returns all the disk SKU pages for a given region.
func (c *Client) GetDiskSkus(ctx context.Context, region string) ([]*armcompute.ResourceSKU, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armcompute.NewResourceSKUsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource SKUs client: %w", err)
	}
	// See https://issues.redhat.com/browse/OCPBUGS-29469 before changing this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var skus []*armcompute.ResourceSKU
	filter := fmt.Sprintf("location eq '%s'", region)
	pager := client.NewListPager(&armcompute.ResourceSKUsClientListOptions{
		Filter:                   &filter,
		IncludeExtendedLocations: to.Ptr("false"),
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, sku := range page.Value {
			if sku.Locations != nil {
				for _, diskRegion := range sku.Locations {
					if strings.EqualFold(*diskRegion, region) {
						skus = append(skus, sku)
						break
					}
				}
			}
		}
	}

	if len(skus) != 0 {
		return skus, nil
	}

	return nil, fmt.Errorf("no disks for specified subscription in region %s", region)
}
```

### Step 6.4: Update GetVirtualMachineSku method

- [ ] **Step 6.4.1: Replace legacy client with V2 client**

Replace:
```go
// GetVirtualMachineSku retrieves the resource SKU of a specified virtual machine SKU in the specified region.
func (c *Client) GetVirtualMachineSku(ctx context.Context, name, region string) (*azenc.ResourceSku, error) {
	client := azenc.NewResourceSkusClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer

	// See https://issues.redhat.com/browse/OCPBUGS-29469 before chaging this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	filter := fmt.Sprintf("location eq '%s'", region)
	// This has to be initialized outside the `for` because we need access to
	// `err`. If initialized in the loop and the API call fails right away,
	// `page.NotDone()` will return `false` and we'll never check for the error
	page, err := client.List(ctx, filter, "false")
	if err != nil {
		return nil, fmt.Errorf("failed to list SKUs: %w", err)
	}
	for ; page.NotDone(); err = page.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, sku := range page.Values() {
			// Filter out resources that are not virtualMachines
			if !strings.EqualFold("virtualMachines", *sku.ResourceType) {
				continue
			}
			// Filter out resources that do not match the provided name
			if !strings.EqualFold(name, *sku.Name) {
				continue
			}
			// Return the resource from the provided region
			for _, location := range to.StringSlice(sku.Locations) {
				if strings.EqualFold(location, region) {
					return &sku, nil
				}
			}
		}
	}

	return nil, nil
}
```

With:
```go
// GetVirtualMachineSku retrieves the resource SKU of a specified virtual machine SKU in the specified region.
func (c *Client) GetVirtualMachineSku(ctx context.Context, name, region string) (*armcompute.ResourceSKU, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armcompute.NewResourceSKUsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource SKUs client: %w", err)
	}

	// See https://issues.redhat.com/browse/OCPBUGS-29469 before changing this timeout
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	filter := fmt.Sprintf("location eq '%s'", region)
	pager := client.NewListPager(&armcompute.ResourceSKUsClientListOptions{
		Filter:                   &filter,
		IncludeExtendedLocations: to.Ptr("false"),
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching SKU pages: %w", err)
		}
		for _, sku := range page.Value {
			// Filter out resources that are not virtualMachines
			if sku.ResourceType == nil || !strings.EqualFold("virtualMachines", *sku.ResourceType) {
				continue
			}
			// Filter out resources that do not match the provided name
			if sku.Name == nil || !strings.EqualFold(name, *sku.Name) {
				continue
			}
			// Return the resource from the provided region
			if sku.Locations != nil {
				for _, location := range sku.Locations {
					if strings.EqualFold(*location, region) {
						return sku, nil
					}
				}
			}
		}
	}

	return nil, nil
}
```

### Step 6.5: Update GetDiskEncryptionSet method

- [ ] **Step 6.5.1: Replace legacy client with V2 client**

Replace:
```go
// GetDiskEncryptionSet retrieves the specified disk encryption set.
func (c *Client) GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName, diskEncryptionSetName string) (*azenc.DiskEncryptionSet, error) {
	if !strings.EqualFold(c.ssn.Credentials.SubscriptionID, subscriptionID) {
		return nil, fmt.Errorf("different subscription from resource group subscription. Azure does not support cross subscription encryption sets")
	}
	client := azenc.NewDiskEncryptionSetsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	diskEncryptionSet, err := client.Get(ctx, groupName, diskEncryptionSetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk encryption set: %w", err)
	}
	return &diskEncryptionSet, nil
}
```

With:
```go
// GetDiskEncryptionSet retrieves the specified disk encryption set.
func (c *Client) GetDiskEncryptionSet(ctx context.Context, subscriptionID, groupName, diskEncryptionSetName string) (*armcompute.DiskEncryptionSet, error) {
	if !strings.EqualFold(c.ssn.Credentials.SubscriptionID, subscriptionID) {
		return nil, fmt.Errorf("different subscription from resource group subscription. Azure does not support cross subscription encryption sets")
	}
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armcompute.NewDiskEncryptionSetsClient(subscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create disk encryption sets client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, groupName, diskEncryptionSetName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk encryption set: %w", err)
	}
	return &resp.DiskEncryptionSet, nil
}
```

### Step 6.6: Update GetVirtualMachineFamily method

- [ ] **Step 6.6.1: Update to use V2 types**

Replace:
```go
// GetVirtualMachineFamily retrieves the VM family of an instance type.
func (c *Client) GetVirtualMachineFamily(ctx context.Context, name, region string) (string, error) {
	typeMeta, err := c.GetVirtualMachineSku(ctx, name, region)
	if err != nil {
		return "", fmt.Errorf("error connecting to Azure client: %w", err)
	}
	if typeMeta == nil {
		return "", fmt.Errorf("not found in region %s", region)
	}
	if typeMeta.Family == nil {
		return "", fmt.Errorf("error getting resource family")
	}

	return to.String(typeMeta.Family), nil
}
```

With:
```go
// GetVirtualMachineFamily retrieves the VM family of an instance type.
func (c *Client) GetVirtualMachineFamily(ctx context.Context, name, region string) (string, error) {
	typeMeta, err := c.GetVirtualMachineSku(ctx, name, region)
	if err != nil {
		return "", fmt.Errorf("error connecting to Azure client: %w", err)
	}
	if typeMeta == nil {
		return "", fmt.Errorf("not found in region %s", region)
	}
	if typeMeta.Family == nil {
		return "", fmt.Errorf("error getting resource family")
	}

	return *typeMeta.Family, nil
}
```

### Step 6.7: Update GetVMCapabilities method

- [ ] **Step 6.7.1: Update to use V2 types**

Replace:
```go
// GetVMCapabilities retrieves the capabilities of an instant type in a specific region. Returns these values
// in a map with the capability name as the key and the corresponding value.
func (c *Client) GetVMCapabilities(ctx context.Context, instanceType, region string) (map[string]string, error) {
	typeMeta, err := c.GetVirtualMachineSku(ctx, instanceType, region)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Azure client: %w", err)
	}
	if typeMeta == nil {
		return nil, fmt.Errorf("not found in region %s", region)
	}
	capabilities := make(map[string]string)
	for _, capability := range *typeMeta.Capabilities {
		capabilities[to.String(capability.Name)] = to.String(capability.Value)
	}

	return capabilities, nil
}
```

With:
```go
// GetVMCapabilities retrieves the capabilities of an instant type in a specific region. Returns these values
// in a map with the capability name as the key and the corresponding value.
func (c *Client) GetVMCapabilities(ctx context.Context, instanceType, region string) (map[string]string, error) {
	typeMeta, err := c.GetVirtualMachineSku(ctx, instanceType, region)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Azure client: %w", err)
	}
	if typeMeta == nil {
		return nil, fmt.Errorf("not found in region %s", region)
	}
	capabilities := make(map[string]string)
	if typeMeta.Capabilities != nil {
		for _, capability := range typeMeta.Capabilities {
			if capability.Name != nil && capability.Value != nil {
				capabilities[*capability.Name] = *capability.Value
			}
		}
	}

	return capabilities, nil
}
```

### Step 6.8: Update GetMarketplaceImage method

- [ ] **Step 6.8.1: Replace legacy client with V2 client**

Replace:
```go
// GetMarketplaceImage get the specified marketplace VM image.
func (c *Client) GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (azenc.VirtualMachineImage, error) {
	client := azenc.NewVirtualMachineImagesClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	image, err := client.Get(ctx, region, publisher, offer, sku, version)
	if err != nil {
		return image, fmt.Errorf("could not get marketplace image: %w", err)
	}
	return image, nil
}
```

With:
```go
// GetMarketplaceImage get the specified marketplace VM image.
func (c *Client) GetMarketplaceImage(ctx context.Context, region, publisher, offer, sku, version string) (*armcompute.VirtualMachineImage, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armcompute.NewVirtualMachineImagesClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual machine images client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, region, publisher, offer, sku, version, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get marketplace image: %w", err)
	}
	return &resp.VirtualMachineImage, nil
}
```

### Step 6.9: Update GetAvailabilityZones method

- [ ] **Step 6.9.1: Update to use V2 types**

Replace:
```go
// GetAvailabilityZones retrieves a list of availability zones for the given region, and instance type.
func (c *Client) GetAvailabilityZones(ctx context.Context, region string, instanceType string) ([]string, error) {
	locationInfo, err := c.GetLocationInfo(ctx, region, instanceType)
	if err != nil {
		return nil, err
	}
	if locationInfo != nil {
		return to.StringSlice(locationInfo.Zones), nil
	}

	return nil, fmt.Errorf("error retrieving availability zones for %s in %s", instanceType, region)
}
```

With:
```go
// GetAvailabilityZones retrieves a list of availability zones for the given region, and instance type.
func (c *Client) GetAvailabilityZones(ctx context.Context, region string, instanceType string) ([]string, error) {
	locationInfo, err := c.GetLocationInfo(ctx, region, instanceType)
	if err != nil {
		return nil, err
	}
	if locationInfo != nil && locationInfo.Zones != nil {
		zones := make([]string, 0, len(locationInfo.Zones))
		for _, z := range locationInfo.Zones {
			if z != nil {
				zones = append(zones, *z)
			}
		}
		return zones, nil
	}

	return nil, fmt.Errorf("error retrieving availability zones for %s in %s", instanceType, region)
}
```

### Step 6.10: Update GetLocationInfo method

- [ ] **Step 6.10.1: Replace legacy client with V2 client**

Replace:
```go
// GetLocationInfo retrieves the location info associated with the instance type in region
func (c *Client) GetLocationInfo(ctx context.Context, region string, instanceType string) (*azenc.ResourceSkuLocationInfo, error) {
	client := azenc.NewResourceSkusClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer

	// Only supported filter atm is `location`
	filter := fmt.Sprintf("location eq '%s'", region)
	res, err := client.List(ctx, filter, "false")
	if err != nil {
		return nil, fmt.Errorf("failed to list SKUs: %w", err)
	}
	for ; res.NotDone(); err = res.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}

		for _, resSku := range res.Values() {
			if !strings.EqualFold(to.String(resSku.ResourceType), "virtualMachines") {
				continue
			}
			if strings.EqualFold(to.String(resSku.Name), instanceType) {
				for _, locationInfo := range *resSku.LocationInfo {
					return &locationInfo, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("location information not found for %s in %s", instanceType, region)
}
```

With:
```go
// GetLocationInfo retrieves the location info associated with the instance type in region
func (c *Client) GetLocationInfo(ctx context.Context, region string, instanceType string) (*armcompute.ResourceSKULocationInfo, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armcompute.NewResourceSKUsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource SKUs client: %w", err)
	}

	// Only supported filter atm is `location`
	filter := fmt.Sprintf("location eq '%s'", region)
	pager := client.NewListPager(&armcompute.ResourceSKUsClientListOptions{
		Filter:                   &filter,
		IncludeExtendedLocations: to.Ptr("false"),
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, resSku := range page.Value {
			if resSku.ResourceType == nil || !strings.EqualFold(*resSku.ResourceType, "virtualMachines") {
				continue
			}
			if resSku.Name != nil && strings.EqualFold(*resSku.Name, instanceType) {
				if resSku.LocationInfo != nil && len(resSku.LocationInfo) > 0 {
					return resSku.LocationInfo[0], nil
				}
			}
		}
	}

	return nil, fmt.Errorf("location information not found for %s in %s", instanceType, region)
}
```

### Step 6.11: Remove legacy to import and fix usages

- [ ] **Step 6.11.1: Remove legacy to import**

Remove:
```go
"github.com/Azure/go-autorest/autorest/to"
```

Update remaining usages to use `to.Ptr()` from `azcore/to` instead of `to.String()`, `to.StringSlice()`, etc.

- [ ] **Step 6.12: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 6.13: Commit**

```bash
git add pkg/asset/installconfig/azure/client.go
git commit -m "$(cat <<'EOF'
azure: migrate client.go compute clients to V2 SDK

Convert compute-related methods from legacy SDK to V2:
- Replace azenc (profiles/latest/compute) with armcompute/v5
- Update ResourceSKUsClient, DiskEncryptionSetsClient, VirtualMachineImagesClient to V2
- Update API interface method signatures for compute types
- Update pagination to use V2 pager pattern
- Replace autorest/to with azcore/to

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: Migrate client.go Marketplace client to V2 SDK

**Files:**
- Modify: `pkg/asset/installconfig/azure/client.go:384-401` (marketplace methods)
- Add new import: `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering`

### Step 7.1: Add marketplace import

- [ ] **Step 7.1.1: Add armmarketplaceordering import**

Add to imports:
```go
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
```

### Step 7.2: Remove legacy marketplace import

- [ ] **Step 7.2.1: Remove azmarketplace import**

Remove:
```go
azmarketplace "github.com/Azure/azure-sdk-for-go/profiles/latest/marketplaceordering/mgmt/marketplaceordering"
```

### Step 7.3: Update AreMarketplaceImageTermsAccepted method

- [ ] **Step 7.3.1: Replace legacy client with V2 client**

Replace:
```go
// AreMarketplaceImageTermsAccepted tests whether the terms have been accepted for the specified marketplace VM image.
func (c *Client) AreMarketplaceImageTermsAccepted(ctx context.Context, publisher, offer, sku string) (bool, error) {
	client := azmarketplace.NewMarketplaceAgreementsClientWithBaseURI(c.ssn.Environment.ResourceManagerEndpoint, c.ssn.Credentials.SubscriptionID)
	client.Authorizer = c.ssn.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	terms, err := client.Get(ctx, publisher, offer, sku)
	if err != nil {
		return false, err
	}

	if terms.AgreementProperties == nil {
		return false, errors.New("no agreement properties for image")
	}

	return terms.AgreementProperties.Accepted != nil && *terms.AgreementProperties.Accepted, nil
}
```

With:
```go
// AreMarketplaceImageTermsAccepted tests whether the terms have been accepted for the specified marketplace VM image.
func (c *Client) AreMarketplaceImageTermsAccepted(ctx context.Context, publisher, offer, sku string) (bool, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.ssn.CloudConfig,
		},
	}
	client, err := armmarketplaceordering.NewMarketplaceAgreementsClient(c.ssn.Credentials.SubscriptionID, c.ssn.TokenCreds, clientOpts)
	if err != nil {
		return false, fmt.Errorf("failed to create marketplace agreements client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, armmarketplaceordering.OfferTypeVirtualmachine, publisher, offer, sku, nil)
	if err != nil {
		return false, err
	}

	if resp.Properties == nil {
		return false, errors.New("no agreement properties for image")
	}

	return resp.Properties.Accepted != nil && *resp.Properties.Accepted, nil
}
```

- [ ] **Step 7.4: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 7.5: Commit**

```bash
git add pkg/asset/installconfig/azure/client.go
git commit -m "$(cat <<'EOF'
azure: migrate client.go marketplace client to V2 SDK

Convert marketplace-related methods from legacy SDK to V2:
- Replace azmarketplace (profiles/latest/marketplaceordering) with armmarketplaceordering
- Update AreMarketplaceImageTermsAccepted to use V2 client

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: Update validation.go to use V2 types

**Files:**
- Modify: `pkg/asset/installconfig/azure/validation.go:1-1112` (all legacy type references)

### Step 8.1: Update imports

- [ ] **Step 8.1.1: Update imports to use V2 packages**

Replace:
```go
aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
azenc "github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
"github.com/Azure/go-autorest/autorest/to"
```

With:
```go
"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
```

### Step 8.2: Update validateConfidentialDiskEncryptionSet

- [ ] **Step 8.2.1: Update encryption type comparison**

Replace:
```go
} else if resp == nil || resp.EncryptionSetProperties == nil || resp.EncryptionSetProperties.EncryptionType != azenc.ConfidentialVMEncryptedWithCustomerKey {
	return fmt.Errorf("the disk encryption set should be created with type %s", azenc.ConfidentialVMEncryptedWithCustomerKey)
}
```

With:
```go
} else if resp == nil || resp.Properties == nil || resp.Properties.EncryptionType == nil || *resp.Properties.EncryptionType != armcompute.DiskEncryptionSetTypeConfidentialVMEncryptedWithCustomerKey {
	return fmt.Errorf("the disk encryption set should be created with type %s", armcompute.DiskEncryptionSetTypeConfidentialVMEncryptedWithCustomerKey)
}
```

### Step 8.3: Update validateSubnet function

- [ ] **Step 8.3.1: Update subnet type references**

Replace:
```go
func validateSubnet(client API, fieldPath *field.Path, subnet *aznetwork.Subnet, subnetName string, networks []types.MachineNetworkEntry) field.ErrorList {
```

With:
```go
func validateSubnet(client API, fieldPath *field.Path, subnet *armnetwork.Subnet, subnetName string, networks []types.MachineNetworkEntry) field.ErrorList {
```

### Step 8.4: Update validateSubnet property access

- [ ] **Step 8.4.1: Update SubnetPropertiesFormat access**

Replace:
```go
if subnet == nil || subnet.SubnetPropertiesFormat == nil {
```

With:
```go
if subnet == nil || subnet.Properties == nil {
```

### Step 8.5: Update validateSubnet address prefix access

- [ ] **Step 8.5.1: Update address prefix property access**

Replace:
```go
switch {
case subnet.AddressPrefix != nil:
	addressPrefix = *subnet.AddressPrefix
// NOTE: if the subscription has the `AllowMultipleAddressPrefixesOnSubnet` feature, the Azure API will return a
// `addressPrefixes` field with a slice of addresses instead of a single value via `addressPrefix`.
case subnet.AddressPrefixes != nil && len(*subnet.AddressPrefixes) > 0:
	addressPrefix = (*subnet.AddressPrefixes)[0]
```

With:
```go
switch {
case subnet.Properties.AddressPrefix != nil:
	addressPrefix = *subnet.Properties.AddressPrefix
// NOTE: if the subscription has the `AllowMultipleAddressPrefixesOnSubnet` feature, the Azure API will return a
// `addressPrefixes` field with a slice of addresses instead of a single value via `addressPrefix`.
case subnet.Properties.AddressPrefixes != nil && len(subnet.Properties.AddressPrefixes) > 0:
	addressPrefix = *subnet.Properties.AddressPrefixes[0]
```

### Step 8.6: Update validateRegion function

- [ ] **Step 8.6.1: Update location property access**

Replace:
```go
for _, location := range *locations {
	availableRegions[to.String(location.Name)] = to.String(location.DisplayName)
}
```

With:
```go
for _, location := range locations {
	if location.Name != nil && location.DisplayName != nil {
		availableRegions[*location.Name] = *location.DisplayName
	}
}
```

### Step 8.7: Update validateRegion provider property access

- [ ] **Step 8.7.1: Update provider resource types access**

Replace:
```go
for _, resType := range *provider.ResourceTypes {
	if *resType.ResourceType == "resourceGroups" {
		for _, resourceCapableRegion := range *resType.Locations {
```

With:
```go
for _, resType := range provider.ResourceTypes {
	if resType.ResourceType != nil && *resType.ResourceType == "resourceGroups" {
		if resType.Locations != nil {
			for _, resourceCapableRegion := range resType.Locations {
```

### Step 8.8: Update validateResourceGroup function

- [ ] **Step 8.8.1: Update group property access**

Replace:
```go
normalizedRegion := strings.Replace(strings.ToLower(to.String(group.Location)), " ", "", -1)
```

With:
```go
var location string
if group.Location != nil {
	location = *group.Location
}
normalizedRegion := strings.Replace(strings.ToLower(location), " ", "", -1)
```

### Step 8.9: Update validateSubnetNatGateway function

- [ ] **Step 8.9.1: Update subnet type reference**

Replace:
```go
func validateSubnetNatGateway(client API, fieldPath *field.Path, subnet *aznetwork.Subnet, outboundType aztypes.OutboundType, role capz.SubnetRole, resourceGroup, virtualNetwork string) field.ErrorList {
```

With:
```go
func validateSubnetNatGateway(client API, fieldPath *field.Path, subnet *armnetwork.Subnet, outboundType aztypes.OutboundType, role capz.SubnetRole, resourceGroup, virtualNetwork string) field.ErrorList {
```

### Step 8.10: Update validateCustomSubnets function

- [ ] **Step 8.10.1: Update vnetSubnetList type**

Replace:
```go
vnetSubnetList := map[string]*aznetwork.Subnet{}
```

With:
```go
vnetSubnetList := map[string]*armnetwork.Subnet{}
```

### Step 8.11: Update validateCustomSubnets property access

- [ ] **Step 8.11.1: Update VirtualNetworkPropertiesFormat access**

Replace:
```go
if existingVnet.VirtualNetworkPropertiesFormat != nil && existingVnet.Subnets != nil {
	for _, subnet := range *existingVnet.Subnets {
		vnetSubnetList[*subnet.Name] = &subnet
	}
}
```

With:
```go
if existingVnet.Properties != nil && existingVnet.Properties.Subnets != nil {
	for _, subnet := range existingVnet.Properties.Subnets {
		if subnet.Name != nil {
			vnetSubnetList[*subnet.Name] = subnet
		}
	}
}
```

### Step 8.12: Update validateUltraSSD function

- [ ] **Step 8.12.1: Update locationInfo property access**

Replace:
```go
if locationInfo == nil || len(to.StringSlice(locationInfo.Zones)) == 0 {
```

With:
```go
if locationInfo == nil || locationInfo.Zones == nil || len(locationInfo.Zones) == 0 {
```

And replace:
```go
allZones := to.StringSlice(locationInfo.Zones)
```

With:
```go
allZones := make([]string, 0, len(locationInfo.Zones))
for _, z := range locationInfo.Zones {
	if z != nil {
		allZones = append(allZones, *z)
	}
}
```

### Step 8.13: Update validateUltraSSD zone details access

- [ ] **Step 8.13.1: Update ZoneDetails property access**

Replace:
```go
for _, zoneDetails := range *locationInfo.ZoneDetails {
	for _, capability := range *zoneDetails.Capabilities {
		if !strings.EqualFold(to.String(capability.Name), "UltraSSDAvailable") {
			continue
		}
		if strings.EqualFold(to.String(capability.Value), "True") {
			// ...
			capZones := to.StringSlice(zoneDetails.Name)
```

With:
```go
if locationInfo.ZoneDetails != nil {
	for _, zoneDetails := range locationInfo.ZoneDetails {
		if zoneDetails.Capabilities == nil {
			continue
		}
		for _, capability := range zoneDetails.Capabilities {
			if capability.Name == nil || !strings.EqualFold(*capability.Name, "UltraSSDAvailable") {
				continue
			}
			if capability.Value != nil && strings.EqualFold(*capability.Value, "True") {
				// ...
				capZones := make([]string, 0)
				if zoneDetails.Name != nil {
					for _, n := range zoneDetails.Name {
						if n != nil {
							capZones = append(capZones, *n)
						}
					}
				}
```

- [ ] **Step 8.14: Verify compilation**

Run: `go build ./pkg/asset/installconfig/azure/...`
Expected: No compilation errors

- [ ] **Step 8.15: Run tests**

Run: `go test ./pkg/asset/installconfig/azure/... -v -short`
Expected: Tests pass

- [ ] **Step 8.16: Commit**

```bash
git add pkg/asset/installconfig/azure/validation.go
git commit -m "$(cat <<'EOF'
azure: migrate validation.go to V2 SDK types

Update validation.go to use V2 SDK types:
- Replace aznetwork profile with armnetwork/v2
- Replace azenc profile with armcompute/v5
- Replace autorest/to with azcore/to
- Update all property access patterns for V2 struct layout
- Update DiskEncryptionSetType comparison

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: Regenerate mocks

**Files:**
- Regenerate: `pkg/asset/installconfig/azure/mock/azureclient_generated.go`

- [ ] **Step 9.1: Regenerate mocks**

Run: `go generate ./pkg/asset/installconfig/azure/...`
Expected: Mock file regenerated with new V2 types

- [ ] **Step 9.2: Verify mock compilation**

Run: `go build ./pkg/asset/installconfig/azure/mock/...`
Expected: No compilation errors

- [ ] **Step 9.3: Run all tests**

Run: `go test ./pkg/asset/installconfig/azure/... -v`
Expected: All tests pass

- [ ] **Step 9.4: Commit**

```bash
git add pkg/asset/installconfig/azure/mock/azureclient_generated.go
git commit -m "$(cat <<'EOF'
azure: regenerate mocks for V2 SDK API interface

Regenerate mock file after API interface changes to use V2 SDK types.

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: Migrate destroy/azure.go to V2 SDK

**Files:**
- Modify: `pkg/destroy/azure/azure.go:1-894` (DNS and resources clients)

### Step 10.1: Update imports

- [ ] **Step 10.1.1: Remove legacy imports**

Remove:
```go
azurestackdns "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/dns/mgmt/dns"
"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
"github.com/Azure/go-autorest/autorest"
azureenv "github.com/Azure/go-autorest/autorest/azure"
"github.com/Azure/go-autorest/autorest/to"
```

Add:
```go
"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
```

### Step 10.2: Update ClusterUninstaller struct

- [ ] **Step 10.2.1: Update client types**

Replace:
```go
resourceGroupsClient    resources.GroupsClient
zonesClient             dns.ZonesClient
recordsClient           dns.RecordSetsClient
privateRecordSetsClient privatedns.RecordSetsClient
privateZonesClient      privatedns.PrivateZonesClient
```

With:
```go
resourceGroupsClient    *armresources.ResourceGroupsClient
zonesClient             *armdns.ZonesClient
recordsClient           *armdns.RecordSetsClient
privateRecordSetsClient *armprivatedns.RecordSetsClient
privateZonesClient      *armprivatedns.PrivateZonesClient
```

### Step 10.3: Update configureClients method

- [ ] **Step 10.3.1: Replace legacy client creation with V2**

Replace:
```go
func (o *ClusterUninstaller) configureClients() error {
	subscriptionID := o.Session.Credentials.SubscriptionID
	endpoint := o.Session.Environment.ResourceManagerEndpoint

	o.resourceGroupsClient = resources.NewGroupsClientWithBaseURI(endpoint, subscriptionID)
	o.resourceGroupsClient.Authorizer = o.Session.Authorizer

	o.zonesClient = dns.NewZonesClientWithBaseURI(endpoint, subscriptionID)
	o.zonesClient.Authorizer = o.Session.Authorizer

	o.recordsClient = dns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionID)
	o.recordsClient.Authorizer = o.Session.Authorizer

	o.privateZonesClient = privatedns.NewPrivateZonesClientWithBaseURI(endpoint, subscriptionID)
	o.privateZonesClient.Authorizer = o.Session.Authorizer

	o.privateRecordSetsClient = privatedns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionID)
	o.privateRecordSetsClient.Authorizer = o.Session.Authorizer
```

With:
```go
func (o *ClusterUninstaller) configureClients() error {
	subscriptionID := o.Session.Credentials.SubscriptionID

	clientOpts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: o.Session.CloudConfig,
		},
	}

	var err error
	o.resourceGroupsClient, err = armresources.NewResourceGroupsClient(subscriptionID, o.Session.TokenCreds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create resource groups client: %w", err)
	}

	o.zonesClient, err = armdns.NewZonesClient(subscriptionID, o.Session.TokenCreds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create DNS zones client: %w", err)
	}

	o.recordsClient, err = armdns.NewRecordSetsClient(subscriptionID, o.Session.TokenCreds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create DNS record sets client: %w", err)
	}

	o.privateZonesClient, err = armprivatedns.NewPrivateZonesClient(subscriptionID, o.Session.TokenCreds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create private DNS zones client: %w", err)
	}

	o.privateRecordSetsClient, err = armprivatedns.NewRecordSetsClient(subscriptionID, o.Session.TokenCreds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create private DNS record sets client: %w", err)
	}
```

(Continue with remaining destroy.go updates - deletePublicRecords, deleteResourceGroup, etc.)

- [ ] **Step 10.4: Verify compilation**

Run: `go build ./pkg/destroy/azure/...`
Expected: No compilation errors

- [ ] **Step 10.5: Commit**

```bash
git add pkg/destroy/azure/azure.go
git commit -m "$(cat <<'EOF'
azure: migrate destroy/azure.go to V2 SDK

Convert destroy operations from legacy SDK to V2:
- Replace dns, privatedns, resources service clients with V2
- Update ClusterUninstaller client types
- Update all pagination to use V2 pager pattern
- Update error handling for V2 response patterns

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: Final verification and cleanup

- [ ] **Step 11.1: Run full build**

Run: `hack/build.sh`
Expected: Build succeeds

- [ ] **Step 11.2: Run unit tests**

Run: `hack/go-test.sh`
Expected: All tests pass

- [ ] **Step 11.3: Run linter**

Run: `hack/go-lint.sh $(go list -f '{{ .ImportPath }}' ./pkg/asset/installconfig/azure/...) $(go list -f '{{ .ImportPath }}' ./pkg/destroy/azure/...)`
Expected: No lint errors

- [ ] **Step 11.4: Verify no legacy imports remain**

Run: `grep -r "github.com/Azure/azure-sdk-for-go/profiles" pkg/asset/installconfig/azure/ pkg/destroy/azure/`
Expected: No matches

Run: `grep -r "github.com/Azure/azure-sdk-for-go/services" pkg/asset/installconfig/azure/ pkg/destroy/azure/`
Expected: No matches (except azurestackdns for Azure Stack compatibility)

- [ ] **Step 11.5: Final commit**

```bash
git add -A
git commit -m "$(cat <<'EOF'
azure: complete V2 SDK migration for validation and destroy

Complete migration of Azure SDK from legacy profiles/services to V2:
- All validation layer code now uses V2 SDK
- All destroy operations now use V2 SDK
- Updated mocks for new API interface
- Verified build and tests pass

Legacy SDK usage remains only for Azure Stack DNS operations
which require specific API version compatibility.

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

---

## Summary

This plan migrates the following legacy Azure SDK packages to V2:

| Legacy Package | V2 Replacement |
|---------------|----------------|
| `profiles/2018-03-01/dns/mgmt/dns` | `resourcemanager/dns/armdns` |
| `profiles/2020-09-01/network/mgmt/network` | `resourcemanager/network/armnetwork/v2` |
| `profiles/latest/compute/mgmt/compute` | `resourcemanager/compute/armcompute/v5` |
| `profiles/2018-03-01/resources/mgmt/resources` | `resourcemanager/resources/armresources` |
| `profiles/2018-03-01/resources/mgmt/subscriptions` | `resourcemanager/resources/armsubscriptions` |
| `profiles/latest/marketplaceordering/mgmt/marketplaceordering` | `resourcemanager/marketplaceordering/armmarketplaceordering` |
| `services/preview/dns/mgmt/2018-03-01-preview/dns` | `resourcemanager/dns/armdns` |
| `services/privatedns/mgmt/2018-09-01/privatedns` | `resourcemanager/privatedns/armprivatedns` |
| `services/resources/mgmt/2018-05-01/resources` | `resourcemanager/resources/armresources` |
| `go-autorest/autorest/to` | `azcore/to` |

**Note:** Azure Stack DNS operations (`azurestackdns`) may need to remain on legacy SDK if the V2 SDK doesn't support the required API version override.
