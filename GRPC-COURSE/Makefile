gen-cal:
	protoc calculator/calculatorpb/calculator.proto --go_out=. --go_opt=module=github.com/ngnguyen512/GRPC-COURSE --go-grpc_out=. --go-grpc_opt=module=github.com/ngnguyen512/GRPC-COURSE
run-server:
	go run calculator/server/server.go
run-client:
	go run calculator/client/client.go