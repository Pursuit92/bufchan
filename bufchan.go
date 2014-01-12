package bufchan

import (
	"container/list"
	"sync"
)

type Element list.Element

type BufChan struct {
	buf *list.List
	// need two mutexes - one for channel blocking and one for operation locking
	mut * sync.Mutex
	status * sync.Mutex
	closed bool
}

func New() (ch *BufChan) {
	var mut,status sync.Mutex

	ch = new(BufChan)
	ch.mut = &mut
	ch.status = &status
	ch.buf = list.New()

	// start out empty and "open"
	ch.status.Lock()
	ch.closed = false
	return ch
}

func (ch BufChan) Send(v interface{}) {
	ch.mut.Lock()
	defer ch.mut.Unlock()
	// if it's not empty, we need to wait till any other operations are done
	if ch.buf.Front() != nil {
		ch.status.Lock()
	}
	ch.buf.PushBack(v)
	ch.status.Unlock()
}

func (ch BufChan) Receive() (v *list.Element) {
	ch.mut.Lock()
	defer ch.mut.Unlock()
	if ch.closed && ch.buf.Front() == nil {
		return nil
	}
	ch.status.Lock()
	v = ch.buf.Front()
	ch.buf.Remove(v)
	newv := ch.buf.Front()
	// if it's not empty, need to unlock
	if newv != nil {
		ch.status.Unlock()
	}
	return v
}

func (ch BufChan) Empty() (ret bool) {
	ch.mut.Lock()
	defer ch.mut.Unlock()
	if ch.buf.Front() == nil {
		ret = true
	}
	return ret
}

func (ch BufChan) Closed() bool {
	ch.mut.Lock()
	defer ch.mut.Unlock()
	closed := ch.closed
	return closed
}

func (ch *BufChan) Close() {
	ch.mut.Lock()
	defer ch.mut.Unlock()
	ch.closed = true
}
