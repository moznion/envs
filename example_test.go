package envs

import (
	"fmt"
	"os"
)

func ExampleUnmarshal() {
	_ = os.Setenv("FOO", "string-value")
	_ = os.Setenv("BAR", "65535")
	_ = os.Setenv("BUZ", "123.456")
	_ = os.Setenv("QUX", "true")

	type StructuredEnv struct {
		Foo     string  `envs:"FOO"`
		Bar     int64   `envs:"BAR"`
		Buz     float64 `envs:"BUZ"`
		Qux     bool    `envs:"QUX"`
		FooBar  string  `envs:"FOOBAR,allowempty"`
		Nothing string
	}

	var e StructuredEnv
	err := Unmarshal(&e)
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

	type PtrStructuredEnv struct {
		Foo     *string  `envs:"FOO"`
		Bar     *int64   `envs:"BAR"`
		Buz     *float64 `envs:"BUZ"`
		Qux     *bool    `envs:"QUX"`
		FooBar  *string  `envs:"FOOBAR,allowempty"`
		Nothing *string
	}
	var pe PtrStructuredEnv
	err = Unmarshal(&pe)
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
	// structured envvar:
	//     Foo     => "string-value"
	//     Bar     => 65535
	//     Buz     => 123.456000
	//     Qux     => true
	//     FooBar  => ""
	//     Nothing => ""
	// pointer based structured envvar:
	//     Foo     => "string-value"
	//     Bar     => 65535
	//     Buz     => 123.456000
	//     Qux     => true
	//     FooBar  => <nil>
	//     Nothing => <nil>
}
