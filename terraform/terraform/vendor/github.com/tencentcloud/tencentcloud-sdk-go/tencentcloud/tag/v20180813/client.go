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
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-08-13"

type Client struct {
    common.Client
}

// Deprecated
func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
    cpf := profile.NewClientProfile()
    client = &Client{}
    client.Init(region).WithSecretId(secretId, secretKey).WithProfile(cpf)
    return
}

func NewClient(credential common.CredentialIface, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
    client = &Client{}
    client.Init(region).
        WithCredential(credential).
        WithProfile(clientProfile)
    return
}


func NewAddResourceTagRequest() (request *AddResourceTagRequest) {
    request = &AddResourceTagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "AddResourceTag")
    return
}

func NewAddResourceTagResponse() (response *AddResourceTagResponse) {
    response = &AddResourceTagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AddResourceTag
// 本接口用于给标签关联资源
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_RESERVEDTAGKEY = "InvalidParameter.ReservedTagKey"
//  INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGKEYLENGTHEXCEEDED = "InvalidParameterValue.TagKeyLengthExceeded"
//  INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUELENGTHEXCEEDED = "InvalidParameterValue.TagValueLengthExceeded"
//  LIMITEXCEEDED_RESOURCEATTACHEDTAGS = "LimitExceeded.ResourceAttachedTags"
//  LIMITEXCEEDED_TAGKEY = "LimitExceeded.TagKey"
//  LIMITEXCEEDED_TAGVALUE = "LimitExceeded.TagValue"
//  RESOURCEINUSE_TAGKEYATTACHED = "ResourceInUse.TagKeyAttached"
func (c *Client) AddResourceTag(request *AddResourceTagRequest) (response *AddResourceTagResponse, err error) {
    if request == nil {
        request = NewAddResourceTagRequest()
    }
    response = NewAddResourceTagResponse()
    err = c.Send(request, response)
    return
}

func NewAttachResourcesTagRequest() (request *AttachResourcesTagRequest) {
    request = &AttachResourcesTagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "AttachResourcesTag")
    return
}

func NewAttachResourcesTagResponse() (response *AttachResourcesTagResponse) {
    response = &AttachResourcesTagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AttachResourcesTag
// 给多个资源关联某个标签
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEAPPIDNOTSAME = "FailedOperation.ResourceAppIdNotSame"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_RESERVEDTAGKEY = "InvalidParameter.ReservedTagKey"
//  INVALIDPARAMETERVALUE_REGIONINVALID = "InvalidParameterValue.RegionInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"
//  INVALIDPARAMETERVALUE_RESOURCEIDINVALID = "InvalidParameterValue.ResourceIdInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEPREFIXINVALID = "InvalidParameterValue.ResourcePrefixInvalid"
//  INVALIDPARAMETERVALUE_SERVICETYPEINVALID = "InvalidParameterValue.ServiceTypeInvalid"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGKEYEMPTY = "InvalidParameterValue.TagKeyEmpty"
//  INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
//  LIMITEXCEEDED_RESOURCEATTACHEDTAGS = "LimitExceeded.ResourceAttachedTags"
//  LIMITEXCEEDED_RESOURCENUMPERREQUEST = "LimitExceeded.ResourceNumPerRequest"
//  OPERATIONDENIED = "OperationDenied"
//  RESOURCEINUSE_TAGKEYATTACHED = "ResourceInUse.TagKeyAttached"
//  RESOURCENOTFOUND_TAGNONEXIST = "ResourceNotFound.TagNonExist"
func (c *Client) AttachResourcesTag(request *AttachResourcesTagRequest) (response *AttachResourcesTagResponse, err error) {
    if request == nil {
        request = NewAttachResourcesTagRequest()
    }
    response = NewAttachResourcesTagResponse()
    err = c.Send(request, response)
    return
}

func NewCreateTagRequest() (request *CreateTagRequest) {
    request = &CreateTagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "CreateTag")
    return
}

func NewCreateTagResponse() (response *CreateTagResponse) {
    response = &CreateTagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateTag
// 本接口用于创建一对标签键和标签值
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_RESERVEDTAGKEY = "InvalidParameter.ReservedTagKey"
//  INVALIDPARAMETERVALUE_RESERVEDTAGKEY = "InvalidParameterValue.ReservedTagKey"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGKEYEMPTY = "InvalidParameterValue.TagKeyEmpty"
//  INVALIDPARAMETERVALUE_TAGKEYLENGTHEXCEEDED = "InvalidParameterValue.TagKeyLengthExceeded"
//  INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUELENGTHEXCEEDED = "InvalidParameterValue.TagValueLengthExceeded"
//  LIMITEXCEEDED_TAGKEY = "LimitExceeded.TagKey"
//  LIMITEXCEEDED_TAGVALUE = "LimitExceeded.TagValue"
//  RESOURCEINUSE_TAGDUPLICATE = "ResourceInUse.TagDuplicate"
func (c *Client) CreateTag(request *CreateTagRequest) (response *CreateTagResponse, err error) {
    if request == nil {
        request = NewCreateTagRequest()
    }
    response = NewCreateTagResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteResourceTagRequest() (request *DeleteResourceTagRequest) {
    request = &DeleteResourceTagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DeleteResourceTag")
    return
}

func NewDeleteResourceTagResponse() (response *DeleteResourceTagResponse) {
    response = &DeleteResourceTagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteResourceTag
// 本接口用于解除标签和资源的关联关系
//
// 可能返回的错误码:
//  INVALIDPARAMETER_RESERVEDTAGKEY = "InvalidParameter.ReservedTagKey"
//  INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGKEYEMPTY = "InvalidParameterValue.TagKeyEmpty"
//  INVALIDPARAMETERVALUE_TAGKEYLENGTHEXCEEDED = "InvalidParameterValue.TagKeyLengthExceeded"
//  RESOURCENOTFOUND_ATTACHEDTAGKEYNOTFOUND = "ResourceNotFound.AttachedTagKeyNotFound"
func (c *Client) DeleteResourceTag(request *DeleteResourceTagRequest) (response *DeleteResourceTagResponse, err error) {
    if request == nil {
        request = NewDeleteResourceTagRequest()
    }
    response = NewDeleteResourceTagResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteTagRequest() (request *DeleteTagRequest) {
    request = &DeleteTagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DeleteTag")
    return
}

func NewDeleteTagResponse() (response *DeleteTagResponse) {
    response = &DeleteTagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteTag
// 本接口用于删除一对标签键和标签值
//
// 可能返回的错误码:
//  FAILEDOPERATION_TAGATTACHEDRESOURCE = "FailedOperation.TagAttachedResource"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUELENGTHEXCEEDED = "InvalidParameterValue.TagValueLengthExceeded"
//  RESOURCENOTFOUND_TAGNONEXIST = "ResourceNotFound.TagNonExist"
func (c *Client) DeleteTag(request *DeleteTagRequest) (response *DeleteTagResponse, err error) {
    if request == nil {
        request = NewDeleteTagRequest()
    }
    response = NewDeleteTagResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeResourceTagsRequest() (request *DescribeResourceTagsRequest) {
    request = &DescribeResourceTagsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeResourceTags")
    return
}

func NewDescribeResourceTagsResponse() (response *DescribeResourceTagsResponse) {
    response = &DescribeResourceTagsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeResourceTags
// 查询资源关联标签
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
func (c *Client) DescribeResourceTags(request *DescribeResourceTagsRequest) (response *DescribeResourceTagsResponse, err error) {
    if request == nil {
        request = NewDescribeResourceTagsRequest()
    }
    response = NewDescribeResourceTagsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeResourceTagsByResourceIdsRequest() (request *DescribeResourceTagsByResourceIdsRequest) {
    request = &DescribeResourceTagsByResourceIdsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeResourceTagsByResourceIds")
    return
}

func NewDescribeResourceTagsByResourceIdsResponse() (response *DescribeResourceTagsByResourceIdsResponse) {
    response = &DescribeResourceTagsByResourceIdsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeResourceTagsByResourceIds
// 用于批量查询已有资源关联的标签键值对
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_REGIONINVALID = "InvalidParameterValue.RegionInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEIDINVALID = "InvalidParameterValue.ResourceIdInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEPREFIXINVALID = "InvalidParameterValue.ResourcePrefixInvalid"
//  INVALIDPARAMETERVALUE_SERVICETYPEINVALID = "InvalidParameterValue.ServiceTypeInvalid"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeResourceTagsByResourceIds(request *DescribeResourceTagsByResourceIdsRequest) (response *DescribeResourceTagsByResourceIdsResponse, err error) {
    if request == nil {
        request = NewDescribeResourceTagsByResourceIdsRequest()
    }
    response = NewDescribeResourceTagsByResourceIdsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeResourceTagsByResourceIdsSeqRequest() (request *DescribeResourceTagsByResourceIdsSeqRequest) {
    request = &DescribeResourceTagsByResourceIdsSeqRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeResourceTagsByResourceIdsSeq")
    return
}

func NewDescribeResourceTagsByResourceIdsSeqResponse() (response *DescribeResourceTagsByResourceIdsSeqResponse) {
    response = &DescribeResourceTagsByResourceIdsSeqResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeResourceTagsByResourceIdsSeq
// 按顺序查看资源关联的标签
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_OFFSETINVALID = "InvalidParameterValue.OffsetInvalid"
//  INVALIDPARAMETERVALUE_REGIONINVALID = "InvalidParameterValue.RegionInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEIDINVALID = "InvalidParameterValue.ResourceIdInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEPREFIXINVALID = "InvalidParameterValue.ResourcePrefixInvalid"
//  INVALIDPARAMETERVALUE_SERVICETYPEINVALID = "InvalidParameterValue.ServiceTypeInvalid"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeResourceTagsByResourceIdsSeq(request *DescribeResourceTagsByResourceIdsSeqRequest) (response *DescribeResourceTagsByResourceIdsSeqResponse, err error) {
    if request == nil {
        request = NewDescribeResourceTagsByResourceIdsSeqRequest()
    }
    response = NewDescribeResourceTagsByResourceIdsSeqResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeResourceTagsByTagKeysRequest() (request *DescribeResourceTagsByTagKeysRequest) {
    request = &DescribeResourceTagsByTagKeysRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeResourceTagsByTagKeys")
    return
}

func NewDescribeResourceTagsByTagKeysResponse() (response *DescribeResourceTagsByTagKeysResponse) {
    response = &DescribeResourceTagsByTagKeysResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeResourceTagsByTagKeys
// 根据标签键获取资源标签
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_TAG = "InvalidParameter.Tag"
//  INVALIDPARAMETERVALUE_REGIONINVALID = "InvalidParameterValue.RegionInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEIDINVALID = "InvalidParameterValue.ResourceIdInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEPREFIXINVALID = "InvalidParameterValue.ResourcePrefixInvalid"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGKEYEMPTY = "InvalidParameterValue.TagKeyEmpty"
//  INVALIDPARAMETERVALUE_TAGKEYLENGTHEXCEEDED = "InvalidParameterValue.TagKeyLengthExceeded"
func (c *Client) DescribeResourceTagsByTagKeys(request *DescribeResourceTagsByTagKeysRequest) (response *DescribeResourceTagsByTagKeysResponse, err error) {
    if request == nil {
        request = NewDescribeResourceTagsByTagKeysRequest()
    }
    response = NewDescribeResourceTagsByTagKeysResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeResourcesByTagsRequest() (request *DescribeResourcesByTagsRequest) {
    request = &DescribeResourcesByTagsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeResourcesByTags")
    return
}

func NewDescribeResourcesByTagsResponse() (response *DescribeResourcesByTagsResponse) {
    response = &DescribeResourcesByTagsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeResourcesByTags
// 通过标签查询资源列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_TAGFILTERS = "InvalidParameterValue.TagFilters"
//  INVALIDPARAMETERVALUE_TAGFILTERSLENGTHEXCEEDED = "InvalidParameterValue.TagFiltersLengthExceeded"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeResourcesByTags(request *DescribeResourcesByTagsRequest) (response *DescribeResourcesByTagsResponse, err error) {
    if request == nil {
        request = NewDescribeResourcesByTagsRequest()
    }
    response = NewDescribeResourcesByTagsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeResourcesByTagsUnionRequest() (request *DescribeResourcesByTagsUnionRequest) {
    request = &DescribeResourcesByTagsUnionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeResourcesByTagsUnion")
    return
}

func NewDescribeResourcesByTagsUnionResponse() (response *DescribeResourcesByTagsUnionResponse) {
    response = &DescribeResourcesByTagsUnionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeResourcesByTagsUnion
// 通过标签查询资源列表并集
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_TAGFILTERS = "InvalidParameterValue.TagFilters"
//  INVALIDPARAMETERVALUE_TAGFILTERSLENGTHEXCEEDED = "InvalidParameterValue.TagFiltersLengthExceeded"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeResourcesByTagsUnion(request *DescribeResourcesByTagsUnionRequest) (response *DescribeResourcesByTagsUnionResponse, err error) {
    if request == nil {
        request = NewDescribeResourcesByTagsUnionRequest()
    }
    response = NewDescribeResourcesByTagsUnionResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTagKeysRequest() (request *DescribeTagKeysRequest) {
    request = &DescribeTagKeysRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeTagKeys")
    return
}

func NewDescribeTagKeysResponse() (response *DescribeTagKeysResponse) {
    response = &DescribeTagKeysResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTagKeys
// 用于查询已建立的标签列表中的标签键。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeTagKeys(request *DescribeTagKeysRequest) (response *DescribeTagKeysResponse, err error) {
    if request == nil {
        request = NewDescribeTagKeysRequest()
    }
    response = NewDescribeTagKeysResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTagValuesRequest() (request *DescribeTagValuesRequest) {
    request = &DescribeTagValuesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeTagValues")
    return
}

func NewDescribeTagValuesResponse() (response *DescribeTagValuesResponse) {
    response = &DescribeTagValuesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTagValues
// 用于查询已建立的标签列表中的标签值。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeTagValues(request *DescribeTagValuesRequest) (response *DescribeTagValuesResponse, err error) {
    if request == nil {
        request = NewDescribeTagValuesRequest()
    }
    response = NewDescribeTagValuesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTagValuesSeqRequest() (request *DescribeTagValuesSeqRequest) {
    request = &DescribeTagValuesSeqRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeTagValuesSeq")
    return
}

func NewDescribeTagValuesSeqResponse() (response *DescribeTagValuesSeqResponse) {
    response = &DescribeTagValuesSeqResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTagValuesSeq
// 用于查询已建立的标签列表中的标签值。
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_OFFSETINVALID = "InvalidParameterValue.OffsetInvalid"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeTagValuesSeq(request *DescribeTagValuesSeqRequest) (response *DescribeTagValuesSeqResponse, err error) {
    if request == nil {
        request = NewDescribeTagValuesSeqRequest()
    }
    response = NewDescribeTagValuesSeqResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTagsRequest() (request *DescribeTagsRequest) {
    request = &DescribeTagsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeTags")
    return
}

func NewDescribeTagsResponse() (response *DescribeTagsResponse) {
    response = &DescribeTagsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTags
// 用于查询已建立的标签列表。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeTags(request *DescribeTagsRequest) (response *DescribeTagsResponse, err error) {
    if request == nil {
        request = NewDescribeTagsRequest()
    }
    response = NewDescribeTagsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTagsSeqRequest() (request *DescribeTagsSeqRequest) {
    request = &DescribeTagsSeqRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DescribeTagsSeq")
    return
}

func NewDescribeTagsSeqResponse() (response *DescribeTagsSeqResponse) {
    response = &DescribeTagsSeqResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTagsSeq
// 用于查询已建立的标签列表。
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_OFFSETINVALID = "InvalidParameterValue.OffsetInvalid"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
func (c *Client) DescribeTagsSeq(request *DescribeTagsSeqRequest) (response *DescribeTagsSeqResponse, err error) {
    if request == nil {
        request = NewDescribeTagsSeqRequest()
    }
    response = NewDescribeTagsSeqResponse()
    err = c.Send(request, response)
    return
}

func NewDetachResourcesTagRequest() (request *DetachResourcesTagRequest) {
    request = &DetachResourcesTagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "DetachResourcesTag")
    return
}

func NewDetachResourcesTagResponse() (response *DetachResourcesTagResponse) {
    response = &DetachResourcesTagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DetachResourcesTag
// 解绑多个资源关联的某个标签
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEAPPIDNOTSAME = "FailedOperation.ResourceAppIdNotSame"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_REGIONINVALID = "InvalidParameterValue.RegionInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"
//  INVALIDPARAMETERVALUE_RESOURCEIDINVALID = "InvalidParameterValue.ResourceIdInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEPREFIXINVALID = "InvalidParameterValue.ResourcePrefixInvalid"
//  INVALIDPARAMETERVALUE_SERVICETYPEINVALID = "InvalidParameterValue.ServiceTypeInvalid"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  LIMITEXCEEDED_RESOURCENUMPERREQUEST = "LimitExceeded.ResourceNumPerRequest"
//  RESOURCENOTFOUND_ATTACHEDTAGKEYNOTFOUND = "ResourceNotFound.AttachedTagKeyNotFound"
func (c *Client) DetachResourcesTag(request *DetachResourcesTagRequest) (response *DetachResourcesTagResponse, err error) {
    if request == nil {
        request = NewDetachResourcesTagRequest()
    }
    response = NewDetachResourcesTagResponse()
    err = c.Send(request, response)
    return
}

func NewModifyResourceTagsRequest() (request *ModifyResourceTagsRequest) {
    request = &ModifyResourceTagsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "ModifyResourceTags")
    return
}

func NewModifyResourceTagsResponse() (response *ModifyResourceTagsResponse) {
    response = &ModifyResourceTagsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyResourceTags
// 本接口用于修改资源关联的所有标签
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_RESERVEDTAGKEY = "InvalidParameter.ReservedTagKey"
//  INVALIDPARAMETER_TAG = "InvalidParameter.Tag"
//  INVALIDPARAMETERVALUE_DELETETAGSPARAMERROR = "InvalidParameterValue.DeleteTagsParamError"
//  INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGKEYEMPTY = "InvalidParameterValue.TagKeyEmpty"
//  INVALIDPARAMETERVALUE_TAGKEYLENGTHEXCEEDED = "InvalidParameterValue.TagKeyLengthExceeded"
//  INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUELENGTHEXCEEDED = "InvalidParameterValue.TagValueLengthExceeded"
//  LIMITEXCEEDED_RESOURCEATTACHEDTAGS = "LimitExceeded.ResourceAttachedTags"
//  LIMITEXCEEDED_TAGKEY = "LimitExceeded.TagKey"
//  LIMITEXCEEDED_TAGVALUE = "LimitExceeded.TagValue"
func (c *Client) ModifyResourceTags(request *ModifyResourceTagsRequest) (response *ModifyResourceTagsResponse, err error) {
    if request == nil {
        request = NewModifyResourceTagsRequest()
    }
    response = NewModifyResourceTagsResponse()
    err = c.Send(request, response)
    return
}

func NewModifyResourcesTagValueRequest() (request *ModifyResourcesTagValueRequest) {
    request = &ModifyResourcesTagValueRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "ModifyResourcesTagValue")
    return
}

func NewModifyResourcesTagValueResponse() (response *ModifyResourcesTagValueResponse) {
    response = &ModifyResourcesTagValueResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyResourcesTagValue
// 修改多个资源关联的某个标签键对应的标签值
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEAPPIDNOTSAME = "FailedOperation.ResourceAppIdNotSame"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_REGIONINVALID = "InvalidParameterValue.RegionInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"
//  INVALIDPARAMETERVALUE_RESOURCEIDINVALID = "InvalidParameterValue.ResourceIdInvalid"
//  INVALIDPARAMETERVALUE_RESOURCEPREFIXINVALID = "InvalidParameterValue.ResourcePrefixInvalid"
//  INVALIDPARAMETERVALUE_SERVICETYPEINVALID = "InvalidParameterValue.ServiceTypeInvalid"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"
//  INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"
//  LIMITEXCEEDED_RESOURCENUMPERREQUEST = "LimitExceeded.ResourceNumPerRequest"
//  RESOURCENOTFOUND_ATTACHEDTAGKEYNOTFOUND = "ResourceNotFound.AttachedTagKeyNotFound"
//  RESOURCENOTFOUND_TAGNONEXIST = "ResourceNotFound.TagNonExist"
func (c *Client) ModifyResourcesTagValue(request *ModifyResourcesTagValueRequest) (response *ModifyResourcesTagValueResponse, err error) {
    if request == nil {
        request = NewModifyResourcesTagValueRequest()
    }
    response = NewModifyResourcesTagValueResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateResourceTagValueRequest() (request *UpdateResourceTagValueRequest) {
    request = &UpdateResourceTagValueRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tag", APIVersion, "UpdateResourceTagValue")
    return
}

func NewUpdateResourceTagValueResponse() (response *UpdateResourceTagValueResponse) {
    response = &UpdateResourceTagValueResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateResourceTagValue
// 本接口用于修改资源已关联的标签值（标签键不变）
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"
//  INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"
//  INVALIDPARAMETERVALUE_TAGVALUELENGTHEXCEEDED = "InvalidParameterValue.TagValueLengthExceeded"
//  LIMITEXCEEDED_TAGVALUE = "LimitExceeded.TagValue"
//  RESOURCENOTFOUND_ATTACHEDTAGKEYNOTFOUND = "ResourceNotFound.AttachedTagKeyNotFound"
func (c *Client) UpdateResourceTagValue(request *UpdateResourceTagValueRequest) (response *UpdateResourceTagValueResponse, err error) {
    if request == nil {
        request = NewUpdateResourceTagValueRequest()
    }
    response = NewUpdateResourceTagValueResponse()
    err = c.Send(request, response)
    return
}
