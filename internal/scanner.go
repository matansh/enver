package internal

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/matansh/enver/errs"
)

const (
	base10    = 10
	bitSize64 = 64
)

func handleDuration(strVal string) (time.Duration, error) {
	durationVal, err := time.ParseDuration(strVal) // https://pkg.go.dev/time#ParseDuration
	if err != nil {
		// `DURATION=10` needs to result in 10 seconds
		if strings.Contains(err.Error(), "time: missing unit in duration") {
			intVal, err := strconv.ParseInt(strVal, base10, bitSize64)
			if err != nil {
				return 0, fmt.Errorf("%w: failed to parse '%s' as time.Duration", err, strVal)
			}
			durationVal = time.Duration(intVal) * time.Second

			return durationVal, nil
		}

		return 0, fmt.Errorf("%w: failed to parse '%s' as time.Duration", err, strVal)
	}

	return durationVal, nil
}

func scanFieldFromEnv(fieldValue reflect.Value, envVar string) error {
	strVal, ok := os.LookupEnv(envVar)
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrEnvVarNotFound, envVar)
	}

	// cleaning up common "extras" when reading strings from the env
	strVal = strings.TrimSpace(strVal)
	strVal = strings.Trim(strVal, "\"")

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(strVal)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// time.Duration is an alias of int64, and as such will have a reflected value kind of int64
		// therefor we need to compare its reflected type, which will be time.Duration
		if fieldValue.Type() == reflect.TypeOf(time.Second) {
			durationVal, err := handleDuration(strVal)
			if err != nil {
				return fmt.Errorf("%w: failed to parse var '%s' with value '%v' as time.Duration", err, envVar, fieldValue)
			}
			fieldValue.Set(reflect.ValueOf(durationVal))

			return nil
		}

		intVal, err := strconv.ParseInt(strVal, base10, bitSize64)
		if err != nil {
			return fmt.Errorf("%w: failed to parse var '%s' with value '%v' as int64", err, envVar, fieldValue)
		}
		fieldValue.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intVal, err := strconv.ParseUint(strVal, base10, bitSize64)
		if err != nil {
			return fmt.Errorf("%w: failed to parse var '%s' with value '%v' as int64", err, envVar, fieldValue)
		}
		fieldValue.SetUint(intVal)

	case reflect.Bool:
		// possible values are: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
		boolVal, err := strconv.ParseBool(strVal)
		if err != nil {
			return fmt.Errorf("%w: failed to parse var '%s' with value '%v' as bool", err, envVar, fieldValue)
		}
		fieldValue.SetBool(boolVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(strVal, bitSize64)
		if err != nil {
			return fmt.Errorf("%w: failed to parse var '%s' with value '%v' as float64", err, envVar, fieldValue)
		}
		fieldValue.SetFloat(floatVal)

	default:
		return fmt.Errorf("%w: %s is unsupported", errs.ErrUnsupportedType, fieldValue.Kind())
	}

	return nil
}

func ScanStructFromEnv(cfg interface{}, tag string) []error {
	vValue := reflect.Indirect(reflect.ValueOf(cfg))
	var errs []error
	fieldCount := vValue.NumField()
	for i := 0; i < fieldCount; i++ {
		fieldValue := vValue.Field(i)
		if fieldValue.Kind() == reflect.Struct {
			// walking down nested structs by recursively de-reflecting the pointer address of the field value
			errs = append(errs, ScanStructFromEnv(fieldValue.Addr().Interface(), tag)...)

			continue
		}
		vType := reflect.TypeOf(cfg)
		if reflect.TypeOf(cfg).Kind() == reflect.Ptr {
			vType = vType.Elem()
		}
		tagValue, hasTag := vType.Field(i).Tag.Lookup(tag)
		// struct fields that do not explicitly specify the provided tag are ignored
		if !hasTag {
			continue
		}
		// ignoring `env:""` & `env:"-"`
		if tagValue == "" || tagValue == "-" {
			continue
		}
		if err := scanFieldFromEnv(fieldValue, tagValue); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
