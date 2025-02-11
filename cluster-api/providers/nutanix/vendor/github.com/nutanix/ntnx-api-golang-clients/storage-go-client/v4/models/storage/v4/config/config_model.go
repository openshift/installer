/*
 * Generated file models/storage/v4/config/config_model.go.
 *
 * Product version: 4.0.2-alpha-3
 *
 * Part of the Nutanix Storage Versioned APIs
 *
 * (c) 2023 Nutanix Inc.  All rights reserved
 *
 */

/*
  Configure storage entities such as containers
*/
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	import2 "github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/models/common/v1/response"
	import4 "github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/models/prism/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/storage-go-client/v4/models/storage/v4/error"
)

type DataStore struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  Maximum capacity of the storage container.
	*/
	Capacity *int64 `json:"capacity,omitempty"`
	/**
	  Uuid of the storage container.
	*/
	ContainerExtId *string `json:"containerExtId,omitempty"`
	/**
	  Id of the storage container instance.
	*/
	ContainerId *string `json:"containerId,omitempty"`
	/**
	  Name of the storage container.
	*/
	ContainerName *string `json:"containerName"`
	/**
	  Name of the datastore.
	*/
	DatastoreName *string `json:"datastoreName,omitempty"`
	/**
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/**
	  freeSpace in the datastore.
	*/
	FreeSpace *int64 `json:"freeSpace,omitempty"`
	/**
	  Uuid of the host for datastore.
	*/
	HostExtId *string `json:"hostExtId,omitempty"`
	/**
	  Uuid of the host for datastore.
	*/
	HostId *string `json:"hostId,omitempty"`
	/**
	  Ip of the host for datastore.
	*/
	HostIpAddress *string `json:"hostIpAddress,omitempty"`
	/**
	  A HATEOAS style link for the response.  Each link contains a user friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/**
	  A globally unique identifier that represents the tenant that owns this entity.  It is automatically assigned by the system and is immutable from an API consumer perspective (some use cases may cause this Id to change - for instance a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/**
	  List of VMs name in the datastore.
	*/
	VmNames []string `json:"vmNames,omitempty"`
}

func (p *DataStore) MarshalJSON() ([]byte, error) {
	type DataStoreProxy DataStore
	return json.Marshal(struct {
		*DataStoreProxy
		ContainerName *string `json:"containerName,omitempty"`
	}{
		DataStoreProxy: (*DataStoreProxy)(p),
		ContainerName:  p.ContainerName,
	})
}

func NewDataStore() *DataStore {
	p := new(DataStore)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "storage.v4.config.DataStore"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.DataStore"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/**
create NFS datastores on the ESX hosts.
*/
type DataStoreMount struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  Name of the storage container.
	*/
	ContainerName *string `json:"containerName"`
	/**
	  Name of the datastore.
	*/
	DatastoreName *string `json:"datastoreName,omitempty"`
	/**
	  The Uuids of the nodes where the NFS datastore have to be created.
	*/
	NodeExtIds []string `json:"nodeExtIds,omitempty"`
	/**
	  The Zeus config ids of the nodes where the NFS datastore have to be created.
	*/
	NodeIds []string `json:"nodeIds,omitempty"`
	/**
	  if the host system have only read-only access to the NFS share (container).
	*/
	ReadOnly *bool `json:"readOnly,omitempty"`
	/**
	  The target path on which to mount the NFS datastore. KVM-only.
	*/
	TargetPath *string `json:"targetPath,omitempty"`
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
	*p.ObjectType_ = "storage.v4.config.DataStoreMount"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.DataStoreMount"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/**
REST response for all response codes in api path /storage/v4.0.a3/config/storage-containers/datastores Get operation
*/
type DataStoreResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDataStoreResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDataStoreResponse() *DataStoreResponse {
	p := new(DataStoreResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "storage.v4.config.DataStoreResponse"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.DataStoreResponse"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DataStoreResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DataStoreResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDataStoreResponseData()
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

/**
REST response for all response codes in api path /storage/v4.0.a3/config/storage-containers/{containerExtId}/$actions/unmount Post operation
*/
type DataStoreTaskResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfDataStoreTaskResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewDataStoreTaskResponse() *DataStoreTaskResponse {
	p := new(DataStoreTaskResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "storage.v4.config.DataStoreTaskResponse"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.DataStoreTaskResponse"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *DataStoreTaskResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *DataStoreTaskResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfDataStoreTaskResponseData()
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

/**
create NFS datastores on the ESX hosts.
*/
type DataStoreUnmount struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  Name of the datastore.
	*/
	DatastoreName *string `json:"datastoreName"`
	/**
	  The Uuids of the nodes where the NFS datastore have to be created.
	*/
	NodeExtIds []string `json:"nodeExtIds,omitempty"`
	/**
	  The Zeus config ids of the nodes where the NFS datastore have to be created.
	*/
	NodeIds []string `json:"nodeIds,omitempty"`
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
	*p.ObjectType_ = "storage.v4.config.DataStoreUnmount"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.DataStoreUnmount"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/**
Indicates possbile values for erasure code.
*/
type ErasureCodeStatus int

const (
	ERASURECODESTATUS_UNKNOWN  ErasureCodeStatus = 0
	ERASURECODESTATUS_REDACTED ErasureCodeStatus = 1
	ERASURECODESTATUS_NONE     ErasureCodeStatus = 2
	ERASURECODESTATUS_FALSE    ErasureCodeStatus = 3
	ERASURECODESTATUS_TRUE     ErasureCodeStatus = 4
)

// returns the name of the enum given an ordinal number
func (e *ErasureCodeStatus) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"false",
		"true",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// returns the enum type given a string value
func (e *ErasureCodeStatus) index(name string) ErasureCodeStatus {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"false",
		"true",
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

/**
Indicates possbile values which can be set to finger print on write.
*/
type FingerPrintOnWrite int

const (
	FINGERPRINTONWRITE_UNKNOWN  FingerPrintOnWrite = 0
	FINGERPRINTONWRITE_REDACTED FingerPrintOnWrite = 1
	FINGERPRINTONWRITE_NONE     FingerPrintOnWrite = 2
	FINGERPRINTONWRITE_FALSE    FingerPrintOnWrite = 3
	FINGERPRINTONWRITE_TRUE     FingerPrintOnWrite = 4
)

// returns the name of the enum given an ordinal number
func (e *FingerPrintOnWrite) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"false",
		"true",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// returns the enum type given a string value
func (e *FingerPrintOnWrite) index(name string) FingerPrintOnWrite {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"false",
		"true",
	}
	for idx := range names {
		if names[idx] == name {
			return FingerPrintOnWrite(idx)
		}
	}
	return FINGERPRINTONWRITE_UNKNOWN
}

func (e *FingerPrintOnWrite) UnmarshalJSON(b []byte) error {
	var enumStr string
	if err := json.Unmarshal(b, &enumStr); err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal for FingerPrintOnWrite:%s", err))
	}
	*e = e.index(enumStr)
	return nil
}

func (e *FingerPrintOnWrite) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(e.name(int(*e)))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

func (e FingerPrintOnWrite) Ref() *FingerPrintOnWrite {
	return &e
}

/**
Indicates possbile values for on disk deduplication.
*/
type OnDiskDedup int

const (
	ONDISKDEDUP_UNKNOWN      OnDiskDedup = 0
	ONDISKDEDUP_REDACTED     OnDiskDedup = 1
	ONDISKDEDUP_NONE         OnDiskDedup = 2
	ONDISKDEDUP_FALSE        OnDiskDedup = 3
	ONDISKDEDUP_POST_PROCESS OnDiskDedup = 4
)

// returns the name of the enum given an ordinal number
func (e *OnDiskDedup) name(index int) string {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"false",
		"POST_PROCESS",
	}
	if index < 0 || index >= len(names) {
		return "$UNKNOWN"
	}
	return names[index]
}

// returns the enum type given a string value
func (e *OnDiskDedup) index(name string) OnDiskDedup {
	names := [...]string{
		"$UNKNOWN",
		"$REDACTED",
		"NONE",
		"false",
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

type StorageContainer struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**
	  Total advertised capacity of the storage container.
	*/
	AdvertisedCapacity *int64 `json:"advertisedCapacity,omitempty"`
	/**
	  Affinity host id for RF 1 container.
	*/
	AffinityHostUuid *string `json:"affinityHostUuid,omitempty"`
	/**
	  Owning cluster uuid of storage container.
	*/
	ClusterExtId *string `json:"clusterExtId,omitempty"`
	/**
	  Compression delay in seconds.
	*/
	CompressionDelayInSecs *int `json:"compressionDelayInSecs,omitempty"`
	/**
	  Whether compression is enabled.
	*/
	CompressionEnabled *bool `json:"compressionEnabled,omitempty"`
	/**
	  Uuid of the storage container.
	*/
	ContainerExtId *string `json:"containerExtId,omitempty"`
	/**
	  Id of the storage container instance.
	*/
	ContainerId *string `json:"containerId,omitempty"`
	/**
	  Map of down migrate time in seconds for random io preference tier.
	*/
	DownMigrateTimesInSecs map[string]int `json:"downMigrateTimesInSecs,omitempty"`
	/**
	  Whether container to enable software encryption.
	*/
	EnableSoftwareEncryption *bool `json:"enableSoftwareEncryption,omitempty"`
	/**
	  Whether container is encrypted or not.
	*/
	Encrypted *bool `json:"encrypted,omitempty"`

	ErasureCode *ErasureCodeStatus `json:"erasureCode,omitempty"`
	/**
	  Erasure code delay in seconds.
	*/
	ErasureCodeDelaySecs *int `json:"erasureCodeDelaySecs,omitempty"`
	/**
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`

	FingerPrintOnWrite *FingerPrintOnWrite `json:"fingerPrintOnWrite,omitempty"`
	/**
	  Whether inline erasure coding is enabled.
	*/
	InlineEcEnabled *bool `json:"inlineEcEnabled,omitempty"`
	/**
	  Whether Nfs whitelist inherited from global config.
	*/
	IsNfsWhitelistInherited *bool `json:"isNfsWhitelistInherited,omitempty"`
	/**
	  Whether container is managed by nutanix.
	*/
	IsNutanixManaged *bool `json:"isNutanixManaged,omitempty"`
	/**
	  A HATEOAS style link for the response.  Each link contains a user friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/**
	  Map of remote containers.
	*/
	MappedRemoteContainers map[string]string `json:"mappedRemoteContainers,omitempty"`
	/**
	  Whether storage container is marked for removal.
	*/
	MarkedForRemoval *bool `json:"markedForRemoval,omitempty"`
	/**
	  Maximum capacity of the storage container.
	*/
	MaxCapacity *int64 `json:"maxCapacity,omitempty"`
	/**
	  Name of the storage container.
	*/
	Name *string `json:"name"`
	/**
	  List of Nfs addresses which needs to be whitelisted.
	*/
	NfsWhitelistAddress []import2.IPAddressOrFQDN `json:"nfsWhitelistAddress,omitempty"`

	OnDiskDedup *OnDiskDedup `json:"onDiskDedup,omitempty"`
	/**
	  Oplog replication factor of the storage container.
	*/
	OplogReplicationFactor *int `json:"oplogReplicationFactor,omitempty"`
	/**
	  Uuid of the storage container.
	*/
	OwnerUuid *string `json:"ownerUuid,omitempty"`
	/**
	  Whether to prefer higher erasure code fault domain.
	*/
	PreferHigherECFaultDomain *bool `json:"preferHigherECFaultDomain,omitempty"`
	/**
	  List of random IO preference tier.
	*/
	RandomIoPreference []string `json:"randomIoPreference,omitempty"`
	/**
	  Replication factor of the storage container.
	*/
	ReplicationFactor *int `json:"replicationFactor,omitempty"`
	/**
	  List of sequential IO preference tier.
	*/
	SeqIoPreference []string `json:"seqIoPreference,omitempty"`
	/**
	  Owning storage pool uuid of the container instance.
	*/
	StoragePoolUuid *string `json:"storagePoolUuid,omitempty"`
	/**
	  A globally unique identifier that represents the tenant that owns this entity.  It is automatically assigned by the system and is immutable from an API consumer perspective (some use cases may cause this Id to change - for instance a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
	/**
	  Total explicit reserved capacity of the storage container.
	*/
	TotalExplicitReservedCapacity *int64 `json:"totalExplicitReservedCapacity,omitempty"`
	/**
	  Total implicit reserved capacity of the storage container.
	*/
	TotalImplicitReservedCapacity *int64 `json:"totalImplicitReservedCapacity,omitempty"`
	/**
	  List of volume stores in the container.
	*/
	VstoreNameList []string `json:"vstoreNameList,omitempty"`
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
	*p.ObjectType_ = "storage.v4.config.StorageContainer"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.StorageContainer"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/**
REST response for all response codes in api path /storage/v4.0.a3/config/storage-containers/{containerExtId} Get operation
*/
type StorageContainerResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfStorageContainerResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewStorageContainerResponse() *StorageContainerResponse {
	p := new(StorageContainerResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "storage.v4.config.StorageContainerResponse"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.StorageContainerResponse"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *StorageContainerResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *StorageContainerResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfStorageContainerResponseData()
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

/**
REST response for all response codes in api path /storage/v4.0.a3/config/storage-containers/{containerExtId}/$actions/mount Post operation
*/
type StorageContainerTaskResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfStorageContainerTaskResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewStorageContainerTaskResponse() *StorageContainerTaskResponse {
	p := new(StorageContainerTaskResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "storage.v4.config.StorageContainerTaskResponse"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.StorageContainerTaskResponse"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *StorageContainerTaskResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *StorageContainerTaskResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfStorageContainerTaskResponseData()
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

/**
REST response for all response codes in api path /storage/v4.0.a3/config/storage-containers Get operation
*/
type StorageContainersResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/**

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfStorageContainersResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewStorageContainersResponse() *StorageContainersResponse {
	p := new(StorageContainersResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "storage.v4.config.StorageContainersResponse"
	p.Reserved_ = map[string]interface{}{"$fqObjectType": "storage.v4.r0.a3.config.StorageContainersResponse"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *StorageContainersResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *StorageContainersResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfStorageContainersResponseData()
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

type OneOfStorageContainerResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType1    []import2.Message      `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
	oneOfType0    *StorageContainer      `json:"-"`
}

func NewOneOfStorageContainerResponseData() *OneOfStorageContainerResponseData {
	p := new(OneOfStorageContainerResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfStorageContainerResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfStorageContainerResponseData is nil"))
	}
	switch v.(type) {
	case []import2.Message:
		p.oneOfType1 = v.([]import2.Message)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<common.v1.config.Message>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<common.v1.config.Message>"
	case import3.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import3.ErrorResponse)
		}
		*p.oneOfType400 = v.(import3.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case StorageContainer:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(StorageContainer)
		}
		*p.oneOfType0 = v.(StorageContainer)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfStorageContainerResponseData) GetValue() interface{} {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return p.oneOfType1
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfStorageContainerResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType1 := new([]import2.Message)
	if err := json.Unmarshal(b, vOneOfType1); err == nil {
		if len(*vOneOfType1) == 0 || "common.v1.config.Message" == *((*vOneOfType1)[0].ObjectType_) {
			p.oneOfType1 = *vOneOfType1
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<common.v1.config.Message>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<common.v1.config.Message>"
			return nil

		}
	}
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "storage.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import3.ErrorResponse)
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
	vOneOfType0 := new(StorageContainer)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "storage.v4.config.StorageContainer" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(StorageContainer)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfStorageContainerResponseData"))
}

func (p *OneOfStorageContainerResponseData) MarshalJSON() ([]byte, error) {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return json.Marshal(p.oneOfType1)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfStorageContainerResponseData")
}

type OneOfStorageContainerTaskResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType1    []import2.Message      `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
	oneOfType0    *import4.TaskReference `json:"-"`
}

func NewOneOfStorageContainerTaskResponseData() *OneOfStorageContainerTaskResponseData {
	p := new(OneOfStorageContainerTaskResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfStorageContainerTaskResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfStorageContainerTaskResponseData is nil"))
	}
	switch v.(type) {
	case []import2.Message:
		p.oneOfType1 = v.([]import2.Message)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<common.v1.config.Message>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<common.v1.config.Message>"
	case import3.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import3.ErrorResponse)
		}
		*p.oneOfType400 = v.(import3.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import4.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import4.TaskReference)
		}
		*p.oneOfType0 = v.(import4.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfStorageContainerTaskResponseData) GetValue() interface{} {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return p.oneOfType1
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfStorageContainerTaskResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType1 := new([]import2.Message)
	if err := json.Unmarshal(b, vOneOfType1); err == nil {
		if len(*vOneOfType1) == 0 || "common.v1.config.Message" == *((*vOneOfType1)[0].ObjectType_) {
			p.oneOfType1 = *vOneOfType1
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<common.v1.config.Message>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<common.v1.config.Message>"
			return nil

		}
	}
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "storage.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import3.ErrorResponse)
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
	vOneOfType0 := new(import4.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import4.TaskReference)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfStorageContainerTaskResponseData"))
}

func (p *OneOfStorageContainerTaskResponseData) MarshalJSON() ([]byte, error) {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return json.Marshal(p.oneOfType1)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfStorageContainerTaskResponseData")
}

type OneOfDataStoreTaskResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType1    []import2.Message      `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
	oneOfType0    *import4.TaskReference `json:"-"`
}

func NewOneOfDataStoreTaskResponseData() *OneOfDataStoreTaskResponseData {
	p := new(OneOfDataStoreTaskResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDataStoreTaskResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDataStoreTaskResponseData is nil"))
	}
	switch v.(type) {
	case []import2.Message:
		p.oneOfType1 = v.([]import2.Message)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<common.v1.config.Message>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<common.v1.config.Message>"
	case import3.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import3.ErrorResponse)
		}
		*p.oneOfType400 = v.(import3.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case import4.TaskReference:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(import4.TaskReference)
		}
		*p.oneOfType0 = v.(import4.TaskReference)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType0.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType0.ObjectType_
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDataStoreTaskResponseData) GetValue() interface{} {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return p.oneOfType1
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfDataStoreTaskResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType1 := new([]import2.Message)
	if err := json.Unmarshal(b, vOneOfType1); err == nil {
		if len(*vOneOfType1) == 0 || "common.v1.config.Message" == *((*vOneOfType1)[0].ObjectType_) {
			p.oneOfType1 = *vOneOfType1
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<common.v1.config.Message>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<common.v1.config.Message>"
			return nil

		}
	}
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "storage.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import3.ErrorResponse)
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
	vOneOfType0 := new(import4.TaskReference)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "prism.v4.config.TaskReference" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(import4.TaskReference)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDataStoreTaskResponseData"))
}

func (p *OneOfDataStoreTaskResponseData) MarshalJSON() ([]byte, error) {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return json.Marshal(p.oneOfType1)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfDataStoreTaskResponseData")
}

type OneOfDataStoreResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType1    []import2.Message      `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
	oneOfType0    []DataStore            `json:"-"`
}

func NewOneOfDataStoreResponseData() *OneOfDataStoreResponseData {
	p := new(OneOfDataStoreResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfDataStoreResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfDataStoreResponseData is nil"))
	}
	switch v.(type) {
	case []import2.Message:
		p.oneOfType1 = v.([]import2.Message)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<common.v1.config.Message>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<common.v1.config.Message>"
	case import3.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import3.ErrorResponse)
		}
		*p.oneOfType400 = v.(import3.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case []DataStore:
		p.oneOfType0 = v.([]DataStore)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<storage.v4.config.DataStore>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<storage.v4.config.DataStore>"
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfDataStoreResponseData) GetValue() interface{} {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return p.oneOfType1
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if "List<storage.v4.config.DataStore>" == *p.Discriminator {
		return p.oneOfType0
	}
	return nil
}

func (p *OneOfDataStoreResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType1 := new([]import2.Message)
	if err := json.Unmarshal(b, vOneOfType1); err == nil {
		if len(*vOneOfType1) == 0 || "common.v1.config.Message" == *((*vOneOfType1)[0].ObjectType_) {
			p.oneOfType1 = *vOneOfType1
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<common.v1.config.Message>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<common.v1.config.Message>"
			return nil

		}
	}
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "storage.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import3.ErrorResponse)
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
	vOneOfType0 := new([]DataStore)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if len(*vOneOfType0) == 0 || "storage.v4.config.DataStore" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<storage.v4.config.DataStore>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<storage.v4.config.DataStore>"
			return nil

		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfDataStoreResponseData"))
}

func (p *OneOfDataStoreResponseData) MarshalJSON() ([]byte, error) {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return json.Marshal(p.oneOfType1)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if "List<storage.v4.config.DataStore>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfDataStoreResponseData")
}

type OneOfStorageContainersResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType1    []import2.Message      `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
	oneOfType0    []StorageContainer     `json:"-"`
}

func NewOneOfStorageContainersResponseData() *OneOfStorageContainersResponseData {
	p := new(OneOfStorageContainersResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfStorageContainersResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfStorageContainersResponseData is nil"))
	}
	switch v.(type) {
	case []import2.Message:
		p.oneOfType1 = v.([]import2.Message)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<common.v1.config.Message>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<common.v1.config.Message>"
	case import3.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import3.ErrorResponse)
		}
		*p.oneOfType400 = v.(import3.ErrorResponse)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType400.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType400.ObjectType_
	case []StorageContainer:
		p.oneOfType0 = v.([]StorageContainer)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<storage.v4.config.StorageContainer>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<storage.v4.config.StorageContainer>"
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfStorageContainersResponseData) GetValue() interface{} {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return p.oneOfType1
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if "List<storage.v4.config.StorageContainer>" == *p.Discriminator {
		return p.oneOfType0
	}
	return nil
}

func (p *OneOfStorageContainersResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType1 := new([]import2.Message)
	if err := json.Unmarshal(b, vOneOfType1); err == nil {
		if len(*vOneOfType1) == 0 || "common.v1.config.Message" == *((*vOneOfType1)[0].ObjectType_) {
			p.oneOfType1 = *vOneOfType1
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<common.v1.config.Message>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<common.v1.config.Message>"
			return nil

		}
	}
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "storage.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import3.ErrorResponse)
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
	vOneOfType0 := new([]StorageContainer)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if len(*vOneOfType0) == 0 || "storage.v4.config.StorageContainer" == *((*vOneOfType0)[0].ObjectType_) {
			p.oneOfType0 = *vOneOfType0
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<storage.v4.config.StorageContainer>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<storage.v4.config.StorageContainer>"
			return nil

		}
	}
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfStorageContainersResponseData"))
}

func (p *OneOfStorageContainersResponseData) MarshalJSON() ([]byte, error) {
	if "List<common.v1.config.Message>" == *p.Discriminator {
		return json.Marshal(p.oneOfType1)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if "List<storage.v4.config.StorageContainer>" == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfStorageContainersResponseData")
}
