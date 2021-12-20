package iampapv1

import (
	"github.com/IBM-Cloud/bluemix-go/client"
)

type IAMService interface {
	GetServiceName(serviceDispName string) (string, error)
	GetServiceDispalyName(serviceName string) (string, error)
}

type iamservice struct {
	client *client.Client
}

func newIAMServiceAPI(c *client.Client) IAMService {
	return &iamservice{
		client: c,
	}
}

//GetServiceName ...
func (r *iamservice) GetServiceName(serviceDispName string) (string, error) {
	serviceMap := make(map[string]string)
	serviceMap["IBM Bluemix Container Service"] = "containers-kubernetes"
	serviceMap["All Identity and Access enabled services"] = "All Identity and Access enabled services"
	//rawURL := "/acms/v1/services"
	//resp, err := r.client.Get(rawURL, &services)
	return serviceMap[serviceDispName], nil
}

//GetServiceDisplayName ...
func (r *iamservice) GetServiceDispalyName(serviceName string) (string, error) {
	serviceMap := make(map[string]string)
	serviceMap["containers-kubernetes"] = "IBM Bluemix Container Service"
	serviceMap["All Identity and Access enabled services"] = "All Identity and Access enabled services"
	//rawURL := "/acms/v1/services"
	//resp, err := r.client.Get(rawURL, &services)
	return serviceMap[serviceName], nil
}
