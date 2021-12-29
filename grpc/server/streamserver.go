package main

import (
	"errors"
	"io"
	"log"
	"net"

	pb "github.com/labs/grpc/proto"
	"google.golang.org/grpc"
)

const (
	port = "5100"
)

type StreamServer struct {
}

func main() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamServer{})

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	if err := server.Serve(l); err != nil {
		log.Fatalf("server finish on err: %v", err)
	}
}

func (s *StreamServer) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(i),
			},
		})
		if err != nil {
			log.Printf("stream.Send err: %v\n", err)
		}
	}
	return nil
}

func (s *StreamServer) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		log.Printf("record recv request: %v", r)
	}
	err := stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{
		Name:  "StreamServer Record",
		Value: 1,
	}})
	return err
}

func (s StreamServer) Route(stream pb.StreamService_RouteServer) error {
	r, err := stream.Recv()
	if err != nil {
		return err
	}
	err = stream.Send(&pb.StreamResponse{Pt: &pb.StreamPoint{
		Name:  "StreamServer Route",
		Value: r.Pt.Value + 1,
	}})
	return err
}
