package gota

import (
	"fmt"
)

type CommonRepository[T any] struct{}

func (r *CommonRepository[T]) Create(c Context, model ...*T) error {
	return GetDB(c).Create(model).Error
}

func (r *CommonRepository[T]) GetByParam(c Context, param T) (*T, error) {
	var model T
	err := GetDB(c).Where(param).First(&model).Error
	return &model, err
}

func (r *CommonRepository[T]) ListByParam(c Context, param T) ([]T, error) {
	var models []T
	err := GetDB(c).Where(param).Find(&models).Error
	return models, err
}

func (r *CommonRepository[T]) CountByParam(c Context, param T) (int64, error) {
	var count int64
	var model T
	err := GetDB(c).Model(&model).Where(param).Count(&count).Error
	return count, err
}

func (r *CommonRepository[T]) FindIn(c Context, field string, value interface{}) ([]T, error) {
	var models []T
	err := GetDB(c).Where(fmt.Sprintf("%s in ?", field), value).Find(&models).Error
	return models, err
}

func (r *CommonRepository[T]) Update(c Context, selector T, model *T) error {
	return GetDB(c).Where(selector).Updates(model).Error
}
