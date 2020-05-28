module jonog/grpc-flatbuffers-example/bookmarks-client

go 1.14

require (
	github.com/google/flatbuffers v1.12.0
	github.com/jonog/grpc-flatbuffers-example/bookmarks v0.0.0
	google.golang.org/grpc v1.29.1
)

replace github.com/jonog/grpc-flatbuffers-example/bookmarks => ../bookmarks
