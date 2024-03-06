package validations

import (
	"fmt"
	"regexp"
)

var KmsArnRE = regexp.MustCompile(
	`^arn:aws[\w-]*:kms:[\w-]+:\d{12}:key\/(mrk-[0-9a-f]{32}$|[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$)`,
)

func ValidateKMSKeyARN(kmsKeyARN *string) error {
	if kmsKeyARN == nil || *kmsKeyARN == "" {
		return nil
	}

	if !KmsArnRE.MatchString(*kmsKeyARN) {
		return fmt.Errorf("expected the kms-key-arn: %s to match %s", *kmsKeyARN, KmsArnRE)
	}
	return nil
}
