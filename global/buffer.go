package global

import (
	"bytes"

	"github.com/LagrangeDev/LagrangeGo/utils/binary"
	"github.com/RomiChan/syncx"
)

var bufferTable syncx.Map[*bytes.Buffer, *binary.Builder]

// NewBuffer 从池中获取新 bytes.Buffer
func NewBuffer() *bytes.Buffer {
	builder := binary.NewBuilder()
	buffer := bytes.NewBuffer(builder.ToBytes())
	bufferTable.Store(buffer, builder)
	return buffer
}

// PutBuffer 将 Buffer放入池中
func PutBuffer(buf *bytes.Buffer) {
	if _, ok := bufferTable.LoadAndDelete(buf); ok {
		// binary.PutBuilder(v)
	}
}
