package osclients

import (
	"context"
	"iter"

	"github.com/gophercloud/gophercloud/v2/pagination"
)

type result[T any] struct {
	ok  T
	err error
}

func (r result[T]) Ok() T      { return r.ok }
func (r result[T]) Err() error { return r.err }

// Result carries either a result or a non-nil error.
type Result[T any] interface {
	Ok() T
	Err() error
}

func NewResultOk[T any](ok T) result[T] {
	return result[T]{ok: ok}
}

func NewResultErr[T any](err error) result[T] {
	return result[T]{err: err}
}

type ResourceFilter[osResourceT any] func(*osResourceT) bool

func Filter[osResourceT any](in iter.Seq2[*osResourceT, error], filters ...ResourceFilter[osResourceT]) iter.Seq2[*osResourceT, error] {
	return func(yield func(*osResourceT, error) bool) {
	next:
		for osResource, err := range in {
			if err != nil {
				yield(nil, err)
				return
			}
			for _, filter := range filters {
				if !filter(osResource) {
					continue next
				}
			}
			if !yield(osResource, nil) {
				return
			}
		}
	}

}

func JustOne[T any, R Result[*T]](in <-chan R, duplicateError error) (*T, error) {
	var found *T
	for result := range in {
		if err := result.Err(); err != nil {
			return nil, err
		}
		if found != nil {
			return nil, duplicateError
		}
		ok := result.Ok()
		found = ok
	}
	return found, nil
}

func yieldPage[osResourcePT *osResourceT, osResourceT any](extracter func(pagination.Page) ([]osResourceT, error), yield func(osResourcePT, error) bool) func(context.Context, pagination.Page) (bool, error) {
	return func(ctx context.Context, page pagination.Page) (bool, error) {
		pageItems, err := extracter(page)
		if err != nil {
			_ = yield(nil, err)
			return false, err
		}
		for i := range pageItems {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				_ = yield(nil, err)
				return false, err
			default:
				if !yield(&pageItems[i], nil) {
					return false, nil
				}
			}
		}
		return true, nil
	}
}
