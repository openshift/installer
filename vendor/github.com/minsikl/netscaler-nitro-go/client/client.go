package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/minsikl/netscaler-nitro-go/datatypes"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// Nitro Client
type NitroClient struct {
	Protocol  string
	IpAddress string
	Mode      string
	User      string
	Password  string
	Debug     bool
}

func (n *NitroClient) Add(req interface{}, options ...string) error {
	resource, err := getResourceStringByObject(req)
	if err != nil {
		return err
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}
	requestQuery := resource + getOptions(options)
	responseBody, _, err := HTTPRequest(n, requestQuery, "POST", reqJson)
	if err != nil {
		return fmt.Errorf("Error in POST 's'", err.Error())
	}
	if len(responseBody) > 0 {
		res := datatypes.BaseRes{}
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return fmt.Errorf("Error in Unmarshal '%s'", err.Error())
		}
		if *res.Severity == "ERROR" {
			return fmt.Errorf("Error in POST : Errorcode '%d' Message '%s' Severity '%s'\r\n", *res.Errorcode, *res.Message, *res.Severity)
		}
	}
	return nil
}

func (n *NitroClient) Update(req interface{}, options ...string) error {
	resource, err := getResourceStringByObject(req)
	if err != nil {
		return err
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}
	requestQuery := resource + getOptions(options)
	responseBody, _, err := HTTPRequest(n, requestQuery, "PUT", reqJson)
	if err != nil {
		return fmt.Errorf("Error in PUT 's'", err.Error())
	}
	if len(responseBody) > 0 {
		res := datatypes.BaseRes{}
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return fmt.Errorf("Error in Unmarshal '%s'", err.Error())
		}
		if *res.Severity == "ERROR" {
			return fmt.Errorf("Error in POST : Errorcode '%d' Message '%s' Severity '%s'\r\n", *res.Errorcode, *res.Message, *res.Severity)
		}
	}
	return nil
}

func (n *NitroClient) Get(res interface{}, resourceName string, options ...string) error {
	resource, err := getResourceStringByObject(res)
	if err != nil {
		return err
	}

	requestQuery := resource + "/" + resourceName + getOptions(options)
	responseBody, _, err := HTTPRequest(n, requestQuery, "GET", nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, res)
	if err != nil {
		return fmt.Errorf("Error in Unmarshal '%s'", err.Error())
	}
	resMessage := datatypes.BaseRes{}
	err = json.Unmarshal(responseBody, &resMessage)
	if err != nil {
		return fmt.Errorf("Error in Unmarshal '%s'", err.Error())
	}
	if *resMessage.Severity == "ERROR" {
		return fmt.Errorf("Error in POST : Errorcode '%d' Message '%s' Severity '%s'\r\n", *resMessage.Errorcode, *resMessage.Message, *resMessage.Severity)
	}

	return nil
}

func (n *NitroClient) Delete(req interface{}, resourceName string, options ...string) error {
	resource, err := getResourceStringByObject(req)
	if err != nil {
		return err
	}

	requestQuery := resource + "/" + resourceName + getOptions(options)
	responseBody, _, err := HTTPRequest(n, requestQuery, "DELETE", nil)
	if err != nil {
		return err
	}
	resMessage := datatypes.BaseRes{}
	err = json.Unmarshal(responseBody, &resMessage)
	if *resMessage.Severity == "ERROR" {
		return fmt.Errorf("Error in POST : Errorcode '%d' Message '%s' Severity '%s'\r\n", *resMessage.Errorcode, *resMessage.Message, *resMessage.Severity)
	}

	return nil
}

func (n *NitroClient) Enable(req interface{}, enable bool) error {
	resource, err := getResourceStringByObject(req)
	if err != nil {
		return err
	}
	reqJson, err := json.Marshal(req)
	log.Printf(string(reqJson))
	if err != nil {
		return err
	}
	action := "/?action=enable"
	if enable == false {
		action = "/?action=disable"
	}
	query := resource+action
	log.Println("QUERY : " + query)
	responseBody, _, err := HTTPRequest(n, query, "POST", reqJson)
	if err != nil {
		return fmt.Errorf("Error in POST '%s' for Enable", err.Error())
	}
	if len(responseBody) > 0 {
		res := datatypes.BaseRes{}
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			return fmt.Errorf("Error in Unmarshal '%s'", err.Error())
		}
		if *res.Severity == "ERROR" {
			return fmt.Errorf("Error in POST : Errorcode '%d' Message '%s' Severity '%s'\r\n", *res.Errorcode, *res.Message, *res.Severity)
		}
	}
	return nil
}

func NewNitroClient(protocol string, ipAddress string, mode string, user string, password string, debug bool) *NitroClient {
	nClient := NitroClient{
		Protocol:  protocol,
		IpAddress: ipAddress,
		Mode:      mode,
		User:      user,
		Password:  password,
		Debug:     debug,
	}
	return &nClient
}

func HTTPRequest(nClient *NitroClient, requestQuery string, requestType string, requestBody []byte) ([]byte, int, error) {

	// Create a request
	Url := nClient.Protocol + "://" + nClient.IpAddress + "/nitro/v1/" + nClient.Mode + "/" + requestQuery
	requestBodyBuffer := bytes.NewBuffer(requestBody)
	req, err := http.NewRequest(requestType, Url, requestBodyBuffer)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-NITRO-USER", nClient.User)
	req.Header.Set("X-NITRO-PASS", nClient.Password)

	if nClient.Debug {
		log.Println("[DEBUG] Nitro Request Path: ", requestType, req.URL)
		log.Println("[DEBUG] Nitro Request Parameters: ", requestBodyBuffer.String())
	}

	// Execute http request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if nClient.Debug {
		log.Println("[DEBUG] Nitro Response: ", string(responseBody))
	}
	return responseBody, resp.StatusCode, nil
}

func getResourceStringByObject(obj interface{}) (string, error) {
	resourceType := reflect.TypeOf(obj).Elem().Name()
	if len(resourceType) < 4 || (!strings.Contains(resourceType, "Req") && !strings.Contains(resourceType, "Res")) {
		return "", fmt.Errorf("Unable to get resource name from '%s'", resourceType)
	}
	resourceName := resourceType[:len(resourceType)-3]
	resourceBytes := make([]byte, 0)
	for index, character := range []byte(resourceName) {
		if index > 0 && character < 97 {
			resourceBytes = append(resourceBytes, []byte("_"+string(character+32))...)
		} else if character < 97 {
			resourceBytes = append(resourceBytes, character + 32)
		} else {
			resourceBytes = append(resourceBytes, character)
		}
	}
	return string(resourceBytes), nil
}

func getOptions(options []string) string {
	res := ""
	if len(options) > 0 {
		for index, option := range options {
			if index == 0 {
				res = "?" + option
			} else {
				res = res + "&" + option
			}
		}
	}
	return res
}