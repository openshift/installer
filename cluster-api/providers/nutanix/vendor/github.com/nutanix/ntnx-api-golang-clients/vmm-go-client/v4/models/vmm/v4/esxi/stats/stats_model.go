/*
 * Generated file models/vmm/v4/esxi/stats/stats_model.go.
 *
 * Product version: 4.0.1-beta-1
 *
 * Part of the Nutanix VMM APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module vmm.v4.esxi.stats of Nutanix VMM APIs
*/
package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/error"
	"time"
)

/*
REST response for all response codes in API path /vmm/v4.0.b1/esxi/stats/vms/{extId} Get operation
*/
type GetVmStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetVmStatsApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetVmStatsApiResponse() *GetVmStatsApiResponse {
	p := new(GetVmStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.esxi.stats.GetVmStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetVmStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetVmStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetVmStatsApiResponseData()
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
REST response for all response codes in API path /vmm/v4.0.b1/esxi/stats/vms Get operation
*/
type ListVmStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfListVmStatsApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewListVmStatsApiResponse() *ListVmStatsApiResponse {
	p := new(ListVmStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.esxi.stats.ListVmStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *ListVmStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *ListVmStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfListVmStatsApiResponseData()
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
A collection of VM stats.
*/
type VmStats struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The timestamp of a specific VM stats response data point.
	*/
	Stats []VmStatsTuple `json:"stats,omitempty"`
	/*
	  The VM external ID associated with the VM stats.
	*/
	VmExtId *string `json:"vmExtId,omitempty"`
}

func NewVmStats() *VmStats {
	p := new(VmStats)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.esxi.stats.VmStats"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
A collection of VM stats.
*/
type VmStatsTuple struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  The VM NCC health check score.
	*/
	CheckScore *int64 `json:"checkScore,omitempty"`
	/*
	  The UUID of the cluster on which the VM resides.
	*/
	Cluster *string `json:"cluster,omitempty"`
	/*
	  The VM controller average I/O latency in microseconds.
	*/
	ControllerAvgIoLatencyMicros *int64 `json:"controllerAvgIoLatencyMicros,omitempty"`
	/*
	  The VM controller average read I/O latency in microseconds.
	*/
	ControllerAvgReadIoLatencyMicros *int64 `json:"controllerAvgReadIoLatencyMicros,omitempty"`
	/*
	  The VM controller average read I/O size in kilobytes.
	*/
	ControllerAvgReadIoSizeKb *int64 `json:"controllerAvgReadIoSizeKb,omitempty"`
	/*
	  The VM controller average write I/O latency in microseconds.
	*/
	ControllerAvgWriteIoLatencyMicros *int64 `json:"controllerAvgWriteIoLatencyMicros,omitempty"`
	/*
	  The VM controller average write I/O size in kilobytes.
	*/
	ControllerAvgWriteIoSizeKb *int64 `json:"controllerAvgWriteIoSizeKb,omitempty"`
	/*
	  The VM controller I/O bandwidth in kilobytes per second.
	*/
	ControllerIoBandwidthKbps *int64 `json:"controllerIoBandwidthKbps,omitempty"`
	/*
	  The VM controller number of I/O requests.
	*/
	ControllerNumIo *int64 `json:"controllerNumIo,omitempty"`
	/*
	  The VM controller number of I/O operations per second.
	*/
	ControllerNumIops *int64 `json:"controllerNumIops,omitempty"`
	/*
	  The VM controller number of random I/O.
	*/
	ControllerNumRandomIo *int64 `json:"controllerNumRandomIo,omitempty"`
	/*
	  The VM controller number of read I/O.
	*/
	ControllerNumReadIo *int64 `json:"controllerNumReadIo,omitempty"`
	/*
	  The VM controller number of read I/O operations per second.
	*/
	ControllerNumReadIops *int64 `json:"controllerNumReadIops,omitempty"`
	/*
	  The VM controller number of sequential I/Os.
	*/
	ControllerNumSeqIo *int64 `json:"controllerNumSeqIo,omitempty"`
	/*
	  The VM controller number of write I/O.
	*/
	ControllerNumWriteIo *int64 `json:"controllerNumWriteIo,omitempty"`
	/*
	  The VM controller number of write I/O operations per second.
	*/
	ControllerNumWriteIops *int64 `json:"controllerNumWriteIops,omitempty"`
	/*
	  The VM controller number of random I/O PPM.
	*/
	ControllerRandomIoPpm *int64 `json:"controllerRandomIoPpm,omitempty"`
	/*
	  The VM controller number of read I/O bandwidth in kilobytes per second.
	*/
	ControllerReadIoBandwidthKbps *int64 `json:"controllerReadIoBandwidthKbps,omitempty"`
	/*
	  The VM controller number of read I/O PPM.
	*/
	ControllerReadIoPpm *int64 `json:"controllerReadIoPpm,omitempty"`
	/*
	  The VM controller number of sequential I/O PPM.
	*/
	ControllerSeqIoPpm *int64 `json:"controllerSeqIoPpm,omitempty"`
	/*
	  The VM controller total usage on SSD tier for the VM.
	*/
	ControllerStorageTierSsdUsageBytes *int64 `json:"controllerStorageTierSsdUsageBytes,omitempty"`
	/*
	  The VM controller timespan in microseconds.
	*/
	ControllerTimespanMicros *int64 `json:"controllerTimespanMicros,omitempty"`
	/*
	  The VM controller number of total I/O size in kilobytes.
	*/
	ControllerTotalIoSizeKb *int64 `json:"controllerTotalIoSizeKb,omitempty"`
	/*
	  The VM controller number of total I/O time in microseconds.
	*/
	ControllerTotalIoTimeMicros *int64 `json:"controllerTotalIoTimeMicros,omitempty"`
	/*
	  The VM controller number of total read I/O size in kilobytes.
	*/
	ControllerTotalReadIoSizeKb *int64 `json:"controllerTotalReadIoSizeKb,omitempty"`
	/*
	  The VM controller number of total read I/O time in microseconds.
	*/
	ControllerTotalReadIoTimeMicros *int64 `json:"controllerTotalReadIoTimeMicros,omitempty"`
	/*
	  The VM controller number of total transformed usage in bytes.
	*/
	ControllerTotalTransformedUsageBytes *int64 `json:"controllerTotalTransformedUsageBytes,omitempty"`
	/*
	  The VM controller user bytes.
	*/
	ControllerUserBytes *int64 `json:"controllerUserBytes,omitempty"`
	/*
	  The VM controller write I/O bandwidth in kilobytes per second.
	*/
	ControllerWriteIoBandwidthKbps *int64 `json:"controllerWriteIoBandwidthKbps,omitempty"`
	/*
	  The VM controller percentage of write I/O in parts per million.
	*/
	ControllerWriteIoPpm *int64 `json:"controllerWriteIoPpm,omitempty"`
	/*
	  The average I/O latency of the VM in microseconds
	*/
	HypervisorAvgIoLatencyMicros *int64 `json:"hypervisorAvgIoLatencyMicros,omitempty"`
	/*
	  Percentage of time that the VM was ready, but could not get scheduled to run.
	*/
	HypervisorCpuReadyTimePpm *int64 `json:"hypervisorCpuReadyTimePpm,omitempty"`
	/*
	  The CPU usage of the VM in parts per million.
	*/
	HypervisorCpuUsagePpm *int64 `json:"hypervisorCpuUsagePpm,omitempty"`
	/*
	  The I/O bandwidth of the VM in kilobytes per second.
	*/
	HypervisorIoBandwidthKbps *int64 `json:"hypervisorIoBandwidthKbps,omitempty"`
	/*
	  Consolidated guest memory usage in percentage.
	*/
	HypervisorMemoryUsagePpm *int64 `json:"hypervisorMemoryUsagePpm,omitempty"`
	/*
	  The number of I/O by the VM.
	*/
	HypervisorNumIo *int64 `json:"hypervisorNumIo,omitempty"`
	/*
	  The number of I/O operations by the VM per second.
	*/
	HypervisorNumIops *int64 `json:"hypervisorNumIops,omitempty"`
	/*
	  The number of read I/O operations by the VM.
	*/
	HypervisorNumReadIo *int64 `json:"hypervisorNumReadIo,omitempty"`
	/*
	  The number of read I/O operations by the VM per second.
	*/
	HypervisorNumReadIops *int64 `json:"hypervisorNumReadIops,omitempty"`
	/*
	  The number of bytes received by the VM.
	*/
	HypervisorNumReceivedBytes *int64 `json:"hypervisorNumReceivedBytes,omitempty"`
	/*
	  The number of bytes transmitted by the VM.
	*/
	HypervisorNumTransmittedBytes *int64 `json:"hypervisorNumTransmittedBytes,omitempty"`
	/*
	  The number of write I/O by the VM.
	*/
	HypervisorNumWriteIo *int64 `json:"hypervisorNumWriteIo,omitempty"`
	/*
	  The number of write I/O operations by the VM per second.
	*/
	HypervisorNumWriteIops *int64 `json:"hypervisorNumWriteIops,omitempty"`
	/*
	  The number of read I/O bandwidth of the VM in kilobytes per second.
	*/
	HypervisorReadIoBandwidthKbps *int64 `json:"hypervisorReadIoBandwidthKbps,omitempty"`
	/*
	  The timespan of the VM in microseconds.
	*/
	HypervisorTimespanMicros *int64 `json:"hypervisorTimespanMicros,omitempty"`
	/*
	  The total I/O size of the VM in kilobytes.
	*/
	HypervisorTotalIoSizeKb *int64 `json:"hypervisorTotalIoSizeKb,omitempty"`
	/*
	  The total I/O time of the VM in microseconds.
	*/
	HypervisorTotalIoTimeMicros *int64 `json:"hypervisorTotalIoTimeMicros,omitempty"`
	/*
	  The total read I/O size of the VM in kilobytes.
	*/
	HypervisorTotalReadIoSizeKb *int64 `json:"hypervisorTotalReadIoSizeKb,omitempty"`
	/*
	  Hypervisor type of the VM.
	*/
	HypervisorType *string `json:"hypervisorType,omitempty"`
	/*
	  The write I/O bandwidth of the VM in kilobytes per second.
	*/
	HypervisorWriteIoBandwidthKbps *int64 `json:"hypervisorWriteIoBandwidthKbps,omitempty"`
	/*
	  The VM memory usage in PPM.
	*/
	MemoryUsagePpm *int64 `json:"memoryUsagePpm,omitempty"`
	/*
	  The timestamp of a specific VM stats response data point.
	*/
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

func NewVmStatsTuple() *VmStatsTuple {
	p := new(VmStatsTuple)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "vmm.v4.esxi.stats.VmStatsTuple"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type OneOfListVmStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 []VmStats              `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
}

func NewOneOfListVmStatsApiResponseData() *OneOfListVmStatsApiResponseData {
	p := new(OneOfListVmStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfListVmStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfListVmStatsApiResponseData is nil"))
	}
	switch v.(type) {
	case []VmStats:
		p.oneOfType2001 = v.([]VmStats)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = "List<vmm.v4.esxi.stats.VmStats>"
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = "List<vmm.v4.esxi.stats.VmStats>"
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
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

func (p *OneOfListVmStatsApiResponseData) GetValue() interface{} {
	if "List<vmm.v4.esxi.stats.VmStats>" == *p.Discriminator {
		return p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfListVmStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new([]VmStats)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {

		if len(*vOneOfType2001) == 0 || "vmm.v4.esxi.stats.VmStats" == *((*vOneOfType2001)[0].ObjectType_) {
			p.oneOfType2001 = *vOneOfType2001
			if nil == p.Discriminator {
				p.Discriminator = new(string)
			}
			*p.Discriminator = "List<vmm.v4.esxi.stats.VmStats>"
			if nil == p.ObjectType_ {
				p.ObjectType_ = new(string)
			}
			*p.ObjectType_ = "List<vmm.v4.esxi.stats.VmStats>"
			return nil

		}
	}
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfListVmStatsApiResponseData"))
}

func (p *OneOfListVmStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if "List<vmm.v4.esxi.stats.VmStats>" == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfListVmStatsApiResponseData")
}

type OneOfGetVmStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType2001 *VmStats               `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
}

func NewOneOfGetVmStatsApiResponseData() *OneOfGetVmStatsApiResponseData {
	p := new(OneOfGetVmStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetVmStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetVmStatsApiResponseData is nil"))
	}
	switch v.(type) {
	case VmStats:
		if nil == p.oneOfType2001 {
			p.oneOfType2001 = new(VmStats)
		}
		*p.oneOfType2001 = v.(VmStats)
		if nil == p.Discriminator {
			p.Discriminator = new(string)
		}
		*p.Discriminator = *p.oneOfType2001.ObjectType_
		if nil == p.ObjectType_ {
			p.ObjectType_ = new(string)
		}
		*p.ObjectType_ = *p.oneOfType2001.ObjectType_
	case import1.ErrorResponse:
		if nil == p.oneOfType400 {
			p.oneOfType400 = new(import1.ErrorResponse)
		}
		*p.oneOfType400 = v.(import1.ErrorResponse)
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

func (p *OneOfGetVmStatsApiResponseData) GetValue() interface{} {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return *p.oneOfType2001
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	return nil
}

func (p *OneOfGetVmStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType2001 := new(VmStats)
	if err := json.Unmarshal(b, vOneOfType2001); err == nil {
		if "vmm.v4.esxi.stats.VmStats" == *vOneOfType2001.ObjectType_ {
			if nil == p.oneOfType2001 {
				p.oneOfType2001 = new(VmStats)
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
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "vmm.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
			if nil == p.oneOfType400 {
				p.oneOfType400 = new(import1.ErrorResponse)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetVmStatsApiResponseData"))
}

func (p *OneOfGetVmStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType2001 != nil && *p.oneOfType2001.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType2001)
	}
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	return nil, errors.New("No value to marshal for OneOfGetVmStatsApiResponseData")
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
