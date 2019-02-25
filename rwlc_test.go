package rwlc

import "testing"

func TestReadWriteLineCloser(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		rw := New()

		rw.WriteLine("test 1")
		rw.WriteLine("test 2")
		rw.WriteLine("test 3")

		s, err := rw.ReadLine()
		assertNoError(t, err)
		assertEqual(t, "test 1", s)

		s, err = rw.ReadLine()
		assertNoError(t, err)
		assertEqual(t, "test 2", s)

		rw.WriteLine("test 4")

		s, err = rw.ReadLine()
		assertNoError(t, err)
		assertEqual(t, "test 3", s)

		rw.WriteLine("test 5")
		rw.WriteLine("test 6")

		s, err = rw.ReadLine()
		assertNoError(t, err)
		assertEqual(t, "test 4", s)

		rw.Close()

		err = rw.WriteLine("test X")
		assertEqual(t, ErrClosed, err)

		_, err = rw.ReadLine()
		assertEqual(t, ErrClosed, err)
	})

	t.Run("case 2", func(t *testing.T) {
		rw := New()

		lines := make(chan string, 100) // huge buffer to never block from ReadLine's goroutine
		go func() {
			for {
				s, err := rw.ReadLine()
				if err == ErrClosed {
					break
				}
				lines <- s
			}
			close(lines)
		}()

		rw.WriteLine("test 1")
		rw.WriteLine("test 2")
		rw.WriteLine("test 3")

		s := <-lines
		assertEqual(t, "test 1", s)

		s = <-lines
		assertEqual(t, "test 2", s)

		rw.WriteLine("test 4")

		s = <-lines
		assertEqual(t, "test 3", s)

		rw.WriteLine("test 5")

		s = <-lines
		assertEqual(t, "test 4", s)

		rw.Close()

		err := rw.WriteLine("test X")
		assertEqual(t, ErrClosed, err)

		s = <-lines
		assertEqual(t, "test 5", s)

		s = <-lines
		assertEqual(t, "", s)
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
