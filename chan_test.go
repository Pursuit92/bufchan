package bufchan

import "testing"

func TestChan(t *testing.T) {
	ch := New()
	ch.Send(1)
	ch.Send(2)
	ch.Send(3)
	ch.Send(4)
	ch.Send(5)

	for i := 0; i < 5; i++ {
		go func() {
			v, _ := ch.Receive()
			if v != nil {
				return
			} else {
				t.Fatal("chan returned nil")
			}
		}()
	}
}

func TestPair(t *testing.T) {
	in, out := NewPair()
	for i := 0; i < 5; i++ {
		in <- i
	}
	close(in)
	println("Sent everything")
	for i := 0; i < 5; i++ {
		<-out
	}
	_, ok := <-out
	println(ok)
}
