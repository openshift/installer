# GCI

GCI, a tool that controls Go package import order and makes it always deterministic.

The desired output format is highly configurable and allows for more custom formatting than `goimport` does.

GCI considers a import block based on AST as below:

```
Doc
Name Path Comment
```

All comments will keep as they were, except the isolated comment blocks.

The isolated comment blocks like below:

```
import (
  "fmt"
  // this line is isolated comment

  // those lines belong to one
  // isolated comment blocks

  "github.com/daixiang0/gci"
)
```

GCI splits all import blocks into different sections, now support six section type:

- standard: Go official imports, like "fmt"
- custom: Custom section, use full and the longest match (match full string first, if multiple matches, use the longest one)
- default: All rest import blocks
- blank: Put blank imports together in a separate group
- dot: Put dot imports together in a separate group
- alias: Put alias imports together in a separate group
- localmodule: Put imports from local packages in a separate group

The priority is standard > default > custom > blank > dot > alias > localmodule, all sections sort alphabetically inside.
By default, blank, dot, and alias sections are not used, and the corresponding lines end up in the other groups.

All import blocks use one TAB(`\t`) as Indent.

Since v0.9.0, GCI always puts C import block as the first.

**Note**:

`nolint` is hard to handle at section level, GCI will consider it as a single comment.

### LocalModule

Local module detection is done via reading the module name from the `go.mod`
file in *the directory where `gci` is invoked*. This means:

  - This mode works when `gci` is invoked from a module root (i.e. directory
    containing `go.mod`)
  - This mode doesn't work with a multi-module setup, i.e. when `gci` is invoked
    from a directory containing `go.work` (though it would work if invoked from
    within one of the modules in the workspace)

## Installation

To download and install the highest available release version -

```shell
go install github.com/daixiang0/gci@latest
```

You may also specify a specific version, for example:

```shell
go install github.com/daixiang0/gci@v0.11.2
```

## Usage

Now GCI provides two command line methods, mainly for backward compatibility.

### New style

GCI supports three modes of operation

> **Note**
>
> Since v0.10.0, the `-s` and `--section` flag can only be used multiple times to specify multiple sections.
> For example, you could use `-s standard,default` before, but now you must use `-s standard -s default`.
> This breaking change makes it possible for the project to support specifying multiple custom prefixes. (Please see below.)

```shell
$ gci print -h
Print outputs the formatted file. If you want to apply the changes to a file use write instead!

Usage:
  gci print path... [flags]

Aliases:
  print, output

Flags:
      --custom-order          Enable custom order of sections
  -d, --debug                 Enables debug output from the formatter
  -h, --help                  help for print
  -s, --section stringArray   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot > alias > localmodule. The default value is [standard,default].
                              standard - standard section that Go provides officially, like "fmt"
                              Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
                              default - default section, contains all rest imports
                              blank - blank section, contains all blank imports.
                              dot - dot section, contains all dot imports. (default [standard,default])
                              alias - alias section, contains all alias imports.
                              localmodule: localmodule section, contains all imports from local packages
      --skip-generated        Skip generated files
      --skip-vendor           Skip files inside vendor directory
```

```shell
$ gci write -h
Write modifies the specified files in-place

Usage:
  gci write path... [flags]

Aliases:
  write, overwrite

Flags:
      --custom-order          Enable custom order of sections
  -d, --debug                 Enables debug output from the formatter
  -h, --help                  help for write
  -s, --section stringArray   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot > alias > localmodule. The default value is [standard,default].
                              standard - standard section that Go provides officially, like "fmt"
                              Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
                              default - default section, contains all rest imports
                              blank - blank section, contains all blank imports.
                              dot - dot section, contains all dot imports. (default [standard,default])
                              alias - alias section, contains all alias imports.
                              localmodule: localmodule section, contains all imports from local packages
      --skip-generated        Skip generated files
      --skip-vendor           Skip files inside vendor directory
```

```shell
$ gci list -h
Prints the filenames that need to be formatted. If you want to show the diff use diff instead, and if you want to apply the changes use write instead

Usage:
  gci list path... [flags]

Flags:
      --custom-order          Enable custom order of sections
  -d, --debug                 Enables debug output from the formatter
  -h, --help                  help for list
  -s, --section stringArray   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot > alias > localmodule. The default value is [standard,default].
                              standard - standard section that Go provides officially, like "fmt"
                              Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
                              default - default section, contains all rest imports
                              blank - blank section, contains all blank imports.
                              dot - dot section, contains all dot imports. (default [standard,default])
                              alias - alias section, contains all alias imports.
                              localmodule: localmodule section, contains all imports from local packages
      --skip-generated        Skip generated files
      --skip-vendor           Skip files inside vendor directory
```

```shell
$ gci diff -h
Diff prints a patch in the style of the diff tool that contains the required changes to the file to make it adhere to the specified formatting.

Usage:
  gci diff path... [flags]

Flags:
      --custom-order          Enable custom order of sections
  -d, --debug                 Enables debug output from the formatter
  -h, --help                  help for diff
  -s, --section stringArray   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot > alias > localmodule. The default value is [standard,default].
                              standard - standard section that Go provides officially, like "fmt"
                              Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
                              default - default section, contains all rest imports
                              blank - blank section, contains all blank imports.
                              dot - dot section, contains all dot imports. (default [standard,default])
                              alias - alias section, contains all alias imports.
                              localmodule: localmodule section, contains all imports from local packages
      --skip-generated        Skip generated files
      --skip-vendor           Skip files inside vendor directory
```

### Old style

```shell
Gci enables automatic formatting of imports in a deterministic manner
If you want to integrate this as part of your CI take a look at golangci-lint.

Usage:
  gci [-diff | -write] [--local localPackageURLs] path... [flags]
  gci [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  diff        Prints a git style diff to STDOUT
  help        Help about any command
  list        Prints filenames that need to be formatted to STDOUT
  print       Outputs the formatted file to STDOUT
  write       Formats the specified files in-place

Flags:
  -d, --diff            display diffs instead of rewriting files
  -h, --help            help for gci
  -l, --local strings   put imports beginning with this string after 3rd-party packages, separate imports by comma
  -v, --version         version for gci
  -w, --write           write result to (source) file instead of stdout

Use "gci [command] --help" for more information about a command.
```

**Note**::

The old style is only for local tests, will be deprecated, please uses new style, `golangci-lint` uses new style as well.

## Examples

Run `gci write -s standard -s default -s "prefix(github.com/daixiang0/gci)" main.go` and you will handle following cases:

### simple case

```go
package main
import (
  "golang.org/x/tools"

  "fmt"

  "github.com/daixiang0/gci"
)
```

to

```go
package main
import (
    "fmt"

    "golang.org/x/tools"

    "github.com/daixiang0/gci"
)
```

### with alias

```go
package main
import (
  "fmt"
  go "github.com/golang"
  "github.com/daixiang0/gci"
)
```

to

```go
package main
import (
  "fmt"

  go "github.com/golang"

  "github.com/daixiang0/gci"
)
```

### with blank and dot grouping enabled

```go
package main
import (
  "fmt"
  go "github.com/golang"
  _ "github.com/golang/blank"
  . "github.com/golang/dot"
  "github.com/daixiang0/gci"
  _ "github.com/daixiang0/gci/blank"
  . "github.com/daixiang0/gci/dot"
)
```

to

```go
package main
import (
  "fmt"

  go "github.com/golang"

  "github.com/daixiang0/gci"

  _ "github.com/daixiang0/gci/blank"
  _ "github.com/golang/blank"

  . "github.com/daixiang0/gci/dot"
  . "github.com/golang/dot"
)
```

### with alias grouping enabled

```go
package main

import (
	testing "github.com/daixiang0/test"
	"fmt"

	g "github.com/golang"

	"github.com/daixiang0/gci"
	"github.com/daixiang0/gci/subtest"
)
```

to

```go
package main

import (
	"fmt"

	"github.com/daixiang0/gci"
	"github.com/daixiang0/gci/subtest"

	testing "github.com/daixiang0/test"
	g "github.com/golang"
)
```

### with localmodule grouping enabled

Assuming this is run on the root of this repo (i.e. where
`github.com/daixiang0/gci` is a local module)

```go
package main

import (
	"os"
	"github.com/daixiang0/gci/cmd/gci"
)
```

to

```go
package main

import (
	"os"

	"github.com/daixiang0/gci/cmd/gci"
)
```

## TODO

- Ensure only one blank between `Name` and `Path` in an import block
- Ensure only one blank between `Path` and `Comment` in an import block
- Format comments
- Add more testcases
- Support imports completion (please use `goimports` first then use GCI)
- Optimize comments
- Remove Analyzer layer and fully use analyzer syntax
