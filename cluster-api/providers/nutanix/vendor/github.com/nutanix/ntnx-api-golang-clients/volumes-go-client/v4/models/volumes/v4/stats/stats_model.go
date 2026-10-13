/*
 * Generated file models/volumes/v4/stats/stats_model.go.
 *
 * Product version: 4.0.1-beta-1
 *
 * Part of the Nutanix Volumes Versioned APIs
 *
 * (c) 2024 Nutanix Inc.  All rights reserved
 *
 */

/*
  Module volumes.v4.stats of Nutanix Volumes Versioned APIs
*/
package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	import2 "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/error"
	"time"
)

/*
REST response for all response codes in API path /volumes/v4.0.b1/stats/volume-groups/{volumeGroupExtId}/disks/{extId} Get operation
*/
type GetVolumeDiskStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetVolumeDiskStatsApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetVolumeDiskStatsApiResponse() *GetVolumeDiskStatsApiResponse {
	p := new(GetVolumeDiskStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "volumes.v4.stats.GetVolumeDiskStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetVolumeDiskStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetVolumeDiskStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetVolumeDiskStatsApiResponseData()
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
REST response for all response codes in API path /volumes/v4.0.b1/stats/volume-groups/{extId} Get operation
*/
type GetVolumeGroupStatsApiResponse struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*

	 */
	DataItemDiscriminator_ *string `json:"$dataItemDiscriminator,omitempty"`

	Data *OneOfGetVolumeGroupStatsApiResponseData `json:"data,omitempty"`

	Metadata *import2.ApiResponseMetadata `json:"metadata,omitempty"`
}

func NewGetVolumeGroupStatsApiResponse() *GetVolumeGroupStatsApiResponse {
	p := new(GetVolumeGroupStatsApiResponse)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "volumes.v4.stats.GetVolumeGroupStatsApiResponse"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

func (p *GetVolumeGroupStatsApiResponse) GetData() interface{} {
	if nil == p.Data {
		return nil
	}
	return p.Data.GetValue()
}

func (p *GetVolumeGroupStatsApiResponse) SetData(v interface{}) error {
	if nil == p.Data {
		p.Data = NewOneOfGetVolumeGroupStatsApiResponseData()
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

type TimeValuePair struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Timestamp is returned in Epoch format.
	*/
	Timestamp *time.Time `json:"timestamp,omitempty"`
	/*
	  Value of the stat at the corresponding timestamp value represented in Int64 format.
	*/
	Value *int64 `json:"value,omitempty"`
}

func NewTimeValuePair() *TimeValuePair {
	p := new(TimeValuePair)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "volumes.v4.stats.TimeValuePair"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
A model that stores the Volume disk stats.
*/
type VolumeDiskStats struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Controller average I/O latency measured in microseconds.
	*/
	ControllerAvgIOLatencyUsecs []TimeValuePair `json:"controllerAvgIOLatencyUsecs,omitempty"`
	/*
	  Controller average read I/O latency measured in microseconds.
	*/
	ControllerAvgReadIOLatencyUsecs []TimeValuePair `json:"controllerAvgReadIOLatencyUsecs,omitempty"`
	/*
	  Controller average write I/O latency measured in microseconds.
	*/
	ControllerAvgWriteIOLatencyUsecs []TimeValuePair `json:"controllerAvgWriteIOLatencyUsecs,omitempty"`
	/*
	  Controller I/O bandwidth measured in Kbps.
	*/
	ControllerIOBandwidthKBps []TimeValuePair `json:"controllerIOBandwidthKBps,omitempty"`
	/*
	  Controller I/O rate measured in iops.
	*/
	ControllerNumIOPS []TimeValuePair `json:"controllerNumIOPS,omitempty"`
	/*
	  Controller read I/O measured in iops.
	*/
	ControllerNumReadIOPS []TimeValuePair `json:"controllerNumReadIOPS,omitempty"`
	/*
	  Controller write I/O measured in iops.
	*/
	ControllerNumWriteIOPS []TimeValuePair `json:"controllerNumWriteIOPS,omitempty"`
	/*
	  Controller read I/O bandwidth measured in Kbps.
	*/
	ControllerReadIOBandwidthKBps []TimeValuePair `json:"controllerReadIOBandwidthKBps,omitempty"`
	/*
	  Controller user bytes.
	*/
	ControllerUserBytes []TimeValuePair `json:"controllerUserBytes,omitempty"`
	/*
	  Controller write I/O bandwidth measured in Kbps.
	*/
	ControllerWriteIOBandwidthKBps []TimeValuePair `json:"controllerWriteIOBandwidthKBps,omitempty"`
	/*
	  Uuid of the Volume Disk.
	*/
	VolumeDiskExtId *string `json:"volumeDiskExtId,omitempty"`
}

func NewVolumeDiskStats() *VolumeDiskStats {
	p := new(VolumeDiskStats)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "volumes.v4.stats.VolumeDiskStats"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

/*
A model that stores the Volume group stats.
*/
type VolumeGroupStats struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Controller average I/O latency measured in microseconds.
	*/
	ControllerAvgIOLatencyUsecs []TimeValuePair `json:"controllerAvgIOLatencyUsecs,omitempty"`
	/*
	  Controller average read I/O latency measured in microseconds.
	*/
	ControllerAvgReadIOLatencyUsecs []TimeValuePair `json:"controllerAvgReadIOLatencyUsecs,omitempty"`
	/*
	  Controller average write I/O latency measured in microseconds.
	*/
	ControllerAvgWriteIOLatencyUsecs []TimeValuePair `json:"controllerAvgWriteIOLatencyUsecs,omitempty"`
	/*
	  Controller I/O bandwidth measured in Kbps.
	*/
	ControllerIOBandwidthKBps []TimeValuePair `json:"controllerIOBandwidthKBps,omitempty"`
	/*
	  Controller I/O rate measured in iops.
	*/
	ControllerNumIOPS []TimeValuePair `json:"controllerNumIOPS,omitempty"`
	/*
	  Controller read I/O measured in iops.
	*/
	ControllerNumReadIOPS []TimeValuePair `json:"controllerNumReadIOPS,omitempty"`
	/*
	  Controller write I/O measured in iops.
	*/
	ControllerNumWriteIOPS []TimeValuePair `json:"controllerNumWriteIOPS,omitempty"`
	/*
	  Controller read I/O bandwidth measured in Kbps.
	*/
	ControllerReadIOBandwidthKBps []TimeValuePair `json:"controllerReadIOBandwidthKBps,omitempty"`
	/*
	  Controller user bytes.
	*/
	ControllerUserBytes []TimeValuePair `json:"controllerUserBytes,omitempty"`
	/*
	  Controller write I/O bandwidth measured in Kbps.
	*/
	ControllerWriteIOBandwidthKBps []TimeValuePair `json:"controllerWriteIOBandwidthKBps,omitempty"`
	/*
	  Uuid of the Volume Group.
	*/
	VolumeGroupExtId *string `json:"volumeGroupExtId,omitempty"`
}

func NewVolumeGroupStats() *VolumeGroupStats {
	p := new(VolumeGroupStats)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "volumes.v4.stats.VolumeGroupStats"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type VolumeGroupStatsProjection struct {
	ObjectType_ *string `json:"$objectType,omitempty"`

	Reserved_ map[string]interface{} `json:"$reserved,omitempty"`

	UnknownFields_ map[string]interface{} `json:"$unknownFields,omitempty"`
	/*
	  Controller average I/O latency measured in microseconds.
	*/
	ControllerAvgIOLatencyUsecs []TimeValuePair `json:"controllerAvgIOLatencyUsecs,omitempty"`
	/*
	  Controller average read I/O latency measured in microseconds.
	*/
	ControllerAvgReadIOLatencyUsecs []TimeValuePair `json:"controllerAvgReadIOLatencyUsecs,omitempty"`
	/*
	  Controller average write I/O latency measured in microseconds.
	*/
	ControllerAvgWriteIOLatencyUsecs []TimeValuePair `json:"controllerAvgWriteIOLatencyUsecs,omitempty"`
	/*
	  Controller I/O bandwidth measured in Kbps.
	*/
	ControllerIOBandwidthKBps []TimeValuePair `json:"controllerIOBandwidthKBps,omitempty"`
	/*
	  Controller I/O rate measured in iops.
	*/
	ControllerNumIOPS []TimeValuePair `json:"controllerNumIOPS,omitempty"`
	/*
	  Controller read I/O measured in iops.
	*/
	ControllerNumReadIOPS []TimeValuePair `json:"controllerNumReadIOPS,omitempty"`
	/*
	  Controller write I/O measured in iops.
	*/
	ControllerNumWriteIOPS []TimeValuePair `json:"controllerNumWriteIOPS,omitempty"`
	/*
	  Controller read I/O bandwidth measured in Kbps.
	*/
	ControllerReadIOBandwidthKBps []TimeValuePair `json:"controllerReadIOBandwidthKBps,omitempty"`
	/*
	  Controller user bytes.
	*/
	ControllerUserBytes []TimeValuePair `json:"controllerUserBytes,omitempty"`
	/*
	  Controller write I/O bandwidth measured in Kbps.
	*/
	ControllerWriteIOBandwidthKBps []TimeValuePair `json:"controllerWriteIOBandwidthKBps,omitempty"`
	/*
	  Uuid of the Volume Group.
	*/
	VolumeGroupExtId *string `json:"volumeGroupExtId,omitempty"`
}

func NewVolumeGroupStatsProjection() *VolumeGroupStatsProjection {
	p := new(VolumeGroupStatsProjection)
	p.ObjectType_ = new(string)
	*p.ObjectType_ = "volumes.v4.stats.VolumeGroupStatsProjection"
	p.Reserved_ = map[string]interface{}{"$fv": "v4.r0.b1"}
	p.UnknownFields_ = map[string]interface{}{}

	return p
}

type OneOfGetVolumeDiskStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
	oneOfType0    *VolumeDiskStats       `json:"-"`
}

func NewOneOfGetVolumeDiskStatsApiResponseData() *OneOfGetVolumeDiskStatsApiResponseData {
	p := new(OneOfGetVolumeDiskStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetVolumeDiskStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetVolumeDiskStatsApiResponseData is nil"))
	}
	switch v.(type) {
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
	case VolumeDiskStats:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(VolumeDiskStats)
		}
		*p.oneOfType0 = v.(VolumeDiskStats)
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

func (p *OneOfGetVolumeDiskStatsApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfGetVolumeDiskStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "volumes.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
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
	vOneOfType0 := new(VolumeDiskStats)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "volumes.v4.stats.VolumeDiskStats" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(VolumeDiskStats)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetVolumeDiskStatsApiResponseData"))
}

func (p *OneOfGetVolumeDiskStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfGetVolumeDiskStatsApiResponseData")
}

type OneOfGetVolumeGroupStatsApiResponseData struct {
	Discriminator *string                `json:"-"`
	ObjectType_   *string                `json:"-"`
	oneOfType400  *import1.ErrorResponse `json:"-"`
	oneOfType0    *VolumeGroupStats      `json:"-"`
}

func NewOneOfGetVolumeGroupStatsApiResponseData() *OneOfGetVolumeGroupStatsApiResponseData {
	p := new(OneOfGetVolumeGroupStatsApiResponseData)
	p.Discriminator = new(string)
	p.ObjectType_ = new(string)
	return p
}

func (p *OneOfGetVolumeGroupStatsApiResponseData) SetValue(v interface{}) error {
	if nil == p {
		return errors.New(fmt.Sprintf("OneOfGetVolumeGroupStatsApiResponseData is nil"))
	}
	switch v.(type) {
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
	case VolumeGroupStats:
		if nil == p.oneOfType0 {
			p.oneOfType0 = new(VolumeGroupStats)
		}
		*p.oneOfType0 = v.(VolumeGroupStats)
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

func (p *OneOfGetVolumeGroupStatsApiResponseData) GetValue() interface{} {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return *p.oneOfType400
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return *p.oneOfType0
	}
	return nil
}

func (p *OneOfGetVolumeGroupStatsApiResponseData) UnmarshalJSON(b []byte) error {
	vOneOfType400 := new(import1.ErrorResponse)
	if err := json.Unmarshal(b, vOneOfType400); err == nil {
		if "volumes.v4.error.ErrorResponse" == *vOneOfType400.ObjectType_ {
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
	vOneOfType0 := new(VolumeGroupStats)
	if err := json.Unmarshal(b, vOneOfType0); err == nil {
		if "volumes.v4.stats.VolumeGroupStats" == *vOneOfType0.ObjectType_ {
			if nil == p.oneOfType0 {
				p.oneOfType0 = new(VolumeGroupStats)
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
	return errors.New(fmt.Sprintf("Unable to unmarshal for OneOfGetVolumeGroupStatsApiResponseData"))
}

func (p *OneOfGetVolumeGroupStatsApiResponseData) MarshalJSON() ([]byte, error) {
	if p.oneOfType400 != nil && *p.oneOfType400.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType400)
	}
	if p.oneOfType0 != nil && *p.oneOfType0.ObjectType_ == *p.Discriminator {
		return json.Marshal(p.oneOfType0)
	}
	return nil, errors.New("No value to marshal for OneOfGetVolumeGroupStatsApiResponseData")
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
