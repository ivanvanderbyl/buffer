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
