package cos

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	token "github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam/token"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMCOSBucketReplicationConfiguration() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCOSBucketReplicationConfigurationCreate,
		Read:     resourceIBMCOSBucketReplicationConfigurationRead,
		Update:   resourceIBMCOSBucketReplicationConfigurationUpdate,
		Delete:   resourceIBMCOSBucketReplicationConfigurationDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS bucket CRN",
			},
			"bucket_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS bucket location",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private", "direct"}),
				Description:  "COS endpoint type: public, private, direct",
				Default:      "public",
			},
			"replication_rule": {
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1000,
				Description: "Replicate objects between buckets, replicate across source and destination. A container for replication rules can add up to 1,000 rules. The maximum size of a replication configuration is 2 MB.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
							Description:  "A unique identifier for the rule. The maximum value is 255 characters.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "A priority is associated with each rule. There may be cases where multiple rules may be applicable to an object that is uploaded. ",
						},
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable or disable an replication rule for a bucket",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The rule applies to any objects with keys that match this prefix",
						},
						"deletemarker_replication_status": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether to replicate delete markers. It should be either Enable or Disable",
						},
						"destination_bucket_crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Cloud Resource Name (CRN) of the bucket where you want COS to store the results",
						},
					},
				},
				Set: resourceIBMCOSReplicationReuleHash,
			},
		},
	}
}

func replicationRuleSet(replicateList []interface{}) []*s3.ReplicationRule {

	var rules []*s3.ReplicationRule
	for _, l := range replicateList {
		bkt_replication_rule := s3.ReplicationRule{}
		replicateMap, _ := l.(map[string]interface{})

		//Rule ID
		if rule_idSet, exist := replicateMap["rule_id"]; exist {
			id := rule_idSet.(string)
			bkt_replication_rule.ID = aws.String(id)
		}

		//Status Enable/Disable
		if statusSet, exist := replicateMap["enable"]; exist {
			StatusEnabled := statusSet.(bool)
			if StatusEnabled == true {
				bkt_replication_rule.Status = aws.String("Enabled")
			} else {
				bkt_replication_rule.Status = aws.String("Disabled")
			}
		}
		//Days
		if priorSet, exist := replicateMap["priority"]; exist {
			replicate_priority := int64(priorSet.(int))
			bkt_replication_rule.Priority = aws.Int64(replicate_priority)
		}
		//Replication Prefix
		if PrefixClassSet, exist := replicateMap["prefix"]; exist {
			prefix_check := PrefixClassSet.(string)
			bkt_replication_rule.Filter = &s3.ReplicationRuleFilter{Prefix: aws.String(prefix_check)}

		}
		//DeleteMarkerReplicationStatus
		if delMarkerStatusSet, exist := replicateMap["deletemarker_replication_status"]; exist {
			del_marker_status_value := delMarkerStatusSet.(bool)
			if del_marker_status_value == true {
				bkt_replication_rule.DeleteMarkerReplication = &s3.DeleteMarkerReplication{Status: aws.String("Enabled")}
			} else {
				bkt_replication_rule.DeleteMarkerReplication = &s3.DeleteMarkerReplication{Status: aws.String("Disabled")}
			}
		}
		//Destination CRN
		if dest_bucket_CrnSet, exist := replicateMap["destination_bucket_crn"]; exist {
			dest_bucketcrn_value := dest_bucket_CrnSet.(string)
			bkt_replication_rule.Destination = &s3.Destination{Bucket: aws.String(dest_bucketcrn_value)}
		}

		rules = append(rules, &bkt_replication_rule)

	}
	return rules
}

func resourceIBMCOSBucketReplicationConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])

	bucketLocation := d.Get("bucket_location").(string)
	endpointType := d.Get("endpoint_type").(string)

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	var rules []*s3.ReplicationRule

	replication, ok := d.GetOk("replication_rule")
	if ok {
		rules = append(rules, replicationRuleSet(replication.(*schema.Set).List())...)

	}
	putBucketReplicationInput := &s3.PutBucketReplicationInput{
		Bucket: aws.String(bucketName),
		ReplicationConfiguration: &s3.ReplicationConfiguration{
			Rules: rules,
		},
	}

	_, err = s3Client.PutBucketReplication(putBucketReplicationInput)

	if err != nil {
		return fmt.Errorf("failed to create the replication rule on COS bucket %s, %v", bucketName, err)
	}

	//Generating a fake id which contains every information about to get the bucket via s3 api
	bktID := fmt.Sprintf("%s:%s:%s:meta:%s:%s", strings.Replace(instanceCRN, "::", "", -1), "bucket", bucketName, bucketLocation, endpointType)
	d.SetId(bktID)

	return resourceIBMCOSBucketReplicationConfigurationRead(d, meta)
}

func resourceIBMCOSBucketReplicationConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])

	bucketLocation := d.Get("bucket_location").(string)
	endpointType := d.Get("endpoint_type").(string)

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return err
	}

	if d.HasChange("replication_rule") {
		var rules []*s3.ReplicationRule

		replication, ok := d.GetOk("replication_rule")
		if ok {
			rules = append(rules, replicationRuleSet(replication.(*schema.Set).List())...)

		}
		putBucketReplication := &s3.PutBucketReplicationInput{
			Bucket: aws.String(bucketName),
			ReplicationConfiguration: &s3.ReplicationConfiguration{
				Rules: rules,
			},
		}

		_, err = s3Client.PutBucketReplication(putBucketReplication)

		if err != nil {
			return fmt.Errorf("failed to update the replication rule on COS bucket %s, %v", bucketName, err)
		}

	}
	return resourceIBMCOSBucketReplicationConfigurationRead(d, meta)
}

func resourceIBMCOSBucketReplicationConfigurationRead(d *schema.ResourceData, meta interface{}) error {

	bucketCRN := parseBucketReplId(d.Id(), "bucketCRN")
	bucketName := parseBucketReplId(d.Id(), "bucketName")
	bucketLocation := parseBucketReplId(d.Id(), "bucketLocation")
	instanceCRN := parseBucketReplId(d.Id(), "instanceCRN")
	endpointType := parseBucketReplId(d.Id(), "endpointType")

	d.Set("bucket_crn", bucketCRN)
	d.Set("bucket_location", bucketLocation)
	if endpointType != "" {
		d.Set("endpoint_type", endpointType)
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return err
	}

	getBucketReplicationInput := &s3.GetBucketReplicationInput{
		Bucket: aws.String(bucketName),
	}

	replicationptr, err := s3Client.GetBucketReplication(getBucketReplicationInput)

	if err != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return err
	}

	if replicationptr != nil {
		replicationRules := flex.ReplicationRuleGet(replicationptr.ReplicationConfiguration)
		if len(replicationRules) > 0 {
			d.Set("replication_rule", replicationRules)
		}
	}

	return nil
}

func resourceIBMCOSBucketReplicationConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	bucketName := parseBucketReplId(d.Id(), "bucketName")
	bucketLocation := parseBucketReplId(d.Id(), "bucketLocation")
	instanceCRN := parseBucketReplId(d.Id(), "instanceCRN")
	endpointType := parseBucketReplId(d.Id(), "endpointType")

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return err
	}

	deleteBucketReplicationInput := &s3.DeleteBucketReplicationInput{
		Bucket: aws.String(bucketName),
	}

	_, err = s3Client.DeleteBucketReplication(deleteBucketReplicationInput)

	if err != nil {
		return err
	}
	return nil
}

func parseBucketReplId(id string, info string) string {
	bucketCRN := strings.Split(id, ":meta:")[0]
	meta := strings.Split(id, ":meta:")[1]

	if info == "bucketName" {
		return strings.Split(bucketCRN, ":bucket:")[1]
	}
	if info == "instanceCRN" {
		return fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	}
	if info == "bucketCRN" {
		return bucketCRN
	}
	if info == "bucketLocation" {
		return strings.Split(meta, ":")[0]
	}
	if info == "endpointType" {
		return strings.Split(meta, ":")[1]
	}

	return parseBucketId(bucketCRN, info)
}

func getCosEndpointType(bucketLocation string, endpointType string) string {

	if bucketLocation != "" {
		hostUrl := "cloud-object-storage.appdomain.cloud"
		switch endpointType {
		case "public":
			return fmt.Sprintf("s3.%s.%s", bucketLocation, hostUrl)
		case "private":
			return fmt.Sprintf("s3.private.%s.%s", bucketLocation, hostUrl)
		case "direct":
			return fmt.Sprintf("s3.direct.%s.%s", bucketLocation, hostUrl)
		default:
			return fmt.Sprintf("s3.%s.%s", bucketLocation, hostUrl)
		}
	}

	return ""
}

func getS3ClientSession(bxSession *bxsession.Session, bucketLocation string, endpointType string, instanceCRN string) (*s3.S3, error) {
	var s3Conf *aws.Config

	visibility := endpointType
	if endpointType == "direct" {
		visibility = "private"
	}
	apiEndpoint := getCosEndpointType(bucketLocation, endpointType)
	apiEndpoint = conns.FileFallBack(bxSession.Config.EndpointsFile, visibility, "IBMCLOUD_COS_ENDPOINT", bucketLocation, apiEndpoint)
	apiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	if apiEndpoint == "" {
		return nil, fmt.Errorf("the endpoint doesn't exists for given location %s and endpoint type %s", bucketLocation, endpointType)
	}

	authEndpoint, err := bxSession.Config.EndpointLocator.IAMEndpoint()
	if err != nil {
		return nil, err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")
	apiKey := bxSession.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, instanceCRN)).WithS3ForcePathStyle(true)
	}
	iamAccessToken := bxSession.Config.IAMAccessToken
	if iamAccessToken != "" {
		initFunc := func() (*token.Token, error) {
			return &token.Token{
				AccessToken:  bxSession.Config.IAMAccessToken,
				RefreshToken: bxSession.Config.IAMRefreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    int64((time.Hour * 248).Seconds()) * -1,
				Expiration:   time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		}
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, instanceCRN)).WithS3ForcePathStyle(true)
	}
	s3Sess := session.Must(session.NewSession())
	return s3.New(s3Sess, s3Conf), nil
}

func resourceIBMCOSReplicationReuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-",
		m["destination_bucket_crn"].(string)))

	return conns.String(buf.String())
}
