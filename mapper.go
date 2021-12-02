package envs

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "envs"

var (
	ErrNonPointerVesselType            = errors.New("non-pointer vessel type is not allowed")
	ErrNilVessel                       = errors.New("nil vessel is not allowed")
	ErrEnvironmentVariableNameIsEmpty  = errors.New("given tag value for environment variable's name is empty")
	ErrEmptyEnvironmentVariable        = errors.New("environment variable is empty; if it would like to allow the empty value, please consider using `allowempty` option")
	ErrEnvironmentVariableParseInt64   = errors.New("failed to parse an environment variable to int64")
	ErrEnvironmentVariableParseFloat64 = errors.New("failed to parse an environment variable to float64")
	ErrUnsupportedFieldName            = errors.New("unsupported field type; supported types are string, int64, float64, bool, and the pointer for them")
)

// Unmarshal maps the environment variables to the given structure following the definition.
// A parameter of this function must be non-nil pointer type.
//
// For example, if the following structure has given to this function,
//
//	type StructuredEnv struct {
//		Foo string `envs:"FOO"`
//	}
//
// The environment variable of `FOO` is mapped into `StructuredEnv#Foo` field according to the field type.
// Currently, it supports the following field types: string, int64, float64, bool, and the pointer for them.
//
// A tag of `envs` is a specifier which environment variable to map to the field.
// As default, if the environment variable is empty it raises an error but if `allowempty` is put in a tag (e.g. `envs:"FOO,allowempty"`),
// it accepts empty (i.e. default) value.
func Unmarshal(vessel interface{}) error {
	if vessel == nil {
		return ErrNilVessel
	}

	rv := reflect.ValueOf(vessel)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("%s: %w", reflect.TypeOf(vessel), ErrNonPointerVesselType)
	}

	elem := rv.Elem()
	for i := 0; i < elem.NumField(); i++ {
		typeField := elem.Type().Field(i)
		tag := typeField.Tag

		tagValue, ok := tag.Lookup(tagName)
		if !ok { // nothing to do
			continue
		}
		splitTagValue := strings.Split(tagValue, ",")
		envVarName := strings.TrimSpace(splitTagValue[0])
		if envVarName == "" {
			return ErrEnvironmentVariableNameIsEmpty
		}
		allowEmpty := false
		if len(splitTagValue) >= 2 && strings.TrimSpace(splitTagValue[1]) == "allowempty" { // XXX
			allowEmpty = true
		}

		envVar := os.Getenv(envVarName)
		if envVar == "" {
			if allowEmpty {
				continue
			}
			return fmt.Errorf("environment variable name is %s: %w", envVarName, ErrEmptyEnvironmentVariable)
		}

		field := elem.Field(i)
		fieldKind := field.Kind()
		switch fieldKind {
		case reflect.String:
			field.SetString(envVar)
		case reflect.Int64:
			int64EnvVar, err := strconv.ParseInt(envVar, 10, 64)
			if err != nil {
				return fmt.Errorf("environment variable `%s` => `%s`: %w", envVarName, envVar, ErrEnvironmentVariableParseInt64)
			}
			field.SetInt(int64EnvVar)
		case reflect.Float64:
			float64Envvar, err := strconv.ParseFloat(envVar, 64)
			if err != nil {
				return fmt.Errorf("environment variable `%s` => `%s`: %w", envVarName, envVar, ErrEnvironmentVariableParseFloat64)
			}
			field.SetFloat(float64Envvar)
		case reflect.Bool:
			field.SetBool(strings.EqualFold(envVar, "true"))
		case reflect.Ptr:
			var v interface{}
			fieldKind := field.Type().Elem().Kind()
			switch fieldKind {
			case reflect.String:
				v = &envVar
			case reflect.Int64:
				int64EnvVar, err := strconv.ParseInt(envVar, 10, 64)
				if err != nil {
					return fmt.Errorf("environment variable `%s` => `%s`: %w", envVarName, envVar, ErrEnvironmentVariableParseInt64)
				}
				v = &int64EnvVar
			case reflect.Float64:
				float64Envvar, err := strconv.ParseFloat(envVar, 64)
				if err != nil {
					return fmt.Errorf("environment variable `%s` => `%s`: %w", envVarName, envVar, ErrEnvironmentVariableParseFloat64)
				}
				v = &float64Envvar
			case reflect.Bool:
				b := strings.EqualFold(envVar, "true")
				v = &b
			default:
				return fmt.Errorf("field type is %s: %w", fieldKind, ErrUnsupportedFieldName)
			}
			field.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("field type is %s: %w", fieldKind, ErrUnsupportedFieldName)
		}
	}

	return nil
}
