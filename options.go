package buffer

import (
	"errors"
	"fmt"
	"time"
)

const (
	invalidSize     = "size cannot be zero"
	invalidFlusher  = "flusher cannot be nil"
	invalidInterval = "interval must be greater than zero (%s)"
	invalidTimeout  = "timeout cannot be negative (%s)"
)

type (
	// Configuration options.
	Options struct {
		Size          uint
		Flusher       FlushFunc
		FlushInterval time.Duration
		PushTimeout   time.Duration
		FlushTimeout  time.Duration
		CloseTimeout  time.Duration
	}

	// FlushFunc represents a flush function.
	FlushFunc func([]interface{})

	// Option setter.
	Option func(*Options)
)

// WithSize sets the size of the buffer.
func WithSize(size uint) Option {
	return func(options *Options) {
		options.Size = size
	}
}

// WithFlusher sets the function to be called when the buffer is flushed.
func WithFlusher(flusher FlushFunc) Option {
	return func(options *Options) {
		options.Flusher = flusher
	}
}

// WithFlushInterval sets the interval between automatic flushes.
func WithFlushInterval(interval time.Duration) Option {
	return func(options *Options) {
		options.FlushInterval = interval
	}
}

// WithPushTimeout sets how long a push should wait before giving up.
func WithPushTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.PushTimeout = timeout
	}
}

// WithFlushTimeout sets how long a manual flush should wait before giving up.
func WithFlushTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.FlushTimeout = timeout
	}
}

// WithCloseTimeout sets how long
func WithCloseTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.CloseTimeout = timeout
	}
}

func validateOptions(options *Options) error {
	if options.Size == 0 {
		return errors.New(invalidSize)
	}
	if options.Flusher == nil {
		return errors.New(invalidFlusher)
	}
	if options.FlushInterval < 0 {
		return fmt.Errorf(invalidInterval, "FlushInterval")
	}
	if options.PushTimeout < 0 {
		return fmt.Errorf(invalidTimeout, "PushTimeout")
	}
	if options.FlushInterval < 0 {
		return fmt.Errorf(invalidTimeout, "FlushTimeout")
	}
	if options.CloseTimeout < 0 {
		return fmt.Errorf(invalidTimeout, "CloseTimeout")
	}

	return nil
}
