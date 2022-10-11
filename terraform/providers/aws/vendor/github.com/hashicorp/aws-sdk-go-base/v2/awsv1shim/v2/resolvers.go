package awsv1shim

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"

	configv2 "github.com/aws/aws-sdk-go-v2/config"
)

func resolveCustomCABundle(ctx context.Context, configSources []interface{}) (value io.Reader, found bool, err error) {
	for _, source := range configSources {
		switch cfg := source.(type) {
		case configv2.LoadOptions:
			value, found, err = loadOptionsGetCustomCABundle(ctx, cfg)
		case configv2.EnvConfig:
			value, found, err = envConfigGetCustomCABundle(ctx, cfg)
		case configv2.SharedConfig:
			value, found, err = sharedConfigGetCustomCABundle(ctx, cfg)
		default:
			log.Printf("[WARN] Unrecognized config source: %T", source)
			continue
		}
		if err != nil || found {
			break
		}
	}

	return
}

// Copied from https://github.com/aws/aws-sdk-go-v2/blob/889e1da2776ae5bd6d056cf44f6ce6d043237769/config/load_options.go#L334-L340
func loadOptionsGetCustomCABundle(_ context.Context, o configv2.LoadOptions) (io.Reader, bool, error) { //nolint:unparam
	if o.CustomCABundle == nil {
		return nil, false, nil
	}

	return o.CustomCABundle, true, nil
}

// Copied from https://github.com/aws/aws-sdk-go-v2/blob/889e1da2776ae5bd6d056cf44f6ce6d043237769/config/env_config.go#L463-L473
func envConfigGetCustomCABundle(_ context.Context, c configv2.EnvConfig) (io.Reader, bool, error) {
	if len(c.CustomCABundle) == 0 {
		return nil, false, nil
	}

	b, err := ioutil.ReadFile(c.CustomCABundle)
	if err != nil {
		return nil, false, err
	}
	return bytes.NewReader(b), true, nil
}

// Copied from https://github.com/aws/aws-sdk-go-v2/blob/889e1da2776ae5bd6d056cf44f6ce6d043237769/config/shared_config.go#L350-L360
func sharedConfigGetCustomCABundle(_ context.Context, c configv2.SharedConfig) (io.Reader, bool, error) {
	if len(c.CustomCABundle) == 0 {
		return nil, false, nil
	}

	b, err := ioutil.ReadFile(c.CustomCABundle)
	if err != nil {
		return nil, false, err
	}
	return bytes.NewReader(b), true, nil
}
