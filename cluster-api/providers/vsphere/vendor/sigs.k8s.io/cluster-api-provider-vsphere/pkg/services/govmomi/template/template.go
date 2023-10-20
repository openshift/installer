/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package template

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

type tplContext interface {
	context.Context
	GetLogger() logr.Logger
	GetSession() *session.Session
}

// FindTemplate finds a template based either on a UUID or name.
func FindTemplate(ctx tplContext, templateID string) (*object.VirtualMachine, error) {
	tpl, err := findTemplateByInstanceUUID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	if tpl != nil {
		return tpl, nil
	}
	return findTemplateByName(ctx, templateID)
}

func findTemplateByInstanceUUID(ctx tplContext, templateID string) (*object.VirtualMachine, error) {
	if !isValidUUID(templateID) {
		return nil, nil
	}
	ctx.GetLogger().V(6).Info("find template by instance uuid", "instance-uuid", templateID)
	ref, err := ctx.GetSession().FindByInstanceUUID(ctx, templateID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying template by instance UUID")
	}
	if ref != nil {
		return object.NewVirtualMachine(ctx.GetSession().Client.Client, ref.Reference()), nil
	}
	return nil, nil
}

func findTemplateByName(ctx tplContext, templateID string) (*object.VirtualMachine, error) {
	ctx.GetLogger().V(6).Info("find template by name", "name", templateID)
	tpl, err := ctx.GetSession().Finder.VirtualMachine(ctx, templateID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find template by name %q", templateID)
	}
	return tpl, nil
}

func isValidUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}
