package main

import (
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/chrislusf/seaweedfs/weed/messaging/msgclient"
	"github.com/chrislusf/seaweedfs/weed/pb/messaging_pb"
	"github.com/chrislusf/seaweedfs/weed/util/grace"
)

var (
	topic      = flag.String("topic", "topic_load", "topic name")
	namespace  = flag.String("namespace", "ns1", "topic namespace")
	subscriber = flag.String("subscriber", "topic_load", "subscriber identification")
	partitionId = flag.Int("partitionId", -1, "partition id")
	subStart     = flag.Duration("timeAgo", 0, "start time before now. \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\"")
)

func main() {

	flag.Parse()

	mc := msgclient.NewMessagingClient("localhost:17777")

	sub, err := mc.NewSubscriber("subscriber1", *namespace, *topic, *partitionId, time.Now().Add(-*subStart))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	var finalCount, finalSize int64
	var isStopping bool
	grace.OnInterrupt(func() {
		isStopping = true
		fmt.Printf("message count: %d\n", finalCount)
		fmt.Printf("message total : %d byte\n", finalSize)
		sub.Shutdown()
	})

	var totalCount, totalSize int64

	go func() {
		var startTime = time.Now()
		for !isStopping {
			time.Sleep(time.Second)
			t := time.Now()
			elapsed := t.Sub(startTime).Seconds()
			fmt.Printf("message throuput: %.2f/s %.2f MB/s\n",
				float64(totalCount)/elapsed,
				float64(totalSize)/1024.0/1024.0/elapsed)

			startTime = t
			totalCount = 0
			totalSize = 0
		}
	}()

	sub.Subscribe(func(m *messaging_pb.Message) {
		atomic.AddInt64(&totalCount, 1)
		atomic.AddInt64(&totalSize, int64(len(m.Value)))
		atomic.AddInt64(&finalCount, 1)
		atomic.AddInt64(&finalSize, int64(len(m.Value)))
	})

}
