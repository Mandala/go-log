// Buffer-like byte slice
// Copyright (c) 2017 Fadhli Dzil Ikram

package buffer

// Buffer type wrap up byte slice built-in type
type Buffer []byte

// Appender implements append method to buffer
type Appender interface {
	Append(data []byte)
	AppendByte(data byte)
}

// Reset buffer position to start
func (b *Buffer) Reset() {
	*b = Buffer([]byte(*b)[:0])
}

// Append byte slice to buffer
func (b *Buffer) Append(data []byte) {
	*b = append(*b, data...)
}

// AppendByte to buffer
func (b *Buffer) AppendByte(data byte) {
	*b = append(*b, data)
}
