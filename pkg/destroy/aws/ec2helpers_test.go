package aws

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	awsconfig "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
)

func TestDeleteEC2NetworkInterfaceUnauthorized(t *testing.T) {
	const eniID = "eni-0123456789abcdef0"

	tests := []struct {
		name              string
		deleteCode        string
		describeResponse  string
		describeStatus    int
		wantErrorCode     string
		wantErrorContains string
		wantDescribeCalls int
	}{
		{
			name:              "confirmed missing interface is success",
			deleteCode:        "UnauthorizedOperation",
			describeResponse:  ec2ErrorResponse("InvalidNetworkInterfaceID.NotFound", "interface does not exist"),
			describeStatus:    http.StatusBadRequest,
			wantDescribeCalls: 1,
		},
		{
			name:              "existing interface remains an error",
			deleteCode:        "UnauthorizedOperation",
			describeResponse:  fmt.Sprintf("<DescribeNetworkInterfacesResponse><networkInterfaceSet><item><networkInterfaceId>%s</networkInterfaceId></item></networkInterfaceSet></DescribeNetworkInterfacesResponse>", eniID),
			describeStatus:    http.StatusOK,
			wantErrorCode:     "UnauthorizedOperation",
			wantErrorContains: "denied",
			wantDescribeCalls: 1,
		},
		{
			name:              "empty describe response remains an error",
			deleteCode:        "UnauthorizedOperation",
			describeResponse:  "<DescribeNetworkInterfacesResponse><networkInterfaceSet></networkInterfaceSet></DescribeNetworkInterfacesResponse>",
			describeStatus:    http.StatusOK,
			wantErrorCode:     "UnauthorizedOperation",
			wantErrorContains: "denied",
			wantDescribeCalls: 1,
		},
		{
			name:              "read denied remains an error",
			deleteCode:        "UnauthorizedOperation",
			describeResponse:  ec2ErrorResponse("UnauthorizedOperation", "describe denied"),
			describeStatus:    http.StatusForbidden,
			wantErrorCode:     "UnauthorizedOperation",
			wantErrorContains: "confirming its existence failed",
			wantDescribeCalls: 1,
		},
		{
			name:              "other delete error does not perform a read",
			deleteCode:        "DependencyViolation",
			describeResponse:  "",
			describeStatus:    0,
			wantErrorCode:     "DependencyViolation",
			wantErrorContains: "denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteCalls := 0
			describeCalls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Fatalf("parsing request form: %v", err)
				}
				switch r.Form.Get("Action") {
				case "DeleteNetworkInterface":
					deleteCalls++
					w.Header().Set("Content-Type", "text/xml")
					w.WriteHeader(http.StatusForbidden)
					_, _ = w.Write([]byte(ec2ErrorResponse(tt.deleteCode, "denied")))
				case "DescribeNetworkInterfaces":
					describeCalls++
					if r.Form.Get("NetworkInterfaceId.1") != eniID {
						t.Errorf("described interface %q, want %q", r.Form.Get("NetworkInterfaceId.1"), eniID)
					}
					w.Header().Set("Content-Type", "text/xml")
					w.WriteHeader(tt.describeStatus)
					_, _ = w.Write([]byte(tt.describeResponse))
				default:
					t.Errorf("unexpected EC2 action %q", r.Form.Get("Action"))
					http.Error(w, "unexpected action", http.StatusBadRequest)
				}
			}))
			defer server.Close()

			client := ec2v2.NewFromConfig(awsconfig.Config{
				Region:      "us-east-1",
				Credentials: credentials.NewStaticCredentialsProvider("test", "test", ""),
				HTTPClient:  server.Client(),
			}, func(options *ec2v2.Options) {
				options.BaseEndpoint = awsconfig.String(server.URL)
			})

			err := deleteEC2NetworkInterface(context.Background(), client, eniID, logrus.New())
			if got := HandleErrorCode(err); got != tt.wantErrorCode {
				t.Errorf("HandleErrorCode(err) = %q, want %q (err %v)", got, tt.wantErrorCode, err)
			}
			if tt.wantErrorContains == "" && err != nil {
				t.Errorf("deleteEC2NetworkInterface() error = %v, want nil", err)
			}
			if tt.wantErrorContains != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErrorContains)) {
				t.Errorf("deleteEC2NetworkInterface() error = %v, want text %q", err, tt.wantErrorContains)
			}
			if deleteCalls != 1 {
				t.Errorf("DeleteNetworkInterface called %d times, want 1", deleteCalls)
			}
			if describeCalls != tt.wantDescribeCalls {
				t.Errorf("DescribeNetworkInterfaces called %d times, want %d", describeCalls, tt.wantDescribeCalls)
			}
		})
	}
}

func ec2ErrorResponse(code, message string) string {
	return fmt.Sprintf("<Response><Errors><Error><Code>%s</Code><Message>%s</Message></Error></Errors><RequestID>test</RequestID></Response>", code, message)
}
