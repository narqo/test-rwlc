package rwlc

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

const n = 5

func ExampleReadWriteLineCloser_case1() {
	rw := New()

	for i := 0; i < n; i++ {
		rw.WriteLine(fmt.Sprintf("1 line %d", i))
	}

	for i := 0; i < n; i++ {
		line, err := rw.ReadLine() // <read all written lines in fifo order>
		if err != nil {
			panic(err)
		}
		fmt.Println(line)
	}

	// Output:
	// 1 line 0
	// 1 line 1
	// 1 line 2
	// 1 line 3
	// 1 line 4
}

func ExampleReadWriteLineCloser_case2() {
	rw := New()

	done := make(chan struct{})
	written := uint32(0)

	// read in goroutine
	go func() {
		for {
			line, err := rw.ReadLine() // <read all written lines in fifo order>
			if err == ErrClosed {
				break
			}
			fmt.Println(line)
			atomic.AddUint32(&written, 1)
		}
	}()

	go func() {
		// wait lines to be observed by ReadLine's goroutine
		for atomic.LoadUint32(&written) < n {
			runtime.Gosched()
		}
		close(done)
	}()

	// then write
	for i := 0; i < n; i++ {
		rw.WriteLine(fmt.Sprintf("2 line %d", i))
	}

	<-done

	// Output:
	// 2 line 0
	// 2 line 1
	// 2 line 2
	// 2 line 3
	// 2 line 4
}
