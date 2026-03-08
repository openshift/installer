package logforwarding

import (
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func LogForwarderObjectAsString(logForwarder *cmv1.LogForwarder) string {
	out := "\n"

	if logForwarder.S3() != nil { // S3 Log Forwarder
		out += fmt.Sprintf("S3 Bucket Prefix:                     %s\n", logForwarder.S3().BucketPrefix())
		out += fmt.Sprintf("S3 Bucket Name:                      %s\n", logForwarder.S3().BucketName())
	} else if logForwarder.Cloudwatch() != nil { // Cloudwatch Log Forwarder
		out += fmt.Sprintf("Cloudwatch Log Group Name:           %s\n",
			logForwarder.Cloudwatch().LogGroupName())
		out += fmt.Sprintf("Cloudwatch Log Distribution Role Arn: %s\n",
			logForwarder.Cloudwatch().LogDistributionRoleArn())
	}

	if logForwarder.Applications() != nil && len(logForwarder.Applications()) > 0 {
		applicationsStr := ""
		for i, app := range logForwarder.Applications() {
			if i > 0 {
				applicationsStr += " "
			}
			applicationsStr += app
		}
		out += fmt.Sprintf("Applications:                        %s\n", applicationsStr)
	}

	if logForwarder.Groups() != nil && len(logForwarder.Groups()) > 0 {
		groupsStr := ""
		for i, group := range logForwarder.Groups() {
			if i > 0 {
				groupsStr += " "
			}
			groupsStr += fmt.Sprintf("(%s,v%s)", group.ID(), group.Version())
		}
		out += fmt.Sprintf("Groups:                              %s\n", groupsStr)
	}

	if logForwarder.Status() != nil {
		if logForwarder.Status().Message() != "" {
			out += fmt.Sprintf("Status Message:                      %s\n", logForwarder.Status().Message())
		}

		if logForwarder.Status().ResolvedApplications() != nil && len(logForwarder.Status().ResolvedApplications()) > 0 {
			resolvedAppsStr := ""
			for i, app := range logForwarder.Status().ResolvedApplications() {
				if i > 0 {
					resolvedAppsStr += " "
				}
				resolvedAppsStr += app
			}
			out += fmt.Sprintf("Resolved Applications:               %s\n", resolvedAppsStr)
		}
	}

	return out
}
