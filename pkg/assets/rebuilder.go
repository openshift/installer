package assets

import (
	"context"
)

// ConstantDataRebuilder returns a Rebuild function which sets the
// data to a constant value.
func ConstantDataRebuilder(ctx context.Context, name string, data []byte, frozen bool) Rebuild {
	return func(ctx context.Context, getByName GetByString) (asset *Asset, err error) {
		return &Asset{
			Name:          name,
			Data:          data,
			Frozen:        frozen,
			RebuildHelper: ConstantDataRebuilder(ctx, name, data, frozen),
		}, nil
	}
}
