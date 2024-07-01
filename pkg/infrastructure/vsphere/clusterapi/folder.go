package clusterapi

import (
	"context"
	"path"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func folderExists(ctx context.Context, dir string, session *session.Session) (*object.Folder, error) {
	/* scenarios:
	 * 1. folder exists and returns
	 * 2. folder does not exist and err and folder are nil
	 * 3. finder.Folder fails and returns folder nil and error
	 */
	var notFoundError *find.NotFoundError
	folder, err := session.Finder.Folder(ctx, dir)

	// scenario two
	if folder == nil && errors.As(err, &notFoundError) {
		return nil, nil
	}
	// scenario three
	if err != nil {
		return nil, err
	}
	// scenario one
	return folder, nil
}

func createFolder(ctx context.Context, fullpath string, session *session.Session) (*object.Folder, error) {
	var folder *object.Folder
	var err error

	dir := path.Dir(fullpath)
	base := path.Base(fullpath)

	// if folder is nil the fullpath does not exist
	if folder, err = folderExists(ctx, dir, session); err == nil && folder == nil {
		folder, err = createFolder(ctx, dir, session)
		if err != nil {
			return nil, err
		}
	}

	if folder != nil && err == nil {
		return folder.CreateFolder(ctx, base)
	}
	return folder, err
}
