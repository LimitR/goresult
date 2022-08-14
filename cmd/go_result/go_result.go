package goresult

import (
	"errors"
	"log"
)

type Result[T any] struct {
	Value T
	Err   error
}

func NewResult(value any) *Result[any] {
	return &Result[any]{
		Value: value,
		Err:   nil,
	}
}

func (s *Result[T]) Unwrap() T {
	if s.Err != nil {
		log.Fatal(s.Err)
	}
	return s.Value
}

func (s *Result[T]) Some(value T) *Result[T] {
	s.Value = value
	s.Err = nil
	return s
}

func (s *Result[T]) Error(value string) *Result[T] {
	s.Value = s.Value
	s.Err = errors.New(value)
	return s
}

func (s *Result[T]) IsOk() bool {
	if s.Err == nil {
		return true
	} else {
		return false
	}
}
