package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-08-08 11:27:13.137476 +0200 CEST m=+0.018111948

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUint16(t *testing.T) {
	t.Run("uint16", func() func(t *testing.T) {
		input := signal.Allocator{
			Channels: 3,
			Capacity: 3,
			Length:   3,
		}.Uint16(signal.BitDepth16)
		signal.WriteStripedUint16(
			[][]uint16{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		return testOk(
			signal.Allocator{
				Channels: 3,
				Capacity: 2,
			}.Uint16(signal.BitDepth16).Append(input.Slice(1, 3)),
			expected{
				length:   2,
				capacity: 2,
				data: [][]uint16{
					{0, 0},
					{2, 3},
					{12, 13},
				},
			},
		)
	}())
}
