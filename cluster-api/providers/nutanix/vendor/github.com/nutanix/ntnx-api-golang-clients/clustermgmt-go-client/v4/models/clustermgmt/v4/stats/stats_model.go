/*
 * Generated file models/clustermgmt/v4/stats/stats_model.go.
 *
 * Product version: 4.0.1-beta-2
 *
 * Part of the Nutanix Clustermgmt Versioned APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module clustermgmt.v4.stats of Nutanix Clustermgmt Versioned APIs
*/
package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	import3 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/error"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/response"
	import2 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/stats"
	"time"
)

type ClusterStats struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpm []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpm,omitempty"`
	/*
	  Lower Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmLowerBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmUpperBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmUpperBuf,omitempty"`
	/*
	  Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecs []TimeValuePair `json:"controllerAvgIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecs []TimeValuePair `json:"controllerAvgReadIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecs []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller IOPS Number.
	*/
	ControllerNumIops []TimeValuePair `json:"controllerNumIops,omitempty"`
	/*
	  Lower Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsLowerBuf []TimeValuePair `json:"controllerNumIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsUpperBuf []TimeValuePair `json:"controllerNumIopsUpperBuf,omitempty"`
	/*
	  Number of controller read IOPS.
	*/
	ControllerNumReadIops []TimeValuePair `json:"controllerNumReadIops,omitempty"`
	/*
	  Lower Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsLowerBuf []TimeValuePair `json:"controllerNumReadIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsUpperBuf []TimeValuePair `json:"controllerNumReadIopsUpperBuf,omitempty"`
	/*
	  Number of controller write IOPS.
	*/
	ControllerNumWriteIops []TimeValuePair `json:"controllerNumWriteIops,omitempty"`
	/*
	  Lower Buf value of Number of controller write IoPS.
	*/
	ControllerNumWriteIopsLowerBuf []TimeValuePair `json:"controllerNumWriteIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller write IOPS.
	*/
	ControllerNumWriteIopsUpperBuf []TimeValuePair `json:"controllerNumWriteIopsUpperBuf,omitempty"`
	/*
	  Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbps []TimeValuePair `json:"controllerReadIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbps []TimeValuePair `json:"controllerWriteIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Free physical space(bytes).
	*/
	FreePhysicalStorageBytes []TimeValuePair `json:"freePhysicalStorageBytes,omitempty"`
	/*
	  Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpm []TimeValuePair `json:"hypervisorCpuUsagePpm,omitempty"`
	/*
	  Lower Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmLowerBuf []TimeValuePair `json:"hypervisorCpuUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmUpperBuf []TimeValuePair `json:"hypervisorCpuUsagePpmUpperBuf,omitempty"`
	/*
	  Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbps []TimeValuePair `json:"ioBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsLowerBuf []TimeValuePair `json:"ioBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsUpperBuf []TimeValuePair `json:"ioBandwidthKbpsUpperBuf,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  Logical storage usage(bytes).
	*/
	LogicalStorageUsageBytes []TimeValuePair `json:"logicalStorageUsageBytes,omitempty"`
	/*
	  Overall memory usage(bytes).
	*/
	OverallMemoryUsageBytes []TimeValuePair `json:"overallMemoryUsageBytes,omitempty"`

	StatType *import2.DownSamplingOperator `json:"statType,omitempty"`
	/*
	  Storage capacity(bytes).
	*/
	StorageCapacityBytes []TimeValuePair `json:"storageCapacityBytes,omitempty"`
	/*
	  Storage usage(bytes).
	*/
	StorageUsageBytes []TimeValuePair `json:"storageUsageBytes,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewClusterStats() *ClusterStats {
	p := new(ClusterStats)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.ClusterStats"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/stats/clusters/{clusterExtId} Get operation
*/
type ClusterStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfClusterStatsApiResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewClusterStatsApiResponse() *ClusterStatsApiResponse {
	p := new(ClusterStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.ClusterStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ClusterStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ClusterStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfClusterStatsApiResponseData()
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

type ClusterStatsProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpm []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpm,omitempty"`
	/*
	  Lower Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmLowerBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmUpperBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmUpperBuf,omitempty"`
	/*
	  Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecs []TimeValuePair `json:"controllerAvgIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecs []TimeValuePair `json:"controllerAvgReadIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecs []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller IOPS Number.
	*/
	ControllerNumIops []TimeValuePair `json:"controllerNumIops,omitempty"`
	/*
	  Lower Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsLowerBuf []TimeValuePair `json:"controllerNumIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsUpperBuf []TimeValuePair `json:"controllerNumIopsUpperBuf,omitempty"`
	/*
	  Number of controller read IOPS.
	*/
	ControllerNumReadIops []TimeValuePair `json:"controllerNumReadIops,omitempty"`
	/*
	  Lower Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsLowerBuf []TimeValuePair `json:"controllerNumReadIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsUpperBuf []TimeValuePair `json:"controllerNumReadIopsUpperBuf,omitempty"`
	/*
	  Number of controller write IOPS.
	*/
	ControllerNumWriteIops []TimeValuePair `json:"controllerNumWriteIops,omitempty"`
	/*
	  Lower Buf value of Number of controller write IoPS.
	*/
	ControllerNumWriteIopsLowerBuf []TimeValuePair `json:"controllerNumWriteIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller write IOPS.
	*/
	ControllerNumWriteIopsUpperBuf []TimeValuePair `json:"controllerNumWriteIopsUpperBuf,omitempty"`
	/*
	  Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbps []TimeValuePair `json:"controllerReadIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbps []TimeValuePair `json:"controllerWriteIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Free physical space(bytes).
	*/
	FreePhysicalStorageBytes []TimeValuePair `json:"freePhysicalStorageBytes,omitempty"`
	/*
	  Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpm []TimeValuePair `json:"hypervisorCpuUsagePpm,omitempty"`
	/*
	  Lower Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmLowerBuf []TimeValuePair `json:"hypervisorCpuUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmUpperBuf []TimeValuePair `json:"hypervisorCpuUsagePpmUpperBuf,omitempty"`
	/*
	  Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbps []TimeValuePair `json:"ioBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsLowerBuf []TimeValuePair `json:"ioBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsUpperBuf []TimeValuePair `json:"ioBandwidthKbpsUpperBuf,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  Logical storage usage(bytes).
	*/
	LogicalStorageUsageBytes []TimeValuePair `json:"logicalStorageUsageBytes,omitempty"`
	/*
	  Overall memory usage(bytes).
	*/
	OverallMemoryUsageBytes []TimeValuePair `json:"overallMemoryUsageBytes,omitempty"`

	StatType *import2.DownSamplingOperator `json:"statType,omitempty"`
	/*
	  Storage capacity(bytes).
	*/
	StorageCapacityBytes []TimeValuePair `json:"storageCapacityBytes,omitempty"`
	/*
	  Storage usage(bytes).
	*/
	StorageUsageBytes []TimeValuePair `json:"storageUsageBytes,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewClusterStatsProjection() *ClusterStatsProjection {
	p := new(ClusterStatsProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.ClusterStatsProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type DiskStats struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Average IO latency.
	*/
	DiskAvgIoLatencyMicrosec []import2.TimeIntValuePair `json:"diskAvgIoLatencyMicrosec,omitempty"`
	/*
	  Lower limit of data transfer that a Disk can handler per second.
	*/
	DiskBaseIoBandwidthkbps []import2.TimeIntValuePair `json:"diskBaseIoBandwidthkbps,omitempty"`
	/*
	  Lower limit of the latency of I/O operations that Disk can handle without exceeding its standard latency level.
	*/
	DiskBaseIoLatencyMicrosec []import2.TimeIntValuePair `json:"diskBaseIoLatencyMicrosec,omitempty"`
	/*
	  Lower limit of I/O operations that a Disk can perform per second.
	*/
	DiskBaseNumIops []import2.TimeIntValuePair `json:"diskBaseNumIops,omitempty"`
	/*
	  Lower buffer capacity average read IO latency, measured in microseconds (usecs).
	*/
	DiskBaseReadIoAvgLatencyMicrosec []import2.TimeIntValuePair `json:"diskBaseReadIoAvgLatencyMicrosec,omitempty"`
	/*
	  Lower buffer capacity for the amount of IO bandwidth that a Disk can handle read operations.
	*/
	DiskBaseReadIoBandwidthkbps []import2.TimeIntValuePair `json:"diskBaseReadIoBandwidthkbps,omitempty"`
	/*
	  Lower buffer capacity for the number of read IOPS that a Disk can handle.
	*/
	DiskBaseReadIops []import2.TimeIntValuePair `json:"diskBaseReadIops,omitempty"`
	/*
	  Lower buffer capacity average write IO latency, measured in microseconds (usecs).
	*/
	DiskBaseWriteIoAvgLatencyMicrosec []import2.TimeIntValuePair `json:"diskBaseWriteIoAvgLatencyMicrosec,omitempty"`
	/*
	  Lower buffer capacity for the amount of IO bandwidth that a Disk can handle write operations.
	*/
	DiskBaseWriteIoBandwidthkbps []import2.TimeIntValuePair `json:"diskBaseWriteIoBandwidthkbps,omitempty"`
	/*
	  Lower buffer capacity of number of write IO per second.
	*/
	DiskBaseWriteIops []import2.TimeIntValuePair `json:"diskBaseWriteIops,omitempty"`
	/*
	  Total amount of storage capacity of a device in bytes.
	*/
	DiskCapacityBytes []import2.TimeIntValuePair `json:"diskCapacityBytes,omitempty"`
	/*
	  Free storage space available on the Disk, measured in bytes.
	*/
	DiskFreeBytes []import2.TimeIntValuePair `json:"diskFreeBytes,omitempty"`
	/*
	  IO bandwidth - KB per second.
	*/
	DiskIoBandwidthkbps []import2.TimeIntValuePair `json:"diskIoBandwidthkbps,omitempty"`
	/*
	  Number of IO operations that a Disk perform per second.
	*/
	DiskNumIops []import2.TimeIntValuePair `json:"diskNumIops,omitempty"`
	/*
	  Upper limit of data transfer that a Disk can handle per second.
	*/
	DiskPeakIoBandwidthkbps []import2.TimeIntValuePair `json:"diskPeakIoBandwidthkbps,omitempty"`
	/*
	  Upper limit of the latency of I/O operations that Disk can handle without exceeding its standard latency level.
	*/
	DiskPeakIoLatencyMicrosec []import2.TimeIntValuePair `json:"diskPeakIoLatencyMicrosec,omitempty"`
	/*
	  Upper limit of I/O operations that a Disk performs per second.
	*/
	DiskPeakNumIops []import2.TimeIntValuePair `json:"diskPeakNumIops,omitempty"`
	/*
	  Upper buffer capacity average read IO latency, measured in microseconds (usecs).
	*/
	DiskPeakReadIoAvgLatencyMicrosec []import2.TimeIntValuePair `json:"diskPeakReadIoAvgLatencyMicrosec,omitempty"`
	/*
	  Upper buffer capacity for the amount of IO bandwidth that a Disk can handle read operations.
	*/
	DiskPeakReadIoBandwidthkbps []import2.TimeIntValuePair `json:"diskPeakReadIoBandwidthkbps,omitempty"`
	/*
	  Upper buffer capacity for the number of read IOPS that a Disk can handle.
	*/
	DiskPeakReadIops []import2.TimeIntValuePair `json:"diskPeakReadIops,omitempty"`
	/*
	  Upper buffer capacity average write IO latency, measured in microseconds (usecs).
	*/
	DiskPeakWriteIoAvgLatencyMicrosec []import2.TimeIntValuePair `json:"diskPeakWriteIoAvgLatencyMicrosec,omitempty"`
	/*
	  Upper buffer capacity for the amount of IO bandwidth that a Disk can handle write operations.
	*/
	DiskPeakWriteIoBandwidthkbps []import2.TimeIntValuePair `json:"diskPeakWriteIoBandwidthkbps,omitempty"`
	/*
	  Upper buffer capacity of number of write IO per second.
	*/
	DiskPeakWriteIops []import2.TimeIntValuePair `json:"diskPeakWriteIops,omitempty"`
	/*
	  Average read IO latency, measured in microseconds (usecs).
	*/
	DiskReadIoAvgLatencyMicrosec []import2.TimeIntValuePair `json:"diskReadIoAvgLatencyMicrosec,omitempty"`
	/*
	  Number of Disk read IO per second reported by Stargate.
	*/
	DiskReadIoBandwidthkbps []import2.TimeIntValuePair `json:"diskReadIoBandwidthkbps,omitempty"`
	/*
	  Disk read IO, expressed in parts per million.
	*/
	DiskReadIoPpm []import2.TimeIntValuePair `json:"diskReadIoPpm,omitempty"`
	/*
	  Number of read IO per second.
	*/
	DiskReadIops []import2.TimeIntValuePair `json:"diskReadIops,omitempty"`
	/*
	  Amount of storage currently being used, measured in bytes.
	*/
	DiskUsageBytes []import2.TimeIntValuePair `json:"diskUsageBytes,omitempty"`
	/*
	  Disk space used on a storage device, expressed in parts per million (ppm).
	*/
	DiskUsagePpm []import2.TimeIntValuePair `json:"diskUsagePpm,omitempty"`
	/*
	  Average write IO latency, measured in microseconds (usecs).
	*/
	DiskWriteIoAvgLatencyMicrosec []import2.TimeIntValuePair `json:"diskWriteIoAvgLatencyMicrosec,omitempty"`
	/*
	  Number of Disk write IO per second reported by Stargate.
	*/
	DiskWriteIoBandwidthkbps []import2.TimeIntValuePair `json:"diskWriteIoBandwidthkbps,omitempty"`
	/*
	  Disk write IO, expressed in parts per million.
	*/
	DiskWriteIoPpm []import2.TimeIntValuePair `json:"diskWriteIoPpm,omitempty"`
	/*
	  Number of write IO per second.
	*/
	DiskWriteIops []import2.TimeIntValuePair `json:"diskWriteIops,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewDiskStats() *DiskStats {
	p := new(DiskStats)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.DiskStats"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/stats/disks/{extId} Get operation
*/
type GetDiskStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetDiskStatsApiResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetDiskStatsApiResponse() *GetDiskStatsApiResponse {
	p := new(GetDiskStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.GetDiskStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetDiskStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetDiskStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetDiskStatsApiResponseData()
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
REST response for all response codes in API path /clustermgmt/v4.0.b2/stats/storage-containers/{extId} Get operation
*/
type GetStorageContainerStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetStorageContainerStatsApiResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetStorageContainerStatsApiResponse() *GetStorageContainerStatsApiResponse {
	p := new(GetStorageContainerStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.GetStorageContainerStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetStorageContainerStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetStorageContainerStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetStorageContainerStatsApiResponseData()
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

type HostStats struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpm []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpm,omitempty"`
	/*
	  Lower Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmLowerBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmUpperBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmUpperBuf,omitempty"`
	/*
	  Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecs []TimeValuePair `json:"controllerAvgIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecs []TimeValuePair `json:"controllerAvgReadIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecs []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller IOPS Number.
	*/
	ControllerNumIops []TimeValuePair `json:"controllerNumIops,omitempty"`
	/*
	  Lower Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsLowerBuf []TimeValuePair `json:"controllerNumIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsUpperBuf []TimeValuePair `json:"controllerNumIopsUpperBuf,omitempty"`
	/*
	  Number of controller read IOPS.
	*/
	ControllerNumReadIops []TimeValuePair `json:"controllerNumReadIops,omitempty"`
	/*
	  Lower Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsLowerBuf []TimeValuePair `json:"controllerNumReadIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsUpperBuf []TimeValuePair `json:"controllerNumReadIopsUpperBuf,omitempty"`
	/*
	  Number of controller write IOPS.
	*/
	ControllerNumWriteIops []TimeValuePair `json:"controllerNumWriteIops,omitempty"`
	/*
	  Lower Buf value of Number of controller write IoPS.
	*/
	ControllerNumWriteIopsLowerBuf []TimeValuePair `json:"controllerNumWriteIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller write IOPS.
	*/
	ControllerNumWriteIopsUpperBuf []TimeValuePair `json:"controllerNumWriteIopsUpperBuf,omitempty"`
	/*
	  Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbps []TimeValuePair `json:"controllerReadIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbps []TimeValuePair `json:"controllerWriteIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  CPU capacity in Hz.
	*/
	CpuCapacityHz []TimeValuePair `json:"cpuCapacityHz,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Free physical space(bytes).
	*/
	FreePhysicalStorageBytes []TimeValuePair `json:"freePhysicalStorageBytes,omitempty"`
	/*
	  Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpm []TimeValuePair `json:"hypervisorCpuUsagePpm,omitempty"`
	/*
	  Lower Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmLowerBuf []TimeValuePair `json:"hypervisorCpuUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmUpperBuf []TimeValuePair `json:"hypervisorCpuUsagePpmUpperBuf,omitempty"`
	/*
	  Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbps []TimeValuePair `json:"ioBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsLowerBuf []TimeValuePair `json:"ioBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsUpperBuf []TimeValuePair `json:"ioBandwidthKbpsUpperBuf,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  Size of memory(in bytes).
	*/
	MemoryCapacityBytes []TimeValuePair `json:"memoryCapacityBytes,omitempty"`
	/*
	  Overall memory usage(ppm).
	*/
	OverallMemoryUsagePpm []TimeValuePair `json:"overallMemoryUsagePpm,omitempty"`
	/*
	  Lower Buf value of overall memory usage(ppm).
	*/
	OverallMemoryUsagePpmLowerBuf []TimeValuePair `json:"overallMemoryUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of overall memory usage(ppm).
	*/
	OverallMemoryUsagePpmUpperBuf []TimeValuePair `json:"overallMemoryUsagePpmUpperBuf,omitempty"`

	StatType *import2.DownSamplingOperator `json:"statType,omitempty"`
	/*
	  Storage capacity(bytes).
	*/
	StorageCapacityBytes []TimeValuePair `json:"storageCapacityBytes,omitempty"`
	/*
	  Storage usage(bytes).
	*/
	StorageUsageBytes []TimeValuePair `json:"storageUsageBytes,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewHostStats() *HostStats {
	p := new(HostStats)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.HostStats"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
REST response for all response codes in API path /clustermgmt/v4.0.b2/stats/clusters/{clusterExtId}/hosts/{hostExtId} Get operation
*/
type HostStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfHostStatsApiResponseData `json:"data,omitempty"`

	Metadata *import1.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewHostStatsApiResponse() *HostStatsApiResponse {
	p := new(HostStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.HostStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *HostStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *HostStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfHostStatsApiResponseData()
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

type HostStatsProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpm []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpm,omitempty"`
	/*
	  Lower Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmLowerBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Aggregate Hypervisor Memory Usage(ppm).
	*/
	AggregateHypervisorMemoryUsagePpmUpperBuf []TimeValuePair `json:"aggregateHypervisorMemoryUsagePpmUpperBuf,omitempty"`
	/*
	  Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecs []TimeValuePair `json:"controllerAvgIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average IO Latency(usecs).
	*/
	ControllerAvgIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecs []TimeValuePair `json:"controllerAvgReadIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Read IO Latency(usecs).
	*/
	ControllerAvgReadIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgReadIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecs []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecs,omitempty"`
	/*
	  Lower Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsLowerBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Average Write IO Latency(usecs).
	*/
	ControllerAvgWriteIoLatencyUsecsUpperBuf []TimeValuePair `json:"controllerAvgWriteIoLatencyUsecsUpperBuf,omitempty"`
	/*
	  Controller IOPS Number.
	*/
	ControllerNumIops []TimeValuePair `json:"controllerNumIops,omitempty"`
	/*
	  Lower Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsLowerBuf []TimeValuePair `json:"controllerNumIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IOPS Number.
	*/
	ControllerNumIopsUpperBuf []TimeValuePair `json:"controllerNumIopsUpperBuf,omitempty"`
	/*
	  Number of controller read IOPS.
	*/
	ControllerNumReadIops []TimeValuePair `json:"controllerNumReadIops,omitempty"`
	/*
	  Lower Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsLowerBuf []TimeValuePair `json:"controllerNumReadIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller read IOPS.
	*/
	ControllerNumReadIopsUpperBuf []TimeValuePair `json:"controllerNumReadIopsUpperBuf,omitempty"`
	/*
	  Number of controller write IOPS.
	*/
	ControllerNumWriteIops []TimeValuePair `json:"controllerNumWriteIops,omitempty"`
	/*
	  Lower Buf value of Number of controller write IoPS.
	*/
	ControllerNumWriteIopsLowerBuf []TimeValuePair `json:"controllerNumWriteIopsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Number of controller write IOPS.
	*/
	ControllerNumWriteIopsUpperBuf []TimeValuePair `json:"controllerNumWriteIopsUpperBuf,omitempty"`
	/*
	  Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbps []TimeValuePair `json:"controllerReadIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Read IO Bandwidth(kBps).
	*/
	ControllerReadIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerReadIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbps []TimeValuePair `json:"controllerWriteIoBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsLowerBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller Write IO Bandwidth(kBps).
	*/
	ControllerWriteIoBandwidthKbpsUpperBuf []TimeValuePair `json:"controllerWriteIoBandwidthKbpsUpperBuf,omitempty"`
	/*
	  CPU capacity in Hz.
	*/
	CpuCapacityHz []TimeValuePair `json:"cpuCapacityHz,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Free physical space(bytes).
	*/
	FreePhysicalStorageBytes []TimeValuePair `json:"freePhysicalStorageBytes,omitempty"`
	/*
	  Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpm []TimeValuePair `json:"hypervisorCpuUsagePpm,omitempty"`
	/*
	  Lower Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmLowerBuf []TimeValuePair `json:"hypervisorCpuUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of Hypervisor CPU Usage(ppm).
	*/
	HypervisorCpuUsagePpmUpperBuf []TimeValuePair `json:"hypervisorCpuUsagePpmUpperBuf,omitempty"`
	/*
	  Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbps []TimeValuePair `json:"ioBandwidthKbps,omitempty"`
	/*
	  Lower Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsLowerBuf []TimeValuePair `json:"ioBandwidthKbpsLowerBuf,omitempty"`
	/*
	  Upper Buf value of Controller IO Bandwidth(kBps).
	*/
	IoBandwidthKbpsUpperBuf []TimeValuePair `json:"ioBandwidthKbpsUpperBuf,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  Size of memory(in bytes).
	*/
	MemoryCapacityBytes []TimeValuePair `json:"memoryCapacityBytes,omitempty"`
	/*
	  Overall memory usage(ppm).
	*/
	OverallMemoryUsagePpm []TimeValuePair `json:"overallMemoryUsagePpm,omitempty"`
	/*
	  Lower Buf value of overall memory usage(ppm).
	*/
	OverallMemoryUsagePpmLowerBuf []TimeValuePair `json:"overallMemoryUsagePpmLowerBuf,omitempty"`
	/*
	  Upper Buf value of overall memory usage(ppm).
	*/
	OverallMemoryUsagePpmUpperBuf []TimeValuePair `json:"overallMemoryUsagePpmUpperBuf,omitempty"`

	StatType *import2.DownSamplingOperator `json:"statType,omitempty"`
	/*
	  Storage capacity(bytes).
	*/
	StorageCapacityBytes []TimeValuePair `json:"storageCapacityBytes,omitempty"`
	/*
	  Storage usage(bytes).
	*/
	StorageUsageBytes []TimeValuePair `json:"storageUsageBytes,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewHostStatsProjection() *HostStatsProjection {
	p := new(HostStatsProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.HostStatsProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type StorageContainerStats struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  extId of the Storage Container.
	*/
	ContainerExtId *string `json:"containerExtId,omitempty"`
	/*
	  Average I/O latency in micro secs.
	*/
	ControllerAvgIoLatencyuSecs []import2.TimeIntValuePair `json:"controllerAvgIoLatencyuSecs,omitempty"`
	/*
	  Average read I/O latency in microseconds.
	*/
	ControllerAvgReadIoLatencyuSecs []import2.TimeIntValuePair `json:"controllerAvgReadIoLatencyuSecs,omitempty"`
	/*
	  Average read I/O latency in microseconds.
	*/
	ControllerAvgWriteIoLatencyuSecs []import2.TimeIntValuePair `json:"controllerAvgWriteIoLatencyuSecs,omitempty"`
	/*
	  Total I/O bandwidth - kB per second.
	*/
	ControllerIoBandwidthkBps []import2.TimeIntValuePair `json:"controllerIoBandwidthkBps,omitempty"`
	/*
	  Number of I/O per second.
	*/
	ControllerNumIops []import2.TimeIntValuePair `json:"controllerNumIops,omitempty"`
	/*
	  Number of read I/O per second.
	*/
	ControllerNumReadIops []import2.TimeIntValuePair `json:"controllerNumReadIops,omitempty"`
	/*
	  Number of write I/O per second.
	*/
	ControllerNumWriteIops []import2.TimeIntValuePair `json:"controllerNumWriteIops,omitempty"`
	/*
	  Read I/O bandwidth - kB per second.
	*/
	ControllerReadIoBandwidthkBps []import2.TimeIntValuePair `json:"controllerReadIoBandwidthkBps,omitempty"`
	/*
	  Ratio of read I/O to total I/O in PPM.
	*/
	ControllerReadIoRatioPpm []import2.TimeIntValuePair `json:"controllerReadIoRatioPpm,omitempty"`
	/*
	  Write I/O bandwidth - kB per second.
	*/
	ControllerWriteIoBandwidthkBps []import2.TimeIntValuePair `json:"controllerWriteIoBandwidthkBps,omitempty"`
	/*
	  Ratio of read I/O to total I/O in PPM.
	*/
	ControllerWriteIoRatioPpm []import2.TimeIntValuePair `json:"controllerWriteIoRatioPpm,omitempty"`
	/*
	  Saving ratio in PPM as a result of the Cloning technique.
	*/
	DataReductionCloneSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionCloneSavingRatioPpm,omitempty"`
	/*
	  Saving ratio in PPM as a result of the Compression technique.
	*/
	DataReductionCompressionSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionCompressionSavingRatioPpm,omitempty"`
	/*
	  Saving ratio in PPM as a result of the Deduplication technique.
	*/
	DataReductionDedupSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionDedupSavingRatioPpm,omitempty"`
	/*
	  Saving ratio in PPM as a result of the Erasure Coding technique.
	*/
	DataReductionErasureCodingSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionErasureCodingSavingRatioPpm,omitempty"`
	/*
	  Usage in bytes after reduction of Deduplication, Compression, Erasure Coding, Cloning, and Thin provisioning.
	*/
	DataReductionOverallPostReductionBytes []import2.TimeIntValuePair `json:"dataReductionOverallPostReductionBytes,omitempty"`
	/*
	  Usage in bytes before reduction of Deduplication, Compression, Erasure Coding, Cloning, and Thin provisioning.
	*/
	DataReductionOverallPreReductionBytes []import2.TimeIntValuePair `json:"dataReductionOverallPreReductionBytes,omitempty"`
	/*
	  Storage savings in bytes as a result of all the techniques.
	*/
	DataReductionSavedBytes []import2.TimeIntValuePair `json:"dataReductionSavedBytes,omitempty"`
	/*
	  Saving ratio in PPM as a result of Deduplication, compression and Erasure Coding.
	*/
	DataReductionSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionSavingRatioPpm,omitempty"`
	/*
	  Saving ratio in PPM as a result of Snapshot technique.
	*/
	DataReductionSnapshotSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionSnapshotSavingRatioPpm,omitempty"`
	/*
	  Saving ratio in PPM as a result of the Thin Provisioning technique.
	*/
	DataReductionThinProvisionSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionThinProvisionSavingRatioPpm,omitempty"`
	/*
	  Saving ratio in PPM consisting of Deduplication, Compression, Erasure Coding, Cloning, and Thin Provisioning.
	*/
	DataReductionTotalSavingRatioPpm []import2.TimeIntValuePair `json:"dataReductionTotalSavingRatioPpm,omitempty"`
	/*
	  Total amount of savings in bytes as a result of zero writes.
	*/
	DataReductionZeroWriteSavingsBytes []import2.TimeIntValuePair `json:"dataReductionZeroWriteSavingsBytes,omitempty"`
	/*
	  A globally unique identifier of an instance that is suitable for external consumption.
	*/
	ExtId *string `json:"extId,omitempty"`
	/*
	  Health of the container is represented by an integer value in the range 0-100. Higher value is indicative of better health.
	*/
	Health []import2.TimeIntValuePair `json:"health,omitempty"`
	/*
	  A HATEOAS style link for the response.  Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
	*/
	Links []import1.ApiLink `json:"links,omitempty"`
	/*
	  Actual physical disk usage of the container without accounting for the reservation.
	*/
	StorageActualPhysicalUsageBytes []import2.TimeIntValuePair `json:"storageActualPhysicalUsageBytes,omitempty"`
	/*
	  Storage capacity in bytes.
	*/
	StorageCapacityBytes []import2.TimeIntValuePair `json:"storageCapacityBytes,omitempty"`
	/*
	  Free storage in bytes.
	*/
	StorageFreeBytes []import2.TimeIntValuePair `json:"storageFreeBytes,omitempty"`
	/*
	  Replication factor of Container.
	*/
	StorageReplicationFactor []import2.TimeIntValuePair `json:"storageReplicationFactor,omitempty"`
	/*
	  Implicit physical reserved capacity(aggregated on vDisk level due to thick provisioning) in bytes.
	*/
	StorageReservedCapacityBytes []import2.TimeIntValuePair `json:"storageReservedCapacityBytes,omitempty"`
	/*
	  Total usage on HDD tier for the Container in bytes.
	*/
	StorageTierDasSataUsageBytes []import2.TimeIntValuePair `json:"storageTierDasSataUsageBytes,omitempty"`
	/*
	  Total usage on SDD tier for the Container in bytes
	*/
	StorageTierSsdUsageBytes []import2.TimeIntValuePair `json:"storageTierSsdUsageBytes,omitempty"`
	/*
	  Used storage in bytes.
	*/
	StorageUsageBytes []import2.TimeIntValuePair `json:"storageUsageBytes,omitempty"`
	/*
	  A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
	*/
	TenantId *string `json:"tenantId,omitempty"`
}

func NewStorageContainerStats() *StorageContainerStats {
	p := new(StorageContainerStats)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.StorageContainerStats"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
Time - Value pair for time-series stat attributes.
*/
type TimeValuePair struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Timestamp for given stat attribute(in ISO-8601 format).
	*/
	Timestamp *time.Time `json:"timestamp,omitempty"`
	/*
	  Value of stat at given timestamp.
	*/
	Value *int64 `json:"value,omitempty"`
}

func NewTimeValuePair() *TimeValuePair {
	p := new(TimeValuePair)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "clustermgmt.v4.stats.TimeValuePair"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b2"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type OneOfGetDiskStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *DiskStats             `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
}

func NewOneOfGetDiskStatsApiResponseData() *OneOfGetDiskStatsApiResponseData {
	p := new(OneOfGetDiskStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetDiskStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetDiskStatsApiResponseData is nil"))
	}
	switch v.(type) {
	case DiskStats:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(DiskStats)
		}
		*p.oneOfType2001 = v.(DiskStats)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
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
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetDiskStatsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetDiskStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(DiskStats)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.stats.DiskStats" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(DiskStats)
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
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetDiskStatsApiResponseData"))
}

func (p *OneOfGetDiskStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetDiskStatsApiResponseData")
}

type OneOfClusterStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *ClusterStats          `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
}

func NewOneOfClusterStatsApiResponseData() *OneOfClusterStatsApiResponseData {
	p := new(OneOfClusterStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfClusterStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfClusterStatsApiResponseData is nil"))
	}
	switch v.(type) {
	case ClusterStats:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(ClusterStats)
		}
		*p.oneOfType2001 = v.(ClusterStats)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
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
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfClusterStatsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfClusterStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(ClusterStats)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.stats.ClusterStats" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(ClusterStats)
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
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfClusterStatsApiResponseData"))
}

func (p *OneOfClusterStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfClusterStatsApiResponseData")
}

type OneOfHostStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *HostStats             `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
}

func NewOneOfHostStatsApiResponseData() *OneOfHostStatsApiResponseData {
	p := new(OneOfHostStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfHostStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfHostStatsApiResponseData is nil"))
	}
	switch v.(type) {
	case HostStats:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(HostStats)
		}
		*p.oneOfType2001 = v.(HostStats)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
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
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfHostStatsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfHostStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(HostStats)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.stats.HostStats" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(HostStats)
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
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfHostStatsApiResponseData"))
}

func (p *OneOfHostStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfHostStatsApiResponseData")
}

type OneOfGetStorageContainerStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *StorageContainerStats `json:"-"`
	oneOfType400  *import3.ErrorResponse `json:"-"`
}

func NewOneOfGetStorageContainerStatsApiResponseData() *OneOfGetStorageContainerStatsApiResponseData {
	p := new(OneOfGetStorageContainerStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetStorageContainerStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetStorageContainerStatsApiResponseData is nil"))
	}
	switch v.(type) {
	case StorageContainerStats:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(StorageContainerStats)
		}
		*p.oneOfType2001 = v.(StorageContainerStats)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
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
	default:
		return errors.New(fmt.Sprintf("%T(%v) is not expected type", v, v))
	}
	return nil
}

func (p *OneOfGetStorageContainerStatsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetStorageContainerStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(StorageContainerStats)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "clustermgmt.v4.stats.StorageContainerStats" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(StorageContainerStats)
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
	vOneOfType400 := new(import3.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "clustermgmt.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetStorageContainerStatsApiResponseData"))
}

func (p *OneOfGetStorageContainerStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetStorageContainerStatsApiResponseData")
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
