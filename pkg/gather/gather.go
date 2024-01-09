package gather

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	"github.com/openshift/installer/pkg/gather/providers"
)

// New returns a Gather based on `metadata.json` in `rootDir`.
func New(logger logrus.FieldLogger, serialLogBundle string, bootstrap string, masters []string, rootDir string) (providers.Gather, error) {
	metadata, err := metadata.Load(rootDir)
	if err != nil {
		return nil, err
	}

	platform := metadata.Platform()
	if platform == "" {
		return nil, errors.New("no platform configured in metadata")
	}

	creator, ok := providers.Registry[platform]
	if !ok {
		return nil, errors.Errorf("no gather methods registered for %q", platform)
	}
	return creator(logger, serialLogBundle, bootstrap, masters, metadata)
}

// CreateArchive creates a gzipped tar file.
func CreateArchive(files []string, archiveName string) error {
	file, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filename := range files {
		err := addToArchive(tarWriter, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func addToArchive(tarWriter *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	st, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(st, st.Name())
	if err != nil {
		return err
	}

	header.Name = filename
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}

	return nil
}

// CombineArchives creates a single gzipped tar file from multiple archives.
// archiveName is the target gzipped tar file. archives maps the existing
// gzipped tar files to a subdirectory in the new gzipped tar file.
func CombineArchives(archiveName string, archives map[string]string) error {
	suffix := ".tar.gz"

	combinedArchive, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer combinedArchive.Close()

	combinedGzipWriter := gzip.NewWriter(combinedArchive)
	defer combinedGzipWriter.Close()

	combinedTarWriter := tar.NewWriter(combinedGzipWriter)
	defer combinedTarWriter.Close()

	combinedDirectory := strings.TrimSuffix(archiveName, suffix)
	if archiveName[0] == '.' || archiveName[0] == '/' {
		combinedDirectory = strings.TrimSuffix(filepath.Base(archiveName), suffix)
	}

	for archive, subDirectory := range archives {
		_, err := os.Stat(archive)
		if err != nil {
			logrus.Warnf("Unable to stat %s, skipping", archive)
			continue
		}

		file, err := os.Open(archive)
		if err != nil {
			return err
		}
		defer file.Close()

		directory := strings.TrimSuffix(archive, suffix) + "/"
		if subDirectory != "" && !strings.HasSuffix(subDirectory, "/") {
			subDirectory += "/"
		}

		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		tarReader := tar.NewReader(gzipReader)

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			newHeaderName := strings.Replace(header.Name, directory, subDirectory, 1)
			// Do not nest `log-bundle-XXXX` directories
			if !strings.HasPrefix(newHeaderName, combinedDirectory) {
				newHeaderName = filepath.Join(combinedDirectory, newHeaderName)
			}
			header.Name = newHeaderName

			err = combinedTarWriter.WriteHeader(header)
			if err != nil {
				return err
			}

			_, err = io.Copy(combinedTarWriter, tarReader)
			if err != nil {
				return err
			}
		}

		// The files are now part of the combined archive, so clean it up
		if err := os.Remove(archive); err != nil {
			logrus.Warnf("Could not remove %s: %v\n", archive, err)
		}
	}

	return nil
}

// DeleteArchiveDirectory deletes an archive directory
func DeleteArchiveDirectory(archiveDirectory string) error {
	if archiveDirectory == "" {
		return nil
	}

	_, err := os.Stat(archiveDirectory)
	if err == nil && !strings.HasPrefix(archiveDirectory, ".") && archiveDirectory != "/" {
		err := os.RemoveAll(archiveDirectory)
		if err != nil {
			return err
		}
	}

	return nil
}
