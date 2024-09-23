package util

import "github.com/pkg/errors"

func WrapErr(err, er error) error {
	if err == nil {
		return er
	}
	return errors.Wrap(err, er.Error())
}
