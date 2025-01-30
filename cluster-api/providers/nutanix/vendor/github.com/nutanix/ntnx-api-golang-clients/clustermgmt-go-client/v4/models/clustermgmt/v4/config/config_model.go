/*
 * Generated file models/clustermgmt/v4/config/config_model.go.
 *
 * Product version: 4.0.1-beta-2
 *
 * Part of the Nutanix Clustermgmt Versioned APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Configure Hosts, Clusters and other Infrastructure
*/
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	import2 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/error"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	"time"
)

/*
Status of Acropolis connection to hypervisor.
*/
type AcropolisConnectionState int

const (
	ACROPOLISCONNECTIONSTATE_UNKNOWN      AcropolisConnectionState = 0
	ACROPOLISCONNECTIONSTATE_REDACTED     AcropolisConnectionState = 1
	ACROPOLISCONNECTIONSTATE_CONNECTED    AcropolisConnectionState = 2
	ACROPOLISCONNECTIONSTATE_DISCONNECTED AcropolisConnectionState = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *AcropolisConnectionState) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CONNECTED",
		"DISCONNECTED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e AcropolisConnectionState) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CONNECTED",
		"DISCONNECTED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *AcropolisConnectionState) index(name string) AcropolisConnectionState {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CONNECTED",
		"DISCONNECTED",
	}
	for idx := range names {
		if names[idx] == name {
			return AcropolisConnectionState(idx)
		}
	}
	return ACROPOLISCONNECTIONSTATE_UNKNOWN
}

func (e *AcropolisConnectionState) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for AcropolisConnectionState:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *AcropolisConnectionState) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e AcropolisConnectionState) Ref() *AcropolisConnectionState {
	return &e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/snmp/$actions/add-transports Post operation
*/
type AddSnmpTransportsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfAddSnmpTransportsApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewAddSnmpTransportsApiResponse() *AddSnmpTransportsApiResponse {
	p := new(AddSnmpTransportsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.AddSnmpTransportsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *AddSnmpTransportsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *AddSnmpTransportsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfAddSnmpTransportsApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Indicates whether the address type is IPV4/IPV6.
*/
type AddressType int

const (
	ADDRESSTYPE_UNKNOWN  AddressType = 0
	ADDRESSTYPE_REDACTED AddressType = 1
	ADDRESSTYPE_IPV4     AddressType = 2
	ADDRESSTYPE_IPV6     AddressType = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *AddressType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"IPV4",
		"IPV6",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e AddressType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"IPV4",
		"IPV6",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *AddressType) index(name string) AddressType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"IPV4",
		"IPV6",
	}
	for idx := range names {
		if names[idx] == name {
			return AddressType(idx)
		}
	}
	return ADDRESSTYPE_UNKNOWN
}

func (e *AddressType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for AddressType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *AddressType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e AddressType) Ref() *AddressType {
	return &e
}

/*
Attribute item information.
*/
type AttributeItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Tolerance message attribute key.
	*/
	Attribute *string `json:"attribute,omitempty"`
	/*
	  Tolerance message attribute value.
	*/
	Value *string `json:"value,omitempty"`
}

func NewAttributeItem() *AttributeItem {
	p := new(AttributeItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.AttributeItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Params associated to the backplane network segmentation. This is part of payload for cluster create operation only.
*/
type BackplaneNetworkParams struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Flag to indicate if the backplane segmentation needs to be enabled or not.
	*/
	IsSegmentationEnabled *bool `json:"isSegmentationEnabled,omitempty"`

	Netmask *import4.IPv4Address `json:"netmask,omitempty"`

	Subnet *import4.IPv4Address `json:"subnet,omitempty"`
	/*
	  VLAN Id tagged to the backplane network on the cluster. This is part of cluster create payload.
	*/
	VlanTag *int64 `json:"vlanTag,omitempty"`
}

func NewBackplaneNetworkParams() *BackplaneNetworkParams {
	p := new(BackplaneNetworkParams)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.BackplaneNetworkParams"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Block item containing block serial and rack name.
*/
type BlockItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rackable unit serial name.
	*/
	BlockId *string `json:"blockId,omitempty"`
	/*
	  Rack name.
	*/
	RackName *string `json:"rackName,omitempty"`
}

func NewBlockItem() *BlockItem {
	p := new(BlockItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.BlockItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Build information details.
*/
type BuildReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Software build type.
	*/
	BuildType *string `json:"buildType,omitempty"`
	/*
	  Commit Id used for version.
	*/
	CommitId *string `json:"commitId,omitempty"`
	/*
	  Full name of software version.
	*/
	FullVersion *string `json:"fullVersion,omitempty"`
	/*
	  Short commit Id used for version.
	*/
	ShortCommitId *string `json:"shortCommitId,omitempty"`
	/*
	  Software version.
	*/
	Version *string `json:"version,omitempty"`
}

func NewBuildReference() *BuildReference {
	p := new(BuildReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.BuildReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Hypervisor bundle information.
*/
type BundleInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Name of the hypervisor bundle.
	*/
	Name *string `json:"name,omitempty"`
}

func NewBundleInfo() *BundleInfo {
	p := new(BundleInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.BundleInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
ISO attributes to validate compatibility.
*/
type BundleParam struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	BundleInfo *BundleInfo `json:"bundleInfo"`
	/*
	  List of node attributes for validating bundle compatibility.
	*/
	NodeList []NodeInfo `json:"nodeList"`
}

func (p *BundleParam) MarshalJSON() ([]byte, error) {
	type BundleParamProxy BundleParam
	return json.Marshal(struct {
		*BundleParamProxy
		BundleInfo *BundleInfo `json:"bundleInfo,omitempty"`
		NodeList   []NodeInfo  `json:"nodeList,omitempty"`
	}{
		BundleParamProxy: (*BundleParamProxy)(p),
		BundleInfo:       p.BundleInfo,
		NodeList:         p.NodeList,
	})
}

func NewBundleParam() *BundleParam {
	p := new(BundleParam)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.BundleParam"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Indicates the current status of Cache Deduplication for the Container.
*/
type CacheDeduplication int

const (
	CACHEDEDUPLICATION_UNKNOWN  CacheDeduplication = 0
	CACHEDEDUPLICATION_REDACTED CacheDeduplication = 1
	CACHEDEDUPLICATION_NONE     CacheDeduplication = 2
	CACHEDEDUPLICATION_OFF      CacheDeduplication = 3
	CACHEDEDUPLICATION_ON       CacheDeduplication = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *CacheDeduplication) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"ON",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e CacheDeduplication) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"ON",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *CacheDeduplication) index(name string) CacheDeduplication {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"ON",
	}
	for idx := range names {
		if names[idx] == name {
			return CacheDeduplication(idx)
		}
	}
	return CACHEDEDUPLICATION_UNKNOWN
}

func (e *CacheDeduplication) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for CacheDeduplication:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *CacheDeduplication) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e CacheDeduplication) Ref() *CacheDeduplication {
	return &e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/$actions/check-hypervisor-requirements Post operation
*/
type CheckHypervisorRequirementsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCheckHypervisorRequirementsApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCheckHypervisorRequirementsApiResponse() *CheckHypervisorRequirementsApiResponse {
	p := new(CheckHypervisorRequirementsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.CheckHypervisorRequirementsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CheckHypervisorRequirementsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CheckHypervisorRequirementsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCheckHypervisorRequirementsApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

type Cluster struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Config *ClusterConfigReference `json:"config,omitempty"`
	/*
	  The name of the default container created as part of cluster creation. This is part of payload for cluster create operation only.
	*/
	ContainerName *string `json:"containerName,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Number of inefficient VMs.
	*/
	InefficientVmCount *int64 `json:"inefficientVmCount,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Cluster name. This is part of payload for both cluster create & update operations.
	*/
	Name *string `json:"name,omitempty"`

	Network *ClusterNetworkReference `json:"network,omitempty"`

	Nodes *NodeReference `json:"nodes,omitempty"`
	/*
	  Flag determines whether only prechecks must be run and create cluster operation be skipped.This is an optional param for cluster create operation only.
	*/
	RunPrechecksOnly *bool `json:"runPrechecksOnly,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`

	UpgradeStatus *UpgradeStatus `json:"upgradeStatus,omitempty"`
	/*
	  Number of VMs.
	*/
	VmCount *int64 `json:"vmCount,omitempty"`
}

func NewCluster() *Cluster {
	p := new(Cluster)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.Cluster"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Cluster arch.
*/
type ClusterArchReference int

const (
	CLUSTERARCHREFERENCE_UNKNOWN  ClusterArchReference = 0
	CLUSTERARCHREFERENCE_REDACTED ClusterArchReference = 1
	CLUSTERARCHREFERENCE_X86_64   ClusterArchReference = 2
	CLUSTERARCHREFERENCE_PPC64LE  ClusterArchReference = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ClusterArchReference) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"X86_64",
		"PPC64LE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ClusterArchReference) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"X86_64",
		"PPC64LE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ClusterArchReference) index(name string) ClusterArchReference {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"X86_64",
		"PPC64LE",
	}
	for idx := range names {
		if names[idx] == name {
			return ClusterArchReference(idx)
		}
	}
	return CLUSTERARCHREFERENCE_UNKNOWN
}

func (e *ClusterArchReference) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ClusterArchReference:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ClusterArchReference) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ClusterArchReference) Ref() *ClusterArchReference {
	return &e
}

/*
Cluster configuration details.
*/
type ClusterConfigReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Public ssh key details. This is part of payload for cluster update operation only.
	*/
	AuthorizedPublicKeyList []PublicKey `json:"authorizedPublicKeyList,omitempty"`

	BuildInfo *BuildReference `json:"buildInfo,omitempty"`

	ClusterArch *ClusterArchReference `json:"clusterArch,omitempty"`
	/*
	  Cluster function. This is part of payload for cluster create operation only (allowed enum values for creation are AOS, ONE_NODE & TWO_NODE only).
	*/
	ClusterFunction []ClusterFunctionRef `json:"clusterFunction,omitempty"`
	/*
	  Cluster software version details.
	*/
	ClusterSoftwareMap []SoftwareMapReference `json:"clusterSoftwareMap,omitempty"`

	EncryptionInTransitStatus *EncryptionStatus `json:"encryptionInTransitStatus,omitempty"`
	/*
	  Encryption option.
	*/
	EncryptionOption []EncryptionOptionInfo `json:"encryptionOption,omitempty"`
	/*
	  Encryption scope.
	*/
	EncryptionScope []EncryptionScopeInfo `json:"encryptionScope,omitempty"`

	FaultToleranceState *FaultToleranceState `json:"faultToleranceState,omitempty"`
	/*
	  Hypervisor type.
	*/
	HypervisorTypes []HypervisorType `json:"hypervisorTypes,omitempty"`
	/*
	  Cluster incarnation Id. This is part of payload for cluster update operation only.
	*/
	IncarnationId *int64 `json:"incarnationId,omitempty"`
	/*
	  Indicates whether the release is categorized as Long-term or not.
	*/
	IsLts *bool `json:"isLts,omitempty"`

	OperationMode *OperationMode `json:"operationMode,omitempty"`
	/*
	  Indicates whether the password ssh into the cluster is enabled or not.
	*/
	PasswordRemoteLoginEnabled *bool `json:"passwordRemoteLoginEnabled,omitempty"`
	/*
	  Redundancy factor of a cluster. This is part of payload for both cluster create & update operations.
	*/
	RedundancyFactor *int64 `json:"redundancyFactor,omitempty"`
	/*
	  Remote support status.
	*/
	RemoteSupport *bool `json:"remoteSupport,omitempty"`
	/*
	  Time zone on a cluster.
	*/
	Timezone *string `json:"timezone,omitempty"`
}

func NewClusterConfigReference() *ClusterConfigReference {
	p := new(ClusterConfigReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ClusterConfigReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Cluster function. This is part of payload for cluster create operation only (allowed enum values for creation are AOS, ONE_NODE & TWO_NODE only).
*/
type ClusterFunctionRef int

const (
	CLUSTERFUNCTIONREF_UNKNOWN            ClusterFunctionRef = 0
	CLUSTERFUNCTIONREF_REDACTED           ClusterFunctionRef = 1
	CLUSTERFUNCTIONREF_AOS                ClusterFunctionRef = 2
	CLUSTERFUNCTIONREF_PRISM_CENTRAL      ClusterFunctionRef = 3
	CLUSTERFUNCTIONREF_CLOUD_DATA_GATEWAY ClusterFunctionRef = 4
	CLUSTERFUNCTIONREF_AFS                ClusterFunctionRef = 5
	CLUSTERFUNCTIONREF_ONE_NODE           ClusterFunctionRef = 6
	CLUSTERFUNCTIONREF_TWO_NODE           ClusterFunctionRef = 7
	CLUSTERFUNCTIONREF_ANALYTICS_PLATFORM ClusterFunctionRef = 8
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ClusterFunctionRef) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"AOS",
		"PRISM_CENTRAL",
		"CLOUD_DATA_GATEWAY",
		"AFS",
		"ONE_NODE",
		"TWO_NODE",
		"ANALYTICS_PLATFORM",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ClusterFunctionRef) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"AOS",
		"PRISM_CENTRAL",
		"CLOUD_DATA_GATEWAY",
		"AFS",
		"ONE_NODE",
		"TWO_NODE",
		"ANALYTICS_PLATFORM",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ClusterFunctionRef) index(name string) ClusterFunctionRef {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"AOS",
		"PRISM_CENTRAL",
		"CLOUD_DATA_GATEWAY",
		"AFS",
		"ONE_NODE",
		"TWO_NODE",
		"ANALYTICS_PLATFORM",
	}
	for idx := range names {
		if names[idx] == name {
			return ClusterFunctionRef(idx)
		}
	}
	return CLUSTERFUNCTIONREF_UNKNOWN
}

func (e *ClusterFunctionRef) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ClusterFunctionRef:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ClusterFunctionRef) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ClusterFunctionRef) Ref() *ClusterFunctionRef {
	return &e
}

/*
Network details of a cluster.
*/
type ClusterNetworkReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Backplane *BackplaneNetworkParams `json:"backplane,omitempty"`

	ExternalAddress *import4.IPAddress `json:"externalAddress,omitempty"`

	ExternalDataServiceIp *import4.IPAddress `json:"externalDataServiceIp,omitempty"`
	/*
	  Cluster external subnet address.
	*/
	ExternalSubnet *string `json:"externalSubnet,omitempty"`
	/*
	  Cluster fully qualified domain name. This is part of payload for cluster update operation only.
	*/
	Fqdn *string `json:"fqdn,omitempty"`
	/*
	  Cluster internal subnet address.
	*/
	InternalSubnet *string `json:"internalSubnet,omitempty"`

	KeyManagementServerType *KeyManagementServerType `json:"keyManagementServerType,omitempty"`

	ManagementServer *ManagementServerRef `json:"managementServer,omitempty"`

	MasqueradingIp *import4.IPAddress `json:"masqueradingIp,omitempty"`
	/*
	  The port to connect to the cluster when using masquerading IP.
	*/
	MasqueradingPort *int `json:"masqueradingPort,omitempty"`
	/*
	  List of name servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
	*/
	NameServerIpList []import4.IPAddressOrFQDN `json:"nameServerIpList,omitempty"`
	/*
	  NFS subnet whitelist addresses. This is part of payload for cluster update operation only.
	*/
	NfsSubnetWhitelist []string `json:"nfsSubnetWhitelist,omitempty"`
	/*
	  List of NTP servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
	*/
	NtpServerIpList []import4.IPAddressOrFQDN `json:"ntpServerIpList,omitempty"`

	SmtpServer *SmtpServerRef `json:"smtpServer,omitempty"`
}

func NewClusterNetworkReference() *ClusterNetworkReference {
	p := new(ClusterNetworkReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ClusterNetworkReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Cluster reference for an entity.
*/
type ClusterReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Cluster name. This is part of payload for both cluster create & update operations.
	*/
	Name *string `json:"name,omitempty"`
	/*
	  Cluster UUID.
	*/
	Uuid *string `json:"uuid,omitempty"`
}

func NewClusterReference() *ClusterReference {
	p := new(ClusterReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ClusterReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Fault tolerance information of a component.
*/
type ComponentFaultTolerance struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	DetailMessage *ToleranceMessage `json:"detailMessage,omitempty"`
	/*
	  Time of last update.
	*/
	LastUpdatesSecs *time.Time `json:"lastUpdatesSecs,omitempty"`
	/*
	  Maximum fault tolerance.
	*/
	MaxFaultsTolerated *int `json:"maxFaultsTolerated,omitempty"`

	Type *ComponentType `json:"type,omitempty"`
	/*
	  Indicates whether the tolerance computation is in progress or not.
	*/
	UnderComputation *bool `json:"underComputation,omitempty"`
}

func NewComponentFaultTolerance() *ComponentFaultTolerance {
	p := new(ComponentFaultTolerance)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ComponentFaultTolerance"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Type of component.
*/
type ComponentType int

const (
	COMPONENTTYPE_UNKNOWN                 ComponentType = 0
	COMPONENTTYPE_REDACTED                ComponentType = 1
	COMPONENTTYPE_EXTENT_GROUP_REPLICAS   ComponentType = 2
	COMPONENTTYPE_OPLOG_EPISODES          ComponentType = 3
	COMPONENTTYPE_CASSANDRA_RING          ComponentType = 4
	COMPONENTTYPE_ZOOKEPER_INSTANCES      ComponentType = 5
	COMPONENTTYPE_FREE_SPACE              ComponentType = 6
	COMPONENTTYPE_STATIC_CONFIG           ComponentType = 7
	COMPONENTTYPE_ERASURE_CODE_STRIP_SIZE ComponentType = 8
	COMPONENTTYPE_STARGATE_HEALTH         ComponentType = 9
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ComponentType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"EXTENT_GROUP_REPLICAS",
		"OPLOG_EPISODES",
		"CASSANDRA_RING",
		"ZOOKEPER_INSTANCES",
		"FREE_SPACE",
		"STATIC_CONFIG",
		"ERASURE_CODE_STRIP_SIZE",
		"STARGATE_HEALTH",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ComponentType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"EXTENT_GROUP_REPLICAS",
		"OPLOG_EPISODES",
		"CASSANDRA_RING",
		"ZOOKEPER_INSTANCES",
		"FREE_SPACE",
		"STATIC_CONFIG",
		"ERASURE_CODE_STRIP_SIZE",
		"STARGATE_HEALTH",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ComponentType) index(name string) ComponentType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"EXTENT_GROUP_REPLICAS",
		"OPLOG_EPISODES",
		"CASSANDRA_RING",
		"ZOOKEPER_INSTANCES",
		"FREE_SPACE",
		"STATIC_CONFIG",
		"ERASURE_CODE_STRIP_SIZE",
		"STARGATE_HEALTH",
	}
	for idx := range names {
		if names[idx] == name {
			return ComponentType(idx)
		}
	}
	return COMPONENTTYPE_UNKNOWN
}

func (e *ComponentType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ComponentType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ComponentType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ComponentType) Ref() *ComponentType {
	return &e
}

/*
Compute node details.
*/
type ComputeNodeItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rackable unit Id in which node resides.
	*/
	BlockId *string `json:"blockId,omitempty"`
	/*
	  List of objects containing digital_certificate_base64 and key_management_server_uuid fields for key management server.
	*/
	DigitalCertificateMapList []DigitalCertificateMapReference `json:"digitalCertificateMapList,omitempty"`
	/*
	  Name of the host.
	*/
	HypervisorHostname *string `json:"hypervisorHostname,omitempty"`

	HypervisorIp *import4.IPAddress `json:"hypervisorIp,omitempty"`

	IpmiIp *import4.IPAddress `json:"ipmiIp,omitempty"`
	/*
	  Rackable unit model type.
	*/
	Model *string `json:"model,omitempty"`
	/*
	  Position of a node in a rackable unit.
	*/
	NodePosition *string `json:"nodePosition,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
}

func NewComputeNodeItem() *ComputeNodeItem {
	p := new(ComputeNodeItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ComputeNodeItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Config parameters.
*/
type ConfigParams struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Hyperv *HypervCredentials `json:"hyperv,omitempty"`
	/*
	  Indicates whether the node is compute only or not.
	*/
	IsComputeOnly *bool `json:"isComputeOnly,omitempty"`
	/*
	  Indicates if node is compatible or not.
	*/
	IsNosCompatible *bool `json:"isNosCompatible,omitempty"`
	/*
	  Indicates whether the node is marked to be never schedulable or not.
	*/
	NeverScheduleable *bool `json:"neverScheduleable,omitempty"`
	/*
	  Indicates if node discovery need to be skipped or not.
	*/
	SkipDiscovery *bool `json:"skipDiscovery,omitempty"`
	/*
	  Indicates if node imaging needs to be skipped or not.
	*/
	SkipImaging *bool `json:"skipImaging,omitempty"`
	/*
	  Target hypervisor.
	*/
	TargetHypervisor *string `json:"targetHypervisor,omitempty"`
	/*
	  Indicates if rack awareness needs to be validated or not.
	*/
	ValidateRackAwareness *bool `json:"validateRackAwareness,omitempty"`
}

func NewConfigParams() *ConfigParams {
	p := new(ConfigParams)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ConfigParams"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Host entity with its attributes.
*/
type ControllerVmReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	BackplaneAddress *import4.IPAddress `json:"backplaneAddress,omitempty"`

	ExternalAddress *import4.IPAddress `json:"externalAddress,omitempty"`
	/*
	  Controller VM Id.
	*/
	Id *int64 `json:"id,omitempty"`

	Ipmi *IpmiReference `json:"ipmi,omitempty"`
	/*
	  Maintenance mode status.
	*/
	MaintenanceMode *bool `json:"maintenanceMode,omitempty"`

	NatIp *import4.IPAddress `json:"natIp,omitempty"`
	/*
	  NAT port.
	*/
	NatPort *int `json:"natPort,omitempty"`
	/*
	  Rackable unit UUID.
	*/
	RackableUnitUuid *string `json:"rackableUnitUuid,omitempty"`
	/*
	  RDMA backplane address.
	*/
	RdmaBackplaneAddress []import4.IPAddress `json:"rdmaBackplaneAddress,omitempty"`
}

func NewControllerVmReference() *ControllerVmReference {
	p := new(ControllerVmReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ControllerVmReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters Post operation
*/
type CreateClusterApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreateClusterApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreateClusterApiResponse() *CreateClusterApiResponse {
	p := new(CreateClusterApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.CreateClusterApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreateClusterApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreateClusterApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreateClusterApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/rsyslog-servers Post operation
*/
type CreateRsyslogServerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreateRsyslogServerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreateRsyslogServerApiResponse() *CreateRsyslogServerApiResponse {
	p := new(CreateRsyslogServerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.CreateRsyslogServerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreateRsyslogServerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreateRsyslogServerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreateRsyslogServerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/snmp/traps Post operation
*/
type CreateSnmpTrapApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreateSnmpTrapApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreateSnmpTrapApiResponse() *CreateSnmpTrapApiResponse {
	p := new(CreateSnmpTrapApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.CreateSnmpTrapApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreateSnmpTrapApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreateSnmpTrapApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreateSnmpTrapApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/snmp/users Post operation
*/
type CreateSnmpUserApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreateSnmpUserApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreateSnmpUserApiResponse() *CreateSnmpUserApiResponse {
	p := new(CreateSnmpUserApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.CreateSnmpUserApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreateSnmpUserApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreateSnmpUserApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreateSnmpUserApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/storage-containers Post operation
*/
type CreateStorageContainerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfCreateStorageContainerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewCreateStorageContainerApiResponse() *CreateStorageContainerApiResponse {
	p := new(CreateStorageContainerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.CreateStorageContainerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *CreateStorageContainerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *CreateStorageContainerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfCreateStorageContainerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

type DataStore struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Maximum physical capacity of the Storage Container in bytes.
	*/
	CapacityBytes *int64 `json:"capacityBytes,omitempty"`
	/*
	  extId of the Storage Container.
	*/
	ContainerExtId *string `json:"containerExtId,omitempty"`
	/*
	  Name of the Storage Container. Note that the name of Storage Container should be unique per cluster.
	*/
	ContainerName *string `json:"containerName,omitempty"`
	/*
	  Name of the Data Store.
	*/
	DatastoreName *string `json:"datastoreName,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  The free space in the Data Store.
	*/
	FreeSpaceBytes *int64 `json:"freeSpaceBytes,omitempty"`
	/*
	  ExtId of the host for Data Store.
	*/
	HostExtId *string `json:"hostExtId,omitempty"`
	/*
	  Host IP for Data Store.
	*/
	HostIpAddress *string `json:"hostIpAddress,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/*
	  List of VM names in the Data Store.
	*/
	VmNames []string `json:"vmNames,omitempty"`
}

func NewDataStore() *DataStore {
	p := new(DataStore)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DataStore"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type DataStoreMount struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Name of the Storage Container. Note that the name of Storage Container should be unique per cluster.
	*/
	ContainerName *string `json:"containerName"`
	/*
	  Name of the Data Store.
	*/
	DatastoreName *string `json:"datastoreName,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  The UUIDs of the nodes where the NFS Data Store has to be created.
	*/
	NodeExtIds []string `json:"nodeExtIds,omitempty"`
	/*
	  Indicates if the host system has only read-only access to the NFS share.
	*/
	ReadOnly *bool `json:"readOnly,omitempty"`
	/*
	  The target path on which to mount the NFS Data Store.
	*/
	TargetPath *string `json:"targetPath,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func (p *DataStoreMount) MarshalJSON() ([]byte, error) {
	type DataStoreMountProxy DataStoreMount
	return json.Marshal(struct {
		*DataStoreMountProxy
		ContainerName *string `json:"containerName,omitempty"`
	}{
		DataStoreMountProxy: (*DataStoreMountProxy)(p),
		ContainerName:       p.ContainerName,
	})
}

func NewDataStoreMount() *DataStoreMount {
	p := new(DataStoreMount)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DataStoreMount"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type DataStoreUnmount struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Name of the Data Store.
	*/
	DatastoreName *string `json:"datastoreName"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  The UUIDs of the nodes where the NFS Data Store has to be created.
	*/
	NodeExtIds []string `json:"nodeExtIds,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func (p *DataStoreUnmount) MarshalJSON() ([]byte, error) {
	type DataStoreUnmountProxy DataStoreUnmount
	return json.Marshal(struct {
		*DataStoreUnmountProxy
		DatastoreName *string `json:"datastoreName,omitempty"`
	}{
		DataStoreUnmountProxy: (*DataStoreUnmountProxy)(p),
		DatastoreName:         p.DatastoreName,
	})
}

func NewDataStoreUnmount() *DataStoreUnmount {
	p := new(DataStoreUnmount)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DataStoreUnmount"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId} Delete operation
*/
type DeleteClusterApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeleteClusterApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeleteClusterApiResponse() *DeleteClusterApiResponse {
	p := new(DeleteClusterApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DeleteClusterApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeleteClusterApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeleteClusterApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeleteClusterApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/rsyslog-servers/{extId} Delete operation
*/
type DeleteRsyslogServerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeleteRsyslogServerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeleteRsyslogServerApiResponse() *DeleteRsyslogServerApiResponse {
	p := new(DeleteRsyslogServerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DeleteRsyslogServerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeleteRsyslogServerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeleteRsyslogServerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeleteRsyslogServerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/snmp/traps/{extId} Delete operation
*/
type DeleteSnmpTrapApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeleteSnmpTrapApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeleteSnmpTrapApiResponse() *DeleteSnmpTrapApiResponse {
	p := new(DeleteSnmpTrapApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DeleteSnmpTrapApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeleteSnmpTrapApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeleteSnmpTrapApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeleteSnmpTrapApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/snmp/users/{extId} Delete operation
*/
type DeleteSnmpUserApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeleteSnmpUserApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeleteSnmpUserApiResponse() *DeleteSnmpUserApiResponse {
	p := new(DeleteSnmpUserApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DeleteSnmpUserApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeleteSnmpUserApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeleteSnmpUserApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeleteSnmpUserApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/storage-containers/{extId} Delete operation
*/
type DeleteStorageContainerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDeleteStorageContainerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDeleteStorageContainerApiResponse() *DeleteStorageContainerApiResponse {
	p := new(DeleteStorageContainerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DeleteStorageContainerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DeleteStorageContainerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DeleteStorageContainerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDeleteStorageContainerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Object containing digital_certificate_base64 and key_management_server_uuid fields for key management server.
*/
type DigitalCertificateMapReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Field containing digital_certificate_base64 and key_management_server_uuid for key management server.
	*/
	Key *string `json:"key,omitempty"`
	/*
	  Value for the fields digital_certificate_base64 and key_management_server_uuid for key management server.
	*/
	Value *string `json:"value,omitempty"`
}

func NewDigitalCertificateMapReference() *DigitalCertificateMapReference {
	p := new(DigitalCertificateMapReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DigitalCertificateMapReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/$actions/discover-unconfigured-nodes Post operation
*/
type DiscoverUnconfiguredNodesApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDiscoverUnconfiguredNodesApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDiscoverUnconfiguredNodesApiResponse() *DiscoverUnconfiguredNodesApiResponse {
	p := new(DiscoverUnconfiguredNodesApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DiscoverUnconfiguredNodesApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DiscoverUnconfiguredNodesApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DiscoverUnconfiguredNodesApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDiscoverUnconfiguredNodesApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
A model that represents Disk resources.
*/
type Disk struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The external identifier of the Disk Cluster.
	*/
	ClusterExtId *string `json:"clusterExtId,omitempty"`
	/*
	  Cluster name of the Disk it belongs to.
	*/
	ClusterName *string `json:"clusterName,omitempty"`

	CvmIpAddress *import4.IPAddress `json:"cvmIpAddress,omitempty"`

	DiskAdvanceConfig *DiskAdvanceConfig `json:"diskAdvanceConfig,omitempty"`
	/*
	  Size of the Disk in bytes.
	*/
	DiskSizeBytes *int64 `json:"diskSizeBytes,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Represents the current firmware version.
	*/
	FirmwareVersion *string `json:"firmwareVersion,omitempty"`
	/*
	  Host name of the Disk it belongs to.
	*/
	HostName *string `json:"hostName,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Indicates the location of Disk in a node.
	*/
	Location *int64 `json:"location,omitempty"`
	/*
	  Represents the Disk model.
	*/
	Model *string `json:"model,omitempty"`
	/*
	  Represents the mount path of the Disk.
	*/
	MountPath *string `json:"mountPath,omitempty"`
	/*
	  The external identifier of the node.
	*/
	NodeExtId *string `json:"nodeExtId,omitempty"`

	NodeIpAddress *import4.IPAddress `json:"nodeIpAddress,omitempty"`
	/*
	  Indicates the PCIe path of NVMe devices.
	*/
	NvmePciePath *string `json:"nvmePciePath,omitempty"`
	/*
	  Physical capacity of the Disk in bytes.
	*/
	PhysicalCapacityBytes *int64 `json:"physicalCapacityBytes,omitempty"`
	/*
	  Represents the Disk serial number.
	*/
	SerialNumber *string `json:"serialNumber,omitempty"`
	/*
	  ServiceVMId of the node.
	*/
	ServiceVMId *string `json:"serviceVMId,omitempty"`

	Status *DiskStatus `json:"status,omitempty"`
	/*
	  The external identifier of storage pool.
	*/
	StoragePoolExtId *string `json:"storagePoolExtId,omitempty"`

	StorageTier *StorageTier `json:"storageTier,omitempty"`
	/*
	  Represents the target firmware version.
	*/
	TargetFirmwareVersion *string `json:"targetFirmwareVersion,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/*
	  Represents the Disk vendor.
	*/
	Vendor *string `json:"vendor,omitempty"`
}

func NewDisk() *Disk {
	p := new(Disk)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.Disk"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
A model that represents Disk advance config properties.
*/
type DiskAdvanceConfig struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Indicates if the Disk is boot only and no Disk operation to be run on it.
	*/
	HasBootPartitionsOnly *bool `json:"hasBootPartitionsOnly,omitempty"`
	/*
	  Indicates if Disk is a boot Disk.
	*/
	IsBootDisk *bool `json:"isBootDisk,omitempty"`
	/*
	  Indicates if data migration is completed for the Disk.
	*/
	IsDataMigrated *bool `json:"isDataMigrated,omitempty"`
	/*
	  Indicates the Disk diagnostic info along with device related statics is present.
	*/
	IsDiagnosticInfoAvailable *bool `json:"isDiagnosticInfoAvailable,omitempty"`
	/*
	  Indicates Disk error is seen on the Disk in kernel logs or not.
	*/
	IsErrorFoundInLog *bool `json:"isErrorFoundInLog,omitempty"`
	/*
	  Indicates if the Disk is marked for removal.
	*/
	IsMarkedForRemoval *bool `json:"isMarkedForRemoval,omitempty"`
	/*
	  Indicates if the Disk is mounted.
	*/
	IsMounted *bool `json:"isMounted,omitempty"`
	/*
	  Indicates whether Disk is online or offline.
	*/
	IsOnline *bool `json:"isOnline,omitempty"`
	/*
	  Indicates if the Disk is password protected.
	*/
	IsPasswordProtected *bool `json:"isPasswordProtected,omitempty"`
	/*
	  Indicates if diagnostics are running on the Disk.
	*/
	IsPlannedOutage *bool `json:"isPlannedOutage,omitempty"`
	/*
	  Indicates whether the Disk has self encryption enabled.
	*/
	IsSelfEncryptingDrive *bool `json:"isSelfEncryptingDrive,omitempty"`
	/*
	  Indicates if the NVMe Disk is self managed and no host/CVM reboot is required.
	*/
	IsSelfManagedNvme *bool `json:"isSelfManagedNvme,omitempty"`
	/*
	  Indicates if NVMe device is managed by storage performance development kit(SPDK).
	*/
	IsSpdkManaged *bool `json:"isSpdkManaged,omitempty"`
	/*
	  Indicates if the Disk is suspected to be unhealthy.
	*/
	IsSuspectedUnhealthy *bool `json:"isSuspectedUnhealthy,omitempty"`
	/*
	  Indicates if the Disk is under diagnosis.
	*/
	IsUnderDiagnosis *bool `json:"isUnderDiagnosis,omitempty"`
	/*
	  Indicates if the Disk is unhealthy.
	*/
	IsUnhealthy *bool `json:"isUnhealthy,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewDiskAdvanceConfig() *DiskAdvanceConfig {
	p := new(DiskAdvanceConfig)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DiskAdvanceConfig"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Disk details attached to a host.
*/
type DiskReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Disk mount path.
	*/
	MountPath *string `json:"mountPath,omitempty"`
	/*
	  Disk serial Id.
	*/
	SerialId *string `json:"serialId,omitempty"`
	/*
	  Disk size.
	*/
	SizeInBytes *int64 `json:"sizeInBytes,omitempty"`

	StorageTier *StorageTierReference `json:"storageTier,omitempty"`
	/*
	  Disk UUID.
	*/
	Uuid *string `json:"uuid,omitempty"`
}

func NewDiskReference() *DiskReference {
	p := new(DiskReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DiskReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Indicates the status of the Disk.
*/
type DiskStatus int

const (
	DISKSTATUS_UNKNOWN                               DiskStatus = 0
	DISKSTATUS_REDACTED                              DiskStatus = 1
	DISKSTATUS_NORMAL                                DiskStatus = 2
	DISKSTATUS_MARKED_FOR_REMOVAL_BUT_NOT_DETACHABLE DiskStatus = 3
	DISKSTATUS_DETACHABLE                            DiskStatus = 4
	DISKSTATUS_DATA_MIGRATION_INITIATED              DiskStatus = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *DiskStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"MARKED_FOR_REMOVAL_BUT_NOT_DETACHABLE",
		"DETACHABLE",
		"DATA_MIGRATION_INITIATED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e DiskStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"MARKED_FOR_REMOVAL_BUT_NOT_DETACHABLE",
		"DETACHABLE",
		"DATA_MIGRATION_INITIATED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *DiskStatus) index(name string) DiskStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"MARKED_FOR_REMOVAL_BUT_NOT_DETACHABLE",
		"DETACHABLE",
		"DATA_MIGRATION_INITIATED",
	}
	for idx := range names {
		if names[idx] == name {
			return DiskStatus(idx)
		}
	}
	return DISKSTATUS_UNKNOWN
}

func (e *DiskStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for DiskStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *DiskStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e DiskStatus) Ref() *DiskStatus {
	return &e
}

/*
Domain awareness level corresponds to unit of cluster group. This is part of payload for both cluster create & update operations.
*/
type DomainAwarenessLevel int

const (
	DOMAINAWARENESSLEVEL_UNKNOWN  DomainAwarenessLevel = 0
	DOMAINAWARENESSLEVEL_REDACTED DomainAwarenessLevel = 1
	DOMAINAWARENESSLEVEL_NODE     DomainAwarenessLevel = 2
	DOMAINAWARENESSLEVEL_BLOCK    DomainAwarenessLevel = 3
	DOMAINAWARENESSLEVEL_RACK     DomainAwarenessLevel = 4
	DOMAINAWARENESSLEVEL_DISK     DomainAwarenessLevel = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *DomainAwarenessLevel) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NODE",
		"BLOCK",
		"RACK",
		"DISK",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e DomainAwarenessLevel) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NODE",
		"BLOCK",
		"RACK",
		"DISK",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *DomainAwarenessLevel) index(name string) DomainAwarenessLevel {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NODE",
		"BLOCK",
		"RACK",
		"DISK",
	}
	for idx := range names {
		if names[idx] == name {
			return DomainAwarenessLevel(idx)
		}
	}
	return DOMAINAWARENESSLEVEL_UNKNOWN
}

func (e *DomainAwarenessLevel) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for DomainAwarenessLevel:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *DomainAwarenessLevel) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e DomainAwarenessLevel) Ref() *DomainAwarenessLevel {
	return &e
}

/*
Domain fault tolerance configuration.
*/
type DomainFaultTolerance struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of components in a domain.
	*/
	ComponentStatus []ComponentFaultTolerance `json:"componentStatus,omitempty"`

	Type *DomainType `json:"type,omitempty"`
}

func NewDomainFaultTolerance() *DomainFaultTolerance {
	p := new(DomainFaultTolerance)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.DomainFaultTolerance"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Type of domain
*/
type DomainType int

const (
	DOMAINTYPE_UNKNOWN       DomainType = 0
	DOMAINTYPE_REDACTED      DomainType = 1
	DOMAINTYPE_CUSTOM        DomainType = 2
	DOMAINTYPE_DISK          DomainType = 3
	DOMAINTYPE_NODE          DomainType = 4
	DOMAINTYPE_RACKABLE_UNIT DomainType = 5
	DOMAINTYPE_RACK          DomainType = 6
	DOMAINTYPE_CLUSTER       DomainType = 7
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *DomainType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CUSTOM",
		"DISK",
		"NODE",
		"RACKABLE_UNIT",
		"RACK",
		"CLUSTER",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e DomainType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CUSTOM",
		"DISK",
		"NODE",
		"RACKABLE_UNIT",
		"RACK",
		"CLUSTER",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *DomainType) index(name string) DomainType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CUSTOM",
		"DISK",
		"NODE",
		"RACKABLE_UNIT",
		"RACK",
		"CLUSTER",
	}
	for idx := range names {
		if names[idx] == name {
			return DomainType(idx)
		}
	}
	return DOMAINTYPE_UNKNOWN
}

func (e *DomainType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for DomainType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *DomainType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e DomainType) Ref() *DomainType {
	return &e
}

/*
Encryption option.
*/
type EncryptionOptionInfo int

const (
	ENCRYPTIONOPTIONINFO_UNKNOWN               EncryptionOptionInfo = 0
	ENCRYPTIONOPTIONINFO_REDACTED              EncryptionOptionInfo = 1
	ENCRYPTIONOPTIONINFO_SOFTWARE              EncryptionOptionInfo = 2
	ENCRYPTIONOPTIONINFO_HARDWARE              EncryptionOptionInfo = 3
	ENCRYPTIONOPTIONINFO_SOFTWARE_AND_HARDWARE EncryptionOptionInfo = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *EncryptionOptionInfo) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SOFTWARE",
		"HARDWARE",
		"SOFTWARE_AND_HARDWARE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e EncryptionOptionInfo) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SOFTWARE",
		"HARDWARE",
		"SOFTWARE_AND_HARDWARE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *EncryptionOptionInfo) index(name string) EncryptionOptionInfo {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SOFTWARE",
		"HARDWARE",
		"SOFTWARE_AND_HARDWARE",
	}
	for idx := range names {
		if names[idx] == name {
			return EncryptionOptionInfo(idx)
		}
	}
	return ENCRYPTIONOPTIONINFO_UNKNOWN
}

func (e *EncryptionOptionInfo) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for EncryptionOptionInfo:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *EncryptionOptionInfo) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e EncryptionOptionInfo) Ref() *EncryptionOptionInfo {
	return &e
}

/*
Encryption scope.
*/
type EncryptionScopeInfo int

const (
	ENCRYPTIONSCOPEINFO_UNKNOWN   EncryptionScopeInfo = 0
	ENCRYPTIONSCOPEINFO_REDACTED  EncryptionScopeInfo = 1
	ENCRYPTIONSCOPEINFO_CLUSTER   EncryptionScopeInfo = 2
	ENCRYPTIONSCOPEINFO_CONTAINER EncryptionScopeInfo = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *EncryptionScopeInfo) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CLUSTER",
		"CONTAINER",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e EncryptionScopeInfo) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CLUSTER",
		"CONTAINER",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *EncryptionScopeInfo) index(name string) EncryptionScopeInfo {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CLUSTER",
		"CONTAINER",
	}
	for idx := range names {
		if names[idx] == name {
			return EncryptionScopeInfo(idx)
		}
	}
	return ENCRYPTIONSCOPEINFO_UNKNOWN
}

func (e *EncryptionScopeInfo) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for EncryptionScopeInfo:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *EncryptionScopeInfo) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e EncryptionScopeInfo) Ref() *EncryptionScopeInfo {
	return &e
}

/*
Encryption in transit Status.
*/
type EncryptionStatus int

const (
	ENCRYPTIONSTATUS_UNKNOWN  EncryptionStatus = 0
	ENCRYPTIONSTATUS_REDACTED EncryptionStatus = 1
	ENCRYPTIONSTATUS_ENABLED  EncryptionStatus = 2
	ENCRYPTIONSTATUS_DISABLED EncryptionStatus = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *EncryptionStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ENABLED",
		"DISABLED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e EncryptionStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ENABLED",
		"DISABLED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *EncryptionStatus) index(name string) EncryptionStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ENABLED",
		"DISABLED",
	}
	for idx := range names {
		if names[idx] == name {
			return EncryptionStatus(idx)
		}
	}
	return ENCRYPTIONSTATUS_UNKNOWN
}

func (e *EncryptionStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for EncryptionStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *EncryptionStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e EncryptionStatus) Ref() *EncryptionStatus {
	return &e
}

/*
Indicates the current status value for Erasure Coding for the Container.
*/
type ErasureCodeStatus int

const (
	ERASURECODESTATUS_UNKNOWN  ErasureCodeStatus = 0
	ERASURECODESTATUS_REDACTED ErasureCodeStatus = 1
	ERASURECODESTATUS_NONE     ErasureCodeStatus = 2
	ERASURECODESTATUS_OFF      ErasureCodeStatus = 3
	ERASURECODESTATUS_ON       ErasureCodeStatus = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ErasureCodeStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"ON",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ErasureCodeStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"ON",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ErasureCodeStatus) index(name string) ErasureCodeStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"ON",
	}
	for idx := range names {
		if names[idx] == name {
			return ErasureCodeStatus(idx)
		}
	}
	return ERASURECODESTATUS_UNKNOWN
}

func (e *ErasureCodeStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ErasureCodeStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ErasureCodeStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ErasureCodeStatus) Ref() *ErasureCodeStatus {
	return &e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/$actions/expand-cluster Post operation
*/
type ExpandClusterApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfExpandClusterApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewExpandClusterApiResponse() *ExpandClusterApiResponse {
	p := new(ExpandClusterApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ExpandClusterApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ExpandClusterApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ExpandClusterApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfExpandClusterApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Property of the node to be added.
*/
type ExpandClusterParams struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	ConfigParams *ConfigParams `json:"configParams,omitempty"`

	NodeParams *NodeParam `json:"nodeParams"`
	/*
	  Indicates if node addition can be skipped.
	*/
	SkipAddNode *bool `json:"skipAddNode,omitempty"`
	/*
	  Indicates if pre-expand checks can be skipped for node addition.
	*/
	SkipPreExpandChecks *bool `json:"skipPreExpandChecks,omitempty"`
}

func (p *ExpandClusterParams) MarshalJSON() ([]byte, error) {
	type ExpandClusterParamsProxy ExpandClusterParams
	return json.Marshal(struct {
		*ExpandClusterParamsProxy
		NodeParams *NodeParam `json:"nodeParams,omitempty"`
	}{
		ExpandClusterParamsProxy: (*ExpandClusterParamsProxy)(p),
		NodeParams:               p.NodeParams,
	})
}

func NewExpandClusterParams() *ExpandClusterParams {
	p := new(ExpandClusterParams)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ExpandClusterParams"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Fault tolerant state of cluster.
*/
type FaultToleranceState struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Maximum fault tolerance that is supported currently.
	*/
	CurrentMaxFaultTolerance *int `json:"currentMaxFaultTolerance,omitempty"`
	/*
	  Maximum fault tolerance desired.
	*/
	DesiredMaxFaultTolerance *int `json:"desiredMaxFaultTolerance,omitempty"`

	DomainAwarenessLevel *DomainAwarenessLevel `json:"domainAwarenessLevel,omitempty"`
}

func NewFaultToleranceState() *FaultToleranceState {
	p := new(FaultToleranceState)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.FaultToleranceState"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/$actions/fetch-node-networking-details Post operation
*/
type FetchNodeNetworkingDetailsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfFetchNodeNetworkingDetailsApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewFetchNodeNetworkingDetailsApiResponse() *FetchNodeNetworkingDetailsApiResponse {
	p := new(FetchNodeNetworkingDetailsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.FetchNodeNetworkingDetailsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *FetchNodeNetworkingDetailsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *FetchNodeNetworkingDetailsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfFetchNodeNetworkingDetailsApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/tasks/{extId}/$actions/fetch-task-response Post operation
*/
type FetchTaskApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfFetchTaskApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewFetchTaskApiResponse() *FetchTaskApiResponse {
	p := new(FetchTaskApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.FetchTaskApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *FetchTaskApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *FetchTaskApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfFetchTaskApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId} Get operation
*/
type GetClusterApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetClusterApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetClusterApiResponse() *GetClusterApiResponse {
	p := new(GetClusterApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetClusterApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetClusterApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetClusterApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetClusterApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/disks/{extId} Get operation
*/
type GetDiskApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetDiskApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetDiskApiResponse() *GetDiskApiResponse {
	p := new(GetDiskApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetDiskApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetDiskApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetDiskApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetDiskApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/fault-tolerance-status Get operation
*/
type GetFaultToleranceStatusByClusterIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetFaultToleranceStatusByClusterIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetFaultToleranceStatusByClusterIdApiResponse() *GetFaultToleranceStatusByClusterIdApiResponse {
	p := new(GetFaultToleranceStatusByClusterIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetFaultToleranceStatusByClusterIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetFaultToleranceStatusByClusterIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetFaultToleranceStatusByClusterIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetFaultToleranceStatusByClusterIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/hosts/{extId} Get operation
*/
type GetHostApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetHostApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetHostApiResponse() *GetHostApiResponse {
	p := new(GetHostApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetHostApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetHostApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetHostApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetHostApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/hosts/{hostExtId}/host-nics/{extId} Get operation
*/
type GetHostNicApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetHostNicApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetHostNicApiResponse() *GetHostNicApiResponse {
	p := new(GetHostNicApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetHostNicApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetHostNicApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetHostNicApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetHostNicApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/rackable-units/{extId} Get operation
*/
type GetRackableUnitApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetRackableUnitApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetRackableUnitApiResponse() *GetRackableUnitApiResponse {
	p := new(GetRackableUnitApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetRackableUnitApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetRackableUnitApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetRackableUnitApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetRackableUnitApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/rsyslog-servers/{extId} Get operation
*/
type GetRsyslogServerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetRsyslogServerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetRsyslogServerApiResponse() *GetRsyslogServerApiResponse {
	p := new(GetRsyslogServerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetRsyslogServerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetRsyslogServerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetRsyslogServerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetRsyslogServerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/snmp Get operation
*/
type GetSnmpConfigByClusterIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetSnmpConfigByClusterIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetSnmpConfigByClusterIdApiResponse() *GetSnmpConfigByClusterIdApiResponse {
	p := new(GetSnmpConfigByClusterIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetSnmpConfigByClusterIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetSnmpConfigByClusterIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetSnmpConfigByClusterIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetSnmpConfigByClusterIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/snmp/traps/{extId} Get operation
*/
type GetSnmpTrapApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetSnmpTrapApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetSnmpTrapApiResponse() *GetSnmpTrapApiResponse {
	p := new(GetSnmpTrapApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetSnmpTrapApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetSnmpTrapApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetSnmpTrapApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetSnmpTrapApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/snmp/users/{extId} Get operation
*/
type GetSnmpUserApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetSnmpUserApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetSnmpUserApiResponse() *GetSnmpUserApiResponse {
	p := new(GetSnmpUserApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetSnmpUserApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetSnmpUserApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetSnmpUserApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetSnmpUserApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/storage-containers/{extId} Get operation
*/
type GetStorageContainerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetStorageContainerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetStorageContainerApiResponse() *GetStorageContainerApiResponse {
	p := new(GetStorageContainerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetStorageContainerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetStorageContainerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetStorageContainerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetStorageContainerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/vcenter-extensions/{extId} Get operation
*/
type GetVcenterExtensionApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetVcenterExtensionApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetVcenterExtensionApiResponse() *GetVcenterExtensionApiResponse {
	p := new(GetVcenterExtensionApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetVcenterExtensionApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetVcenterExtensionApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetVcenterExtensionApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetVcenterExtensionApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/hosts/{hostExtId}/virtual-nics/{extId} Get operation
*/
type GetVirtualNicApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetVirtualNicApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetVirtualNicApiResponse() *GetVirtualNicApiResponse {
	p := new(GetVirtualNicApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GetVirtualNicApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetVirtualNicApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetVirtualNicApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetVirtualNicApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
GPU configuration details.
*/
type GpuConfig struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  GPU assignable.
	*/
	Assignable *int64 `json:"assignable,omitempty"`
	/*
	  Device Id.
	*/
	DeviceId *string `json:"deviceId,omitempty"`
	/*
	  Device name.
	*/
	DeviceName *string `json:"deviceName,omitempty"`
	/*
	  GPU fraction.
	*/
	Fraction *int64 `json:"fraction,omitempty"`
	/*
	  Frame buffer size in bytes.
	*/
	FrameBufferSizeBytes *int64 `json:"frameBufferSizeBytes,omitempty"`
	/*
	  Guest driver version.
	*/
	GuestDriverVersion *string `json:"guestDriverVersion,omitempty"`
	/*
	  GPU in use.
	*/
	InUse *bool `json:"inUse,omitempty"`
	/*
	  GPU license list.
	*/
	LicenseList []string `json:"licenseList,omitempty"`
	/*
	  Maximum resolution per display heads.
	*/
	MaxResolution *string `json:"maxResolution,omitempty"`

	Mode *GpuMode `json:"mode,omitempty"`
	/*
	  NUMA node.
	*/
	NumaNode *string `json:"numaNode,omitempty"`
	/*
	  Number of virtual display heads.
	*/
	NumberOfVirtualDisplayHeads *int64 `json:"numberOfVirtualDisplayHeads,omitempty"`
	/*
	  SBDF address.
	*/
	Sbdf *string `json:"sbdf,omitempty"`

	Type *GpuType `json:"type,omitempty"`
	/*
	  Vendor name.
	*/
	VendorName *string `json:"vendorName,omitempty"`
}

func NewGpuConfig() *GpuConfig {
	p := new(GpuConfig)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GpuConfig"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
GPU mode.
*/
type GpuMode int

const (
	GPUMODE_UNKNOWN              GpuMode = 0
	GPUMODE_REDACTED             GpuMode = 1
	GPUMODE_UNUSED               GpuMode = 2
	GPUMODE_USED_FOR_PASSTHROUGH GpuMode = 3
	GPUMODE_USED_FOR_VIRTUAL     GpuMode = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *GpuMode) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UNUSED",
		"USED_FOR_PASSTHROUGH",
		"USED_FOR_VIRTUAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e GpuMode) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UNUSED",
		"USED_FOR_PASSTHROUGH",
		"USED_FOR_VIRTUAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *GpuMode) index(name string) GpuMode {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UNUSED",
		"USED_FOR_PASSTHROUGH",
		"USED_FOR_VIRTUAL",
	}
	for idx := range names {
		if names[idx] == name {
			return GpuMode(idx)
		}
	}
	return GPUMODE_UNKNOWN
}

func (e *GpuMode) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for GpuMode:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *GpuMode) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e GpuMode) Ref() *GpuMode {
	return &e
}

/*
GPU Profile.
*/
type GpuProfile struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of UUIDs of virtual machines with an allocated GPU belonging to this profile.
	*/
	AllocatedVmUuids []string `json:"allocatedVmUuids,omitempty"`
	/*
	  Device Id.
	*/
	DeviceId *string `json:"deviceId,omitempty"`

	GpuConfig *GpuConfig `json:"gpuConfig,omitempty"`
}

func NewGpuProfile() *GpuProfile {
	p := new(GpuProfile)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.GpuProfile"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
GPU type.
*/
type GpuType int

const (
	GPUTYPE_UNKNOWN               GpuType = 0
	GPUTYPE_REDACTED              GpuType = 1
	GPUTYPE_PASS_THROUGH_GRAPHICS GpuType = 2
	GPUTYPE_PASS_THROUGH_COMPUTE  GpuType = 3
	GPUTYPE_VIRTUAL               GpuType = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *GpuType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PASS_THROUGH_GRAPHICS",
		"PASS_THROUGH_COMPUTE",
		"VIRTUAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e GpuType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PASS_THROUGH_GRAPHICS",
		"PASS_THROUGH_COMPUTE",
		"VIRTUAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *GpuType) index(name string) GpuType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PASS_THROUGH_GRAPHICS",
		"PASS_THROUGH_COMPUTE",
		"VIRTUAL",
	}
	for idx := range names {
		if names[idx] == name {
			return GpuType(idx)
		}
	}
	return GPUTYPE_UNKNOWN
}

func (e *GpuType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for GpuType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *GpuType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e GpuType) Ref() *GpuType {
	return &e
}

/*
Host entity with its attributes.
*/
type Host struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rackable unit model name.
	*/
	BlockModel *string `json:"blockModel,omitempty"`
	/*
	  Rackable unit serial name.
	*/
	BlockSerial *string `json:"blockSerial,omitempty"`
	/*
	  Boot time in secs.
	*/
	BootTimeUsecs *int64 `json:"bootTimeUsecs,omitempty"`

	Cluster *ClusterReference `json:"cluster,omitempty"`

	ControllerVm *ControllerVmReference `json:"controllerVm,omitempty"`
	/*
	  CPU capacity in Hz.
	*/
	CpuCapacityHz *int64 `json:"cpuCapacityHz,omitempty"`
	/*
	  CPU frequency in Hz.
	*/
	CpuFrequencyHz *int64 `json:"cpuFrequencyHz,omitempty"`
	/*
	  CPU model name.
	*/
	CpuModel *string `json:"cpuModel,omitempty"`
	/*
	  Default VHD container UUID.
	*/
	DefaultVhdContainerUuid *string `json:"defaultVhdContainerUuid,omitempty"`
	/*
	  Default VHD location.
	*/
	DefaultVhdLocation *string `json:"defaultVhdLocation,omitempty"`
	/*
	  Default VM container UUID.
	*/
	DefaultVmContainerUuid *string `json:"defaultVmContainerUuid,omitempty"`
	/*
	  Default VM location.
	*/
	DefaultVmLocation *string `json:"defaultVmLocation,omitempty"`
	/*
	  Disks attached to host.
	*/
	Disk []DiskReference `json:"disk,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Failover cluster FQDN.
	*/
	FailoverClusterFqdn *string `json:"failoverClusterFqdn,omitempty"`
	/*
	  Failover cluster node status.
	*/
	FailoverClusterNodeStatus *string `json:"failoverClusterNodeStatus,omitempty"`
	/*
	  GPU driver version.
	*/
	GpuDriverVersion *string `json:"gpuDriverVersion,omitempty"`
	/*
	  GPU attached list.
	*/
	GpuList []string `json:"gpuList,omitempty"`
	/*
	  Certificate signing request status.
	*/
	HasCsr *bool `json:"hasCsr,omitempty"`
	/*
	  Name of the host.
	*/
	HostName *string `json:"hostName,omitempty"`

	HostType *HostTypeEnum `json:"hostType,omitempty"`

	Hypervisor *HypervisorReference `json:"hypervisor,omitempty"`
	/*
	  Node degraded status.
	*/
	IsDegraded *bool `json:"isDegraded,omitempty"`
	/*
	  Indicates whether the hardware is virtualized or not.
	*/
	IsHardwareVirtualized *bool `json:"isHardwareVirtualized,omitempty"`
	/*
	  Secure boot status.
	*/
	IsSecureBooted *bool `json:"isSecureBooted,omitempty"`
	/*
	  Mapping of key management device to certificate status list.
	*/
	KeyManagementDeviceToCertStatus []KeyManagementDeviceToCertStatusInfo `json:"keyManagementDeviceToCertStatus,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Host Maintenance State.
	*/
	MaintenanceState *string `json:"maintenanceState,omitempty"`
	/*
	  Memory size in bytes.
	*/
	MemorySizeBytes *int64 `json:"memorySizeBytes,omitempty"`

	NodeStatus *NodeStatus `json:"nodeStatus,omitempty"`
	/*
	  Number of CPU cores.
	*/
	NumberOfCpuCores *int64 `json:"numberOfCpuCores,omitempty"`
	/*
	  Number of CPU sockets.
	*/
	NumberOfCpuSockets *int64 `json:"numberOfCpuSockets,omitempty"`
	/*
	  Number of CPU threads.
	*/
	NumberOfCpuThreads *int64 `json:"numberOfCpuThreads,omitempty"`
	/*
	  Reboot pending status.
	*/
	RebootPending *bool `json:"rebootPending,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewHost() *Host {
	p := new(Host)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.Host"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Host GPU details.
*/
type HostGpu struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Cluster *ClusterReference `json:"cluster,omitempty"`

	Config *GpuConfig `json:"config,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Controller VM Id.
	*/
	NodeId *string `json:"nodeId,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  Number of vGPUs allocated.
	*/
	NumberOfVgpusAllocated *int64 `json:"numberOfVgpusAllocated,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewHostGpu() *HostGpu {
	p := new(HostGpu)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HostGpu"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Host rename parameters.
*/
type HostNameParam struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Name of the host.
	*/
	Name *string `json:"name"`
}

func (p *HostNameParam) MarshalJSON() ([]byte, error) {
	type HostNameParamProxy HostNameParam
	return json.Marshal(struct {
		*HostNameParamProxy
		Name *string `json:"name,omitempty"`
	}{
		HostNameParamProxy: (*HostNameParamProxy)(p),
		Name:               p.Name,
	})
}

func NewHostNameParam() *HostNameParam {
	p := new(HostNameParam)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HostNameParam"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Host NIC details.
*/
type HostNic struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of network switch interfaces attached to the host NIC.
	*/
	AttachedSwitchInterfaceList []NetworkSwitchInterface `json:"attachedSwitchInterfaceList,omitempty"`
	/*
	  Status of DHCP protocol.
	*/
	DhcpEnabled *bool `json:"dhcpEnabled,omitempty"`
	/*
	  Network discovery protocol (either LLDP or None).
	*/
	DiscoveryProtocol *string `json:"discoveryProtocol,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Host description.
	*/
	HostDescription *string `json:"hostDescription,omitempty"`
	/*
	  Operational status of the interface to the port associated with the NIC entity.
	*/
	InterfaceStatus *string `json:"interfaceStatus,omitempty"`
	/*
	  List of IPv4 addresses associated with the NIC entity for the network connection.
	*/
	Ipv4Addresses []import4.IPAddress `json:"ipv4Addresses,omitempty"`
	/*
	  List of IPv6 addresses associated with the NIC entity for the network connection.
	*/
	Ipv6Addresses []import4.IPAddress `json:"ipv6Addresses,omitempty"`
	/*
	  Link speed in Kbps.
	*/
	LinkSpeedInKbps *int64 `json:"linkSpeedInKbps,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Host Mac address.
	*/
	MacAddress *string `json:"macAddress,omitempty"`
	/*
	  Maximum transmission unit in bytes.
	*/
	MtuInBytes *int64 `json:"mtuInBytes,omitempty"`
	/*
	  Name of the host NIC.
	*/
	Name *string `json:"name,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  Size of configured buffer (in bytes) to the port associated with NIC, storing the network packets received through the port.
	*/
	RxRingSizeInBytes *int64 `json:"rxRingSizeInBytes,omitempty"`
	/*
	  Switch device Id learned through the discovery protocol.
	*/
	SwitchDeviceId *string `json:"switchDeviceId,omitempty"`
	/*
	  Switch Mac address
	*/
	SwitchMacAddress *string `json:"switchMacAddress,omitempty"`

	SwitchManagementIp *import4.IPAddress `json:"switchManagementIp,omitempty"`
	/*
	  Switch port Id learned through the discovery protocol.
	*/
	SwitchPortId *string `json:"switchPortId,omitempty"`
	/*
	  Switch vendor information learned through the discovery protocol.
	*/
	SwitchVendorInfo *string `json:"switchVendorInfo,omitempty"`
	/*
	  Switch VLAN Id learned through the discovery protocol.
	*/
	SwitchVlanId *string `json:"switchVlanId,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/*
	  Size of configured buffer (in bytes) to the port associated with NIC, storing the network packets that would be transmitted through the port.
	*/
	TxRingSizeInBytes *int64 `json:"txRingSizeInBytes,omitempty"`
}

func NewHostNic() *HostNic {
	p := new(HostNic)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HostNic"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Type of the host.
*/
type HostTypeEnum int

const (
	HOSTTYPEENUM_UNKNOWN         HostTypeEnum = 0
	HOSTTYPEENUM_REDACTED        HostTypeEnum = 1
	HOSTTYPEENUM_HYPER_CONVERGED HostTypeEnum = 2
	HOSTTYPEENUM_COMPUTE_ONLY    HostTypeEnum = 3
	HOSTTYPEENUM_STORAGE_ONLY    HostTypeEnum = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *HostTypeEnum) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"HYPER_CONVERGED",
		"COMPUTE_ONLY",
		"STORAGE_ONLY",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e HostTypeEnum) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"HYPER_CONVERGED",
		"COMPUTE_ONLY",
		"STORAGE_ONLY",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *HostTypeEnum) index(name string) HostTypeEnum {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"HYPER_CONVERGED",
		"COMPUTE_ONLY",
		"STORAGE_ONLY",
	}
	for idx := range names {
		if names[idx] == name {
			return HostTypeEnum(idx)
		}
	}
	return HOSTTYPEENUM_UNKNOWN
}

func (e *HostTypeEnum) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for HostTypeEnum:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *HostTypeEnum) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e HostTypeEnum) Ref() *HostTypeEnum {
	return &e
}

/*
HyperV Credentials.
*/
type HypervCredentials struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	DomainDetails *UserInfo `json:"domainDetails,omitempty"`

	FailoverClusterDetails *UserInfo `json:"failoverClusterDetails,omitempty"`
}

func NewHypervCredentials() *HypervCredentials {
	p := new(HypervCredentials)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HypervCredentials"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Map containing key as hypervisor type and value as md5sum of ISO.
*/
type HypervisorIsoMap struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Md5sum of ISO.
	*/
	Md5Sum *string `json:"md5Sum,omitempty"`

	Type *HypervisorType `json:"type,omitempty"`
}

func NewHypervisorIsoMap() *HypervisorIsoMap {
	p := new(HypervisorIsoMap)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HypervisorIsoMap"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Hypervisor details.
*/
type HypervisorReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	AcropolisConnectionState *AcropolisConnectionState `json:"acropolisConnectionState,omitempty"`

	ExternalAddress *import4.IPAddress `json:"externalAddress,omitempty"`
	/*
	  Hypervisor full name.
	*/
	FullName *string `json:"fullName,omitempty"`
	/*
	  Number of VMs.
	*/
	NumberOfVms *int64 `json:"numberOfVms,omitempty"`

	State *HypervisorState `json:"state,omitempty"`

	Type *HypervisorType `json:"type,omitempty"`
	/*
	  Hypervisor user name.
	*/
	UserName *string `json:"userName,omitempty"`
}

func NewHypervisorReference() *HypervisorReference {
	p := new(HypervisorReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HypervisorReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Hypervisor state.
*/
type HypervisorState int

const (
	HYPERVISORSTATE_UNKNOWN                                    HypervisorState = 0
	HYPERVISORSTATE_REDACTED                                   HypervisorState = 1
	HYPERVISORSTATE_ACROPOLIS_NORMAL                           HypervisorState = 2
	HYPERVISORSTATE_ENTERING_MAINTENANCE_MODE                  HypervisorState = 3
	HYPERVISORSTATE_ENTERED_MAINTENANCE_MODE                   HypervisorState = 4
	HYPERVISORSTATE_RESERVED_FOR_HA_FAILOVER                   HypervisorState = 5
	HYPERVISORSTATE_ENTERING_MAINTENANCE_MODE_FROM_HA_FAILOVER HypervisorState = 6
	HYPERVISORSTATE_RESERVING_FOR_HA_FAILOVER                  HypervisorState = 7
	HYPERVISORSTATE_HA_FAILOVER_SOURCE                         HypervisorState = 8
	HYPERVISORSTATE_HA_FAILOVER_TARGET                         HypervisorState = 9
	HYPERVISORSTATE_HA_HEALING_SOURCE                          HypervisorState = 10
	HYPERVISORSTATE_HA_HEALING_TARGET                          HypervisorState = 11
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *HypervisorState) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACROPOLIS_NORMAL",
		"ENTERING_MAINTENANCE_MODE",
		"ENTERED_MAINTENANCE_MODE",
		"RESERVED_FOR_HA_FAILOVER",
		"ENTERING_MAINTENANCE_MODE_FROM_HA_FAILOVER",
		"RESERVING_FOR_HA_FAILOVER",
		"HA_FAILOVER_SOURCE",
		"HA_FAILOVER_TARGET",
		"HA_HEALING_SOURCE",
		"HA_HEALING_TARGET",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e HypervisorState) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACROPOLIS_NORMAL",
		"ENTERING_MAINTENANCE_MODE",
		"ENTERED_MAINTENANCE_MODE",
		"RESERVED_FOR_HA_FAILOVER",
		"ENTERING_MAINTENANCE_MODE_FROM_HA_FAILOVER",
		"RESERVING_FOR_HA_FAILOVER",
		"HA_FAILOVER_SOURCE",
		"HA_FAILOVER_TARGET",
		"HA_HEALING_SOURCE",
		"HA_HEALING_TARGET",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *HypervisorState) index(name string) HypervisorState {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"ACROPOLIS_NORMAL",
		"ENTERING_MAINTENANCE_MODE",
		"ENTERED_MAINTENANCE_MODE",
		"RESERVED_FOR_HA_FAILOVER",
		"ENTERING_MAINTENANCE_MODE_FROM_HA_FAILOVER",
		"RESERVING_FOR_HA_FAILOVER",
		"HA_FAILOVER_SOURCE",
		"HA_FAILOVER_TARGET",
		"HA_HEALING_SOURCE",
		"HA_HEALING_TARGET",
	}
	for idx := range names {
		if names[idx] == name {
			return HypervisorState(idx)
		}
	}
	return HYPERVISORSTATE_UNKNOWN
}

func (e *HypervisorState) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for HypervisorState:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *HypervisorState) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e HypervisorState) Ref() *HypervisorState {
	return &e
}

/*
Hypervisor type.
*/
type HypervisorType int

const (
	HYPERVISORTYPE_UNKNOWN  HypervisorType = 0
	HYPERVISORTYPE_REDACTED HypervisorType = 1
	HYPERVISORTYPE_AHV      HypervisorType = 2
	HYPERVISORTYPE_ESX      HypervisorType = 3
	HYPERVISORTYPE_HYPERV   HypervisorType = 4
	HYPERVISORTYPE_XEN      HypervisorType = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *HypervisorType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"AHV",
		"ESX",
		"HYPERV",
		"XEN",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e HypervisorType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"AHV",
		"ESX",
		"HYPERV",
		"XEN",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *HypervisorType) index(name string) HypervisorType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"AHV",
		"ESX",
		"HYPERV",
		"XEN",
	}
	for idx := range names {
		if names[idx] == name {
			return HypervisorType(idx)
		}
	}
	return HYPERVISORTYPE_UNKNOWN
}

func (e *HypervisorType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for HypervisorType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *HypervisorType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e HypervisorType) Ref() *HypervisorType {
	return &e
}

/*
Hypervisor upload required information.
*/
type HypervisorUploadInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Error message.
	*/
	ErrorMessage *string `json:"errorMessage,omitempty"`
	/*
	  Node list containing upload information.
	*/
	UploadInfoNodeList []UploadInfoNodeItem `json:"uploadInfoNodeList,omitempty"`
}

func NewHypervisorUploadInfo() *HypervisorUploadInfo {
	p := new(HypervisorUploadInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HypervisorUploadInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Individual node item details for checking whether hypervisor ISO upload is required or not.
*/
type HypervisorUploadNodeListItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rackable unit Id in which node resides.
	*/
	BlockId *string `json:"blockId,omitempty"`

	HypervisorType *HypervisorType `json:"hypervisorType,omitempty"`
	/*
	  Host version of the node.
	*/
	HypervisorVersion *string `json:"hypervisorVersion,omitempty"`
	/*
	  Indicates whether the node is light compute or not.
	*/
	IsLightCompute *bool `json:"isLightCompute,omitempty"`
	/*
	  Indicates if node is minimum compute or not.
	*/
	MinimumComputeNode *bool `json:"minimumComputeNode,omitempty"`
	/*
	  Rackable unit model type.
	*/
	Model *string `json:"model,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  NOS software version of a node.
	*/
	NosVersion *string `json:"nosVersion,omitempty"`
	/*
	  Indicates whether the hypervisor is robo mixed or not.
	*/
	RoboMixedHypervisor *bool `json:"roboMixedHypervisor,omitempty"`
}

func NewHypervisorUploadNodeListItem() *HypervisorUploadNodeListItem {
	p := new(HypervisorUploadNodeListItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HypervisorUploadNodeListItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Parameters to get information on whether hypervisor ISO upload is required or not.
*/
type HypervisorUploadParam struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of node details for checking whether hypervisor ISO upload is required or not.
	*/
	NodeList []HypervisorUploadNodeListItem `json:"nodeList"`
}

func (p *HypervisorUploadParam) MarshalJSON() ([]byte, error) {
	type HypervisorUploadParamProxy HypervisorUploadParam
	return json.Marshal(struct {
		*HypervisorUploadParamProxy
		NodeList []HypervisorUploadNodeListItem `json:"nodeList,omitempty"`
	}{
		HypervisorUploadParamProxy: (*HypervisorUploadParamProxy)(p),
		NodeList:                   p.NodeList,
	})
}

func NewHypervisorUploadParam() *HypervisorUploadParam {
	p := new(HypervisorUploadParam)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.HypervisorUploadParam"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
IPMI reference.
*/
type IpmiReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Ip *import4.IPAddress `json:"ip,omitempty"`
	/*
	  IPMI username.
	*/
	Username *string `json:"username,omitempty"`
}

func NewIpmiReference() *IpmiReference {
	p := new(IpmiReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.IpmiReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Mapping of key management device to certificate status list.
*/
type KeyManagementDeviceToCertStatusInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Key management server name.
	*/
	KeyManagementServerName *string `json:"keyManagementServerName,omitempty"`
	/*
	  Certificate status.
	*/
	Status *bool `json:"status,omitempty"`
}

func NewKeyManagementDeviceToCertStatusInfo() *KeyManagementDeviceToCertStatusInfo {
	p := new(KeyManagementDeviceToCertStatusInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.KeyManagementDeviceToCertStatusInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Management server type.
*/
type KeyManagementServerType int

const (
	KEYMANAGEMENTSERVERTYPE_UNKNOWN       KeyManagementServerType = 0
	KEYMANAGEMENTSERVERTYPE_REDACTED      KeyManagementServerType = 1
	KEYMANAGEMENTSERVERTYPE_LOCAL         KeyManagementServerType = 2
	KEYMANAGEMENTSERVERTYPE_PRISM_CENTRAL KeyManagementServerType = 3
	KEYMANAGEMENTSERVERTYPE_EXTERNAL      KeyManagementServerType = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *KeyManagementServerType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"LOCAL",
		"PRISM_CENTRAL",
		"EXTERNAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e KeyManagementServerType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"LOCAL",
		"PRISM_CENTRAL",
		"EXTERNAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *KeyManagementServerType) index(name string) KeyManagementServerType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"LOCAL",
		"PRISM_CENTRAL",
		"EXTERNAL",
	}
	for idx := range names {
		if names[idx] == name {
			return KeyManagementServerType(idx)
		}
	}
	return KEYMANAGEMENTSERVERTYPE_UNKNOWN
}

func (e *KeyManagementServerType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for KeyManagementServerType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *KeyManagementServerType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e KeyManagementServerType) Ref() *KeyManagementServerType {
	return &e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters Get operation
*/
type ListClustersApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListClustersApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListClustersApiResponse() *ListClustersApiResponse {
	p := new(ListClustersApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListClustersApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListClustersApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListClustersApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListClustersApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/storage-containers/datastores Get operation
*/
type ListDataStoresByClusterIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListDataStoresByClusterIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListDataStoresByClusterIdApiResponse() *ListDataStoresByClusterIdApiResponse {
	p := new(ListDataStoresByClusterIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListDataStoresByClusterIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListDataStoresByClusterIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListDataStoresByClusterIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListDataStoresByClusterIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/disks Get operation
*/
type ListDisksApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListDisksApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListDisksApiResponse() *ListDisksApiResponse {
	p := new(ListDisksApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListDisksApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListDisksApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListDisksApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListDisksApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/hosts/{extId}/host-nics Get operation
*/
type ListHostNicsByHostIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListHostNicsByHostIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListHostNicsByHostIdApiResponse() *ListHostNicsByHostIdApiResponse {
	p := new(ListHostNicsByHostIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListHostNicsByHostIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListHostNicsByHostIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListHostNicsByHostIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListHostNicsByHostIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/hosts Get operation
*/
type ListHostsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListHostsApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListHostsApiResponse() *ListHostsApiResponse {
	p := new(ListHostsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListHostsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListHostsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListHostsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListHostsApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/hosts Get operation
*/
type ListHostsByClusterIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListHostsByClusterIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListHostsByClusterIdApiResponse() *ListHostsByClusterIdApiResponse {
	p := new(ListHostsByClusterIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListHostsByClusterIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListHostsByClusterIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListHostsByClusterIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListHostsByClusterIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/rackable-units Get operation
*/
type ListRackableUnitsByClusterIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListRackableUnitsByClusterIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListRackableUnitsByClusterIdApiResponse() *ListRackableUnitsByClusterIdApiResponse {
	p := new(ListRackableUnitsByClusterIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListRackableUnitsByClusterIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListRackableUnitsByClusterIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListRackableUnitsByClusterIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListRackableUnitsByClusterIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/rsyslog-servers Get operation
*/
type ListRsyslogServersByClusterIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListRsyslogServersByClusterIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListRsyslogServersByClusterIdApiResponse() *ListRsyslogServersByClusterIdApiResponse {
	p := new(ListRsyslogServersByClusterIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListRsyslogServersByClusterIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListRsyslogServersByClusterIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListRsyslogServersByClusterIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListRsyslogServersByClusterIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/storage-containers Get operation
*/
type ListStorageContainersApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListStorageContainersApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListStorageContainersApiResponse() *ListStorageContainersApiResponse {
	p := new(ListStorageContainersApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListStorageContainersApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListStorageContainersApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListStorageContainersApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListStorageContainersApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/vcenter-extensions Get operation
*/
type ListVcenterExtensionsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListVcenterExtensionsApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListVcenterExtensionsApiResponse() *ListVcenterExtensionsApiResponse {
	p := new(ListVcenterExtensionsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListVcenterExtensionsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListVcenterExtensionsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListVcenterExtensionsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListVcenterExtensionsApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/hosts/{extId}/virtual-nics Get operation
*/
type ListVirtualNicsByHostIdApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListVirtualNicsByHostIdApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListVirtualNicsByHostIdApiResponse() *ListVirtualNicsByHostIdApiResponse {
	p := new(ListVirtualNicsByHostIdApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ListVirtualNicsByHostIdApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListVirtualNicsByHostIdApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListVirtualNicsByHostIdApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListVirtualNicsByHostIdApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Management server information.
*/
type ManagementServerRef struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Indicates whether it is DRS enabled or not.
	*/
	DrsEnabled *bool `json:"drsEnabled,omitempty"`
	/*
	  Indicates whether the host is managed by an entity or not.
	*/
	InUse *bool `json:"inUse,omitempty"`

	Ip *import4.IPAddress `json:"ip,omitempty"`
	/*
	  Indicates whether it is registered or not.
	*/
	IsRegistered *bool `json:"isRegistered,omitempty"`

	Type *ManagementServerType `json:"type,omitempty"`
}

func NewManagementServerRef() *ManagementServerRef {
	p := new(ManagementServerRef)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ManagementServerRef"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Management server type.
*/
type ManagementServerType int

const (
	MANAGEMENTSERVERTYPE_UNKNOWN  ManagementServerType = 0
	MANAGEMENTSERVERTYPE_REDACTED ManagementServerType = 1
	MANAGEMENTSERVERTYPE_VCENTER  ManagementServerType = 2
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *ManagementServerType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"VCENTER",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e ManagementServerType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"VCENTER",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *ManagementServerType) index(name string) ManagementServerType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"VCENTER",
	}
	for idx := range names {
		if names[idx] == name {
			return ManagementServerType(idx)
		}
	}
	return MANAGEMENTSERVERTYPE_UNKNOWN
}

func (e *ManagementServerType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for ManagementServerType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *ManagementServerType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e ManagementServerType) Ref() *ManagementServerType {
	return &e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/storage-containers/{extId}/$actions/mount Post operation
*/
type MountStorageContainerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfMountStorageContainerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewMountStorageContainerApiResponse() *MountStorageContainerApiResponse {
	p := new(MountStorageContainerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.MountStorageContainerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *MountStorageContainerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *MountStorageContainerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfMountStorageContainerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Interface name and mac address.
*/
type NameMacRef struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Mac address.
	*/
	Mac *string `json:"mac,omitempty"`
	/*
	  Interface name.
	*/
	Name *string `json:"name,omitempty"`
}

func NewNameMacRef() *NameMacRef {
	p := new(NameMacRef)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NameMacRef"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Name and network information.
*/
type NameNetworkRef struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	HypervisorType *HypervisorType `json:"hypervisorType,omitempty"`
	/*
	  Interface name.
	*/
	Name *string `json:"name,omitempty"`
	/*
	  List of networks for interface.
	*/
	Networks []string `json:"networks,omitempty"`
}

func NewNameNetworkRef() *NameNetworkRef {
	p := new(NameNetworkRef)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NameNetworkRef"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Network information of HCI and SO nodes.
*/
type NetworkInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Network information of HCI nodes.
	*/
	Hci []NameNetworkRef `json:"hci,omitempty"`
	/*
	  Network information of SO nodes.
	*/
	So []NameNetworkRef `json:"so,omitempty"`
}

func NewNetworkInfo() *NetworkInfo {
	p := new(NetworkInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NetworkInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Network switch interface details.
*/
type NetworkSwitchInterface struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of host NIC UUIDs connected to this interface.
	*/
	AttachedHostNicUuids []string `json:"attachedHostNicUuids,omitempty"`
	/*
	  UUID of the host connected to the interface.
	*/
	AttachedHostUuid *string `json:"attachedHostUuid,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Network switch interface index.
	*/
	Index *int64 `json:"index,omitempty"`
	/*
	  Timestamp when the interface state was last changed or modified.
	*/
	LastChangeTime *time.Time `json:"lastChangeTime,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Host Mac address.
	*/
	MacAddress *string `json:"macAddress,omitempty"`
	/*
	  Maximum transmission unit in bytes.
	*/
	MtuInBytes *int64 `json:"mtuInBytes,omitempty"`
	/*
	  Network switch interface port number.
	*/
	Port *int64 `json:"port,omitempty"`
	/*
	  Network switch interface link speed in Kbps.
	*/
	SpeedInKbps *int64 `json:"speedInKbps,omitempty"`
	/*
	  Network switch interface description.
	*/
	SwitchInterfaceDescription *string `json:"switchInterfaceDescription,omitempty"`
	/*
	  Network switch interface name.
	*/
	SwitchInterfaceName *string `json:"switchInterfaceName,omitempty"`
	/*
	  Network switch interface type.
	*/
	SwitchInterfaceType *string `json:"switchInterfaceType,omitempty"`

	SwitchManagementAddress *import4.IPAddress `json:"switchManagementAddress,omitempty"`
	/*
	  UUID of the switch.
	*/
	SwitchUuid *string `json:"switchUuid,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewNetworkSwitchInterface() *NetworkSwitchInterface {
	p := new(NetworkSwitchInterface)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NetworkSwitchInterface"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Request type and networking details for nodes.
*/
type NodeDetails struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Node specific details required to fetch node networking information.
	*/
	NodeList []NodeListNetworkingDetails `json:"nodeList"`
	/*
	  Type of request, either it can be expand_cluster or npe.
	*/
	RequestType *string `json:"requestType,omitempty"`
}

func (p *NodeDetails) MarshalJSON() ([]byte, error) {
	type NodeDetailsProxy NodeDetails
	return json.Marshal(struct {
		*NodeDetailsProxy
		NodeList []NodeListNetworkingDetails `json:"nodeList,omitempty"`
	}{
		NodeDetailsProxy: (*NodeDetailsProxy)(p),
		NodeList:         p.NodeList,
	})
}

func NewNodeDetails() *NodeDetails {
	p := new(NodeDetails)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeDetails"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Discover unconfigured node details.
*/
type NodeDiscoveryParams struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	AddressType *AddressType `json:"addressType,omitempty"`
	/*
	  Interface name that is used for packet broadcasting.
	*/
	InterfaceFilterList []string `json:"interfaceFilterList,omitempty"`
	/*
	  IP addresses of the unconfigured nodes.
	*/
	IpFilterList []import4.IPAddress `json:"ipFilterList,omitempty"`
	/*
	  Indicates if the discovery is manual or not.
	*/
	IsManualDiscovery *bool `json:"isManualDiscovery,omitempty"`
	/*
	  Timeout for discovering nodes in seconds.
	*/
	Timeout *int64 `json:"timeout,omitempty"`
	/*
	  Unconfigured node UUIDs.
	*/
	UuidFilterList []string `json:"uuidFilterList,omitempty"`
}

func NewNodeDiscoveryParams() *NodeDiscoveryParams {
	p := new(NodeDiscoveryParams)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeDiscoveryParams"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Node item containing attributes of node.
*/
type NodeInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rackable unit serial name.
	*/
	BlockId *string `json:"blockId,omitempty"`
	/*
	  Current network interface of a node.
	*/
	CurrentNetworkInterface *string `json:"currentNetworkInterface,omitempty"`

	CvmIp *import4.IPAddress `json:"cvmIp,omitempty"`
	/*
	  List of objects containing digital_certificate_base64 and key_management_server_uuid fields for key management server.
	*/
	DigitalCertificateMapList []DigitalCertificateMapReference `json:"digitalCertificateMapList,omitempty"`
	/*
	  Name of the host.
	*/
	HypervisorHostname *string `json:"hypervisorHostname,omitempty"`

	HypervisorIp *import4.IPAddress `json:"hypervisorIp,omitempty"`

	HypervisorType *HypervisorType `json:"hypervisorType,omitempty"`
	/*
	  Host version of the node.
	*/
	HypervisorVersion *string `json:"hypervisorVersion,omitempty"`

	IpmiIp *import4.IPAddress `json:"ipmiIp,omitempty"`
	/*
	  Indicates whether the node is light compute or not.
	*/
	IsLightCompute *bool `json:"isLightCompute,omitempty"`
	/*
	  Rackable unit model name.
	*/
	Model *string `json:"model,omitempty"`
	/*
	  Position of a node in a rackable unit.
	*/
	NodePosition *string `json:"nodePosition,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  NOS software version of a node.
	*/
	NosVersion *string `json:"nosVersion,omitempty"`
	/*
	  Indicates whether the hypervisor is robo mixed or not.
	*/
	RoboMixedHypervisor *bool `json:"roboMixedHypervisor,omitempty"`
}

func NewNodeInfo() *NodeInfo {
	p := new(NodeInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Node item containing attributes of node.
*/
type NodeItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rackable unit serial name.
	*/
	BlockId *string `json:"blockId,omitempty"`
	/*
	  Current network interface of a node.
	*/
	CurrentNetworkInterface *string `json:"currentNetworkInterface,omitempty"`

	CvmIp *import4.IPAddress `json:"cvmIp,omitempty"`
	/*
	  List of objects containing digital_certificate_base64 and key_management_server_uuid fields for key management server.
	*/
	DigitalCertificateMapList []DigitalCertificateMapReference `json:"digitalCertificateMapList,omitempty"`
	/*
	  Name of the host.
	*/
	HypervisorHostname *string `json:"hypervisorHostname,omitempty"`

	HypervisorIp *import4.IPAddress `json:"hypervisorIp,omitempty"`

	HypervisorType *HypervisorType `json:"hypervisorType,omitempty"`
	/*
	  Host version of the node.
	*/
	HypervisorVersion *string `json:"hypervisorVersion,omitempty"`

	IpmiIp *import4.IPAddress `json:"ipmiIp,omitempty"`
	/*
	  Indicates whether the node is light compute or not.
	*/
	IsLightCompute *bool `json:"isLightCompute,omitempty"`
	/*
	  Rackable unit model name.
	*/
	Model *string `json:"model,omitempty"`
	/*
	  Active and standby uplink information of the target nodes.
	*/
	Networks []UplinkNetworkItem `json:"networks,omitempty"`
	/*
	  Position of a node in a rackable unit.
	*/
	NodePosition *string `json:"nodePosition,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  NOS software version of a node.
	*/
	NosVersion *string `json:"nosVersion,omitempty"`
	/*
	  Indicates whether the hypervisor is robo mixed or not.
	*/
	RoboMixedHypervisor *bool `json:"roboMixedHypervisor,omitempty"`
}

func NewNodeItem() *NodeItem {
	p := new(NodeItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
List of nodes in a cluster.
*/
type NodeListItemReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	ControllerVmIp *import4.IPAddress `json:"controllerVmIp,omitempty"`

	HostIp *import4.IPAddress `json:"hostIp,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
}

func NewNodeListItemReference() *NodeListItemReference {
	p := new(NodeListItemReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeListItemReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Node specific details required to fetch node networking information.
*/
type NodeListNetworkingDetails struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rackable unit Id in which node resides.
	*/
	BlockId *string `json:"blockId,omitempty"`
	/*
	  Current network interface of a node.
	*/
	CurrentNetworkInterface *string `json:"currentNetworkInterface,omitempty"`

	CvmIp *import4.IPAddress `json:"cvmIp,omitempty"`
	/*
	  List of objects containing digital_certificate_base64 and key_management_server_uuid fields for key management server.
	*/
	DigitalCertificateMapList []DigitalCertificateMapReference `json:"digitalCertificateMapList,omitempty"`

	HypervisorIp *import4.IPAddress `json:"hypervisorIp,omitempty"`

	HypervisorType *HypervisorType `json:"hypervisorType,omitempty"`
	/*
	  Host version of the node.
	*/
	HypervisorVersion *string `json:"hypervisorVersion,omitempty"`

	IpmiIp *import4.IPAddress `json:"ipmiIp,omitempty"`
	/*
	  Indicates whether the node is compute only or not.
	*/
	IsComputeOnly *bool `json:"isComputeOnly,omitempty"`
	/*
	  Indicates whether the node is light compute or not.
	*/
	IsLightCompute *bool `json:"isLightCompute,omitempty"`
	/*
	  Rackable unit model name.
	*/
	Model *string `json:"model,omitempty"`
	/*
	  Position of a node in a rackable unit.
	*/
	NodePosition *string `json:"nodePosition,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  NOS software version of a node.
	*/
	NosVersion *string `json:"nosVersion,omitempty"`
	/*
	  Indicates whether the hypervisor is robo mixed or not.
	*/
	RoboMixedHypervisor *bool `json:"roboMixedHypervisor,omitempty"`
}

func NewNodeListNetworkingDetails() *NodeListNetworkingDetails {
	p := new(NodeListNetworkingDetails)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeListNetworkingDetails"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Network details of nodes.
*/
type NodeNetworkingDetails struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	NetworkInfo *NetworkInfo `json:"networkInfo,omitempty"`
	/*
	  List of uplinks information for each CVM IP.
	*/
	Uplinks []UplinkInfo `json:"uplinks,omitempty"`
	/*
	  List of warning messages.
	*/
	Warnings []string `json:"warnings,omitempty"`
}

func NewNodeNetworkingDetails() *NodeNetworkingDetails {
	p := new(NodeNetworkingDetails)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeNetworkingDetails"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Parameters of the node to be added.
*/
type NodeParam struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Block list of a cluster.
	*/
	BlockList []BlockItem `json:"blockList,omitempty"`

	BundleInfo *BundleInfo `json:"bundleInfo,omitempty"`
	/*
	  List of compute only nodes.
	*/
	ComputeNodeList []ComputeNodeItem `json:"computeNodeList,omitempty"`
	/*
	  Hyperv SKU.
	*/
	HypervSku *string `json:"hypervSku,omitempty"`
	/*
	  Hypervisor type to md5sum map.
	*/
	HypervisorIsos []HypervisorIsoMap `json:"hypervisorIsos,omitempty"`
	/*
	  List of nodes in a cluster.
	*/
	NodeList []NodeItem `json:"nodeList,omitempty"`
	/*
	  Indicates if the host networking needs to be skipped or not.
	*/
	SkipHostNetworking *bool `json:"skipHostNetworking,omitempty"`
}

func NewNodeParam() *NodeParam {
	p := new(NodeParam)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeParam"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Node reference for a cluster.
*/
type NodeReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of nodes in a cluster.
	*/
	NodeList []NodeListItemReference `json:"nodeList,omitempty"`
	/*
	  Number of nodes in a cluster.
	*/
	NumberOfNodes *int `json:"numberOfNodes,omitempty"`
}

func NewNodeReference() *NodeReference {
	p := new(NodeReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Extra parameters for node addition.
*/
type NodeRemovalExtraParam struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Indicates if add node check need to be skip or not.
	*/
	SkipAddCheck *bool `json:"skipAddCheck,omitempty"`
	/*
	  Indicates if space check needs to be skip or not.
	*/
	SkipSpaceCheck *bool `json:"skipSpaceCheck,omitempty"`
	/*
	  Indicates if upgrade check needs to be skip or not.
	*/
	SkipUpgradeCheck *bool `json:"skipUpgradeCheck,omitempty"`
}

func NewNodeRemovalExtraParam() *NodeRemovalExtraParam {
	p := new(NodeRemovalExtraParam)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeRemovalExtraParam"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Parameters to remove nodes from cluster.
*/
type NodeRemovalParams struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	ExtraParams *NodeRemovalExtraParam `json:"extraParams,omitempty"`
	/*
	  List of node UUIDs to remove.
	*/
	NodeUuids []string `json:"nodeUuids"`
	/*
	  Indicates if prechecks can be skipped for node removal.
	*/
	SkipPrechecks *bool `json:"skipPrechecks,omitempty"`
	/*
	  Indicates if node removal can be skipped.
	*/
	SkipRemove *bool `json:"skipRemove,omitempty"`
}

func (p *NodeRemovalParams) MarshalJSON() ([]byte, error) {
	type NodeRemovalParamsProxy NodeRemovalParams
	return json.Marshal(struct {
		*NodeRemovalParamsProxy
		NodeUuids []string `json:"nodeUuids,omitempty"`
	}{
		NodeRemovalParamsProxy: (*NodeRemovalParamsProxy)(p),
		NodeUuids:              p.NodeUuids,
	})
}

func NewNodeRemovalParams() *NodeRemovalParams {
	p := new(NodeRemovalParams)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.NodeRemovalParams"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Node status.
*/
type NodeStatus int

const (
	NODESTATUS_UNKNOWN            NodeStatus = 0
	NODESTATUS_REDACTED           NodeStatus = 1
	NODESTATUS_NORMAL             NodeStatus = 2
	NODESTATUS_TO_BE_REMOVED      NodeStatus = 3
	NODESTATUS_OK_TO_BE_REMOVED   NodeStatus = 4
	NODESTATUS_NEW_NODE           NodeStatus = 5
	NODESTATUS_TO_BE_PREPROTECTED NodeStatus = 6
	NODESTATUS_PREPROTECTED       NodeStatus = 7
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *NodeStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"TO_BE_REMOVED",
		"OK_TO_BE_REMOVED",
		"NEW_NODE",
		"TO_BE_PREPROTECTED",
		"PREPROTECTED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e NodeStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"TO_BE_REMOVED",
		"OK_TO_BE_REMOVED",
		"NEW_NODE",
		"TO_BE_PREPROTECTED",
		"PREPROTECTED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *NodeStatus) index(name string) NodeStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"TO_BE_REMOVED",
		"OK_TO_BE_REMOVED",
		"NEW_NODE",
		"TO_BE_PREPROTECTED",
		"PREPROTECTED",
	}
	for idx := range names {
		if names[idx] == name {
			return NodeStatus(idx)
		}
	}
	return NODESTATUS_UNKNOWN
}

func (e *NodeStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for NodeStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *NodeStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e NodeStatus) Ref() *NodeStatus {
	return &e
}

/*
Indicates the current status of Disk Deduplication for the Container.
*/
type OnDiskDedup int

const (
	ONDISKDEDUP_UNKNOWN      OnDiskDedup = 0
	ONDISKDEDUP_REDACTED     OnDiskDedup = 1
	ONDISKDEDUP_NONE         OnDiskDedup = 2
	ONDISKDEDUP_OFF          OnDiskDedup = 3
	ONDISKDEDUP_POST_PROCESS OnDiskDedup = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *OnDiskDedup) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"POST_PROCESS",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e OnDiskDedup) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"POST_PROCESS",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *OnDiskDedup) index(name string) OnDiskDedup {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"OFF",
		"POST_PROCESS",
	}
	for idx := range names {
		if names[idx] == name {
			return OnDiskDedup(idx)
		}
	}
	return ONDISKDEDUP_UNKNOWN
}

func (e *OnDiskDedup) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for OnDiskDedup:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *OnDiskDedup) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e OnDiskDedup) Ref() *OnDiskDedup {
	return &e
}

/*
Cluster operation mode. This is part of payload for cluster update operation only.
*/
type OperationMode int

const (
	OPERATIONMODE_UNKNOWN            OperationMode = 0
	OPERATIONMODE_REDACTED           OperationMode = 1
	OPERATIONMODE_NORMAL             OperationMode = 2
	OPERATIONMODE_READ_ONLY          OperationMode = 3
	OPERATIONMODE_STAND_ALONE        OperationMode = 4
	OPERATIONMODE_SWITCH_TO_TWO_NODE OperationMode = 5
	OPERATIONMODE_OVERRIDE           OperationMode = 6
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *OperationMode) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"READ_ONLY",
		"STAND_ALONE",
		"SWITCH_TO_TWO_NODE",
		"OVERRIDE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e OperationMode) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"READ_ONLY",
		"STAND_ALONE",
		"SWITCH_TO_TWO_NODE",
		"OVERRIDE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *OperationMode) index(name string) OperationMode {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NORMAL",
		"READ_ONLY",
		"STAND_ALONE",
		"SWITCH_TO_TWO_NODE",
		"OVERRIDE",
	}
	for idx := range names {
		if names[idx] == name {
			return OperationMode(idx)
		}
	}
	return OPERATIONMODE_UNKNOWN
}

func (e *OperationMode) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for OperationMode:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *OperationMode) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e OperationMode) Ref() *OperationMode {
	return &e
}

/*
Public ssh key details. This is part of payload for cluster update operation only.
*/
type PublicKey struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Ssh key value.
	*/
	Key *string `json:"key"`
	/*
	  Ssh key name.
	*/
	Name *string `json:"name"`
}

func (p *PublicKey) MarshalJSON() ([]byte, error) {
	type PublicKeyProxy PublicKey
	return json.Marshal(struct {
		*PublicKeyProxy
		Key  *string `json:"key,omitempty"`
		Name *string `json:"name,omitempty"`
	}{
		PublicKeyProxy: (*PublicKeyProxy)(p),
		Key:            p.Key,
		Name:           p.Name,
	})
}

func NewPublicKey() *PublicKey {
	p := new(PublicKey)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.PublicKey"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Rack reference for the block.
*/
type RackReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Rack Id.
	*/
	Id *int64 `json:"id,omitempty"`
	/*
	  Rack UUID.
	*/
	Uuid *string `json:"uuid,omitempty"`
}

func NewRackReference() *RackReference {
	p := new(RackReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RackReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Rackable Unit configuration.
*/
type RackableUnit struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Rackable unit Id.
	*/
	Id *int64 `json:"id,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`

	Model *RackableUnitModel `json:"model,omitempty"`
	/*
	  Rackable unit model name.
	*/
	ModelName *string `json:"modelName,omitempty"`
	/*
	  List of node information registered to the block.
	*/
	Nodes []RackableUnitNode `json:"nodes,omitempty"`

	Rack *RackReference `json:"rack,omitempty"`
	/*
	  Rackable unit serial name.
	*/
	Serial *string `json:"serial,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewRackableUnit() *RackableUnit {
	p := new(RackableUnit)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RackableUnit"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Rackable unit model type.
*/
type RackableUnitModel int

const (
	RACKABLEUNITMODEL_UNKNOWN   RackableUnitModel = 0
	RACKABLEUNITMODEL_REDACTED  RackableUnitModel = 1
	RACKABLEUNITMODEL_DESKTOP   RackableUnitModel = 2
	RACKABLEUNITMODEL_NX2000    RackableUnitModel = 3
	RACKABLEUNITMODEL_NX3000    RackableUnitModel = 4
	RACKABLEUNITMODEL_NX3050    RackableUnitModel = 5
	RACKABLEUNITMODEL_NX6050    RackableUnitModel = 6
	RACKABLEUNITMODEL_NX6070    RackableUnitModel = 7
	RACKABLEUNITMODEL_NX1050    RackableUnitModel = 8
	RACKABLEUNITMODEL_NX3060    RackableUnitModel = 9
	RACKABLEUNITMODEL_NX6060    RackableUnitModel = 10
	RACKABLEUNITMODEL_NX6080    RackableUnitModel = 11
	RACKABLEUNITMODEL_NX6020    RackableUnitModel = 12
	RACKABLEUNITMODEL_NX7110    RackableUnitModel = 13
	RACKABLEUNITMODEL_NX1020    RackableUnitModel = 14
	RACKABLEUNITMODEL_NX9040    RackableUnitModel = 15
	RACKABLEUNITMODEL_USELAYOUT RackableUnitModel = 16
	RACKABLEUNITMODEL_NULLVALUE RackableUnitModel = 17
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *RackableUnitModel) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"DESKTOP",
		"NX2000",
		"NX3000",
		"NX3050",
		"NX6050",
		"NX6070",
		"NX1050",
		"NX3060",
		"NX6060",
		"NX6080",
		"NX6020",
		"NX7110",
		"NX1020",
		"NX9040",
		"USELAYOUT",
		"NULLVALUE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e RackableUnitModel) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"DESKTOP",
		"NX2000",
		"NX3000",
		"NX3050",
		"NX6050",
		"NX6070",
		"NX1050",
		"NX3060",
		"NX6060",
		"NX6080",
		"NX6020",
		"NX7110",
		"NX1020",
		"NX9040",
		"USELAYOUT",
		"NULLVALUE",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *RackableUnitModel) index(name string) RackableUnitModel {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"DESKTOP",
		"NX2000",
		"NX3000",
		"NX3050",
		"NX6050",
		"NX6070",
		"NX1050",
		"NX3060",
		"NX6060",
		"NX6080",
		"NX6020",
		"NX7110",
		"NX1020",
		"NX9040",
		"USELAYOUT",
		"NULLVALUE",
	}
	for idx := range names {
		if names[idx] == name {
			return RackableUnitModel(idx)
		}
	}
	return RACKABLEUNITMODEL_UNKNOWN
}

func (e *RackableUnitModel) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for RackableUnitModel:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *RackableUnitModel) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e RackableUnitModel) Ref() *RackableUnitModel {
	return &e
}

/*
Node information registered to this rackable unit.
*/
type RackableUnitNode struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Position of a node in a rackable unit.
	*/
	Position *int `json:"position,omitempty"`
	/*
	  Controller VM Id.
	*/
	SvmId *int64 `json:"svmId,omitempty"`
	/*
	  UUID of the host.
	*/
	Uuid *string `json:"uuid,omitempty"`
}

func NewRackableUnitNode() *RackableUnitNode {
	p := new(RackableUnitNode)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RackableUnitNode"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/vcenter-extensions/{extId}/$actions/register Post operation
*/
type RegisterVcenterExtensionApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfRegisterVcenterExtensionApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewRegisterVcenterExtensionApiResponse() *RegisterVcenterExtensionApiResponse {
	p := new(RegisterVcenterExtensionApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RegisterVcenterExtensionApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *RegisterVcenterExtensionApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *RegisterVcenterExtensionApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfRegisterVcenterExtensionApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/$actions/remove-node Post operation
*/
type RemoveNodeApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfRemoveNodeApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewRemoveNodeApiResponse() *RemoveNodeApiResponse {
	p := new(RemoveNodeApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RemoveNodeApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *RemoveNodeApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *RemoveNodeApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfRemoveNodeApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/snmp/$actions/remove-transports Post operation
*/
type RemoveSnmpTransportsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfRemoveSnmpTransportsApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewRemoveSnmpTransportsApiResponse() *RemoveSnmpTransportsApiResponse {
	p := new(RemoveSnmpTransportsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RemoveSnmpTransportsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *RemoveSnmpTransportsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *RemoveSnmpTransportsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfRemoveSnmpTransportsApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/hosts/{extId}/$actions/rename-host Post operation
*/
type RenameHostApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfRenameHostApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewRenameHostApiResponse() *RenameHostApiResponse {
	p := new(RenameHostApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RenameHostApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *RenameHostApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *RenameHostApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfRenameHostApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
RSYSLOG Module information.
*/
type RsyslogModuleItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	LogSeverityLevel *RsyslogModuleLogSeverityLevel `json:"logSeverityLevel"`
	/*
	  Option to log, monitor/output files of a module.
	*/
	Monitor *bool `json:"monitor"`

	Name *RsyslogModuleName `json:"name"`
}

func (p *RsyslogModuleItem) MarshalJSON() ([]byte, error) {
	type RsyslogModuleItemProxy RsyslogModuleItem
	return json.Marshal(struct {
		*RsyslogModuleItemProxy
		LogSeverityLevel *RsyslogModuleLogSeverityLevel `json:"logSeverityLevel,omitempty"`
		Monitor          *bool                          `json:"monitor,omitempty"`
		Name             *RsyslogModuleName             `json:"name,omitempty"`
	}{
		RsyslogModuleItemProxy: (*RsyslogModuleItemProxy)(p),
		LogSeverityLevel:       p.LogSeverityLevel,
		Monitor:                p.Monitor,
		Name:                   p.Name,
	})
}

func NewRsyslogModuleItem() *RsyslogModuleItem {
	p := new(RsyslogModuleItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RsyslogModuleItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
RSYSLOG module log severity level.
*/
type RsyslogModuleLogSeverityLevel int

const (
	RSYSLOGMODULELOGSEVERITYLEVEL_UNKNOWN   RsyslogModuleLogSeverityLevel = 0
	RSYSLOGMODULELOGSEVERITYLEVEL_REDACTED  RsyslogModuleLogSeverityLevel = 1
	RSYSLOGMODULELOGSEVERITYLEVEL_EMERGENCY RsyslogModuleLogSeverityLevel = 2
	RSYSLOGMODULELOGSEVERITYLEVEL_ALERT     RsyslogModuleLogSeverityLevel = 3
	RSYSLOGMODULELOGSEVERITYLEVEL_CRITICAL  RsyslogModuleLogSeverityLevel = 4
	RSYSLOGMODULELOGSEVERITYLEVEL_ERROR     RsyslogModuleLogSeverityLevel = 5
	RSYSLOGMODULELOGSEVERITYLEVEL_WARNING   RsyslogModuleLogSeverityLevel = 6
	RSYSLOGMODULELOGSEVERITYLEVEL_NOTICE    RsyslogModuleLogSeverityLevel = 7
	RSYSLOGMODULELOGSEVERITYLEVEL_INFO      RsyslogModuleLogSeverityLevel = 8
	RSYSLOGMODULELOGSEVERITYLEVEL_DEBUG     RsyslogModuleLogSeverityLevel = 9
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *RsyslogModuleLogSeverityLevel) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"EMERGENCY",
		"ALERT",
		"CRITICAL",
		"ERROR",
		"WARNING",
		"NOTICE",
		"INFO",
		"DEBUG",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e RsyslogModuleLogSeverityLevel) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"EMERGENCY",
		"ALERT",
		"CRITICAL",
		"ERROR",
		"WARNING",
		"NOTICE",
		"INFO",
		"DEBUG",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *RsyslogModuleLogSeverityLevel) index(name string) RsyslogModuleLogSeverityLevel {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"EMERGENCY",
		"ALERT",
		"CRITICAL",
		"ERROR",
		"WARNING",
		"NOTICE",
		"INFO",
		"DEBUG",
	}
	for idx := range names {
		if names[idx] == name {
			return RsyslogModuleLogSeverityLevel(idx)
		}
	}
	return RSYSLOGMODULELOGSEVERITYLEVEL_UNKNOWN
}

func (e *RsyslogModuleLogSeverityLevel) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for RsyslogModuleLogSeverityLevel:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *RsyslogModuleLogSeverityLevel) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e RsyslogModuleLogSeverityLevel) Ref() *RsyslogModuleLogSeverityLevel {
	return &e
}

/*
RSYSLOG module name.
*/
type RsyslogModuleName int

const (
	RSYSLOGMODULENAME_UNKNOWN           RsyslogModuleName = 0
	RSYSLOGMODULENAME_REDACTED          RsyslogModuleName = 1
	RSYSLOGMODULENAME_CASSANDRA         RsyslogModuleName = 2
	RSYSLOGMODULENAME_CEREBRO           RsyslogModuleName = 3
	RSYSLOGMODULENAME_CURATOR           RsyslogModuleName = 4
	RSYSLOGMODULENAME_GENESIS           RsyslogModuleName = 5
	RSYSLOGMODULENAME_PRISM             RsyslogModuleName = 6
	RSYSLOGMODULENAME_STARGATE          RsyslogModuleName = 7
	RSYSLOGMODULENAME_SYSLOG_MODULE     RsyslogModuleName = 8
	RSYSLOGMODULENAME_ZOOKEEPER         RsyslogModuleName = 9
	RSYSLOGMODULENAME_UHARA             RsyslogModuleName = 10
	RSYSLOGMODULENAME_LAZAN             RsyslogModuleName = 11
	RSYSLOGMODULENAME_API_AUDIT         RsyslogModuleName = 12
	RSYSLOGMODULENAME_AUDIT             RsyslogModuleName = 13
	RSYSLOGMODULENAME_CALM              RsyslogModuleName = 14
	RSYSLOGMODULENAME_EPSILON           RsyslogModuleName = 15
	RSYSLOGMODULENAME_ACROPOLIS         RsyslogModuleName = 16
	RSYSLOGMODULENAME_MINERVA_CVM       RsyslogModuleName = 17
	RSYSLOGMODULENAME_FLOW              RsyslogModuleName = 18
	RSYSLOGMODULENAME_FLOW_SERVICE_LOGS RsyslogModuleName = 19
	RSYSLOGMODULENAME_LCM               RsyslogModuleName = 20
	RSYSLOGMODULENAME_APLOS             RsyslogModuleName = 21
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *RsyslogModuleName) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CASSANDRA",
		"CEREBRO",
		"CURATOR",
		"GENESIS",
		"PRISM",
		"STARGATE",
		"SYSLOG_MODULE",
		"ZOOKEEPER",
		"UHARA",
		"LAZAN",
		"API_AUDIT",
		"AUDIT",
		"CALM",
		"EPSILON",
		"ACROPOLIS",
		"MINERVA_CVM",
		"FLOW",
		"FLOW_SERVICE_LOGS",
		"LCM",
		"APLOS",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e RsyslogModuleName) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CASSANDRA",
		"CEREBRO",
		"CURATOR",
		"GENESIS",
		"PRISM",
		"STARGATE",
		"SYSLOG_MODULE",
		"ZOOKEEPER",
		"UHARA",
		"LAZAN",
		"API_AUDIT",
		"AUDIT",
		"CALM",
		"EPSILON",
		"ACROPOLIS",
		"MINERVA_CVM",
		"FLOW",
		"FLOW_SERVICE_LOGS",
		"LCM",
		"APLOS",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *RsyslogModuleName) index(name string) RsyslogModuleName {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"CASSANDRA",
		"CEREBRO",
		"CURATOR",
		"GENESIS",
		"PRISM",
		"STARGATE",
		"SYSLOG_MODULE",
		"ZOOKEEPER",
		"UHARA",
		"LAZAN",
		"API_AUDIT",
		"AUDIT",
		"CALM",
		"EPSILON",
		"ACROPOLIS",
		"MINERVA_CVM",
		"FLOW",
		"FLOW_SERVICE_LOGS",
		"LCM",
		"APLOS",
	}
	for idx := range names {
		if names[idx] == name {
			return RsyslogModuleName(idx)
		}
	}
	return RSYSLOGMODULENAME_UNKNOWN
}

func (e *RsyslogModuleName) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for RsyslogModuleName:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *RsyslogModuleName) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e RsyslogModuleName) Ref() *RsyslogModuleName {
	return &e
}

/*
RSYSLOG server protocol type.
*/
type RsyslogNetworkProtocol int

const (
	RSYSLOGNETWORKPROTOCOL_UNKNOWN  RsyslogNetworkProtocol = 0
	RSYSLOGNETWORKPROTOCOL_REDACTED RsyslogNetworkProtocol = 1
	RSYSLOGNETWORKPROTOCOL_UDP      RsyslogNetworkProtocol = 2
	RSYSLOGNETWORKPROTOCOL_TCP      RsyslogNetworkProtocol = 3
	RSYSLOGNETWORKPROTOCOL_RELP     RsyslogNetworkProtocol = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *RsyslogNetworkProtocol) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UDP",
		"TCP",
		"RELP",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e RsyslogNetworkProtocol) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UDP",
		"TCP",
		"RELP",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *RsyslogNetworkProtocol) index(name string) RsyslogNetworkProtocol {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UDP",
		"TCP",
		"RELP",
	}
	for idx := range names {
		if names[idx] == name {
			return RsyslogNetworkProtocol(idx)
		}
	}
	return RSYSLOGNETWORKPROTOCOL_UNKNOWN
}

func (e *RsyslogNetworkProtocol) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for RsyslogNetworkProtocol:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *RsyslogNetworkProtocol) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e RsyslogNetworkProtocol) Ref() *RsyslogNetworkProtocol {
	return &e
}

type RsyslogServer struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`

	IpAddress *import4.IPAddress `json:"ipAddress"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  List of modules registered to RSYSLOG server.
	*/
	Modules []RsyslogModuleItem `json:"modules,omitempty"`

	NetworkProtocol *RsyslogNetworkProtocol `json:"networkProtocol"`
	/*
	  RSYSLOG server port.
	*/
	Port *int `json:"port"`
	/*
	  RSYSLOG server name.
	*/
	ServerName *string `json:"serverName"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func (p *RsyslogServer) MarshalJSON() ([]byte, error) {
	type RsyslogServerProxy RsyslogServer
	return json.Marshal(struct {
		*RsyslogServerProxy
		IpAddress       *import4.IPAddress      `json:"ipAddress,omitempty"`
		NetworkProtocol *RsyslogNetworkProtocol `json:"networkProtocol,omitempty"`
		Port            *int                    `json:"port,omitempty"`
		ServerName      *string                 `json:"serverName,omitempty"`
	}{
		RsyslogServerProxy: (*RsyslogServerProxy)(p),
		IpAddress:          p.IpAddress,
		NetworkProtocol:    p.NetworkProtocol,
		Port:               p.Port,
		ServerName:         p.ServerName,
	})
}

func NewRsyslogServer() *RsyslogServer {
	p := new(RsyslogServer)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.RsyslogServer"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Search parameters.
*/
type SearchParams struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	SearchType *SearchType `json:"searchType"`
}

func (p *SearchParams) MarshalJSON() ([]byte, error) {
	type SearchParamsProxy SearchParams
	return json.Marshal(struct {
		*SearchParamsProxy
		SearchType *SearchType `json:"searchType,omitempty"`
	}{
		SearchParamsProxy: (*SearchParamsProxy)(p),
		SearchType:        p.SearchType,
	})
}

func NewSearchParams() *SearchParams {
	p := new(SearchParams)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SearchParams"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Task Response which is one of node-discovery, networking-details, hypervisor-upload information.
*/
type SearchResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	ResponseItemDiscriminator_ *string `json:"$responseItemDiscriminator,omitempty"`

	Response *OneOfSearchResponseResponse `json:"response,omitempty"`

	SearchType *SearchType `json:"searchType,omitempty"`
}

func NewSearchResponse() *SearchResponse {
	p := new(SearchResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SearchResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *SearchResponse) GetResponse() interface{} {
	if nil == p.Response {
		return nil
	}
	return p.Response.GetValue()
}

func (p *SearchResponse) SetResponse(v interface{}) error {
	if nil == p.Response {
		p.Response = NewOneOfSearchResponseResponse()
	}
	e := p.Response.SetValue(v)
	if nil == e {
		if nil == p.ResponseItemDiscriminator_ {
			p.ResponseItemDiscriminator_ = new(string)
		}
		*p.ResponseItemDiscriminator_ = *p.Response.Discriminator
	}
	return e
}

/*
Task Response search type.
*/
type SearchType int

const (
	SEARCHTYPE_UNKNOWN                SearchType = 0
	SEARCHTYPE_REDACTED               SearchType = 1
	SEARCHTYPE_UNCONFIGURED_NODES     SearchType = 2
	SEARCHTYPE_NETWORKING_DETAILS     SearchType = 3
	SEARCHTYPE_HYPERVISOR_UPLOAD_INFO SearchType = 4
	SEARCHTYPE_VALIDATE_BUNDLE_INFO   SearchType = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *SearchType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UNCONFIGURED_NODES",
		"NETWORKING_DETAILS",
		"HYPERVISOR_UPLOAD_INFO",
		"VALIDATE_BUNDLE_INFO",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e SearchType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UNCONFIGURED_NODES",
		"NETWORKING_DETAILS",
		"HYPERVISOR_UPLOAD_INFO",
		"VALIDATE_BUNDLE_INFO",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *SearchType) index(name string) SearchType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UNCONFIGURED_NODES",
		"NETWORKING_DETAILS",
		"HYPERVISOR_UPLOAD_INFO",
		"VALIDATE_BUNDLE_INFO",
	}
	for idx := range names {
		if names[idx] == name {
			return SearchType(idx)
		}
	}
	return SEARCHTYPE_UNKNOWN
}

func (e *SearchType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for SearchType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *SearchType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e SearchType) Ref() *SearchType {
	return &e
}

/*
SMTP network details.
*/
type SmtpNetwork struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	IpAddress *import4.IPAddress `json:"ipAddress"`
	/*
	  SMTP server password.
	*/
	Password *string `json:"password,omitempty"`
	/*
	  SMTP port.
	*/
	Port *int `json:"port,omitempty"`
	/*
	  SMTP server user name.
	*/
	Username *string `json:"username,omitempty"`
}

func (p *SmtpNetwork) MarshalJSON() ([]byte, error) {
	type SmtpNetworkProxy SmtpNetwork
	return json.Marshal(struct {
		*SmtpNetworkProxy
		IpAddress *import4.IPAddress `json:"ipAddress,omitempty"`
	}{
		SmtpNetworkProxy: (*SmtpNetworkProxy)(p),
		IpAddress:        p.IpAddress,
	})
}

func NewSmtpNetwork() *SmtpNetwork {
	p := new(SmtpNetwork)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SmtpNetwork"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
SMTP servers on a cluster. This is part of payload for cluster update operation only.
*/
type SmtpServerRef struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  SMTP email address.
	*/
	EmailAddress *string `json:"emailAddress"`

	Server *SmtpNetwork `json:"server"`

	Type *SmtpType `json:"type,omitempty"`
}

func (p *SmtpServerRef) MarshalJSON() ([]byte, error) {
	type SmtpServerRefProxy SmtpServerRef
	return json.Marshal(struct {
		*SmtpServerRefProxy
		EmailAddress *string      `json:"emailAddress,omitempty"`
		Server       *SmtpNetwork `json:"server,omitempty"`
	}{
		SmtpServerRefProxy: (*SmtpServerRefProxy)(p),
		EmailAddress:       p.EmailAddress,
		Server:             p.Server,
	})
}

func NewSmtpServerRef() *SmtpServerRef {
	p := new(SmtpServerRef)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SmtpServerRef"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Type of SMTP server.
*/
type SmtpType int

const (
	SMTPTYPE_UNKNOWN  SmtpType = 0
	SMTPTYPE_REDACTED SmtpType = 1
	SMTPTYPE_PLAIN    SmtpType = 2
	SMTPTYPE_STARTTLS SmtpType = 3
	SMTPTYPE_SSL      SmtpType = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *SmtpType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PLAIN",
		"STARTTLS",
		"SSL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e SmtpType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PLAIN",
		"STARTTLS",
		"SSL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *SmtpType) index(name string) SmtpType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PLAIN",
		"STARTTLS",
		"SSL",
	}
	for idx := range names {
		if names[idx] == name {
			return SmtpType(idx)
		}
	}
	return SMTPTYPE_UNKNOWN
}

func (e *SmtpType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for SmtpType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *SmtpType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e SmtpType) Ref() *SmtpType {
	return &e
}

/*
SNMP user authentication type.
*/
type SnmpAuthType int

const (
	SNMPAUTHTYPE_UNKNOWN  SnmpAuthType = 0
	SNMPAUTHTYPE_REDACTED SnmpAuthType = 1
	SNMPAUTHTYPE_MD5      SnmpAuthType = 2
	SNMPAUTHTYPE_SHA      SnmpAuthType = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *SnmpAuthType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"MD5",
		"SHA",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e SnmpAuthType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"MD5",
		"SHA",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *SnmpAuthType) index(name string) SnmpAuthType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"MD5",
		"SHA",
	}
	for idx := range names {
		if names[idx] == name {
			return SnmpAuthType(idx)
		}
	}
	return SNMPAUTHTYPE_UNKNOWN
}

func (e *SnmpAuthType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for SnmpAuthType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *SnmpAuthType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e SnmpAuthType) Ref() *SnmpAuthType {
	return &e
}

/*
SNMP information.
*/
type SnmpConfig struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  SNMP status.
	*/
	Status *bool `json:"status,omitempty"`
	/*
	  SNMP transport details.
	*/
	Transports []SnmpTransport `json:"transports,omitempty"`
	/*
	  SNMP trap details.
	*/
	Traps []SnmpTrap `json:"traps,omitempty"`
	/*
	  SNMP user information.
	*/
	Users []SnmpUser `json:"users,omitempty"`
}

func NewSnmpConfig() *SnmpConfig {
	p := new(SnmpConfig)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SnmpConfig"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
SNMP user encryption type.
*/
type SnmpPrivType int

const (
	SNMPPRIVTYPE_UNKNOWN  SnmpPrivType = 0
	SNMPPRIVTYPE_REDACTED SnmpPrivType = 1
	SNMPPRIVTYPE_DES      SnmpPrivType = 2
	SNMPPRIVTYPE_AES      SnmpPrivType = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *SnmpPrivType) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"DES",
		"AES",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e SnmpPrivType) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"DES",
		"AES",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *SnmpPrivType) index(name string) SnmpPrivType {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"DES",
		"AES",
	}
	for idx := range names {
		if names[idx] == name {
			return SnmpPrivType(idx)
		}
	}
	return SNMPPRIVTYPE_UNKNOWN
}

func (e *SnmpPrivType) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for SnmpPrivType:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *SnmpPrivType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e SnmpPrivType) Ref() *SnmpPrivType {
	return &e
}

/*
SNMP protocol type.
*/
type SnmpProtocol int

const (
	SNMPPROTOCOL_UNKNOWN  SnmpProtocol = 0
	SNMPPROTOCOL_REDACTED SnmpProtocol = 1
	SNMPPROTOCOL_UDP      SnmpProtocol = 2
	SNMPPROTOCOL_UDP6     SnmpProtocol = 3
	SNMPPROTOCOL_TCP      SnmpProtocol = 4
	SNMPPROTOCOL_TCP6     SnmpProtocol = 5
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *SnmpProtocol) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UDP",
		"UDP6",
		"TCP",
		"TCP6",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e SnmpProtocol) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UDP",
		"UDP6",
		"TCP",
		"TCP6",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *SnmpProtocol) index(name string) SnmpProtocol {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"UDP",
		"UDP6",
		"TCP",
		"TCP6",
	}
	for idx := range names {
		if names[idx] == name {
			return SnmpProtocol(idx)
		}
	}
	return SNMPPROTOCOL_UNKNOWN
}

func (e *SnmpProtocol) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for SnmpProtocol:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *SnmpProtocol) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e SnmpProtocol) Ref() *SnmpProtocol {
	return &e
}

/*
SNMP status.
*/
type SnmpStatusParam struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  SNMP user information.
	*/
	Status *bool `json:"status"`
}

func (p *SnmpStatusParam) MarshalJSON() ([]byte, error) {
	type SnmpStatusParamProxy SnmpStatusParam
	return json.Marshal(struct {
		*SnmpStatusParamProxy
		Status *bool `json:"status,omitempty"`
	}{
		SnmpStatusParamProxy: (*SnmpStatusParamProxy)(p),
		Status:               p.Status,
	})
}

func NewSnmpStatusParam() *SnmpStatusParam {
	p := new(SnmpStatusParam)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SnmpStatusParam"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
SNMP transport details.
*/
type SnmpTransport struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  SNMP port.
	*/
	Port *int `json:"port"`

	Protocol *SnmpProtocol `json:"protocol"`
}

func (p *SnmpTransport) MarshalJSON() ([]byte, error) {
	type SnmpTransportProxy SnmpTransport
	return json.Marshal(struct {
		*SnmpTransportProxy
		Port     *int          `json:"port,omitempty"`
		Protocol *SnmpProtocol `json:"protocol,omitempty"`
	}{
		SnmpTransportProxy: (*SnmpTransportProxy)(p),
		Port:               p.Port,
		Protocol:           p.Protocol,
	})
}

func NewSnmpTransport() *SnmpTransport {
	p := new(SnmpTransport)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SnmpTransport"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type SnmpTrap struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	Address *import4.IPAddress `json:"address"`
	/*
	  Community string(plaintext) for SNMP version 2.0.
	*/
	CommunityString *string `json:"communityString,omitempty"`
	/*
	  SNMP engine Id.
	*/
	EngineId *string `json:"engineId,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  SNMP information status.
	*/
	Inform *bool `json:"inform,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  SNMP port.
	*/
	Port *int `json:"port,omitempty"`

	Protocol *SnmpProtocol `json:"protocol,omitempty"`
	/*
	  SNMP receiver name.
	*/
	RecieverName *string `json:"recieverName,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/*
	  SNMP username. For SNMP trap v3 version, SNMP username is required parameter.
	*/
	Username *string `json:"username,omitempty"`

	Version *SnmpTrapVersion `json:"version"`
}

func (p *SnmpTrap) MarshalJSON() ([]byte, error) {
	type SnmpTrapProxy SnmpTrap
	return json.Marshal(struct {
		*SnmpTrapProxy
		Address *import4.IPAddress `json:"address,omitempty"`
		Version *SnmpTrapVersion   `json:"version,omitempty"`
	}{
		SnmpTrapProxy: (*SnmpTrapProxy)(p),
		Address:       p.Address,
		Version:       p.Version,
	})
}

func NewSnmpTrap() *SnmpTrap {
	p := new(SnmpTrap)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SnmpTrap"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
SNMP version.
*/
type SnmpTrapVersion int

const (
	SNMPTRAPVERSION_UNKNOWN  SnmpTrapVersion = 0
	SNMPTRAPVERSION_REDACTED SnmpTrapVersion = 1
	SNMPTRAPVERSION_V2       SnmpTrapVersion = 2
	SNMPTRAPVERSION_V3       SnmpTrapVersion = 3
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *SnmpTrapVersion) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"V2",
		"V3",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e SnmpTrapVersion) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"V2",
		"V3",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *SnmpTrapVersion) index(name string) SnmpTrapVersion {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"V2",
		"V3",
	}
	for idx := range names {
		if names[idx] == name {
			return SnmpTrapVersion(idx)
		}
	}
	return SNMPTRAPVERSION_UNKNOWN
}

func (e *SnmpTrapVersion) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for SnmpTrapVersion:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *SnmpTrapVersion) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e SnmpTrapVersion) Ref() *SnmpTrapVersion {
	return &e
}

/*
SNMP user information.
*/
type SnmpUser struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  SNMP user authentication key.
	*/
	AuthKey *string `json:"authKey"`

	AuthType *SnmpAuthType `json:"authType"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  SNMP user encryption key.
	*/
	PrivKey *string `json:"privKey,omitempty"`

	PrivType *SnmpPrivType `json:"privType,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/*
	  SNMP username. For SNMP trap v3 version, SNMP username is required parameter.
	*/
	Username *string `json:"username"`
}

func (p *SnmpUser) MarshalJSON() ([]byte, error) {
	type SnmpUserProxy SnmpUser
	return json.Marshal(struct {
		*SnmpUserProxy
		AuthKey  *string       `json:"authKey,omitempty"`
		AuthType *SnmpAuthType `json:"authType,omitempty"`
		Username *string       `json:"username,omitempty"`
	}{
		SnmpUserProxy: (*SnmpUserProxy)(p),
		AuthKey:       p.AuthKey,
		AuthType:      p.AuthType,
		Username:      p.Username,
	})
}

func NewSnmpUser() *SnmpUser {
	p := new(SnmpUser)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SnmpUser"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Cluster software version details.
*/
type SoftwareMapReference struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	SoftwareType *SoftwareTypeRef `json:"softwareType,omitempty"`
	/*
	  Software version.
	*/
	Version *string `json:"version,omitempty"`
}

func NewSoftwareMapReference() *SoftwareMapReference {
	p := new(SoftwareMapReference)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.SoftwareMapReference"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Software type.
*/
type SoftwareTypeRef int

const (
	SOFTWARETYPEREF_UNKNOWN       SoftwareTypeRef = 0
	SOFTWARETYPEREF_REDACTED      SoftwareTypeRef = 1
	SOFTWARETYPEREF_NOS           SoftwareTypeRef = 2
	SOFTWARETYPEREF_NCC           SoftwareTypeRef = 3
	SOFTWARETYPEREF_PRISM_CENTRAL SoftwareTypeRef = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *SoftwareTypeRef) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NOS",
		"NCC",
		"PRISM_CENTRAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e SoftwareTypeRef) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NOS",
		"NCC",
		"PRISM_CENTRAL",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *SoftwareTypeRef) index(name string) SoftwareTypeRef {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NOS",
		"NCC",
		"PRISM_CENTRAL",
	}
	for idx := range names {
		if names[idx] == name {
			return SoftwareTypeRef(idx)
		}
	}
	return SOFTWARETYPEREF_UNKNOWN
}

func (e *SoftwareTypeRef) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for SoftwareTypeRef:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *SoftwareTypeRef) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e SoftwareTypeRef) Ref() *SoftwareTypeRef {
	return &e
}

type StorageContainer struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Affinity host extId for RF 1 Storage Container.
	*/
	AffinityHostExtId *string `json:"affinityHostExtId,omitempty"`

	CacheDeduplication *CacheDeduplication `json:"cacheDeduplication,omitempty"`
	/*
	  extId of the Cluster owning the Storage Container.
	*/
	ClusterExtId *string `json:"clusterExtId,omitempty"`
	/*
	  Corresponding name of the Cluster owning the Storage Container instance.
	*/
	ClusterName *string `json:"clusterName,omitempty"`
	/*
	  The compression delay in seconds.
	*/
	CompressionDelaySecs *int `json:"compressionDelaySecs,omitempty"`
	/*
	  extId of the Storage Container.
	*/
	ContainerExtId *string `json:"containerExtId,omitempty"`

	ErasureCode *ErasureCodeStatus `json:"erasureCode,omitempty"`
	/*
	  Delay in performing ErasureCode for the current Container instance.
	*/
	ErasureCodeDelaySecs *int `json:"erasureCodeDelaySecs,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Indicates whether to prefer a higher Erasure Code fault domain.
	*/
	HasHigherEcFaultDomainPreference *bool `json:"hasHigherEcFaultDomainPreference,omitempty"`
	/*
	  Indicates whether the compression is enabled for the Container.
	*/
	IsCompressionEnabled *bool `json:"isCompressionEnabled,omitempty"`
	/*
	  Indicates whether the Container is encrypted or not.
	*/
	IsEncrypted *bool `json:"isEncrypted,omitempty"`
	/*
	  Indicates whether data written to this container should be inline erasure coded or not. This field is only considered when ErasureCoding is enabled.
	*/
	IsInlineEcEnabled *bool `json:"isInlineEcEnabled,omitempty"`
	/*
	  Indicates whether the Container is internal and is managed by Nutanix.
	*/
	IsInternal *bool `json:"isInternal,omitempty"`
	/*
	  Indicates if the Storage Container is marked for removal. This field is set when the Storage Container is about to be destroyed.
	*/
	IsMarkedForRemoval *bool `json:"isMarkedForRemoval,omitempty"`
	/*
	  Indicates whether the NFS whitelist is inherited from global config.
	*/
	IsNfsWhitelistInherited *bool `json:"isNfsWhitelistInherited,omitempty"`
	/*
	  Indicates whether the Container instance has software encryption enabled.
	*/
	IsSoftwareEncryptionEnabled *bool `json:"isSoftwareEncryptionEnabled,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Max capacity of the Container as defined by the user.
	*/
	LogicalAdvertisedCapacityBytes *int64 `json:"logicalAdvertisedCapacityBytes,omitempty"`
	/*
	  Total reserved size (in bytes) of the container (set by Admin). This also accounts for the container's replication factor. The actual reserved capacity of the container will be the maximum of explicitReservedCapacity and implicitReservedCapacity.
	*/
	LogicalExplicitReservedCapacityBytes *int64 `json:"logicalExplicitReservedCapacityBytes,omitempty"`
	/*
	  This is the summation of reservations provisioned on all vdisks in the container. The actual reserved capacity of the container will be the maximum of explicitReservedCapacity and implicitReservedCapacity.
	*/
	LogicalImplicitReservedCapacityBytes *int64 `json:"logicalImplicitReservedCapacityBytes,omitempty"`
	/*
	  Maximum physical capacity of the Storage Container in bytes.
	*/
	MaxCapacityBytes *int64 `json:"maxCapacityBytes,omitempty"`
	/*
	  Name of the Storage Container. Note that the name of Storage Container should be unique per cluster.
	*/
	Name *string `json:"name"`
	/*
	  List of NFS addresses which need to be whitelisted.
	*/
	NfsWhitelistAddress []import4.IPAddressOrFQDN `json:"nfsWhitelistAddress,omitempty"`

	OnDiskDedup *OnDiskDedup `json:"onDiskDedup,omitempty"`
	/*
	  extId of the owner.
	*/
	OwnerExtId *string `json:"ownerExtId,omitempty"`
	/*
	  Replication factor of the Storage Container.
	*/
	ReplicationFactor *int `json:"replicationFactor,omitempty"`
	/*
	  extId of the Storage Pool owning the Storage Container instance.
	*/
	StoragePoolExtId *string `json:"storagePoolExtId,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func (p *StorageContainer) MarshalJSON() ([]byte, error) {
	type StorageContainerProxy StorageContainer
	return json.Marshal(struct {
		*StorageContainerProxy
		Name *string `json:"name,omitempty"`
	}{
		StorageContainerProxy: (*StorageContainerProxy)(p),
		Name:                  p.Name,
	})
}

func NewStorageContainer() *StorageContainer {
	p := new(StorageContainer)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.StorageContainer"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Represents the Disk storage tier.
*/
type StorageTier int

const (
	STORAGETIER_UNKNOWN      StorageTier = 0
	STORAGETIER_REDACTED     StorageTier = 1
	STORAGETIER_SSD_PCIE     StorageTier = 2
	STORAGETIER_SSD_SATA     StorageTier = 3
	STORAGETIER_DAS_SATA     StorageTier = 4
	STORAGETIER_CLOUD        StorageTier = 5
	STORAGETIER_SSD_MEM_NVME StorageTier = 6
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *StorageTier) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SSD_PCIE",
		"SSD_SATA",
		"DAS_SATA",
		"CLOUD",
		"SSD_MEM_NVME",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e StorageTier) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SSD_PCIE",
		"SSD_SATA",
		"DAS_SATA",
		"CLOUD",
		"SSD_MEM_NVME",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *StorageTier) index(name string) StorageTier {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"SSD_PCIE",
		"SSD_SATA",
		"DAS_SATA",
		"CLOUD",
		"SSD_MEM_NVME",
	}
	for idx := range names {
		if names[idx] == name {
			return StorageTier(idx)
		}
	}
	return STORAGETIER_UNKNOWN
}

func (e *StorageTier) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for StorageTier:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *StorageTier) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e StorageTier) Ref() *StorageTier {
	return &e
}

/*
Disk storage Tier type.
*/
type StorageTierReference int

const (
	STORAGETIERREFERENCE_UNKNOWN  StorageTierReference = 0
	STORAGETIERREFERENCE_REDACTED StorageTierReference = 1
	STORAGETIERREFERENCE_PCIE_SSD StorageTierReference = 2
	STORAGETIERREFERENCE_SATA_SSD StorageTierReference = 3
	STORAGETIERREFERENCE_HDD      StorageTierReference = 4
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *StorageTierReference) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PCIE_SSD",
		"SATA_SSD",
		"HDD",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e StorageTierReference) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PCIE_SSD",
		"SATA_SSD",
		"HDD",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *StorageTierReference) index(name string) StorageTierReference {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PCIE_SSD",
		"SATA_SSD",
		"HDD",
	}
	for idx := range names {
		if names[idx] == name {
			return StorageTierReference(idx)
		}
	}
	return STORAGETIERREFERENCE_UNKNOWN
}

func (e *StorageTierReference) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for StorageTierReference:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *StorageTierReference) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e StorageTierReference) Ref() *StorageTierReference {
	return &e
}

/*
Message contains the component domain fault tolerance text details.
*/
type ToleranceMessage struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of tolerance message attributes.
	*/
	AttributeList []AttributeItem `json:"attributeList,omitempty"`
	/*
	  Message Id.
	*/
	Id *string `json:"id,omitempty"`
}

func NewToleranceMessage() *ToleranceMessage {
	p := new(ToleranceMessage)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ToleranceMessage"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
List of unconfigured nodes.
*/
type UnconfigureNodeDetails struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  List of unconfigured nodes.
	*/
	NodeList []UnconfiguredNodeListItem `json:"nodeList,omitempty"`
}

func NewUnconfigureNodeDetails() *UnconfigureNodeDetails {
	p := new(UnconfigureNodeDetails)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UnconfigureNodeDetails"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Map providing additional node attributes for the unconfigured node.
*/
type UnconfiguredNodeAttributeMap struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Indicates if cvm interface can work with 1 GIG NIC or not.
	*/
	CanWorkWith1GNic *bool `json:"canWorkWith1GNic,omitempty"`
	/*
	  Default workload.
	*/
	DefaultWorkload *string `json:"defaultWorkload,omitempty"`
	/*
	  Indicates whether the model is supported or not.
	*/
	IsModelSupported *bool `json:"isModelSupported,omitempty"`
	/*
	  LCM family name.
	*/
	LcmFamily *string `json:"lcmFamily,omitempty"`
	/*
	  Indicates whether the hypervisor is robo mixed or not.
	*/
	RoboMixedHypervisor *bool `json:"roboMixedHypervisor,omitempty"`
}

func NewUnconfiguredNodeAttributeMap() *UnconfiguredNodeAttributeMap {
	p := new(UnconfiguredNodeAttributeMap)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UnconfiguredNodeAttributeMap"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Unconfigured node details.
*/
type UnconfiguredNodeListItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Cluster arch.
	*/
	Arch *string `json:"arch,omitempty"`

	Attributes *UnconfiguredNodeAttributeMap `json:"attributes,omitempty"`
	/*
	  Cluster ID.
	*/
	ClusterId *string `json:"clusterId,omitempty"`
	/*
	  CPU type.
	*/
	CpuType []string `json:"cpuType,omitempty"`
	/*
	  Current CVM VLAN tag.
	*/
	CurrentCvmVlanTag *string `json:"currentCvmVlanTag,omitempty"`
	/*
	  Current network interface of a node.
	*/
	CurrentNetworkInterface *string `json:"currentNetworkInterface,omitempty"`

	CvmIp *import4.IPAddress `json:"cvmIp,omitempty"`
	/*
	  Foundation version.
	*/
	FoundationVersion *string `json:"foundationVersion,omitempty"`

	HostType *HostTypeEnum `json:"hostType,omitempty"`

	HypervisorIp *import4.IPAddress `json:"hypervisorIp,omitempty"`

	HypervisorType *HypervisorType `json:"hypervisorType,omitempty"`
	/*
	  Host version of the node.
	*/
	HypervisorVersion *string `json:"hypervisorVersion,omitempty"`
	/*
	  Interface IPV6 address.
	*/
	InterfaceIpv6 *string `json:"interfaceIpv6,omitempty"`

	IpmiIp *import4.IPAddress `json:"ipmiIp,omitempty"`
	/*
	  Secure boot status.
	*/
	IsSecureBooted *bool `json:"isSecureBooted,omitempty"`
	/*
	  Position of a node in a rackable unit.
	*/
	NodePosition *string `json:"nodePosition,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  NOS software version of a node.
	*/
	NosVersion *string `json:"nosVersion,omitempty"`
	/*
	  Maximum number of nodes in rackable-unit.
	*/
	RackableUnitMaxNodes *int64 `json:"rackableUnitMaxNodes,omitempty"`
	/*
	  Rackable unit model type.
	*/
	RackableUnitModel *string `json:"rackableUnitModel,omitempty"`
	/*
	  Rackable unit serial name.
	*/
	RackableUnitSerial *string `json:"rackableUnitSerial,omitempty"`
}

func NewUnconfiguredNodeListItem() *UnconfiguredNodeListItem {
	p := new(UnconfiguredNodeListItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UnconfiguredNodeListItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/storage-containers/{extId}/$actions/unmount Post operation
*/
type UnmountStorageContainerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUnmountStorageContainerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUnmountStorageContainerApiResponse() *UnmountStorageContainerApiResponse {
	p := new(UnmountStorageContainerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UnmountStorageContainerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UnmountStorageContainerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UnmountStorageContainerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUnmountStorageContainerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/vcenter-extensions/{extId}/$actions/unregister Post operation
*/
type UnregisterVcenterExtensionApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUnregisterVcenterExtensionApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUnregisterVcenterExtensionApiResponse() *UnregisterVcenterExtensionApiResponse {
	p := new(UnregisterVcenterExtensionApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UnregisterVcenterExtensionApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UnregisterVcenterExtensionApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UnregisterVcenterExtensionApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUnregisterVcenterExtensionApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId} Put operation
*/
type UpdateClusterApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateClusterApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateClusterApiResponse() *UpdateClusterApiResponse {
	p := new(UpdateClusterApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UpdateClusterApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateClusterApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateClusterApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateClusterApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/rsyslog-servers/{extId} Put operation
*/
type UpdateRsyslogServerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateRsyslogServerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateRsyslogServerApiResponse() *UpdateRsyslogServerApiResponse {
	p := new(UpdateRsyslogServerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UpdateRsyslogServerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateRsyslogServerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateRsyslogServerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateRsyslogServerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/snmp/$actions/update-status Post operation
*/
type UpdateSnmpStatusApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateSnmpStatusApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateSnmpStatusApiResponse() *UpdateSnmpStatusApiResponse {
	p := new(UpdateSnmpStatusApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UpdateSnmpStatusApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateSnmpStatusApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateSnmpStatusApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateSnmpStatusApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/snmp/traps/{extId} Put operation
*/
type UpdateSnmpTrapApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateSnmpTrapApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateSnmpTrapApiResponse() *UpdateSnmpTrapApiResponse {
	p := new(UpdateSnmpTrapApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UpdateSnmpTrapApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateSnmpTrapApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateSnmpTrapApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateSnmpTrapApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{clusterExtId}/snmp/users/{extId} Put operation
*/
type UpdateSnmpUserApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateSnmpUserApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateSnmpUserApiResponse() *UpdateSnmpUserApiResponse {
	p := new(UpdateSnmpUserApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UpdateSnmpUserApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateSnmpUserApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateSnmpUserApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateSnmpUserApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/storage-containers/{extId} Put operation
*/
type UpdateStorageContainerApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfUpdateStorageContainerApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewUpdateStorageContainerApiResponse() *UpdateStorageContainerApiResponse {
	p := new(UpdateStorageContainerApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UpdateStorageContainerApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *UpdateStorageContainerApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *UpdateStorageContainerApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfUpdateStorageContainerApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Upgrade status of a cluster.
*/
type UpgradeStatus int

const (
	UPGRADESTATUS_UNKNOWN     UpgradeStatus = 0
	UPGRADESTATUS_REDACTED    UpgradeStatus = 1
	UPGRADESTATUS_PENDING     UpgradeStatus = 2
	UPGRADESTATUS_DOWNLOADING UpgradeStatus = 3
	UPGRADESTATUS_QUEUED      UpgradeStatus = 4
	UPGRADESTATUS_PREUPGRADE  UpgradeStatus = 5
	UPGRADESTATUS_UPGRADING   UpgradeStatus = 6
	UPGRADESTATUS_SUCCEEDED   UpgradeStatus = 7
	UPGRADESTATUS_FAILED      UpgradeStatus = 8
	UPGRADESTATUS_CANCELLED   UpgradeStatus = 9
	UPGRADESTATUS_SCHEDULED   UpgradeStatus = 10
)

// Returns the name of the enum given an ordinal number
//
// Deprecated: Please use GetName instead of name
func (e *UpgradeStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PENDING",
		"DOWNLOADING",
		"QUEUED",
		"PREUPGRADE",
		"UPGRADING",
		"SUCCEEDED",
		"FAILED",
		"CANCELLED",
		"SCHEDULED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the name of the enum
func (e UpgradeStatus) GetName() string {
	index := int(e)
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PENDING",
		"DOWNLOADING",
		"QUEUED",
		"PREUPGRADE",
		"UPGRADING",
		"SUCCEEDED",
		"FAILED",
		"CANCELLED",
		"SCHEDULED",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// Returns the enum type given a string value
func (e *UpgradeStatus) index(name string) UpgradeStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"PENDING",
		"DOWNLOADING",
		"QUEUED",
		"PREUPGRADE",
		"UPGRADING",
		"SUCCEEDED",
		"FAILED",
		"CANCELLED",
		"SCHEDULED",
	}
	for idx := range names {
		if names[idx] == name {
			return UpgradeStatus(idx)
		}
	}
	return UPGRADESTATUS_UNKNOWN
}

func (e *UpgradeStatus) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for UpgradeStatus:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *UpgradeStatus) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e UpgradeStatus) Ref() *UpgradeStatus {
	return &e
}

/*
Uplink information for controller VM.
*/
type UplinkInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	CvmIp *import4.IPAddress `json:"cvmIp,omitempty"`
	/*
	  Uplink details for a controller VM.
	*/
	UplinkList []NameMacRef `json:"uplinkList,omitempty"`
}

func NewUplinkInfo() *UplinkInfo {
	p := new(UplinkInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UplinkInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Active and standby uplink information of the target nodes.
*/
type UplinkNetworkItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Name of the uplink.
	*/
	Name *string `json:"name,omitempty"`
	/*
	  List of network types.
	*/
	Networks []string `json:"networks,omitempty"`

	Uplinks *Uplinks `json:"uplinks,omitempty"`
}

func NewUplinkNetworkItem() *UplinkNetworkItem {
	p := new(UplinkNetworkItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UplinkNetworkItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Uplink information of the target nodes.
*/
type UplinkNode struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`

	CvmIp *import4.IPAddress `json:"cvmIp"`

	HypervisorIp *import4.IPAddress `json:"hypervisorIp,omitempty"`
	/*
	  Active and standby uplink information of the target nodes.
	*/
	Networks []UplinkNetworkItem `json:"networks"`
}

func (p *UplinkNode) MarshalJSON() ([]byte, error) {
	type UplinkNodeProxy UplinkNode
	return json.Marshal(struct {
		*UplinkNodeProxy
		CvmIp    *import4.IPAddress  `json:"cvmIp,omitempty"`
		Networks []UplinkNetworkItem `json:"networks,omitempty"`
	}{
		UplinkNodeProxy: (*UplinkNodeProxy)(p),
		CvmIp:           p.CvmIp,
		Networks:        p.Networks,
	})
}

func NewUplinkNode() *UplinkNode {
	p := new(UplinkNode)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UplinkNode"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Active and standby uplink information of the target nodes.
*/
type Uplinks struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Active uplink information.
	*/
	Active []UplinksField `json:"active,omitempty"`
	/*
	  Standby uplink information.
	*/
	Standby []UplinksField `json:"standby,omitempty"`
}

func NewUplinks() *Uplinks {
	p := new(Uplinks)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.Uplinks"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Properties of active and standby uplink.
*/
type UplinksField struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Mac address.
	*/
	Mac *string `json:"mac,omitempty"`
	/*
	  Interface name.
	*/
	Name *string `json:"name,omitempty"`
	/*
	  Interface value.
	*/
	Value *string `json:"value,omitempty"`
}

func NewUplinksField() *UplinksField {
	p := new(UplinksField)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UplinksField"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Upload information for a node.
*/
type UploadInfoNodeItem struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Error message if any, for available hypervisor ISO.
	*/
	AvailableHypervisorIsoError *string `json:"availableHypervisorIsoError,omitempty"`
	/*
	  Name of the hypervisor bundle.
	*/
	BundleName *string `json:"bundleName,omitempty"`
	/*
	  Provides information on whether hypervisor ISO upload is required or not. This API is not supported for XEN hypervisor type.
	*/
	HypervisorUploadRequired *bool `json:"hypervisorUploadRequired,omitempty"`
	/*
	  Indicates if imaging is required or not.
	*/
	IsImagingMandatory *bool `json:"isImagingMandatory,omitempty"`
	/*
	  Indicates if node is compatible or not.
	*/
	IsNodeCompatible *bool `json:"isNodeCompatible,omitempty"`
	/*
	  Md5sum of ISO.
	*/
	Md5Sum *string `json:"md5Sum,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`

	RequiredHypervisorType *HypervisorType `json:"requiredHypervisorType,omitempty"`
}

func NewUploadInfoNodeItem() *UploadInfoNodeItem {
	p := new(UploadInfoNodeItem)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UploadInfoNodeItem"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
UserName and Password model.
*/
type UserInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Cluster name. This is part of payload for both cluster create & update operations.
	*/
	ClusterName *string `json:"clusterName,omitempty"`
	/*
	  Password.
	*/
	Password *string `json:"password,omitempty"`
	/*
	  Username.
	*/
	UserName *string `json:"userName,omitempty"`
}

func NewUserInfo() *UserInfo {
	p := new(UserInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.UserInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Information of the bundle to verify.
*/
type ValidateBundleInfo struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Md5sum of ISO.
	*/
	Md5Sum *string `json:"md5Sum,omitempty"`
}

func NewValidateBundleInfo() *ValidateBundleInfo {
	p := new(ValidateBundleInfo)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ValidateBundleInfo"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/config/clusters/{extId}/$actions/validate-node Post operation
*/
type ValidateNodeApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfValidateNodeApiResponseData `json:"data,omitempty"`

	Metadata *import3.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewValidateNodeApiResponse() *ValidateNodeApiResponse {
	p := new(ValidateNodeApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ValidateNodeApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ValidateNodeApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ValidateNodeApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfValidateNodeApiResponseData()
	}
	e := p.Data.SetValue(v)
	if nil == e {
		if nil == p.DataItemDiscriminator_ {
			p.DataItemDiscriminator_ = new(string)
		}
		*p.DataItemDiscriminator_ = *p.Data.Discriminator
	}
	return e
}

/*
Request body for node validation. It can be OneOf between hypervisor bundle and node uplinks.
*/
type ValidateNodeParam struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	SpecItemDiscriminator_ *string `json:"$specItemDiscriminator,omitempty"`

	Spec *OneOfValidateNodeParamSpec `json:"spec"`
}

func (p *ValidateNodeParam) MarshalJSON() ([]byte, error) {
	type ValidateNodeParamProxy ValidateNodeParam
	return json.Marshal(struct {
		*ValidateNodeParamProxy
		Spec *OneOfValidateNodeParamSpec `json:"spec,omitempty"`
	}{
		ValidateNodeParamProxy: (*ValidateNodeParamProxy)(p),
		Spec:                   p.Spec,
	})
}

func NewValidateNodeParam() *ValidateNodeParam {
	p := new(ValidateNodeParam)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.ValidateNodeParam"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ValidateNodeParam) GetSpec() interface{} {
	if nil == p.Spec {
		return nil
	}
	return p.Spec.GetValue()
}

func (p *ValidateNodeParam) SetSpec(v interface{}) error {
	if nil == p.Spec {
		p.Spec = NewOneOfValidateNodeParamSpec()
	}
	e := p.Spec.SetValue(v)
	if nil == e {
		if nil == p.SpecItemDiscriminator_ {
			p.SpecItemDiscriminator_ = new(string)
		}
		*p.SpecItemDiscriminator_ = *p.Spec.Discriminator
	}
	return e
}

/*
Credentials for registering/unregistering vCenter Server extension for Nutanix to manage VMs of the cluster.
*/
type VcenterCredentials struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Password for vCenter Server extension registration/unregistration.
	*/
	Password *string `json:"password"`
	/*
	  vCenter port to connect for registering/unregistering extension.
	*/
	Port *int `json:"port,omitempty"`
	/*
	  Username for vCenter Server extension registration/unregistration.
	*/
	Username *string `json:"username"`
}

func (p *VcenterCredentials) MarshalJSON() ([]byte, error) {
	type VcenterCredentialsProxy VcenterCredentials
	return json.Marshal(struct {
		*VcenterCredentialsProxy
		Password *string `json:"password,omitempty"`
		Username *string `json:"username,omitempty"`
	}{
		VcenterCredentialsProxy: (*VcenterCredentialsProxy)(p),
		Password:                p.Password,
		Username:                p.Username,
	})
}

func NewVcenterCredentials() *VcenterCredentials {
	p := new(VcenterCredentials)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.VcenterCredentials"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	p.Port = new(int)
	*p.Port = 443

	return p
}

/*
vCenter Server extension information of the cluster. Nutanix Prism requires registering vCenter Server extension keys to be able to perform VM Management and other operations.
*/
type VcenterExtension struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The globally unique identifier of cluster instance. It should be of type UUID.
	*/
	ClusterExtId *string `json:"clusterExtId,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  IP Address of vCenter.
	*/
	IpAddress *string `json:"ipAddress,omitempty"`
	/*
	  Indicates whether the vCenter Server extension is registered for the cluster.
	*/
	IsRegistered *bool `json:"isRegistered,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewVcenterExtension() *VcenterExtension {
	p := new(VcenterExtension)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.VcenterExtension"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Virtual NIC details.
*/
type VirtualNic struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Status of DHCP protocol.
	*/
	DhcpEnabled *bool `json:"dhcpEnabled,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Host description.
	*/
	HostDescription *string `json:"hostDescription,omitempty"`
	/*
	  List of host NIC UUID associated with the host virtual NIC.
	*/
	HostNicsUuids []string `json:"hostNicsUuids,omitempty"`
	/*
	  Operational status of the interface to the port associated with the NIC entity.
	*/
	InterfaceStatus *string `json:"interfaceStatus,omitempty"`
	/*
	  List of IPv4 addresses associated with the NIC entity for the network connection.
	*/
	Ipv4Addresses []import4.IPAddress `json:"ipv4Addresses,omitempty"`
	/*
	  List of IPv6 addresses associated with the NIC entity for the network connection.
	*/
	Ipv6Addresses []import4.IPAddress `json:"ipv6Addresses,omitempty"`
	/*
	  Link speed in Kbps.
	*/
	LinkSpeedInKbps *int64 `json:"linkSpeedInKbps,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import3.ApiLink `json:"links,omitempty"`
	/*
	  Host Mac address.
	*/
	MacAddress *string `json:"macAddress,omitempty"`
	/*
	  Maximum transmission unit in bytes.
	*/
	MtuInBytes *int64 `json:"mtuInBytes,omitempty"`
	/*
	  Virtual NIC name.
	*/
	Name *string `json:"name,omitempty"`
	/*
	  UUID of the host.
	*/
	NodeUuid *string `json:"nodeUuid,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/*
	  VLAN Id.
	*/
	VlanId *int64 `json:"vlanId,omitempty"`
}

func NewVirtualNic() *VirtualNic {
	p := new(VirtualNic)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.config.VirtualNic"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type OneOfGetDiskApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType0    *Disk                  `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetDiskApiResponseData() *OneOfGetDiskApiResponseData {
	p := new(OneOfGetDiskApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetDiskApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetDiskApiResponseData is nil"))
	}
	switch v.(type) {
	case Disk:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(Disk)
		}
		*p.oneOfType0 = v.(Disk)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetDiskApiResponseData) GetValue() interface{} {
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetDiskApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType0 := new(Disk)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "clustermgmt.v4.config.Disk" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(Disk)
			}
			*p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType0.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType0.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetDiskApiResponseData"))
}

func (p *OneOfGetDiskApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetDiskApiResponseData")
}

type OneOfDiscoverUnconfiguredNodesApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfDiscoverUnconfiguredNodesApiResponseData() *OneOfDiscoverUnconfiguredNodesApiResponseData {
	p := new(OneOfDiscoverUnconfiguredNodesApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDiscoverUnconfiguredNodesApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDiscoverUnconfiguredNodesApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDiscoverUnconfiguredNodesApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfDiscoverUnconfiguredNodesApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDiscoverUnconfiguredNodesApiResponseData"))
}

func (p *OneOfDiscoverUnconfiguredNodesApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfDiscoverUnconfiguredNodesApiResponseData")
}

type OneOfMountStorageContainerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfMountStorageContainerApiResponseData() *OneOfMountStorageContainerApiResponseData {
	p := new(OneOfMountStorageContainerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfMountStorageContainerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfMountStorageContainerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfMountStorageContainerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfMountStorageContainerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfMountStorageContainerApiResponseData"))
}

func (p *OneOfMountStorageContainerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfMountStorageContainerApiResponseData")
}

type OneOfSearchResponseResponse struct {
	Discriminator *string                 `json:"-"`
	ObjectType_   *string                 `json:"-"`
	oneOfType2002 *NodeNetworkingDetails  `json:"-"`
	oneOfType2004 *ValidateBundleInfo     `json:"-"`
	oneOfType2003 *HypervisorUploadInfo   `json:"-"`
	oneOfType2001 *UnconfigureNodeDetails `json:"-"`
}

func NewOneOfSearchResponseResponse() *OneOfSearchResponseResponse {
	p := new(OneOfSearchResponseResponse)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfSearchResponseResponse) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfSearchResponseResponse is nil"))
	}
	switch v.(type) {
	case NodeNetworkingDetails:
		if nil == p.oneOfType2002 {
			p.oneOfType2002 = new(NodeNetworkingDetails)
		}
		*p.oneOfType2002 = v.(NodeNetworkingDetails)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2002.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2002.ObjectType_
	case ValidateBundleInfo:
		if nil == p.oneOfType2004 {
			p.oneOfType2004 = new(ValidateBundleInfo)
		}
		*p.oneOfType2004 = v.(ValidateBundleInfo)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2004.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2004.ObjectType_
	case HypervisorUploadInfo:
		if nil == p.oneOfType2003 {
			p.oneOfType2003 = new(HypervisorUploadInfo)
		}
		*p.oneOfType2003 = v.(HypervisorUploadInfo)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2003.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2003.ObjectType_
	case UnconfigureNodeDetails:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(UnconfigureNodeDetails)
		}
		*p.oneOfType2001 = v.(UnconfigureNodeDetails)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfSearchResponseResponse) GetValue() interface{} {
	if p.oneOfType2002 != nil && *p.oneOfType2002.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2002
	}
	if p.oneOfType2004 != nil && *p.oneOfType2004.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2004
	}
	if p.oneOfType2003 != nil && *p.oneOfType2003.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2003
	}
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	return nil
}

func (p *OneOfSearchResponseResponse) UnmarshalJSON(b []byte) error {
	vOneOfType2002 := new(NodeNetworkingDetails)
	if err := json.Unmarshal(b, vOneOfType2002); err == nil {
		if "clustermgmt.v4.config.NodeNetworkingDetails" == *vOneOfType2002.ObjectType_ {
			if nil == p.oneOfType2002 {
				p.oneOfType2002 = new(NodeNetworkingDetails)
			}
			*p.oneOfType2002 = *vOneOfType2002
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2002.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2002.ObjectType_
			return nil
		}
	}
	vOneOfType2004 := new(ValidateBundleInfo)
	if err := json.Unmarshal(b, vOneOfType2004); err == nil {
		if "clustermgmt.v4.config.ValidateBundleInfo" == *vOneOfType2004.ObjectType_ {
			if nil == p.oneOfType2004 {
				p.oneOfType2004 = new(ValidateBundleInfo)
			}
			*p.oneOfType2004 = *vOneOfType2004
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2004.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2004.ObjectType_
			return nil
		}
	}
	vOneOfType2003 := new(HypervisorUploadInfo)
	if err := json.Unmarshal(b, vOneOfType2003); err == nil {
		if "clustermgmt.v4.config.HypervisorUploadInfo" == *vOneOfType2003.ObjectType_ {
			if nil == p.oneOfType2003 {
				p.oneOfType2003 = new(HypervisorUploadInfo)
			}
			*p.oneOfType2003 = *vOneOfType2003
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2003.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2003.ObjectType_
			return nil
		}
	}
	vOneOfType2001 := new(UnconfigureNodeDetails)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.UnconfigureNodeDetails" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(UnconfigureNodeDetails)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfSearchResponseResponse"))
}

func (p *OneOfSearchResponseResponse) MarshalJSON() ([]byte, error) {
	if p.oneOfType2002 != nil && *p.oneOfType2002.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2002)
	}
	if p.oneOfType2004 != nil && *p.oneOfType2004.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2004)
	}
	if p.oneOfType2003 != nil && *p.oneOfType2003.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2003)
	}
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	return nil, errors.New("No value to marshal for OneOfSearchResponseResponse")
}

type OneOfListClustersApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []Cluster              `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListClustersApiResponseData() *OneOfListClustersApiResponseData {
	p := new(OneOfListClustersApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListClustersApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListClustersApiResponseData is nil"))
	}
	switch v.(type) {
	case []Cluster:
		p.oneOfType2001 = v.([]Cluster)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.Cluster>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.Cluster>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListClustersApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.Cluster>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListClustersApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]Cluster)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.Cluster" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.Cluster>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.Cluster>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListClustersApiResponseData"))
}

func (p *OneOfListClustersApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.Cluster>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListClustersApiResponseData")
}

type OneOfListVcenterExtensionsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []VcenterExtension     `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListVcenterExtensionsApiResponseData() *OneOfListVcenterExtensionsApiResponseData {
	p := new(OneOfListVcenterExtensionsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListVcenterExtensionsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListVcenterExtensionsApiResponseData is nil"))
	}
	switch v.(type) {
	case []VcenterExtension:
		p.oneOfType2001 = v.([]VcenterExtension)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.VcenterExtension>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.VcenterExtension>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListVcenterExtensionsApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.VcenterExtension>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListVcenterExtensionsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]VcenterExtension)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.VcenterExtension" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.VcenterExtension>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.VcenterExtension>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListVcenterExtensionsApiResponseData"))
}

func (p *OneOfListVcenterExtensionsApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.VcenterExtension>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListVcenterExtensionsApiResponseData")
}

type OneOfCreateClusterApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfCreateClusterApiResponseData() *OneOfCreateClusterApiResponseData {
	p := new(OneOfCreateClusterApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreateClusterApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreateClusterApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCreateClusterApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfCreateClusterApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreateClusterApiResponseData"))
}

func (p *OneOfCreateClusterApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfCreateClusterApiResponseData")
}

type OneOfUpdateSnmpTrapApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUpdateSnmpTrapApiResponseData() *OneOfUpdateSnmpTrapApiResponseData {
	p := new(OneOfUpdateSnmpTrapApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateSnmpTrapApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateSnmpTrapApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdateSnmpTrapApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUpdateSnmpTrapApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateSnmpTrapApiResponseData"))
}

func (p *OneOfUpdateSnmpTrapApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateSnmpTrapApiResponseData")
}

type OneOfListHostsByClusterIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []Host                 `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListHostsByClusterIdApiResponseData() *OneOfListHostsByClusterIdApiResponseData {
	p := new(OneOfListHostsByClusterIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListHostsByClusterIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListHostsByClusterIdApiResponseData is nil"))
	}
	switch v.(type) {
	case []Host:
		p.oneOfType2001 = v.([]Host)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.Host>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.Host>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListHostsByClusterIdApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.Host>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListHostsByClusterIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]Host)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.Host" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.Host>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.Host>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListHostsByClusterIdApiResponseData"))
}

func (p *OneOfListHostsByClusterIdApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.Host>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListHostsByClusterIdApiResponseData")
}

type OneOfUpdateSnmpUserApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUpdateSnmpUserApiResponseData() *OneOfUpdateSnmpUserApiResponseData {
	p := new(OneOfUpdateSnmpUserApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateSnmpUserApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateSnmpUserApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdateSnmpUserApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUpdateSnmpUserApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateSnmpUserApiResponseData"))
}

func (p *OneOfUpdateSnmpUserApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateSnmpUserApiResponseData")
}

type OneOfUpdateRsyslogServerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUpdateRsyslogServerApiResponseData() *OneOfUpdateRsyslogServerApiResponseData {
	p := new(OneOfUpdateRsyslogServerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateRsyslogServerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateRsyslogServerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdateRsyslogServerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUpdateRsyslogServerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateRsyslogServerApiResponseData"))
}

func (p *OneOfUpdateRsyslogServerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateRsyslogServerApiResponseData")
}

type OneOfListRsyslogServersByClusterIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []RsyslogServer        `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListRsyslogServersByClusterIdApiResponseData() *OneOfListRsyslogServersByClusterIdApiResponseData {
	p := new(OneOfListRsyslogServersByClusterIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListRsyslogServersByClusterIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListRsyslogServersByClusterIdApiResponseData is nil"))
	}
	switch v.(type) {
	case []RsyslogServer:
		p.oneOfType2001 = v.([]RsyslogServer)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.RsyslogServer>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.RsyslogServer>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListRsyslogServersByClusterIdApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.RsyslogServer>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListRsyslogServersByClusterIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]RsyslogServer)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.RsyslogServer" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.RsyslogServer>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.RsyslogServer>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListRsyslogServersByClusterIdApiResponseData"))
}

func (p *OneOfListRsyslogServersByClusterIdApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.RsyslogServer>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListRsyslogServersByClusterIdApiResponseData")
}

type OneOfGetSnmpConfigByClusterIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *SnmpConfig            `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetSnmpConfigByClusterIdApiResponseData() *OneOfGetSnmpConfigByClusterIdApiResponseData {
	p := new(OneOfGetSnmpConfigByClusterIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetSnmpConfigByClusterIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetSnmpConfigByClusterIdApiResponseData is nil"))
	}
	switch v.(type) {
	case SnmpConfig:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(SnmpConfig)
		}
		*p.oneOfType2001 = v.(SnmpConfig)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetSnmpConfigByClusterIdApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetSnmpConfigByClusterIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(SnmpConfig)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.SnmpConfig" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(SnmpConfig)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetSnmpConfigByClusterIdApiResponseData"))
}

func (p *OneOfGetSnmpConfigByClusterIdApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetSnmpConfigByClusterIdApiResponseData")
}

type OneOfDeleteClusterApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfDeleteClusterApiResponseData() *OneOfDeleteClusterApiResponseData {
	p := new(OneOfDeleteClusterApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeleteClusterApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeleteClusterApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDeleteClusterApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfDeleteClusterApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeleteClusterApiResponseData"))
}

func (p *OneOfDeleteClusterApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfDeleteClusterApiResponseData")
}

type OneOfAddSnmpTransportsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfAddSnmpTransportsApiResponseData() *OneOfAddSnmpTransportsApiResponseData {
	p := new(OneOfAddSnmpTransportsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfAddSnmpTransportsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfAddSnmpTransportsApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfAddSnmpTransportsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfAddSnmpTransportsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfAddSnmpTransportsApiResponseData"))
}

func (p *OneOfAddSnmpTransportsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfAddSnmpTransportsApiResponseData")
}

type OneOfCreateSnmpUserApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfCreateSnmpUserApiResponseData() *OneOfCreateSnmpUserApiResponseData {
	p := new(OneOfCreateSnmpUserApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreateSnmpUserApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreateSnmpUserApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCreateSnmpUserApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfCreateSnmpUserApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreateSnmpUserApiResponseData"))
}

func (p *OneOfCreateSnmpUserApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfCreateSnmpUserApiResponseData")
}

type OneOfFetchTaskApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *SearchResponse        `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfFetchTaskApiResponseData() *OneOfFetchTaskApiResponseData {
	p := new(OneOfFetchTaskApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfFetchTaskApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfFetchTaskApiResponseData is nil"))
	}
	switch v.(type) {
	case SearchResponse:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(SearchResponse)
		}
		*p.oneOfType2001 = v.(SearchResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfFetchTaskApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfFetchTaskApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(SearchResponse)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.SearchResponse" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(SearchResponse)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfFetchTaskApiResponseData"))
}

func (p *OneOfFetchTaskApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfFetchTaskApiResponseData")
}

type OneOfValidateNodeApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfValidateNodeApiResponseData() *OneOfValidateNodeApiResponseData {
	p := new(OneOfValidateNodeApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfValidateNodeApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfValidateNodeApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfValidateNodeApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfValidateNodeApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfValidateNodeApiResponseData"))
}

func (p *OneOfValidateNodeApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfValidateNodeApiResponseData")
}

type OneOfCreateSnmpTrapApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfCreateSnmpTrapApiResponseData() *OneOfCreateSnmpTrapApiResponseData {
	p := new(OneOfCreateSnmpTrapApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreateSnmpTrapApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreateSnmpTrapApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCreateSnmpTrapApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfCreateSnmpTrapApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreateSnmpTrapApiResponseData"))
}

func (p *OneOfCreateSnmpTrapApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfCreateSnmpTrapApiResponseData")
}

type OneOfUpdateSnmpStatusApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUpdateSnmpStatusApiResponseData() *OneOfUpdateSnmpStatusApiResponseData {
	p := new(OneOfUpdateSnmpStatusApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateSnmpStatusApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateSnmpStatusApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdateSnmpStatusApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUpdateSnmpStatusApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateSnmpStatusApiResponseData"))
}

func (p *OneOfUpdateSnmpStatusApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateSnmpStatusApiResponseData")
}

type OneOfCreateRsyslogServerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfCreateRsyslogServerApiResponseData() *OneOfCreateRsyslogServerApiResponseData {
	p := new(OneOfCreateRsyslogServerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreateRsyslogServerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreateRsyslogServerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCreateRsyslogServerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfCreateRsyslogServerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreateRsyslogServerApiResponseData"))
}

func (p *OneOfCreateRsyslogServerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfCreateRsyslogServerApiResponseData")
}

type OneOfGetFaultToleranceStatusByClusterIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []DomainFaultTolerance `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetFaultToleranceStatusByClusterIdApiResponseData() *OneOfGetFaultToleranceStatusByClusterIdApiResponseData {
	p := new(OneOfGetFaultToleranceStatusByClusterIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetFaultToleranceStatusByClusterIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetFaultToleranceStatusByClusterIdApiResponseData is nil"))
	}
	switch v.(type) {
	case []DomainFaultTolerance:
		p.oneOfType2001 = v.([]DomainFaultTolerance)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.DomainFaultTolerance>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.DomainFaultTolerance>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetFaultToleranceStatusByClusterIdApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.DomainFaultTolerance>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetFaultToleranceStatusByClusterIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]DomainFaultTolerance)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.DomainFaultTolerance" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.DomainFaultTolerance>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.DomainFaultTolerance>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetFaultToleranceStatusByClusterIdApiResponseData"))
}

func (p *OneOfGetFaultToleranceStatusByClusterIdApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.DomainFaultTolerance>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetFaultToleranceStatusByClusterIdApiResponseData")
}

type OneOfDeleteSnmpUserApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfDeleteSnmpUserApiResponseData() *OneOfDeleteSnmpUserApiResponseData {
	p := new(OneOfDeleteSnmpUserApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeleteSnmpUserApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeleteSnmpUserApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDeleteSnmpUserApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfDeleteSnmpUserApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeleteSnmpUserApiResponseData"))
}

func (p *OneOfDeleteSnmpUserApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfDeleteSnmpUserApiResponseData")
}

type OneOfGetVcenterExtensionApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *VcenterExtension      `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetVcenterExtensionApiResponseData() *OneOfGetVcenterExtensionApiResponseData {
	p := new(OneOfGetVcenterExtensionApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetVcenterExtensionApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetVcenterExtensionApiResponseData is nil"))
	}
	switch v.(type) {
	case VcenterExtension:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(VcenterExtension)
		}
		*p.oneOfType2001 = v.(VcenterExtension)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetVcenterExtensionApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetVcenterExtensionApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(VcenterExtension)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.VcenterExtension" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(VcenterExtension)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetVcenterExtensionApiResponseData"))
}

func (p *OneOfGetVcenterExtensionApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetVcenterExtensionApiResponseData")
}

type OneOfGetRackableUnitApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *RackableUnit          `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetRackableUnitApiResponseData() *OneOfGetRackableUnitApiResponseData {
	p := new(OneOfGetRackableUnitApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetRackableUnitApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetRackableUnitApiResponseData is nil"))
	}
	switch v.(type) {
	case RackableUnit:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(RackableUnit)
		}
		*p.oneOfType2001 = v.(RackableUnit)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetRackableUnitApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetRackableUnitApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(RackableUnit)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.RackableUnit" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(RackableUnit)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetRackableUnitApiResponseData"))
}

func (p *OneOfGetRackableUnitApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetRackableUnitApiResponseData")
}

type OneOfRenameHostApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfRenameHostApiResponseData() *OneOfRenameHostApiResponseData {
	p := new(OneOfRenameHostApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfRenameHostApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfRenameHostApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfRenameHostApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfRenameHostApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfRenameHostApiResponseData"))
}

func (p *OneOfRenameHostApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfRenameHostApiResponseData")
}

type OneOfUpdateStorageContainerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUpdateStorageContainerApiResponseData() *OneOfUpdateStorageContainerApiResponseData {
	p := new(OneOfUpdateStorageContainerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateStorageContainerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateStorageContainerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdateStorageContainerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUpdateStorageContainerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateStorageContainerApiResponseData"))
}

func (p *OneOfUpdateStorageContainerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateStorageContainerApiResponseData")
}

type OneOfCreateStorageContainerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfCreateStorageContainerApiResponseData() *OneOfCreateStorageContainerApiResponseData {
	p := new(OneOfCreateStorageContainerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCreateStorageContainerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCreateStorageContainerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCreateStorageContainerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfCreateStorageContainerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCreateStorageContainerApiResponseData"))
}

func (p *OneOfCreateStorageContainerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfCreateStorageContainerApiResponseData")
}

type OneOfGetHostApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *Host                  `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetHostApiResponseData() *OneOfGetHostApiResponseData {
	p := new(OneOfGetHostApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetHostApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetHostApiResponseData is nil"))
	}
	switch v.(type) {
	case Host:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(Host)
		}
		*p.oneOfType2001 = v.(Host)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetHostApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetHostApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(Host)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.Host" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(Host)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetHostApiResponseData"))
}

func (p *OneOfGetHostApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetHostApiResponseData")
}

type OneOfListDisksApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType0    []Disk                 `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListDisksApiResponseData() *OneOfListDisksApiResponseData {
	p := new(OneOfListDisksApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListDisksApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListDisksApiResponseData is nil"))
	}
	switch v.(type) {
	case []Disk:
		p.oneOfType0 = v.([]Disk)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.Disk>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.Disk>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListDisksApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.Disk>" == *p.Discriminator {
		return p.oneOfType0
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListDisksApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType0 := new([]Disk)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {

		if len(*vOneOfType0) == 0 || "clustermgmt.v4.config.Disk" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.Disk>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.Disk>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListDisksApiResponseData"))
}

func (p *OneOfListDisksApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.Disk>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListDisksApiResponseData")
}

type OneOfGetSnmpTrapApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *SnmpTrap              `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetSnmpTrapApiResponseData() *OneOfGetSnmpTrapApiResponseData {
	p := new(OneOfGetSnmpTrapApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetSnmpTrapApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetSnmpTrapApiResponseData is nil"))
	}
	switch v.(type) {
	case SnmpTrap:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(SnmpTrap)
		}
		*p.oneOfType2001 = v.(SnmpTrap)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetSnmpTrapApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetSnmpTrapApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(SnmpTrap)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.SnmpTrap" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(SnmpTrap)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetSnmpTrapApiResponseData"))
}

func (p *OneOfGetSnmpTrapApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetSnmpTrapApiResponseData")
}

type OneOfFetchNodeNetworkingDetailsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfFetchNodeNetworkingDetailsApiResponseData() *OneOfFetchNodeNetworkingDetailsApiResponseData {
	p := new(OneOfFetchNodeNetworkingDetailsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfFetchNodeNetworkingDetailsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfFetchNodeNetworkingDetailsApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfFetchNodeNetworkingDetailsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfFetchNodeNetworkingDetailsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfFetchNodeNetworkingDetailsApiResponseData"))
}

func (p *OneOfFetchNodeNetworkingDetailsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfFetchNodeNetworkingDetailsApiResponseData")
}

type OneOfGetSnmpUserApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *SnmpUser              `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetSnmpUserApiResponseData() *OneOfGetSnmpUserApiResponseData {
	p := new(OneOfGetSnmpUserApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetSnmpUserApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetSnmpUserApiResponseData is nil"))
	}
	switch v.(type) {
	case SnmpUser:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(SnmpUser)
		}
		*p.oneOfType2001 = v.(SnmpUser)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetSnmpUserApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetSnmpUserApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(SnmpUser)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.SnmpUser" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(SnmpUser)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetSnmpUserApiResponseData"))
}

func (p *OneOfGetSnmpUserApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetSnmpUserApiResponseData")
}

type OneOfListRackableUnitsByClusterIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []RackableUnit         `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListRackableUnitsByClusterIdApiResponseData() *OneOfListRackableUnitsByClusterIdApiResponseData {
	p := new(OneOfListRackableUnitsByClusterIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListRackableUnitsByClusterIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListRackableUnitsByClusterIdApiResponseData is nil"))
	}
	switch v.(type) {
	case []RackableUnit:
		p.oneOfType2001 = v.([]RackableUnit)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.RackableUnit>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.RackableUnit>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListRackableUnitsByClusterIdApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.RackableUnit>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListRackableUnitsByClusterIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]RackableUnit)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.RackableUnit" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.RackableUnit>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.RackableUnit>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListRackableUnitsByClusterIdApiResponseData"))
}

func (p *OneOfListRackableUnitsByClusterIdApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.RackableUnit>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListRackableUnitsByClusterIdApiResponseData")
}

type OneOfGetClusterApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *Cluster               `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetClusterApiResponseData() *OneOfGetClusterApiResponseData {
	p := new(OneOfGetClusterApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetClusterApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetClusterApiResponseData is nil"))
	}
	switch v.(type) {
	case Cluster:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(Cluster)
		}
		*p.oneOfType2001 = v.(Cluster)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetClusterApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetClusterApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(Cluster)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.Cluster" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(Cluster)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetClusterApiResponseData"))
}

func (p *OneOfGetClusterApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetClusterApiResponseData")
}

type OneOfDeleteStorageContainerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfDeleteStorageContainerApiResponseData() *OneOfDeleteStorageContainerApiResponseData {
	p := new(OneOfDeleteStorageContainerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeleteStorageContainerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeleteStorageContainerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDeleteStorageContainerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfDeleteStorageContainerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeleteStorageContainerApiResponseData"))
}

func (p *OneOfDeleteStorageContainerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfDeleteStorageContainerApiResponseData")
}

type OneOfListDataStoresByClusterIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []DataStore            `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListDataStoresByClusterIdApiResponseData() *OneOfListDataStoresByClusterIdApiResponseData {
	p := new(OneOfListDataStoresByClusterIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListDataStoresByClusterIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListDataStoresByClusterIdApiResponseData is nil"))
	}
	switch v.(type) {
	case []DataStore:
		p.oneOfType2001 = v.([]DataStore)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.DataStore>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.DataStore>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListDataStoresByClusterIdApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.DataStore>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListDataStoresByClusterIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]DataStore)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.DataStore" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.DataStore>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.DataStore>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListDataStoresByClusterIdApiResponseData"))
}

func (p *OneOfListDataStoresByClusterIdApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.DataStore>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListDataStoresByClusterIdApiResponseData")
}

type OneOfListStorageContainersApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []StorageContainer     `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListStorageContainersApiResponseData() *OneOfListStorageContainersApiResponseData {
	p := new(OneOfListStorageContainersApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListStorageContainersApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListStorageContainersApiResponseData is nil"))
	}
	switch v.(type) {
	case []StorageContainer:
		p.oneOfType2001 = v.([]StorageContainer)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.StorageContainer>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.StorageContainer>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListStorageContainersApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.StorageContainer>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListStorageContainersApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]StorageContainer)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.StorageContainer" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.StorageContainer>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.StorageContainer>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListStorageContainersApiResponseData"))
}

func (p *OneOfListStorageContainersApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.StorageContainer>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListStorageContainersApiResponseData")
}

type OneOfUnregisterVcenterExtensionApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUnregisterVcenterExtensionApiResponseData() *OneOfUnregisterVcenterExtensionApiResponseData {
	p := new(OneOfUnregisterVcenterExtensionApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUnregisterVcenterExtensionApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUnregisterVcenterExtensionApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUnregisterVcenterExtensionApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUnregisterVcenterExtensionApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUnregisterVcenterExtensionApiResponseData"))
}

func (p *OneOfUnregisterVcenterExtensionApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUnregisterVcenterExtensionApiResponseData")
}

type OneOfGetRsyslogServerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *RsyslogServer         `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetRsyslogServerApiResponseData() *OneOfGetRsyslogServerApiResponseData {
	p := new(OneOfGetRsyslogServerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetRsyslogServerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetRsyslogServerApiResponseData is nil"))
	}
	switch v.(type) {
	case RsyslogServer:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(RsyslogServer)
		}
		*p.oneOfType2001 = v.(RsyslogServer)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetRsyslogServerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetRsyslogServerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(RsyslogServer)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.RsyslogServer" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(RsyslogServer)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetRsyslogServerApiResponseData"))
}

func (p *OneOfGetRsyslogServerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetRsyslogServerApiResponseData")
}

type OneOfDeleteRsyslogServerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfDeleteRsyslogServerApiResponseData() *OneOfDeleteRsyslogServerApiResponseData {
	p := new(OneOfDeleteRsyslogServerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeleteRsyslogServerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeleteRsyslogServerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDeleteRsyslogServerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfDeleteRsyslogServerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeleteRsyslogServerApiResponseData"))
}

func (p *OneOfDeleteRsyslogServerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfDeleteRsyslogServerApiResponseData")
}

type OneOfRemoveSnmpTransportsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfRemoveSnmpTransportsApiResponseData() *OneOfRemoveSnmpTransportsApiResponseData {
	p := new(OneOfRemoveSnmpTransportsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfRemoveSnmpTransportsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfRemoveSnmpTransportsApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfRemoveSnmpTransportsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfRemoveSnmpTransportsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfRemoveSnmpTransportsApiResponseData"))
}

func (p *OneOfRemoveSnmpTransportsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfRemoveSnmpTransportsApiResponseData")
}

type OneOfRemoveNodeApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfRemoveNodeApiResponseData() *OneOfRemoveNodeApiResponseData {
	p := new(OneOfRemoveNodeApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfRemoveNodeApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfRemoveNodeApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfRemoveNodeApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfRemoveNodeApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfRemoveNodeApiResponseData"))
}

func (p *OneOfRemoveNodeApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfRemoveNodeApiResponseData")
}

type OneOfRegisterVcenterExtensionApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfRegisterVcenterExtensionApiResponseData() *OneOfRegisterVcenterExtensionApiResponseData {
	p := new(OneOfRegisterVcenterExtensionApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfRegisterVcenterExtensionApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfRegisterVcenterExtensionApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfRegisterVcenterExtensionApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfRegisterVcenterExtensionApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfRegisterVcenterExtensionApiResponseData"))
}

func (p *OneOfRegisterVcenterExtensionApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfRegisterVcenterExtensionApiResponseData")
}

type OneOfCheckHypervisorRequirementsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfCheckHypervisorRequirementsApiResponseData() *OneOfCheckHypervisorRequirementsApiResponseData {
	p := new(OneOfCheckHypervisorRequirementsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfCheckHypervisorRequirementsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfCheckHypervisorRequirementsApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfCheckHypervisorRequirementsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfCheckHypervisorRequirementsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfCheckHypervisorRequirementsApiResponseData"))
}

func (p *OneOfCheckHypervisorRequirementsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfCheckHypervisorRequirementsApiResponseData")
}

type OneOfGetStorageContainerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *StorageContainer      `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetStorageContainerApiResponseData() *OneOfGetStorageContainerApiResponseData {
	p := new(OneOfGetStorageContainerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetStorageContainerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetStorageContainerApiResponseData is nil"))
	}
	switch v.(type) {
	case StorageContainer:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(StorageContainer)
		}
		*p.oneOfType2001 = v.(StorageContainer)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetStorageContainerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetStorageContainerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(StorageContainer)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.StorageContainer" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(StorageContainer)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetStorageContainerApiResponseData"))
}

func (p *OneOfGetStorageContainerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetStorageContainerApiResponseData")
}

type OneOfDeleteSnmpTrapApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfDeleteSnmpTrapApiResponseData() *OneOfDeleteSnmpTrapApiResponseData {
	p := new(OneOfDeleteSnmpTrapApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDeleteSnmpTrapApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDeleteSnmpTrapApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDeleteSnmpTrapApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfDeleteSnmpTrapApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDeleteSnmpTrapApiResponseData"))
}

func (p *OneOfDeleteSnmpTrapApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfDeleteSnmpTrapApiResponseData")
}

type OneOfListHostNicsByHostIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []HostNic              `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListHostNicsByHostIdApiResponseData() *OneOfListHostNicsByHostIdApiResponseData {
	p := new(OneOfListHostNicsByHostIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListHostNicsByHostIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListHostNicsByHostIdApiResponseData is nil"))
	}
	switch v.(type) {
	case []HostNic:
		p.oneOfType2001 = v.([]HostNic)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.HostNic>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.HostNic>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListHostNicsByHostIdApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.HostNic>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListHostNicsByHostIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]HostNic)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.HostNic" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.HostNic>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.HostNic>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListHostNicsByHostIdApiResponseData"))
}

func (p *OneOfListHostNicsByHostIdApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.HostNic>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListHostNicsByHostIdApiResponseData")
}

type OneOfUnmountStorageContainerApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUnmountStorageContainerApiResponseData() *OneOfUnmountStorageContainerApiResponseData {
	p := new(OneOfUnmountStorageContainerApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUnmountStorageContainerApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUnmountStorageContainerApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUnmountStorageContainerApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUnmountStorageContainerApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUnmountStorageContainerApiResponseData"))
}

func (p *OneOfUnmountStorageContainerApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUnmountStorageContainerApiResponseData")
}

type OneOfListHostsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []Host                 `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListHostsApiResponseData() *OneOfListHostsApiResponseData {
	p := new(OneOfListHostsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListHostsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListHostsApiResponseData is nil"))
	}
	switch v.(type) {
	case []Host:
		p.oneOfType2001 = v.([]Host)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.Host>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.Host>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListHostsApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.Host>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListHostsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]Host)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.Host" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.Host>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.Host>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListHostsApiResponseData"))
}

func (p *OneOfListHostsApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.Host>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListHostsApiResponseData")
}

type OneOfGetVirtualNicApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *VirtualNic            `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetVirtualNicApiResponseData() *OneOfGetVirtualNicApiResponseData {
	p := new(OneOfGetVirtualNicApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetVirtualNicApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetVirtualNicApiResponseData is nil"))
	}
	switch v.(type) {
	case VirtualNic:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(VirtualNic)
		}
		*p.oneOfType2001 = v.(VirtualNic)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetVirtualNicApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetVirtualNicApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(VirtualNic)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.VirtualNic" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(VirtualNic)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetVirtualNicApiResponseData"))
}

func (p *OneOfGetVirtualNicApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetVirtualNicApiResponseData")
}

type OneOfUpdateClusterApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfUpdateClusterApiResponseData() *OneOfUpdateClusterApiResponseData {
	p := new(OneOfUpdateClusterApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfUpdateClusterApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfUpdateClusterApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfUpdateClusterApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfUpdateClusterApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfUpdateClusterApiResponseData"))
}

func (p *OneOfUpdateClusterApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfUpdateClusterApiResponseData")
}

type OneOfExpandClusterApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *import1.TaskReference `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfExpandClusterApiResponseData() *OneOfExpandClusterApiResponseData {
	p := new(OneOfExpandClusterApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfExpandClusterApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfExpandClusterApiResponseData is nil"))
	}
	switch v.(type) {
	case import1.TaskReference:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(import1.TaskReference)
		}
		*p.oneOfType2001 = v.(import1.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfExpandClusterApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfExpandClusterApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(import1.TaskReference)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(import1.TaskReference)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfExpandClusterApiResponseData"))
}

func (p *OneOfExpandClusterApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfExpandClusterApiResponseData")
}

type OneOfListVirtualNicsByHostIdApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []VirtualNic           `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfListVirtualNicsByHostIdApiResponseData() *OneOfListVirtualNicsByHostIdApiResponseData {
	p := new(OneOfListVirtualNicsByHostIdApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListVirtualNicsByHostIdApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListVirtualNicsByHostIdApiResponseData is nil"))
	}
	switch v.(type) {
	case []VirtualNic:
		p.oneOfType2001 = v.([]VirtualNic)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.VirtualNic>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.VirtualNic>"
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfListVirtualNicsByHostIdApiResponseData) GetValue() interface{} {
	if "List<clustermgmt.v4.config.VirtualNic>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListVirtualNicsByHostIdApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]VirtualNic)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "clustermgmt.v4.config.VirtualNic" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.VirtualNic>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.VirtualNic>"
			return nil

		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListVirtualNicsByHostIdApiResponseData"))
}

func (p *OneOfListVirtualNicsByHostIdApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<clustermgmt.v4.config.VirtualNic>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListVirtualNicsByHostIdApiResponseData")
}

type OneOfValidateNodeParamSpec struct {
	Discriminator *string      `json:"-"`
	ObjectType_   *string      `json:"-"`
	oneOfType2001 *BundleParam `json:"-"`
	oneOfType2002 []UplinkNode `json:"-"`
}

func NewOneOfValidateNodeParamSpec() *OneOfValidateNodeParamSpec {
	p := new(OneOfValidateNodeParamSpec)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfValidateNodeParamSpec) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfValidateNodeParamSpec is nil"))
	}
	switch v.(type) {
	case BundleParam:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(BundleParam)
		}
		*p.oneOfType2001 = v.(BundleParam)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case []UplinkNode:
		p.oneOfType2002 = v.([]UplinkNode)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<clustermgmt.v4.config.UplinkNode>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<clustermgmt.v4.config.UplinkNode>"
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfValidateNodeParamSpec) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if "List<clustermgmt.v4.config.UplinkNode>" == *p.Discriminator {
		return p.oneOfType2002
	}
	return nil
}

func (p *OneOfValidateNodeParamSpec) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(BundleParam)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.BundleParam" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(BundleParam)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType2002 := new([]UplinkNode)
	if err := json.Unmarshal(b, vOneOfType2002); err == nil {

		if len(*vOneOfType2002) == 0 || "clustermgmt.v4.config.UplinkNode" == *((*vOneOfType2002)[0].ObjectType_) {
			p.oneOfType2002 = *vOneOfType2002
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<clustermgmt.v4.config.UplinkNode>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<clustermgmt.v4.config.UplinkNode>"
			return nil

		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfValidateNodeParamSpec"))
}

func (p *OneOfValidateNodeParamSpec) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if "List<clustermgmt.v4.config.UplinkNode>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2002)
	}
	return nil, errors.New("No value to marshal for OneOfValidateNodeParamSpec")
}

type OneOfGetHostNicApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *HostNic               `json:"-"`
	oneOfType400  *import2.ErrorResponse `json:"-"`
}

func NewOneOfGetHostNicApiResponseData() *OneOfGetHostNicApiResponseData {
	p := new(OneOfGetHostNicApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetHostNicApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetHostNicApiResponseData is nil"))
	}
	switch v.(type) {
	case HostNic:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(HostNic)
		}
		*p.oneOfType2001 = v.(HostNic)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import2.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import2.ErrorResponse)
		}
		*p.oneOfType400 = v.(import2.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetHostNicApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetHostNicApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(HostNic)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.config.HostNic" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(HostNic)
			}
			*p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType2001.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType2001.ObjectType_
			return nil
		}
	}
	vOneOfType400 := new(import2.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import2.ErrorResponse)
			}
			*p.oneOfType400 = *vOneOfType400
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = *p.oneOfType400.ObjectType_
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = *p.oneOfType400.ObjectType_
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetHostNicApiResponseData"))
}

func (p *OneOfGetHostNicApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetHostNicApiResponseData")
}

type FileDetail struct {
	Path        *string `json:"-"`
	ObjectType_ *string `json:"-"`
}

func NewFileDetail() *FileDetail {
	p := new(FileDetail)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "FileDetail"

	return p
}
