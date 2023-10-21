// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"

	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

func testStepTaint(ctx context.Context, step TestStep, wd *plugintest.WorkingDir) error {
	if len(step.Taint) == 0 {
		return nil
	}

	logging.HelperResourceTrace(ctx, fmt.Sprintf("Using TestStep Taint: %v", step.Taint))

	for _, p := range step.Taint {
		err := wd.Taint(ctx, p)
		if err != nil {
			return fmt.Errorf("error tainting resource: %s", err)
		}
	}
	return nil
}
