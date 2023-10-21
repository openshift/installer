// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"context"
)

type loggerKeyT string

const loggerKey loggerKeyT = "logger-key"

func RegisterLogger(ctx context.Context, logger TfLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func RetrieveLogger(ctx context.Context) TfLogger {
	logger, ok := ctx.Value(loggerKey).(TfLogger)
	if !ok {
		return TfLogger("")
	}
	return logger
}
