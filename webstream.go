package webstream

import (
	"net/http"
	"time"
)

// WebStream represents an HTTP connection that is long-lived where data can be sent back and forth
// between the server and the client as long as the connection is open. It is important to note that
// WebStream does not use websockets and instead uses a standard HTTP connection where the data
// buffer is flushed after each message sent from the server and data is read asynchronously from
// the client.
type WebStream struct {
	responseWriter http.ResponseWriter
	flusher        http.Flusher
	request        *http.Request
	ready          bool
	closer         chan bool
	closed         bool
	timeout        time.Duration
	timeoutSet     bool
}

// Handler is the HTTP handler function that sets up a WebStream to work over an HTTP connection.
// This function will wait until a timeout is reached or until the WebStream is closed.
func (ws *WebStream) Handler(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	ws.responseWriter = rw
	ws.flusher = flusher
	ws.request = req
	ws.ready = true

	// Monitor if a timeout is ever set and close the connection if we ever reach that time
	go func(ws *WebStream) {
		for ws.ready {
			if ws.timeoutSet {
				select {
				case <-time.After(ws.timeout):
					ws.closer <- true
				}
			}
		}
	}(ws)

	select {
	case <-ws.closer:
		ws.closed = true
		return
	}
}

// Write sends data immediately by calling the connections responseWriter.Write method and immediately flushing.
func (ws *WebStream) Write(data []byte) {
	// make sure we haven't set this webstream as closed before
	if ws.closed {
		return
	}

	// wait until the webstream is ready
	for !ws.ready {
	}

	// due to the nature of async, if a webstream loses scope, this becomes nil
	// so attempting to access its members is a no-no
	if ws != nil {
		// TODO: What if ws becomes nil at any point in here?
		ws.responseWriter.Write(data)
		ws.flusher.Flush()
	}
}

// SetTimeout sets a timeout for the request connection. By default, there is no timeout.
// The entire connection will drop if the timeout is reached. For example, if you set the timeout
// to (1 * time.Second), then a connection will only last for one second. It should be noted that
// if a timeout is set later on in the connection lifespan, it will close after that timeout has
// been reached from the time that you call SetTimeout.
func (ws *WebStream) SetTimeout(d time.Duration) {
	ws.timeout = d
	ws.timeoutSet = true
}

// Close closes a webstream and its HTTP connection. Any pending writes will not go through after this
// is called.
func (ws *WebStream) Close() {
	ws.closer <- true
}

func NewStream() *WebStream {
	return &WebStream{
		responseWriter: nil,
		flusher:        nil,
		request:        nil,
		ready:          false,
		closer:         make(chan bool),
		timeout:        0,
		timeoutSet:     false,
	}
}
