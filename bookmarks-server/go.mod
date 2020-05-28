module jonog/grpc-flatbuffers-example/bookmarks-server

go 1.14

require (
	github.com/google/flatbuffers v1.12.0
	github.com/jonog/grpc-flatbuffers-example v0.0.0-20170912123016-8f2cb1fcf3d7
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
	google.golang.org/grpc v1.29.1
)

replace github.com/jonog/grpc-flatbuffers-example/bookmarks => ../bookmarks
