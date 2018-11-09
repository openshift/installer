package httpcache

import (
	"errors"
	"io"
)

// streamReader fans out a source ReadCloser into two consumers.
//
// https://github.com/golang/go/issues/9051#issue-51289031
type streamReader struct {
	sourceCloser  io.Closer
	teeReader     io.Reader
	pipeReader    *io.PipeReader
	pipeWriter    *io.PipeWriter
	handlerErrors chan error
}

var closedEarly = errors.New("closed before EOF")

// Read pulls data from the source reader.  Besides being returned
// here, it will also be pushed into the pipe for the handler to pick
// up.
func (r *streamReader) Read(data []byte) (n int, err error) {
	n, err = r.teeReader.Read(data)
	if err == io.EOF {
		r.pipeWriter.Close()
		r.pipeWriter = nil
		err = <-r.handlerErrors
		if err != nil {
			return n, err
		}
		return n, io.EOF
	}
	return n, err
}

// Close closes the source reader and the pipe.  The handler reader FIXME
func (r *streamReader) Close() (err error) {
	err = r.sourceCloser.Close()
	if err != nil {
		if r.pipeWriter != nil {
			r.pipeWriter.CloseWithError(err)
		}
		return err
	}

	if r.pipeWriter != nil {
		err = r.pipeWriter.CloseWithError(closedEarly)
		if err != nil {
			return err
		}

		return <-r.handlerErrors
	}

	return nil
}

func stream(readCloser io.ReadCloser, handler func(io.ReadCloser) error) *streamReader {
	pipeReader, pipeWriter := io.Pipe()
	handlerErrors := make(chan error, 1)
	go func() {
		handlerErrors <- handler(pipeReader)
	}()
	return &streamReader{
		sourceCloser:  readCloser,
		teeReader:     io.TeeReader(readCloser, pipeWriter),
		pipeReader:    pipeReader,
		pipeWriter:    pipeWriter,
		handlerErrors: handlerErrors,
	}
}
