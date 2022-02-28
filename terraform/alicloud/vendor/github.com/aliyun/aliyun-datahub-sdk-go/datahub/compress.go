package datahub

import (
    "bytes"
    "compress/zlib"
    "errors"
    "github.com/pierrec/lz4"
    "io"
)

// compress type
type CompressorType string

const (
    NOCOMPRESS CompressorType = ""
    LZ4        CompressorType = "lz4"
    DEFLATE    CompressorType = "deflate"
    ZLIB       CompressorType = "zlib"
)

// validate that the type is valid
func validateCompressorType(ct CompressorType) bool {
    switch ct {
    case NOCOMPRESS, LZ4, DEFLATE, ZLIB:
        return true
    }
    return false
}

func getCompressTypeFromValue(value int) CompressorType {
    switch value {
    case 0:
        return NOCOMPRESS
    case 1:
        return DEFLATE
    case 2:
        return LZ4
    case 3:
        return ZLIB
    default:
        return NOCOMPRESS
    }
}

func (ct *CompressorType) String() string {
    return string(*ct)
}

func (ct *CompressorType) toValue() int {
    switch *ct {
    case NOCOMPRESS:
        return 0
    case DEFLATE:
        return 1
    case LZ4:
        return 2
    case ZLIB:
        return 3
    default:
        return 0
    }
}

// Compressor is a interface for the compress
type compressor interface {
    Compress(data []byte) ([]byte, error)
    DeCompress(data []byte, rawSize int64) ([]byte, error)
}

type lz4Compressor struct {
}

func (lc *lz4Compressor) Compress(data []byte) ([]byte, error) {
    if len(data) == 0 {
        return nil, nil
    }
    buf := make([]byte, len(data))
    ht := make([]int, 64<<10) // buffer for the compression table
    n, err := lz4.CompressBlock(data, buf, ht)
    if err != nil {
        return nil, err
    }
    if n >= len(data) || n == 0 {
        return nil, errors.New("data is not compressible")
    }
    buf = buf[:n] // compressed data
    return buf, nil
}

/*func (lc *Lz4Compressor) Compress(data []byte) ([]byte, error) {
    buf := bytes.NewBuffer(nil)
    writer := lz4.NewWriter(buf)

    defer writer.Close()
    // 写入待压缩内容
    if _, err := writer.Write(data); err != nil {
        return nil, err
    }
    if err := writer.Flush(); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}*/

/*func (lc *Lz4Compressor) DeCompress(data []byte, rawSize int64) ([]byte, error) {
    //get the maximum size of data when not compressible.
    buffer := bytes.NewBuffer(data)
    reader := lz4.NewReader(buffer)
    buf, err := ioutil.ReadAll(reader)

    if err != nil {
        return nil, err
    }
    return buf, nil
}*/

func (lc *lz4Compressor) DeCompress(data []byte, rawSize int64) ([]byte, error) {
    // Allocated a very large buffer for decompression.
    buf := make([]byte, rawSize)
    _, err := lz4.UncompressBlock(data, buf)
    if err != nil {
        return nil, err
    }
    return buf, nil
}

type deflateCompressor struct {
}

/*func (dc *DeflateCompressor) Compress(data []byte) ([]byte, error) {

    // 一个缓存区压缩的内容
    buf := bytes.NewBuffer(nil)
    // 创建一个flate.Writer

    deflateWriter, err := flate.NewWriter(buf, flate.DefaultCompression)
    if err != nil {
        return nil, err
    }
    defer deflateWriter.Close()
    // 写入待压缩内容
    if _, err := deflateWriter.Write(data); err != nil {
        return nil, err
    }
    if err := deflateWriter.Flush(); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil

}*/

/*func (dc *DeflateCompressor)Compress(data []byte) ([]byte, error) {
    var bufs bytes.Buffer
    w, _ := flate.NewWriter(&bufs, flate.DefaultCompression)
    if _, err := w.Write([]byte(data)); err != nil {
        return nil, err
    }
    if err := w.Flush(); err != nil {
        return nil, err
    }
    defer w.Close()

    return bufs.Bytes(), nil
}*/

func (dc *deflateCompressor) Compress(data []byte) ([]byte, error) {
    var buf bytes.Buffer
    w := zlib.NewWriter(&buf)
    if _, err := w.Write(data); err != nil {
        return nil, err
    }
    if err := w.Close(); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

/*func (dc *DeflateCompressor)DeCompress(data []byte,rawSize int64)([]byte,error)  {

    r :=flate.NewReader(bytes.NewReader(data))
    defer r.Close()
    out, err := ioutil.ReadAll(r)
    if err !=nil {
        return nil,err
    }
    return out,nil
    //fmt.Println(out)
}*/

/*func (dc *DeflateCompressor) DeCompress(data []byte, rawSize int64) ([]byte, error) {
    buffer := bytes.NewBuffer(data)
    deflateReader := flate.NewReader(buffer)
    defer deflateReader.Close()

    buf := bytes.NewBuffer(nil)

    if _, err := io.CopyN(buf, deflateReader, rawSize); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}*/

func (dc *deflateCompressor) DeCompress(data []byte, rawSize int64) ([]byte, error) {
    b := bytes.NewReader(data)
    var buf bytes.Buffer
    r, _ := zlib.NewReader(b)
    if _, err := io.Copy(&buf, r); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

type zlibCompressor struct {
}

func (zc *zlibCompressor) Compress(data []byte) ([]byte, error) {
    var buf bytes.Buffer
    w := zlib.NewWriter(&buf)
    if _, err := w.Write(data); err != nil {
        return nil, err
    }
    if err := w.Close(); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func (zc *zlibCompressor) DeCompress(data []byte, rawSize int64) ([]byte, error) {
    b := bytes.NewReader(data)
    var buf bytes.Buffer
    r, _ := zlib.NewReader(b)
    if _, err := io.Copy(&buf, r); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

var compressorMap map[CompressorType]compressor = map[CompressorType]compressor{
    LZ4:     &lz4Compressor{},
    DEFLATE: &deflateCompressor{},
    ZLIB:    &zlibCompressor{},
}

func newCompressor(c CompressorType) compressor {
    switch CompressorType(c) {
    case LZ4:
        return &lz4Compressor{}
    case DEFLATE:
        return &deflateCompressor{}
    case ZLIB:
        return &zlibCompressor{}
    default:
        return nil
    }
}

func getCompressor(c CompressorType) compressor {
    if c == NOCOMPRESS {
        return nil
    }
    ret, ok := compressorMap[c]
    if !ok {
        com := newCompressor(c)
        compressorMap[c] = com
    }
    return ret
}
