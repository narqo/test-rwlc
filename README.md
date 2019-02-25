> Last night I thought up a fun _advanced_ Go challenge (because why not)™️

Implement a `ReadWriteLineCloser` interface, that works as following:

```go
type ReadWriteLineCloser interface {
    ReadLine() (s string, err error)
    WriteLine(s string) error
    Close()
}

// case 1
var rw ReadWriteLineCloser = New()

rw.WriteLine("line 1")
rw.WriteLine("line 2") // <unspecified number of write calls>

line, _ := rw.ReadLine()
line, _ = rw.ReadLine() // <read all written lines in fifo order>

rw.Close()

// case 2
var rw ReadWriteLineCloser = New()

go func() {
    for {
        line, err := rw.ReadLine() // <read all written lines in fifo order>
        if err == ErrClosed {
            return
        }
        log.Println(line)
    }
}()

rw.WriteLine("line 1")
rw.WriteLine("line 2") // <unspecified number of write calls>

rw.Close()
```