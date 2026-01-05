package gota

import (
	"sync"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type gota struct {
	server *echo.Echo
	db     *gorm.DB
}

// func (g gotaImpl) getDB() *gorm.DB {
// 	return g.db
// }

// type gota interface {
// 	getDB() *gorm.DB
// }

var gInstance *gota
var o = sync.Once{}

func Setup(server *echo.Echo, db *gorm.DB) {
	o.Do(func() {
		gInstance = &gota{
			server: server,
			db:     db,
		}
	})
}
