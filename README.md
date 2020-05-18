# seaweedfs-messaging

SeaweedFS Messaging offers a simple message queue, which has an unlimited capacity and repeatable reads.

```
# To start filer
weed server -filer
# Also start a message broker. Later it will also be included in the weed server command, so just one command shall be enough.
weed msg.broker
```

The message queue in SeaweedFS is conceptually a remote FIFO file. It can be rewinded at any time, addressed by nano-second timestamp.

There are 2 kinds of examples in this repo: 
* Network channel for Go.
* Pub/Sub to a distributed persisted message queue.

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

	mc := msgclient.NewMessagingClient("localhost:17777")

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

### Network Channel to read
```
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

}

```


## Pub/Sub to a distributed persisted message queue

The pub/sub example is also simple.
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
