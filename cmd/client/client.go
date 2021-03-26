package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Brunorodrigoss/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Rodrigo",
		Email: "r@r.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Rodrigo",
		Email: "r@r.com",
	}

	reponseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := reponseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		fmt.Println("Status:", stream.Status, " - ", stream.User)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "b1",
			Name:  "Bruno 1",
			Email: "b1@example.com",
		},
		&pb.User{
			Id:    "b2",
			Name:  "Bruno 2",
			Email: "b2@example.com",
		},
		&pb.User{
			Id:    "b3",
			Name:  "Bruno 3",
			Email: "b3@example.com",
		},
		&pb.User{
			Id:    "b4",
			Name:  "Bruno 4",
			Email: "b4@example.com",
		},
		&pb.User{
			Id:    "b5",
			Name:  "Bruno 5",
			Email: "b5@example.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "b1",
			Name:  "Bruno 1",
			Email: "b1@example.com",
		},
		&pb.User{
			Id:    "b2",
			Name:  "Bruno 2",
			Email: "b2@example.com",
		},
		&pb.User{
			Id:    "b3",
			Name:  "Bruno 3",
			Email: "b3@example.com",
		},
		&pb.User{
			Id:    "b4",
			Name:  "Bruno 4",
			Email: "b4@example.com",
		},
		&pb.User{
			Id:    "b5",
			Name:  "Bruno 5",
			Email: "b5@example.com",
		},
	}

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	wait := make(chan int)

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
