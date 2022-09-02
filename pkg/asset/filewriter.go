package asset

// FileWriter interface is used to write all the files in the specified location
type FileWriter interface {
	PersistToFile(directory string) error
}

// NewDefaultFileWriter create a new adapter to expose the default implementation as a FileWriter
func NewDefaultFileWriter(a WritableAsset) FileWriter {
	return &fileWriterAdapter{a: a}
}

type fileWriterAdapter struct {
	a WritableAsset
}

// PersistToFile wraps the default implementation
func (fwa *fileWriterAdapter) PersistToFile(directory string) error {
	return PersistToFile(fwa.a, directory)
}
