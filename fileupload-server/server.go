package main

import (
	fb "../fileupload"
	pb "../fileuploadpb"
	"bytes"
	"crypto/sha256"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	_ "io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

type fbServer struct {
	out   bytes.Buffer
	nPkt  int
	size  int
	fname string
	sum   [32]byte
}
type pbServer fbServer

var fbAddr = "0.0.0.0:50051"
var pbAddr = "0.0.0.0:50052"

func (s *fbServer) Upload(context context.Context, in *fb.UploadRequest) (*flatbuffers.Builder, error) {
	b := flatbuffers.NewBuilder(0)
	begin := time.Now()
	flag := in.Flag()

	if flag == fb.FlagFirst {
		s.size = int(in.Size())
		s.fname = string(in.Filename())
		s.out.Grow(s.size)
		//		fmt.Printf("File %s size %d started\n", s.fname, s.size)
	}
	s.nPkt++
	if flag == fb.FlagLast {
		s.out.Truncate(int(s.size))
		s.sum = sha256.Sum256(s.out.Bytes())
		fmt.Printf("file %s sha %x size %d received in %.3fs packet %d\n",
			s.fname, s.sum, s.out.Len(), time.Since(begin).Seconds(), s.nPkt)
		/*
			err := ioutil.WriteFile("hello", s.out.Bytes(), 0644)
			if err != nil {
				log.Fatal(err)
			}
		*/
		s.nPkt = 0
		s.out.Reset()
		s.size = 0
	} else {
		file := in.File()
		s.out.Write(file)
	}
	fb.UploadResponseStart(b)
	b.Finish(fb.UploadResponseEnd(b))
	return b, nil
}

func (s *pbServer) Upload(context context.Context, in *pb.UploadRequest) (*pb.UploadResponse, error) {
	begin := time.Now()
	flag := in.GetFlag()

	if flag == pb.UploadRequest_First {
		s.size = int(in.GetSize())
		s.fname = string(in.GetFilename())
		s.out.Grow(s.size)
		fmt.Printf("File %s size %d started\n", s.fname, s.size)
	}
	s.nPkt++
	if flag == pb.UploadRequest_Last {
		s.out.Truncate(int(s.size))
		s.sum = sha256.Sum256(s.out.Bytes())
		fmt.Printf("file %s sha %x size %d received in %.3fs packet %d\n",
			s.fname, s.sum, s.out.Len(), time.Since(begin).Seconds(), s.nPkt)
		/*
			err := ioutil.WriteFile("hello", s.out.Bytes(), 0644)
			if err != nil {
				log.Fatal(err)
			}
		*/
		s.nPkt = 0
		s.out.Reset()
		s.size = 0
	} else {
		file := in.GetFile()
		s.out.Write(file)
	}
	return &pb.UploadResponse{}, nil
}

func main() {
	done := make(chan bool, 1)
	go func() {
		r := http.NewServeMux()
		// Register pprof handlers
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		http.ListenAndServe(":9090", r)
	}()

	go func() {
		lis, err := net.Listen("tcp", fbAddr)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		ser := grpc.NewServer(grpc.CustomCodec(flatbuffers.FlatbuffersCodec{}))
		fb.RegisterFileUploadServer(ser, &fbServer{})
		if err := ser.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}

	}()

	go func() {
		lis, err := net.Listen("tcp", pbAddr)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		ser := grpc.NewServer()
		pb.RegisterFileUploadServer(ser, &pbServer{})
		if err := ser.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	<-done

}
