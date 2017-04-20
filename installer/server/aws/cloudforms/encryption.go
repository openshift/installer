package cloudforms

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// SecretAssets are secret assets as raw bytes.
type SecretAssets struct {
	CACert     []byte
	ClientCert []byte
	ClientKey  []byte
}

// CompactSecretAssets as secrets assets which have been gzipped and base64
// encoded.
type compactSecretAssets struct {
	CACert     string
	ClientCert string
	ClientKey  string
}

func compressData(d []byte) (string, error) {
	var buff bytes.Buffer
	gzw := gzip.NewWriter(&buff)
	if _, err := gzw.Write(d); err != nil {
		return "", err
	}
	if err := gzw.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buff.Bytes()), nil
}

type encryptService interface {
	Encrypt(*kms.EncryptInput) (*kms.EncryptOutput, error)
}

func (r *SecretAssets) compact(cfg *Config, kmsSvc encryptService) (*compactSecretAssets, error) {
	var err error
	compact := func(data []byte) string {
		if err != nil {
			return ""
		}

		encryptInput := kms.EncryptInput{
			KeyId:     aws.String(cfg.KMSKeyARN),
			Plaintext: data,
		}

		var encryptOutput *kms.EncryptOutput
		if encryptOutput, err = kmsSvc.Encrypt(&encryptInput); err != nil {
			err = fmt.Errorf("KMS encryption error: %v", maybeAwsErr(err))
			return ""
		}
		data = encryptOutput.CiphertextBlob

		var out string
		if out, err = compressData(data); err != nil {
			return ""
		}
		return out
	}
	compactAssets := compactSecretAssets{
		CACert:     compact(r.CACert),
		ClientCert: compact(r.ClientCert),
		ClientKey:  compact(r.ClientKey),
	}
	if err != nil {
		return nil, err
	}
	return &compactAssets, nil
}
