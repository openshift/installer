package macie2

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/macie2"
)

// findMemberNotAssociated Return a list of members not associated and compare with account ID
func findMemberNotAssociated(ctx context.Context, conn *macie2.Macie2, accountID string) (*macie2.Member, error) {
	input := &macie2.ListMembersInput{
		OnlyAssociated: aws.String("false"),
	}
	var result *macie2.Member

	err := conn.ListMembersPagesWithContext(ctx, input, func(page *macie2.ListMembersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, member := range page.Members {
			if member == nil {
				continue
			}

			if aws.StringValue(member.AccountId) == accountID {
				result = member
				return false
			}
		}

		return !lastPage
	})

	return result, err
}
