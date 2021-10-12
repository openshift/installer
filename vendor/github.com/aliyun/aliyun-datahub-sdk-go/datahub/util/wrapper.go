package util

import (
    "encoding/binary"
    "errors"
    "fmt"
    "hash/crc32"
    "reflect"
)

func WrapMessage(data []byte) []byte {
    h1 := []byte("DHUB")

    crc32c := crc32.MakeTable(crc32.Castagnoli)
    crc := crc32.Checksum(data, crc32c)
    h2 := make([]byte, 4)
    binary.BigEndian.PutUint32(h2, crc)

    h3 := make([]byte, 4)
    binary.BigEndian.PutUint32(h3, uint32(len(data)))

    //buf := make([]byte,0,len(h1)+len(h2)+len(h3)+len(data))
    buf := append(h1, h2...)
    buf = append(buf, h3...)
    buf = append(buf, data...)
    return buf

}

func UnwrapMessage(data []byte) ([]byte, error) {

    crc := data[4:8]

    body := data[12:]
    crc32c := crc32.MakeTable(crc32.Castagnoli)
    cs := crc32.Checksum(body, crc32c)
    computedCrc := make([]byte, 4)
    binary.BigEndian.PutUint32(computedCrc, cs)

    if !reflect.DeepEqual(crc, computedCrc) {
        return nil, errors.New(fmt.Sprintf("Parse pb response body fail, error: crc check error. crc: %s, compute crc: %s", crc, computedCrc))
    }

    return body, nil

}