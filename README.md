# envs [![.github/workflows/check.yml](https://github.com/moznion/envs/actions/workflows/check.yml/badge.svg)](https://github.com/moznion/envs/actions/workflows/check.yml) [![codecov](https://codecov.io/gh/moznion/envs/branch/main/graph/badge.svg?token=81AO4XSLSH)](https://codecov.io/gh/moznion/envs)

a mapper of ENVironment variables to Structure for Go.

This library maps the environment variables to the struct according to the fields' types and tags.  
Currently, it supports the following field types: string, int64, float64, bool, and the pointer for them.

## Synopsis

Basic usage:

```go
import (
	"fmt"
	"os"

	"github.com/moznion/envs"
)

type StructuredEnv struct {
	Foo     string  `envs:"FOO"`
	Bar     int64   `envs:"BAR"`
	Buz     float64 `envs:"BUZ"`
	Qux     bool    `envs:"QUX"`
	FooBar  string  `envs:"FOOBAR,allowempty"`
	Nothing string
}

func main() {
	_ = os.Setenv("FOO", "string-value")
	_ = os.Setenv("BAR", "65535")
	_ = os.Setenv("BUZ", "123.456")
	_ = os.Setenv("QUX", "true")

	var e StructuredEnv
	err := envs.Unmarshal(&e)
	if err != nil {
		panic(err)
	}
	fmt.Printf("structured envvar:\n")
	fmt.Printf("    Foo     => \"%s\"\n", e.Foo)
	fmt.Printf("    Bar     => %d\n", e.Bar)
	fmt.Printf("    Buz     => %f\n", e.Buz)
	fmt.Printf("    Qux     => %v\n", e.Qux)
	fmt.Printf("    FooBar  => \"%s\"\n", e.FooBar)
	fmt.Printf("    Nothing => \"%s\"\n", e.Nothing)

	// Output:
	// structured envvar:
	//     Foo     => "string-value"
	//     Bar     => 65535
	//     Buz     => 123.456000
	//     Qux     => true
	//     FooBar  => ""
	//     Nothing => ""
}
```

Pointer based usage:

```go
import (
	"fmt"
	"os"

	"github.com/moznion/envs"
)

type PtrStructuredEnv struct {
	Foo     *string  `envs:"FOO"`
	Bar     *int64   `envs:"BAR"`
	Buz     *float64 `envs:"BUZ"`
	Qux     *bool    `envs:"QUX"`
	FooBar  *string  `envs:"FOOBAR,allowempty"`
	Nothing *string
}

func main() {
	_ = os.Setenv("FOO", "string-value")
	_ = os.Setenv("BAR", "65535")
	_ = os.Setenv("BUZ", "123.456")
	_ = os.Setenv("QUX", "true")

	var pe PtrStructuredEnv
	err = envs.Unmarshal(&pe)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pointer based structured envvar:\n")
	fmt.Printf("    Foo     => \"%s\"\n", *pe.Foo)
	fmt.Printf("    Bar     => %d\n", *pe.Bar)
	fmt.Printf("    Buz     => %f\n", *pe.Buz)
	fmt.Printf("    Qux     => %v\n", *pe.Qux)
	fmt.Printf("    FooBar  => %v\n", pe.FooBar)
	fmt.Printf("    Nothing => %v\n", pe.Nothing)

	// Output:
	// pointer based structured envvar:
	//     Foo     => "string-value"
	//     Bar     => 65535
	//     Buz     => 123.456000
	//     Qux     => true
	//     FooBar  => <nil>
	//     Nothing => <nil>
}
```

and examples are [here](./example_test.go)

## Documentations

[![GoDoc](https://godoc.org/github.com/moznion/envs?status.svg)](https://godoc.org/github.com/moznion/envs)

## Author

moznion (<moznion@mail.moznion.net>)

