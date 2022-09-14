> This is a fork of <https://github.com/globocom/go-buffer> using Generics.

<p align="center">
  <img src="gopher.png">
</p>
<p align="center">
  <img src="https://img.shields.io/github/workflow/status/globocom/go-buffer/Go?style=flat-square">
  <a href="https://github.com/globocom/go-buffer/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/globocom/go-buffer?color=blue&style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/globocom/go-buffer?style=flat-square">
  <a href="https://pkg.go.dev/github.com/globocom/go-buffer">
    <img src="https://img.shields.io/badge/Go-reference-blue?style=flat-square">
  </a>
</p>

# buffer

`buffer` represents a buffer that asynchronously flushes its contents. It is useful for applications that need to aggregate data before writing it to an external storage. A buffer is flushed manually, or automatically when it becomes full or after an interval has elapsed, whichever comes first.

## Installation

```bash
go get github.com/ivanvanderbyl/buffer
```

## Examples

### Size-triggered flush

```golang
package main

import (
  "time"

  "github.com/ivanvanderbyl/buffer"
)

func main() {
  buff := buffer.New(
  // call this function when the buffer needs flushing
  buffer.Fn(func(items []string) {
   println("flushing", len(items), "items")
   for _, item := range items {
    println(item)
   }
  }),
  // buffer can hold up to 5 items
  buffer.WithSize(5),
 )
 // ensure the buffer
 defer buff.Close()

 buff.Push("item 1")
 buff.Push("item 2")
 buff.Push("item 3")
 buff.Push("item 4")
 buff.Push("item 5")
 buff.Push("item 6") // This item will be flushed by the Closer.

 println("exiting...")
}
```

### Interval-triggered flush

```golang
package main

import (
  "time"

  "github.com/ivanvanderbyl/buffer"
)

func main() {
  buff := buffer.New(
  // call this function when the buffer needs flushing
  buffer.Fn(func(items []string) {
   println("flushing", len(items), "items")
   for _, item := range items {
    println(item)
   }
  }),
  // buffer can hold up to 3 items
  buffer.WithSize(3),
  buffer.WithFlushInterval(time.Second),
 )
 // ensure the buffer
 defer buff.Close()

 buff.Push("item 1") // Flushed on timeout of 1 second
 buff.Push("item 2") // Flushed on timeout of 1 second
 time.Sleep(2 * time.Second)
 buff.Push("item 3")
 buff.Push("item 4")
 buff.Push("item 5")
 buff.Push("item 6") // Flushed on close

 println("exiting...")
}
```

### Manual flush

```golang
package main

import (
 "github.com/ivanvanderbyl/buffer"
)

func main() {
 buff := buffer.New(
  buffer.Fn(func(items []string) {
   println("flushing", len(items), "items")
   for _, item := range items {
    println(item)
   }
  }),
  // buffer can hold up to 5 items
  buffer.WithSize(5),
 )
 defer buff.Close()

 buff.Push("item 1")
 buff.Push("item 2")
 buff.Push("item 3")
 buff.Flush()

 println("done")
}
```

### Custom Flusher

```golang
package main

import (
 "github.com/ivanvanderbyl/buffer"
)

type CustomFlusher struct{}

func (f CustomFlusher) Write(items []string) {
 println("flushing", len(items), "items")
 for _, item := range items {
  println(item)
 }
}

// Verify that CustomFlusher implements the buffer.Flusher interface
var _ buffer.Flusher[string] = (*CustomFlusher)(nil)

func main() {
 flusher := CustomFlusher{}

 buff := buffer.New[string](flusher, buffer.WithSize(5))
 defer buff.Close()

 buff.Push("item 1")
 buff.Push("item 2")
 buff.Push("item 3")
}
```

## Documentation

Visit [Pkg.go.dev](https://pkg.go.dev/github.com/globocom/go-buffer) for full documentation.

## License

[MIT License](/LICENSE)
