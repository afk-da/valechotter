package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// result type
type Result struct {
	Ok  interface{}
	Err interface{}
}

func (r Result) Error() string {
	return ""
}
func Ok(value interface{}) Result {
	return Result{Ok: value, Err: nil}
}
func Err(value interface{}) Result {
	return Result{Ok: nil, Err: value}
}

// vallector type
type (
	VallectorContext struct {
		c          echo.Context
		validators map[string]VallectorValidator
	}
	VallectorValidator struct {
		kind    string
		handler func(string) Result
	}
)

func Vallector(c echo.Context) VallectorContext {
	return VallectorContext{c: c, validators: map[string]VallectorValidator{}}
}
func (r VallectorContext) Query(key string, f func(string) Result) VallectorContext {
	r.validators[key] = VallectorValidator{kind: "query", handler: f}
	return r
}
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}
	structFieldValue.Set(val)
	return nil
}
func (r VallectorContext) Validate(req interface{}) error {
	kvs := map[string]string{}
	for key, validator := range r.validators {
		var value string
		switch validator.kind {
		case "query":
			value = r.c.QueryParam(key)
		case "path":
			value = r.c.Param(key)
		}
		result := validator.handler(value)
		if result.Err != nil {
			return result
		}
		kvs[key] = value
	}
	for key, value := range kvs {
		err := SetField(req, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}
func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		var req struct {
			MerchantId string
			Sku        string
		}
		EchoValidator(&c).UuidV4Validator(query, "MerchantId")

		return c.JSON(http.StatusOK, req)
	})
	// e.GET("/", func(c echo.Context) error {

	// 	err := Vallector(c).
	// 		Query("MerchantId", func(s string) Result {
	// 			return Ok(s)
	// 		}).
	// 		Query("Sku", func(s string) Result {
	// 			return Ok(s)
	// 		}).
	// 		Validate(&req)
	// 	if err != nil {
	// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	// 	}
	// 	return c.JSON(http.StatusOK, req)
	// })
	e.Logger.Fatal(e.Start(":1323"))
}

// HandlerFunc func(c Context) error
type MyEcho struct {
	origEcho *echo.Context
}

func EchoValidator(c *echo.Context) *MyEcho {
	return &MyEcho{origEcho: c}
}

type Options int8

const (
	query = iota // 0
	body
)

func (c *MyEcho) UuidV4Validator(option Options, params ...string) error {
	localContext := *c.origEcho
	switch {
	case option == query:
		for _, param := range params {
			p := localContext.QueryParam(param)
			if strings.TrimSpace(p) != "" {
				id, err := uuid.Parse(p)
				if err != nil {
					localContext.JSON(http.StatusBadRequest, fmt.Sprintf("Given %s : %v query param is not valid.", param, id))
				}
			}
		}
	// case option == body:
	// 	for _, param := range params {
	// 		p := localContext.Bind()
	// 	}
	default:
		return nil
	}
	return nil
}
