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

// Package template has tools for finding VM templates.
package template

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

// FindTemplate finds a template based either on a UUID or name.
func FindTemplate(ctx context.Context, session *session.Session, templateID string) (*object.VirtualMachine, error) {
	tpl, err := findTemplateByInstanceUUID(ctx, session, templateID)
	if err != nil {
		return nil, err
	}
	if tpl != nil {
		return tpl, nil
	}
	return findTemplateByName(ctx, session, templateID)
}

func findTemplateByInstanceUUID(ctx context.Context, session *session.Session, templateID string) (*object.VirtualMachine, error) {
	log := ctrl.LoggerFrom(ctx)

	if !isValidUUID(templateID) {
		return nil, nil
	}
	log.V(5).Info("Find template by instanceUUID", "instanceUUID", templateID)
	ref, err := session.FindByInstanceUUID(ctx, templateID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying template by instance UUID")
	}
	if ref != nil {
		return object.NewVirtualMachine(session.Client.Client, ref.Reference()), nil
	}
	return nil, nil
}

func findTemplateByName(ctx context.Context, session *session.Session, templateID string) (*object.VirtualMachine, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(5).Info("Find template by name", "name", templateID)
	tpl, err := session.Finder.VirtualMachine(ctx, templateID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find template by name %q", templateID)
	}
	return tpl, nil
}

func isValidUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}
