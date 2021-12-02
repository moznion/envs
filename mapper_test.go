package envs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	type StructuredEnv struct {
		Foo     string  `envs:"FOO"`
		Bar     int64   `envs:"BAR"`
		Buz     float64 `envs:"BUZ"`
		Qux     bool    `envs:"QUX"`
		FooBar  string  `envs:"FOOBAR,allowempty"`
		Nothing string
	}

	_ = os.Setenv("FOO", "string-value")
	_ = os.Setenv("BAR", "65535")
	_ = os.Setenv("BUZ", "123.456")
	_ = os.Setenv("QUX", "true")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.NoError(t, err)
	assert.EqualValues(t, StructuredEnv{
		Foo:     "string-value",
		Bar:     65535,
		Buz:     123.456,
		Qux:     true,
		FooBar:  "",
		Nothing: "",
	}, e)
}

func TestUnmarshalWithPointerBasedStructure(t *testing.T) {
	type PtrBasedStructuredEnv struct {
		Foo    *string  `envs:"FOO"`
		Bar    *int64   `envs:"BAR"`
		Buz    *float64 `envs:"BUZ"`
		Qux    *bool    `envs:"QUX"`
		FooBar *string  `envs:"FOOBAR,allowempty"`
	}

	_ = os.Setenv("FOO", "string-value")
	_ = os.Setenv("BAR", "65535")
	_ = os.Setenv("BUZ", "123.456")
	_ = os.Setenv("QUX", "true")

	var e PtrBasedStructuredEnv
	err := Unmarshal(&e)
	assert.NoError(t, err)

	assert.EqualValues(t, PtrBasedStructuredEnv{
		Foo: func() *string {
			v := "string-value"
			return &v
		}(),
		Bar: func() *int64 {
			v := int64(65535)
			return &v
		}(),
		Buz: func() *float64 {
			v := 123.456
			return &v
		}(),
		Qux: func() *bool {
			v := true
			return &v
		}(),
		FooBar: nil,
	}, e)
}

func TestUnmarshal_WithNilValue(t *testing.T) {
	err := Unmarshal(nil)
	assert.ErrorIs(t, err, ErrNilVessel)
}

func TestUnmarshal_WithNonPointerValue(t *testing.T) {
	type StructuredEnv struct{}
	var e StructuredEnv
	err := Unmarshal(e)
	assert.ErrorIs(t, err, ErrNonPointerVesselType)
}

func TestUnmarshal_WithEmptyTagName(t *testing.T) {
	type StructuredEnv struct {
		Foo string `envs:""`
	}
	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrEnvironmentVariableNameIsEmpty)
}

func TestUnmarshal_WithEmptyEnvVar(t *testing.T) {
	type StructuredEnv struct {
		Foo string `envs:"FOO"`
	}

	_ = os.Unsetenv("FOO")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrEmptyEnvironmentVariable)
}

func TestUnmarshal_WithUnparsableInt64(t *testing.T) {
	type StructuredEnv struct {
		Foo int64 `envs:"FOO"`
	}

	_ = os.Setenv("FOO", "string-value")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrEnvironmentVariableParseInt64)
}

func TestUnmarshal_WithUnparsablePtrInt64(t *testing.T) {
	type StructuredEnv struct {
		Foo *int64 `envs:"FOO"`
	}

	_ = os.Setenv("FOO", "string-value")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrEnvironmentVariableParseInt64)
}

func TestUnmarshal_WithUnparsableFloat64(t *testing.T) {
	type StructuredEnv struct {
		Foo float64 `envs:"FOO"`
	}

	_ = os.Setenv("FOO", "string-value")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrEnvironmentVariableParseFloat64)
}

func TestUnmarshal_WithUnparsablePtrFloat64(t *testing.T) {
	type StructuredEnv struct {
		Foo *float64 `envs:"FOO"`
	}

	_ = os.Setenv("FOO", "string-value")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrEnvironmentVariableParseFloat64)
}

func TestUnmarshal_WithUnsupportedType(t *testing.T) {
	type StructuredEnv struct {
		Foo time.Time `envs:"FOO"`
	}

	_ = os.Setenv("FOO", "string-value")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrUnsupportedFieldName)
}

func TestUnmarshal_WithUnsupportedPtrType(t *testing.T) {
	type StructuredEnv struct {
		Foo *time.Time `envs:"FOO"`
	}

	_ = os.Setenv("FOO", "string-value")

	var e StructuredEnv
	err := Unmarshal(&e)
	assert.ErrorIs(t, err, ErrUnsupportedFieldName)
}
