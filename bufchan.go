package bufchan

import (
	"container/list"
	"sync"
)

type Element list.Element

type BufChan struct {
	buf *list.List
	// need two mutexes - one for channel blocking and one for operation locking
	mut    *sync.Mutex
	empty  *sync.Mutex
	closed bool
}

func New() (ch *BufChan) {
	var mut, empty sync.Mutex

	ch = &BufChan{
		mut:   &mut,
		empty: &empty,
		buf:   list.New(),
	}

	// start out empty
	ch.empty.Lock()
	return ch
}

func (ch BufChan) Send(v interface{}) {
	ch.mut.Lock()
	// if it's empty, adding something should unblock the empty mutex.
	if ch.closed {
		ch.mut.Unlock()
		return
	}
	if ch.buf.Len() == 0 {
		ch.empty.Unlock()
	}
	ch.buf.PushBack(v)
	ch.mut.Unlock()
}

func (ch BufChan) Receive() (v interface{}, ok bool) {
	for {
		ch.empty.Lock()
		ch.mut.Lock()
		if ch.buf.Len() != 0 {
			break
		} else {
			if ch.closed {
				ch.empty.Unlock()
				ch.mut.Unlock()
				return nil, false
			}
		}
		ch.empty.Unlock()
		ch.mut.Unlock()
	}
	v = ch.buf.Remove(ch.buf.Front())
	if ch.buf.Len() != 0 || ch.closed {
		ch.empty.Unlock()
	}
	ch.mut.Unlock()
	return v, true
}

func (ch *BufChan) Close() {
	ch.mut.Lock()
	println("Setting closed to true")
	ch.closed = true
	if ch.buf.Len() == 0 {
		println("unlcoking the empty mut")
		ch.empty.Unlock()
	}
	ch.mut.Unlock()
}
