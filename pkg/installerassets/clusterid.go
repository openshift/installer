package installerassets

import (
	"context"

	"github.com/pborman/uuid"
)

func getUUID(ctx context.Context) (data []byte, err error) {
	return []byte(uuid.New()), nil
}

func init() {
	Defaults["cluster-id"] = getUUID
}
