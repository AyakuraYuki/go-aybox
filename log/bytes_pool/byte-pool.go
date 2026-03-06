package bytes_pool

import (
	"strconv"
	"sync"
)

// Pool is type-safed sync.Pool
type Pool struct {
	p *sync.Pool
}

func NewPool() Pool {
	return Pool{
		p: &sync.Pool{
			New: func() interface{} {
				return &Buffer{bs: make([]byte, 0)}
			},
		},
	}
}

func (p Pool) Get() *Buffer {
	buf := p.p.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p Pool) put(buf *Buffer) {
	p.p.Put(buf)
}

// Buffer is a thin wrapper around a byte slice. It's intended to be pooled, so
// the only way to construct one is via a Pool.
type Buffer struct {
	bs   []byte
	pool Pool
}

func (b *Buffer) AppendByte(v byte)     { b.bs = append(b.bs, v) }
func (b *Buffer) AppendString(s string) { b.bs = append(b.bs, s...) }
func (b *Buffer) AppendInt(i int64)     { b.bs = strconv.AppendInt(b.bs, i, 10) }
func (b *Buffer) AppendUint(i uint64)   { b.bs = strconv.AppendUint(b.bs, i, 10) }
func (b *Buffer) AppendBool(v bool)     { b.bs = strconv.AppendBool(b.bs, v) }
func (b *Buffer) AppendFloat(f float64, bitSize int) {
	b.bs = strconv.AppendFloat(b.bs, f, 'f', -1, bitSize)
}

// Len returns the length of the underlying byte slice.
func (b *Buffer) Len() int {
	return len(b.bs)
}

// Cap returns the capacity of the underlying byte slice.
func (b *Buffer) Cap() int {
	return cap(b.bs)
}

// Bytes returns a mutable reference to the underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.bs
}

// String returns a string copy of the underlying byte slice.
func (b *Buffer) String() string {
	return string(b.bs)
}

// Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// backing array.
func (b *Buffer) Reset() {
	b.bs = b.bs[:0]
}

// Write implements io.Writer.
func (b *Buffer) Write(bs []byte) (int, error) {
	b.bs = append(b.bs, bs...)
	return len(bs), nil
}

// Free returns the Buffer to its Pool.
//
// Callers must not retain references to the Buffer after calling Free.
func (b *Buffer) Free() {
	b.pool.put(b)
}
