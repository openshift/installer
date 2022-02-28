package datahub

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "reflect"
)

// BaseRecord
type BaseRecord struct {
    ShardId      string                 `json:"ShardId,omitempty"`
    PartitionKey string                 `json:"PartitionKey,omitempty"`
    HashKey      string                 `json:"HashKey,omitempty"`
    SystemTime   int64                  `json:"SystemTime,omitempty"`
    Sequence     int64                  `json:"Sequence,omitempty"`
    Cursor       string                 `json:"Cursor,omitempty"`
    NextCursor   string                 `json:"NextCursor,omitempty"`
    Serial       int64                  `json:"Serial,omitempty"`
    Attributes   map[string]interface{} `json:"Attributes,omitempty"`
}

func (br *BaseRecord) GetSystemTime() int64 {
    return br.SystemTime
}

func (br *BaseRecord) GetSequence() int64 {
    return br.Sequence
}

// SetAttribute set or modify(if exist) attribute
func (br *BaseRecord) SetAttribute(key string, val interface{}) {
    if br.Attributes == nil {
        br.Attributes = make(map[string]interface{})
    }
    br.Attributes[key] = val
}

func (br *BaseRecord) GetAttributes() map[string]interface{} {
    return br.Attributes
}

//RecordEntry
type RecordEntry struct {
    Data interface{} `json:"Data"`
    BaseRecord
}

// IRecord record interface
type IRecord interface {
    fmt.Stringer
    GetSystemTime() int64
    GetSequence() int64
    GetData() interface{}
    FillData(data interface{}) error
    GetBaseRecord() BaseRecord
    SetBaseRecord(baseRecord BaseRecord)
    SetAttribute(key string, val interface{})
    GetAttributes() map[string]interface{}
}

// BlobRecord blob type record
type BlobRecord struct {
    RawData []byte
    BaseRecord
}

// NewBlobRecord new a tuple type record from given record schema
func NewBlobRecord(bytedata []byte, systemTime int64) *BlobRecord {
    br := &BlobRecord{}
    br.RawData = bytedata
    br.Attributes = make(map[string]interface{})
    br.SystemTime = systemTime
    return br
}

func (br *BlobRecord) String() string {
    record := struct {
        Data       string                 `json:"Data"`
        Attributes map[string]interface{} `json:"Attributes"`
    }{
        Data:       string(br.RawData),
        Attributes: br.Attributes,
    }
    byts, _ := json.Marshal(record)
    return string(byts)
}

// FillData implement of IRecord interface
func (br *BlobRecord) FillData(data interface{}) error {
    switch data.(type) {
    case string:
        bytedata, err := base64.StdEncoding.DecodeString(data.(string))
        if err != nil {
            return err
        }
        br.RawData = bytedata
    case []byte:
        br.RawData = data.([]byte)
    default:
        return errors.New(fmt.Sprintf("invalid data type: %s", reflect.TypeOf(data)))
    }
    return nil
}

// GetData implement of IRecord interface
func (br *BlobRecord) GetData() interface{} {
    return br.RawData
}

// GetBaseRecord get base record enbry
func (br *BlobRecord) GetBaseRecord() BaseRecord {
    return br.BaseRecord
}

func (br *BlobRecord) SetBaseRecord(baseRecord BaseRecord) {
    br.BaseRecord = baseRecord
}

// TupleRecord tuple type record
type TupleRecord struct {
    RecordSchema *RecordSchema
    Values       []DataType
    BaseRecord
}

// NewTupleRecord new a tuple type record from given record schema
func NewTupleRecord(schema *RecordSchema, systemTime int64) *TupleRecord {
    tr := &TupleRecord{}
    if schema != nil {
        tr.RecordSchema = schema
        tr.Values = make([]DataType, schema.Size())
    }
    tr.Attributes = make(map[string]interface{})
    tr.SystemTime = systemTime
    for idx := range tr.Values {
        tr.Values[idx] = nil
    }
    return tr
}

func (tr *TupleRecord) String() string {
    record := struct {
        RecordSchema *RecordSchema          `json:"RecordSchema"`
        Values       []DataType             `json:"Values"`
        Attributes   map[string]interface{} `json:"Attributes"`
    }{
        RecordSchema: tr.RecordSchema,
        Values:       tr.Values,
        Attributes:   tr.Attributes,
    }
    byts, _ := json.Marshal(record)
    return string(byts)
}

// SetValueByIdx set a value by idx
func (tr *TupleRecord) SetValueByIdx(idx int, val interface{}) *TupleRecord {
    if idx < 0 || idx >= tr.RecordSchema.Size() {
        panic(fmt.Sprintf("index[%d] out range", idx))
    }

    field := tr.RecordSchema.Fields[idx]
    if val == nil && !field.AllowNull {
        panic(fmt.Sprintf("[%s] not allow null", field.Name))
    }

    v, err := validateFieldValue(field.Type, val)
    if err != nil {
        panic(err)
    }
    tr.Values[idx] = v
    return tr
}

// SetValueByName set a value by name
func (tr *TupleRecord) SetValueByName(name string, val interface{}) *TupleRecord {
    idx := tr.RecordSchema.GetFieldIndex(name)
    return tr.SetValueByIdx(idx, val)
}

func (tr *TupleRecord) GetValueByIdx(idx int) DataType {
    return tr.Values[idx]
}

func (tr *TupleRecord) GetValueByName(name string) DataType {
    idx := tr.RecordSchema.GetFieldIndex(name)
    return tr.GetValueByIdx(idx)
}

func (tr *TupleRecord) GetValues() map[string]DataType {
    values := make(map[string]DataType)
    for i, f := range tr.RecordSchema.Fields {
        values[f.Name] = tr.Values[i]
    }
    return values
}

// SetValues batch set values
func (tr *TupleRecord) SetValues(values []DataType) *TupleRecord {
    if fsize := tr.RecordSchema.Size(); fsize != len(values) {
        panic(fmt.Sprintf("values size not match field size(field.size=%d, values.size=%d)", fsize, len(values)))
    }
    for idx, val := range values {
        v, err := validateFieldValue(tr.RecordSchema.Fields[idx].Type, val)
        if err != nil {
            panic(err)
        }
        tr.Values[idx] = v
    }
    return tr
}

// FillData implement of IRecord interface
func (tr *TupleRecord) FillData(data interface{}) error {
    datas, ok := data.([]interface{})
    if !ok {
        return errors.New("data must be array")
    }
    //else if fsize := tr.RecordSchema.Size(); len(datas) != fsize {
    //    return errors.New(fmt.Sprintf("data array size not match field size(field.size=%d, values.size=%d)", fsize, len(datas)))
    //}
    for idx, v := range datas {
        if v != nil {
            s, ok := v.(string)
            if !ok {
                return errors.New(fmt.Sprintf("data value type[%T] illegal", v))
            }
            tv, err := castValueFromString(s, tr.RecordSchema.Fields[idx].Type)
            if err != nil {
                return err
            }
            tr.Values[idx] = tv
        }
    }
    return nil
}

// GetData implement of IRecord interface
func (tr *TupleRecord) GetData() interface{} {
    result := make([]interface{}, len(tr.Values))
    for idx, val := range tr.Values {
        if val != nil {
            result[idx] = fmt.Sprintf("%s", val)
        } else {
            result[idx] = nil
        }
    }
    return result
}

// GetBaseRecord get base record entry
func (tr *TupleRecord) GetBaseRecord() BaseRecord {
    return tr.BaseRecord
}

func (tr *TupleRecord) SetBaseRecord(baseRecord BaseRecord) {
    tr.BaseRecord = baseRecord
}

type FailedRecord struct {
    Index        int    `json:"Index"`
    ErrorCode    string `json:"ErrorCode"`
    ErrorMessage string `json:"ErrorMessage"`
}
