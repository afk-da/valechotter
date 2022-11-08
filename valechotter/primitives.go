package valechotter

import (
	"errors"
	"strconv"
)

func String(data interface{}) error {
	switch data.(type) {
	case string:
		return nil
	default:
		return errors.New("expected type string")
	}
}

func Float(data interface{}) error {
	switch data.(type) {
	case float64:
		return nil
	default:
		return errors.New("expected type float64")
	}
}

func Int(data interface{}) error {
	switch data.(type) {
	case int:
		return nil
	default:
		return errors.New("expected type int")
	}
}

func Bool(data interface{}) error {
	switch data.(type) {
	case bool:
		return nil
	default:
		return errors.New("expected type bool")
	}
}

func BoolExtended(data interface{}) error {
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
