package alicloud

import (
	"log"
	"regexp"
	"strings"
	"time"

	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOssBucketObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOssBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_control": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_disposition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_encoding": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_md5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_side_encryption": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sse_kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOssBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	bucketName := d.Get("bucket_name").(string)

	// List bucket objects
	var initialOptions []oss.Option
	if v, ok := d.GetOk("key_prefix"); ok && v.(string) != "" {
		keyPrefix := v.(string)
		initialOptions = append(initialOptions, oss.Prefix(keyPrefix))
	}
	var requestInfo *oss.Client
	var allObjects []oss.ObjectProperties
	nextMarker := ""
	for {
		var options []oss.Option
		options = append(options, initialOptions...)
		if nextMarker != "" {
			options = append(options, oss.Marker(nextMarker))
		}

		raw, err := client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			requestInfo = &bucket.Client
			return bucket.ListObjects(options...)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_oss_bucket_object", "ListObjects", AliyunOssGoSdk)
		}
		if debugOn() {
			addDebug("ListObjects", raw, requestInfo, map[string]interface{}{"options": options})
		}
		response, _ := raw.(oss.ListObjectsResult)

		if response.Objects == nil || len(response.Objects) < 1 {
			break
		}

		if len(response.Objects) < 1 {
			break
		}

		allObjects = append(allObjects, response.Objects...)

		nextMarker = response.NextMarker
		if nextMarker == "" {
			break
		}
	}

	var filteredObjectsTemp []oss.ObjectProperties
	keyRegex, ok := d.GetOk("key_regex")
	if ok && keyRegex.(string) != "" {
		var r *regexp.Regexp
		if keyRegex != "" {
			r = regexp.MustCompile(keyRegex.(string))
		}
		for _, object := range allObjects {
			if r != nil && !r.MatchString(object.Key) {
				continue
			}
			filteredObjectsTemp = append(filteredObjectsTemp, object)
		}
	} else {
		filteredObjectsTemp = allObjects
	}

	return bucketObjectsDescriptionAttributes(d, bucketName, filteredObjectsTemp, meta)
}

func bucketObjectsDescriptionAttributes(d *schema.ResourceData, bucketName string, objects []oss.ObjectProperties, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var s []map[string]interface{}
	var requestInfo *oss.Client
	for _, object := range objects {
		mapping := map[string]interface{}{
			"key":                    object.Key,
			"etag":                   strings.Trim(object.ETag, `"`),
			"storage_class":          object.StorageClass,
			"last_modification_time": object.LastModified.Format(time.RFC3339),
		}

		// Add metadata information
		raw, err := client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			requestInfo = &bucket.Client
			return bucket.GetObjectDetailedMeta(object.Key)
		})
		if err != nil {
			log.Printf("[ERROR] Unable to get metadata for the object %s: %v", object.Key, err)
		} else {
			objectHeader, _ := raw.(http.Header)
			mapping["content_type"] = objectHeader.Get("Content-Type")
			mapping["content_length"] = objectHeader.Get("Content-Length")
			mapping["cache_control"] = objectHeader.Get("Cache-Control")
			mapping["content_disposition"] = objectHeader.Get("Content-Disposition")
			mapping["content_encoding"] = objectHeader.Get("Content-Encoding")
			mapping["content_md5"] = objectHeader.Get("Content-Md5")
			mapping["expires"] = objectHeader.Get("Expires")
			mapping["server_side_encryption"] = objectHeader.Get(oss.HTTPHeaderOssServerSideEncryption)
			mapping["sse_kms_key_id"] = objectHeader.Get(oss.HTTPHeaderOssServerSideEncryptionKeyID)
		}
		if debugOn() {
			addDebug("GetObjectDetailedMeta", raw, requestInfo, map[string]string{"objectKey": object.Key})
		}
		// Add ACL information
		raw, err = client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			requestInfo = &bucket.Client
			return bucket.GetObjectACL(object.Key)
		})
		if err != nil {
			log.Printf("[ERROR] Unable to get ACL for the object %s: %v", object.Key, err)
		} else {
			objectACL, _ := raw.(oss.GetObjectACLResult)
			mapping["acl"] = objectACL.ACL
		}
		if debugOn() {
			addDebug("GetObjectACL", raw, requestInfo, map[string]string{"objectKey": object.Key})
		}

		ids = append(ids, object.Key)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("objects", s); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
