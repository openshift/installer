// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20180813

import (
    "encoding/json"
    tcerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AddResourceTagRequest struct {
	*tchttp.BaseRequest

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// [ 资源六段式描述 ](https://cloud.tencent.com/document/product/598/10606)
	Resource *string `json:"Resource,omitempty" name:"Resource"`
}

func (r *AddResourceTagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AddResourceTagRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKey")
	delete(f, "TagValue")
	delete(f, "Resource")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AddResourceTagRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type AddResourceTagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *AddResourceTagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AddResourceTagResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type AttachResourcesTagRequest struct {
	*tchttp.BaseRequest

	// 资源所属业务名称（资源六段式中的第三段）
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源ID数组，资源个数最多为50
	ResourceIds []*string `json:"ResourceIds,omitempty" name:"ResourceIds"`

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// 资源所在地域，不区分地域的资源不需要传入该字段，区分地域的资源必填
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 资源前缀（资源六段式中最后一段"/"前面的部分），cos存储桶不需要传入该字段，其他云资源必填
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`
}

func (r *AttachResourcesTagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AttachResourcesTagRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ServiceType")
	delete(f, "ResourceIds")
	delete(f, "TagKey")
	delete(f, "TagValue")
	delete(f, "ResourceRegion")
	delete(f, "ResourcePrefix")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AttachResourcesTagRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type AttachResourcesTagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *AttachResourcesTagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AttachResourcesTagResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type CreateTagRequest struct {
	*tchttp.BaseRequest

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`
}

func (r *CreateTagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateTagRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKey")
	delete(f, "TagValue")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateTagRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type CreateTagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateTagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateTagResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DeleteResourceTagRequest struct {
	*tchttp.BaseRequest

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// [ 资源六段式描述 ](https://cloud.tencent.com/document/product/598/10606)
	Resource *string `json:"Resource,omitempty" name:"Resource"`
}

func (r *DeleteResourceTagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteResourceTagRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKey")
	delete(f, "Resource")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteResourceTagRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DeleteResourceTagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteResourceTagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteResourceTagResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DeleteTagRequest struct {
	*tchttp.BaseRequest

	// 需要删除的标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 需要删除的标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`
}

func (r *DeleteTagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteTagRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKey")
	delete(f, "TagValue")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteTagRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DeleteTagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteTagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteTagResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsByResourceIdsRequest struct {
	*tchttp.BaseRequest

	// 业务类型
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源前缀
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`

	// 资源ID数组，大小不超过50
	ResourceIds []*string `json:"ResourceIds,omitempty" name:"ResourceIds"`

	// 资源所在地域
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`
}

func (r *DescribeResourceTagsByResourceIdsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsByResourceIdsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ServiceType")
	delete(f, "ResourcePrefix")
	delete(f, "ResourceIds")
	delete(f, "ResourceRegion")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeResourceTagsByResourceIdsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsByResourceIdsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 标签列表
		Tags []*TagResource `json:"Tags,omitempty" name:"Tags"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeResourceTagsByResourceIdsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsByResourceIdsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsByResourceIdsSeqRequest struct {
	*tchttp.BaseRequest

	// 业务类型
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源前缀
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`

	// 资源唯一标记
	ResourceIds []*string `json:"ResourceIds,omitempty" name:"ResourceIds"`

	// 资源所在地域
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`
}

func (r *DescribeResourceTagsByResourceIdsSeqRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsByResourceIdsSeqRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ServiceType")
	delete(f, "ResourcePrefix")
	delete(f, "ResourceIds")
	delete(f, "ResourceRegion")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeResourceTagsByResourceIdsSeqRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsByResourceIdsSeqResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 标签列表
		Tags []*TagResource `json:"Tags,omitempty" name:"Tags"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeResourceTagsByResourceIdsSeqResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsByResourceIdsSeqResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsByTagKeysRequest struct {
	*tchttp.BaseRequest

	// 业务类型
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源前缀
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`

	// 资源地域
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 资源唯一标识
	ResourceIds []*string `json:"ResourceIds,omitempty" name:"ResourceIds"`

	// 资源标签键
	TagKeys []*string `json:"TagKeys,omitempty" name:"TagKeys"`

	// 每页大小，默认为 400
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`
}

func (r *DescribeResourceTagsByTagKeysRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsByTagKeysRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ServiceType")
	delete(f, "ResourcePrefix")
	delete(f, "ResourceRegion")
	delete(f, "ResourceIds")
	delete(f, "TagKeys")
	delete(f, "Limit")
	delete(f, "Offset")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeResourceTagsByTagKeysRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsByTagKeysResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 资源标签
		Rows []*ResourceIdTag `json:"Rows,omitempty" name:"Rows"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeResourceTagsByTagKeysResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsByTagKeysResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsRequest struct {
	*tchttp.BaseRequest

	// 创建者uin
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 资源所在地域
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 业务类型
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源前缀
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`

	// 资源唯一标识。只输入ResourceId进行查询可能会查询较慢，或者无法匹配到结果，建议在输入ResourceId的同时也输入ServiceType、ResourcePrefix和ResourceRegion（不区分地域的资源可忽略该参数）
	ResourceId *string `json:"ResourceId,omitempty" name:"ResourceId"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 是否是cos的资源（0或者1），输入的ResourceId为cos资源时必填
	CosResourceId *uint64 `json:"CosResourceId,omitempty" name:"CosResourceId"`
}

func (r *DescribeResourceTagsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "CreateUin")
	delete(f, "ResourceRegion")
	delete(f, "ServiceType")
	delete(f, "ResourcePrefix")
	delete(f, "ResourceId")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "CosResourceId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeResourceTagsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourceTagsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
	// 注意：此字段可能返回 null，表示取不到有效值。
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 资源标签
		Rows []*TagResource `json:"Rows,omitempty" name:"Rows"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeResourceTagsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourceTagsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourcesByTagsRequest struct {
	*tchttp.BaseRequest

	// 标签过滤数组
	TagFilters []*TagFilter `json:"TagFilters,omitempty" name:"TagFilters"`

	// 创建标签者uin
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 资源前缀
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`

	// 资源唯一标记
	ResourceId *string `json:"ResourceId,omitempty" name:"ResourceId"`

	// 资源所在地域
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 业务类型
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`
}

func (r *DescribeResourcesByTagsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourcesByTagsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagFilters")
	delete(f, "CreateUin")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "ResourcePrefix")
	delete(f, "ResourceId")
	delete(f, "ResourceRegion")
	delete(f, "ServiceType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeResourcesByTagsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourcesByTagsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
	// 注意：此字段可能返回 null，表示取不到有效值。
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 资源标签
		Rows []*ResourceTag `json:"Rows,omitempty" name:"Rows"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeResourcesByTagsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourcesByTagsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourcesByTagsUnionRequest struct {
	*tchttp.BaseRequest

	// 标签过滤数组
	TagFilters []*TagFilter `json:"TagFilters,omitempty" name:"TagFilters"`

	// 创建标签者uin
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 资源前缀
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`

	// 资源唯一标记
	ResourceId *string `json:"ResourceId,omitempty" name:"ResourceId"`

	// 资源所在地域
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 业务类型
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`
}

func (r *DescribeResourcesByTagsUnionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourcesByTagsUnionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagFilters")
	delete(f, "CreateUin")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "ResourcePrefix")
	delete(f, "ResourceId")
	delete(f, "ResourceRegion")
	delete(f, "ServiceType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeResourcesByTagsUnionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeResourcesByTagsUnionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 资源标签
		Rows []*ResourceTag `json:"Rows,omitempty" name:"Rows"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeResourcesByTagsUnionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeResourcesByTagsUnionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagKeysRequest struct {
	*tchttp.BaseRequest

	// 创建者用户 Uin，不传或为空只将 Uin 作为条件查询
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 是否展现项目
	ShowProject *uint64 `json:"ShowProject,omitempty" name:"ShowProject"`
}

func (r *DescribeTagKeysRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagKeysRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "CreateUin")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "ShowProject")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeTagKeysRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagKeysResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 标签列表
		Tags []*string `json:"Tags,omitempty" name:"Tags"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTagKeysResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagKeysResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagValuesRequest struct {
	*tchttp.BaseRequest

	// 标签键列表
	TagKeys []*string `json:"TagKeys,omitempty" name:"TagKeys"`

	// 创建者用户 Uin，不传或为空只将 Uin 作为条件查询
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`
}

func (r *DescribeTagValuesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagValuesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKeys")
	delete(f, "CreateUin")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeTagValuesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagValuesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 标签列表
		Tags []*Tag `json:"Tags,omitempty" name:"Tags"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTagValuesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagValuesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagValuesSeqRequest struct {
	*tchttp.BaseRequest

	// 标签键列表
	TagKeys []*string `json:"TagKeys,omitempty" name:"TagKeys"`

	// 创建者用户 Uin，不传或为空只将 Uin 作为条件查询
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`
}

func (r *DescribeTagValuesSeqRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagValuesSeqRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKeys")
	delete(f, "CreateUin")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeTagValuesSeqRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagValuesSeqResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 标签列表
		Tags []*Tag `json:"Tags,omitempty" name:"Tags"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTagValuesSeqResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagValuesSeqResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagsRequest struct {
	*tchttp.BaseRequest

	// 标签键,与标签值同时存在或同时不存在，不存在时表示查询该用户所有标签
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值,与标签键同时存在或同时不存在，不存在时表示查询该用户所有标签
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 创建者用户 Uin，不传或为空只将 Uin 作为条件查询
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 标签键数组,与标签值同时存在或同时不存在，不存在时表示查询该用户所有标签,当与TagKey同时传递时只取本值
	TagKeys []*string `json:"TagKeys,omitempty" name:"TagKeys"`

	// 是否展现项目标签
	ShowProject *uint64 `json:"ShowProject,omitempty" name:"ShowProject"`
}

func (r *DescribeTagsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKey")
	delete(f, "TagValue")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "CreateUin")
	delete(f, "TagKeys")
	delete(f, "ShowProject")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeTagsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 标签列表
		Tags []*TagWithDelete `json:"Tags,omitempty" name:"Tags"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTagsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagsSeqRequest struct {
	*tchttp.BaseRequest

	// 标签键,与标签值同时存在或同时不存在，不存在时表示查询该用户所有标签
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值,与标签键同时存在或同时不存在，不存在时表示查询该用户所有标签
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// 数据偏移量，默认为 0, 必须为Limit参数的整数倍
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 每页大小，默认为 15
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 创建者用户 Uin，不传或为空只将 Uin 作为条件查询
	CreateUin *uint64 `json:"CreateUin,omitempty" name:"CreateUin"`

	// 标签键数组,与标签值同时存在或同时不存在，不存在时表示查询该用户所有标签,当与TagKey同时传递时只取本值
	TagKeys []*string `json:"TagKeys,omitempty" name:"TagKeys"`

	// 是否展现项目标签
	ShowProject *uint64 `json:"ShowProject,omitempty" name:"ShowProject"`
}

func (r *DescribeTagsSeqRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagsSeqRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKey")
	delete(f, "TagValue")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "CreateUin")
	delete(f, "TagKeys")
	delete(f, "ShowProject")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeTagsSeqRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTagsSeqResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 结果总数
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 数据位移偏量
		Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

		// 每页大小
		Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

		// 标签列表
		Tags []*TagWithDelete `json:"Tags,omitempty" name:"Tags"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTagsSeqResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeTagsSeqResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DetachResourcesTagRequest struct {
	*tchttp.BaseRequest

	// 资源所属业务名称（资源六段式中的第三段）
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源ID数组，资源个数最多为50
	ResourceIds []*string `json:"ResourceIds,omitempty" name:"ResourceIds"`

	// 需要解绑的标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 资源所在地域，不区分地域的资源不需要传入该字段，区分地域的资源必填
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 资源前缀（资源六段式中最后一段"/"前面的部分），cos存储桶不需要传入该字段，其他云资源必填
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`
}

func (r *DetachResourcesTagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DetachResourcesTagRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ServiceType")
	delete(f, "ResourceIds")
	delete(f, "TagKey")
	delete(f, "ResourceRegion")
	delete(f, "ResourcePrefix")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DetachResourcesTagRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type DetachResourcesTagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DetachResourcesTagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DetachResourcesTagResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type ModifyResourceTagsRequest struct {
	*tchttp.BaseRequest

	// [ 资源六段式描述 ](https://cloud.tencent.com/document/product/598/10606)
	Resource *string `json:"Resource,omitempty" name:"Resource"`

	// 需要增加或修改的标签集合。如果Resource描述的资源未关联输入的标签键，则增加关联；若已关联，则将该资源关联的键对应的标签值修改为输入值。本接口中ReplaceTags和DeleteTags二者必须存在其一，且二者不能包含相同的标签键。可以不传该参数，但不能是空数组。
	ReplaceTags []*Tag `json:"ReplaceTags,omitempty" name:"ReplaceTags"`

	// 需要解关联的标签集合。本接口中ReplaceTags和DeleteTags二者必须存在其一，且二者不能包含相同的标签键。可以不传该参数，但不能是空数组。
	DeleteTags []*TagKeyObject `json:"DeleteTags,omitempty" name:"DeleteTags"`
}

func (r *ModifyResourceTagsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyResourceTagsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Resource")
	delete(f, "ReplaceTags")
	delete(f, "DeleteTags")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyResourceTagsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type ModifyResourceTagsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyResourceTagsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyResourceTagsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type ModifyResourcesTagValueRequest struct {
	*tchttp.BaseRequest

	// 资源所属业务名称（资源六段式中的第三段）
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源ID数组，资源个数最多为50
	ResourceIds []*string `json:"ResourceIds,omitempty" name:"ResourceIds"`

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// 资源所在地域，不区分地域的资源不需要传入该字段，区分地域的资源必填
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 资源前缀（资源六段式中最后一段"/"前面的部分），cos存储桶不需要传入该字段，其他云资源必填
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`
}

func (r *ModifyResourcesTagValueRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyResourcesTagValueRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ServiceType")
	delete(f, "ResourceIds")
	delete(f, "TagKey")
	delete(f, "TagValue")
	delete(f, "ResourceRegion")
	delete(f, "ResourcePrefix")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyResourcesTagValueRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type ModifyResourcesTagValueResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyResourcesTagValueResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyResourcesTagValueResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type ResourceIdTag struct {

	// 资源唯一标识
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResourceId *string `json:"ResourceId,omitempty" name:"ResourceId"`

	// 标签键值对
	// 注意：此字段可能返回 null，表示取不到有效值。
	TagKeyValues []*Tag `json:"TagKeyValues,omitempty" name:"TagKeyValues"`
}

type ResourceTag struct {

	// 资源所在地域
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResourceRegion *string `json:"ResourceRegion,omitempty" name:"ResourceRegion"`

	// 业务类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`

	// 资源前缀
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResourcePrefix *string `json:"ResourcePrefix,omitempty" name:"ResourcePrefix"`

	// 资源唯一标记
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResourceId *string `json:"ResourceId,omitempty" name:"ResourceId"`

	// 资源标签
	// 注意：此字段可能返回 null，表示取不到有效值。
	Tags []*Tag `json:"Tags,omitempty" name:"Tags"`
}

type Tag struct {

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`
}

type TagFilter struct {

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值数组 多个值的话是或的关系
	TagValue []*string `json:"TagValue,omitempty" name:"TagValue"`
}

type TagKeyObject struct {

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`
}

type TagResource struct {

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// 资源ID
	ResourceId *string `json:"ResourceId,omitempty" name:"ResourceId"`

	// 标签键MD5值
	TagKeyMd5 *string `json:"TagKeyMd5,omitempty" name:"TagKeyMd5"`

	// 标签值MD5值
	TagValueMd5 *string `json:"TagValueMd5,omitempty" name:"TagValueMd5"`

	// 资源类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	ServiceType *string `json:"ServiceType,omitempty" name:"ServiceType"`
}

type TagWithDelete struct {

	// 标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// 是否可以删除
	CanDelete *uint64 `json:"CanDelete,omitempty" name:"CanDelete"`
}

type UpdateResourceTagValueRequest struct {
	*tchttp.BaseRequest

	// 资源关联的标签键
	TagKey *string `json:"TagKey,omitempty" name:"TagKey"`

	// 修改后的标签值
	TagValue *string `json:"TagValue,omitempty" name:"TagValue"`

	// [ 资源六段式描述 ](https://cloud.tencent.com/document/product/598/10606)
	Resource *string `json:"Resource,omitempty" name:"Resource"`
}

func (r *UpdateResourceTagValueRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateResourceTagValueRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TagKey")
	delete(f, "TagValue")
	delete(f, "Resource")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "UpdateResourceTagValueRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

type UpdateResourceTagValueResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateResourceTagValueResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateResourceTagValueResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}
