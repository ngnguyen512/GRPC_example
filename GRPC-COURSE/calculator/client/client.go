package main

import (
	"context"
	"github/ngnguyen512/GRPC-COURSE/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.NewClient("localhost:50069", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error while dialing: ", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	// callSum(client)
	// callPND(client)
	// callAverage(client)
	callFindMax(client)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})
	if err != nil {
		log.Fatal("Call sum API error: ", err)
	}

	log.Printf("Sum API response: %d", resp.GetResult())
}

func callPND(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.PrimeNumber(context.Background(), &calculatorpb.PNRequest{
		Number: 120,
	})
	if err != nil {
		log.Fatal("Call PND error: ", err)
	}

	for {
		resp, recErr := stream.Recv()
		if recErr == io.EOF {
			log.Println("Server finished streaming")
			break
		}
		if recErr != nil {
			log.Fatal("Receive error: ", recErr)
		}
		log.Printf("Prime number: %d", resp.GetResult())
	}
}

func callAverage(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatal("call average err", err)
	}
	listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Num: 5,
		},
		calculatorpb.AverageRequest{
			Num: 10,
		},
		calculatorpb.AverageRequest{
			Num: 12,
		},
		calculatorpb.AverageRequest{
			Num: 3,
		},
		calculatorpb.AverageRequest{
			Num: 4.2,
		},
	}
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatal("send average request err", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receive average response err", err)
	}
	log.Printf("average response", resp)
}

func callFindMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling find max...")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call find max err", err)
	}
	waitc := make(chan struct{})
	go func() {
		listReq := []calculatorpb.FindMaxRequest{
			calculatorpb.FindMaxRequest{
				Num: 5,
			},
			calculatorpb.FindMaxRequest{
				Num: 10,
			},
			calculatorpb.FindMaxRequest{
				Num: 12,
			},
			calculatorpb.FindMaxRequest{
				Num: 3,
			},
			calculatorpb.FindMaxRequest{
				Num: 4,
			},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatal("send find max request err", err)
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("err", err)
				break
			}

			log.Println("max", resp.GetMax())
		}
		close(waitc)

	}()

	<-waitc
}
