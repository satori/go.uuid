package uuid

import (
	"io"
	"sync/atomic"
	"unsafe"
)

type buffer struct {
	data    []byte
	readPos int64
}

type bytesBuffer struct {
	bufferPointer unsafe.Pointer
	bufferSize    int
	status        uint32
	src           io.Reader
}

func newBytesBuffer(bufferSize int, src io.Reader) *bytesBuffer {
	if bufferSize <= 0 || src == nil {
		return nil
	}
	return &bytesBuffer{
		bufferSize:    bufferSize,
		bufferPointer: unsafe.Pointer(&buffer{}),
		src:           src,
	}
}

func (bb *bytesBuffer) reset(p []byte) {
	atomic.StorePointer(&bb.bufferPointer, unsafe.Pointer(&buffer{
		data:    p,
		readPos: 0,
	}))
}

func (bb *bytesBuffer) getBuffer() *buffer {
	point := atomic.LoadPointer(&bb.bufferPointer)
	return (*buffer)(point)
}

func (bb *bytesBuffer) read(dst []byte) (n int, err error) {
	buff := bb.getBuffer()
	readSize := (int64)(len(dst))
	readEnd := atomic.AddInt64(&buff.readPos, readSize)
	if readEnd > (int64)(3*len(buff.data)/4) {
		go bb.readAhead()
		if readEnd >= int64(len(buff.data)) {
			return -1, io.EOF
		}
	}

	return copy(dst, buff.data[readEnd-readSize:readEnd]), nil

}
func (bb *bytesBuffer) Read(dst []byte) (n int, err error) {
	n, err = bb.read(dst)
	if err != nil || n != len(dst) {
		n, err = io.ReadFull(bb.src, dst)
	}
	return
}

func (bb *bytesBuffer) readAhead() (err error) {
	if !atomic.CompareAndSwapUint32(&bb.status, 0, 1) {
		return nil
	}
	defer atomic.StoreUint32(&bb.status, 0)

	buf := make([]byte, bb.bufferSize)
	readN, err := io.ReadFull(bb.src, buf[:])
	if err != nil {
		return err
	}
	if readN != bb.bufferSize {
		return io.ErrUnexpectedEOF
	}

	bb.reset(buf)
	return nil
}
