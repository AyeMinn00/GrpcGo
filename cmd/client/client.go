package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-faker/faker/v4"
	proto "github.com/gosample/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"time"
)

var (
	addr = flag.String("addr", "localhost:5000", "the address connect to")
)

func init() {
	rand.Seed(time.Now().UnixNano()) // Initialize the random number generator
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("cannot close connection %s ", err)
		}
	}(conn)
	//c := proto.NewGreeterClient(conn)
	client := proto.NewChatClient(conn)

	// Contact the server and print out the response
	//ctx, cancel := context.With(context.Background())
	//defer cancel()

	stream, err := client.SendMessage(context.Background())
	if err != nil {
		log.Fatalf("joined chat failed, error %s ", err)
	}
	// Send chat messages
	for i := 1; i <= 1000; i++ {
		msg := &proto.ChatMessage{
			UserId: "User1",
			Text:   faker.Sentence(),
		}
		if err := stream.Send(msg); err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Receive chat messages
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}
		fmt.Printf("%s (%s): %s\n", msg.UserId, msg.Text)
	}

}
