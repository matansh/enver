package config

import (
	"fmt"
	"reflect"

	"github.com/matansh/envloader/errs"
	"github.com/matansh/envloader/internal"
)

const defaultStructTag = "env"

func LoadEnvFromTag(cfg interface{}, tag string) []error {
	if tag == "" {
		tag = defaultStructTag
	}
	// if the config is not a pointer then any alterations will be unique to our namespace
	if reflect.TypeOf(cfg).Kind() != reflect.Ptr {
		return []error{fmt.Errorf("%w: cfg is not a pointer: %T", errs.ErrNotPtr, cfg)}
	}
	// in order to populate struct fields we need to be provided a struct
	if reflect.TypeOf(cfg).Elem().Kind() != reflect.Struct {
		return []error{fmt.Errorf("%w: cfg is not a struct: %T", errs.ErrNotStruct, cfg)}
	}

	return internal.ScanStructFromEnv(cfg, tag)
}

// LoadEnv is a convenience method for using the default struct tag.
func LoadEnv(cfg interface{}) []error {
	return LoadEnvFromTag(cfg, defaultStructTag)
}
