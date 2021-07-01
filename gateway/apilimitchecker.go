package gateway

import (
	"errors"
)

func APILimitChecker(key, path string) error {
	for _, field := range APIData[key] {
		if field.Template.JoinPath() == path {
			switch max := field.Max.(type) {
			case int:
				if field.Num >= max {
					return errors.New("limit exceeded")
				}
			case string:
				if max != "-" {
					return errors.New("unexpected limit value")
				}
			default:
				return errors.New("unexpected limit value")
			}

			return nil
		}
	}

	return nil
}
