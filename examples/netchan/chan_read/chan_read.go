package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chrislusf/seaweedfs/weed/messaging/msgclient"
)

var (
	subSubscriberId = flag.String("subscriberId", "sub1", "a unique subscriber id for persisting progress")
)

func main() {

	flag.Parse()

	mc := msgclient.NewMessagingClient("localhost:17777")

	subChan, err := mc.NewSubChannel(*subSubscriberId, "some_chan")
	if err != nil {
		log.Fatalf("fail to create publish channel: %v\n", err)
	}

	for m := range subChan.Channel() {
		fmt.Println(string(m))
	}

	fmt.Printf("sent md5 %X\n", subChan.Md5())

}
