# [embed empty interface{} into explicit struct and then reflect the containing struct] #go #golang #embed #struct #interface #empty #reflect

[View original Gist on GitHub](https://gist.github.com/Integralist/0bfa1c1ab84ef06a657cd6aa6b2ce25e)

**Tags:** #go #golang #embed #struct #interface #empty #reflect

## embed empty interface{} into explicit struct and then reflect the containing struct.go

```go
package main

import (
	"fmt"
	"os"
	"reflect"
)

type random struct {
	Foo int    `json_env:"FOO"`
	Bar string `json_env:"BAR"`
	Baz bool   `json_env:"BAZ"`
}

type defaults struct {
	Random interface{}
	Beep   string `json_env:"BEEP"`
	Boop   string `json_env:"BOOP"`
}

func inspect(i interface{}) {
	fmt.Printf("inspect func:\n%+v\n\n", i)

	v := reflect.ValueOf(i)
	s := v.Elem()
	t := s.Type()

	fmt.Printf("v:\n%+v\n\n", v)
	fmt.Printf("s:\n%+v\n\n", s)
	fmt.Printf("t:\n%+v\n\n", t)

	fmt.Printf("number of top level fields:\n%+v\n\n", s.NumField())

	for i := 0; i < s.NumField(); i++ {
		sf := s.Field(i)
		tf := t.Field(i)

		fmt.Printf("sf:\n%+v\n\n", sf)
		fmt.Printf("tf:\n%+v\n\n", tf)

		if tf.Name == "Random" && tf.Type.Kind() == reflect.Interface {
			fmt.Printf("we got a random interface: %+v\n\n", tf.Type.Kind())

			v := reflect.ValueOf(sf)
			
			/*
			I GET HERE AND THEN I GET STUCK 😬
			*/

			// create new pointer
			ptr := reflect.New(reflect.TypeOf(i))

			// create variable to value of pointer
			s := ptr.Elem()

			//s := v.Elem()
			//t := s.Type()

			fmt.Printf("nested v:\n%+v\n\n", v)
			fmt.Printf("nested s:\n%+v\n\n", s)
			//fmt.Printf("nested t:\n%+v\n\n", t)
		}

		envName := tf.Tag.Get("json_env")
		if len(envName) == 0 {
			fmt.Println("skipping!")
			continue
		}

		val := []byte(os.Getenv(envName))

		fmt.Printf("\nenvName:\n%+v\n", val)

		fptr := sf.Addr().Interface()

		fmt.Printf("fptr:\n%+v\n", fptr)
	}
}

func accept(r interface{}) {
	d := defaults{Random: r, Beep: "X", Boop: "Y"}

	var i interface{}
	i = &defaults{Random: r, Beep: "X", Boop: "Y"}

	fmt.Printf("random user struct:\n%+v\n\ndefault struct with user struct embedded:\n%+v\n\ni:\n%+v\n\n", r, d, i)

	inspect(i)
}

func main() {
	os.Setenv("BEEP", "1")
	os.Setenv("BOOP", "2")
	os.Setenv("FOO", "3")
	os.Setenv("BAR", "4")
	os.Setenv("BAZ", "5")

	r := &random{}
	accept(r)

}

```

