package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/labs/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = "5100"

	caRoot = "../ca/"
)

type StreamServer struct {
}

func main() {
	cert, err := tls.LoadX509KeyPair(caRoot+"server.crt", caRoot+"server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caRoot + "ca.crt")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("add cert failed")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	server := grpc.NewServer(grpc.Creds(creds))
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
