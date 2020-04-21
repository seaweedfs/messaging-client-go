package main

import (
	"fmt"

	"github.com/chrislusf/seaweedfs/weed/messaging/client"
	"github.com/chrislusf/seaweedfs/weed/pb/messaging_pb"
)

func main() {
	mc, err := client.NewMessagingClient([]string{"localhost:9777"})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	sub, err := mc.NewSubscriber("subscriber1", "ns1", "topic1")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	sub.Subscribe(func(m *messaging_pb.Message) {
		fmt.Printf("> %s\n", string(m.Value))
	})

	select {}
}
