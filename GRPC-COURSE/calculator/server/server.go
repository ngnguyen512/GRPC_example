package main

import (
	"context"
	"fmt"
	"github/ngnguyen512/GRPC-COURSE/calculator/calculatorpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("sum called")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (*server) PrimeNumber(req *calculatorpb.PNRequest, stream calculatorpb.CalculatorService_PrimeNumberServer) error {
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			stream.Send(&calculatorpb.PNResponse{
				Result: k,
			})
		} else {
			k++
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(resp)

		}
		if err != nil {
			log.Fatal("err while Recv Average", err)
			return err
		}
		log.Println("receive req", req)
		total += req.GetNum()
		count++
	}
}

func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	log.Println("Find max called...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatal("err while Recv FindMax", err)
			return err
		}
		num := req.GetNum()
		if num > max {
			max = num
		}
		err = stream.Send(&calculatorpb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatal("send max err", err)
			return err
		}
	}
}

// func (*server) SumwithDeadline(ctx context.Context, req *calculatorpb.SumRequest) *calculatorpb.SumResponse {
// 	log.Println("sum with called")
// 	fo i = 0; i <3; i++ {
// 		time.Sleep(1 *time.Second)
// 	}
// 	resp := &calculatorpb.SumResponse{
// 		Result: req.GetNum1() + req.GetNum2(),
// 	}

// 	return resp, nil
// }

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")

	if err != nil {
		log.Fatal("err while create listen", err)
	}
	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	fmt.Println("calculator is running...")

	err = s.Serve(lis)

	if err != nil {
		log.Fatal("err while serve", err)
	}
}
