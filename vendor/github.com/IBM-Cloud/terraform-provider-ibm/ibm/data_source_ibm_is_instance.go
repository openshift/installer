// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/ScaleFT/sshkeys"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/crypto/ssh"
)

const (
	isInstancePEM          = "private_key"
	isInstancePassphrase   = "passphrase"
	isInstanceInitPassword = "password"
	isInstanceInitKeys     = "keys"
)

func dataSourceIBMISInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceRead,

		Schema: map[string]*schema.Schema{

			isInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name",
			},

			isInstancePEM: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance Private Key file",
			},

			isInstancePassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase for Instance Private Key file",
			},

			isInstanceInitPassword: {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "password for Windows Instance",
			},

			isInstanceInitKeys: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance keys",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance key id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance key name",
						},
					},
				},
			},

			isInstanceVPC: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC id",
			},

			isInstanceZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone name",
			},

			isInstanceProfile: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Profile info",
			},

			isInstanceTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         resourceIBMVPCHash,
				Description: "list of tags for the instance",
			},
			isInstanceBootVolume: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance Boot Volume",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume name",
						},
						"device": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume device",
						},
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume id",
						},
						"volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume name",
						},
						"volume_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume CRN",
						},
					},
				},
			},

			isInstanceVolumeAttachments: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance Volume Attachments",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Volume Attachment id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Volume Attachment name",
						},
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume id",
						},
						"volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume name",
						},
						"volume_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume CRN",
						},
					},
				},
			},

			isInstancePrimaryNetworkInterface: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface id",
						},
						isInstanceNicName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface name",
						},
						isInstanceNicPortSpeed: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance Primary Network Interface port speed",
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface IPV4 Address",
						},
						isInstanceNicSecurityGroups: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Instance Primary Network Interface Security groups",
						},
						isInstanceNicSubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface subnet",
						},
					},
				},
			},

			isInstanceNetworkInterfaces: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface id",
						},
						isInstanceNicName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface name",
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface IPV4 Address",
						},
						isInstanceNicSecurityGroups: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Instance Network Interface Security Groups",
						},
						isInstanceNicSubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface subnet",
						},
					},
				},
			},

			isInstanceImage: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance Image",
			},

			isInstanceVolumes: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of volumes",
			},

			isInstanceResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance resource group",
			},

			isInstanceCPU: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance vCPU",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceCPUArch: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance vCPU Architecture",
						},
						isInstanceCPUCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance vCPU count",
						},
					},
				},
			},

			isInstanceGpu: {
				Type:        schema.TypeList,
				Computed:    true,
				Deprecated:  "This field is deprecated",
				Description: "Instance GPU",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceGpuCores: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance GPU Cores",
						},
						isInstanceGpuCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance GPU Count",
						},
						isInstanceGpuMemory: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance GPU Memory",
						},
						isInstanceGpuManufacturer: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance GPU Manufacturer",
						},
						isInstanceGpuModel: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance GPU Model",
						},
					},
				},
			},

			isInstanceMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance memory",
			},

			isInstanceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance status",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			isInstanceDisks: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the instance's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the disk was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance disk.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance disk.",
						},
						"interface_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this disk.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes).",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISInstanceRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	name := d.Get(isInstanceName).(string)
	if userDetails.generation == 1 {
		err := classicInstanceGetByName(d, meta, name)
		if err != nil {
			return err
		}
	} else {
		err := instanceGetByName(d, meta, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicInstanceGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.Instance{}
	for {
		listInstancesOptions := &vpcclassicv1.ListInstancesOptions{}
		if start != "" {
			listInstancesOptions.Start = &start
		}
		instances, response, err := sess.ListInstances(listInstancesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Instances %s\n%s", err, response)
		}
		start = GetNext(instances.Next)
		allrecs = append(allrecs, instances.Instances...)
		if start == "" {
			break
		}
	}
	for _, instance := range allrecs {
		if *instance.Name == name {
			d.SetId(*instance.ID)
			id := *instance.ID
			d.Set(isInstanceName, *instance.Name)
			if instance.Profile != nil {
				d.Set(isInstanceProfile, *instance.Profile.Name)
			}
			cpuList := make([]map[string]interface{}, 0)
			if instance.Vcpu != nil {
				currentCPU := map[string]interface{}{}
				currentCPU[isInstanceCPUArch] = *instance.Vcpu.Architecture
				currentCPU[isInstanceCPUCount] = *instance.Vcpu.Count
				cpuList = append(cpuList, currentCPU)
			}
			d.Set(isInstanceCPU, cpuList)

			d.Set(isInstanceMemory, *instance.Memory)
			gpuList := make([]map[string]interface{}, 0)
			d.Set(isInstanceGpu, gpuList)

			if instance.PrimaryNetworkInterface != nil {
				primaryNicList := make([]map[string]interface{}, 0)
				currentPrimNic := map[string]interface{}{}
				currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
				currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name
				currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
				getnicoptions := &vpcclassicv1.GetInstanceNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         instance.PrimaryNetworkInterface.ID,
				}
				insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
				if err != nil {
					return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
				}
				currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
				if len(insnic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(insnic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
				}

				primaryNicList = append(primaryNicList, currentPrimNic)
				d.Set(isInstancePrimaryNetworkInterface, primaryNicList)
			}

			if instance.NetworkInterfaces != nil {
				interfacesList := make([]map[string]interface{}, 0)
				for _, intfc := range instance.NetworkInterfaces {
					if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
						currentNic := map[string]interface{}{}
						currentNic["id"] = *intfc.ID
						currentNic[isInstanceNicName] = *intfc.Name
						currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
						getnicoptions := &vpcclassicv1.GetInstanceNetworkInterfaceOptions{
							InstanceID: &id,
							ID:         intfc.ID,
						}
						insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
						if err != nil {
							return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
						}
						currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
						if len(insnic.SecurityGroups) != 0 {
							secgrpList := []string{}
							for i := 0; i < len(insnic.SecurityGroups); i++ {
								secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
							}
							currentNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
						}
						interfacesList = append(interfacesList, currentNic)

					}
				}

				d.Set(isInstanceNetworkInterfaces, interfacesList)
			}

			var rsaKey *rsa.PrivateKey
			if instance.Image != nil {
				d.Set(isInstanceImage, *instance.Image.ID)
				image := *instance.Image.Name
				res := strings.Contains(image, "windows")
				if res {
					if privatekey, ok := d.GetOk(isInstancePEM); ok {
						keyFlag := privatekey.(string)
						keybytes := []byte(keyFlag)

						if keyFlag != "" {
							block, err := pem.Decode(keybytes)
							if block == nil {
								return fmt.Errorf("Failed to load the private key from the given key contents. Instead of the key file path, please make sure the private key is pem format")
							}
							isEncrypted := false
							switch block.Type {
							case "RSA PRIVATE KEY":
								isEncrypted = x509.IsEncryptedPEMBlock(block)
							case "OPENSSH PRIVATE KEY":
								var err error
								isEncrypted, err = isOpenSSHPrivKeyEncrypted(block.Bytes)
								if err != nil {
									return fmt.Errorf("Failed to check if the provided open ssh key is encrypted or not %s", err)
								}
							default:
								return fmt.Errorf("PEM and OpenSSH private key formats with RSA key type are supported, can not support this key file type: %s", err)
							}
							passphrase := ""
							var privateKey interface{}
							if isEncrypted {
								if pass, ok := d.GetOk(isInstancePassphrase); ok {
									passphrase = pass.(string)
								} else {
									return fmt.Errorf("Mandatory field 'passphrase' not provided")
								}
								var err error
								privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, []byte(passphrase))
								if err != nil {
									return fmt.Errorf("Fail to decrypting the private key: %s", err)
								}
							} else {
								var err error
								privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, nil)
								if err != nil {
									return fmt.Errorf("Fail to decrypting the private key: %s", err)
								}
							}
							var ok bool
							rsaKey, ok = privateKey.(*rsa.PrivateKey)
							if !ok {
								return fmt.Errorf("Failed to convert to RSA private key")
							}
						}
					}
				}
			}

			getInstanceInitializationOptions := &vpcclassicv1.GetInstanceInitializationOptions{
				ID: &id,
			}
			initParms, response, err := sess.GetInstanceInitialization(getInstanceInitializationOptions)
			if err != nil {
				return fmt.Errorf("Error Getting instance Initialization: %s\n%s", err, response)
			}
			if initParms.Keys != nil {
				initKeyList := make([]map[string]interface{}, 0)
				for _, key := range initParms.Keys {
					key := key.(*vpcclassicv1.KeyReferenceInstanceInitializationContext)
					initKey := map[string]interface{}{}
					id := ""
					if key.ID != nil {
						id = *key.ID
					}
					initKey["id"] = id
					name := ""
					if key.Name != nil {
						name = *key.Name
					}
					initKey["name"] = name
					initKeyList = append(initKeyList, initKey)
					break

				}
				d.Set(isInstanceInitKeys, initKeyList)
			}
			if initParms.Password != nil && initParms.Password.EncryptedPassword != nil {
				ciphertext := *initParms.Password.EncryptedPassword
				password := base64.StdEncoding.EncodeToString(ciphertext)
				if rsaKey != nil {
					rng := rand.Reader
					clearPassword, err := rsa.DecryptOAEP(sha256.New(), rng, rsaKey, ciphertext, nil)
					if err != nil {
						return fmt.Errorf("Can not decrypt the password with the given key, %s", err)
					}
					password = string(clearPassword)
				}
				d.Set(isInstanceInitPassword, password)
			}

			d.Set(isInstanceStatus, *instance.Status)
			d.Set(isInstanceVPC, *instance.VPC.ID)
			d.Set(isInstanceZone, *instance.Zone.Name)

			var volumes []string
			volumes = make([]string, 0)
			if instance.VolumeAttachments != nil {
				for _, volume := range instance.VolumeAttachments {
					if volume.Volume != nil && *volume.Volume.ID != *instance.BootVolumeAttachment.Volume.ID {
						volumes = append(volumes, *volume.Volume.ID)
					}
				}
			}
			d.Set(isInstanceVolumes, newStringSet(schema.HashString, volumes))
			if instance.VolumeAttachments != nil {
				volList := make([]map[string]interface{}, 0)
				for _, volume := range instance.VolumeAttachments {
					vol := map[string]interface{}{}
					if volume.Volume != nil {
						vol["id"] = *volume.ID
						vol["volume_id"] = *volume.Volume.ID
						vol["name"] = *volume.Name
						vol["volume_name"] = *volume.Volume.Name
						vol["volume_crn"] = *volume.Volume.CRN
						volList = append(volList, vol)
					}
				}
				d.Set(isInstanceVolumeAttachments, volList)
			}
			if instance.BootVolumeAttachment != nil {
				bootVolList := make([]map[string]interface{}, 0)
				bootVol := map[string]interface{}{}
				bootVol["id"] = *instance.BootVolumeAttachment.ID
				bootVol["name"] = *instance.BootVolumeAttachment.Name
				if instance.BootVolumeAttachment.Device != nil {
					bootVol["device"] = *instance.BootVolumeAttachment.Device.ID
				}
				if instance.BootVolumeAttachment.Volume != nil {
					bootVol["volume_name"] = *instance.BootVolumeAttachment.Volume.Name
					bootVol["volume_id"] = *instance.BootVolumeAttachment.Volume.ID
					bootVol["volume_crn"] = *instance.BootVolumeAttachment.Volume.CRN
				}
				bootVolList = append(bootVolList, bootVol)
				d.Set(isInstanceBootVolume, bootVolList)
			}
			tags, err := GetTagsUsingCRN(meta, *instance.CRN)
			if err != nil {
				log.Printf(
					"Error on get of resource vpc Instance (%s) tags: %s", d.Id(), err)
			}
			d.Set(isInstanceTags, tags)

			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc/compute/vs")
			d.Set(ResourceName, instance.Name)
			d.Set(ResourceCRN, instance.CRN)
			d.Set(ResourceStatus, instance.Status)
			if instance.ResourceGroup != nil {
				rsMangClient, err := meta.(ClientSession).ResourceManagementAPIv2()
				if err != nil {
					return err
				}
				grp, err := rsMangClient.ResourceGroup().Get(*instance.ResourceGroup.ID)
				if err != nil {
					return err
				}
				d.Set(ResourceGroupName, grp.Name)
				d.Set(isInstanceResourceGroup, instance.ResourceGroup.ID)
			}
			return nil
		}
	}
	return fmt.Errorf("No Instance found with name %s", name)
}

func instanceGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.Instance{}
	for {
		listInstancesOptions := &vpcv1.ListInstancesOptions{}
		if start != "" {
			listInstancesOptions.Start = &start
		}
		instances, response, err := sess.ListInstances(listInstancesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Instances %s\n%s", err, response)
		}
		start = GetNext(instances.Next)
		allrecs = append(allrecs, instances.Instances...)
		if start == "" {
			break
		}
	}
	for _, instance := range allrecs {
		if *instance.Name == name {
			d.SetId(*instance.ID)
			id := *instance.ID
			d.Set(isInstanceName, *instance.Name)
			if instance.Profile != nil {
				d.Set(isInstanceProfile, *instance.Profile.Name)
			}
			cpuList := make([]map[string]interface{}, 0)
			if instance.Vcpu != nil {
				currentCPU := map[string]interface{}{}
				currentCPU[isInstanceCPUArch] = *instance.Vcpu.Architecture
				currentCPU[isInstanceCPUCount] = *instance.Vcpu.Count
				cpuList = append(cpuList, currentCPU)
			}
			d.Set(isInstanceCPU, cpuList)

			d.Set(isInstanceMemory, *instance.Memory)
			gpuList := make([]map[string]interface{}, 0)
			d.Set(isInstanceGpu, gpuList)

			if instance.Disks != nil {
				err = d.Set(isInstanceDisks, dataSourceInstanceFlattenDisks(instance.Disks))
				if err != nil {
					return fmt.Errorf("Error setting disks %s", err)
				}
			}

			if instance.PrimaryNetworkInterface != nil {
				primaryNicList := make([]map[string]interface{}, 0)
				currentPrimNic := map[string]interface{}{}
				currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
				currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name
				currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
				getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         instance.PrimaryNetworkInterface.ID,
				}
				insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
				if err != nil {
					return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
				}
				currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
				if len(insnic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(insnic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
				}

				primaryNicList = append(primaryNicList, currentPrimNic)
				d.Set(isInstancePrimaryNetworkInterface, primaryNicList)
			}

			if instance.NetworkInterfaces != nil {
				interfacesList := make([]map[string]interface{}, 0)
				for _, intfc := range instance.NetworkInterfaces {
					if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
						currentNic := map[string]interface{}{}
						currentNic["id"] = *intfc.ID
						currentNic[isInstanceNicName] = *intfc.Name
						currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
						getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
							InstanceID: &id,
							ID:         intfc.ID,
						}
						insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
						if err != nil {
							return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
						}
						currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
						if len(insnic.SecurityGroups) != 0 {
							secgrpList := []string{}
							for i := 0; i < len(insnic.SecurityGroups); i++ {
								secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
							}
							currentNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
						}
						interfacesList = append(interfacesList, currentNic)

					}
				}

				d.Set(isInstanceNetworkInterfaces, interfacesList)
			}

			var rsaKey *rsa.PrivateKey
			if instance.Image != nil {
				d.Set(isInstanceImage, *instance.Image.ID)
				image := *instance.Image.Name
				res := strings.Contains(image, "windows")
				if res {
					if privatekey, ok := d.GetOk(isInstancePEM); ok {
						keyFlag := privatekey.(string)
						keybytes := []byte(keyFlag)

						if keyFlag != "" {
							block, err := pem.Decode(keybytes)
							if block == nil {
								return fmt.Errorf("Failed to load the private key from the given key contents. Instead of the key file path, please make sure the private key is pem format")
							}
							isEncrypted := false
							switch block.Type {
							case "RSA PRIVATE KEY":
								isEncrypted = x509.IsEncryptedPEMBlock(block)
							case "OPENSSH PRIVATE KEY":
								var err error
								isEncrypted, err = isOpenSSHPrivKeyEncrypted(block.Bytes)
								if err != nil {
									return fmt.Errorf("Failed to check if the provided open ssh key is encrypted or not %s", err)
								}
							default:
								return fmt.Errorf("PEM and OpenSSH private key formats with RSA key type are supported, can not support this key file type: %s", err)
							}
							passphrase := ""
							var privateKey interface{}
							if isEncrypted {
								if pass, ok := d.GetOk(isInstancePassphrase); ok {
									passphrase = pass.(string)
								} else {
									return fmt.Errorf("Mandatory field 'passphrase' not provided")
								}
								var err error
								privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, []byte(passphrase))
								if err != nil {
									return fmt.Errorf("Fail to decrypting the private key: %s", err)
								}
							} else {
								var err error
								privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, nil)
								if err != nil {
									return fmt.Errorf("Fail to decrypting the private key: %s", err)
								}
							}
							var ok bool
							rsaKey, ok = privateKey.(*rsa.PrivateKey)
							if !ok {
								return fmt.Errorf("Failed to convert to RSA private key")
							}
						}
					}
				}
			}

			getInstanceInitializationOptions := &vpcv1.GetInstanceInitializationOptions{
				ID: &id,
			}
			initParms, response, err := sess.GetInstanceInitialization(getInstanceInitializationOptions)
			if err != nil {
				return fmt.Errorf("Error Getting instance Initialization: %s\n%s", err, response)
			}
			if initParms.Keys != nil {
				initKeyList := make([]map[string]interface{}, 0)
				for _, key := range initParms.Keys {
					key := key.(*vpcv1.KeyReferenceInstanceInitializationContext)
					initKey := map[string]interface{}{}
					id := ""
					if key.ID != nil {
						id = *key.ID
					}
					initKey["id"] = id
					name := ""
					if key.Name != nil {
						name = *key.Name
					}
					initKey["name"] = name
					initKeyList = append(initKeyList, initKey)
					break

				}
				d.Set(isInstanceInitKeys, initKeyList)
			}
			if initParms.Password != nil && initParms.Password.EncryptedPassword != nil {
				ciphertext := *initParms.Password.EncryptedPassword
				password := base64.StdEncoding.EncodeToString(ciphertext)
				if rsaKey != nil {
					rng := rand.Reader
					clearPassword, err := rsa.DecryptPKCS1v15(rng, rsaKey, ciphertext)
					if err != nil {
						return fmt.Errorf("Can not decrypt the password with the given key, %s", err)
					}
					password = string(clearPassword)
				}
				d.Set(isInstanceInitPassword, password)
			}

			d.Set(isInstanceStatus, *instance.Status)
			d.Set(isInstanceVPC, *instance.VPC.ID)
			d.Set(isInstanceZone, *instance.Zone.Name)

			var volumes []string
			volumes = make([]string, 0)
			if instance.VolumeAttachments != nil {
				for _, volume := range instance.VolumeAttachments {
					if volume.Volume != nil && *volume.Volume.ID != *instance.BootVolumeAttachment.Volume.ID {
						volumes = append(volumes, *volume.Volume.ID)
					}
				}
			}
			d.Set(isInstanceVolumes, newStringSet(schema.HashString, volumes))
			if instance.VolumeAttachments != nil {
				volList := make([]map[string]interface{}, 0)
				for _, volume := range instance.VolumeAttachments {
					vol := map[string]interface{}{}
					if volume.Volume != nil {
						vol["id"] = *volume.ID
						vol["volume_id"] = *volume.Volume.ID
						vol["name"] = *volume.Name
						vol["volume_name"] = *volume.Volume.Name
						vol["volume_crn"] = *volume.Volume.CRN
						volList = append(volList, vol)
					}
				}
				d.Set(isInstanceVolumeAttachments, volList)
			}
			if instance.BootVolumeAttachment != nil {
				bootVolList := make([]map[string]interface{}, 0)
				bootVol := map[string]interface{}{}
				bootVol["id"] = *instance.BootVolumeAttachment.ID
				bootVol["name"] = *instance.BootVolumeAttachment.Name
				if instance.BootVolumeAttachment.Device != nil {
					bootVol["device"] = *instance.BootVolumeAttachment.Device.ID
				}
				if instance.BootVolumeAttachment.Volume != nil {
					bootVol["volume_name"] = *instance.BootVolumeAttachment.Volume.Name
					bootVol["volume_id"] = *instance.BootVolumeAttachment.Volume.ID
					bootVol["volume_crn"] = *instance.BootVolumeAttachment.Volume.CRN
				}
				bootVolList = append(bootVolList, bootVol)
				d.Set(isInstanceBootVolume, bootVolList)
			}
			tags, err := GetTagsUsingCRN(meta, *instance.CRN)
			if err != nil {
				log.Printf(
					"Error on get of resource vpc Instance (%s) tags: %s", d.Id(), err)
			}
			d.Set(isInstanceTags, tags)

			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc-ext/compute/vs")
			d.Set(ResourceName, instance.Name)
			d.Set(ResourceCRN, instance.CRN)
			d.Set(ResourceStatus, instance.Status)
			if instance.ResourceGroup != nil {
				d.Set(isInstanceResourceGroup, instance.ResourceGroup.ID)
				d.Set(ResourceGroupName, instance.ResourceGroup.Name)
			}
			return nil
		}
	}
	return fmt.Errorf("No Instance found with name %s", name)
}

const opensshv1Magic = "openssh-key-v1"

type opensshPrivateKey struct {
	CipherName   string
	KdfName      string
	KdfOpts      string
	NumKeys      uint32
	PubKey       string
	PrivKeyBlock string
}

func isOpenSSHPrivKeyEncrypted(data []byte) (bool, error) {
	magic := append([]byte(opensshv1Magic), 0)
	if !bytes.Equal(magic, data[0:len(magic)]) {
		return false, errors.New("Invalid openssh private key format")
	}
	content := data[len(magic):]

	privKey := opensshPrivateKey{}

	if err := ssh.Unmarshal(content, &privKey); err != nil {
		return false, err
	}

	if privKey.KdfName == "none" && privKey.CipherName == "none" {
		return false, nil
	}
	return true, nil
}

func dataSourceInstanceFlattenDisks(result []vpcv1.InstanceDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceInstanceDisksToMap(disksItem))
	}

	return disks
}

func dataSourceInstanceDisksToMap(disksItem vpcv1.InstanceDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.CreatedAt != nil {
		disksMap["created_at"] = disksItem.CreatedAt.String()
	}
	if disksItem.Href != nil {
		disksMap["href"] = disksItem.Href
	}
	if disksItem.ID != nil {
		disksMap["id"] = disksItem.ID
	}
	if disksItem.InterfaceType != nil {
		disksMap["interface_type"] = disksItem.InterfaceType
	}
	if disksItem.Name != nil {
		disksMap["name"] = disksItem.Name
	}
	if disksItem.ResourceType != nil {
		disksMap["resource_type"] = disksItem.ResourceType
	}
	if disksItem.Size != nil {
		disksMap["size"] = disksItem.Size
	}

	return disksMap
}
