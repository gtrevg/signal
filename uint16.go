package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-09-12 15:44:22.504697 +0200 CEST m=+0.018587604

import "math"

// Uint16 is uint16 Unsigned fixed-point signal.
type Uint16 struct {
	buffer []uint16
	channels
	bitDepth
}

// Uint16 allocates a new sequential uint16 signal buffer.
func (a Allocator) Uint16(bd BitDepth) Unsigned {
	return Uint16{
		buffer:   make([]uint16, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth16),
	}
}

func (s Uint16) setBitDepth(bd BitDepth) Unsigned {
	s.bitDepth = limitBitDepth(bd, BitDepth16)
	return s
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Uint16) AppendSample(value uint64) Unsigned {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, uint16(s.BitDepth().UnsignedValue(value)))
	return s
}

// SetSample sets sample value for provided index.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Uint16) SetSample(i int, value uint64) {
	s.buffer[i] = uint16(s.BitDepth().UnsignedValue(value))
}

// GetUint16 selects a new sequential uint16 signal buffer.
// from the pool.
func (p PoolAllocator) GetUint16(bd BitDepth) Unsigned {
	return Uint16{
		buffer:   p.u16.Get().([]uint16)[:p.Channels*p.Length],
		channels: channels(p.Channels),
		bitDepth: limitBitDepth(bd, BitDepth16),
	}
}

// PutUint16 places signal buffer back to the pool. If a type of
// provided buffer isn't Uint16 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p PoolAllocator) PutUint16(s Unsigned) {
	mustSameCapacity(s.Cap(), p.Channels*p.Capacity)
	if sig, ok := s.(Uint16); ok {
		buf := sig.buffer[:sig.Cap()]
		for i := range buf {
			buf[i] = 0
		}
		p.u16.Put(sig.buffer)
	} else {
		panic("pool put uint16 invalid type")
	}
}

// Capacity returns capacity of a single channel.
func (s Uint16) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Uint16) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s Uint16) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Uint16) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided channel and index.
func (s Uint16) Sample(i int) uint64 {
	return uint64(s.buffer[i])
}

// Append appends [0:Length] samples from src to current buffer and returns
// new Unsigned buffer. Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func (s Uint16) Append(src Unsigned) Unsigned {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	offset := s.Len()
	if s.Cap() < s.Len()+src.Len() {
		s.buffer = append(s.buffer, make([]uint16, src.Len())...)
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
func (s Uint16) Slice(start, end int) Unsigned {
	start = s.BufferIndex(0, start)
	end = s.BufferIndex(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// ReadUint16 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadUint16(src Unsigned, dst []uint16) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = uint16(BitDepth16.UnsignedValue(src.Sample(i)))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedUint16 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedUint16(src Unsigned, dst [][]uint16) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = uint16(BitDepth16.UnsignedValue(src.Sample(src.BufferIndex(c, i))))
		}
	}
	return
}

// WriteUint16 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteUint16(src []uint16, dst Unsigned) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, uint64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedUint16 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedUint16(src [][]uint16, dst Unsigned) (written int) {
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
				dst.SetSample(dst.BufferIndex(c, i), uint64(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}
