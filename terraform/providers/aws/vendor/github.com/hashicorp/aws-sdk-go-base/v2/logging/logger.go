// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func New(ctx context.Context, name string) (context.Context, TfLogger) {
	ctx = tflog.NewSubsystem(ctx, name, tflog.WithRootFields())
	logger := TfLogger(name)

	return ctx, logger
}

type TfLogger string

func (l TfLogger) Warn(ctx context.Context, msg string, fields ...map[string]any) {
	if l == "" {
		tflog.Warn(ctx, msg, fields...)
	} else {
		tflog.SubsystemWarn(ctx, string(l), msg, fields...)
	}
}

func (l TfLogger) Info(ctx context.Context, msg string, fields ...map[string]any) {
	if l == "" {
		tflog.Info(ctx, msg, fields...)
	} else {
		tflog.SubsystemInfo(ctx, string(l), msg, fields...)
	}
}

func (l TfLogger) Debug(ctx context.Context, msg string, fields ...map[string]any) {
	if l == "" {
		tflog.Debug(ctx, msg, fields...)
	} else {
		tflog.SubsystemDebug(ctx, string(l), msg, fields...)
	}
}

func (l TfLogger) Trace(ctx context.Context, msg string, fields ...map[string]any) {
	if l == "" {
		tflog.Trace(ctx, msg, fields...)
	} else {
		tflog.SubsystemTrace(ctx, string(l), msg, fields...)
	}
}

func (l TfLogger) SetField(ctx context.Context, key string, value any) context.Context {
	if l == "" {
		return tflog.SetField(ctx, key, value)
	} else {
		return tflog.SubsystemSetField(ctx, string(l), key, value)
	}
}
