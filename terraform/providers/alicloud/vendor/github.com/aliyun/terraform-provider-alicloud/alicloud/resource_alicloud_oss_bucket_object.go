package alicloud

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func resourceAlicloudOssBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOssBucketObjectPut,
		Read:   resourceAlicloudOssBucketObjectRead,
		Update: resourceAlicloudOssBucketObjectPut,
		Delete: resourceAlicloudOssBucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},

			"acl": {
				Type:         schema.TypeString,
				Default:      oss.ACLPrivate,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"content_length": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_md5": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"expires": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"server_side_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ServerSideEncryptionKMS), string(ServerSideEncryptionAes256),
				}, false),
				Default: ServerSideEncryptionAes256,
			},

			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return ServerSideEncryptionKMS != d.Get("server_side_encryption").(string)
				},
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudOssBucketObjectPut(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.Bucket(d.Get("bucket").(string))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_object", "Bucket", AliyunOssGoSdk)
	}
	addDebug("Bucket", raw, requestInfo, map[string]string{"bucketName": d.Get("bucket").(string)})
	bucket, _ := raw.(*oss.Bucket)
	var filePath string
	var body io.Reader

	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return WrapError(err)
		}

		filePath = path
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	} else {
		return WrapError(Error("[ERROR] Must specify \"source\" or \"content\" field"))
	}

	key := d.Get("key").(string)
	options, err := buildObjectHeaderOptions(d)

	if v, ok := d.GetOk("server_side_encryption"); ok {
		options = append(options, oss.ServerSideEncryption(v.(string)))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		options = append(options, oss.ServerSideEncryptionKeyID(v.(string)))
	}

	if err != nil {
		return WrapError(err)
	}
	if filePath != "" {
		err = bucket.PutObjectFromFile(key, filePath, options...)
	}

	if body != nil {
		err = bucket.PutObject(key, body, options...)
	}

	if err != nil {
		return WrapError(Error("Error putting object in Oss bucket (%#v): %s", bucket, err))
	}

	d.SetId(key)
	return resourceAlicloudOssBucketObjectRead(d, meta)
}

func resourceAlicloudOssBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.Bucket(d.Get("bucket").(string))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "Bucket", AliyunOssGoSdk)
	}
	addDebug("Bucket", raw, requestInfo, map[string]string{"bucketName": d.Get("bucket").(string)})
	bucket, _ := raw.(*oss.Bucket)
	options, err := buildObjectHeaderOptions(d)
	if err != nil {
		return WrapError(err)
	}

	object, err := bucket.GetObjectDetailedMeta(d.Get("key").(string), options...)
	if err != nil {
		if IsExpectedErrors(err, []string{"404 Not Found"}) {
			d.SetId("")
			return WrapError(Error("To get the Object: %#v but it is not exist in the specified bucket %s.", d.Get("key").(string), d.Get("bucket").(string)))
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetObjectDetailedMeta", AliyunOssGoSdk)
	}
	addDebug("GetObjectDetailedMeta", object, requestInfo, map[string]interface{}{
		"objectKey": d.Get("key").(string),
		"options":   options,
	})

	d.Set("content_type", object.Get("Content-Type"))
	d.Set("content_length", object.Get("Content-Length"))
	d.Set("cache_control", object.Get("Cache-Control"))
	d.Set("content_disposition", object.Get("Content-Disposition"))
	d.Set("content_encoding", object.Get("Content-Encoding"))
	d.Set("expires", object.Get("Expires"))
	d.Set("etag", strings.Trim(object.Get("ETag"), `"`))
	d.Set("version_id", object.Get("x-oss-version-id"))

	return nil
}

func resourceAlicloudOssBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.Bucket(d.Get("bucket").(string))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "Bucket", AliyunOssGoSdk)
	}
	addDebug("Bucket", raw, requestInfo, map[string]string{"bucketName": d.Get("bucket").(string)})
	bucket, _ := raw.(*oss.Bucket)

	err = bucket.DeleteObject(d.Id())
	if err != nil {
		if IsExpectedErrors(err, []string{"No Content", "Not Found"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteObject", AliyunOssGoSdk)
	}

	return WrapError(ossService.WaitForOssBucketObject(bucket, d.Id(), Deleted, DefaultTimeoutMedium))

}

func buildObjectHeaderOptions(d *schema.ResourceData) (options []oss.Option, err error) {

	if v, ok := d.GetOk("acl"); ok {
		options = append(options, oss.ObjectACL(oss.ACLType(v.(string))))
	}

	if v, ok := d.GetOk("content_type"); ok {
		options = append(options, oss.ContentType(v.(string)))
	}

	if v, ok := d.GetOk("cache_control"); ok {
		options = append(options, oss.CacheControl(v.(string)))
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		options = append(options, oss.ContentDisposition(v.(string)))
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		options = append(options, oss.ContentEncoding(v.(string)))
	}

	if v, ok := d.GetOk("content_md5"); ok {
		options = append(options, oss.ContentMD5(v.(string)))
	}

	if v, ok := d.GetOk("expires"); ok {
		expires := v.(string)
		expiresTime, err := time.Parse(time.RFC1123, expires)
		if err != nil {
			return nil, fmt.Errorf("expires format must respect the RFC1123 standard (current value: %s)", expires)
		}
		options = append(options, oss.Expires(expiresTime))
	}

	if options == nil || len(options) == 0 {
		log.Printf("[WARN] Object header options is nil.")
	}
	return options, nil
}
