package isoeditor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/openshift/assisted-image-service/pkg/overlay"
)

const (
	defaultGrubFilePath     = "/EFI/redhat/grub.cfg"
	defaultIsolinuxFilePath = "/isolinux/isolinux.cfg"
	kargsConfigFilePath     = "/coreos/kargs.json"
)

type FileReader func(isoPath, filePath string) ([]byte, error)

func kargsFiles(isoPath string, fileReader FileReader) ([]string, error) {
	kargsData, err := fileReader(isoPath, kargsConfigFilePath)
	if err != nil {
		// If the kargs file is not found, it is probably iso for old iso version which the file does not exist.  Therefore,
		// default is returned
		return []string{defaultGrubFilePath, defaultIsolinuxFilePath}, nil
	}
	var kargsConfig struct {
		Files []struct {
			Path *string
		}
	}
	if err := json.Unmarshal(kargsData, &kargsConfig); err != nil {
		return nil, err
	}
	var ret []string
	for _, file := range kargsConfig.Files {
		if file.Path != nil {
			ret = append(ret, *file.Path)
		}
	}
	return ret, nil
}

func KargsFiles(isoPath string) ([]string, error) {
	return kargsFiles(isoPath, ReadFileFromISO)
}

func kargsEmbedAreaBoundariesFinder(isoPath, filePath string, fileBoundariesFinder BoundariesFinder, fileReader FileReader) (int64, int64, error) {
	start, _, err := fileBoundariesFinder(filePath, isoPath)
	if err != nil {
		return 0, 0, err
	}

	b, err := fileReader(isoPath, filePath)
	if err != nil {
		return 0, 0, err
	}

	re := regexp.MustCompile(`(\n#*)# COREOS_KARG_EMBED_AREA`)
	submatchIndexes := re.FindSubmatchIndex(b)
	if len(submatchIndexes) != 4 {
		return 0, 0, errors.New("failed to find COREOS_KARG_EMBED_AREA")
	}
	return start + int64(submatchIndexes[2]), int64(submatchIndexes[3] - submatchIndexes[2]), nil
}

func createKargsEmbedAreaBoundariesFinder() BoundariesFinder {
	return func(filePath, isoPath string) (int64, int64, error) {
		return kargsEmbedAreaBoundariesFinder(isoPath, filePath, GetISOFileInfo, ReadFileFromISO)
	}
}

func readerForKargsContent(isoPath string, filePath string, base io.ReadSeeker, contentReader *bytes.Reader) (overlay.OverlayReader, error) {
	return readerForContent(isoPath, filePath, base, contentReader, createKargsEmbedAreaBoundariesFinder())
}

type kernelArgument struct {
	// The operation to apply on the kernel argument.
	// Enum: [append replace delete]
	Operation string `json:"operation,omitempty"`

	// Kernel argument can have the form <parameter> or <parameter>=<value>. The following examples should
	// be supported:
	// rd.net.timeout.carrier=60
	// isolcpus=1,2,10-20,100-2000:2/25
	// quiet
	// The parsing by the command line parser in linux kernel is much looser and this pattern follows it.
	Value string `json:"value,omitempty"`
}

type kernelArguments []*kernelArgument

func KargsToStr(args []string) (string, error) {
	var kargs kernelArguments
	for _, s := range args {
		kargs = append(kargs, &kernelArgument{
			Operation: "append",
			Value:     s,
		})
	}
	b, err := json.Marshal(&kargs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal kernel arguments %v", err)
	}
	return string(b), nil
}

func StrToKargs(kargsStr string) ([]string, error) {
	var kargs kernelArguments
	if err := json.Unmarshal([]byte(kargsStr), &kargs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kernel arguments %v", err)
	}
	var args []string
	for _, arg := range kargs {
		if arg.Operation != "append" {
			return nil, fmt.Errorf("only 'append' operation is allowed.  got %s", arg.Operation)
		}
		args = append(args, arg.Value)
	}
	return args, nil
}
