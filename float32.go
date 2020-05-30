package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-29 23:12:40.774145 +0200 CEST m=+0.021472626

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

// SetSample sets sample value for provided position.
func (s Float32) SetSample(pos int, value float64) {
	s.buffer[pos] = float32(value)
}

// GetFloat32 selects a new sequential float32 signal buffer.
// from the pool.
func (p *Pool) GetFloat32() Floating {
	if p == nil {
		return nil
	}
	return p.f32.Get().(Floating)
}

// PutFloat32 places signal buffer back to the pool. If a type of
// provided buffer isn't Float32 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p *Pool) PutFloat32(s Floating) {
	if p == nil {
		return
	}
	if _, ok := s.(Float32); !ok {
		panic("pool put float32 invalid type")
	}
	mustSameCapacity(s.Capacity(), p.allocator.Capacity)
	p.f32.Put(s.Slice(0, p.allocator.Length))
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

// Sample returns signal value for provided channel and position.
func (s Float32) Sample(pos int) float64 {
	return float64(s.buffer[pos])
}

// Append appends [0:Length] samples from src to current buffer and returns new
// Floating buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s Float32) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())

	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]float32, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Floating(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// Slice slices buffer with respect to channels.
func (s Float32) Slice(start, end int) Floating {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// ReadFloat32 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadFloat32(src Floating, dst []float32) int {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = float32(src.Sample(pos))
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
	for channel := 0; channel < src.Channels(); channel++ {
		length := min(len(dst[channel]), src.Length())
		if length > read {
			read = length
		}
		for pos := 0; pos < length; pos++ {
			dst[channel][pos] = float32(src.Sample(src.ChannelPos(channel, pos)))
		}
	}
	return
}

// WriteFloat32 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteFloat32(src []float32, dst Floating) int {
	length := min(dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst.SetSample(pos, float64(src[pos]))
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
	for channel := 0; channel < dst.Channels(); channel++ {
		for pos := 0; pos < written; pos++ {
			if pos < len(src[channel]) {
				dst.SetSample(dst.ChannelPos(channel, pos), float64(src[channel][pos]))
			} else {
				dst.SetSample(dst.ChannelPos(channel, pos), 0)
			}
		}
	}
	return
}
