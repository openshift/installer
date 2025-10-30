package tshelpers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cavaliercoder/go-cpio"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/diskfs/go-diskfs"
	"github.com/go-openapi/errors"
	"github.com/pkg/diff"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/vincent-petithory/dataurl"
)

// Command				   | Short description																	   | Usage example
// ------------------------|---------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------
// IgnitionImgContains     | checks if the specified file is stored within /images/ignition.img                    | ignitionImgContains agent.x86_64.iso config.ign
// IsoIgnitionUser 		   | checks if the ignition file extracted from the ISO contains the specified user with   | isoIgnitionUser node.x86_64.iso core my-sshKey
//						   | the given `authKey`.																   |
// IsoIgnitionContains     | checks that the file extracted from the ISO embedded configuration file exists        | isoIgnitionContains node.x86_64.iso /usr/local/bin/add-node.sh
// IsoCmp                  | check that the content of the file extracted from the ISO embedded configuration file | isoCmp agent.x86_64.iso /etc/assisted/manifests/infraenv.yaml expected/infraenv.yaml
//                         | matches the content of the local file                                                 |
// IsoCmpRegEx			   | Same as IsoCmp, but the expected file can contain a regex pattern that will be applied| isoCmpRegEx agent.x86_64.iso /etc/assisted/manifests/infraenv.yaml expected/infraenv.yaml
// 						   | during the comparison.																   |
// IsoFileCmpRegEx 		   | checks that file context extracted directly from the ISO matches the content of the   | isoFileCmpRegEx node.x86_64.iso /EFI/redhat/grub.cfg expected/grub.cfg
//						   | local file, by applying a regex comparison.										   |
// IsoSizeMin 		   | checks that the ISO is greater than a minimum number of megabytes                         | isoSizeMin agent.x86_64.iso 100
// IsoSizeMax 		   | checks that the ISO is less than a maximum number of megabytes                            | isoSizeMax agent.x86_64.iso 1000
// InitrdImgContains       | check if the specified file is stored within a compressed cpio archive by scanning the| initrdImgContains agent.x86_64.iso /agent-files/agent-tui
//                         | content of /images/ignition.img archive in the ISO                                    |
// ------------------------|---------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------
// ConfigImgContains       | check if the specified file is stored within the config image ISO                     | configImgContains agentconfig.noarch.iso /config.gz
// UnconfiguredIgnContains | check if the specified file is stored within the unconfigured ignition Storage Files  | unconfiguredIgnContains /etc/assisted/manifests/infraenv.yaml
// IgnitionStorageContains | check if the specified file is stored within the ignition Storage Files               | - (note: works directly on a ignition file)
// UnconfiguredIgnCmp      | check that the content extracted from the unconfigured ignition configuration file    | unconfiguredIgnCmp /etc/assisted/manifests/infraenv.yaml expected/infraenv.yaml
//                         | matches the content of the local file                                                 |
// IsoContains             | check if the specified is stored within the ISO                                       | isoContains imagebasedconfig.iso /cluster-configuration/manifest.json

// IgnitionImgContains `isoPath` `file` check if the specified file `file`
// is stored within /images/ignition.img archive in the ISO `isoPath` image.
func IgnitionImgContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: ignitionImgContains isoPath file")
	}

	workDir := ts.Getenv("WORK")
	isoPath, eFilePath := args[0], args[1]
	isoPathAbs := filepath.Join(workDir, isoPath)

	_, err := extractArchiveFile(isoPathAbs, "/images/ignition.img", eFilePath)
	ts.Check(err)
}

// IsoIgnitionContains `isoPath` `file` checks that the file
// `isoFile` - extracted from the ISO embedded configuration file
// referenced by `isoPath` - exists.
func IsoIgnitionContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: isoIgnitionContains isoPath")
	}

	workDir := ts.Getenv("WORK")
	isoPath, eFilePath := args[0], args[1]
	isoPathAbs := filepath.Join(workDir, isoPath)

	archiveFile, ignitionFile, err := archiveFileNames(isoPath)
	if err != nil {
		ts.Check(err)
	}

	_, err = readFileFromISO(isoPathAbs, archiveFile, ignitionFile, eFilePath)
	ts.Check(err)
}

// IsoIgnitionUser `isoPath` `user` `authKey` checks if the ignition file extracted
// from the ISO `isoPath` contains the specified `user` with the given `authKey`.
func IsoIgnitionUser(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 3 {
		ts.Fatalf("usage: IsoIgnitionUser isoPath user authKey")
	}

	workDir := ts.Getenv("WORK")
	isoPath, eUser, eAuthKey := args[0], args[1], args[2]
	isoPathAbs := filepath.Join(workDir, isoPath)

	archiveFile, ignitionFile, err := archiveFileNames(isoPath)
	if err != nil {
		ts.Check(err)
	}

	ignition, err := extractIgnition(isoPathAbs, archiveFile, ignitionFile)
	ts.Check(err)

	for _, user := range ignition.Passwd.Users {
		if user.Name != eUser {
			continue
		}
		if len(user.SSHAuthorizedKeys) == 0 || user.SSHAuthorizedKeys[0] != igntypes.SSHAuthorizedKey(eAuthKey) {
			continue
		}
		return
	}
	ts.Fatalf("expected user '%s' with SSH auth key '%s' not found", eUser, eAuthKey)
}

// IsoCmp `isoPath` `isoFile` `expectedFile` check that the content of the file
// `isoFile` - extracted from the ISO embedded configuration file referenced
// by `isoPath` - matches the content of the local file `expectedFile`.
// Environment variables in `expectedFile` are substituted before the comparison.
func IsoCmp(ts *testscript.TestScript, neg bool, args []string) {
	isoCmpInternal(ts, neg, args, byteCompare)
}

// IsoCmpRegEx `isoPath` `isoFile` `expectedFile` works as `IsoCmp`,
// but the expected file can contain a regex pattern that will be
// applied during the comparison.
func IsoCmpRegEx(ts *testscript.TestScript, neg bool, args []string) {
	isoCmpInternal(ts, neg, args, byteCompareRegEx)
}

func isoCmpInternal(ts *testscript.TestScript, neg bool, args []string, cmp func(ts *testscript.TestScript, neg bool, aData, eData []byte, aFilePath, eFilePath string)) {
	if len(args) != 3 {
		ts.Fatalf("usage: isocmp isoPath file1 file2")
	}

	workDir := ts.Getenv("WORK")
	isoPath, aFilePath, eFilePath := args[0], args[1], args[2]
	isoPathAbs := filepath.Join(workDir, isoPath)

	archiveFile, ignitionFile, err := archiveFileNames(isoPath)
	if err != nil {
		ts.Check(err)
	}

	aData, err := readFileFromISO(isoPathAbs, archiveFile, ignitionFile, aFilePath)
	ts.Check(err)

	eFilePathAbs := filepath.Join(workDir, eFilePath)
	eData, err := os.ReadFile(eFilePathAbs)
	ts.Check(err)

	cmp(ts, neg, aData, eData, aFilePath, eFilePath)
}

// IsoFileCmpRegEx `isoPath` `isoFile` `expectedFile` check that the content of the ISO
// `isoFile` - matches the content of the local file `expectedFile` (using a regex
// comparison). Environment variables in `expectedFile` are substituted before the comparison.
func IsoFileCmpRegEx(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 3 {
		ts.Fatalf("usage: isofilecmpregex isoPath file1 file2")
	}

	workDir := ts.Getenv("WORK")
	isoPath, aFilePath, eFilePath := args[0], args[1], args[2]
	isoPathAbs := filepath.Join(workDir, isoPath)

	disk, err := diskfs.Open(isoPathAbs, diskfs.WithOpenMode(diskfs.ReadOnly))
	ts.Check(err)

	fs, err := disk.GetFilesystem(0)
	ts.Check(err)

	aFile, err := fs.OpenFile(aFilePath, os.O_RDONLY)
	ts.Check(err)
	defer aFile.Close()

	aData, err := io.ReadAll(aFile)
	ts.Check(err)

	eFilePathAbs := filepath.Join(workDir, eFilePath)
	eData, err := os.ReadFile(eFilePathAbs)
	ts.Check(err)

	byteCompareRegEx(ts, neg, aData, eData, aFilePath, eFilePath)
}

// InitrdImgContains `isoPath` `file` check if the specified file `file`
// is stored within a compressed cpio archive by scanning the content of
// /images/ignition.img archive in the ISO `isoPath` image (note: plain cpio
// archives are ignored).
func InitrdImgContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: initrdImgContains isoPath file")
	}

	workDir := ts.Getenv("WORK")
	isoPath, eFilePath := args[0], args[1]
	isoPathAbs := filepath.Join(workDir, isoPath)

	err := checkFileFromInitrdImg(isoPathAbs, eFilePath)
	ts.Check(err)
}

// ConfigImgContains `isoPath` `file` check if the specified file `file`
// is stored within the config image ISO.
func ConfigImgContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: configImgContains isoPath file")
	}

	workDir := ts.Getenv("WORK")
	isoPath, eFilePath := args[0], args[1]
	isoPathAbs := filepath.Join(workDir, isoPath)

	_, err := extractArchiveFile(isoPathAbs, eFilePath, "")
	ts.Check(err)
}

// UnconfiguredIgnContains `file` check if the specified file `file`
// is stored within the unconfigured ignition Storage Files.
func UnconfiguredIgnContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 1 {
		ts.Fatalf("usage: unconfiguredIgnContains file")
	}
	IgnitionStorageContains(ts, neg, []string{"unconfigured-agent.ign", args[0]})
}

// IgnitionStorageContains `ignPath` `file` check if the specified file `file`
// is stored within the ignition Storage Files.
func IgnitionStorageContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: ignitionStorageContains ignPath file")
	}

	workDir := ts.Getenv("WORK")
	ignPath, eFilePath := args[0], args[1]
	ignPathAbs := filepath.Join(workDir, ignPath)

	config, err := readIgnition(ts, ignPathAbs)
	ts.Check(err)

	found := false
	for _, f := range config.Storage.Files {
		if f.Path == eFilePath {
			found = true
		}
	}

	if !found && !neg {
		ts.Fatalf("%s does not contain %s", ignPath, eFilePath)
	}

	if neg && found {
		ts.Fatalf("%s should not contain %s", ignPath, eFilePath)
	}
}

// UnconfiguredIgnCmp `fileInIgn` `expectedFile` check that the content
// of the file `fileInIgn` extracted from the unconfigured ignition
// configuration file matches the content of the local file `expectedFile`.
// Environment variables in in `expectedFile` are substituted before the comparison.
func UnconfiguredIgnCmp(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: iunconfiguredIgnCmp file1 file2")
	}
	argsNext := []string{"unconfigured-agent.ign", args[0], args[1]}
	ignitionStorageCmp(ts, neg, argsNext)
}

// IsoContains `isoPath` `file` check if the specified `file` is stored
// within the ISO `isoPath` image.
func IsoContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: isoContains isoPath file")
	}

	workDir := ts.Getenv("WORK")
	isoPath, filePath := args[0], args[1]
	isoPathAbs := filepath.Join(workDir, isoPath)

	disk, err := diskfs.Open(isoPathAbs, diskfs.WithOpenMode(diskfs.ReadOnly))
	ts.Check(err)

	fs, err := disk.GetFilesystem(0)
	ts.Check(err)

	_, err = fs.OpenFile(filePath, os.O_RDONLY)
	ts.Check(err)
}

// ExpandFile `file...` can be used to substitute environment variables
// references for each file specified.
func ExpandFile(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 1 {
		ts.Fatalf("usage: expandFile file...")
	}

	workDir := ts.Getenv("WORK")
	for _, f := range args {
		fileName := filepath.Join(workDir, f)
		data, err := os.ReadFile(fileName)
		ts.Check(err)

		newData := expand(ts, data)
		err = os.WriteFile(fileName, []byte(newData), 0)
		ts.Check(err)
	}
}

// IsoSizeMin `isoPath` `size` checks if the specified ISO is larger
// than the specified number of bytes.
func IsoSizeMin(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: IsoSizeMin isoPath size")
	}

	isoPath := args[0]
	size, err := strconv.ParseInt(args[1], 10, 64)
	ts.Check(err)
	size *= 1000000

	fileSize := isoSize(ts, isoPath)
	if fileSize < size {
		ts.Fatalf("%s of size %d bytes is less than %d", isoPath, fileSize, size)
	}
}

// IsoSizeMax `isoPath` `size` checks if the specified ISO is smaller
// than the specified number of bytes.
func IsoSizeMax(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: IsoSizeMax isoPath size")
	}

	isoPath := args[0]
	size, err := strconv.ParseInt(args[1], 10, 64)
	ts.Check(err)
	size *= 1000000

	fileSize := isoSize(ts, isoPath)
	if fileSize > size {
		ts.Fatalf("%s of size %d bytes is greater than %d", isoPath, fileSize, size)
	}
}

func isoSize(ts *testscript.TestScript, path string) int64 {
	workDir := ts.Getenv("WORK")
	isoPathAbs := filepath.Join(workDir, path)
	fileInfo, err := os.Stat(isoPathAbs)
	ts.Check(err)

	return fileInfo.Size()
}

// ignitionStorageCmp `ignPath` `ignFile` `expectedFile` check that the content of the file
// `ignFile` - extracted from the ignition configuration file referenced
// by `ignPath` - matches the content of the local file `expectedFile`.
// Environment variables in in `expectedFile` are substituted before the comparison.
func ignitionStorageCmp(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 3 {
		ts.Fatalf("usage: ignitionStorageCmp ignPath file1 file2")
	}

	workDir := ts.Getenv("WORK")
	ignPath, aFilePath, eFilePath := args[0], args[1], args[2]
	ignPathAbs := filepath.Join(workDir, ignPath)

	config, err := readIgnition(ts, ignPathAbs)
	ts.Check(err)

	aData, err := readFileFromIgnitionCfg(&config, aFilePath)
	ts.Check(err)

	eFilePathAbs := filepath.Join(workDir, eFilePath)
	eData, err := os.ReadFile(eFilePathAbs)
	ts.Check(err)

	byteCompare(ts, neg, aData, eData, aFilePath, eFilePath)
}

func readIgnition(ts *testscript.TestScript, ignPath string) (config igntypes.Config, err error) {
	rawIgn, err := os.ReadFile(ignPath)
	ts.Check(err)
	err = json.Unmarshal(rawIgn, &config)
	return config, err
}

// archiveFileNames `isoPath` get the names of the archive files to use
// based on the name of the ISO image.
func archiveFileNames(isoPath string) (string, string, error) {
	if strings.HasPrefix(isoPath, "agent.") || strings.HasPrefix(isoPath, "node.") {
		return "/images/ignition.img", "config.ign", nil
	} else if strings.HasPrefix(isoPath, "agentconfig.") {
		return "/config.gz", "", nil
	}

	return "", "", errors.NotFound(fmt.Sprintf("ISO %s has unrecognized prefix", isoPath))
}

func expand(ts *testscript.TestScript, s []byte) string {
	return os.Expand(string(s), func(key string) string {
		if key == "$" {
			return "$"
		}
		return ts.Getenv(key)
	})
}

func byteCompare(ts *testscript.TestScript, neg bool, aData, eData []byte, aFilePath, eFilePath string) {
	byteCompareInternal(ts, neg, aData, eData, aFilePath, eFilePath, func(aText, eText string) (bool, error) {
		return aText == eText, nil
	})
}

func byteCompareRegEx(ts *testscript.TestScript, neg bool, aData, eData []byte, aFilePath, eFilePath string) {
	byteCompareInternal(ts, neg, aData, eData, aFilePath, eFilePath, func(aText, eText string) (bool, error) {
		return regexp.MatchString(eText, aText)
	})
}

func byteCompareInternal(ts *testscript.TestScript, neg bool, aData, eData []byte, aFilePath, eFilePath string, cmp func(string, string) (bool, error)) {
	aText := string(aData)
	eText := expand(ts, eData)

	eq, err := cmp(aText, eText)
	if err != nil {
		ts.Fatalf("unexpected error while comparing strings: %v", err)
	}
	if neg {
		if eq {
			ts.Fatalf("%s and %s do not differ", aFilePath, eFilePath)
		}
		return
	}
	if eq {
		return
	}

	ts.Logf("%s", aText)

	var sb strings.Builder
	if err := diff.Text(eFilePath, aFilePath, eText, aText, &sb); err != nil {
		ts.Check(err)
	}

	ts.Logf("%s", sb.String())
	ts.Fatalf("%s and %s differ", eFilePath, aFilePath)
}

func readFileFromISO(isoPath, archiveFile, ignitionFile, nodePath string) ([]byte, error) {
	config, err := extractCfgStorageData(isoPath, archiveFile, ignitionFile, nodePath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readFileFromIgnitionCfg(config *igntypes.Config, nodePath string) ([]byte, error) {
	for _, f := range config.Storage.Files {
		if f.Node.Path == nodePath {
			actualData, err := dataurl.DecodeString(*f.FileEmbedded1.Contents.Source)
			if err != nil {
				return nil, err
			}
			return actualData.Data, nil
		}
	}

	return nil, errors.NotFound(nodePath)
}

func extractArchiveFile(isoPath, archive, fileName string) ([]byte, error) {
	disk, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return nil, err
	}

	fs, err := disk.GetFilesystem(0)
	if err != nil {
		return nil, err
	}

	ignitionImg, err := fs.OpenFile(archive, os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(ignitionImg)
	if err != nil {
		return nil, err
	}

	cpioReader := cpio.NewReader(gzipReader)

	for {
		header, err := cpioReader.Next()
		if err == io.EOF { //nolint:errorlint
			// end of cpio archive
			break
		}
		if err != nil {
			return nil, err
		}

		// If the file is not in ignition return it directly
		if fileName == "" || header.Name == fileName {
			rawContent, err := io.ReadAll(cpioReader)
			if err != nil {
				return nil, err
			}
			return rawContent, nil
		}
	}

	return nil, errors.NotFound(fmt.Sprintf("File %s not found within the %s archive", fileName, archive))
}

func extractIgnition(isoPath, archiveFile, ignitionFile string) (igntypes.Config, error) {
	var config igntypes.Config

	rawContent, err := extractArchiveFile(isoPath, archiveFile, ignitionFile)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(rawContent, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func extractCfgStorageData(isoPath, archiveFile, ignitionFile, nodePath string) ([]byte, error) {
	if ignitionFile == "" {
		// If the archive is not part of an ignition file return the archive data
		rawContent, err := extractArchiveFile(isoPath, archiveFile, nodePath)
		if err != nil {
			return nil, err
		}
		return rawContent, nil
	}

	config, err := extractIgnition(isoPath, archiveFile, ignitionFile)
	if err != nil {
		return nil, err
	}
	for _, f := range config.Storage.Files {
		if f.Node.Path == nodePath {
			actualData, err := dataurl.DecodeString(*f.FileEmbedded1.Contents.Source)
			if err != nil {
				return nil, err
			}
			return actualData.Data, nil
		}
	}

	return nil, errors.NotFound(fmt.Sprintf("File %s not found within the %s archive", nodePath, archiveFile))
}

func checkFileFromInitrdImg(isoPath string, fileName string) error {
	disk, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return err
	}

	fs, err := disk.GetFilesystem(0)
	if err != nil {
		return err
	}

	initRdImg, err := fs.OpenFile("/images/pxeboot/initrd.img", os.O_RDONLY)
	if err != nil {
		return err
	}
	defer initRdImg.Close()

	const (
		gzipID1     = 0x1f
		gzipID2     = 0x8b
		gzipDeflate = 0x08
	)

	buff := make([]byte, 4096)
	for {
		_, err := initRdImg.Read(buff)
		if err == io.EOF { //nolint:errorlint
			break
		}

		foundAt := -1
		for idx := 0; idx < len(buff)-2; idx++ {
			// scan the buffer for a potential gzip header
			if buff[idx+0] == gzipID1 && buff[idx+1] == gzipID2 && buff[idx+2] == gzipDeflate {
				foundAt = idx
				break
			}
		}

		if foundAt >= 0 {
			// check if it's really a compressed cpio archive
			delta := int64(foundAt - len(buff))
			newPos, err := initRdImg.Seek(delta, io.SeekCurrent)
			if err != nil {
				break
			}

			files, err := lookForCpioFiles(initRdImg)
			if err != nil {
				if _, err := initRdImg.Seek(newPos+2, io.SeekStart); err != nil {
					break
				}
				continue
			}

			// check if the current cpio files match the required ones
			for _, f := range files {
				matched, err := filepath.Match(fileName, f)
				if err != nil {
					return err
				}
				if matched {
					return nil
				}
			}
		}
	}

	return errors.NotFound(fmt.Sprintf("File %s not found within the /images/pxeboot/initrd.img archive", fileName))
}

func lookForCpioFiles(r io.Reader) ([]string, error) {
	var files []string

	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	// skip in case of garbage
	if gr.OS != 255 && gr.OS >= 13 {
		return nil, fmt.Errorf("unknown OS code: %v", gr.Header.OS)
	}

	cr := cpio.NewReader(gr)
	for {
		h, err := cr.Next()
		if err != nil {
			break
		}

		files = append(files, h.Name)
	}

	return files, nil
}
