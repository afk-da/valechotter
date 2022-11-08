package valechotter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	MyEcho struct {
		query map[string]Validator
		body  Validator
	}

	Validator func(interface{}) error

	M map[string]Validator
	A []Validator
)

func EchoValidator() *MyEcho {
	return &MyEcho{}
}

func (c *MyEcho) Query(query map[string]Validator) *MyEcho {
	c.query = query
	return c
}

func (c *MyEcho) Body(body Validator) *MyEcho {
	c.body = body
	return c
}

func (m *MyEcho) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := m.Validate(c)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			return next(c)
		}
	}
}

func (m *MyEcho) Validate(c echo.Context) error {
	// validate query params
	for key, validator := range m.query {
		value := c.QueryParam(key)
		err := validator(value)
		if err != nil {
			return err
		}
	}

	// validate body
	if m.body != nil {
		if c.Request().Body == nil {
			return errors.New("body is required but is not present")
		}

		var buf bytes.Buffer
		_, err := buf.ReadFrom(c.Request().Body)
		if err != nil {
			return err
		}

		var data interface{}
		err = json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&data)
		if err != nil {
			return err
		}

		err = m.body(data)
		if err != nil {
			return err
		}

		c.Request().Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	}

	return nil
}

func Object(fields map[string]Validator) Validator {
	return func(data interface{}) error {
		switch data := data.(type) {
		case map[string]interface{}:
			for key, validator := range fields {
				err := validator(data[key])
				if err != nil {
					return fmt.Errorf("%s: %v", key, err)
				}
			}
			return nil

		default:
			return errors.New("expected type object")
		}
	}
}

func Array(validator Validator) Validator {
	return func(data interface{}) error {
		switch data := data.(type) {
		case []interface{}:
			for i, item := range data {
				err := validator(item)
				if err != nil {
					return fmt.Errorf("at index %d: %v", i, err)
				}
			}
			return nil

		default:
			return errors.New("expected type array")
		}
	}
}

func Date() dateValidatorBuilder {
	return dateValidatorBuilder{}
}
