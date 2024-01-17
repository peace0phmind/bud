package structure

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

func init() {
	RegisterMapper[string, time.Duration](string2durationMapper)
	RegisterMapper[string, url.URL](string2urlMapper)
}

func string2durationMapper(from reflect.Value, to reflect.Value) error {
	s, err := time.ParseDuration(from.String())
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse duration: %v", err))
	}
	to.Set(reflect.ValueOf(s))
	return nil
}

func string2urlMapper(from reflect.Value, to reflect.Value) error {
	u, err := url.Parse(from.String())
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse URL: %v", err))
	}
	to.Set(reflect.ValueOf(*u))
	return nil
}
