package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/tknoptn/avg-stuff-grpc/proto"
)

func callSayHelloBidirectonalStream(client pb.MeetServiceClient, names *pb.NamesList) {
	log.Printf("Bidirectional STREAMING STARTED")

	stream, err := client.SayHelloBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("cpuld not send names : %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v ", err)
		}
		time.Sleep(2 * time.Second)
	}

	stream.CloseSend()

	<-waitc

	log.Printf("BIdirectional streaming fininshed")

}
