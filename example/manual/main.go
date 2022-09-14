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
