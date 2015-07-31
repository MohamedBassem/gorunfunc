#GoRunFunc

[![asciicast](https://asciinema.org/a/8bjpolygkl1za9tx8ojqozy7l.png)](https://asciinema.org/a/8bjpolygkl1za9tx8ojqozy7l)

A tool to run certain go functions in a package.

#Installation

```bash
# If you already have goimports, ignore the following command
go get golang.org/x/tools/cmd/goimports

go get github.com/MohamedBassem/gorunfunc
```

#Usage

In the package root folder run:
```bash
gorunfunc '<function call>'
```

###Example

Assuming we have a package file containing:
```go
package hello

import "fmt"

func helloWorld() {
	fmt.Println("Hello World!")
}

func otherThings() {
	// ....
}
```

If we want to run the `helloWorld` function, for debugging purposes for example, execute:
```bash
$ gorunfunc 'helloWorld()'
Hello World!
````

#How It Works

- The package name is extracted from other files in the directory.
- A go test file is generated with the input function as the test body.
- `goimports` is invoked on the test file to generate the missing imports.
- The test is executed.
- The test file is removed.

To check the generated test file without execution, run the command with `--dry-run` flag. For instance, the dry run of the example would print:
```bash
package hello

import "testing"

func Test_jBE07UXZ58(t *testing.T) {
        helloWorld()
}
```
