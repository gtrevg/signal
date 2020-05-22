// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-22 22:11:58.6906 +0200 CEST m=+0.016077027
package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestInt64(t *testing.T) {
	t.Run("int64", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Int64(signal.BitDepth64).
			Append(signal.WriteStripedInt64(
				[][]int64{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.Int64(signal.BitDepth64)),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]int64{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}