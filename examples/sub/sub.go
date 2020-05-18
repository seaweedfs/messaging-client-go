package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/chrislusf/seaweedfs/weed/messaging/msgclient"
	"github.com/chrislusf/seaweedfs/weed/pb/messaging_pb"
)

var (
	subSubscriberId = flag.String("subscriberId", "sub1", "a unique subscriber id for persisting progress")
	subBroker       = flag.String("broker", "localhost:17777", "comma-separated broker list in hostname:port")
	subNamespace    = flag.String("ns", "ns1", "namespace")
	subTopic        = flag.String("topic", "topic1", "topic name")
	subStart        = flag.Duration("timeAgo", 0, "start time before now. \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\"")
)

func main() {
	flag.Parse()

	mc := msgclient.NewMessagingClient(strings.Split(*subBroker, ",")...)
	sub, err := mc.NewSubscriber(*subSubscriberId, *subNamespace, *subTopic, -1, time.Now().Add(-*subStart))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	sub.Subscribe(func(m *messaging_pb.Message) {
		fmt.Printf("> %s\n", string(m.Value))
	})

}
