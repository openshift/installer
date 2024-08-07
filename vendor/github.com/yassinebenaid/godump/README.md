<div align="center">

<div width="50px" height="50px">

![binoculars (3)](https://github.com/yassinebenaid/godump/assets/101285507/f2d40c7a-6f5c-4dd9-9580-3accc74efeb4)

</div>

<h1> godump </h1>
</div>

<div align="center">

[![Tests](https://github.com/yassinebenaid/godump/actions/workflows/test.yml/badge.svg)](https://github.com/yassinebenaid/godump/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/yassinebenaid/godump/graph/badge.svg?token=EAZNA85AIS)](https://codecov.io/github/yassinebenaid/godump)
[![Go Report Card](https://goreportcard.com/badge/github.com/yassinebenaid/godump)](https://goreportcard.com/report/github.com/yassinebenaid/godump)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9241/badge)](https://www.bestpractices.dev/projects/9241)
[![Version](https://badge.fury.io/gh/yassinebenaid%2Fgodump.svg)](https://badge.fury.io/gh/yassinebenaid%2Fgodump)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENCE)
[![Go Reference](https://pkg.go.dev/badge/github.com/yassinebenaid/godump.svg)](https://pkg.go.dev/github.com/yassinebenaid/godump)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go?tab=readme-ov-file#parsersencodersdecoders)


</div>

A versatile Go library designed to output any Go variable in a structured and colored format.

This library is especially useful for debugging and testing when the standard `fmt` library falls short in displaying arbitrary data effectively. It can also serve as a powerful logging adapter, providing clear and readable logs for both development and production environments.

`godump` is not here to replace the `fmt` package. Instead, it provides an extension to what the `fmt.Printf("%#v")` can do.

## Why godump

- ability to pretty print values of all types
- well formatted output
- unexported structs are dumped too
- pointers are followed, and recursive pointers are taken in mind ([see examples](#example-3))
- customizable, you have full control over the output, **you can even generate HTML if you'd like to**, [see examples](#example-4)
- zero dependencies

## Get Started

Install the library:

```bash
go get -u github.com/yassinebenaid/godump
```

Then use the **Dump** function:

```go
package main

import (
	"github.com/yassinebenaid/godump"
)

func main() {
	godump.Dump("Anything")
}

```

## Customization

If you need more control over the output. Use the `Dumper`

```go
package main

import (
	"os"

	"github.com/yassinebenaid/godump"
)

func main() {

	var v = "Foo Bar"
	var d = godump.Dumper{
		Indentation:       "  ",
		HidePrivateFields: false,
		ShowPrimitiveNamedTypes = false
		Theme: godump.Theme{
			String: godump.RGB{R: 138, G: 201, B: 38},
			// ...
		},
	}

	d.Print(v)
	d.Println(v)
	d.Fprint(os.Stdout, v)
	d.Fprintln(os.Stdout, v)
	d.Sprint(v)
	d.Sprintln(v)
}

```

## Demo

### Example 1.

```go
package main

import (
	"os"

	"github.com/yassinebenaid/godump"
)

func main() {
	godump.Dump(os.Stdout)
}

```

Output:

![stdout](./demo/stdout.png)

### Example 2.

```go
package main

import (
	"net"

	"github.com/yassinebenaid/godump"
)

func main() {
	godump.Dump(net.Dialer{})
}

```

Output:

![dialer](./demo/dialer.png)

### Example 3.

This example shows how recursive pointers are handled.

```go
package main

import (
	"github.com/yassinebenaid/godump"
)

func main() {
	type User struct {
		Name       string
		age        int
		BestFriend *User
	}

	me := User{
		Name: "yassinebenaid",
		age:  22,
	}

    // This creates a ring
	me.BestFriend = &me

	godump.Dump(me)
}
```

Output:

![pointer](./demo/pointer.png)

### Example 4.

This example emphasizes how you can generate HTML

```go
package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/yassinebenaid/godump"
)

// Define your custom style implementation
type CSSColor struct {
	R, G, B int
}

func (c CSSColor) Apply(s string) string {
	return fmt.Sprintf(`<div style="color: rgb(%d, %d, %d); display: inline-block">%s</div>`, c.R, c.G, c.B, s)
}

func main() {

	var d godump.Dumper

	d.Theme = godump.Theme{
		String:        CSSColor{138, 201, 38},// edit the theme to use your implementation
		Quotes:        CSSColor{112, 214, 255},
		Bool:          CSSColor{249, 87, 56},
		Number:        CSSColor{10, 178, 242},
		Types:         CSSColor{0, 150, 199},
		Address:       CSSColor{205, 93, 0},
		PointerTag:    CSSColor{110, 110, 110},
		Nil:           CSSColor{219, 57, 26},
		Func:          CSSColor{160, 90, 220},
		Fields:        CSSColor{189, 176, 194},
		Chan:          CSSColor{195, 154, 76},
		UnsafePointer: CSSColor{89, 193, 180},
		Braces:        CSSColor{185, 86, 86},
	}

	var html = `<pre style="background: #111; padding: 10px; color: white">`
	html += d.Sprint(net.Dialer{})
	html += "<pre>"

    // render it to browsers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	http.ListenAndServe(":8000", nil)

}

```

Output:

![theme](./demo/theme.png)

For more examples, please take a look at [dumper_test](./dumper_test.go) along with [testdata](./testdata)

## Contribution

Please read [CONTRIBUTING guidelines](.github/CONTRIBUTING.md)
