Retry
==================

[![Build Status](https://travis-ci.org/thedevsaddam/retry.svg?branch=master)](https://travis-ci.org/thedevsaddam/retry)
[![Project status](https://img.shields.io/badge/version-1.2-green.svg)](https://github.com/thedevsaddam/retry/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/thedevsaddam/retry)](https://goreportcard.com/report/github.com/thedevsaddam/retry)
[![Coverage Status](https://coveralls.io/repos/github/thedevsaddam/retry/badge.svg?branch=master)](https://coveralls.io/github/thedevsaddam/retry?branch=master)
[![GoDoc](https://godoc.org/github.com/thedevsaddam/retry?status.svg)](https://pkg.go.dev/github.com/thedevsaddam/retry)
[![License](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/thedevsaddam/retry/blob/dev/LICENSE.md)


Simple and easy retry mechanism package for Go

### Installation

Install the package using
```go
$ go get github.com/thedevsaddam/retry
```

### Usage

To use the package import it in your `*.go` code
```go
import "github.com/thedevsaddam/retry"
```

### Example

Simply retry a function to execute for max 10 times with interval of 1 second

```go

package main

import (
	"fmt"
	"time"

	"github.com/thedevsaddam/retry"
)

func main() {
	i := 1 // lets assume we expect i to be a value of 8
	err := retry.DoFunc(10, 1*time.Second, func() error {
		fmt.Printf("trying for: %dth time\n", i)
		i++
		if i > 7 {
			return nil
		}
		return fmt.Errorf("i = %d is still low value", i)
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Got our expected result: ", i)
}

```

We can execute function from other package with arguments and return values

```go

package main

import (
	"errors"
	"log"
	"time"

	"github.com/thedevsaddam/retry"
)

func div(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Can not divide by zero")
	}
	return a / b, nil
}

func main() {
	a := 20.6
	b := 3.7 // if we assign 0.0 to b, it will cause an error and will retry for 3 times
	res, err := retry.Do(3, 5*time.Second, div, a, b)
	if err != nil {
		panic(err)
	}
	log.Println(res)
}

```

### **Contribution**
If you are interested to make the package better please send pull requests or create an issue so that others can fix. Read the [contribution guide here](CONTRIBUTING.md). 

### **License**
The **retry** is an open-source software licensed under the [MIT License](LICENSE.md).
