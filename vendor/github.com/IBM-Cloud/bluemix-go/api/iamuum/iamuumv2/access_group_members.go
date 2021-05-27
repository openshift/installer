package iamuumv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/models"
)

const (
	AccessGroupMemberUser    = "user"
	AccessGroupMemberService = "service"
)

type AddGroupMemberRequestV2 struct {
	Members []models.AccessGroupMemberV2 `json:"members"`
}

type AddGroupMemberResponseV2 struct {
	Members []AddedGroupMemberV2 `json:"members"`
}

type AddedGroupMemberV2 struct {
	ID          string  `json:"iam_id"`
	Type        string  `json:"type"`
	Href        string  `json:"href,omitempty"`
	StatusCode  int     `json:"status_code,omitempty"`
	Trace       string  `json:"trace,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	CreatedByID string  `json:"created_by_id,omitempty"`
	Errors      []Error `json:"errors,omitempty"`
}
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
type GroupMembers struct {
	PaginationFields
	Members []models.AccessGroupMemberV2 `json:"members"`
}

func (gm *GroupMembers) Resources() []interface{} {
	r := make([]interface{}, len(gm.Members))
	for i := range gm.Members {
		r[i] = gm.Members[i]
	}
	return r
}

type AccessGroupMemberRepositoryV2 interface {
	List(groupID string) ([]models.AccessGroupMemberV2, error)
	Add(groupID string, request AddGroupMemberRequestV2) (AddGroupMemberResponseV2, error)
	Remove(groupID string, memberID string) error
}

type accessGroupMemberRepository struct {
	client *client.Client
}

func NewAccessGroupMemberRepository(c *client.Client) AccessGroupMemberRepositoryV2 {
	return &accessGroupMemberRepository{
		client: c,
	}

}

func (r *accessGroupMemberRepository) List(groupID string) ([]models.AccessGroupMemberV2, error) {
	members := []models.AccessGroupMemberV2{}
	_, err := r.client.GetPaginated(fmt.Sprintf("/v2/groups/%s/members", groupID),
		NewPaginatedResourcesHandler(&GroupMembers{}), func(resource interface{}) bool {
			if member, ok := resource.(models.AccessGroupMemberV2); ok {
				members = append(members, member)
				return true
			}
			return false
		})
	if err != nil {
		return []models.AccessGroupMemberV2{}, err
	}
	return members, nil
}

func (r *accessGroupMemberRepository) Add(groupID string, request AddGroupMemberRequestV2) (AddGroupMemberResponseV2, error) {
	res := AddGroupMemberResponseV2{}
	_, err := r.client.Put(fmt.Sprintf("/v2/groups/%s/members", groupID), &request, &res)
	if err != nil {
		return AddGroupMemberResponseV2{}, err
	}
	return res, nil
}

func (r *accessGroupMemberRepository) Remove(groupID string, memberID string) error {
	_, err := r.client.Delete(helpers.Tprintf("/v2/groups/{{.GroupID}}/members/{{.MemberID}}", map[string]interface{}{
		"GroupID":  groupID,
		"MemberID": memberID,
	}))
	return err
}
