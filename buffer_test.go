package buffer_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ivanvanderbyl/buffer"
)

var _ = Describe("Buffer", func() {
	var flusher *MockFlusher

	BeforeEach(func() {
		flusher = NewMockFlusher()
	})

	Context("Constructor", func() {
		It("creates a new Buffer instance", func() {
			// act
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(10),
			)

			// assert
			Expect(sut).NotTo(BeNil())
		})

		Context("invalid options", func() {
			It("panics when provided an invalid size", func() {
				Expect(func() {
					buffer.New[int](
						nil,
						buffer.WithSize(0),
					)
				}).To(Panic())
			})

			It("panics when provided an invalid flush interval", func() {
				Expect(func() {
					buffer.New[int](
						nil,
						buffer.WithSize(1),
						buffer.WithFlushInterval(-1),
					)
				}).To(Panic())
			})

			It("panics when provided an invalid push timeout", func() {
				Expect(func() {
					buffer.New[int](
						flusher,
						buffer.WithSize(1),
						buffer.WithPushTimeout(-1),
					)
				}).To(Panic())
			})

			It("panics when provided an invalid flush timeout", func() {
				Expect(func() {
					buffer.New[int](
						flusher,
						buffer.WithSize(1),
						buffer.WithFlushTimeout(-1),
					)
				}).To(Panic())
			})

			It("panics when provided an invalid close timeout", func() {
				Expect(func() {
					buffer.New[int](
						flusher,
						buffer.WithSize(1),
						buffer.WithCloseTimeout(-1),
					)
				}).To(Panic())
			})
		})
	})

	Context("Pushing", func() {
		It("pushes items into the buffer when Push is called", func() {
			// arrange
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(3),
			)

			// act
			err1 := sut.Push(1)
			err2 := sut.Push(2)
			err3 := sut.Push(3)

			// assert
			Expect(err1).To(Succeed())
			Expect(err2).To(Succeed())
			Expect(err3).To(Succeed())
		})

		It("fails when Push cannot execute in a timely fashion", func() {
			// arrange
			flusher.Func = func() { select {} }
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(2),
				buffer.WithPushTimeout(time.Second),
			)

			// act
			err1 := sut.Push(1)
			err2 := sut.Push(2)
			err3 := sut.Push(3)

			// assert
			Expect(err1).To(Succeed())
			Expect(err2).To(Succeed())
			Expect(err3).To(MatchError(buffer.ErrTimeout))
		})

		It("fails when the buffer is closed", func() {
			// arrange
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(2),
			)
			_ = sut.Close()

			// act
			err := sut.Push(1)

			// assert
			Expect(err).To(MatchError(buffer.ErrClosed))
		})
	})

	Context("Flushing", func() {
		It("flushes the buffer when it fills up", func(done Done) {
			// arrange
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(5),
			)

			// act
			_ = sut.Push(1)
			_ = sut.Push(2)
			_ = sut.Push(3)
			_ = sut.Push(4)
			_ = sut.Push(5)

			// assert
			result := <-flusher.Done
			Expect(result.Items).To(ConsistOf(1, 2, 3, 4, 5))
			close(done)
		})

		It("flushes the buffer when the provided interval has elapsed", func(done Done) {
			// arrange
			interval := 3 * time.Second
			start := time.Now()
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(5),
				buffer.WithFlushInterval(interval),
			)

			// act
			_ = sut.Push(1)

			// assert
			result := <-flusher.Done
			Expect(result.Items).To(ConsistOf(1))
			Expect(result.Time).To(BeTemporally("~", start, interval+time.Second))
			close(done)
		}, 5)

		It("flushes the buffer when Flush is called", func(done Done) {
			// arrange
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(3),
			)
			_ = sut.Push(1)
			_ = sut.Push(2)

			// act
			err := sut.Flush()

			// assert
			result := <-flusher.Done
			Expect(err).To(Succeed())
			Expect(result.Items).To(ConsistOf(1, 2))
			close(done)
		})

		It("fails when Flush cannot execute in a timely fashion", func() {
			// arrange
			flusher.Func = func() { time.Sleep(3 * time.Second) }
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(1),
				buffer.WithFlushTimeout(time.Second),
			)
			_ = sut.Push(1)

			// act
			err := sut.Flush()

			// assert
			Expect(err).To(MatchError(buffer.ErrTimeout))
		})

		It("fails when the buffer is closed", func() {
			// arrange
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(2),
			)
			_ = sut.Close()

			// act
			err := sut.Flush()

			// assert
			Expect(err).To(MatchError(buffer.ErrClosed))
		})
	})

	Context("Closing", func() {
		It("flushes the buffer and closes it when Close is called", func(done Done) {
			// arrange
			sut := buffer.New[int](
				flusher,
				buffer.WithSize(3),
			)
			_ = sut.Push(1)
			_ = sut.Push(2)

			// act
			err := sut.Close()

			// assert
			result := <-flusher.Done
			Expect(err).To(Succeed())
			Expect(result.Items).To(ConsistOf(1, 2))
			close(done)
		})

		It("fails when Close cannot execute in a timely fashion", func() {
			// arrange
			flusher.Func = func() { time.Sleep(2 * time.Second) }

			sut := buffer.New[int](
				flusher,
				buffer.WithSize(1),
				buffer.WithCloseTimeout(time.Second),
			)
			_ = sut.Push(1)

			// act
			err := sut.Close()

			// assert
			Expect(err).To(MatchError(buffer.ErrTimeout))
		})

		It("fails when the buffer is closed", func() {
			// arrange
			flusher.Func = func() { time.Sleep(2 * time.Second) }

			sut := buffer.New[int](
				flusher,
				buffer.WithSize(1),
				buffer.WithCloseTimeout(time.Second),
			)
			_ = sut.Close()

			// act
			err := sut.Close()

			// assert
			Expect(err).To(MatchError(buffer.ErrClosed))
		})

		It("allows Close to be called again if it fails", func() {
			// arrange
			flusher.Func = func() { time.Sleep(2 * time.Second) }

			sut := buffer.New[int](
				flusher,
				buffer.WithSize(1),
				buffer.WithCloseTimeout(time.Second),
			)
			_ = sut.Push(1)

			// act
			err1 := sut.Close()
			time.Sleep(time.Second)
			err2 := sut.Close()

			// assert
			Expect(err1).To(MatchError(buffer.ErrTimeout))
			Expect(err2).To(Succeed())
		})
	})
})

type (
	MockFlusher struct {
		Done chan *WriteCall
		Func func()
	}

	WriteCall struct {
		Time  time.Time
		Items []int
	}
)

func (flusher *MockFlusher) Write(items []int) {
	call := &WriteCall{
		Time:  time.Now(),
		Items: items,
	}

	if flusher.Func != nil {
		flusher.Func()
	}

	flusher.Done <- call
}

func NewMockFlusher() *MockFlusher {
	return &MockFlusher{
		Done: make(chan *WriteCall, 1),
	}
}
