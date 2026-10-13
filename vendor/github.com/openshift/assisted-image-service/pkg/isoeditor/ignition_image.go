package isoeditor

import (
	"bytes"
	"encoding/json"
	"io"
)

// NewIgnitionImageReader returns the filename of the ignition image in the ISO,
// along with a stream of the ignition image with ignition content embedded.
// This can be used to overwrite the ignition image file of an ISO previously
// unpacked by Extract() in order to embed ignition data.
func NewIgnitionImageReader(isoPath string, ignitionContent *IgnitionContent) ([]FileData, error) {
	info, iso, err := ignitionOverlay(isoPath, ignitionContent, true)
	if err != nil {
		return nil, err
	}

	minLength := info.Offset + info.Length
	ignitionImage, expanded, err := isolateISOFile(isoPath, info.File, iso, minLength)
	if err != nil {
		iso.Close()
		return nil, err
	}
	output := []FileData{ignitionImage}

	// output updated igninfo.json if we have expanded the embed area
	if expanded {
		if _, _, err := GetISOFileInfo(ignitionInfoPath, isoPath); err == nil {
			if ignitionInfoData, err := json.Marshal(info); err == nil {
				output = append(output, FileData{
					Filename: ignitionInfoPath,
					Data:     io.NopCloser(bytes.NewReader(ignitionInfoData)),
				})
			} else {
				iso.Close()
				return nil, err
			}
		}
	}

	return output, nil
}
