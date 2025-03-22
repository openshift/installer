// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package mqcloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func ResourceIbmMqcloudKeystoreCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmMqcloudKeystoreCertificateCreate,
		ReadContext:   resourceIbmMqcloudKeystoreCertificateRead,
		UpdateContext: resourceIbmMqcloudKeystoreCertificateUpdate,
		DeleteContext: resourceIbmMqcloudKeystoreCertificateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_keystore_certificate", "service_instance_guid"),
				Description:  "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"queue_manager_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_keystore_certificate", "queue_manager_id"),
				Description:  "The id of the queue manager to retrieve its full details.",
			},
			"label": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_keystore_certificate", "label"),
				Description:  "The label to use for the certificate to be uploaded.",
			},
			"certificate_file": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The filename and path of the certificate to be uploaded.",
			},
			"certificate_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of certificate.",
			},
			"fingerprint_sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fingerprint SHA256.",
			},
			"subject_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subject's Distinguished Name.",
			},
			"subject_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subject's Common Name.",
			},
			"issuer_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer's Distinguished Name.",
			},
			"issuer_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer's Common Name.",
			},
			"issued": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date certificate was issued.",
			},
			"expiry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiry date for the certificate.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is the queue manager's default certificate.",
			},
			"dns_names_total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of dns names.",
			},
			"dns_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DNS names.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this key store certificate.",
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The configuration details for this certificate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ams": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of channels that are configured with this certificate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channels": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of channels that are configured with this certificate.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The name of the channel.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the certificate.",
			},
		},
	}
}

func ResourceIbmMqcloudKeystoreCertificateValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_instance_guid",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "queue_manager_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{32}$`,
			MinValueLength:             32,
			MaxValueLength:             32,
		},
		validate.ValidateSchema{
			Identifier:                 "label",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_.]*$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_mqcloud_keystore_certificate", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmMqcloudKeystoreCertificateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Create Keystore Certificate failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createKeyStorePemCertificateOptions := &mqcloudv1.CreateKeyStorePemCertificateOptions{}

	createKeyStorePemCertificateOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createKeyStorePemCertificateOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))
	createKeyStorePemCertificateOptions.SetLabel(d.Get("label").(string))
	certificateFileBytes, err := base64.StdEncoding.DecodeString(d.Get("certificate_file").(string))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "create", "parse-certificate_file").GetDiag()
	}
	createKeyStorePemCertificateOptions.SetCertificateFile(io.NopCloser(bytes.NewReader(certificateFileBytes)))

	keyStoreCertificateDetails, _, err := mqcloudClient.CreateKeyStorePemCertificateWithContext(context, createKeyStorePemCertificateOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateKeyStorePemCertificateWithContext failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createKeyStorePemCertificateOptions.ServiceInstanceGuid, *createKeyStorePemCertificateOptions.QueueManagerID, *keyStoreCertificateDetails.ID))

	// Update channel after creating the certificate
	setCertificateAmsChannelsOptions := &mqcloudv1.SetCertificateAmsChannelsOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "create", "sep-id-parts").GetDiag()
	}
	setCertificateAmsChannelsOptions.SetServiceInstanceGuid(parts[0])
	setCertificateAmsChannelsOptions.SetQueueManagerID(parts[1])
	setCertificateAmsChannelsOptions.SetCertificateID(parts[2])

	var channels []mqcloudv1.ChannelDetails
	channels = []mqcloudv1.ChannelDetails{}
	for _, v := range d.Get("config.0.ams.0.channels").([]interface{}) {
		value := v.(map[string]interface{})
		channelsItem, err := ResourceIbmMqcloudKeystoreCertificateMapToChannelDetails(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "create", "parse-channels").GetDiag()
		}
		channels = append(channels, *channelsItem)
	}

	setCertificateAmsChannelsOptions.SetChannels(channels)
	setCertificateAmsChannelsOptions.SetUpdateStrategy("replace")

	_, _, err = mqcloudClient.SetCertificateAmsChannelsWithContext(context, setCertificateAmsChannelsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SetCertificateAmsChannelsWithContext failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIbmMqcloudKeystoreCertificateRead(context, d, meta)
}

func resourceIbmMqcloudKeystoreCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getKeyStoreCertificateOptions := &mqcloudv1.GetKeyStoreCertificateOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "sep-id-parts").GetDiag()
	}

	getKeyStoreCertificateOptions.SetServiceInstanceGuid(parts[0])
	getKeyStoreCertificateOptions.SetQueueManagerID(parts[1])
	getKeyStoreCertificateOptions.SetCertificateID(parts[2])

	keyStoreCertificateDetails, response, err := mqcloudClient.GetKeyStoreCertificateWithContext(context, getKeyStoreCertificateOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetKeyStoreCertificateWithContext failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		err = fmt.Errorf("Error setting service_instance_guid: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-service_instance_guid").GetDiag()
	}
	if err = d.Set("queue_manager_id", parts[1]); err != nil {
		err = fmt.Errorf("Error setting queue_manager_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-queue_manager_id").GetDiag()
	}
	if err = d.Set("label", keyStoreCertificateDetails.Label); err != nil {
		err = fmt.Errorf("Error setting label: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-label").GetDiag()
	}
	if err = d.Set("certificate_type", keyStoreCertificateDetails.CertificateType); err != nil {
		err = fmt.Errorf("Error setting certificate_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-certificate_type").GetDiag()
	}
	if err = d.Set("fingerprint_sha256", keyStoreCertificateDetails.FingerprintSha256); err != nil {
		err = fmt.Errorf("Error setting fingerprint_sha256: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-fingerprint_sha256").GetDiag()
	}
	if err = d.Set("subject_dn", keyStoreCertificateDetails.SubjectDn); err != nil {
		err = fmt.Errorf("Error setting subject_dn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-subject_dn").GetDiag()
	}
	if err = d.Set("subject_cn", keyStoreCertificateDetails.SubjectCn); err != nil {
		err = fmt.Errorf("Error setting subject_cn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-subject_cn").GetDiag()
	}
	if err = d.Set("issuer_dn", keyStoreCertificateDetails.IssuerDn); err != nil {
		err = fmt.Errorf("Error setting issuer_dn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-issuer_dn").GetDiag()
	}
	if err = d.Set("issuer_cn", keyStoreCertificateDetails.IssuerCn); err != nil {
		err = fmt.Errorf("Error setting issuer_cn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-issuer_cn").GetDiag()
	}
	if err = d.Set("issued", flex.DateTimeToString(keyStoreCertificateDetails.Issued)); err != nil {
		err = fmt.Errorf("Error setting issued: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-issued").GetDiag()
	}
	if err = d.Set("expiry", flex.DateTimeToString(keyStoreCertificateDetails.Expiry)); err != nil {
		err = fmt.Errorf("Error setting expiry: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-expiry").GetDiag()
	}
	if err = d.Set("is_default", keyStoreCertificateDetails.IsDefault); err != nil {
		err = fmt.Errorf("Error setting is_default: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-is_default").GetDiag()
	}
	if err = d.Set("dns_names_total_count", flex.IntValue(keyStoreCertificateDetails.DnsNamesTotalCount)); err != nil {
		err = fmt.Errorf("Error setting dns_names_total_count: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-dns_names_total_count").GetDiag()
	}
	if err = d.Set("dns_names", keyStoreCertificateDetails.DnsNames); err != nil {
		err = fmt.Errorf("Error setting dns_names: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-dns_names").GetDiag()
	}
	if err = d.Set("href", keyStoreCertificateDetails.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-href").GetDiag()
	}
	configMap, err := ResourceIbmMqcloudKeystoreCertificateCertificateConfigurationToMap(keyStoreCertificateDetails.Config)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "config-to-map").GetDiag()
	}
	// Check if configMap is not empty and contains non-empty lists before setting it
	if len(configMap) > 0 && len(configMap["ams"].([]map[string]interface{})) > 0 && len(configMap["ams"].([]map[string]interface{})[0]["channels"].([]map[string]interface{})) > 0 {
		if err = d.Set("config", []map[string]interface{}{configMap}); err != nil {
			err = fmt.Errorf("Error setting config: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-config").GetDiag()
		}
	}
	if err = d.Set("certificate_id", keyStoreCertificateDetails.ID); err != nil {
		err = fmt.Errorf("Error setting certificate_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "read", "set-certificate_id").GetDiag()
	}

	return nil
}

func resourceIbmMqcloudKeystoreCertificateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Update Keystore Certificate failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	setCertificateAmsChannelsOptions := &mqcloudv1.SetCertificateAmsChannelsOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "update", "sep-id-parts").GetDiag()
	}

	setCertificateAmsChannelsOptions.SetServiceInstanceGuid(parts[0])
	setCertificateAmsChannelsOptions.SetQueueManagerID(parts[1])
	setCertificateAmsChannelsOptions.SetCertificateID(parts[2])

	hasChange := false

	if d.HasChange("queue_manager_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "queue_manager_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_mqcloud_keystore_certificate", "update", "queue_manager_id-forces-new").GetDiag()
	}
	if d.HasChange("service_instance_guid") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "service_instance_guid")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_mqcloud_keystore_certificate", "update", "service_instance_guid-forces-new").GetDiag()
	}
	if d.HasChange("config.0.ams.0") {
		var channels []mqcloudv1.ChannelDetails
		channels = []mqcloudv1.ChannelDetails{}
		for _, v := range d.Get("config.0.ams.0.channels").([]interface{}) {
			value := v.(map[string]interface{})
			channelsItem, err := ResourceIbmMqcloudKeystoreCertificateMapToChannelDetails(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "update", "parse-channels").GetDiag()
			}
			channels = append(channels, *channelsItem)
		}
		setCertificateAmsChannelsOptions.SetChannels(channels)
		setCertificateAmsChannelsOptions.SetUpdateStrategy("replace")
		hasChange = true
	}

	if hasChange {
		_, _, err = mqcloudClient.SetCertificateAmsChannelsWithContext(context, setCertificateAmsChannelsOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SetCertificateAmsChannelsWithContext failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmMqcloudKeystoreCertificateRead(context, d, meta)
}

func resourceIbmMqcloudKeystoreCertificateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Delete Keystore Certificate failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteKeyStoreCertificateOptions := &mqcloudv1.DeleteKeyStoreCertificateOptions{}
	setCertificateAmsChannelsOptions := &mqcloudv1.SetCertificateAmsChannelsOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_keystore_certificate", "delete", "sep-id-parts").GetDiag()
	}

	deleteKeyStoreCertificateOptions.SetServiceInstanceGuid(parts[0])
	deleteKeyStoreCertificateOptions.SetQueueManagerID(parts[1])
	deleteKeyStoreCertificateOptions.SetCertificateID(parts[2])

	setCertificateAmsChannelsOptions.SetServiceInstanceGuid(parts[0])
	setCertificateAmsChannelsOptions.SetQueueManagerID(parts[1])
	setCertificateAmsChannelsOptions.SetCertificateID(parts[2])
	channels := []mqcloudv1.ChannelDetails{}
	setCertificateAmsChannelsOptions.SetChannels(channels)
	setCertificateAmsChannelsOptions.SetUpdateStrategy("replace")

	_, _, err = mqcloudClient.SetCertificateAmsChannelsWithContext(context, setCertificateAmsChannelsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteKeyStoreCertificateWithContext failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	_, err = mqcloudClient.DeleteKeyStoreCertificateWithContext(context, deleteKeyStoreCertificateOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteKeyStoreCertificateWithContext failed: %s", err.Error()), "ibm_mqcloud_keystore_certificate", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmMqcloudKeystoreCertificateMapToChannelDetails(modelMap map[string]interface{}) (*mqcloudv1.ChannelDetails, error) {
	model := &mqcloudv1.ChannelDetails{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIbmMqcloudKeystoreCertificateCertificateConfigurationToMap(model *mqcloudv1.CertificateConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	amsMap, err := ResourceIbmMqcloudKeystoreCertificateChannelsDetailsToMap(model.Ams)
	if err != nil {
		return modelMap, err
	}
	modelMap["ams"] = []map[string]interface{}{amsMap}
	return modelMap, nil
}

func ResourceIbmMqcloudKeystoreCertificateChannelsDetailsToMap(model *mqcloudv1.ChannelsDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	channels := []map[string]interface{}{}
	for _, channelsItem := range model.Channels {
		channelsItem := channelsItem
		channelsItemMap, err := ResourceIbmMqcloudKeystoreCertificateChannelDetailsToMap(&channelsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		channels = append(channels, channelsItemMap)
	}
	modelMap["channels"] = channels
	return modelMap, nil
}

func ResourceIbmMqcloudKeystoreCertificateChannelDetailsToMap(model *mqcloudv1.ChannelDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}
