package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-09-15 15:22:49.764301 +0200 CEST m=+0.008746823

import (
	"math"
	"sync"
)

// Int8 is int8 Signed fixed-point signal.
type Int8 struct {
	buffer []int8
	channels
	bitDepth
}

// Int8 allocates a new sequential int8 signal buffer.
func (a Allocator) Int8(bd BitDepth) Signed {
	return &Int8{
		buffer:   make([]int8, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth8),
	}
}

// GetInt8 selects a new sequential int8 signal buffer.
// from the pool.
func (p *PoolAllocator) GetInt8(bd BitDepth) Signed {
	s := p.i8.Get().(*Int8)
	s.channels = channels(p.Channels)
	s.buffer = s.buffer[:p.Length*p.Channels]
	s.bitDepth = limitBitDepth(bd, BitDepth8)
	return s
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
// Sample values are capped by maximum value of the buffer bit depth.
func (s *Int8) AppendSample(value int64) {
	if len(s.buffer) == cap(s.buffer) {
		return
	}
	s.buffer = append(s.buffer, int8(s.BitDepth().SignedValue(value)))
}

// SetSample sets sample value for provided index.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Int8) SetSample(i int, value int64) {
	s.buffer[i] = int8(s.BitDepth().SignedValue(value))
}

// PutInt8 places signal buffer back to the pool. If a type of
// provided buffer isn't Int8 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (s *Int8) Free(p *PoolAllocator) {
	mustSameCapacity(s.Cap(), p.Channels*p.Capacity)
	for i := range s.buffer {
		s.buffer[i] = 0
	}
	p.i8.Put(s)
}

func int8pool(size int) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &Int8{
				buffer: make([]int8, 0, size),
			}
		},
	}
}

// Capacity returns capacity of a single channel.
func (s *Int8) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s *Int8) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s *Int8) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s *Int8) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided channel and index.
func (s *Int8) Sample(i int) int64 {
	return int64(s.buffer[i])
}

// Append appends [0:Length] samples from src to current buffer and returns
// new Signed buffer. Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func (s *Int8) Append(src Signed) {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	offset := s.Len()
	if s.Cap() < s.Len()+src.Len() {
		s.buffer = append(s.buffer, make([]int8, src.Len())...)
	} else {
		s.buffer = s.buffer[:s.Len()+src.Len()]
	}
	for i := 0; i < src.Len(); i++ {
		s.SetSample(i+offset, src.Sample(i))
	}
	alignCapacity(&s.buffer, s.Channels(), s.Cap())
}

// Slice slices buffer with respect to channels.
func (s *Int8) Slice(start, end int) Signed {
	start = s.BufferIndex(0, start)
	end = s.BufferIndex(0, end)
	return &Int8{
		channels: s.channels,
		buffer: s.buffer[start:end],
		bitDepth: s.bitDepth,
	}
}

// ReadInt8 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadInt8(src Signed, dst []int8) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = int8(BitDepth8.SignedValue(src.Sample(i)))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedInt8 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedInt8(src Signed, dst [][]int8) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = int8(BitDepth8.SignedValue(src.Sample(src.BufferIndex(c, i))))
		}
	}
	return
}

// WriteInt8 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteInt8(src []int8, dst Signed) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, int64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedInt8 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedInt8(src [][]int8, dst Signed) (written int) {
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
				dst.SetSample(dst.BufferIndex(c, i), int64(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}
