package datahub

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/shopspring/decimal"
    "math"
    "strconv"
)

type DataType interface {
    fmt.Stringer
}

// Bigint
type Bigint int64

func (bi Bigint) String() string {
    return strconv.FormatInt(int64(bi), 10)
}

// String
type String string

func (str String) String() string {
    return string(str)
}

// Boolean
type Boolean bool

func (bl Boolean) String() string {
    return strconv.FormatBool(bool(bl))
}

// Double
type Double float64

func (d Double) String() string {
    return strconv.FormatFloat(float64(d), 'f', -1, 64)
}

// Timestamp
type Timestamp uint64

func (t Timestamp) String() string {
    return strconv.FormatUint(uint64(t), 10)
}

// DECIMAL
type Decimal decimal.Decimal

func (d Decimal) String() string {
    return decimal.Decimal(d).String()
}

type Integer int32

func (i Integer) String() string {
    return strconv.FormatInt(int64(i), 10)
}

type Float float32

func (f Float) String() string {
    return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

type Tinyint int8

func (ti Tinyint) String() string {
    return strconv.FormatInt(int64(ti), 10)
}

type Smallint int16

func (si Smallint) String() string {
    return strconv.FormatInt(int64(si), 10)
}

// FieldType
type FieldType string

func (ft FieldType) String() string {
    return string(ft)
}

const (
    // BIGINT 8-bit long signed integer, not include (-9223372036854775808)
    // -9223372036854775807 ~ 9223372036854775807
    BIGINT FieldType = "BIGINT"

    // only support utf-8
    // 1Mb max size
    STRING FieldType = "STRING"

    // BOOLEAN
    // True/False，true/false, 0/1
    BOOLEAN FieldType = "BOOLEAN"

    // DOUBLE 8-bit double
    // -1.0 * 10^308 ~ 1.0 * 10^308
    DOUBLE FieldType = "DOUBLE"

    // TIMESTAMP
    // unit: us
    TIMESTAMP FieldType = "TIMESTAMP"

    // DECIMAL
    // can "only" represent numbers with a maximum of 2^31 digits after the decimal point.
    DECIMAL FieldType = "DECIMAL"

    // 4-byte signed integer
    INTEGER FieldType = "INTEGER"

    // Float type
    FLOAT FieldType = "FLOAT"

    // 1-byte signed integer
    TINYINT FieldType = "TINYINT"

    // 2-byte signed integer
    SMALLINT FieldType = "SMALLINT"
)

// validateFieldType validate field type
func validateFieldType(ft FieldType) bool {
    switch ft {
    case BIGINT, STRING, BOOLEAN, DOUBLE, TIMESTAMP, DECIMAL, INTEGER, FLOAT, TINYINT, SMALLINT:
        return true
    default:
        return false
    }
}

func getIntegerValue(val interface{}) (int64, error) {
    var realval int64
    switch v := val.(type) {
    case int:
        realval = int64(v)
    case int8:
        realval = int64(v)
    case int16:
        realval = int64(v)
    case int32:
        realval = int64(v)
    case int64:
        realval = int64(v)
    case uint:
        realval = int64(v)
    case uint8:
        realval = int64(v)
    case uint16:
        realval = int64(v)
    case uint32:
        realval = int64(v)
    case Bigint:
        realval = int64(v)
    case Integer:
        realval = int64(v)
    case Smallint:
        realval = int64(v)
    case Tinyint:
        realval = int64(v)
    case uint64:
        if v > 9223372036854775807 {
            return 0, errors.New("Integer type field must be in [-9223372036854775807,9223372036854775807]")
        }
        realval = int64(v)
    case json.Number:
        nval, err := v.Int64()
        if err != nil {
            return 0, err
        }
        realval = int64(nval)
    default:
        return 0, errors.New(fmt.Sprintf("value type[%T] not match field type", val))
    }
    return realval, nil
}

// validateFieldValue validate field value
func validateFieldValue(ft FieldType, val interface{}) (DataType, error) {
    switch ft {
    case BIGINT:
        realval, err := getIntegerValue(val)
        if err != nil {
            return nil, err
        }

        if int64(realval) < -9223372036854775807 || int64(realval) > 9223372036854775807 {
            return nil, errors.New("BIGINT type field must be in [-9223372036854775807,9223372036854775807]")
        }
        return Bigint(realval), nil
    case STRING:
        var realval String
        switch v := val.(type) {
        case String:
            realval = v
        case string:
            realval = String(v)
        default:
            return nil, errors.New(fmt.Sprintf("value type[%T] not match field type[STRING]", val))
        }
        return realval, nil
    case BOOLEAN:
        switch v := val.(type) {
        case Boolean:
            return v, nil
        case bool:
            return Boolean(v), nil
        default:
            return nil, errors.New(fmt.Sprintf("value type[%T] not match field type[BOOLEAN]", val))
        }
    case DOUBLE:
        switch v := val.(type) {
        case Double:
            return v, nil
        case float64:
            return Double(v), nil
        case json.Number:
            nval, err := v.Float64()
            if err != nil {
                return nil, err
            }
            return Double(nval), nil
        default:
            return nil, errors.New(fmt.Sprintf("value type[%T] not match field type[DOUBLE]", val))
        }
    case TIMESTAMP:
        var realval Timestamp
        switch v := val.(type) {
        case Timestamp:
            realval = v
        case uint:
            realval = Timestamp(v)
        case uint8:
            realval = Timestamp(v)
        case uint16:
            realval = Timestamp(v)
        case uint32:
            realval = Timestamp(v)
        case uint64:
            realval = Timestamp(v)
        case int:
            if v < 0 {
                return nil, errors.New("TIMESTAMP type field must be in positive")
            }
            realval = Timestamp(v)
        case int8:
            if v < 0 {
                return nil, errors.New("TIMESTAMP type field must be in positive")
            }
            realval = Timestamp(v)
        case int16:
            if v < 0 {
                return nil, errors.New("TIMESTAMP type field must be in positive")
            }
            realval = Timestamp(v)
        case int32:
            if v < 0 {
                return nil, errors.New("TIMESTAMP type field must be in positive")
            }
            realval = Timestamp(v)
        case int64:
            if v < 0 {
                return nil, errors.New("TIMESTAMP type field must be in positive")
            }
            realval = Timestamp(v)
        case json.Number:
            nval, err := v.Int64()
            if err != nil {
                return nil, err
            }
            if nval < 0 {
                return nil, errors.New("TIMESTAMP type field must be in positive")
            }
            realval = Timestamp(nval)
        default:
            return nil, errors.New(fmt.Sprintf("value type[%T] not match field type[TIMESTAMP]", val))
        }
        return realval, nil
    case DECIMAL:
        var realval Decimal
        switch v := val.(type) {
        case decimal.Decimal:
            realval = Decimal(v)
        default:
            return nil, errors.New(fmt.Sprintf("value type[%T] not match field type[DECIMAL]", val))
        }
        return realval, nil
    case INTEGER:
        realval, err := getIntegerValue(val)
        if err != nil {
            return nil, err
        }
        if realval > math.MaxInt32 || realval < math.MinInt32 {
            return nil, errors.New(fmt.Sprintf("%T exceed the range of INTEGER", val))
        }
        return Integer(realval), nil
    case FLOAT:
        switch v := val.(type) {
        case Float:
            return v, nil
        case float32:
            return Float(v), nil
        case json.Number:
            nval, err := v.Float64()
            if err != nil {
                return nil, err
            }
            return Float(nval), nil
        default:
            return nil, errors.New(fmt.Sprintf("value type[%T] not match field type[FLOAT]", val))
        }
    case TINYINT:
        realval, err := getIntegerValue(val)
        if err != nil {
            return nil, err
        }
        if realval > math.MaxInt8 || realval < math.MinInt8 {
            return nil, errors.New(fmt.Sprintf("%T exceed the range of TINYINT", val))
        }
        return Tinyint(realval), nil
    case SMALLINT:
        realval, err := getIntegerValue(val)
        if err != nil {
            return nil, err
        }
        if realval > math.MaxInt16 || realval < math.MinInt16 {
            return nil, errors.New(fmt.Sprintf("%T exceed the range of TINYINT", val))
        }
        return Smallint(realval), nil
    default:
        return nil, errors.New(fmt.Sprintf("field type[%T] is not illegal", ft))
    }
}

// CastValueFromString cast value from string
func castValueFromString(str string, ft FieldType) (DataType, error) {
    switch ft {
    case BIGINT:
        v, err := strconv.ParseInt(str, 10, 64)
        if err == nil {
            return Bigint(v), nil
        }
        return nil, err
    case STRING:
        return String(str), nil
    case BOOLEAN:
        v, err := strconv.ParseBool(str)
        if err == nil {
            return Boolean(v), nil
        }
        return nil, err
    case DOUBLE:
        v, err := strconv.ParseFloat(str, 64)
        if err == nil {
            return Double(v), nil
        }
        return nil, err
    case TIMESTAMP:
        v, err := strconv.ParseUint(str, 10, 64)
        if err == nil {
            return Timestamp(v), nil
        }
        return nil, err
    case DECIMAL:
        v, err := decimal.NewFromString(str)
        if err == nil {
            return Decimal(v), nil
        }
        return nil, err
    case INTEGER:
        v, err := strconv.ParseInt(str, 10, 32)
        if err == nil {
            return Integer(v), nil
        }
        return nil, err
    case FLOAT:
        v, err := strconv.ParseFloat(str, 32)
        if err == nil {
            return Float(v), nil
        }
        return nil, err
    case TINYINT:
        v, err := strconv.ParseInt(str, 10, 32)
        if err == nil {
            return Tinyint(v), nil
        }
        return nil, err
    case SMALLINT:
        v, err := strconv.ParseInt(str, 10, 32)
        if err == nil {
            return Smallint(v), nil
        }
        return nil, err
    default:
        return nil, errors.New(fmt.Sprintf("not support field type %s", string(ft)))
    }
}

// RecordType
type RecordType string

func (rt RecordType) String() string {
    return string(rt)
}

const (
    // BLOB record
    BLOB RecordType = "BLOB"

    // TUPLE record
    TUPLE RecordType = "TUPLE"
)

// validateRecordType validate record type
func validateRecordType(rt RecordType) bool {
    switch rt {
    case BLOB, TUPLE:
        return true
    default:
        return false
    }
}

type TopicStatus string

func (ts TopicStatus) String() string {
    return string(ts)
}

const (
    TOPIC_ON  TopicStatus = "on"
    TOPIC_OFF TopicStatus = "off"
)

type ExpandMode string

func (ft ExpandMode) String() string {
    return string(ft)
}

const (
    SPLIT_EXTEND ExpandMode = ""
    ONLY_SPLIT   ExpandMode = "split"
    ONLY_EXTEND  ExpandMode = "extend"
)

// ShardState
type ShardState string

func (state ShardState) String() string {
    return string(state)
}

const (
    // OPENING shard is creating or fail over, not available
    OPENING ShardState = "OPENING"

    // ACTIVE is available
    ACTIVE ShardState = "ACTIVE"

    // CLOSED read-only
    CLOSED ShardState = "CLOSED"

    // CLOSING shard is closing, not available
    CLOSING ShardState = "CLOSING"
)

// CursorType
type CursorType string

func (ct CursorType) String() string {
    return string(ct)
}

const (
    // OLDEST
    OLDEST CursorType = "OLDEST"

    // LATEST
    LATEST CursorType = "LATEST"

    // SYSTEM_TIME point to first record after system_time
    SYSTEM_TIME CursorType = "SYSTEM_TIME"

    // SEQUENCE point to the specified sequence
    SEQUENCE CursorType = "SEQUENCE"
)

// validateCursorType validate field type
func validateCursorType(ct CursorType) bool {
    switch ct {
    case OLDEST, LATEST, SYSTEM_TIME, SEQUENCE:
        return true
    default:
        return false
    }
}

// SubscriptionType
type SubscriptionType int

const (
    // SUBTYPE_USER
    SUBTYPE_USER SubscriptionType = iota

    // SUBTYPE_SYSTEM
    SUBTYPE_SYSTEM

    // SUBTYPE_TT
    SUBTYPE_TT
)

func (subType SubscriptionType) Value() int {
    return int(subType)
}

// SubscriptionState
type SubscriptionState int

const (
    // SUB_OFFLINE
    SUB_OFFLINE SubscriptionState = iota

    // SUB_ONLINE
    SUB_ONLINE
)

func (subState SubscriptionState) Value() int {
    return int(subState)
}

type OffsetAction string

func (oa OffsetAction) String() string {
    return string(oa)
}

//const (
//    OPENSESSION  OffsetAction = "open"
//    GETOFFSET    OffsetAction = "get"
//    COMMITOFFSET OffsetAction = "commit"
//    RESETOFFSET  OffsetAction = "reset"
//)

const (
    maxWaitingTimeInMs = 600000
    minWaitingTimeInMs = 60000
)
