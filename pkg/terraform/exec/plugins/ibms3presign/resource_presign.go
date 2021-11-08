package ibms3presign

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var storageClasses = []string{
	"standard", "vault", "cold", "flex", "smart",
}

var endpointTypes = []string{
	"public", "private", "direct",
}

func resourcePresign() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMCOSBucketObjectPresignCreate,
		Read:   resourceIBMCOSBucketObjectPresignCreate,
		Delete: resourceIBMCOSBucketObjectPresignDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "access_key_id",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "secret_access_key",
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS Bucket name",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS object key",
			},
			"region_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region Location info.",
			},
			"storage_class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(storageClasses),
				ForceNew:     true,
				Description:  "Storage class info",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(endpointTypes),
				Description:  "public, private or direct",
				Default:      "public",
			},
			"presigned_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
				Description: "Presigned URL",
			},
			"expire": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Set expire time for URL in minutes",
				Default:     60,
			},
		},
	}
}

func validateAllowedStringValue(validValues []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		input := v.(string)
		existed := false
		for _, s := range validValues {
			if s == input {
				existed = true
				break
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a value from %#v, got %q",
				k, validValues, input))
		}
		return

	}
}

// resourceIBMCOSBucketObjectPresignDelete is just a placeholder NOP function to satisfy the schema resource
func resourceIBMCOSBucketObjectPresignDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIBMCOSBucketObjectPresignCreate(d *schema.ResourceData, m interface{}) error {
	storageClass := d.Get("storage_class").(string)
	bucketName := d.Get("bucket_name").(string)
	object := d.Get("key").(string)
	regionLocation := d.Get("region_location").(string)
	endpointType := d.Get("endpoint_type").(string)
	accessKey := d.Get("access_key_id").(string)
	secretKey := d.Get("secret_access_key").(string)
	expire := d.Get("expire").(int)
	var svcEndpoint string
	if endpointType != "public" {
		regionLocation = fmt.Sprintf("%s.%s", endpointType, regionLocation)
	}
	svcEndpoint = fmt.Sprintf("https://s3.%s.cloud-object-storage.appdomain.cloud", regionLocation)
	region := fmt.Sprintf("%s-%s", regionLocation, storageClass)
	conf := aws.NewConfig().
		WithRegion(region).
		WithEndpoint(svcEndpoint).
		WithS3ForcePathStyle(true).
		WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, ""))
	sess := session.Must(session.NewSession()) // Creating a new session
	client := s3.New(sess, conf)
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(object),
	})
	presignedURL, err := req.Presign(time.Duration(expire) * time.Minute)
	if err != nil {
		return err
	}
	d.Set("presigned_url", presignedURL)
	d.SetId(fmt.Sprintf("%s/%s", bucketName, object))
	return nil
}
