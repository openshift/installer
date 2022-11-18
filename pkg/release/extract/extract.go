package extract

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/containers/image/v5/manifest"
	"github.com/opencontainers/go-digest"
	"github.com/openshift/library-go/pkg/image/dockerv1client"
	"golang.org/x/crypto/openpgp"

	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/MakeNowJust/heredoc"
	imageapi "github.com/openshift/api/image/v1"
	imagereference "github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/oc/pkg/cli/admin/release"
	"github.com/openshift/oc/pkg/cli/image/extract"
	"github.com/openshift/oc/pkg/cli/image/imagesource"
	imagemanifest "github.com/openshift/oc/pkg/cli/image/manifest"
)

// extractTarget describes how a file in the release image can be extracted to disk.
type extractTarget struct {
	OS       string
	Arch     string
	Command  string
	Optional bool
	NewArch  bool

	InjectReleaseImage   bool
	InjectReleaseVersion bool

	ArchiveFormat string
	AsArchive     bool
	AsZip         bool
	Readme        string
	LinkTo        []string

	Mapping extract.Mapping
}

const (
	targetReleaseArch = "release-arch"
)

var hashFn = sha256.New
var signer *openpgp.Entity

type ExtractOptions struct {
	genericclioptions.IOStreams

	SecurityOptions imagemanifest.SecurityOptions
	FilterOptions   imagemanifest.FilterOptions
	ParallelOptions imagemanifest.ParallelOptions

	ICSPFile string

	Output string

	FromDir string
	From    string

	Tools                  bool
	Command                string
	CommandOperatingSystem string
	SigningKey             string

	Cloud string

	// GitExtractDir is the path of a root directory to extract the source of a release to.
	GitExtractDir string

	Directory string
	File      string
	FileDir   string

	ExtractManifests bool
	Manifests        []manifest.Manifest

	ImageMetadataCallback func(m *extract.Mapping, dgst, contentDigest digest.Digest, config *dockerv1client.DockerImageConfig, manifestListDigest digest.Digest)
}

// extractTools extracts specific commands out of images referenced by the release image.
// TODO: in the future the metadata this command contains might be loaded from the release
//
//	image, but we must maintain compatibility with older payloads if so
func (o *ExtractOptions) Extract(command string) error {
	// Available targets is treated as a GA API and may not be changed without backwards
	// compatibility of at least N-2 releases.

	fmt.Println("EXTRACTING")
	currentArch := runtime.GOARCH
	currentOS := runtime.GOOS
	if len(o.CommandOperatingSystem) > 0 {
		currentOS = o.CommandOperatingSystem
		if currentOS == "*" {
			currentArch = "*"
		} else {
			currentArch = "amd64"
		}
	}
	if currentOS == "mac" {
		currentOS = "darwin"
	}

	// Select the subset of targets based on command line input
	var willArchive bool
	var targets []extractTarget
	fmt.Println("CHECKING AVAILABLE TARGETS")
	// Filter by command, or gather all non-optional targets
	if len(command) > 0 {
		for _, target := range availableTargets {
			fmt.Println(target, command)
			if target.Command == command {
				targets = append(targets, target)
				fmt.Println("ADDED", targets)
			}
		}
	} else {
		for _, target := range availableTargets {
			if !target.Optional {
				targets = append(targets, target)
			}
		}
	}

	// If the user didn't specify a command, or the operating system is set
	// to '*', we'll produce an archive
	if len(command) == 0 || o.CommandOperatingSystem == "*" {
		for i := range targets {
			targets[i].AsArchive = true
			targets[i].AsZip = targets[i].OS == "windows"
		}
	}
	fmt.Println("TARGETS", targets)
	if len(targets) == 0 {
		switch {
		case len(command) > 0 && currentOS != "*":
			return fmt.Errorf("command %q does not support the operating system %q", o.Command, currentOS)
		case len(command) > 0:
			return fmt.Errorf("the supported commands are 'oc' and 'openshift-install'")
		default:
			return fmt.Errorf("no available commands")
		}
	}
	fmt.Println("Finding release")
	// load the release image
	dir := o.Directory
	infoOptions := release.NewInfoOptions(o.IOStreams)
	infoOptions.SecurityOptions = o.SecurityOptions
	infoOptions.FilterOptions = o.FilterOptions
	infoOptions.FileDir = o.FileDir
	infoOptions.ICSPFile = o.ICSPFile
	release, err := infoOptions.LoadReleaseInfo(o.From, false)
	if err != nil {
		fmt.Println("NO RELEASEEEEEEEEEEEEEEEE", err)
		return err
	}
	releaseArch := release.Config.Architecture
	releaseName := release.PreferredName()
	refExact := release.ImageRef
	refExact.Ref.Tag = ""
	// if the release image is manifestlist image, we'll not change digest with
	// arch based sha. Because we want that the extracted tool can be used for all archs.
	//if !release.IsManifestList {
	//	refExact.Ref.ID = release.Digest.String()
	//}
	exactReleaseImage := refExact.String()
	// resolve target image references to their pull specs
	missing := sets.NewString()
	var validTargets []extractTarget
	for _, target := range targets {
		if currentOS != "*" && target.OS != currentOS {
			klog.V(2).Infof("Skipping %s, does not match current OS %s", target.ArchiveFormat, target.OS)
			continue
		}
		if currentArch != "*" && target.Arch != currentArch {
			if currentArch != releaseArch || target.Arch != targetReleaseArch {
				klog.V(2).Infof("Skipping %s, does not match current architecture %s", target.ArchiveFormat, target.Arch)
				continue
			}
		}
		// TODO: Figure this out
		// if target.OS == "linux" && target.Arch == releaseArch {
		// 	klog.V(2).Infof("Skipping duplicate %s", target.ArchiveFormat)
		// 	continue
		// }

		spec, err := findImageSpec(release.References, target.Mapping.Image, o.From)
		if err != nil && !target.NewArch {
			missing.Insert(target.Mapping.Image)
			continue
		}
		klog.V(2).Infof("Will extract %s from %s", target.Mapping.From, spec)
		ref, err := imagereference.Parse(spec)
		if err != nil {
			return err
		}
		target.Mapping.Image = spec
		target.Mapping.ImageRef = imagesource.TypedImageReference{Ref: ref, Type: imagesource.DestinationRegistry}
		if target.AsArchive {
			willArchive = true
			target.Mapping.Name = fmt.Sprintf(target.ArchiveFormat, releaseName)
			target.Mapping.To = filepath.Join(dir, target.Mapping.Name)
		} else {
			target.Mapping.To = filepath.Join(dir, target.Command)
			target.Mapping.Name = fmt.Sprintf("%s-%s-%s", target.OS, target.Arch, target.Command)
		}
		validTargets = append(validTargets, target)
	}
	fmt.Println("validTargets", validTargets)
	if len(validTargets) == 0 {
		if len(missing) == 1 {
			return fmt.Errorf("the image %q containing the desired command is not available", missing.List()[0])
		}
		return fmt.Errorf("some required images are missing: %s", strings.Join(missing.List(), ", "))
	}
	if len(missing) > 0 {
		fmt.Fprintf(o.ErrOut, "warning: Some commands can not be extracted due to missing images: %s\n", strings.Join(missing.List(), ", "))
	}

	// will extract in parallel
	opts := extract.NewExtractOptions(genericclioptions.IOStreams{Out: o.Out, ErrOut: o.ErrOut})
	opts.ParallelOptions = o.ParallelOptions
	opts.SecurityOptions = o.SecurityOptions
	opts.FilterOptions = o.FilterOptions
	opts.ICSPFile = o.ICSPFile
	opts.OnlyFiles = true

	// create the mapping lookup of the valid targets
	var extractLock sync.Mutex
	targetsByName := make(map[string]extractTarget)
	for _, target := range validTargets {
		targetsByName[target.Mapping.Name] = target
		opts.Mappings = append(opts.Mappings, target.Mapping)
	}
	hashByTargetName := make(map[string]string)

	// ensure to is a directory
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	// as each layer is extracted, take the output binary and write it to disk
	opts.TarEntryCallback = func(hdr *tar.Header, layer extract.LayerInfo, r io.Reader) (bool, error) {
		// ensure we don't process the same mapping twice due to programmer error
		target, ok := func() (extractTarget, bool) {
			extractLock.Lock()
			defer extractLock.Unlock()
			target, ok := targetsByName[layer.Mapping.Name]
			return target, ok
		}()
		if !ok {
			return false, fmt.Errorf("unable to find target with mapping name %s", layer.Mapping.Name)
		}

		// open the file
		f, err := os.OpenFile(layer.Mapping.To, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return false, err
		}

		// if we need to write an archive, wrap the file appropriately to create a single
		// entry
		var w io.Writer = f

		bw := bufio.NewWriterSize(w, 16*1024)
		w = bw

		var hash hash.Hash
		closeFn := func() error { return nil }
		if target.AsArchive {
			text := strings.Replace(target.Readme, `\u0060`, "`", -1)
			hash = hashFn()
			w = io.MultiWriter(hash, w)
			if target.AsZip {
				klog.V(2).Infof("Writing %s as a ZIP archive %s", hdr.Name, layer.Mapping.To)
				zw := zip.NewWriter(w)

				if len(text) > 0 {
					text = strings.Replace(text, "\n", "\r\n", -1)
					zh := &zip.FileHeader{
						Method:             zip.Deflate,
						Name:               "README.md",
						UncompressedSize64: uint64(len(text)),
						Modified:           hdr.ModTime,
					}
					zh.SetMode(os.FileMode(0755))

					fw, err := zw.CreateHeader(zh)
					if err != nil {
						return false, err
					}
					if _, err := fmt.Fprintf(fw, text); err != nil {
						return false, err
					}
				}

				zh := &zip.FileHeader{
					Method:             zip.Deflate,
					Name:               target.Command + ".exe",
					UncompressedSize64: uint64(hdr.Size),
					Modified:           hdr.ModTime,
				}
				zh.SetMode(os.FileMode(0755))

				fw, err := zw.CreateHeader(zh)
				if err != nil {
					return false, err
				}

				w = fw
				closeFn = func() error { return zw.Close() }

			} else {
				klog.V(2).Infof("Writing %s as a tar.gz archive %s", hdr.Name, layer.Mapping.To)
				gw, err := gzip.NewWriterLevel(w, 3)
				if err != nil {
					return false, err
				}
				tw := tar.NewWriter(gw)

				if len(text) > 0 {
					if err := tw.WriteHeader(&tar.Header{
						Name:     "README.md",
						Mode:     int64(os.FileMode(0644).Perm()),
						Size:     int64(len(text)),
						Typeflag: tar.TypeReg,
						ModTime:  hdr.ModTime,
					}); err != nil {
						return false, err
					}
					if _, err := fmt.Fprintf(tw, text); err != nil {
						return false, err
					}
				}

				if err := tw.WriteHeader(&tar.Header{
					Name:     target.Command,
					Mode:     int64(os.FileMode(0755).Perm()),
					Size:     hdr.Size,
					Typeflag: tar.TypeReg,
					ModTime:  hdr.ModTime,
				}); err != nil {
					return false, err
				}

				w = tw
				closeFn = func() error {
					for _, link := range target.LinkTo {
						if err := tw.WriteHeader(&tar.Header{
							Name:     link,
							Mode:     int64(os.FileMode(0755).Perm()),
							Size:     0,
							Typeflag: tar.TypeLink,
							ModTime:  hdr.ModTime,
							Linkname: target.Command,
						}); err != nil {
							return err
						}
					}
					if err := tw.Close(); err != nil {
						return err
					}
					return gw.Close()
				}
			}
		}

		// copy the input to disk
		replacements := []replacement{}
		if target.InjectReleaseImage {
			replacements = append(replacements, replacement{
				name:   "release image",
				marker: append([]byte{0}, []byte(releaseImageMarker[1:])...),
				value:  exactReleaseImage,
			})
		}
		if target.InjectReleaseVersion {
			replacements = append(replacements, replacement{
				name:   "release version",
				marker: append([]byte{0}, []byte(releaseVersionMarker[1:])...),
				value:  releaseName,
			})
		}
		err = copyAndReplace(o.ErrOut, w, r, 4*1024, replacements, target.Command)
		if err != nil {
			closeFn()
			f.Close()
			os.Remove(f.Name())
			return false, err
		}

		// ensure the file is written to disk
		if err := closeFn(); err != nil {
			return false, err
		}
		if err := bw.Flush(); err != nil {
			return false, err
		}
		if err := f.Close(); err != nil {
			return false, err
		}
		if err := os.Chtimes(f.Name(), hdr.ModTime, hdr.ModTime); err != nil {
			klog.V(2).Infof("Unable to set extracted file modification time: %v", err)
		}

		func() {
			extractLock.Lock()
			defer extractLock.Unlock()
			delete(targetsByName, layer.Mapping.Name)
			if hash != nil {
				hashByTargetName[layer.Mapping.To] = hex.EncodeToString(hash.Sum(nil))
			}
		}()

		return false, nil
	}
	fmt.Println("RUNNING")
	if err := opts.Run(); err != nil {
		return err
	}
	fmt.Println("DONE RUNNING")

	if willArchive {
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, heredoc.Doc(`
			Client tools for OpenShift
			--------------------------

			These archives contain the client tooling for [OpenShift](https://docs.openshift.com).

			To verify the contents of this directory, use the 'gpg' and 'shasum' tools to
			ensure the archives you have downloaded match those published from this location.

			The openshift-install binary has been preconfigured to install the following release:

			---

		`))

		hash := hashFn()
		hash.Write(buf.Bytes())
		hashByTargetName["release.txt"] = hex.EncodeToString(hash.Sum(nil))
	}

	// write a checksum of the tar files to disk as sha256sum.txt.asc
	if len(hashByTargetName) > 0 {
		var keys []string
		for k := range hashByTargetName {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var lines []string
		for _, k := range keys {
			hash := hashByTargetName[k]
			lines = append(lines, fmt.Sprintf("%s  %s", hash, filepath.Base(k)))
		}
		// ensure a trailing newline
		if len(lines[len(lines)-1]) != 0 {
			lines = append(lines, "")
		}
		// write the content manifest
		data := []byte(strings.Join(lines, "\n"))
		filename := "sha256sum.txt"
		if err := ioutil.WriteFile(filepath.Join(dir, filename), data, 0644); err != nil {
			return fmt.Errorf("unable to write checksum file: %v", err)
		}
		// sign the content manifest
		if signer != nil {
			buf := &bytes.Buffer{}
			if err := openpgp.ArmoredDetachSign(buf, signer, bytes.NewBuffer(data), nil); err != nil {
				return fmt.Errorf("unable to sign the sha256sum.txt file: %v", err)
			}
			if err := ioutil.WriteFile(filepath.Join(dir, filename+".asc"), buf.Bytes(), 0644); err != nil {
				return fmt.Errorf("unable to write signed manifest: %v", err)
			}
		}
	}

	// if we did not process some targets, report that to the user and error if necessary
	if len(targetsByName) > 0 {
		var missing []string
		for _, target := range targetsByName {
			if target.NewArch {
				continue
			}
			missing = append(missing, target.Mapping.From)
		}
		sort.Strings(missing)
		if len(missing) == 1 {
			return fmt.Errorf("image did not contain %s", missing[0])
		}
		if len(missing) > 1 {
			return fmt.Errorf("unable to find multiple files: %s", strings.Join(missing, ", "))
		}
	}

	return nil
}

const (
	// releaseImageMarker is the placeholder within a binary for the release image pullspec.
	releaseImageMarker = "!_RELEASE_IMAGE_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	// releaseVersionMarker is the placeholder within a binary for the release image version name string.
	releaseVersionMarker = "!_RELEASE_VERSION_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
)

type replacement struct {
	name   string
	marker []byte
	value  string
}

// copyAndReplace performs a targeted replacement for binaries that
// contain special marker strings, replacing the first occurrence of each
// marker with a new string and a NUL terminating byte.  It logs a warning
// if any replacements are not performed.
func copyAndReplace(errorOutput io.Writer, w io.Writer, r io.Reader, bufferSize int, replacements []replacement, name string) error {
	if len(replacements) == 0 {
		_, err := io.Copy(w, r)
		return err
	}

	longestMarker := 0
	for _, replacement := range replacements {
		if len(replacement.value) > len(replacement.marker)-1 {
			return fmt.Errorf("the %s value has %d bytes, but the maximum replacement length is %d", replacement.name, len(replacement.value), len(replacement.marker)-1)
		}

		if len(replacement.marker) > longestMarker {
			longestMarker = len(replacement.marker)
			if bufferSize < longestMarker {
				return fmt.Errorf("the buffer size must be greater than %d bytes to find %s", len(replacement.marker), replacement.name)
			}
		}
	}

	offset := 0
	buf := make([]byte, bufferSize)
	remaining := len(replacements)
	matches := make([]bool, remaining)

	for {
		n, err := io.ReadFull(r, buf[offset:])

		// search in the buffer for the expected match
		end := offset + n
		if n > 0 {
			for i, matched := range matches {
				if !matched {
					replacement := replacements[i]
					index := bytes.Index(buf[:end], replacement.marker)
					if index != -1 {
						klog.V(2).Infof("Found match for %s at %d (len=%d, offset=%d, n=%d)", replacements[i].name, index, len(buf), offset, n)
						// the replacement starts at the beginning of the match, contains the replacement value and a terminating NUL byte
						copy(buf[index:index+len(replacement.value)], []byte(replacement.value))
						buf[index+len(replacement.value)] = 0x00
						matches[i] = true
						remaining -= 1
					}
				}
			}
		}

		// write everything that we have already searched (excluding the end of the buffer that will
		// be checked next pass)
		writeTo := end - longestMarker
		if writeTo < 0 {
			writeTo = 0
		}
		if remaining == 0 || err != nil {
			writeTo = end
		}
		offset = 0
		for offset < writeTo {
			n, wErr := w.Write(buf[offset:writeTo])
			if wErr != nil {
				return wErr
			}
			offset += n
		}

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				remainingNames := make([]string, 0, len(replacements))
				for i, matched := range matches {
					if !matched {
						remainingNames = append(remainingNames, replacements[i].name)
					}
				}
				sort.Strings(remainingNames)
				if len(remainingNames) > 0 {
					fmt.Fprintf(errorOutput, "warning: Unable to make all expected replacements in %s.  Remaining: %s", name, strings.Join(remainingNames, ", "))
				}
				return nil
			}
			return err
		}

		// once we match each replacement, we can copy the rest of the file without processing
		if remaining == 0 {
			_, err := io.Copy(w, r)
			return err
		}

		// ensure the beginning of the buffer matches the end of the current buffer so that we
		// can search for matches that span buffers
		copy(buf[:writeTo], buf[writeTo:end])
		offset = end - writeTo
	}
}

func findImageSpec(image *imageapi.ImageStream, tagName, imageName string) (string, error) {
	for _, tag := range image.Spec.Tags {
		if tag.Name == tagName {
			if tag.From != nil && tag.From.Kind == "DockerImage" && len(tag.From.Name) > 0 {
				return tag.From.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no image tag %q exists in the release image %s", tagName, imageName)
}
