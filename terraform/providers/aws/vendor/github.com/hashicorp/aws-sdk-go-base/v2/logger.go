package awsbase

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/smithy-go/logging"
)

type debugLogger struct{}

func (l debugLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	s = strings.ReplaceAll(s, "\r", "") // Works around https://github.com/jen20/teamcity-go-test/pull/2
	log.Printf("[%s] [aws-sdk-go-v2] %s", classification, s)
}
