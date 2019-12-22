

generate_fbs:
	flatc --go --grpc fileupload.fbs

generate_proto:
	mkdir fileuploadpb
	protoc fileupload.proto --go_out=plugins=grpc:fileuploadpb

all: clean generate_proto generate_fbs compile_fileupload_client compile_fileupload_server

compile_fileupload_client:
	cd fileupload-client && go build -o ../client  && cd ..

compile_fileupload_server:
	cd fileupload-server && go build -o ../server  && cd ..

clean:
	rm -rf server client fileupload fileuploadpb

.PHONY: clean generate_fbs generate_proto compile compile_fileupload_client compile_fileupload_server
