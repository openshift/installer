package alicloud

import (
	"regexp"
	"strings"

	"github.com/alibabacloud-go/tea/tea"

	sls "github.com/aliyun/aliyun-log-go-sdk"

	"fmt"

	"log"
	"runtime"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/common"
)

const (
	// common
	NotFound                = "NotFound"
	ResourceNotfound        = "ResourceNotfound"
	ServiceUnavailable      = "ServiceUnavailable"
	InstanceNotFound        = "Instance.Notfound"
	ForbiddenInstance       = "Forbidden.InstanceNotFound"
	MessageInstanceNotFound = "instance is not found"
	Throttling              = "Throttling"
	ThrottlingUser          = "Throttling.User"

	// RAM Instance Not Found
	RamInstanceNotFound   = "Forbidden.InstanceNotFound"
	AliyunGoClientFailure = "AliyunGoClientFailure"

	LogClientTimeout = "Client.Timeout exceeded while awaiting headers"

	InvalidFileSystemStatus_Ordering = "InvalidFileSystemStatus.Ordering"
)

var SlbIsBusy = []string{"SystemBusy", "OperationBusy", "ServiceIsStopping", "BackendServer.configuring", "ServiceIsConfiguring"}
var EcsNotFound = []string{"InvalidInstanceId.NotFound", "Forbidden.InstanceNotFound"}
var DiskInvalidOperation = []string{"IncorrectDiskStatus", "IncorrectInstanceStatus", "OperationConflict", "InternalError", "InvalidOperation.Conflict", "IncorrectDiskStatus.Initializing"}
var NetworkInterfaceInvalidOperations = []string{"InvalidOperation.InvalidEniState", "InvalidOperation.InvalidEcsState", "OperationConflict", "ServiceUnavailable", "InternalError"}
var OperationDeniedDBStatus = []string{"OperationDenied.DBStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBClusterStatus", "InternalError", "OperationDenied.OutofUsage", "IncorrectDBInstanceState"}
var DBReadInstanceNotReadyStatus = []string{"OperationDenied.ReadDBInstanceStatus", "OperationDenied.MasterDBInstanceState", "ReadDBInstance.Mismatch"}
var NasNotFound = []string{"InvalidMountTarget.NotFound", "InvalidFileSystem.NotFound", "Forbidden.NasNotFound", "InvalidLBid.NotFound", "VolumeUnavailable"}
var SnapshotInvalidOperations = []string{"OperationConflict", "ServiceUnavailable", "InternalError", "SnapshotCreatedDisk", "SnapshotCreatedImage"}
var DiskNotSupportOnlineChangeErrors = []string{"InvalidDiskCategory.NotSupported", "InvalidRegion.NotSupport", "IncorrectInstanceStatus", "IncorrectDiskStatus", "InvalidOperation.InstanceTypeNotSupport"}

// details at: https://help.aliyun.com/document_detail/27300.html
var OtsTableIsTemporarilyUnavailable = []string{"no such host", "OTSServerBusy", "OTSPartitionUnavailable", "OTSInternalServerError",
	"OTSTimeout", "OTSServerUnavailable", "OTSRowOperationConflict", "OTSTableNotReady", "OTSNotEnoughCapacityUnit", "Too frequent table operations."}

// An Error represents a custom error for Terraform failure response
type ProviderError struct {
	errorCode string
	message   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[ERROR] Terraform Alicloud Provider Error: Code: %s Message: %s", e.errorCode, e.message)
}

func (err *ProviderError) ErrorCode() string {
	return err.errorCode
}

func (err *ProviderError) Message() string {
	return err.message
}

func GetNotFoundErrorFromString(str string) error {
	return &ProviderError{
		errorCode: InstanceNotFound,
		message:   str,
	}
}
func NotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*ComplexError); ok {
		if e.Err != nil && strings.HasPrefix(e.Err.Error(), ResourceNotfound) {
			return true
		}
		return NotFoundError(e.Cause)
	}
	if err == nil {
		return false
	}

	if e, ok := err.(*errors.ServerError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*ProviderError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == InstanceNotFound || e.Code == RamInstanceNotFound || e.Code == NotFound || strings.Contains(strings.ToLower(e.Message), MessageInstanceNotFound)
	}

	if e, ok := err.(oss.ServiceError); ok {
		return e.StatusCode == 404 || strings.HasPrefix(e.Code, "NoSuch") || strings.HasPrefix(e.Message, "No Row found")
	}

	return false
}

func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}

	if e, ok := err.(*tea.SDKError); ok {
		for _, code := range expectCodes {
			// The second statement aims to match the tea sdk history bug
			if *e.Code == code || strings.HasPrefix(code, *e.Code) || strings.Contains(*e.Data, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*errors.ServerError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*common.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*sls.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) || strings.Contains(e.String(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(oss.ServiceError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*fc.ServiceError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode == code || strings.Contains(e.ErrorMessage, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*datahub.DatahubClientError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

func NeedRetry(err error) bool {
	if err == nil {
		return false
	}

	postRegex := regexp.MustCompile("^Post [\"]*https://.*")
	if postRegex.MatchString(err.Error()) {
		return true
	}

	throttlingRegex := regexp.MustCompile("^Throttling.*")
	codeRegex := regexp.MustCompile("^code: 5[\\d]{2}")

	if e, ok := err.(*tea.SDKError); ok {
		if strings.Contains(*e.Message, "code: 500, 您已开通过") {
			return false
		}
		if *e.Code == ServiceUnavailable || *e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(*e.Code) || codeRegex.MatchString(*e.Message) {
			return true
		}
	}

	if e, ok := err.(*errors.ServerError); ok {
		return e.ErrorCode() == ServiceUnavailable || e.ErrorCode() == "Rejected.Throttling" || throttlingRegex.MatchString(e.ErrorCode()) || codeRegex.MatchString(e.Message())
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == ServiceUnavailable || e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(e.Code) || codeRegex.MatchString(e.Message)
	}

	return false
}

func IsExpectedErrorCodes(code string, errorCodes []string) bool {
	if code == "" {
		return false
	}
	for _, v := range errorCodes {
		if v == code {
			return true
		}
	}
	return false
}

func IsThrottling(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*errors.ServerError); ok {
		if e.ErrorCode() == Throttling {
			return true
		}
		return false
	}

	if e, ok := err.(*common.Error); ok {
		if e.Code == Throttling {
			return true
		}
		return false
	}
	return false
}

func GetTimeErrorFromString(str string) error {
	return &ProviderError{
		errorCode: "WaitForTimeout",
		message:   str,
	}
}

func GetNotFoundMessage(product, id string) string {
	return fmt.Sprintf("The specified %s %s is not found.", product, id)
}

func GetTimeoutMessage(product, status string) string {
	return fmt.Sprintf("Waitting for %s %s is timeout.", product, status)
}

type ErrorSource string

const (
	AlibabaCloudSdkGoERROR = ErrorSource("[SDK alibaba-cloud-sdk-go ERROR]")
	AliyunLogGoSdkERROR    = ErrorSource("[SDK aliyun-log-go-sdk ERROR]")
	AliyunDatahubSdkGo     = ErrorSource("[SDK aliyun-datahub-sdk-go ERROR]")
	AliyunMaxComputeSdkGo  = ErrorSource("[SDK aliyun-maxcompute-sdk-go ERROR]")
	AliyunOssGoSdk         = ErrorSource("[SDK aliyun-oss-go-sdk ERROR]")
	FcGoSdk                = ErrorSource("[SDK fc-go-sdk ERROR]")
	DenverdinoAliyungo     = ErrorSource("[SDK denverdino/aliyungo ERROR]")
	AliyunTablestoreGoSdk  = ErrorSource("[SDK aliyun-tablestore-go-sdk ERROR]")
	AliMnsERROR            = ErrorSource("[SDK ali_mns ERROR]")
	ProviderERROR          = ErrorSource("[Provider ERROR]")
)

// ComplexError is a format error which including origin error, extra error message, error occurred file and line
// Cause: a error is a origin error that comes from SDK, some exceptions and so on
// Err: a new error is built from extra message
// Path: the file path of error occurred
// Line: the file line of error occurred
type ComplexError struct {
	Cause error
	Err   error
	Path  string
	Line  int
}

func (e ComplexError) Error() string {
	if e.Cause == nil {
		e.Cause = Error("<nil cause>")
	}
	if e.Err == nil {
		return fmt.Sprintf("\u001B[31m[ERROR]\u001B[0m %s:%d:\n%s", e.Path, e.Line, e.Cause.Error())
	}
	return fmt.Sprintf("\u001B[31m[ERROR]\u001B[0m %s:%d: %s:\n%s", e.Path, e.Line, e.Err.Error(), e.Cause.Error())
}

func Error(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

// Return a ComplexError which including error occurred file and path
func WrapError(cause error) error {
	if cause == nil {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("\u001B[31m[ERROR]\u001B[0m runtime.Caller error in WrapError.")
		return WrapComplexError(cause, nil, "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, nil, filepath, line)
}

// Return a ComplexError which including extra error message, error occurred file and path
func WrapErrorf(cause error, msg string, args ...interface{}) error {
	if cause == nil && strings.TrimSpace(msg) == "" {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("\u001B[31m[ERROR]\u001B[0m runtime.Caller error in WrapErrorf.")
		return WrapComplexError(cause, Error(msg), "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	// The second parameter of args is requestId, if the error message is NotFoundMsg the requestId need to be returned.
	if msg == NotFoundMsg && len(args) == 2 {
		msg += RequestIdMsg
	}
	return WrapComplexError(cause, fmt.Errorf(msg, args...), filepath, line)
}

func WrapComplexError(cause, err error, filepath string, fileline int) error {
	return &ComplexError{
		Cause: cause,
		Err:   err,
		Path:  filepath,
		Line:  fileline,
	}
}

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const DefaultErrorMsg = "Resource %s %s Failed!!! %s"
const ResponseCodeMsg = "Resource %s %s Failed!!! %v"
const RequestIdMsg = "RequestId: %s"
const NotFoundMsg = ResourceNotfound + "!!! %s"
const NotFoundWithResponse = ResourceNotfound + "!!! Response: %v"
const DefaultTimeoutMsg = "Resource %s %s Timeout!!! %s"
const DeleteTimeoutMsg = "Resource %s Still Exists. %s Timeout!!! %s"
const WaitTimeoutMsg = "Resource %s %s Timeout In %d Seconds. Got: %s Expected: %s !!! %s"
const DataDefaultErrorMsg = "Datasource %s %s Failed!!! %s"
const SweepDefaultErrorMsg = "Sweep %s %s Failed!!!"
const IdMsg = "Resource id：%s "
const FailedGetAttributeMsg = "Getting resource %s attribute by path %s failed!!! Body: %v."

const DefaultDebugMsg = "\n*************** %s Response *************** \n%s\n%s******************************\n\n"
const FailedToReachTargetStatus = "Failed to reach target status. Current status is %s."
const FailedToReachTargetAttribute = "Failed to reach value for target attribute. Current value is %s."
const RequiredWhenMsg = "attribute '%s' is required when '%s' is %v"
