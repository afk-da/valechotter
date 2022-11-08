package valechotter

import (
	"errors"
	"fmt"
	"time"
)

type (
	dateValidatorBuilder struct {
		before *time.Time
		after  *time.Time
	}
)

func (r dateValidatorBuilder) Before(d time.Time) dateValidatorBuilder {
	r.before = &d
	return r
}

func (r dateValidatorBuilder) After(d time.Time) dateValidatorBuilder {
	r.after = &d
	return r
}

func (r dateValidatorBuilder) Build() Validator {
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
