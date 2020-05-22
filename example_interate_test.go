package signal_test

import (
	"fmt"

	"pipelined.dev/signal"
)

// This example demonstrates how to iterate over the buffer.
func Example_iterate() {
	// allocate int64 buffer with 2 channels and capacity of 8 samples per channel
	buf := signal.Allocator{Channels: 2, Capacity: 8}.Int64(signal.BitDepth64)

	// write striped data
	buf = signal.WriteStripedInt8([][]int8{{1, 1, 1, 1}, {2, 2, 2, 2}}, buf)

	// iterate over buffer interleaved data
	for pos := 0; pos < buf.Len(); pos++ {
		fmt.Printf("%d", buf.Sample(pos))
	}

	for channel := 0; channel < buf.Channels(); channel++ {
		fmt.Println()
		for pos := 0; pos < buf.Length(); pos++ {
			fmt.Printf("%d", buf.Sample(buf.ChannelPos(channel, pos)))
		}
	}

	// Output:
	// 12121212
	// 1111
	// 2222
}
