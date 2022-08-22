package goresult

import (
	"errors"
)

type Result[T any] struct {
	value T
	err   error
}

func NewResult[T any](value T) *Result[T] {
	return &Result[T]{
		value: value,
		err:   nil,
	}
}

func CreateResultFrom[T any](value T, err error) *Result[T] {
	return &Result[T]{
		value: value,
		err:   err,
	}
}

func CheckAll[T any](arrayResults []Result[T]) []T {
	result := []T{}
	for i := 0; i < len(arrayResults); i++ {
		if arrayResults[i].IsOk() {
			result = append(result, arrayResults[i].Unwrap())
		}
	}
	return result
}

func (s *Result[T]) Unwrap() T {
	if s.err != nil {
		panic(s.err)
	}
	return s.value
}

func (s *Result[T]) UnwrapDelay(callback func(res T)) T {
	if s.err != nil {
		defer callback(s.value)
		defer func() {
			s.err = nil
		}()
		panic(s.err)
	}
	return s.value
}

func (s *Result[T]) Expect(messageError string) T {
	if s.err != nil {
		panic(messageError)
	}
	return s.value
}

func (s *Result[T]) UnwrapOrElse(value T) T {
	if s.err != nil {
		s.err = nil
		s.value = value
	}
	return s.value
}

func (s *Result[T]) UnwrapOrOn(callback func(error) T) T {
	if s.err != nil {
		s.err = nil
		res := callback(s.err)
		s.value = res
		return res
	}
	return s.value
}

func (s *Result[T]) AddError(value string) *Result[T] {
	s.err = errors.New(value)
	return s
}

func (s *Result[T]) Match(err error) error {
	if errors.Is(s.err, err) {
		return s.err
	} else {
		return nil
	}
}

func (s *Result[T]) IsOk() bool {
	return s.err == nil
}
