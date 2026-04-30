package rx_helper

import (
	"context"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/constants/error_code"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/reactivex/rxgo/v2"
)

var (
	// Deprecated: This function is deprecated and should not be used.
	RxPanicError = error_helper.New(errors.New("internal rx error"), error_code.UnknownError)
)

// Deprecated: This function is deprecated and should not be used.
type mapFunc func(ctx context.Context, i interface{}) (r interface{}, e error)

// Deprecated: This function is deprecated and should not be used.
type mapRecoverFunc func(ctx context.Context, i interface{}, r *interface{}, e *error)

// MapDeferRecoverFunc use this function to handle panic inside the mapFunc.
//
// mapFunc is the function to be passed into rxgo.Observable Map parameter.
//
// [OPTIONAL] mapRecoverFunc is the function to be run when panic is occurred on the mapFunc.
//
// # Default behaviour of mapRecoverFunc will set RxPanicError on mapFunc return e error
//
// Deprecated: This function is deprecated and should not be used.
func MapDeferRecoverFunc(
	mapFunc mapFunc,
	mapRecoverFunc ...mapRecoverFunc,
) func(ctx context.Context, i interface{}) (r interface{}, e error) {
	return func(ctx context.Context, i interface{}) (r interface{}, e error) {
		defer func() {
			if rc := recover(); rc != nil {
				if mapRecoverFunc != nil {
					mapRecoverFunc[0](ctx, i, &r, &e)
				} else {
					fmt.Println("[rx_helper.MapDeferRecoverFunc]:", rc)
					fmt.Println("Stack trace:", string(debug.Stack()))
					e = RxPanicError
				}
			}
		}()

		return mapFunc(ctx, i)
	}
}

// Deprecated: This function is deprecated and should not be used.
type itemToObservable func(item rxgo.Item) (o rxgo.Observable)

// Deprecated: This function is deprecated and should not be used.
type flatMapRecoverFunc func(item rxgo.Item, o *rxgo.Observable)

// FlatMapDeferRecoverFunc use this function to handle panic inside the flatMapFunc.
//
// flatMapFunc is the function to be passed into rxgo.Observable FlatMap parameter.
//
// [OPTIONAL] flatMapRecoverFunc is the function to be run when panic is occurred on the flatMapFunc.
//
// # Default behaviour of flatMapRecoverFunc will set rxgo.Thrown(RxPanicError) on flatMapFunc return o rxgo.Observable
//
// Deprecated: This function is deprecated and should not be used.
func FlatMapDeferRecoverFunc(
	flatMapFunc itemToObservable,
	flatMapRecoverFunc ...flatMapRecoverFunc,
) rxgo.ItemToObservable {
	return func(item rxgo.Item) (o rxgo.Observable) {
		defer func() {
			if r := recover(); r != nil {
				if flatMapRecoverFunc != nil {
					flatMapRecoverFunc[0](item, &o)
				} else {
					fmt.Println("[rx_helper.FlatMapDeferRecoverFunc]:", r)
					fmt.Println("Stack trace:", string(debug.Stack()))
					o = rxgo.Thrown(RxPanicError)
				}
			}
		}()

		return flatMapFunc(item)
	}
}

// Deprecated: This function is deprecated and should not be used.
type funcN func(i ...interface{}) (r interface{})

// Deprecated: This function is deprecated and should not be used.
type combineLatestRecoverFunc func(i []interface{}, r *interface{})

// CombineLatestDeferRecoverFunc use this function to handle panic inside the funcN.
//
// funcN is the function to be passed into rxgo.Observable CombineLatest parameter.
//
// [OPTIONAL] combineLatestRecoverFunc is the function to be run when panic is occurred on the funcN.
//
// Default behaviour of combineLatestRecoverFunc will set rxgo.Error(RxPanicError) on funcN return r interface{}
//
// # How to handle error from this function is to use RxError func when subscribe to the rxgo.Observable
//
// Deprecated: This function is deprecated and should not be used.
func CombineLatestDeferRecoverFunc(
	funcN funcN,
	combineLatestRecoverFunc ...combineLatestRecoverFunc,
) rxgo.FuncN {
	return func(i ...interface{}) (r interface{}) {
		defer func() {
			if rc := recover(); rc != nil {
				if combineLatestRecoverFunc != nil {
					combineLatestRecoverFunc[0](i, &r)
				} else {
					fmt.Println("[rx_helper.CombineLatestDeferRecoverFunc]:", rc)
					fmt.Println("Stack trace:", string(debug.Stack()))
					r = rxgo.Error(RxPanicError)
				}
			}
		}()

		return funcN(i...)
	}
}

// Deprecated: This function is deprecated and should not be used.
type producerFunc func(ctx context.Context, next chan<- rxgo.Item)

// Deprecated: This function is deprecated and should not be used.
type producerRecoverFunc func(ctx context.Context, next chan<- rxgo.Item)

// ProducerDeferRecoverFunc use this function to handle panic inside the producerFunc.
//
// producerFunc is the function to be passed into rxgo.Producer struct.
//
// [OPTIONAL] producerRecoverFunc is the function to be run when panic is occurred on the producerFunc.
//
// # Default behaviour of producerRecoverFunc will set RxPanicError on producerFunc next channel
//
// Deprecated: This function is deprecated and should not be used.
func ProducerDeferRecoverFunc(producerFunc producerFunc, producerRecoverFunc ...producerRecoverFunc) func(ctx context.Context, next chan<- rxgo.Item) {
	return func(ctx context.Context, next chan<- rxgo.Item) {
		defer func() {
			if rc := recover(); rc != nil {
				if producerRecoverFunc != nil {
					producerRecoverFunc[0](ctx, next)
				} else {
					fmt.Println("[rx_helper.ProducerDeferRecoverFunc]:", rc)
					fmt.Println("Stack trace:", string(debug.Stack()))
					next <- rxgo.Error(RxPanicError)
				}
			}
		}()

		producerFunc(ctx, next)
	}
}

// RxError this function will get any error that exist on rxgo.Item or nested rxgo.Item->rxgo.Item
//
// Deprecated: This function is deprecated and should not be used.
func RxError(item rxgo.Item) error {
	// Handle panic recover case when using CombineLatestDeferRecover function.
	// the default behaviour is return rxgo.Error(error) inside this item.V
	if v, ok := item.V.(rxgo.Item); ok && v.Error() {
		item.E = v.E
	}

	return item.E
}
