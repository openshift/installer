package overlay

import (
	"errors"
	"io"
)

type BaseStream = io.ReadSeeker

type OverlayReader interface {
	BaseStream
	io.ReadSeekCloser
}

type Overlay struct {
	Reader io.ReadSeeker
	Offset int64
	Length int64
}

func (ol Overlay) end() int64 {
	return ol.Offset + ol.Length
}

func (ol Overlay) contains(index int64) bool {
	return ol.Offset <= index && ol.end() > index
}

type overlayReader struct {
	Base    BaseStream
	Overlay Overlay

	readIndex   int64
	totalLength int64
}

func newReader(base BaseStream, overlay Overlay, length int64) (*overlayReader, error) {
	if overlay.end() > length {
		length = overlay.end()
	}

	or := overlayReader{
		Base:        base,
		Overlay:     overlay,
		totalLength: length,
	}

	if _, err := base.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	if _, err := overlay.Reader.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return &or, nil
}

func NewOverlayReader(base BaseStream, overlay Overlay) (OverlayReader, error) {
	length, err := base.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	if overlay.Offset < 0 || overlay.Offset > length {
		return nil, errors.New("Overlay offset is beyond end of base")
	}
	return newReader(base, overlay, length)
}

func NewAppendReader(base BaseStream, reader io.ReadSeeker) (OverlayReader, error) {
	length, err := base.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	appendLength, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	overlay := Overlay{
		Reader: reader,
		Offset: length,
		Length: appendLength,
	}
	return newReader(base, overlay, length)
}

func (or *overlayReader) seek(index int64) (err error) {
	if or.Overlay.contains(index) {
		_, err = or.Overlay.Reader.Seek(index-or.Overlay.Offset, io.SeekStart)
	} else {
		_, err = or.Base.Seek(index, io.SeekStart)
	}
	or.readIndex = index
	return err
}

func (or *overlayReader) Len() int {
	return int(or.totalLength - or.readIndex)
}

func (or *overlayReader) Seek(offset int64, whence int) (int64, error) {
	var start int64
	switch whence {
	case io.SeekStart:
		start = 0
	case io.SeekCurrent:
		start = or.readIndex
	case io.SeekEnd:
		start = or.totalLength
	}

	err := or.seek(start + offset)
	return or.readIndex, err
}

func (or *overlayReader) Read(p []byte) (int, error) {
	if or.readIndex >= or.totalLength {
		return 0, io.EOF
	}

	reader := or.Base
	buffer := p

	overlayBytes := or.Overlay.end() - or.readIndex
	switch {
	case or.Overlay.contains(or.readIndex):
		reader = or.Overlay.Reader
		if int64(len(buffer)) > overlayBytes {
			buffer = p[:overlayBytes]
		}
	case overlayBytes > 0:
		// before the overlay
		baseBytes := or.Overlay.Offset - or.readIndex
		if int64(len(buffer)) > baseBytes {
			buffer = p[:baseBytes]
		}
	default:
		// after the overlay
	}

	bytesRead, readErr := reader.Read(buffer)

	seekErr := or.seek(or.readIndex + int64(bytesRead))

	if readErr != nil && readErr != io.EOF {
		return bytesRead, readErr
	}
	return bytesRead, seekErr
}

func (or *overlayReader) Close() error {
	if closer, hasClose := or.Base.(io.Closer); hasClose {
		return closer.Close()
	}
	return nil
}
