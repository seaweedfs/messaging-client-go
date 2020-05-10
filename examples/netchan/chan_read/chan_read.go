package main

import (
	"fmt"
	"log"

	"github.com/chrislusf/seaweedfs/weed/messaging/msgclient"
)

func main()  {

	mc := msgclient.NewMessagingClient("localhost:17777")

	subChan, err := mc.NewSubChannel("some_chan")
	if err != nil {
		log.Fatalf("fail to create publish channel: %v\n", err)
	}

	for m := range subChan.Channel() {
		fmt.Println(string(m))
	}

	fmt.Printf("sent md5 %X\n", subChan.Md5())

}
