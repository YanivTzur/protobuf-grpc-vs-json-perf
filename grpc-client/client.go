package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "grpc-json-comparison-go/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:50051"
	usage   = "Usage: ./client --num-requests [number of requests] --num-iterations [number of iterations]"
)

func getNumRequests() (uint64, error) {
	var arg string
	if os.Args[1] == "--num-requests" {
		arg = os.Args[2]
	} else {
		arg = os.Args[4]
	}
	return strconv.ParseUint(arg, 10, 32)
}

func getNumIterations() (uint64, error) {
	var arg string
	if os.Args[1] == "--num-iterations" {
		arg = os.Args[2]
	} else {
		arg = os.Args[4]
	}
	return strconv.ParseUint(arg, 10, 32)
}

func parseClientArgs() (uint64, uint64, error) {
	if len(os.Args) < 5 {
		log.Fatalln(usage)
	}
	if (os.Args[1] != "--num-requests" && os.Args[1] != "--num-iterations") ||
		(os.Args[3] != "--num-requests" && os.Args[3] != "--num-iterations") {
		log.Fatalln(usage)
	}
	numRequests, err := getNumRequests()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to parse number of requests: %v", err)
	}
	numIterations, err := getNumIterations()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to parse number of iterations: %v", err)
	}

	return numRequests, numIterations, nil
}

func main() {
	numRequests, numIterations, err := parseClientArgs()
	if err != nil {
		log.Fatalln("Failed to parse command line arguments:", err)
	}

	creds, _ := credentials.NewClientTLSFromFile("server.cert", "")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	var i uint64
	var iterationTimesMs []uint64
	for i = 0; i < numIterations; i++ {
		iterationStartTime := time.Now()
		var j uint64
		for j = 0; j < numRequests; j++ {
			_, err = c.AddUser(ctx, &pb.User{Age: 74, FirstName: "Bill", LastName: "Clinton"})
			if err != nil {
				log.Fatalf("Failed to add user. Error: %v", err)
			}
		}
		iterationTimeElapsed := time.Now().Sub(iterationStartTime).Milliseconds()
		iterationTimesMs = append(iterationTimesMs, uint64(iterationTimeElapsed))
		log.Println("Finished", i+1, "iterations out of", numIterations)
	}

	var sum uint64
	for i = 0; i < numIterations; i++ {
		sum += iterationTimesMs[i]
	}
	averageIterationTimeMs := sum / numIterations

	log.Println(
		"Successfully added users. Number of iterations:",
		numIterations,
		", number of requests:",
		numRequests,
		", average time per iteration in milliseconds:",
		averageIterationTimeMs)
}
