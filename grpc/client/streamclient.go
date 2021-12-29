package main

import (
	"context"
	"errors"
	"io"
	"log"

	pb "github.com/labs/grpc/proto"
	"google.golang.org/grpc"
)

const (
	port = "5100"
)

func main() {
	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dail failed: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	err = printList(client, &pb.StreamRequest{Pt: &pb.StreamPoint{
		Name:  "gRPC Stream Client: List",
		Value: 2021,
	}})
	if err != nil {
		log.Fatalf("printList err: %v", err)
	}

	err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{
		Name:  "gRPC Stream Client: Record",
		Value: 2021,
	}})
	if err != nil {
		log.Fatalf("printRecord err: %v", err)
	}

	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{
		Name:  "gRPC Stream Client: Route",
		Value: 2021,
	}})
	if err != nil {
		log.Fatalf("printRoute err: %v", err)
	}
}

func printList(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			} else {
				return nil
			}
		}
		log.Printf("list recv response: %v\n", res)
	}
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		if err := stream.Send(r); err != nil {
			return err
		}
		r.Pt.Value += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("record recv response: %v", res)
	return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}
	if err := stream.Send(r); err != nil {
		return err
	}
	res, err := stream.Recv()
	if err != nil {
		return err
	}
	log.Printf("route recv response: %v", res)
	return nil
}
