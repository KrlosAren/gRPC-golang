package main

import (
	"context"
	"io"
	"krlosaren/go/grpc/testpbf"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Couldn't connect to %v", err)
	}

	defer cc.Close()

	c := testpbf.NewTestServiceClient(cc)

	// DoUnary(c)
	// DoClientStreaming(c)
	// DoServerStreaming(c)
	DoBidirectionalStreaming(c)
}
func DoUnary(c testpbf.TestServiceClient) {
	req := &testpbf.GetTestRequest{
		Id: "t1",
	}

	res, err := c.GetTest(context.Background(), req)

	if err != nil {
		log.Fatalf("Couldn't get test %v", err)
	}

	log.Printf("Response: %v", res)
}

func DoClientStreaming(c testpbf.TestServiceClient) {
	questions := []*testpbf.Question{
		{
			Id:       "q8t1",
			Answer:   "Azul",
			Question: "color saso",
			TestId:   "t1",
		},
		{
			Id:       "q9t1",
			Answer:   "Azul",
			Question: "color saso!!!",
			TestId:   "t1",
		},
		{
			Id:       "q10t1",
			Answer:   "Azul",
			Question: "color saso!!!!!!",
			TestId:   "t1",
		},
		{
			Id:       "q11t1",
			Answer:   "Azul",
			Question: "color saso!!!!",
			TestId:   "t1",
		},
	}

	stream, err := c.SetQuestion(context.Background())

	if err != nil {

		log.Fatalf("Could not set question stream %v", err)
	}

	for _, q := range questions {
		log.Println("Send question to server ", q.Id)
		stream.Send(q)
		time.Sleep(2 * time.Second)
	}

	msg, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Could not send question to server  %v", err)
	}

	log.Printf("Server responded with %v", msg)
}

func DoServerStreaming(c testpbf.TestServiceClient) {
	req := &testpbf.GetStudentsPerTestRequest{
		TestId: "t1",
	}

	stream, err := c.GetStudentsPerTest(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not get students per test: %v", err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive students: %v", err)
		}

		log.Printf("Got students from server %v", msg)

	}
}

func DoBidirectionalStreaming(c testpbf.TestServiceClient) {
	answer := testpbf.TakeTestRequest{
		Answer: "42",
	}

	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())

	if err != nil {
		log.Fatalf("Could not get students: %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			stream.Send(&answer)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				log.Fatalf("io.EOF returned %v", res)
				break
			}

			if err != nil {
				log.Fatalf("Could not receive students: %v", err)
				break
			}

			log.Printf("Received students: %v", res)
		}
		close(waitChannel)
	}()

	<-waitChannel

}
