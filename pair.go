package bufchan

type pair struct {
	send    chan<- interface{}
	receive <-chan interface{}
	buf     *BufChan
}

func NewPair() (send chan<- interface{}, recv <-chan interface{}) {
	sendC := make(chan interface{})
	recvC := make(chan interface{})
	buf := New()

	go func() {
		for {
			v, ok := <-sendC
			if ok {
				buf.Send(v)
			} else {
				println("sender closed, closing chan")
				buf.Close()
				return
			}
		}
	}()

	go func() {
		for {
			v, ok := buf.Receive()
			if ok {
				recvC <- v
			} else {
				println("bufchan closed, closing receiver")
				close(recvC)
				return
			}
		}
	}()
	return sendC, recvC
}
