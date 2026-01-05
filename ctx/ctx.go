package ctx

import (
	"context"

	"github.com/ekayesorko/gota/types"
	"gorm.io/gorm"
)

type Context struct {
	UserId types.UserId
	C      context.Context
	Limit  int
	Offset int
	Sort   *Sort
	Tx     *gorm.DB
	Params map[string]any
}

type Sort struct {
	Field     string
	Direction string
}

func GetDB(c Context) *gorm.DB {
	panic("todo")
}
