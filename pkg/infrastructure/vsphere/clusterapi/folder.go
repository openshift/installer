package clusterapi

import (
	"context"
	"path"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func createFolder(ctx context.Context, fullpath string, session *session.Session) (*object.Folder, error) {
	dir := path.Dir(fullpath)
	base := path.Base(fullpath)
	finder := session.Finder

	folder, err := finder.Folder(ctx, fullpath)

	if folder == nil {
		folder, err = finder.Folder(ctx, dir)

		var notFoundError *find.NotFoundError
		if errors.As(err, &notFoundError) {
			folder, err = createFolder(ctx, dir, session)
			if err != nil {
				return folder, err
			}
		}

		if folder != nil {
			folder, err = folder.CreateFolder(ctx, base)
			if err != nil {
				return folder, err
			}
		}
	}
	return folder, err
}
