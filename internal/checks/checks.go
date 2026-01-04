package checks

import (
	"context"
)

type EmitFn func(x any)

type CheckFn func(ctx context.Context, in any, emit EmitFn) error

type Checker interface {
	Name() string
	Check() CheckFn
}

// type checker[T any] struct {
// 	name string
// 	fn   func(ctx context.Context, in T, emit EmitFn) error
// }
//
// func (c checker[T]) Name() string { return c.name }
//
// func (c checker[T]) Check() CheckFn {
// 	return func(ctx context.Context, in any, emit EmitFn) error {
// 		v, ok := in.(T)
// 		if !ok {
// 			return fmt.Errorf("check %s: want %T, got %T", c.name, *new(T), in)
// 		}
// 		return c.fn(ctx, v, emit)
// 	}
// }
//
// func New[T any](name string, fn func(ctx context.Context, in T, emit EmitFn) error) Checker {
// 	return checker[T]{name: name, fn: fn}
// }
