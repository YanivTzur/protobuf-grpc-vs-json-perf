package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type user struct {
	FirstName string
	LastName  string
	Age       uint64
}

const port = 8000

func parseCertificateFileName() string {
	if os.Args[1] == "--certificate" {
		return os.Args[2]
	}
	return os.Args[4]
}

func parsePrivateKeyFileName() string {
	if os.Args[1] == "--key" {
		return os.Args[2]
	}
	return os.Args[4]
}

func main() {
	if len(os.Args) < 5 || os.Args[1] != "--certificate" && os.Args[1] != "--key" ||
		os.Args[3] != "--certificate" && os.Args[3] != "--key" {
		log.Fatalln("Usage: ./server --certificate [TLS certificate] --key [TLS private key]")
	}

	listeningAddress := fmt.Sprintf("localhost:%v", port)
	srv := &http.Server{Addr: listeningAddress, Handler: http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatalln("Failed to read request's body")
			}
			newUser := user{}
			if err := json.Unmarshal(body, &newUser); err != nil {
				log.Fatalln("Failed to serialize new user from request body's JSON")
			}

			w.WriteHeader(http.StatusOK)
		})}

	certificateFileName := parseCertificateFileName()
	privateKeyFileName := parsePrivateKeyFileName()
	log.Println("Serving on", listeningAddress)
	log.Fatal(srv.ListenAndServeTLS(certificateFileName, privateKeyFileName))
}
