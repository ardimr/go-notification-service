package cloudstorage

import (
	"context"
)

type CloudStorageInterface interface {
	ListBuckets(ctx context.Context)
}
