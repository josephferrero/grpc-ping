package client

import (
	"context"
	"grpc-ping/grpc/ping"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var serverAddress = "localhost:50051"

func Ping() {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	c := ping.NewPingServiceClient(conn)
	clientId := uuid.New().String()
	for {
		time.Sleep(1 * time.Second)
		sentOn := timestamppb.Now()
		r, err := c.Ping(context.TODO(), &ping.PingRequest{Message: "Ping", ClientId: clientId})
		if err != nil {
			log.Printf("could not ping: %v", err)
			continue
		}
		latency := sentOn.AsTime().Sub(r.RecievedOn.AsTime())
		log.Print(r.GetMessage(), " ", r.GetIndex(), latency)
	}

}

func GetStream() {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	c := ping.NewPingServiceClient(conn)

	stream, err := c.Stream(context.TODO(), &ping.StreamRequest{Message: "give me a stream"})
	if err != nil {
		log.Fatalf("failed to get stream: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Print("end of stream")
			break
		} else if err == nil {
			log.Print(resp.GetMessage())
		}

		if err != nil {
			log.Fatal(err)
		}

	}
}
