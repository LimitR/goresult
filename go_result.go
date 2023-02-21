package goresult

import (
	"errors"
	"time"
)

type Result[T any] struct {
	value T
	err   err
	mode  bool
}

func NewResult[T any](value T, args ...bool) *Result[T] {
	if len(args) > 0 && args[0] {
		return &Result[T]{
			value: value,
			err:   err{},
			mode:  true,
		}
	} else {
		return &Result[T]{
			value: value,
			err:   err{},
			mode:  false,
		}
	}
}

func CreateResultFrom[T any](value T, errs error, args ...bool) *Result[T] {
	t := getTrace(2)
	t.message = errs.Error()
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

func (s *Result[T]) AddTrace() {
	s.err.AddTrace()
}

func (s *Result[T]) Unwrap() T {
	s.err.trace = append(s.err.trace, getTrace(2))
	if s.err.Err != nil {
		errStr := s.err.print()
		panic(errStr)
	}
	return s.value
}

func (s *Result[T]) UnwrapDelay(callback func(res T)) T {
	if s.err.Err != nil {
		defer callback(s.value)
		defer func() {
			s.err.Err = nil
		}()
		panic(s.err)
	}
	return s.value
}

func (s *Result[T]) Expect(messageerr string) T {
	if s.err.Err != nil {
		panic(messageerr)
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
	t := getTrace(2)
	t.message = errs.Error()
	s.err = err{
		trace: append(s.err.trace, t),
		Err:   errs,
	}
	return s
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
