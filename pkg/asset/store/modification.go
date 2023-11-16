package store

import (
	"encoding/json"
	"io/fs"
	"os"
	"time"

	"github.com/openshift/installer/pkg/asset"
)

type modificationChecker interface {
	CheckModified(string, *asset.File) error
}

type fileTime struct {
	time.Time
}

type timestamps map[string]fileTime

type storedData struct {
	ModTime timestamps
}

type empty struct{}

type modificationTracker struct {
	previous  storedData
	current   storedData
	preserved map[string]empty
}

const modTrackerKey = "modificationTracker"

func newModificationTracker() *modificationTracker {
	return &modificationTracker{
		previous:  storedData{timestamps{}},
		current:   storedData{timestamps{}},
		preserved: map[string]empty{},
	}
}

func (m *modificationTracker) CheckModified(path string, file *asset.File) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return m.checkTimestamp(file.Filename, info.ModTime())
}

func (m *modificationTracker) checkTimestamp(name string, modTime time.Time) error {
	if !modTime.After(m.previous.ModTime[name].Time) {
		m.preserved[name] = empty{}
		// The file is unmodified, so tell the Asset that it does not exist
		// on disk. Many assets use os.IsNotExist(err) to check this, so we
		// cannot even wrap it in our own error as we could if we could rely on
		// errors.Is(err, fs.ErrNotExist) being used consistently.
		return fs.ErrNotExist
	}
	m.current.ModTime[name] = fileTime{modTime}
	return nil
}

func (m *modificationTracker) Purge(wa asset.WritableAsset) {
	for _, f := range wa.Files() {
		delete(m.current.ModTime, f.Filename)
	}
}

func (m *modificationTracker) HasPreservedFiles(a asset.Asset) bool {
	if wa, ok := a.(asset.WritableAsset); ok {
		for _, f := range wa.Files() {
			if _, present := m.preserved[f.Filename]; present {
				return true
			}
		}
	}
	return false
}

func (m *modificationTracker) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(m.current, "", "    ")
}

func (m *modificationTracker) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &m.previous); err != nil {
		return err
	}
	return json.Unmarshal(data, &m.current)
}

func (ft fileTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.UnixNano())
}

func (ft *fileTime) UnmarshalJSON(data []byte) error {
	var nsec int64
	err := json.Unmarshal(data, &nsec)
	ft.Time = time.Unix(0, nsec)
	return err
}
