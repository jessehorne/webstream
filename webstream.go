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
	ResponseWriter http.ResponseWriter
	Flusher        http.Flusher
	Request        *http.Request
	Ready          bool
	Closer         chan bool
	Closed         bool
	Timeout        time.Duration
	TimeoutSet     bool
}

// Handler is the HTTP handler function that sets up a WebStream to work over an HTTP connection.
// This function will wait indefinitely until a timeout is reached or until the WebStream is
// closed.
func (ws *WebStream) Handler(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	ws.ResponseWriter = rw
	ws.Flusher = flusher
	ws.Request = req
	ws.Ready = true

	// Monitor if a timeout is ever set and close the connection if we ever reach that time
	go func(ws *WebStream) {
		for ws.Ready {
			if ws.TimeoutSet {
				select {
				case <-time.After(ws.Timeout):
					ws.Closer <- true
				}
			}
		}
	}(ws)

	select {
	case <-ws.Closer:
		ws.Closed = true
		return
	}
}

// WriteString will simply convert the string to an array of bytes and send it over the ResponseWriter
// and then immediately flush the message queue. Messages are sent asynchronously.
func (ws *WebStream) WriteString(data string) {
	// make sure we haven't set this webstream as closed before
	if ws.Closed {
		return
	}

	// if the webstream isn't ready then we hold here indefinitely
	for !ws.Ready {
	}

	// due to the nature of async, if a webstream loses scope, this becomes nil
	// so attempting to access its members is a no-no
	if ws != nil {
		// TODO: What if ws becomes nil at any point in here?
		ws.ResponseWriter.Write([]byte(data))
		ws.Flusher.Flush()
	}
}

// SetTimeout sets a timeout for the request connection. By default, there is no timeout.
// The entire connection will drop if the timeout is reached. For example, if you set the timeout
// to (1 * time.Second), then a connection will only last for one second. It should be noted that
// if a timeout is set later on in the connection lifespan, it will close after that timeout has
// been reached from the time that you call SetTimeout.
func (ws *WebStream) SetTimeout(d time.Duration) {
	ws.Timeout = d
	ws.TimeoutSet = true
}

// Close closes a webstream and its HTTP connection. Any pending Write's will not go through after this
// is called.
func (ws *WebStream) Close() {
	ws.Closer <- true
}

func NewStream() *WebStream {
	return &WebStream{
		ResponseWriter: nil,
		Flusher:        nil,
		Request:        nil,
		Ready:          false,
		Closer:         make(chan bool),
		Timeout:        0,
		TimeoutSet:     false,
	}
}
