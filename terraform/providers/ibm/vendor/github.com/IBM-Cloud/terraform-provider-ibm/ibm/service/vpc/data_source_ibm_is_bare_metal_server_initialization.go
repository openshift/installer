// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/ScaleFT/sshkeys"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerImageName                = "image_name"
	isBareMetalServerUserAccountUserName      = "username"
	isBareMetalServerUserAccountEncryptionKey = "encryption_key"
	isBareMetalServerUserAccountEncPwd        = "encrypted_password"
	isBareMetalServerUserAccountPassword      = "password"
	isBareMetalServerPEM                      = "private_key"
	isBareMetalServerPassphrase               = "passphrase"
	isBareMetalServerUserAccountResourceType  = "resource_type"
)

func DataSourceIBMIsBareMetalServerInitialization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerInitializationRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},

			isBareMetalServerPEM: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Bare Metal Server Private Key file",
			},

			isBareMetalServerPassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase for Bare Metal Server Private Key file",
			},

			isBareMetalServerImage: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the image the bare metal server was provisioned from",
			},

			isBareMetalServerImageName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined or system-provided name for the image the bare metal server was provisioned from",
			},

			isBareMetalServerUserAccounts: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The user accounts that are created at initialization. There can be multiple account types distinguished by the resource_type attribute.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerUserAccountUserName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username for the account created at initialization",
						},
						isBareMetalServerUserAccountEncryptionKey: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for the encryption key",
						},
						isBareMetalServerUserAccountEncPwd: {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The password at initialization, encrypted using encryption_key, and returned base64-encoded",
						},
						isBareMetalServerUserAccountPassword: {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The password at initialization, encrypted using encryption_key, and returned base64-encoded",
						},
						isBareMetalServerUserAccountResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced : [ host_user_account ]",
						},
					},
				},
			},

			isBareMetalServerKeys: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "SSH key Ids for the bare metal server",
			},
		},
	}
}

func dataSourceIBMISBareMetalServerInitializationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetBareMetalServerInitializationOptions{
		ID: &bareMetalServerID,
	}

	initialization, response, err := sess.GetBareMetalServerInitializationWithContext(context, options)
	if err != nil || initialization == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) initialization : %s\n%s", bareMetalServerID, err, response))
	}
	d.SetId(bareMetalServerID)
	if initialization.Image != nil {
		d.Set(isBareMetalServerImage, initialization.Image.ID)
		d.Set(isBareMetalServerImageName, initialization.Image.Name)
	}

	var keys []string
	keys = make([]string, 0)
	if initialization.Keys != nil {
		for _, key := range initialization.Keys {
			keys = append(keys, *key.ID)
		}
	}
	d.Set(isBareMetalServerKeys, flex.NewStringSet(schema.HashString, keys))
	accList := make([]map[string]interface{}, 0)
	if initialization.UserAccounts != nil {

		for _, accIntf := range initialization.UserAccounts {
			acc := accIntf.(*vpcv1.BareMetalServerInitializationUserAccount)
			currAccount := map[string]interface{}{
				isBareMetalServerUserAccountUserName: *acc.Username,
			}
			currAccount[isBareMetalServerUserAccountResourceType] = *acc.ResourceType
			currAccount[isBareMetalServerUserAccountEncryptionKey] = *acc.EncryptionKey.CRN
			encPassword := base64.StdEncoding.EncodeToString(*acc.EncryptedPassword)
			currAccount[isBareMetalServerUserAccountEncPwd] = encPassword

			var rsaKey *rsa.PrivateKey
			if privatekey, ok := d.GetOk(isBareMetalServerPEM); ok {
				keyFlag := privatekey.(string)
				keybytes := []byte(keyFlag)

				if keyFlag != "" {
					block, err := pem.Decode(keybytes)
					if block == nil {
						return diag.FromErr(fmt.Errorf("[ERROR] Failed to load the private key from the given key contents. Instead of the key file path, please make sure the private key is pem format"))
					}
					isEncrypted := false
					switch block.Type {
					case "RSA PRIVATE KEY":
						isEncrypted = x509.IsEncryptedPEMBlock(block)
					case "OPENSSH PRIVATE KEY":
						var err error
						isEncrypted, err = isOpenSSHPrivKeyEncrypted(block.Bytes)
						if err != nil {
							return diag.FromErr(fmt.Errorf("[ERROR] Failed to check if the provided open ssh key is encrypted or not %s", err))
						}
					default:
						return diag.FromErr(fmt.Errorf("[ERROR] PEM and OpenSSH private key formats with RSA key type are supported, can not support this key file type: %s", err))
					}
					passphrase := ""
					var privateKey interface{}
					if isEncrypted {
						if pass, ok := d.GetOk(isBareMetalServerPassphrase); ok {
							passphrase = pass.(string)
						} else {
							return diag.FromErr(fmt.Errorf("[ERROR] Mandatory field 'passphrase' not provided"))
						}
						var err error
						privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, []byte(passphrase))
						if err != nil {
							return diag.FromErr(fmt.Errorf("[ERROR] Fail to decrypting the private key: %s", err))
						}
					} else {
						var err error
						privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, nil)
						if err != nil {
							return diag.FromErr(fmt.Errorf("[ERROR] Fail to decrypting the private key: %s", err))
						}
					}
					var ok bool
					rsaKey, ok = privateKey.(*rsa.PrivateKey)
					if !ok {
						return diag.FromErr(fmt.Errorf("[ERROR] Failed to convert to RSA private key"))
					}
				}
			}

			if acc.EncryptedPassword != nil {
				ciphertext := *acc.EncryptedPassword
				password := base64.StdEncoding.EncodeToString(ciphertext)
				if rsaKey != nil {
					rng := rand.Reader
					clearPassword, err := rsa.DecryptPKCS1v15(rng, rsaKey, ciphertext)
					if err != nil {
						return diag.FromErr(fmt.Errorf("[ERROR] Can not decrypt the password with the given key, %s", err))
					}
					password = string(clearPassword)
				}
				currAccount[isBareMetalServerUserAccountPassword] = password
			}
			accList = append(accList, currAccount)
		}
		d.Set(isBareMetalServerUserAccounts, accList)
	}
	return nil
}
