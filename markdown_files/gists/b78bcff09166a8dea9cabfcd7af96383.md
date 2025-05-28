# Go: sharing behaviours between individual structs and a base struct 

[View original Gist on GitHub](https://gist.github.com/Integralist/b78bcff09166a8dea9cabfcd7af96383)

**Tags:** #go #golang #struct #embedding #polymorphism

## Golang sharing behaviours between individual structs and a base struct.go

```go
package main

import (
	"fmt"
)

type JavaScript struct {
	postProcessing bool
	toolchain      string
}

func (j JavaScript) Validate() {
	fmt.Printf("validating with %s\n", j.toolchain)

	if j.postProcessing {
		fmt.Println("post processing")
	}
}

func (j JavaScript) Build() {
	fmt.Printf("do our own JavaScript thing with %s\n", j.toolchain)
}

func (j *JavaScript) Toolchain(toolchain string) {
	j.toolchain = toolchain
}

type AssemblyScript struct {
	JavaScript
}

func (a AssemblyScript) Build() {
	fmt.Printf("do our own AssemblyScript thing with %s\n", a.toolchain)
}

func main() {
	js := JavaScript{true, "yarn"}
	js.Validate() // validating with yarn + post processing
	js.Build()    // do our own JavaScript thing with yarn

	as := AssemblyScript{JavaScript{false, "npm"}}
	as.Validate() // validating with npm
	as.Build()    // do our own AssemblyScript thing with npm

	as2 := AssemblyScript{JavaScript{false, "npm"}}
	as2.Toolchain("beepboop")
	as2.Validate() // validating with beepboop
	as2.Build()    // do our own AssemblyScript thing with beepboop
}
```

