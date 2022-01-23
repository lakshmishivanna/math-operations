package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	pb "example.com/go-techmgmt-grpc/techmgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50055"
)

type TechMangementServer struct {
	pb.UnimplementedTechMangementServer
}

func (s *TechMangementServer) AddNmbers(ctx context.Context, in *pb.Numbers) (*pb.SumOfIntegers, error) {
	num1 := in.GetNumber1()
	num2 := in.GetNumber2()
	log.Printf("Server: Received: %v,%v", num1, num2)

	if num1 == 0 && num2 == 0 {
		return &pb.SumOfIntegers{ErrorResponse: "Numbers should be greater than 0", ErrCode: 200}, nil
	}
	var addNum uint32 = uint32(in.GetNumber1() + in.GetNumber2())

	byteTemp := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteTemp, addNum)
	cryptHashSum := sha256.Sum256(byteTemp)

	log.Printf("Server: Crytographic hash of Sum %x", cryptHashSum[:])
	return &pb.SumOfIntegers{Hashsum: base64.StdEncoding.EncodeToString(cryptHashSum[:])}, nil
}

func NewPostService() TechMangementServer {
	return TechMangementServer{}
}

func main() {

	runtime.SetBlockProfileRate(1)
	go func() {
		fmt.Println(http.ListenAndServe("localhost:8024", nil))
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Server: failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTechMangementServer(s, &TechMangementServer{})
	log.Printf("Server: server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { // start server
		log.Fatalf("Server: failed to serve: %v", err)
	}

}
