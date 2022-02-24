package icdv4

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

type UserReq struct {
	User User `json:"user"`
}

type User struct {
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type TaskResult struct {
	Task Task `json:"task"`
}

type Task struct {
	Id              string `json:"id"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	DeploymentId    string `json:"deployment_id"`
	ProgressPercent int    `json:"progress_percent"`
	CreatedAt       string `json:"created_at"`
}

type Users interface {
	CreateUser(icdId string, userReq UserReq) (Task, error)
	UpdateUser(icdId string, userName string, userReq UserReq) (Task, error)
	DeleteUser(icdId string, userName string) (Task, error)
}

type users struct {
	client *client.Client
}

func newUsersAPI(c *client.Client) Users {
	return &users{
		client: c,
	}
}

func (r *users) CreateUser(icdId string, userReq UserReq) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/users", utils.EscapeUrlParm(icdId))
	_, err := r.client.Post(rawURL, &userReq, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}

func (r *users) UpdateUser(icdId string, userName string, userReq UserReq) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/users/%s", utils.EscapeUrlParm(icdId), userName)
	_, err := r.client.Patch(rawURL, &userReq, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}

func (r *users) DeleteUser(icdId string, userName string) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/users/%s", utils.EscapeUrlParm(icdId), userName)
	_, err := r.client.DeleteWithResp(rawURL, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}
