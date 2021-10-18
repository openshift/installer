package ignition

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"path/filepath"

	"github.com/clarketm/json"
	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/asset"
)

// Marshal is a helper function to use the marshaler function from "github.com/clarketm/json".
// It supports zero values of structs with the omittempty annotation.
// In effect this excludes empty pointer struct fields from the marshaled data,
// instead of inserting nil values into them.
// This is necessary for ignition configs to pass openAPI validation on fields
// that are not supposed to contain nil pointers, but e.g. strings.
// It can be used as a dropin replacement for "encoding/json".Marshal
func Marshal(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

// FilesFromAsset creates ignition files for each of the files in the specified
// asset.
func FilesFromAsset(pathPrefix string, username string, mode int, asset asset.WritableAsset) []igntypes.File {
	var files []igntypes.File
	for _, f := range asset.Files() {
		files = append(files, FileFromBytes(filepath.Join(pathPrefix, f.Filename), username, mode, f.Data))
	}
	return files
}

// FileFromString creates an ignition-config file with the given contents.
func FileFromString(path string, username string, mode int, contents string) igntypes.File {
	return FileFromBytes(path, username, mode, []byte(contents))
}

// FileFromBytes creates an ignition-config file with the given contents.
func FileFromBytes(path string, username string, mode int, contents []byte) igntypes.File {
	return igntypes.File{
		Node: igntypes.Node{
			Path: path,
			User: igntypes.NodeUser{
				Name: &username,
			},
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &mode,
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(dataurl.EncodeBytes(contents)),
			},
		},
	}
}

// ConvertToRawExtension converts and ignition config to a RawExtension containing the ignition as raw bytes
func ConvertToRawExtension(config igntypes.Config) (runtime.RawExtension, error) {
	rawIgnConfig, err := json.Marshal(config)
	if err != nil {
		return runtime.RawExtension{}, fmt.Errorf("failed to marshal Ignition config: %v", err)
	}

	return runtime.RawExtension{
		Raw: rawIgnConfig,
	}, nil
}

// ConvertToAppendix converts the contents of an ignition file to an appendix.
// In ignition config spec v2 the `Append` boolean value was used to denote whether
// the `Contents` field was an appendix or not. It was also permitted to define
// multiple file configs (appendix or not) that would be merged/overwritten
// sequentially in the order of the json data, which made them non-deterministic.
// In spec v3 this has changed with only one config allowed for each file config,
// and `Append` now being a list of objects that are being appended to `Contents`,
// with the `Contents` field itself never being an appendix.
// This function moves an ignition file's `Contents` object into the `Append` list.
// Since the resulting ignition file of this function has an empty `Contents` field,
// `Overwrite` must be set to false, per the spec.
// The output is an ignition file config that will write a new file with only the
// appendix contents in the case of a file not already existing on disk,
// or append the appendix contents to a file already existing.
func ConvertToAppendix(file *igntypes.File) {
	file.Append = []igntypes.Resource{
		file.Contents,
	}
	file.Contents = igntypes.Resource{}
	file.Overwrite = ignutil.BoolToPtr(false)
}

// parseCertificateBundle loads each certificate in the bundle to the Ingition
// carrier type, ignoring any invisible character before, after and in between
// certificates.
func parseCertificateBundle(userCA []byte) ([]igntypes.Resource, error) {
	userCA = bytes.TrimSpace(userCA)

	var carefs []igntypes.Resource
	for len(userCA) > 0 {
		var block *pem.Block
		block, userCA = pem.Decode(userCA)
		if block == nil {
			return nil, fmt.Errorf("unable to parse certificate, please check the certificates")
		}

		carefs = append(carefs, igntypes.Resource{Source: ignutil.StrToPtr(dataurl.EncodeBytes(pem.EncodeToMemory(block)))})

		userCA = bytes.TrimSpace(userCA)
	}

	return carefs, nil
}

// GenerateIgnitionShim is used to generate an ignition file that contains a user ca bundle
// in its Security section.
func GenerateIgnitionShim(bootstrapConfigURL string, userCA string) ([]byte, error) {
	ign := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Replace: igntypes.Resource{
					Source: ignutil.StrToPtr(bootstrapConfigURL),
				},
			},
		},
	}

	carefs, err := parseCertificateBundle([]byte(userCA))
	if err != nil {
		return nil, err
	}
	if len(carefs) > 0 {
		ign.Ignition.Security = igntypes.Security{
			TLS: igntypes.TLS{
				CertificateAuthorities: carefs,
			},
		}
	}

	data, err := Marshal(ign)
	if err != nil {
		return nil, err
	}

	return data, nil
}
