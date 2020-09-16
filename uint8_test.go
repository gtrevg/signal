package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-09-15 15:22:49.765111 +0200 CEST m=+0.009556873

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUint8(t *testing.T) {
	t.Run("uint8", func() func(t *testing.T) {
		input := signal.Allocator{
			Channels: 3,
			Capacity: 3,
			Length:   3,
		}.Uint8(signal.BitDepth8)
		signal.WriteStripedUint8(
			[][]uint8{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		result := signal.Allocator{
				Channels: 3,
				Capacity: 2,
			}.Uint8(signal.BitDepth8)
		result.Append(input.Slice(1, 3))
		return testOk(
			result,
			expected{
				length:   2,
				capacity: 2,
				data: [][]uint8{
					{0, 0},
					{2, 3},
					{12, 13},
				},
			},
		)
	}())
}
