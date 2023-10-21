package datahub

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "github.com/pkg/errors"
    "hash/crc32"
)

const (
    batchRecordHeaderSize int = 26
)

var (
    batchMagicBytes = []byte{'D', 'H', 'U', 'B'}
    batchMagicNum   = int(binary.LittleEndian.Uint32(batchMagicBytes))
)

func calculateCrc32(buf []byte) uint32 {
    table := crc32.MakeTable(crc32.Castagnoli)
    return crc32.Checksum(buf, table)
}

type respMeta struct {
    cursor     string
    nextCursor string
    sequence   int64
    systemTime int64
    serial     int64
}

type batchRecordHeader struct {
    magic       int
    version     int
    length      int
    rawSize     int
    crc32       uint32
    attributes  int16
    recordCount int
}

type batchRecord struct {
    header  *batchRecordHeader
    records [] *binaryRecord
}

type batchSerializer struct {
    cType        CompressorType
    schemaClient *schemaRegistryClient
    bSerializer  *binaryRecordContextSerializer
}

func newBatchSerializer(project, topic string, cType CompressorType, schemaClient *schemaRegistryClient) *batchSerializer {
    tmpSerializer := &binaryRecordContextSerializer{
        projectName:  project,
        topicName:    topic,
        schemaClient: schemaClient,
    }

    return &batchSerializer{
        cType:        cType,
        schemaClient: schemaClient,
        bSerializer:  tmpSerializer,
    }
}

func (serializer *batchSerializer) serialize(records []IRecord) ([]byte, error) {
    batch, err := serializer.parseBatchRecord(records)
    if err != nil {
        return nil, err
    }
    return serializer.serializeBatchRecord(batch)
}

// dh record list => batch record
func (serializer *batchSerializer) parseBatchRecord(records []IRecord) (*batchRecord, error) {
    batch := &batchRecord{
        records: make([]*binaryRecord, 0, len(records)),
    }

    for _, record := range records {
        bRecord, err := serializer.bSerializer.dhRecord2BinaryRecord(record)
        if err != nil {
            return nil, err
        }
        batch.records = append(batch.records, bRecord)
    }
    return batch, nil
}

func (serializer *batchSerializer) serializeBatchHeader(bHeader *batchRecordHeader) []byte {
    buf := make([]byte, batchRecordHeaderSize)
    copy(buf, batchMagicBytes)
    offset := len(batchMagicBytes)
    binary.LittleEndian.PutUint32(buf[offset:], uint32(bHeader.version))
    binary.LittleEndian.PutUint32(buf[offset+4:], uint32(bHeader.length))
    binary.LittleEndian.PutUint32(buf[offset+8:], uint32(bHeader.rawSize))
    binary.LittleEndian.PutUint32(buf[offset+12:], uint32(bHeader.crc32))
    binary.LittleEndian.PutUint16(buf[offset+16:], uint16(bHeader.attributes))
    binary.LittleEndian.PutUint32(buf[offset+18:], uint32(bHeader.recordCount))
    return buf
}

func (serializer *batchSerializer) serializeBatchRecord(batch *batchRecord) ([]byte, error) {
    calSize := batchRecordHeaderSize
    for _, bRecord := range batch.records {
        calSize += bRecord.getRecordSize()
    }

    writer := &bytes.Buffer{}
    writer.Grow(calSize)
    writer.Write(make([]byte, batchRecordHeaderSize))

    for _, bRecord := range batch.records {
        if err := serializer.bSerializer.serializeBinaryRecord(writer, bRecord); err != nil {
            return nil, err
        }
    }

    data := writer.Bytes()

    if batch.header == nil {
        batch.header = &batchRecordHeader{}
    }

    batch.header.magic = int(batchMagicNum)
    batch.header.version = 0
    batch.header.rawSize = len(data) - batchRecordHeaderSize
    batch.header.length = len(data)
    batch.header.attributes = int16(serializer.cType.toValue() & 3)
    batch.header.recordCount = len(batch.records)

    data, err := serializer.compressIfNeed(data, batch)
    if err != nil {
        return nil, err
    }

    batch.header.crc32 = calculateCrc32(data[batchRecordHeaderSize:])
    copy(data, serializer.serializeBatchHeader(batch.header))
    return data, nil
}

func (serializer *batchSerializer) compressIfNeed(data []byte, batch *batchRecord) ([]byte, error) {
    buf := data
    compressor := getCompressor(serializer.cType)
    if compressor != nil {
        cBuf, err := compressor.Compress(data[batchRecordHeaderSize:])
        if err != nil {
            return nil, err
        }

        buf = append(data[0:batchRecordHeaderSize], cBuf...)
        batch.header.length = len(buf)
    }
    return buf, nil
}

type batchDeserializer struct {
    shardId      string
    schemaClient *schemaRegistryClient
    bSerializer  *binaryRecordContextSerializer
}

func newBatchDeserializer(project, topic, shardId string, schema *RecordSchema, schemaClient *schemaRegistryClient) *batchDeserializer {
    tmpSerializer := &binaryRecordContextSerializer{
        projectName:  project,
        topicName:    topic,
        shardId:      shardId,
        schema:       schema,
        schemaClient: schemaClient,
    }

    return &batchDeserializer{
        schemaClient: schemaClient,
        bSerializer:  tmpSerializer,
    }
}

func (deserializer *batchDeserializer) deserialize(data []byte, meta *respMeta) ([]IRecord, error) {
    batch, err := deserializer.parseBatchRecord(data)
    if err != nil {
        return nil, err
    }

    return deserializer.deserializeBatchRecord(batch, meta)
}

func (deserializer *batchDeserializer) deserializeBatchHeader(data []byte) (*batchRecordHeader, error) {
    if len(data) < batchRecordHeaderSize {
        return nil, errors.New("read batch header fail")
    }

    header := &batchRecordHeader{}
    header.magic = int(binary.LittleEndian.Uint32(data[0:]))
    header.version = int(binary.LittleEndian.Uint32(data[4:]))
    header.length = int(binary.LittleEndian.Uint32(data[8:]))
    header.rawSize = int(binary.LittleEndian.Uint32(data[12:]))
    header.crc32 = binary.LittleEndian.Uint32(data[16:])
    header.attributes = int16(binary.LittleEndian.Uint16(data[20:]))
    header.recordCount = int(binary.LittleEndian.Uint32(data[22:]))

    if header.magic != batchMagicNum {
        return nil, errors.New("Check magic number fail")
    }

    if header.length != len(data) {
        return nil, errors.New("Check payload length fail")
    }

    if header.crc32 != 0 {
        calCrc := calculateCrc32(data[batchRecordHeaderSize:header.length])
        if calCrc != header.crc32 {
            return nil, errors.New(fmt.Sprintf("Check crc fail. expect:%d, real:%d", header.crc32, calCrc))
        }
    }

    return header, nil
}

// []byte => batch record
func (deserializer *batchDeserializer) parseBatchRecord(data []byte) (*batchRecord, error) {
    batchHeader, err := deserializer.deserializeBatchHeader(data)
    if err != nil {
        return nil, err
    }

    // 跳过batch header的部分
    data = data[batchRecordHeaderSize:]

    rawBuf, err := deserializer.decompressIfNeed(batchHeader, data)
    if err != nil {
        return nil, err
    }

    batch := &batchRecord{}
    reader := bytes.NewReader(rawBuf)
    batch.records = make([]*binaryRecord, 0, batchHeader.recordCount)
    for idx := 0; idx < batchHeader.recordCount; idx = idx + 1 {
        bRecord, err := deserializer.bSerializer.deserializeBinaryRecord(reader)
        if err != nil {
            return nil, err
        }
        batch.records = append(batch.records, bRecord)
    }

    return batch, nil
}

func (deserializer *batchDeserializer) deserializeBatchRecord(batch *batchRecord, meta *respMeta) ([]IRecord, error) {
    recordList := make([]IRecord, 0, len(batch.records))
    for _, bRecord := range batch.records {
        record, err := deserializer.bSerializer.binaryRecord2DhRecord(bRecord, meta, bRecord.schema)

        if err != nil {
            return nil, err
        }
        recordList = append(recordList, record)
    }

    return recordList, nil
}

func (deserializer *batchDeserializer) decompressIfNeed(header *batchRecordHeader, data []byte) ([]byte, error) {
    cType := getCompressTypeFromValue(int(header.attributes) & 3)
    compressor := getCompressor(cType)
    if compressor != nil {
        return compressor.DeCompress(data, int64(header.rawSize))
    }
    return data, nil
}
