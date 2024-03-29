package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/PerminovEugene/configs-and-learning/go-course/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client running")

	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Could not connect %w", err)
	}
	defer connection.Close()

	c := greetpb.NewGreetServiceClient(connection)
	fmt.Printf("Created client: %f\n", c)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Eugene",
			LastName: "Perminov",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Misha",
			LastName: "Medved",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Responce from greet many times %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Antoha",
				LastName: "Margnas",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Potap",
				LastName: "Margnas",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Nastya",
				LastName: "Margnas",
			},
		},
	}
	
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreed %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req %v\n\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Microsecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while close and receive %v", err)
	}
	fmt.Printf("Long greet response %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	// create stream

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatal("Error while greet everyone")
		return
	}
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tolik",
				LastName: "Rikardo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Vita",
				LastName: "Rikardo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Kia",
				LastName: "Rikardo",
			},
		},
	}
	waitChannel := make(chan struct{})
	// send messages
	go func() {
		for _, req := range requests {
			fmt.Println("Sending", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitChannel)
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving %v", err)
				break
			}
			fmt.Printf("Received %v\n", res.GetResult())
		}
	}()
	<-waitChannel
	// receive messages

	// block until everything is done
}
