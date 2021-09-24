package datahub

import (
    "fmt"
)

/*
examples errors
*/

// Error codes
const (
    InvalidParameter    = "InvalidParameter"
    InvalidSubscription = "InvalidSubscription"
    InvalidCursor       = "InvalidCursor"
    /**
     * for later arrange error code
     */
    ResourceNotFound   = "ResourceNotFound"
    NoSuchTopic        = "NoSuchTopic"
    NoSuchProject      = "NoSuchProject"
    NoSuchSubscription = "NoSuchSubscription"
    NoSuchShard        = "NoSuchShard"
    NoSuchConnector    = "NoSuchConnector"
    NoSuchMeterInfo    = "NoSuchMeteringInfo"
    /**
     * for later arrange error code
     */
    SeekOutOfRange        = "SeekOutOfRange"
    ResourceAlreadyExist  = "ResourceAlreadyExist"
    ProjectAlreadyExist   = "ProjectAlreadyExist"
    TopicAlreadyExist     = "TopicAlreadyExist"
    ConnectorAlreadyExist = "ConnectorAlreadyExist"
    UnAuthorized          = "Unauthorized"
    NoPermission          = "NoPermission"
    InvalidShardOperation = "InvalidShardOperation"
    OperatorDenied        = "OperationDenied"
    LimitExceed           = "LimitExceeded"
    //ODPSServiceError       = "OdpsServiceError"
    //MysqlServiceError      = "MysqlServiceError"
    //InternalServerErrorS    = "InternalServerError"
    SubscriptionOffline    = "SubscriptionOffline"
    OffsetReseted          = "OffsetReseted"
    OffsetSessionClosed    = "OffsetSessionClosed"
    OffsetSessionChanged   = "OffsetSessionChanged"
    MalformedRecord        = "MalformedRecord"
    NoSuchConsumer         = "NoSuchConsumer"
    ConsumerGroupInProcess = "ConsumerGroupInProcess"
)

const (
    projectNameInvalid   string = "project name should start with letter, only contains [a-zA-Z0-9_], 3 < length < 32"
    commentInvalid       string = "comment can not be empty and length must less than 1024"
    topicNameInvalid     string = "topic name should start with letter, only contains [a-zA-Z0-9_], 1 < length < 128"
    shardIdInvalid       string = "shardId is invalid"
    shardListInvalid     string = "shard list is empty"
    lifecycleInvalid     string = "lifecycle is invalid"
    parameterInvalid     string = "parameter is invalid"
    parameterNull        string = "parameter is nil"
    parameterNumInvalid  string = "parameter num invalid"
    parameterTypeInvalid string = "parameter type is invalid,please check your input parameter type"
    missingRecordSchema  string = "missing record schema for tuple record type"
    recordsInvalid       string = "records is invalid, nil, empty or other invalid reason"
)

// return the specific err type by errCode,
// you can handle the error by type assert
func errorHandler(statusCode int, requestId string, errorCode string, message string) error {

    switch errorCode {
    case InvalidParameter, InvalidSubscription, InvalidCursor:
        return NewInvalidParameterError(statusCode, requestId, errorCode, message)
    case ResourceNotFound, NoSuchTopic, NoSuchProject, NoSuchSubscription, NoSuchShard, NoSuchConnector,
        NoSuchMeterInfo, NoSuchConsumer:
        return NewResourceNotFoundError(statusCode, requestId, errorCode, message)
    case SeekOutOfRange:
        return NewSeekOutOfRangeError(statusCode, requestId, errorCode, message)
    case ResourceAlreadyExist, ProjectAlreadyExist, TopicAlreadyExist, ConnectorAlreadyExist:
        return NewResourceExistError(statusCode, requestId, errorCode, message)
    case UnAuthorized:
        return NewAuthorizationFailedError(statusCode, requestId, errorCode, message)
    case NoPermission:
        return NewNoPermissionError(statusCode, requestId, errorCode, message)
    case OperatorDenied:
        return NewInvalidOperationError(statusCode, requestId, errorCode, message)
    case LimitExceed:
        return NewLimitExceededError(statusCode, requestId, errorCode, message)
    case SubscriptionOffline:
        return NewSubscriptionOfflineError(statusCode, requestId, errorCode, message)
    case OffsetReseted:
        return NewSubscriptionOffsetResetError(statusCode, requestId, errorCode, message)
    case OffsetSessionClosed, OffsetSessionChanged:
        return NewSubscriptionSessionInvalidError(statusCode, requestId, errorCode, message)
    case MalformedRecord:
        return NewMalformedRecordError(statusCode, requestId, errorCode, message)
    case ConsumerGroupInProcess:
        return NewServiceInProcessError(statusCode, requestId, errorCode, message)
    case InvalidShardOperation:
        return NewShardSealedError(statusCode, requestId, errorCode, message)
    }
    return NewDatahubClientError(statusCode, requestId, errorCode, message)
}

// create a new DatahubClientError
func NewDatahubClientError(statusCode int, requestId string, code string, message string) *DatahubClientError {
    return &DatahubClientError{StatusCode: statusCode, RequestId: requestId, Code: code, Message: message}
}

// DatahubError struct
type DatahubClientError struct {
    StatusCode int    `json:"StatusCode"`   // Http status code
    RequestId  string `json:"RequestId"`    // Request-id to trace the request
    Code       string `json:"ErrorCode"`    // Datahub error code
    Message    string `json:"ErrorMessage"` // Error msg of the error code
}

func (err *DatahubClientError) Error() string {
    return fmt.Sprintf("statusCode: %d, requestId: %s, errCode: %s, errMsg: %s",
        err.StatusCode, err.RequestId, err.Code, err.Message)
}

func NewInvalidParameterErrorWithMessage(message string) *InvalidParameterError {
    return &InvalidParameterError{
        DatahubClientError{
            StatusCode: -1,
            RequestId:  "",
            Code:       "",
            Message:    message,
        },
    }
}

func NewInvalidParameterError(statusCode int, requestId string, code string, message string) *InvalidParameterError {
    return &InvalidParameterError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

// InvalidParameterError represent the parameter error
type InvalidParameterError struct {
    DatahubClientError
}

func NewResourceNotFoundError(statusCode int, requestId string, code string, message string) *ResourceNotFoundError {
    return &ResourceNotFoundError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type ResourceNotFoundError struct {
    DatahubClientError
}

func NewResourceExistError(statusCode int, requestId string, code string, message string) *ResourceExistError {
    return &ResourceExistError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type ResourceExistError struct {
    DatahubClientError
}

func NewInvalidOperationError(statusCode int, requestId string, code string, message string) *InvalidOperationError {
    return &InvalidOperationError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type InvalidOperationError struct {
    DatahubClientError
}

func NewLimitExceededError(statusCode int, requestId string, code string, message string) *LimitExceededError {
    return &LimitExceededError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type LimitExceededError struct {
    DatahubClientError
}

func NewAuthorizationFailedError(statusCode int, requestId string, code string, message string) *AuthorizationFailedError {
    return &AuthorizationFailedError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type AuthorizationFailedError struct {
    DatahubClientError
}

//func (afe *AuthorizationFailureError) Error() string {
//    return afe.DatahubClientError.Error()
//}

func NewNoPermissionError(statusCode int, requestId string, code string, message string) *NoPermissionError {
    return &NoPermissionError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type NoPermissionError struct {
    DatahubClientError
}

func NewSeekOutOfRangeError(statusCode int, requestId string, code string, message string) *SeekOutOfRangeError {
    return &SeekOutOfRangeError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type SeekOutOfRangeError struct {
    DatahubClientError
}

func NewSubscriptionOfflineError(statusCode int, requestId string, code string, message string) *SubscriptionOfflineError {
    return &SubscriptionOfflineError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type SubscriptionOfflineError struct {
    DatahubClientError
}

func NewSubscriptionOffsetResetError(statusCode int, requestId string, code string, message string) *SubscriptionOffsetResetError {
    return &SubscriptionOffsetResetError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type SubscriptionOffsetResetError struct {
    DatahubClientError
}

func NewSubscriptionSessionInvalidError(statusCode int, requestId string, code string, message string) *SubscriptionSessionInvalidError {
    return &SubscriptionSessionInvalidError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type SubscriptionSessionInvalidError struct {
    DatahubClientError
}

func NewMalformedRecordError(statusCode int, requestId string, code string, message string) *MalformedRecordError {
    return &MalformedRecordError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type MalformedRecordError struct {
    DatahubClientError
}

func NewServiceInProcessError(statusCode int, requestId string, code string, message string) *ServiceInProcessError {
    return &ServiceInProcessError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type ServiceInProcessError struct {
    DatahubClientError
}

func NewShardSealedError(statusCode int, requestId string, code string, message string) *ShardSealedError {
    return &ShardSealedError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type ShardSealedError struct {
    DatahubClientError
}

func NewServiceTemporaryUnavailableError(message string) *ServiceTemporaryUnavailableError {
    return &ServiceTemporaryUnavailableError{
        DatahubClientError{
            StatusCode: -1,
            RequestId:  "",
            Code:       "",
            Message:    message,
        },
    }
}

func NewServiceTemporaryUnavailableErrorWithCode(statusCode int, requestId string, code string, message string) *ServiceTemporaryUnavailableError {
    return &ServiceTemporaryUnavailableError{
        DatahubClientError{
            StatusCode: statusCode,
            RequestId:  requestId,
            Code:       code,
            Message:    message,
        },
    }
}

type ServiceTemporaryUnavailableError struct {
    DatahubClientError
}
