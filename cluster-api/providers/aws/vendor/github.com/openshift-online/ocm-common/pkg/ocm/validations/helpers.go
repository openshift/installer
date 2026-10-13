package validations

import (
        "fmt"
        "net/url"
        "strings"
)

func ValidateIssuerUrlMatchesAssumePolicyDocument(
        roleArn string, parsedUrl *url.URL, assumePolicyDocument string) error {
        issuerUrl := parsedUrl.Host
        if parsedUrl.Path != "" {
                issuerUrl += parsedUrl.Path
        }
        decodedAssumePolicyDocument, err := url.QueryUnescape(assumePolicyDocument)
        if err != nil {
                return err
        }
        if !strings.Contains(decodedAssumePolicyDocument, issuerUrl) {
                return fmt.Errorf("Operator role '%s' does not have trusted relationship to '%s' issuer URL",
            roleArn, issuerUrl)
        }
        return nil
}
