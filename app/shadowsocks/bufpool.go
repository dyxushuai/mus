package shadowsocks

import (
	"bytes"
)


type BufferPool struct {
	c chan *bytes.Buffer
}


func NewBufferPool(size int) (bp *BufferPool) {
	return &BufferPool{
		c: make(chan *bytes.Buffer, size),
	}
}


func (self *BufferPool) Get() (b *bytes.Buffer) {
	select {
	case b = <-self.c:
		// reuse existing buffer
	default:
		// create new buffer
		b = bytes.NewBuffer([]byte{})
	}
	return
}


func (self *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	self.c <- b
}
