# seaweedfs-messaging

SeaweedFS Messaging offers a simple message queue, which has an unlimited capacity and repeatable reads.

```
# To start message broker
weed server -msgBroker
```

The message queue in SeaweedFS is conceptually a remote FIFO file. It can be rewinded at any time, addressed by nano-second timestamp.

the message broker uses gRPC API to stream read and write messages. More clients can be added following the [SeaweedFS messaging gRPC protobuf definition](https://github.com/chrislusf/seaweedfs/blob/master/weed/pb/messaging.proto). 

There are 2 kinds of examples in this repo: 
* Network channel for Go.
* Pub/Sub to a distributed persisted message queue.

Current examples are in Go. But more examples can be easily added by the gRPC API.

## Network Channel

Originally Go has a netchan package, but it is found hard to implement. Here is just one way to implement it.

For Go, the best tutorial is the source code.

### Network Channel to write
```
package main

import (
	"log"
	"strings"

	"github.com/chrislusf/seaweedfs/weed/messaging/msgclient"
)

func main()  {

	// connect to message broker via gRPC
	mc := msgclient.NewMessagingClient("localhost:17777")

	// write to this channel.
	pubChan, err := mc.NewPubChannel("some_chan")
	if err != nil {
		log.Fatalf("fail to create publish channel: %v\n", err)
	}

	for _, t := range strings.Split(text, "\n") {
		pubChan.Publish([]byte(t))
	}
	pubChan.Close()

}
var text = "..."

```

Note: If the channel is closed, new writes will be rejected.

### Network Channel to read



```
package main

import (
	"fmt"
	"log"

	"github.com/chrislusf/seaweedfs/weed/messaging/msgclient"
)

func main() {

	mc := msgclient.NewMessagingClient("localhost:17777")

	// connect to channel
	// the channel to read does not need to exist beforehand.
	subChan, err := mc.NewSubChannel("subscriber1", "some_chan")
	if err != nil {
		log.Fatalf("fail to create publish channel: %v\n", err)
	}

	// loop the data until the channel is closed by the publishing program
	for m := range subChan.Channel() {
		fmt.Println(string(m))
	}

}

```

Note: The data in the channel can be read multiple times.

After consuming the data, since the messages in the channel is a persisted, you need to delete it explicitly:

```
  mc.DeleteChannel("some_chan")

```

## Pub/Sub to a distributed persisted message queue

The pub/sub example is also simple. The difference from channel is that message queue can not be closed.

### Publish to a distributed persisted message queue
```
	mc := msgclient.NewMessagingClient("localhost:17777")
	pub, err := mc.NewPublisher("publisher1", "ns1", "topic1")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	for _, t := range strings.Split(text, "\n") {
		if err = pub.Publish(&messaging_pb.Message{
			Key:     nil,
			Value:   []byte(t),
			Headers: nil,
		}); err != nil {
			println("err:", err.Error())
		} else {
			println(t)
		}
	}

```

### Subscribe to a distributed persisted message queue
```
	mc := msgclient.NewMessagingClient("localhost:17777")
	sub, err := mc.NewSubscriber("subscriber1", "ns1", "topic1", -1, time.Now())
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	sub.Subscribe(func(m *messaging_pb.Message) {
		fmt.Printf("> %s\n", string(m.Value))
	})

```

The read can be rewinded to any timestamp.

After consuming the data, since the messages in the message queue is a persisted, you need to delete it explicitly:

```
  mc.DeleteTopic("ns1", "topic1")

```
