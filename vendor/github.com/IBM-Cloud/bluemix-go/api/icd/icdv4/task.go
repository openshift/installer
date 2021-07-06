package icdv4

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

// type TaskResult struct {
//       Task Task `json:"task"`
// }

// type Task struct {
//       Id string `json:"id"`
//       Description string `json:"description"`
//       Status string `json:"status"`
//       DeploymentId string `json:"deployment_id"`
//       ProgressPercent int `json:"progress_percent"`
//       CreatedAt string `json:"created_at"`

// }

type Tasks interface {
	GetTask(taskId string) (Task, error)
}

type tasks struct {
	client *client.Client
}

func newTaskAPI(c *client.Client) Tasks {
	return &tasks{
		client: c,
	}
}

func (r *tasks) GetTask(taskId string) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/tasks/%s", utils.EscapeUrlParm(taskId))
	_, err := r.client.Get(rawURL, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}
