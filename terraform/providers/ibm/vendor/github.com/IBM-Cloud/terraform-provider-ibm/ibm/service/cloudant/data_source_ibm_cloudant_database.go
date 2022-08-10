// Copyright IBM Corp. 2021, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudant

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
)

func DataSourceIBMCloudantDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCloudantDatabaseRead,

		Schema: map[string]*schema.Schema{
			"db": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path parameter to specify the database name.",
			},
			"instance_crn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloudant Instance CRN.",
			},
			"cluster": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Database cluster information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replicas": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of replicas of a database in a cluster.",
						},
						"shards": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of shards in a database. Each shard is a partition of the hash value range.",
						},
						"read_quorum": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Read quorum. The number of consistent copies of a document that need to be read before a successful reply.",
						},
						"write_quorum": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Write quorum. The number of copies of a document that need to be written before a successful reply.",
						},
					},
				},
			},
			"committed_update_seq": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An opaque string that describes the committed state of the database.",
			},
			"compact_running": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the database compaction routine is operating on this database.",
			},
			"compacted_seq": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An opaque string that describes the compaction state of the database.",
			},
			"disk_format_version": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the physical format used for the data when it is stored on disk.",
			},
			"doc_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A count of the documents in the specified database.",
			},
			"doc_del_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of deleted documents.",
			},
			"engine": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The engine used for the database.",
			},
			"props": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The database properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partitioned": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The value is `true` for a partitioned database.",
						},
					},
				},
			},
			"sizes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Database size information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The active size of the data in the database, in bytes.",
						},
						"external": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total uncompressed size of the data in the database, in bytes.",
						},
						"file": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total size of the database as stored on disk, in bytes.",
						},
					},
				},
			},
			"update_seq": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An opaque string that describes the state of the database. Do not rely on this string for counting the number of updates.",
			},
			"uuid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the database.",
			},
		},
	}
}

func dataSourceIBMCloudantDatabaseRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceCRN := d.Get("instance_crn").(string)
	cUrl, err := GetCloudantInstanceUrl(instanceCRN, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudantClient, err := GetCloudantClientForUrl(cUrl, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	dbName := d.Get("db").(string)
	getDatabaseInformationOptions := cloudantClient.NewGetDatabaseInformationOptions(dbName)

	databaseInformation, response, err := cloudantClient.GetDatabaseInformationWithContext(context, getDatabaseInformationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetDatabaseInformationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDatabaseInformationWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMCloudantDatabaseID(d))

	if databaseInformation.Cluster != nil {
		err = d.Set("cluster", dataSourceDatabaseInformationFlattenCluster(*databaseInformation.Cluster))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cluster %s", err))
		}
	}

	if databaseInformation.CommittedUpdateSeq != nil {
		d.Set("committed_update_seq", *databaseInformation.CommittedUpdateSeq)
	}
	if databaseInformation.CompactRunning != nil {
		d.Set("compact_running", *databaseInformation.CompactRunning)
	}
	if databaseInformation.CompactedSeq != nil {
		d.Set("compacted_seq", *databaseInformation.CompactedSeq)
	}
	if databaseInformation.DiskFormatVersion != nil {
		d.Set("disk_format_version", *databaseInformation.DiskFormatVersion)
	}
	if databaseInformation.DocCount != nil {
		d.Set("doc_count", *databaseInformation.DocCount)
	}
	if databaseInformation.DocDelCount != nil {
		d.Set("doc_del_count", *databaseInformation.DocDelCount)
	}
	if databaseInformation.Engine != nil {
		d.Set("engine", *databaseInformation.Engine)
	}
	if databaseInformation.Props != nil {
		d.Set("props", dataSourceDatabaseInformationFlattenProps(*databaseInformation.Props))
	}
	if databaseInformation.Sizes != nil {
		d.Set("sizes", dataSourceDatabaseInformationFlattenSizes(*databaseInformation.Sizes))
	}
	if databaseInformation.UpdateSeq != nil {
		d.Set("update_seq", *databaseInformation.UpdateSeq)
	}
	if databaseInformation.UUID != nil {
		d.Set("uuid", *databaseInformation.UUID)
	}

	return nil
}

// dataSourceIBMCloudantDatabaseID returns a reasonable ID for the list.
func dataSourceIBMCloudantDatabaseID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDatabaseInformationFlattenCluster(result cloudantv1.DatabaseInformationCluster) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDatabaseInformationClusterToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDatabaseInformationClusterToMap(clusterItem cloudantv1.DatabaseInformationCluster) (clusterMap map[string]interface{}) {
	clusterMap = map[string]interface{}{}

	if clusterItem.N != nil {
		clusterMap["replicas"] = *clusterItem.N
	}
	if clusterItem.Q != nil {
		clusterMap["shards"] = *clusterItem.Q
	}
	if clusterItem.R != nil {
		clusterMap["read_quorum"] = *clusterItem.R
	}
	if clusterItem.W != nil {
		clusterMap["write_quorum"] = *clusterItem.W
	}

	return clusterMap
}

func dataSourceDatabaseInformationFlattenProps(result cloudantv1.DatabaseInformationProps) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDatabaseInformationPropsToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDatabaseInformationPropsToMap(propsItem cloudantv1.DatabaseInformationProps) (propsMap map[string]interface{}) {
	propsMap = map[string]interface{}{}

	if propsItem.Partitioned != nil {
		propsMap["partitioned"] = *propsItem.Partitioned
	} else {
		propsMap["partitioned"] = false
	}

	return propsMap
}

func dataSourceDatabaseInformationFlattenSizes(result cloudantv1.ContentInformationSizes) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDatabaseInformationSizesToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDatabaseInformationSizesToMap(sizesItem cloudantv1.ContentInformationSizes) (sizesMap map[string]interface{}) {
	sizesMap = map[string]interface{}{}

	if sizesItem.Active != nil {
		sizesMap["active"] = *sizesItem.Active
	}
	if sizesItem.External != nil {
		sizesMap["external"] = *sizesItem.External
	}
	if sizesItem.File != nil {
		sizesMap["file"] = *sizesItem.File
	}

	return sizesMap
}
