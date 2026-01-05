package gota

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ekayesorko/gota/resterr"
	"github.com/ekayesorko/gota/types"
	"github.com/ekayesorko/gota/types/constants"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Context struct {
	UserId types.UserId
	Role   types.Role
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

// todo (c Context) GetDB() *gorm.DB
func GetDB(c Context) *gorm.DB {
	if gInstance == nil {
		panic("Gota uninitialized")
	}
	if c.Tx == nil {
		c.Tx = gInstance.db
	}
	_tx := c.Tx.Limit(c.Limit).Offset(c.Offset)
	if c.Sort != nil {
		_tx = c.Tx.Order(fmt.Sprintf("%s %s", c.Sort.Field, c.Sort.Direction))
	}
	for key, val := range c.Params {
		_tx = _tx.Where(fmt.Sprintf("%s = ?", key), val)
	}
	return _tx
}

func GetContext(c echo.Context, fs ...ModifyContextOption) (Context, *resterr.RestError) {
	contextOption := &contextOptions{}
	for _, f := range fs {
		f(contextOption)
	}

	userId, role, err := getUser(c)
	if err != nil {
		return Context{}, err
	}
	limit, offset, err := getPagination(c)
	if err != nil {
		return Context{}, err
	}
	result := Context{
		UserId: userId,
		Role:   role,
		C:      c.Request().Context(),
		Limit:  limit,
		Offset: offset,
		Params: make(map[string]any),
	}
	if result.Role == constants.RoleAnonymous {
		return result, nil
	}

	for param, paramSpec := range contextOption.ParamSpecs {
		valStr := c.QueryParam(param)
		if valStr == "" {
			continue
		}
		val, err := paramSpec.validate(param, valStr)
		if err != nil {
			return result, err
		}
		result.Params[param] = val
	}
	return result, nil

	// organizationId, organizationRole, projectId, err := GetOrganizationAndProject(c, userId)
	// if err != nil {
	// 	return Context{}, err
	// }
	// result.OrganizationId = organizationId
	// result.OrganizationRole = organizationRole
	// result.ProjectId = projectId
}

func getPagination(c echo.Context) (int, int, *resterr.RestError) {
	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, 0, resterr.BadRequest
	}
	pageSize := c.QueryParam("page_size")
	if pageSize == "" {
		pageSize = "10"
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return 0, 0, resterr.BadRequest
	}
	limit := pageSizeInt
	offset := pageSizeInt * (pageInt - 1)
	return limit, offset, nil
}

type ParamSpecUnit struct {
	Datatype types.Datatype //string, int, float
	Enums    []string
	Min      *float64
	Max      *float64
}

func (p ParamSpecUnit) validate(key, value string) (any, *resterr.RestError) {
	switch p.Datatype {
	case constants.String:
		found := false
		if len(p.Enums) == 0 {
			found = true
		}
		for _, e := range p.Enums {
			if value == e {
				found = true
			}
		}
		if !found {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("invalid value for %s", key))
		}
		return value, nil
	case constants.Float:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, resterr.NewBadRequestError(err.Error())
		}
		if p.Min != nil && val < *p.Min {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("%s lesser than %f", key, *p.Min))
		}
		if p.Max != nil && val > *p.Max {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("%s greater than %f", key, *p.Max))
		}
		return val, nil
	case constants.Int:
		val, err := strconv.Atoi(value)
		if err != nil {
			return nil, resterr.NewBadRequestError(err.Error())
		}
		if p.Min != nil && float64(val) < *p.Min {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("%s lesser than %d", key, int(*p.Min)))
		}
		if p.Max != nil && float64(val) > *p.Max {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("%s greater than %d", key, int(*p.Max)))
		}
		return val, nil
	default:
		return nil, resterr.NewInternalServerError(fmt.Errorf("invalid key"))
	}
}

type contextOptions struct {
	ParamSpecs map[string]ParamSpecUnit
}

type ModifyContextOption func(*contextOptions)

func WithParams(params map[string]ParamSpecUnit) ModifyContextOption {
	return func(opts *contextOptions) {
		opts.ParamSpecs = params
	}
}

func getUser(c echo.Context) (types.UserId, types.Role, *resterr.RestError) {
	userId, _ := c.Get(constants.UserIdContextKey).(types.UserId)
	role, ok := c.Get(constants.RoleContextKey).(types.Role)
	if !ok {
		role = constants.RoleAnonymous
	}
	return userId, role, nil
}
