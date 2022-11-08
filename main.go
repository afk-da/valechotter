package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		v, err := EchoValidator(c)
		merchantId, err = v.UuidV4Validator(Query, "merchantId")
		err = v.UuidV4Validator(Query, "customerId")

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, req)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

// HandlerFunc func(c Context) error
type MyEcho struct {
	origEcho echo.Context
	err      error
}

func EchoValidator(c echo.Context) (*MyEcho, error) {
	return &MyEcho{origEcho: c}, nil
}

type Options int8

const (
	Query = iota // 0
	Body
)

func (c *MyEcho) UuidV4Validator(option Options, params ...string) error {
	if c.err != nil {
		return c.err
	}

	switch {
	case option == Query:
		for _, param := range params {
			p := c.origEcho.QueryParam(param)

			if strings.TrimSpace(p) == "" {
				return c.withError(map[string]interface{}{"error": "param is empty"})
			}

			id, err := uuid.Parse(p)
			if err != nil {
				return c.withError(map[string]interface{}{"error": fmt.Sprintf("Given %s : %v query param is not valid.", param, id)})
			}
		}

	default:
		panic("invalid source")
	}

	return c.err
}

func (c *MyEcho) Validate() error {
	return c.err
}

func (c *MyEcho) withError(response map[string]interface{}) error {
	c.err = c.origEcho.JSON(http.StatusBadRequest, response)
	return c.err
}
