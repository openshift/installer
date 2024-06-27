package vsphere

import (
	"context"
	"errors"
	"time"

	"github.com/vmware/govmomi/find"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(validationCtx *validationContext, folderPath string, fldPath *field.Path) field.ErrorList {
	var notFoundError *find.NotFoundError
	allErrs := field.ErrorList{}
	finder := validationCtx.Finder
	// If no folder is specified, skip this check as the folder will be created.
	if folderPath == "" {
		return allErrs
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	folder, err := finder.Folder(ctx, folderPath)
	if err != nil && !errors.As(err, &notFoundError) {
		return append(allErrs, field.Invalid(fldPath, folderPath, err.Error()))
	}

	// folder was not found so no privilege check can be performed
	if folder == nil {
		return allErrs
	}
	permissionGroup := permissions[permissionFolder]

	err = comparePrivileges(ctx, validationCtx, folder.Reference(), permissionGroup)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}
	return allErrs
}
