package compute

import (
<<<<<<< HEAD
=======
	"crypto/rsa"
	"encoding/base64"
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SSHKeysSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"public_key": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validate.NoEmptyStrings,
				},

				"username": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validate.NoEmptyStrings,
				},
			},
		},
	}
}

func ExpandSSHKeys(input []interface{}) []compute.SSHPublicKey {
	output := make([]compute.SSHPublicKey, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		username := raw["username"].(string)
		output = append(output, compute.SSHPublicKey{
			KeyData: utils.String(raw["public_key"].(string)),
			Path:    utils.String(formatUsernameForAuthorizedKeysPath(username)),
		})
	}

	return output
}

func FlattenSSHKeys(input *compute.SSHConfiguration) (*[]interface{}, error) {
	if input == nil || input.PublicKeys == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)
	for _, v := range *input.PublicKeys {
		if v.KeyData == nil || v.Path == nil {
			continue
		}

		username := parseUsernameFromAuthorizedKeysPath(*v.Path)
		if username == nil {
			return nil, fmt.Errorf("Error parsing username from %q", *v.Path)
		}

		output = append(output, map[string]interface{}{
			"public_key": *v.KeyData,
			"username":   *username,
		})
	}

	return &output, nil
}

// formatUsernameForAuthorizedKeysPath returns the path to the authorized keys file
// for the specified username
func formatUsernameForAuthorizedKeysPath(username string) string {
	return fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
}

// parseUsernameFromAuthorizedKeysPath parses the username out of the authorized keys
// path returned from the Azure API
func parseUsernameFromAuthorizedKeysPath(input string) *string {
	// the Azure VM agent hard-codes this to `/home/username/.ssh/authorized_keys`
	// as such we can hard-code this for a better UX
	compiled, err := regexp.Compile("(/home/)+(?P<username>.*?)(/.ssh/authorized_keys)+")
	if err != nil {
		return nil
	}

	keys := compiled.SubexpNames()
	values := compiled.FindStringSubmatch(input)

	if values == nil {
		return nil
	}

	for i, k := range keys {
		if k == "username" {
			value := values[i]
			return &value
		}
	}

	return nil
}
<<<<<<< HEAD
=======

// ValidateSSHKey performs some basic validation on supplied SSH Keys - Encoded Signature and Key Size are evaluated
// Will require rework if/when other Key Types are supported
func ValidateSSHKey(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string or whitespace", k)}
	}

	keyParts := strings.Fields(v)
	if len(keyParts) > 1 {
		byteStr, err := base64.StdEncoding.DecodeString(keyParts[1])
		if err != nil {
			return nil, []error{fmt.Errorf("Error decoding %q for public key data", k)}
		}
		pubKey, err := ssh.ParsePublicKey(byteStr)
		if err != nil {
			return nil, []error{fmt.Errorf("Error parsing %q as a public key object", k)}
		}

		if pubKey.Type() != ssh.KeyAlgoRSA {
			return nil, []error{fmt.Errorf("Error - the provided %s SSH key is not supported. Only RSA SSH keys are supported by Azure", pubKey.Type())}
		} else {
			rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
			if !ok {
				return nil, []error{fmt.Errorf("Error - could not retrieve the RSA public key from the SSH public key")}
			}
			rsaPubKeyBits := rsaPubKey.Size() * 8
			if rsaPubKeyBits < 2048 {
				return nil, []error{fmt.Errorf("Error - the provided RSA SSH key has %d bits. Only ssh-rsa keys with 2048 bits or higher are supported by Azure", rsaPubKeyBits)}
			}
		}
	} else {
		return nil, []error{fmt.Errorf("Error %q is not a complete SSH2 Public Key", k)}
	}

	return warnings, errors
}
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0
