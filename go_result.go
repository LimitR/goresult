package goresult

import (
	"errors"
	"log"
)

type Result[T any] struct {
	Value T
	Err   error
}

func NewResult[T any](value T) *Result[T] {
	return &Result[T]{
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

func (s *Result[T]) UnwrapOrElse(value T) T {
	if s.Err != nil {
		s.Err = nil
		s.Value = value
	}
	return s.Value
}

func (s *Result[T]) UnwrapOrOn(callback func(error) T) T {
	if s.Err != nil {
		res := callback(s.Err)
		s.Value = res
		return res
	}
	return s.Value
}

func (s *Result[T]) Some(value T) *Result[T] {
	s.Value = value
	s.Err = nil
	return s
}

func (s *Result[T]) Error(value string) *Result[T] {
	s.Err = errors.New(value)
	return s
}

func (s *Result[T]) Match(err error) error {
	if errors.Is(s.Err, err) {
		return s.Err
	} else {
		return nil
	}
}

func (s *Result[T]) IsOk() bool {
	return s.Err == nil
}
