WebStream
===

A lightweight library for handling http2 streams.

# Usage Example

## Server

The following script uses the WebStream library to listen on port `5656` at the endpoint `/counter`.

```shell
go run examples/server_counter.go
```

To connect and view the servers stream of data, use the following command.

```shell
curl --http2 -N localhost:5656/counter
```

It should count to 3 and close the connection.

## more coming soon!