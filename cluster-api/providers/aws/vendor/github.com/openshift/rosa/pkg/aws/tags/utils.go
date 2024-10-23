package tags

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Ec2ResourceHasTag(ec2Tags []ec2types.Tag, tagKey string, tagValue string) bool {
	for _, tag := range ec2Tags {
		if aws.ToString(tag.Key) == tagKey && aws.ToString(tag.Value) == tagValue {
			return true
		}
	}

	return false
}
