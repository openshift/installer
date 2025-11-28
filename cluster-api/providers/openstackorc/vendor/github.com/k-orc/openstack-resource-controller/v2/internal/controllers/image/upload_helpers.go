/*
Copyright 2024 The ORC Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package image

import (
	"compress/bzip2"
	"compress/gzip"
	"crypto/md5"  //nolint:gosec
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"

	"github.com/ulikunitz/xz"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

type progressReporter func(progress int64)

type readerWithProgress struct {
	reader   io.Reader
	progress int64
	reporter progressReporter
}

var _ io.Reader = &readerWithProgress{}

func (r *readerWithProgress) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.progress += int64(n)
	r.reporter(r.progress)
	return n, err
}

func newReaderWithProgress(reader io.Reader, reporter progressReporter) io.Reader {
	return &readerWithProgress{
		reader:   reader,
		reporter: reporter,
	}
}

type hashCompletionHandler func(string) error

type readerWithHash struct {
	algorithm         string
	reader            io.Reader
	hasher            hash.Hash
	completionHandler hashCompletionHandler
}

func (r *readerWithHash) Read(p []byte) (int, error) {
	readerN, readerErr := r.reader.Read(p)
	hasherN, _ := r.hasher.Write(p[:readerN])

	if readerN != hasherN {
		return readerN, errors.Join(fmt.Errorf("hasher did not consume all read bytes"), readerErr)
	}

	if errors.Is(readerErr, io.EOF) {
		hashValue := fmt.Sprintf("%x", r.hasher.Sum(nil))
		if err := r.completionHandler(hashValue); err != nil {
			return readerN, err
		}
	}

	return readerN, readerErr
}

var _ io.Reader = &readerWithHash{}

func newReaderWithHash(reader io.Reader, algorithm orcv1alpha1.ImageHashAlgorithm, completionHandler hashCompletionHandler) (io.Reader, error) {
	hasher := getHasher(algorithm)
	if hasher == nil {
		return nil, errors.New("no registered hash algorithm for " + string(algorithm))
	}
	return &readerWithHash{
		algorithm:         string(algorithm),
		reader:            reader,
		hasher:            hasher,
		completionHandler: completionHandler,
	}, nil
}

func getHasher(algorithm orcv1alpha1.ImageHashAlgorithm) hash.Hash {
	switch algorithm {
	case orcv1alpha1.ImageHashAlgorithmMD5:
		return md5.New() //nolint:gosec
	case orcv1alpha1.ImageHashAlgorithmSHA1:
		return sha1.New() //nolint:gosec
	case orcv1alpha1.ImageHashAlgorithmSHA256:
		return sha256.New()
	case orcv1alpha1.ImageHashAlgorithmSHA512:
		return sha512.New()
	default:
		// Should have been handled by API validation
		return nil
	}
}

func newReaderWithDecompression(reader io.Reader, compression orcv1alpha1.ImageCompression) (io.Reader, error) {
	switch compression {
	case orcv1alpha1.ImageCompressionXZ:
		reader, err := xz.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("error opening xz compression: %w", err)
		}
		return reader, err
	case orcv1alpha1.ImageCompressionGZ:
		reader, err := gzip.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("error opening gz compression: %w", err)
		}
		return reader, err
	case orcv1alpha1.ImageCompressionBZ2:
		return bzip2.NewReader(reader), nil
	default:
		msg := fmt.Sprintf("unsupported compression algorithm: %s", compression)
		return nil, orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, msg)
	}
}
