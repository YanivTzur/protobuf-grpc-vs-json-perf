package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/http2"
)

const (
	port  = 8000
	usage = "Usage: ./client --num-requests [number of requests] --num-iterations [number of iterations]"
)

type user struct {
	FirstName string
	LastName  string
	Age       uint64
}

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

	client := &http.Client{}
	caCert, err := ioutil.ReadFile("server.cert")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	client.Transport = &http2.Transport{
		TLSClientConfig: tlsConfig,
	}

	var i uint64
	var iterationTimesMs []uint64
	for i = 0; i < numIterations; i++ {
		iterationStartTime := time.Now()
		var j uint64
		for j = 0; j < numRequests; j++ {
			jsonValue, _ := json.Marshal(user{FirstName: "Bill", LastName: "Clinton", Age: 74})
			resp, err := client.Post(
				fmt.Sprintf("https://localhost:%v", port),
				"application/json",
				bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Fatalf("Failed post: %s", err)
			}
			resp.Body.Close()
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
