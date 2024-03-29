package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/PerminovEugene/configs-and-learning/go-course/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func(*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func(*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) (error){
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " iteration " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func(*server) LongGreet(stream greetpb.GreetService_LongGreetServer) (error) {
	result := ""
	fmt.Println("Received long greed")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream API %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
		fmt.Println("Updated result")
	}
}

func(*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatal("Error while reading client stream")
			return err;
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatal("Error while sending data to client")
			return err;
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed listen %v", err)
	}
	s := grpc.NewServer()
	
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}