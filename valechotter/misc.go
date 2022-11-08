package valechotter

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

func UuidV4(uuidInterface interface{}) error {
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

func Nullable(validator Validator) Validator {
	return func(data interface{}) error {
		if data == nil {
			return nil
		}
		return validator(data)
	}
}
