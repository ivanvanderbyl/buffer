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
