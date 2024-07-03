package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jessehorne/webstream"
)

func streamHandler(rw http.ResponseWriter, req *http.Request) {
	ws, _ := webstream.NewStream(rw, req)

	go func() {
		for i := 1; i <= 3; i++ {
			ws.Write([]byte(fmt.Sprintf("Count: %d\n", i)))
			time.Sleep(1 * time.Second)
		}
		ws.Close()
	}()

	ws.Wait()
}

func main() {
	http.HandleFunc("/counter", streamHandler)
	if err := http.ListenAndServe(":5656", nil); err != nil {
		panic(err)
	}
}
