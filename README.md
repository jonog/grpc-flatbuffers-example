
# gRPC Flatbuffers Example

A simple bookmarking service defined in the FlatBuffers IDL, and creation of gRPC server interfaces and client stubs. (`fileupload.fbs`)

A Protocol Buffers IDL has also been provided for comparison. (`fileupload.proto`)

### Instructions

#### Compile the FlatBuffers IDL file
```
flatc --go --grpc fileupload.fbs
```

### Compile the Go Server & Client
```
make compile
```

#### Start Server
```
./server
```

#### Send commands via Client
```
./client <filename> 
```

Run `./server`

Run `./client`

### FlatBuffers Compiler Setup

Setup `flatc`:
* Download flatbuffers src via [Github Releases](https://github.com/google/flatbuffers/releases)
* Compile `flatc`. e.g. `cmake -G"Unix Makefiles"` then run `make`

### gRPC Setup
```
go get google.golang.org/grpc
```
