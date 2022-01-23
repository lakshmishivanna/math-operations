package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"log"
	"testing"

	pb "example.com/go-techmgmt-grpc/techmgmt"
	"github.com/stretchr/testify/assert"
)

func TestAddNmbers(t *testing.T) {

	con := new(context.Context)
	serv := NewPostService()

	r, err := serv.AddNmbers(*con, &pb.Numbers{Number1: 10, Number2: 20})

	if err != nil {
		log.Fatalf("Test Package: could not calculate: %v", err)
		assert.Error(t, err)
	}

	retHashSum, _ := base64.StdEncoding.DecodeString(r.GetHashsum())
	log.Printf("Test Package: HashSum: %x", retHashSum[:])

	exectedHashVal := GenerateHash(&pb.Numbers{Number1: 10, Number2: 20})
	hashVal, _ := base64.StdEncoding.DecodeString(exectedHashVal.GetHashsum())

	assert.Equal(t, retHashSum, hashVal, "Hash is not successful")
}

func TestAddNmbersWithNoValue(t *testing.T) {

	con := new(context.Context)
	serv := NewPostService()

	r, err := serv.AddNmbers(*con, &pb.Numbers{})

	if err != nil {
		log.Fatalf("Test Package: could not calculate: %v", err)
		assert.Error(t, err)
	} else if r.GetErrCode() != 0 {
		log.Printf("Test Package: Error code and Response: %d - %v", r.GetErrCode(), r.GetErrorResponse())
	} else {
		hashSum, _ := base64.StdEncoding.DecodeString(r.GetHashsum())
		log.Printf("Test Package: HashSum: %x", hashSum[:])
	}
}

func GenerateHash(in *pb.Numbers) *pb.SumOfIntegers {
	var addNum uint32 = uint32(in.GetNumber1() + in.GetNumber2())

	byteTemp := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteTemp, addNum)
	cryptHashSum := sha256.Sum256(byteTemp)

	return &pb.SumOfIntegers{Hashsum: base64.StdEncoding.EncodeToString(cryptHashSum[:])}
}

func BenchmarkAddNmbers(b *testing.B) {
	con := new(context.Context)
	serv := NewPostService()

	for i := 0; i < b.N; i++ {
		serv.AddNmbers(*con, &pb.Numbers{Number1: 10, Number2: 20})
	}
}
