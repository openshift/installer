package lexmodels

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceBotAlias() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBotAliasRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validBotName,
			},
			"bot_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validBotAliasName,
			},
		},
	}
}

func dataSourceBotAliasRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LexModelsConn

	botName := d.Get("bot_name").(string)
	botAliasName := d.Get("name").(string)
	d.SetId(fmt.Sprintf("%s:%s", botName, botAliasName))

	resp, err := conn.GetBotAlias(&lexmodelbuildingservice.GetBotAliasInput{
		BotName: aws.String(botName),
		Name:    aws.String(botAliasName),
	})
	if err != nil {
		return fmt.Errorf("error reading Lex bot alias (%s): %w", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Service:   "lex",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("bot:%s", d.Id()),
	}
	d.Set("arn", arn.String())

	d.Set("bot_name", resp.BotName)
	d.Set("bot_version", resp.BotVersion)
	d.Set("checksum", resp.Checksum)
	d.Set("created_date", resp.CreatedDate.Format(time.RFC3339))
	d.Set("description", resp.Description)
	d.Set("last_updated_date", resp.LastUpdatedDate.Format(time.RFC3339))
	d.Set("name", resp.Name)

	return nil
}
