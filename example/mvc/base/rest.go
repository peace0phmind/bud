package base

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/peace0phmind/bud/example/mvc/infra"
	"github.com/peace0phmind/bud/stream"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Rest[T any] interface {
	InitNeedAuth(app *fiber.App)
	InitUnAuth(app *fiber.App)

	List(c *fiber.Ctx) error

	Find(c *fiber.Ctx) error
	AfterFind(m *T) error

	Create(c *fiber.Ctx) error
	BeforeCreate(m *T) error
	AfterCreate(m *T) error

	Edit(c *fiber.Ctx) error
	BeforeEdit(m *T) error
	AfterEdit(m *T) error

	Delete(c *fiber.Ctx) error
	BeforeDelete(ids []string) error
	AfterDelete(ids []string) error
}

const DefaultApiPrefix = "/api"

type BaseRest[T any] struct {
	Self     Rest[T]        `wire:"self"`
	Repo     CrudRepoInf[T] `wire:"auto"`
	cfg      *infra.Config  `wire:"auto"`
	RestName string
}

func (br *BaseRest[T]) InitUnAuth(*fiber.App) {

}

func (br *BaseRest[T]) InitNeedAuth(app *fiber.App) {
	app.Get(br.RouterUrl("list"), br.Self.List)
	app.Get(br.RouterUrl("find/:id"), br.Self.Find)
	app.Post(br.RouterUrl("create"), br.Self.Create)
	app.Post(br.RouterUrl("edit"), br.Self.Edit)
	app.Post(br.RouterUrl("delete/:ids"), br.Self.Delete)
}

func (br *BaseRest[T]) RouterUrl(subPath string) string {
	return fmt.Sprintf("%s/%s/%s", DefaultApiPrefix, br.RestName, subPath)
}

func (br *BaseRest[T]) List(c *fiber.Ctx) error {
	return br.ListWithQuery(c, nil)
}

var _defaultQueryParams = &defaultQueryParams{}

type defaultQueryParams struct{}

func (dqp *defaultQueryParams) QueryCondition(db *gorm.DB) *gorm.DB {
	return db
}

func (br *BaseRest[T]) ListWithQuery(c *fiber.Ctx, queryParams QueryCondition) (err error) {
	pageParams := &PageParams{}
	if err = c.QueryParser(pageParams); err != nil {
		return c.JSON(RetErr(RetCodeParamsError, err.Error()))
	}

	if queryParams != nil {
		if err = c.BodyParser(&queryParams); err != nil {
			return c.JSON(RetErr(RetCodeParamsError, err.Error()))
		}
	} else {
		queryParams = _defaultQueryParams
	}

	// 获取total count
	pageParams.Total, err = br.Repo.Count(func(db *gorm.DB) *gorm.DB {
		return queryParams.QueryCondition(db)
	})
	if err != nil {
		return c.JSON(RetErr(RetCodeQueryErr, err.Error()))
	}

	// check page params
	pageParams.CheckDefault()

	var pages []*T
	pages, err = br.Repo.FindAll(func(db *gorm.DB) *gorm.DB {
		return pageParams.QueryCondition(queryParams.QueryCondition(db))
	})
	if err != nil {
		return c.JSON(RetErr(RetCodeQueryErr, err.Error()))
	}

	return c.JSON(RetOkPage(pages, pageParams))

}

func (br *BaseRest[T]) Find(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(RetErr(RetCodeParamsError, "id is empty"))
	}

	m, err := br.Repo.FindById(id)
	if err != nil {
		return c.JSON(RetErr(RetCodeQueryErr, err.Error()))
	}

	if err = br.AfterFind(m); err != nil {
		return c.JSON(RetErr(RetCodeQueryErr, err.Error()))
	}

	if m != nil {
		return c.JSON(RetOk(m))
	} else {
		return c.JSON(RetErr(RetCodeRecordNotExists))
	}
}

func (br *BaseRest[T]) AfterFind(*T) error {
	return nil
}

func (br *BaseRest[T]) Create(c *fiber.Ctx) error {
	m := new(T)

	if err := c.BodyParser(m); err != nil {
		return c.JSON(RetErr(RetCodeParamsError, err.Error()))
	}

	validater := validator.New()
	if err := validater.Struct(m); err != nil {
		return c.JSON(RetErr(RetCodeCreateErr, err.Error()))
	}

	if err := br.Self.BeforeCreate(m); err != nil {
		return c.JSON(RetErr(RetCodeCreateErr, err.Error()))
	}

	mm, err := br.Repo.Create(m)
	if err != nil {
		return c.JSON(RetErr(RetCodeCreateErr, err.Error()))
	}

	if err = br.Self.AfterCreate(mm); err != nil {
		return c.JSON(RetErr(RetCodeCreateErr, err.Error()))
	}

	return c.JSON(RetOk(mm))
}

func (br *BaseRest[T]) BeforeCreate(*T) error {
	return nil
}

func (br *BaseRest[T]) AfterCreate(*T) error {
	return nil
}

func (br *BaseRest[T]) Edit(c *fiber.Ctx) error {
	m := new(T)

	if err := c.BodyParser(m); err != nil {
		return c.JSON(RetErr(RetCodeParamsError, err.Error()))
	}

	validater := validator.New()
	if err := validater.Struct(m); err != nil {
		return c.JSON(RetErr(RetCodeUpdateErr, err.Error()))
	}

	if err := br.Self.BeforeEdit(m); err != nil {
		return c.JSON(RetErr(RetCodeUpdateErr, err.Error()))
	}

	mm, err := br.Repo.Update(m)
	if err != nil {
		return c.JSON(RetErr(RetCodeUpdateErr, err.Error()))
	}

	if err = br.Self.AfterEdit(mm); err != nil {
		return c.JSON(RetErr(RetCodeUpdateErr, err.Error()))
	}

	return c.JSON(RetOk(mm))
}

func (br *BaseRest[T]) BeforeEdit(*T) error {
	return nil
}

func (br *BaseRest[T]) AfterEdit(*T) error {
	return nil
}

func (br *BaseRest[T]) Delete(c *fiber.Ctx) error {
	ids := c.Params("ids")
	if ids == "" {
		return c.JSON(RetErr(RetCodeParamsError, "id params is null"))
	}

	idList := strings.Split(ids, ",")

	if err := br.Self.BeforeDelete(idList); err != nil {
		return c.JSON(RetErr(RetCodeDeleteErr, err.Error()))
	}

	if err := br.Repo.DeleteByIds(stream.Of(idList).MustToAny()); err != nil {
		return c.JSON(RetErr(RetCodeDeleteErr, err.Error()))
	}

	if err := br.Self.AfterDelete(idList); err != nil {
		return c.JSON(RetErr(RetCodeDeleteErr, err.Error()))
	}

	return c.JSON(RetOk[any](nil))
}

func (br *BaseRest[T]) BeforeDelete([]string) error {
	return nil
}

func (br *BaseRest[T]) AfterDelete([]string) error {
	return nil
}

func (br *BaseRest[T]) Unsupported(c *fiber.Ctx, methodName string) error {
	return c.JSON(RetErr(RetCodeMethodNotSupported, methodName))
}

func (br *BaseRest[T]) JwtTokenToUser(c *fiber.Ctx) *User {
	if c.Locals("user") == nil {
		return nil
	}

	if user, ok := c.Locals("user").(*jwt.Token); !ok {
		return nil
	} else {
		claims := user.Claims.(jwt.MapClaims)
		ret := &User{}
		uid := uuid.MustParse(claims["id"].(string))
		ret.Id = &uid
		ret.Username = claims["username"].(string)
		ret.IsAdmin = claims["isAdmin"].(bool)
		ret.Valid = claims["valid"].(bool)

		return ret
	}
}

func (br *BaseRest[T]) JwtUserToToken(user *User, exp time.Time) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id.String()
	claims["username"] = user.Username
	claims["isAdmin"] = user.IsAdmin
	claims["valid"] = user.Valid
	claims["exp"] = exp.Unix()

	// Generate encoded token and send it as response.
	if t, err := token.SignedString([]byte(br.cfg.JwtSigningKey)); err != nil {
		return "", err
	} else {
		return t, nil
	}
}
