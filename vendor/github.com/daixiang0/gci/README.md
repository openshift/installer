# GCI

GCI, a tool that controls golang package import order and makes it always deterministic.

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

GCI splits all import blocks into different sections, now support five section type:
- standard: Golang official imports, like "fmt"
- custom: Custom section, use full and the longest match (match full string first, if multiple matches, use the longest one)
- default: All rest import blocks
- blank: Put blank imports together in a separate group
- dot: Put dot imports together in a separate group

The priority is standard > default > custom > blank > dot, all sections sort alphabetically inside.
By default, blank and dot sections are not used and the corresponding lines end up in the other groups.

All import blocks use one TAB(`\t`) as Indent.

Since v0.9.0, GCI always puts C import block as the first.

**Note**:

`nolint` is hard to handle at section level, GCI will consider it as a single comment.

## Installation

To download and install the highest available release version -

```shell
go install github.com/daixiang0/gci@latest
```

You may also specify a specific version, for example:

```shell
go install github.com/daixiang0/gci@v0.6.0
```

## Usage

Now GCI provides two command line methods, mainly for backward compatibility.

### New style
GCI supports three modes of operation

```shell
$ gci print -h
Print outputs the formatted file. If you want to apply the changes to a file use write instead!

Usage:
  gci print path... [flags]

Aliases:
  print, output

Flags:
  -d, --debug             Enables debug output from the formatter
  -h, --help              help for write
  -s, --section strings   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot. The default value is [standard,default].
                          standard - standard section that Golang provides officially, like "fmt"
                          Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
                          default - default section, contains all rest imports
                          blank - blank section, contains all blank imports.
      --skip-generated    Skip generated files
      --custom-order      Enable custom order of sections. If specified, make the section order the same as your configuration order. The default order is standard > default > custom > blank > dot.
```

```shell
$ gci write -h
Write modifies the specified files in-place

Usage:
  gci write path... [flags]

Aliases:
  write, overwrite

Flags:
  -d, --debug             Enables debug output from the formatter
  -h, --help              help for write
  -s, --section strings   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot. The default value is [standard,default].
                          standard - standard section that Golang provides officially, like "fmt"
                          Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
                          default - default section, contains all rest imports
                          blank - blank section, contains all blank imports.
                          dot - dot section, contains all dot imports.
      --skip-generated    Skip generated files
      --custom-order      Enable custom order of sections. If specified, make the section order the same as your configuration order. The default order is standard > default > custom > blank > dot.
```

```shell
$ gci diff -h
Diff prints a patch in the style of the diff tool that contains the required changes to the file to make it adhere to the specified formatting.

Usage:
  gci diff path... [flags]

Flags:
  -d, --debug             Enables debug output from the formatter
  -h, --help              help for write
  -s, --section strings   Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot. The default value is [standard,default].
                          standard - standard section that Golang provides officially, like "fmt"
                          Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
                          default - default section, contains all rest imports
                          blank - blank section, contains all blank imports.
                          dot - dot section, contains all dot imports.
      --skip-generated    Skip generated files
      --custom-order      Enable custom order of sections. If specified, make the section order the same as your configuration order. The default order is standard > default > custom > blank > dot.
```

### Old style

```shell
Usage:
  gci [-diff | -write] [--local localPackageURLs] path... [flags]

Flags:
  -d, --diff            display diffs instead of rewriting files
  -h, --help            help for gci
  -l, --local strings   put imports beginning with this string after 3rd-party packages, separate imports by comma
  -v, --version         version for gci
  -w, --write           write result to (source) file instead of stdout

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

## TODO

- Ensure only one blank between `Name` and `Path` in an import block
- Ensure only one blank between `Path` and `Comment` in an import block
- Format comments
- Add more testcases
- Support imports completion (please use `goimports` first then use GCI)
- Optimize comments
