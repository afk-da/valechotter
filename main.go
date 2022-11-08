package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	Validator func(interface{}) error
)

func main() {
	e := echo.New()
	var req struct {
		Name bool `json:"name"`
	}

	e.GET("/", func(c echo.Context) error {
		v := EchoValidator()

		err := v.
			Query(map[string]Validator{
				"merchantID": v.UuidV4,
			}).
			Body(v.Object(map[string]Validator{
				"name":    v.Bool,
				"friends": v.Array(v.String),
				"attr": v.Object(map[string]Validator{
					"width":  v.UuidV4,
					"height": v.UuidV4,
				}),
				"birth": v.Date().After(time.Now().Add(-24 * time.Hour)).Build(),
			})).
			Validate(c, &req)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, req)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

// HandlerFunc func(c Context) error
type MyEcho struct {
	query map[string]Validator
	body  Validator
}

func EchoValidator() MyEcho {
	return MyEcho{}
}

func (c *MyEcho) Query(query map[string]Validator) *MyEcho {
	c.query = query
	return c
}

func (c *MyEcho) Body(body Validator) *MyEcho {
	c.body = body
	return c
}

func (c MyEcho) UuidV4(uuidInterface interface{}) error {
	switch uuidStr := uuidInterface.(type) {
	case string:
		if strings.TrimSpace(uuidStr) == "" {
			return errors.New("missing parameter")
		}

		_, err := uuid.Parse(uuidStr)
		return err
	default:
		return errors.New("expected type uuid")
	}
}

func (c MyEcho) String(data interface{}) error {
	switch data.(type) {
	case string:
		return nil
	default:
		return errors.New("expected type string")
	}
}
func (c MyEcho) Float(data interface{}) error {
	switch data.(type) {
	case float64:
		return nil
	default:
		return errors.New("expected type float64")
	}
}
func (c MyEcho) Int(data interface{}) error {
	switch data.(type) {
	case int:
		return nil
	default:
		return errors.New("expected type int")
	}
}
func (c MyEcho) Bool(data interface{}) error {
	switch data.(type) {
	case bool:
		return nil
	default:
		return errors.New("expected type bool")
	}
}
func (c MyEcho) BoolExtended(data interface{}) error {
	switch d := data.(type) {
	case bool:
		return nil
	case string:
		_, err := strconv.ParseBool(d)
		if err != nil {
			return err
		}
		return nil
	case int:
		return nil
	default:
		return errors.New("expected type bool")
	}
}
func (c MyEcho) Nullable(validator Validator) Validator {
	return func(data interface{}) error {
		if data == nil {
			return nil
		}
		return validator(data)
	}
}

func (c MyEcho) Object(fields map[string]Validator) Validator {
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

func (c MyEcho) Array(validator Validator) Validator {
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

type (
	DateValidatorBuilder struct {
		before *time.Time
		after  *time.Time
	}
)

func (r DateValidatorBuilder) Before(d time.Time) DateValidatorBuilder {
	r.before = &d
	return r
}

func (r DateValidatorBuilder) After(d time.Time) DateValidatorBuilder {
	r.after = &d
	return r
}

func (r DateValidatorBuilder) Build() Validator {
	return func(data interface{}) error {
		switch data := data.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, data)
			if err != nil {
				return fmt.Errorf("expected type time.RFC3339, found %s", data)
			}

			if r.before != nil {
				if !t.Before(*r.before) {
					return fmt.Errorf("expected time to be before %s", *r.before)
				}
			}

			if r.after != nil {
				if !t.After(*r.after) {
					return fmt.Errorf("expected time to be after %s", *r.after)
				}
			}

			return nil

		default:
			return errors.New("expected type time")
		}
	}
}

func (c MyEcho) Date() DateValidatorBuilder {
	return DateValidatorBuilder{}
}

func (m *MyEcho) Validate(c echo.Context, req interface{}) error {
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

		var buf1, buf2 bytes.Buffer
		_, err := buf1.ReadFrom(c.Request().Body)
		if err != nil {
			return err
		}
		buf2.Write(buf1.Bytes())

		var data interface{}
		err = json.NewDecoder(&buf1).Decode(&data)
		if err != nil {
			return err
		}

		err = m.body(data)
		if err != nil {
			return err
		}

		err = json.NewDecoder(&buf2).Decode(req)
		return err
	}

	return nil
}
