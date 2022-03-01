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

const (
	// 此产品的特有错误码

	// 未通过CAM鉴权。
	AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"

	// 操作失败。
	FAILEDOPERATION = "FailedOperation"

	// 单次请求的资源appId必须相同。
	FAILEDOPERATION_RESOURCEAPPIDNOTSAME = "FailedOperation.ResourceAppIdNotSame"

	// 已关联资源的标签无法删除。
	FAILEDOPERATION_TAGATTACHEDRESOURCE = "FailedOperation.TagAttachedResource"

	// 参数错误。
	INVALIDPARAMETER = "InvalidParameter"

	// 系统预留标签键 qcloud、tencent和project 禁止创建。
	INVALIDPARAMETER_RESERVEDTAGKEY = "InvalidParameter.ReservedTagKey"

	// Tag参数错误。
	INVALIDPARAMETER_TAG = "InvalidParameter.Tag"

	// DeleteTags中不能包含ReplaceTags或AddTags中元素。
	INVALIDPARAMETERVALUE_DELETETAGSPARAMERROR = "InvalidParameterValue.DeleteTagsParamError"

	// offset error。
	INVALIDPARAMETERVALUE_OFFSETINVALID = "InvalidParameterValue.OffsetInvalid"

	// 地域错误。
	INVALIDPARAMETERVALUE_REGIONINVALID = "InvalidParameterValue.RegionInvalid"

	// 系统预留标签键 qcloud、tencent和project 禁止创建。
	INVALIDPARAMETERVALUE_RESERVEDTAGKEY = "InvalidParameterValue.ReservedTagKey"

	// 资源描述错误。
	INVALIDPARAMETERVALUE_RESOURCEDESCRIPTIONERROR = "InvalidParameterValue.ResourceDescriptionError"

	// 资源Id错误。
	INVALIDPARAMETERVALUE_RESOURCEIDINVALID = "InvalidParameterValue.ResourceIdInvalid"

	// 资源前缀错误。
	INVALIDPARAMETERVALUE_RESOURCEPREFIXINVALID = "InvalidParameterValue.ResourcePrefixInvalid"

	// 业务类型错误。
	INVALIDPARAMETERVALUE_SERVICETYPEINVALID = "InvalidParameterValue.ServiceTypeInvalid"

	// TagFilters参数错误。
	INVALIDPARAMETERVALUE_TAGFILTERS = "InvalidParameterValue.TagFilters"

	// 过滤标签键对应标签值达到上限数 6。
	INVALIDPARAMETERVALUE_TAGFILTERSLENGTHEXCEEDED = "InvalidParameterValue.TagFiltersLengthExceeded"

	// 标签键包含非法字符。
	INVALIDPARAMETERVALUE_TAGKEYCHARACTERILLEGAL = "InvalidParameterValue.TagKeyCharacterIllegal"

	// 标签键不能为空值。
	INVALIDPARAMETERVALUE_TAGKEYEMPTY = "InvalidParameterValue.TagKeyEmpty"

	// 标签键长度超过限制。
	INVALIDPARAMETERVALUE_TAGKEYLENGTHEXCEEDED = "InvalidParameterValue.TagKeyLengthExceeded"

	// 标签值包含非法字符。
	INVALIDPARAMETERVALUE_TAGVALUECHARACTERILLEGAL = "InvalidParameterValue.TagValueCharacterIllegal"

	// 标签值长度超过限制。
	INVALIDPARAMETERVALUE_TAGVALUELENGTHEXCEEDED = "InvalidParameterValue.TagValueLengthExceeded"

	// Uin参数不合法。
	INVALIDPARAMETERVALUE_UININVALID = "InvalidParameterValue.UinInvalid"

	// 资源关联的标签数超过限制。
	LIMITEXCEEDED_RESOURCEATTACHEDTAGS = "LimitExceeded.ResourceAttachedTags"

	// 单次请求的资源数达到上限。
	LIMITEXCEEDED_RESOURCENUMPERREQUEST = "LimitExceeded.ResourceNumPerRequest"

	// 用户创建标签键达到上限数 1000。
	LIMITEXCEEDED_TAGKEY = "LimitExceeded.TagKey"

	// 单个标签键对应标签值达到上限数 1000。
	LIMITEXCEEDED_TAGVALUE = "LimitExceeded.TagValue"

	// 操作被拒绝。
	OPERATIONDENIED = "OperationDenied"

	// 标签已存在。
	RESOURCEINUSE_TAGDUPLICATE = "ResourceInUse.TagDuplicate"

	// 对应的标签键和资源已关联。
	RESOURCEINUSE_TAGKEYATTACHED = "ResourceInUse.TagKeyAttached"

	// 资源关联的标签键不存在。
	RESOURCENOTFOUND_ATTACHEDTAGKEYNOTFOUND = "ResourceNotFound.AttachedTagKeyNotFound"

	// 标签不存在。
	RESOURCENOTFOUND_TAGNONEXIST = "ResourceNotFound.TagNonExist"
)
