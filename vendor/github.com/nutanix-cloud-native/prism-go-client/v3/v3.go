package v3

import (
	"fmt"
	"strings"

	"github.com/nutanix-cloud-native/prism-go-client"
	"github.com/nutanix-cloud-native/prism-go-client/internal"
)

const (
	libraryVersion = "v3"
	absolutePath   = "api/nutanix/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
	clientName     = "prism_central"
)

// Client manages the V3 API
type Client struct {
	client *internal.Client
	V3     Service
}

// NewV3Client return a internal to operate V3 resources
func NewV3Client(credentials prismgoclient.Credentials) (*Client, error) {
	var baseClient *internal.Client

	// check if all required fields are present. Else create an empty internal
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		c, err := internal.NewClient(
			internal.WithCredentials(&credentials),
			internal.WithUserAgent(userAgent),
			internal.WithAbsolutePath(absolutePath))
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Prism Central (PC) Client is missing. "+
			"Please provide required details - %s in provider configuration.", strings.Join(credentials.RequiredFields[clientName], ", "))

		baseClient = &internal.Client{UserAgent: userAgent, ErrorMsg: errorMsg}
	}

	f := &Client{
		client: baseClient,
		V3: Operations{
			client: baseClient,
		},
	}

	// f.internal.OnRequestCompleted(func(req *http.Request, resp *http.Response, v interface{}) {
	// 	if v != nil {
	// 		utils.PrintToJSON(v, "[Debug] FINISHED REQUEST")
	// 		// TBD: How to print responses before all requests.
	// 	}
	// })

	return f, nil
}
