package datahub

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "github.com/shopspring/decimal"
    "io"
    "math"
    "reflect"
)

const (
    fieldCountSize         int = 4
    intByteSize            int = 4
    byteSizePerField       int = 8
    binaryRecordHeaderSize int = 16
)

type binaryRecordHeader struct {
    encodeType    int32
    schemaVersion int32
    totalSize     int32
    attrOffset    int32
}

func newBinaryRecordHeader(schemaVersion, recordSize, nextOffset int) *binaryRecordHeader {
    return &binaryRecordHeader{
        encodeType:    1,
        schemaVersion: int32(schemaVersion),
        totalSize:     int32(recordSize),
        attrOffset:    int32(nextOffset),
    }
}

type binaryRecord struct {
    parsedAttr       bool
    fieldStartOffset int
    nextOffset       int
    fieldCount       int
    attrLen          int
    schemaVersion    int
    schema           *RecordSchema
    header           *binaryRecordHeader
    data             []byte
    attributes       map[string]string
}

func newBinaryRecordForDeserialize(buf []byte, header *binaryRecordHeader, schema *RecordSchema) *binaryRecord {
    bRecord := &binaryRecord{
        data:          buf,
        header:        header,
        schema:        schema,
        schemaVersion: int(header.schemaVersion),
        parsedAttr:    false,
        attrLen:       0,
    }

    bRecord.initVariable()
    return bRecord
}

func newBinaryRecordForSerialize(version int, schema *RecordSchema) *binaryRecord {
    bRecord := &binaryRecord{
        schemaVersion: version,
        schema:        schema,
        parsedAttr:    true,
        attrLen:       0,
    }

    bRecord.initVariable()
    minAllocSize := bRecord.getMinAllocSize(bRecord.fieldCount)
    bRecord.data = make([]byte, minAllocSize)
    bRecord.nextOffset = minAllocSize
    return bRecord
}

func (bRecord *binaryRecord) initVariable() {
    if bRecord.schema != nil {
        bRecord.fieldCount = bRecord.schema.Size()
    } else {
        bRecord.fieldCount = 1
    }
    bRecord.fieldStartOffset = bRecord.getFixHeaderLength(bRecord.fieldCount)
}

func (bRecord *binaryRecord) getFixHeaderLength(fieldCount int) int {
    return binaryRecordHeaderSize + fieldCountSize + (((fieldCount + 63) >> 6) << 3)
}

func (bRecord *binaryRecord) getMinAllocSize(fieldCount int) int {
    return bRecord.getFixHeaderLength(fieldCount) + fieldCount*byteSizePerField
}

func (bRecord *binaryRecord) getRecordSize() int {
    return intByteSize + bRecord.attrLen + bRecord.nextOffset
}

func (bRecord *binaryRecord) addAttribute(key string, value string) {
    if bRecord.attributes == nil {
        bRecord.attributes = make(map[string]string)
    }
    bRecord.attributes[key] = value
    bRecord.attrLen += len(key) + len(value) + intByteSize*2
}

func (bRecord *binaryRecord) getAttributes() map[string]string {
    bRecord.parseAttributesIfNeed()
    return bRecord.attributes
}

func (bRecord *binaryRecord) parseAttributesIfNeed() error {
    if !bRecord.parsedAttr {
        if bRecord.header == nil {
            bRecord.header = &binaryRecordHeader{
                encodeType:    0,
                schemaVersion: int32(bRecord.schemaVersion),
                totalSize:     int32(bRecord.getRecordSize()),
                attrOffset:    int32(bRecord.nextOffset),
            }
        }

        offset := bRecord.header.attrOffset
        attrSize := binary.LittleEndian.Uint32(bRecord.data[offset:])
        if attrSize != 0 && bRecord.attributes == nil {
            bRecord.attributes = make(map[string]string, attrSize)
        }

        if uint32(len(bRecord.data)) < attrSize {
            return errors.New("check data len failed")
        }

        offset = offset + 4
        for i := uint32(0); i < attrSize; i = i + 1 {
            keyLen := int32(binary.LittleEndian.Uint32(bRecord.data[offset:]))
            offset = offset + 4
            key := string(bRecord.data[offset : offset+keyLen])
            offset = offset + keyLen

            valueLen := int32(binary.LittleEndian.Uint32(bRecord.data[offset:]))
            offset = offset + 4
            value := string(bRecord.data[offset : offset+valueLen])

            bRecord.attributes[key] = value
            offset = offset + valueLen
        }
        bRecord.parsedAttr = true
    }
    return nil
}

func (bRecord *binaryRecord) getFieldOffset(index int) int {
    return byteSizePerField*index + bRecord.fieldStartOffset
}

func (bRecord *binaryRecord) setNotNullAt(index int) {
    nullOffset := binaryRecordHeaderSize + fieldCountSize + (index >> 3)
    value := bRecord.data[nullOffset] | byte(uint32(1)<<uint32(index&7))
    bRecord.data[nullOffset] = value
}

func (bRecord *binaryRecord) isFieldNull(index int) bool {
    nullOffset := binaryRecordHeaderSize + fieldCountSize + (index >> 3)
    value := bRecord.data[nullOffset] & byte(1<<uint32(index&7))
    return value == 0
}

func (bRecord *binaryRecord) checkFieldIndex(index int) error {
    if index >= bRecord.fieldCount {
        return errors.New(fmt.Sprintf("Filed index: %d exceed field num: %d", index, bRecord.fieldCount))
    }
    return nil
}

func (bRecord *binaryRecord) alignSize(size int) int {
    return (size + 7) & (^7)
}

func (bRecord *binaryRecord) setField(index int, data interface{}) error {
    if err := bRecord.checkFieldIndex(index); err != nil {
        return err
    }

    if data == nil {
        return nil
    }
    bRecord.setNotNullAt(index)

    offset := bRecord.getFieldOffset(index)

    if bRecord.schema != nil {
        field := bRecord.schema.Fields[index]
        switch field.Type {
        case STRING:
            str, ok := data.(String)
            if ! ok {
                return errors.New(fmt.Sprintf("value type [%v] dismatch field type [STRING]", reflect.TypeOf(data)))
            }
            if err := bRecord.writeStr(offset, []byte(str)); err != nil {
                return err
            }
        case DECIMAL:
            val, ok := data.(Decimal)
            if !ok {
                return errors.New(fmt.Sprintf("value type [%v] dismatch field type [DECIMAL]", reflect.TypeOf(data)))
            }
            if err := bRecord.writeStr(offset, []byte(val.String())); err != nil {
                return err
            }
        case BOOLEAN, TINYINT, SMALLINT, INTEGER, BIGINT, TIMESTAMP, FLOAT, DOUBLE:
            val, err := bRecord.convertToUInt64(data)
            if err != nil {
                return err
            }
            binary.LittleEndian.PutUint64(bRecord.data[offset:], val)
        default:
            return errors.New(fmt.Sprintf("Invalid field type [%v]", field.Type))
        }
    } else {
        buf, ok := data.([]byte)
        if !ok {
            return errors.New("only support write byte[] for no schema")
        }
        if err := bRecord.writeStr(offset, buf); err != nil {
            return err
        }
    }

    return nil
}

func (bRecord *binaryRecord) getField(index int) (interface{}, error) {

    if err := bRecord.checkFieldIndex(index); err != nil {
        return nil, err
    }

    if bRecord.isFieldNull(index) {
        return nil, nil
    }

    offset := bRecord.getFieldOffset(index)
    field := bRecord.schema.Fields[index]
    switch field.Type {
    case STRING:
        str := bRecord.readStr(offset)
        return str, nil
    case DECIMAL:
        str := bRecord.readStr(offset)
        return decimal.NewFromString(str)
    case FLOAT:
        bits := binary.LittleEndian.Uint32(bRecord.data[offset:])
        return math.Float32frombits(bits), nil
    case DOUBLE:
        bits := binary.LittleEndian.Uint64(bRecord.data[offset:])
        return math.Float64frombits(bits), nil
    case BOOLEAN:
        val := binary.LittleEndian.Uint64(bRecord.data[offset:])
        return val == 1, nil
    case TINYINT:
        val := binary.LittleEndian.Uint64(bRecord.data[offset:])
        return int8(val), nil
    case SMALLINT:
        val := binary.LittleEndian.Uint64(bRecord.data[offset:])
        return int16(val), nil
    case INTEGER:
        val := binary.LittleEndian.Uint64(bRecord.data[offset:])
        return int32(val), nil
    case BIGINT, TIMESTAMP:
        val := binary.LittleEndian.Uint64(bRecord.data[offset:])
        return int64(val), nil
    default:
        return nil, errors.New(fmt.Sprintf("Invalid field type [%v]", field.Type))
    }
}

func (bRecord *binaryRecord) writeStr(offset int, data []byte) error {
    length := len(data)
    if length < 7 {
        num := copy(bRecord.data[offset:], data)
        for i := num; i < 7; i = i + 1 {
            bRecord.data[offset+i] = byte(0)
        }
        bRecord.data[offset+7] = byte(0x80) | byte(length)
    } else {
        offsetAndSize := ((int64(bRecord.nextOffset - binaryRecordHeaderSize)) << 32) | int64(length)
        binary.LittleEndian.PutUint64(bRecord.data[offset:], uint64(offsetAndSize))

        buf := bytes.NewBuffer(bRecord.data)
        buf.Write(data)
        needSize := bRecord.alignSize(length)
        padSize := needSize - length

        if padSize > 0 {
            for i := 0; i < padSize; i = i + 1 {
                buf.WriteByte(0)
            }
        }
        bRecord.data = buf.Bytes()
        bRecord.nextOffset += needSize
    }
    return nil
}

func (bRecord *binaryRecord) readStr(offset int) string {

    data := binary.LittleEndian.Uint64(bRecord.data[offset:])

    isLittleStr := (data & (uint64(0x80) << 56)) != 0

    if isLittleStr {
        length := int((data >> 56) & 0x07)
        return string(bRecord.data[offset : offset+length])
    } else {
        strOffset := binaryRecordHeaderSize + int(data>>32)
        strLen := int(data & math.MaxUint32)
        return string(bRecord.data[strOffset : strOffset+strLen])
    }
}

// for TINYINT, SMALLINT, INTEGER, BIGINT, TIMESTAMP, FLOAT, DOUBLE
func (bRecord *binaryRecord) convertToUInt64(data interface{}) (uint64, error) {
    var val uint64
    switch v := data.(type) {
    case Tinyint:
        val = uint64(v)
    case Smallint:
        val = uint64(v)
    case Integer:
        val = uint64(v)
    case Bigint:
        val = uint64(v)
    case Timestamp:
        val = uint64(v)
    case Float:
        fVal := float32(v)
        bits := math.Float32bits(fVal)
        val = uint64(bits)
    case Double:
        fVal := float64(v)
        val = math.Float64bits(fVal)
    case Boolean:
        if v {
            val = uint64(1)
        } else {
            val = uint64(0)
        }
    default:
        return 0, errors.New(fmt.Sprintf("value type[%T] not match field type", reflect.ValueOf(val)))
    }
    return val, nil
}

type binaryRecordContextSerializer struct {
    projectName  string
    topicName    string
    shardId      string
    schema       *RecordSchema
    schemaClient *schemaRegistryClient
}

func (serializer *binaryRecordContextSerializer) serializeRecordHeader(bHeader *binaryRecordHeader) []byte {
    buf := make([]byte, binaryRecordHeaderSize)
    binary.LittleEndian.PutUint32(buf, uint32(bHeader.encodeType))
    binary.LittleEndian.PutUint32(buf[4:], uint32(bHeader.schemaVersion))
    binary.LittleEndian.PutUint32(buf[8:], uint32(bHeader.totalSize))
    binary.LittleEndian.PutUint32(buf[12:], uint32(bHeader.attrOffset))
    return buf
}

// BinaryRecord => []byte
func (serializer *binaryRecordContextSerializer) serializeBinaryRecord(writer *bytes.Buffer, bRecord *binaryRecord) error {
    if bRecord.header == nil {
        bRecord.header = newBinaryRecordHeader(bRecord.schemaVersion, bRecord.getRecordSize(), bRecord.nextOffset)
    }
    headerBuf := serializer.serializeRecordHeader(bRecord.header)
    copy(bRecord.data, headerBuf)

    _, err := writer.Write(bRecord.data)
    if err != nil {
        return err
    }

    if err = binary.Write(writer, binary.LittleEndian, int32(len(bRecord.attributes))); err != nil {
        return err
    }
    for key, val := range bRecord.attributes {
        if err = binary.Write(writer, binary.LittleEndian, int32(len(key))); err != nil {
            return err
        }
        if _, err = writer.WriteString(key); err != nil {
            return err
        }
        if err = binary.Write(writer, binary.LittleEndian, int32(len(val))); err != nil {
            return err
        }
        if _, err = writer.WriteString(val); err != nil {
            return err
        }
    }

    return nil
}

func (serializer *binaryRecordContextSerializer) deserializeRecordHeader(reader *bytes.Reader) (*binaryRecordHeader, error) {
    if reader.Len() < binaryRecordHeaderSize {
        return nil, errors.New(fmt.Sprintf("Data length is not enough for BinaryRecordHeader(%d)", binaryRecordHeaderSize))
    }

    header := &binaryRecordHeader{}
    if err := binary.Read(reader, binary.LittleEndian, &header.encodeType); err != nil {
        return nil, err
    }
    if err := binary.Read(reader, binary.LittleEndian, &header.schemaVersion); err != nil {
        return nil, err
    }
    if err := binary.Read(reader, binary.LittleEndian, &header.totalSize); err != nil {
        return nil, err
    }
    if err := binary.Read(reader, binary.LittleEndian, &header.attrOffset); err != nil {
        return nil, err
    }

    return header, nil
}

// []byte => BinaryRecord
func (serializer *binaryRecordContextSerializer) deserializeBinaryRecord(reader *bytes.Reader) (*binaryRecord, error) {
    bHeader, err := serializer.deserializeRecordHeader(reader)
    if err != nil {
        return nil, err
    }

    // 读取header完成之后重置到读header之前的位点
    if _, err := reader.Seek(-int64(binaryRecordHeaderSize), io.SeekCurrent); err != nil {
        return nil, err
    }

    if reader.Len() < int(bHeader.totalSize) {
        return nil, errors.New(fmt.Sprintf("Check record header length fail, need: %d, real: %d", bHeader.totalSize, reader.Len()))
    }

    var schema *RecordSchema = nil
    if bHeader.schemaVersion != -1 {
        if serializer.schemaClient != nil {
            schema, err = serializer.getSchemeByVersion(int(bHeader.schemaVersion))
            if err != nil {
                return nil, err
            }
        } else {
            schema = serializer.schema
        }
    }

    buf := make([]byte, bHeader.totalSize)
    if _, err = reader.Read(buf); err != nil {
        return nil, err
    }
    return newBinaryRecordForDeserialize(buf, bHeader, schema), nil
}

func (serializer *binaryRecordContextSerializer) getSchemeByVersion(version int) (*RecordSchema, error) {
    if serializer.schemaClient != nil {
        return serializer.schemaClient.getSchemaByVersion(serializer.projectName, serializer.topicName, version)
    }
    return nil, nil
}

func (serializer *binaryRecordContextSerializer) getVersionBySchema(schema *RecordSchema) (int, error) {
    if serializer.schemaClient != nil {
        return serializer.schemaClient.getVersionBySchema(serializer.projectName, serializer.topicName, schema)
    }
    return 0, nil
}

func (serializer *binaryRecordContextSerializer) blob2BinaryRecord(record *BlobRecord) (*binaryRecord, error) {
    bRecord := newBinaryRecordForSerialize(-1, nil)
    if err := bRecord.setField(0, record.RawData); err != nil {
        return nil, err
    }
    return bRecord, nil
}

func (serializer *binaryRecordContextSerializer) tuple2BinaryRecord(record *TupleRecord) (*binaryRecord, error) {
    version := 0
    if serializer.schemaClient != nil {
        val, err := serializer.getVersionBySchema(record.RecordSchema)
        if err != nil {
            return nil, err
        }
        version = val
    }

    bRecord := newBinaryRecordForSerialize(version, record.RecordSchema)
    for idx, val := range record.Values {
        if err := bRecord.setField(idx, val); err != nil {
            return nil, err
        }
    }
    return bRecord, nil
}

func (serializer *binaryRecordContextSerializer) dhRecord2BinaryRecord(record IRecord) (*binaryRecord, error) {

    var err error
    var bRecord *binaryRecord = nil

    switch record.(type) {
    case *TupleRecord:
        bRecord, err = serializer.tuple2BinaryRecord(record.(*TupleRecord))
        if err != nil {
            return nil, err
        }
    case *BlobRecord:
        bRecord, err = serializer.blob2BinaryRecord(record.(*BlobRecord))
        if err != nil {
            return nil, err
        }
    default:
        return nil, errors.New(fmt.Sprintf("Invalid record type %v", reflect.TypeOf(record)))
    }

    attributes := record.GetAttributes()
    if attributes != nil {
        for key, val := range attributes {
            strVal, ok := val.(string)
            if !ok {
                return nil, errors.New("attribute only support map[string]string now")
            }
            bRecord.addAttribute(key, strVal)
        }
    }
    return bRecord, nil
}

func (serializer *binaryRecordContextSerializer) binaryRecord2DhRecord(bRecord *binaryRecord, meta *respMeta, schema *RecordSchema) (IRecord, error) {
    var record IRecord
    var err error

    if schema != nil {
        record, err = serializer.binary2TupleRecord(bRecord, schema)
        if err != nil {
            return nil, err
        }
    } else {
        record, err = serializer.binary2BlobRecord(bRecord)
        if err != nil {
            return nil, err
        }
    }

    baseRecord := BaseRecord{
        ShardId:    serializer.shardId,
        SystemTime: meta.systemTime,
        Sequence:   meta.sequence,
        Cursor:     meta.cursor,
        NextCursor: meta.nextCursor,
        Serial:     meta.serial,
    }
    attributes := bRecord.getAttributes()
    for key, val := range attributes {
        baseRecord.SetAttribute(key, val)
    }
    record.SetBaseRecord(baseRecord)

    return record, nil
}

func (serializer *binaryRecordContextSerializer) binary2TupleRecord(bRecord *binaryRecord, schema *RecordSchema) (*TupleRecord, error) {
    record := NewTupleRecord(schema, 0)
    for idx := 0; idx < schema.Size(); idx = idx + 1 {
        val, err := bRecord.getField(idx)
        if err != nil {
            return nil, err
        }
        if val != nil {
            record.SetValueByIdx(idx, val)
        }
    }
    return record, nil
}

func (serializer *binaryRecordContextSerializer) binary2BlobRecord(bRecord *binaryRecord) (*BlobRecord, error) {
    val, err := bRecord.getField(0)
    if err != nil {
        return nil, err
    }
    data, ok := val.([]byte)
    if !ok {
        return nil, errors.New("only support write byte[] for no schema")
    }
    return NewBlobRecord(data, 0), nil
}
