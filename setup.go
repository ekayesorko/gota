package gota

import (
	"sync"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type gota struct {
	server *echo.Echo
	db     *gorm.DB

	config SetupConfig
}

var gInstance *gota
var o = sync.Once{}

type SetupConfig struct {
	Auth *Auth
}

type Auth struct {
	GoogleClientId string
}

type SetupOption func(c *SetupConfig)

func WithGoogleAuth(gcId string) SetupOption {
	return func(c *SetupConfig) {
		c.Auth = &Auth{
			GoogleClientId: gcId,
		}
	}
}

func Setup(server *echo.Echo, db *gorm.DB, fs ...SetupOption) {
	o.Do(func() {
		config := &SetupConfig{}
		for _, f := range fs {
			f(config)
		}
		gInstance = &gota{
			server: server,
			db:     db,
			config: *config,
		}
		if config.Auth != nil {
			SetupAuth(server.Group("auth"))
		}
	})
}
