package main

import (
	"github.com/go-faker/faker/v4"
	proto "github.com/gosample/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

type MyChatServer struct {
	mu sync.RWMutex
	// Store connected clients
	clients map[proto.Chat_SendMessageServer]bool
	proto.UnimplementedChatServer
}

func (s *MyChatServer) SendMessage(stream proto.Chat_SendMessageServer) error {
	s.mu.Lock()
	s.clients[stream] = true
	s.mu.Unlock()

	// Create a goroutine to monitor the client's context
	go func() {
		<-stream.Context().Done()
		// The client has disconnected
		log.Printf("Client disconnected")
		s.mu.Lock()
		delete(s.clients, stream)
		s.mu.Unlock()
	}()

	// Handle coming messages from client
	for {
		msg, err := stream.Recv()
		if err != nil {
			// Handle errors (e.g., client disconnect, stream error)
			// Remove the client from the list if needed
			s.mu.Lock()
			delete(s.clients, stream)
			s.mu.Unlock()

			// If the error indicates a client disconnect, you can return specific codes or messages
			if errStatus, ok := status.FromError(err); ok && errStatus.Code() == codes.Canceled {
				return nil // Client canceled the stream
			}

			// Handle other errors as needed
			log.Printf("Error receiving message: %v", err)
			return err
		}
		log.Printf("%s send %s", msg.UserId, msg.Text)
		// Process the received message and broadcast to other clients
		s.mu.Lock()
		for client := range s.clients {
			responseMsg := &proto.ResponseMessage{
				MsgId:  faker.UUIDHyphenated(),
				UserId: msg.UserId,
				Text:   msg.Text,
			}
			if err := client.Send(responseMsg); err != nil {
				// Handle send error
				log.Printf("Error sending message: %v", err)
			}
		}
		s.mu.Unlock()
	}
}
