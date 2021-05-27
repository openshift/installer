/*
Code to call the IBM IAM Services and get a session object that will be used by the Power Colo Code.


*/

package ibmpisession

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/IBM-Cloud/power-go-client/power/client"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/power-go-client/utils"
	//"github.com/IBM-Cloud/bluemix-go/crn"
)

const (
	offering                 = "power-iaas"
	crnString                = "crn"
	version                  = "v1"
	service                  = "bluemix"
	serviceType              = "public"
	serviceInstanceSeparator = "/"
	separator                = ":"
)

// IBMPISession ...
type IBMPISession struct {
	IAMToken    string
	IMSToken    string
	Power       *client.PowerIaas
	Timeout     time.Duration
	UserAccount string
	Region      string
	Zone        string
}

func powerJSONConsumer() runtime.Consumer {
	return runtime.ConsumerFunc(func(reader io.Reader, data interface{}) error {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(reader)
		if err != nil {
			return err
		}
		b := buf.Bytes()
		if b != nil {
			dec := json.NewDecoder(bytes.NewReader(b))
			dec.UseNumber() // preserve number formats
			err = dec.Decode(data)
		}
		if string(b) == "null" || err != nil {
			errorRecord, _ := data.(*models.Error)
			log.Printf("The errorrecord is %s ", errorRecord.Error)
			return nil
		}
		return err
	})
}

// New ...
/*
The method takes in the following params
iamtoken : this is the token that is passed from the client
region : Obtained from the terraform template. Every template /resource will be required to have this information
timeout:
useraccount:
*/
func New(iamtoken, region string, debug bool, timeout time.Duration, useraccount string, zone string) (*IBMPISession, error) {
	session := &IBMPISession{
		IAMToken:    iamtoken,
		UserAccount: useraccount,
		Region:      region,
		Zone:        zone,
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	apiEndpointURL := utils.GetPowerEndPoint(region)
	transport := httptransport.New(apiEndpointURL, "/", []string{"https"})
	if debug {
		transport.Debug = debug
	}
	transport.Consumers[runtime.JSONMime] = powerJSONConsumer()
	session.Power = client.New(transport, nil)
	session.Timeout = timeout
	return session, nil
}

// NewAuth ...
func NewAuth(sess *IBMPISession, PowerInstanceID string) runtime.ClientAuthInfoWriter {
	var crndata = crnBuilder(PowerInstanceID, sess.UserAccount, sess.Region, sess.Zone)
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		if err := r.SetHeaderParam("Authorization", sess.IAMToken); err != nil {
			return err
		}
		return r.SetHeaderParam("CRN", crndata)
	})

}

// BearerTokenAndCRN ...
func BearerTokenAndCRN(session *IBMPISession, crn string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		if err := r.SetHeaderParam("Authorization", session.IAMToken); err != nil {
			return err
		}
		return r.SetHeaderParam("CRN", crn)
	})
}

// crnBuilder ...
func crnBuilder(powerinstance, useraccount, region string, zone string) string {
	var crnData string
	if zone == "" {
		crnData = crnString + separator + version + separator + service + separator + serviceType + separator + offering + separator + region + separator + "a" + serviceInstanceSeparator + useraccount + separator + powerinstance + separator + separator
	} else {
		crnData = crnString + separator + version + separator + service + separator + serviceType + separator + offering + separator + zone + separator + "a" + serviceInstanceSeparator + useraccount + separator + powerinstance + separator + separator
	}
	return crnData
}
