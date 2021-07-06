package iamuumv1

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

type AddGroupMemberRequest struct {
	Members []models.AccessGroupMember `json:"members"`
}

type AddGroupMemberResponse struct {
	Members []AddedGroupMember `json:"members"`
}

type AddedGroupMember struct {
	ID string `json:"id"`

	Msg string `json:"msg,omitempty"`

	Type string `json:"type"`

	OK bool `json:"ok,omitempty"`
}

type GroupMembers struct {
	PaginationFields
	Members []models.AccessGroupMember `json:"members"`
}

func (gm *GroupMembers) Resources() []interface{} {
	r := make([]interface{}, len(gm.Members))
	for i := range gm.Members {
		r[i] = gm.Members[i]
	}
	return r
}

type AccessGroupMemberRepository interface {
	List(groupID string) ([]models.AccessGroupMember, error)
	Add(groupID string, request AddGroupMemberRequest) (AddGroupMemberResponse, error)
	Remove(groupID string, memberID string) error
}

type accessGroupMemberRepository struct {
	client *client.Client
}

func NewAccessGroupMemberRepository(c *client.Client) AccessGroupMemberRepository {
	return &accessGroupMemberRepository{
		client: c,
	}

}

func (r *accessGroupMemberRepository) List(groupID string) ([]models.AccessGroupMember, error) {
	members := []models.AccessGroupMember{}
	_, err := r.client.GetPaginated(fmt.Sprintf("/v1/groups/%s/members", groupID),
		NewPaginatedResourcesHandler(&GroupMembers{}), func(resource interface{}) bool {
			if member, ok := resource.(models.AccessGroupMember); ok {
				members = append(members, member)
				return true
			}
			return false
		})
	if err != nil {
		return []models.AccessGroupMember{}, err
	}
	return members, nil
}

func (r *accessGroupMemberRepository) Add(groupID string, request AddGroupMemberRequest) (AddGroupMemberResponse, error) {
	res := AddGroupMemberResponse{}
	_, err := r.client.Put(fmt.Sprintf("/v1/groups/%s/members", groupID), &request, &res)
	if err != nil {
		return AddGroupMemberResponse{}, err
	}
	return res, nil
}

func (r *accessGroupMemberRepository) Remove(groupID string, memberID string) error {
	_, err := r.client.Delete(helpers.Tprintf("/v1/groups/{{.GroupID}}/members/{{.MemberID}}", map[string]interface{}{
		"GroupID":  groupID,
		"MemberID": memberID,
	}))
	return err
}
