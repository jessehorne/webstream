WebStream
===

A lightweight library for handling http2 streams.

# Usage Example

## Server

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jessehorne/webstream"
)

func streamHandler(rw http.ResponseWriter, req *http.Request) {
	ws := webstream.NewStream()
	go func() {
		for i := 1; i <= 3; i++ {
			ws.Write([]byte(fmt.Sprintf("Count: %d\n", i)))
			time.Sleep(1 * time.Second)
		}
		ws.Close()
	}()
	ws.Handler(rw, req)
}

func main() {
	http.HandleFunc("/counter", streamHandler)
	if err := http.ListenAndServe(":5656", nil); err != nil {
		panic(err)
	}
}

```

To connect and view the servers stream of data, use the following command.

```shell
curl --http2 -N localhost:5656/counter
```

It should count to 3 and close the connection.

## more coming soon!