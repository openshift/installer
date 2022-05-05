package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	log "github.com/ovirt/go-ovirt-client-log/v3"
)

func newTerraformLogger() log.Logger {
	return &terraformLogger{}
}

type terraformLogger struct {
	ctx context.Context
}

func (t terraformLogger) Debugf(format string, args ...interface{}) {
	if t.ctx == nil {
		panic(fmt.Errorf("bug: the Terraform logger was used without calling WithContext"))
	}
	tflog.Debug(t.ctx, fmt.Sprintf(format, args...))
}

func (t terraformLogger) Infof(format string, args ...interface{}) {
	if t.ctx == nil {
		panic(fmt.Errorf("bug: the Terraform logger was used without calling WithContext"))
	}
	tflog.Info(t.ctx, fmt.Sprintf(format, args...))
}

func (t terraformLogger) Warningf(format string, args ...interface{}) {
	if t.ctx == nil {
		panic(fmt.Errorf("bug: the Terraform logger was used without calling WithContext"))
	}
	tflog.Warn(t.ctx, fmt.Sprintf(format, args...))
}

func (t terraformLogger) Errorf(format string, args ...interface{}) {
	if t.ctx == nil {
		panic(fmt.Errorf("bug: the Terraform logger was used without calling WithContext"))
	}
	tflog.Error(t.ctx, fmt.Sprintf(format, args...))
}

func (t terraformLogger) WithContext(ctx context.Context) log.Logger {
	return &terraformLogger{
		ctx: ctx,
	}
}
