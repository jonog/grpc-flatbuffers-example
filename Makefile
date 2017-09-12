

generate_fbs:
	flatc --go --grpc bookmarks.fbs

generate_proto:
	protoc -I bookmarkspb/ bookmarkspb/bookmarks.proto --go_out=plugins=grpc:bookmarkspb

clean:
	rm -rf bookmarks

compile: compile_bookmarks_client compile_bookmarks_server

compile_bookmarks_client:
	cd bookmarks-client && go build -o ../client && cd ..

compile_bookmarks_server:
	cd bookmarks-server && go build -o ../server && cd ..

.PHONY: generate_fbs generate_proto clean run_bookmarks_client run_bookmarks_server