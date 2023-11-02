package main

import (
	"fmt"
	proto "github.com/gosample/pb"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"log"
	"net"
)

const port = 5000

func main() {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterRouteGuideServer(s, &sampleRouteGuideServer{})
	proto.RegisterChatServer(s, &MyChatServer{clients: make(map[proto.Chat_SendMessageServer]bool)})
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
