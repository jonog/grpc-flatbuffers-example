package main

import (
	fb "../fileupload"
	pb "../fileuploadpb"
	"context"
	_ "fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

var fbAddr = "0.0.0.0:50051"
var pbAddr = "0.0.0.0:50052"

func fbSendFile(fname string) {
	conn, err := grpc.Dial(fbAddr, grpc.WithInsecure(), grpc.WithCodec(flatbuffers.FlatbuffersCodec{}))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// pack flatbuffer
	b := flatbuffers.NewBuilder(0)

	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("--> failed to open file")
		return
	}
	fInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("--> failed to stat file")
		return
	}

	ctx := context.Background()
	client := fb.NewFileUploadClient(conn)

	buf := make([]byte, 16*1024)
	var start bool
	var last bool
	var nPkts int
	for {
		n, err := file.Read(buf)
		b.Reset()
		nPkts++
		if n > 0 {
			file_pos := b.CreateByteString(buf)
			name_pos := b.CreateString(fname)

			fb.UploadRequestStart(b)
			fb.UploadRequestAddFile(b, file_pos)
			if start == false {
				fb.UploadRequestAddFilename(b, name_pos)
				fb.UploadRequestAddFlag(b, fb.FlagFirst)
				fb.UploadRequestAddSize(b, fInfo.Size())
				start = true
			}
			b.Finish(fb.UploadRequestEnd(b))
			//b.Bytes[b.Head():]
		}

		if err == io.EOF {
			fb.UploadRequestStart(b)
			fb.UploadRequestAddFlag(b, fb.FlagLast)
			b.Finish(fb.UploadRequestEnd(b))
			last = true
		}
		// send over grpc
		_, err = client.Upload(ctx, b)
		if err != nil {
			log.Fatalf("Retrieve client failed: %v", err)
		}
		if last {
			log.Printf("SENT in %d pkts", nPkts)
			nPkts = 0
			break
		}
	}
}

func setSize(r *pb.UploadRequest, size uint32) {
	r.Size = size
}

func setFlag(r *pb.UploadRequest, flag pb.UploadRequest_Flag) {
	r.Flag = flag
}

func setFilename(r *pb.UploadRequest, name string) {
	r.Filename = name
}

func setFile(r *pb.UploadRequest, buf []byte) {
	r.File = buf
}

func pbSendFile(fname string) {
	conn, err := grpc.Dial(pbAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("--> failed to open file")
		return
	}

	fInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("--> failed to stat file")
		return
	}

	ctx := context.Background()
	client := pb.NewFileUploadClient(conn)

	buf := make([]byte, 16*1024)
	var start bool
	var last bool
	var nPkts int
	for {
		n, err := file.Read(buf)
		r := &pb.UploadRequest{}
		nPkts++
		if n > 0 {
			setFile(r, buf)
			if start == false {
				setFilename(r, fname)
				setFlag(r, pb.UploadRequest_First)
				setSize(r, uint32(fInfo.Size()))
				start = true
			}
			//b.Bytes[b.Head():]
		}

		if err == io.EOF {
			setFlag(r, pb.UploadRequest_Last)
			last = true
		}

		_, err = client.Upload(ctx, r)
		if err != nil {
			log.Fatalf("Fail to get the grpc stream, Error:%v", err)
		}
		if last {
			log.Printf("SENT in %d pkts", nPkts)
			nPkts = 0
			break
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("./client [-f | -p] <filename>")
		return
	}
	count := 10
	cmd := os.Args[1]
	fname := os.Args[2]
	if cmd == "-f" {
		for i := 0; i < count; i++ {
			fbSendFile(fname)
		}
	} else if cmd == "-p" {
		for i := 0; i < count; i++ {
			pbSendFile(fname)
		}
	}
}
