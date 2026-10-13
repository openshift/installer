## iso9660
[![GoDoc](https://godoc.org/github.com/kdomanski/iso9660?status.svg)](https://godoc.org/github.com/kdomanski/iso9660)

A package for reading and creating ISO9660

Joliet and Rock Ridge extensions are not supported.

## Examples

### Extracting an ISO

```go
package main

import (
  "log"

  "github.com/kdomanski/iso9660/util"
)

func main() {
  f, err := os.Open("/home/user/myImage.iso")
  if err != nil {
    log.Fatalf("failed to open file: %s", err)
  }
  defer f.Close()

  if err = util.ExtractImageToDirectory(f, "/home/user/target_dir"); err != nil {
    log.Fatalf("failed to extract image: %s", err)
  }
}
```

### Creating an ISO

```go
package main

import (
  "log"
  "os"

  "github.com/kdomanski/iso9660"
)

func main() {
  writer, err := iso9660.NewWriter()
  if err != nil {
    log.Fatalf("failed to create writer: %s", err)
  }
  defer writer.Cleanup()

  f, err := os.Open("/home/user/myFile.txt")
  if err != nil {
    log.Fatalf("failed to open file: %s", err)
  }
  defer f.Close()

  err = writer.AddFile(f, "folder/MYFILE.TXT")
  if err != nil {
    log.Fatalf("failed to add file: %s", err)
  }

  outputFile, err := os.OpenFile("/home/user/output.iso", os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0644)
  if err != nil {
    log.Fatalf("failed to create file: %s", err)
  }

  err = writer.WriteTo(outputFile, "testvol")
  if err != nil {
    log.Fatalf("failed to write ISO image: %s", err)
  }

  err = outputFile.Close()
  if err != nil {
    log.Fatalf("failed to close output file: %s", err)
  }
}
```
