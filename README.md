#GoRunFunc

[![asciicast](https://asciinema.org/a/24190.png)](https://asciinema.org/a/24190)

 Cherry-pick go functions to run from cmd.

##Installation

```bash
# If you already have goimports, ignore the following command
go get golang.org/x/tools/cmd/goimports

go get github.com/MohamedBassem/gorunfunc
```

##Usage

In the package root folder run:
```bash
gorunfunc '<function call or code to execute>'
```

####Example

Assuming we have a package file containing:
```go
package hello

import "fmt"

func helloWorld() {
        fmt.Println("Hello World!")
}

func SayHi() {
        helloWorld()
        fmt.Println("Hi!")
}

func otherThings() {
        // ....
}
```

If we want to run `helloWorld` and `SayHi` functions, for debugging purposes for example, execute:
```bash
$ gorunfunc 'helloWorld()'
Hello World!

$ gorunfunc 'SayHi()'
Hello World!
Hi!
````

Other usages may include running arbitrary go code from the command line, for example:
```bash
$ gorunfunc 'fmt.Println(3*2)'
6
```

##How It Works

- The package name is extracted from other files in the directory.
- A go test file is generated with the input function as the test body.
- `goimports` is invoked on the test file to generate the missing imports.
- The test is executed.
- The test file is removed.

To check the generated test file without execution, run the command with `--dry-run` flag. For instance, the dry run of the example would print:
```go
$ gorunfunc --dry-run 'helloWorld()'

package hello

import "testing"

func Test_jBE07UXZ58(t *testing.T) {
        helloWorld()
}
```

##Contribution

Your contributions and ideas are welcomed through issues and pull requests.
