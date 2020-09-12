package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-09-12 15:44:22.506314 +0200 CEST m=+0.020204799

import "math"

// Float32 is float32 floating-point signal.
type Float32 struct {
	buffer []float32
	channels
}

// Float32 allocates a new sequential float32 signal buffer.
func (a Allocator) Float32() Floating {
	return Float32{
		buffer:   make([]float32, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s Float32) AppendSample(value float64) Floating {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, float32(value))
	return s
}

// SetSample sets sample value for provided index.
func (s Float32) SetSample(i int, value float64) {
	s.buffer[i] = float32(value)
}

// GetFloat32 selects a new sequential float32 signal buffer.
// from the pool.
func (p PoolAllocator) GetFloat32() Floating {
	return Float32{
		buffer:   p.f32.Get().([]float32)[:p.Channels*p.Length],
		channels: channels(p.Channels),
	}
}

// PutFloat32 places signal buffer back to the pool. If a type of
// provided buffer isn't Float32 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p PoolAllocator) PutFloat32(s Floating) {
	mustSameCapacity(s.Cap(), p.Channels*p.Capacity)
	if sig, ok := s.(Float32); ok {
		buf := sig.buffer[:sig.Cap()]
		for i := range buf {
			buf[i] = 0
		}
		p.f32.Put(sig.buffer)
	} else {
		panic("pool put float32 invalid type")
	}
}

// Capacity returns capacity of a single channel.
func (s Float32) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Float32) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s Float32) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Float32) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided channel and index.
func (s Float32) Sample(i int) float64 {
	return float64(s.buffer[i])
}

// Append appends [0:Length] samples from src to current buffer and returns
// new Floating buffer. Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func (s Float32) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())
	offset := s.Len()
	if s.Cap() < s.Len()+src.Len() {
		s.buffer = append(s.buffer, make([]float32, src.Len())...)
	} else {
		s.buffer = s.buffer[:s.Len()+src.Len()]
	}
	for i := 0; i < src.Len(); i++ {
		s.SetSample(i+offset, src.Sample(i))
	}
	alignCapacity(&s.buffer, s.Channels(), s.Cap())
	return s
}

// Slice slices buffer with respect to channels.
func (s Float32) Slice(start, end int) Floating {
	start = s.BufferIndex(0, start)
	end = s.BufferIndex(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// ReadFloat32 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadFloat32(src Floating, dst []float32) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = float32(src.Sample(i))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedFloat32 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedFloat32(src Floating, dst [][]float32) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = float32(src.Sample(src.BufferIndex(c, i)))
		}
	}
	return
}

// WriteFloat32 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteFloat32(src []float32, dst Floating) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, float64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedFloat32 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedFloat32(src [][]float32, dst Floating) (written int) {
	mustSameChannels(dst.Channels(), len(src))
	// determine the length of longest nested slice
	for i := range src {
		if len(src[i]) > written {
			written = len(src[i])
		}
	}
	// limit a number of writes to the length of the buffer
	written = min(written, dst.Length())
	for c := 0; c < dst.Channels(); c++ {
		for i := 0; i < written; i++ {
			if i < len(src[c]) {
				dst.SetSample(dst.BufferIndex(c, i), float64(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}
