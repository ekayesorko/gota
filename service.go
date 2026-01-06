package gota

import (
	"errors"
	"fmt"

	"github.com/ekayesorko/gota/resterr"
	"gorm.io/gorm"
)

type CommonService[T any] struct {
	Repository CommonRepository[T]
}

func (s *CommonService[T]) Create(c Context, model ...*T) *resterr.RestError {
	err := s.Repository.Create(c, model...)
	if err != nil {
		return resterr.NewInternalServerError(err)
	}
	return nil
}

func (s *CommonService[T]) Update(c Context, selector T, model *T) *resterr.RestError {
	err := s.Repository.Update(c, selector, model)
	if err != nil {
		return resterr.NewInternalServerError(err)
	}
	return nil
}

func (s *CommonService[T]) GetByParam(c Context, param T) (*T, *resterr.RestError) {
	item, err := s.Repository.GetByParam(c, param)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, resterr.NewInternalServerError(err)
	}
	if err != nil {
		return nil, resterr.NewNotFoundError(fmt.Sprintf("%T", param))
	}
	return item, nil
}

func (s *CommonService[T]) ListByParam(c Context, param T) (ListResponse[T], *resterr.RestError) {
	items, err := s.Repository.ListByParam(c, param)
	if err != nil {
		return ListResponse[T]{}, resterr.NewInternalServerError(err)
	}
	total, err := s.Repository.CountByParam(c, param)
	if err != nil {
		return ListResponse[T]{}, resterr.NewInternalServerError(err)
	}
	return ListResponse[T]{
		Items: items,
		Total: int64(total),
	}, nil
}

func (s *CommonService[T]) FindIn(c Context, field string, value interface{}) ([]T, *resterr.RestError) {
	res, err := s.Repository.FindIn(c, field, value)
	if err != nil {
		return nil, resterr.NewInternalServerError(err)
	}
	return res, nil
}
