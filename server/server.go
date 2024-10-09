package server

import (
	"context"
	"grpc-ping/grpc/ping"
	"net"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pingIndex map[string]int32
	server    *grpc.Server
	ping.UnimplementedPingServiceServer
}

func New() *Server {
	return &Server{
		pingIndex: make(map[string]int32),
		server:    grpc.NewServer(),
	}
}

func (s *Server) Ping(ctx context.Context, in *ping.PingRequest) (*ping.PingResponse, error) {
	recievedOn := timestamppb.Now()
	log.Printf("Recieved: %v", in.Message)
	clientID := in.GetClientId()
	var index int32 = 0
	if currIndex, ok := s.pingIndex[clientID]; ok {
		index = currIndex + 1
		s.pingIndex[clientID] = index
	} else {
		s.pingIndex[clientID] = index
	}
	res := &ping.PingResponse{
		Index:      index,
		Message:    "Pong",
		RecievedOn: recievedOn,
	}
	return res, nil
}

func (s *Server) Stream(in *ping.StreamRequest, srv ping.PingService_StreamServer) error {
	for i := 0; i < 10; i++ {
		srv.Send(&ping.StreamResponse{Message: "Hello"})
	}
	srv.Context().Done()
	return nil
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ping.RegisterPingServiceServer(s.server, s)
	log.Printf("listening on: %v", lis.Addr())
	if err := s.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
