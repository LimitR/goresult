package goresult

import (
	"context"
	"errors"
	"time"
)

type Result[T any] struct {
	value T
	err   err
	mode  bool
	ctx   context.Context
}

func NewResult[T any](value T, args ...bool) *Result[T] {
	if len(args) > 0 || !args[0] {
		return &Result[T]{
			value: value,
			err:   err{},
			mode:  false,
		}
	} else {
		return &Result[T]{
			value: value,
			err:   err{},
			mode:  true,
		}
	}
}

func CreateResultFrom[T any](value T, errs error, args ...bool) *Result[T] {
	if len(args) > 0 && !args[0] {
		t := trace{}
		t.Message = errs.Error()
		return &Result[T]{
			value: value,
			err: err{
				trace: []trace{
					t,
				},
				TimeStamp: time.Now().Unix(),
				Err:       errs,
			},
		}
	} else {
		t := getTrace(2)
		t.Message = errs.Error()
		return &Result[T]{
			value: value,
			err: err{
				trace: []trace{
					t,
				},
				TimeStamp: time.Now().Unix(),
				Err:       errs,
			},
			mode: true,
		}
	}
}

func CreateResultChannel[T any](ctx context.Context, callback func() (T, error), args ...bool) chan *Result[T] {
	res := make(chan *Result[T])
	t := trace{}
	mode := false
	if len(args) <= 0 {
		t = getTrace(2)
		mode = true
	} else {
		mode = args[0]
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				if len(args) > 0 && !args[0] {
					res <- &Result[T]{
						err: err{
							trace: []trace{
								t,
							},
							Err:       ctx.Err(),
							TimeStamp: time.Now().Unix(),
						},
					}
				} else {
					t.Message = ctx.Err().Error()
					res <- &Result[T]{
						err: err{
							trace: []trace{
								t,
							},
							Err:       ctx.Err(),
							TimeStamp: time.Now().Unix(),
						},
						mode: true,
					}
				}
			default:
				go func() {
					callbackResult, errs := callback()
					if errs != nil {
						t.Message = errs.Error()
					}
					result := &Result[T]{
						value: callbackResult,
						err: err{
							trace: []trace{
								t,
							},
							TimeStamp: time.Now().Unix(),
							Err:       errs,
						},
						mode: mode,
					}
					result.addContext(ctx)
					res <- result
				}()
			}
		}
	}()
	return res
}

func CheckAll[T any](arrayResults []Result[T]) []T {
	result := []T{}
	for i := 0; i < len(arrayResults); i++ {
		if arrayResults[i].IsOk() {
			result = append(result, arrayResults[i].value)
		}
	}
	return result
}

func (s *Result[T]) addContext(ctx context.Context) {
	if s.ctx == nil {
		s.ctx = ctx
	}
}

func (s *Result[T]) AddTrace() {
	if s.mode {
		s.err.AddTrace()
	}
}

func (s *Result[T]) Unwrap() T {
	if s.mode {
		s.err.trace = append(s.err.trace, getTrace(2))
		if s.err.Err != nil {
			errStr := s.err.print()
			panic(errStr)
		}
		return s.value
	} else {
		if s.err.Err != nil {
			panic(s.err.Err.Error())
		}
		return s.value
	}
}

func (s *Result[T]) UnwrapDelay(callback func(res T)) T {
	if s.err.Err != nil {
		defer callback(s.value)
		defer func() {
			s.err.Err = nil
		}()
		errStr := s.err.print()
		panic(errStr)
	}
	return s.value
}

func (s *Result[T]) Expect(messageError string) T {
	if s.err.Err != nil {
		panic(messageError)
	}
	return s.value
}

func (s *Result[T]) UnwrapOrElse(value T) T {
	if s.err.Err != nil {
		s.err.Err = nil
		s.value = value
	}
	return s.value
}

func (s *Result[T]) UnwrapOrOn(callback func(error) T) T {
	if s.err.Err != nil {
		s.err.Err = nil
		res := callback(s.err.Err)
		s.value = res
		return res
	}
	return s.value
}

func (s *Result[T]) AddError(errs error) *Result[T] {
	if s.mode {
		t := getTrace(2)
		t.Message = errs.Error()
		s.err = err{
			trace: append(s.err.trace, t),
			Err:   errs,
		}
		return s
	} else {
		t := trace{}
		t.Message = errs.Error()
		s.err = err{
			trace: append(s.err.trace, t),
			Err:   errs,
		}
		return s
	}
}

func (s *Result[T]) GetErrorTrace() []trace {
	return s.err.trace
}

func (s *Result[T]) Match(errs error) error {
	if errors.Is(s.err.Err, errs) {
		return s.err.Err
	} else {
		return nil
	}
}

func (s *Result[T]) IsOk() bool {
	return s.err.Err == nil
}
