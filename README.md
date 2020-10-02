# Performance Comparison Between Protobuf over gRPC and JSON over HTTP2

## Description
For each type of communication, there is a server and a client.
Communication between each server-client pair for a given communication type is encrypted using SSL.
Therefore:
* In both cases communication is over HTTP/2.
* In both cases communication is encrypted using SSL.
* The difference is that in one case we send a JSON over HTTP/2 directly and in the other case we send a Protobuf message over gRPC.

The general flow is as follows:
* A user starts the server for a particular communication type (HTTP/gRPC).
* The user runs the client while passing as command line arguments the number of iterations to run for and the number of requests to send to the server in each iteration.
* In each iteration, in each request, the user sends the constant message `{firstName: "Bill", lastName: "Clinton", Age: 74}` to the server (while it's written here as JSON, but in Protobuf the same message is sent in the corresponding format).
* All requests are sent sequentially (not concurrently).
* At the end of the client's execution, the total amount of time in milliseconds it took to receive a response from the server is summed up, divided by the number of iterations and the resulting average time in milliseconds is printed back to the user.

## Installation Instructions
1. Create a TLS certificate and key pair. There are instructions on how to do that online. 
1. Compile each file like this:
    1. `go build -o http-server http-server/server.go`
    1. `go build -o http-client http-client/client.go`
    1. `go build -o grpc-server grpc-server/server.go`
    1. `go build -o grpc-client grpc-client/client.go`
    
## Execution Instructions
1. For a particular communication type, first run the server using the command (for HTTP for example): `./http-server --certificate [path to certificate] --key [path to private key]`
1. Run the corresponding client using the command (for HTTP for example): `./http-client --num-requests [number of requests per iteration] --num-iterations [number of iterations]`
