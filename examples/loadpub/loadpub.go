package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/chrislusf/seaweedfs/weed/messaging/client"
	"github.com/chrislusf/seaweedfs/weed/pb/messaging_pb"
)

var (
	messageCount = flag.Int("n", 1000000, "message count")
	messageSize = flag.Int("size", 1024, "message size")
	topic = flag.String("topic", "topic_load", "topic name")
	namespace = flag.String("namespace", "ns1", "topic namespace")
)

func main() {

	flag.Parse()

	mc, err := client.NewMessagingClient([]string{"localhost:9777"})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	pub, err := mc.NewPublisher(*namespace, *topic)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	startTime := time.Now()
	var totalCount, totalSize int64
	var buf = make([]byte, *messageSize)
	for i := 0; i < *messageCount; i++ {
		rand.Read(buf)
		pub.Publish(&messaging_pb.Message{
			Value: buf,
		})
		totalCount++
		totalSize+=int64(*messageSize)
	}

	fmt.Printf("message count: %d\n", totalCount)
	fmt.Printf("message size : %d byte\n", *messageSize)
	fmt.Printf("message total : %d byte\n", totalSize)
	fmt.Printf("message throuput: %.2f MB/s\n", float64(totalSize) / 1024.0 / 1024.0 / time.Now().Sub(startTime).Seconds())

	time.Sleep(time.Second)

}
