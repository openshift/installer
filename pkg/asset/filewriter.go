package asset

// Asset interface used to write all the files in the specified location
type FileWriter interface {
	PersistToFile(directory string) error
}

// Create a new adapter to expose the default implementation as a FileWriter
func NewDefaultFileWriter(a WritableAsset) FileWriter {
	return &fileWriterAdapter{a: a}
}

type fileWriterAdapter struct {
	a WritableAsset
}

func (fwa *fileWriterAdapter) PersistToFile(directory string) error {
	return PersistToFile(fwa.a, directory)
}
