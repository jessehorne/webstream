package main

import (
	"fmt"
	"github.com/jessehorne/webstream"
	"net/http"
	"time"
)

func streamHandler(rw http.ResponseWriter, req *http.Request) {
	ws := webstream.NewStream()

	// If you want to set a connection timeout, try the following.
	// This will force a connection to close 2 seconds after it is created.
	// ws.SetTimeout(2 * time.Second)

	go func() {
		for i := 1; i <= 3; i++ {
			ws.WriteString(fmt.Sprintf("Count: %d\n", i))
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
