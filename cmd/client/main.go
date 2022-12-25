package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"net"
	"time"

	"github.com/mkorolyov/wordofwisdom/pow/hashcash"
	"github.com/mkorolyov/wordofwisdom/pow/transport"
)

var addr = flag.String("addr", ":9991", "")

func main() {
	flag.Parse()

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("net.Dial: %v", err)
	}

	defer func() { _ = conn.Close() }()

	solveChallenge(conn)

	// load quote
	quoteBytes, err := transport.DeserializeSlice(conn)
	if err != nil {
		log.Fatalf("read quote: %v", err)
	}

	log.Printf("received quote from the server: %q", quoteBytes)
}

func solveChallenge(conn net.Conn) {
	hashCashSolver := hashcash.NewDoer(hashcash.DoerConfig{Hasher: sha256.New})

	var challenge transport.HashcashChallenge
	if err := challenge.Deserialize(conn); err != nil {
		log.Fatalf("cant read hashcash challenge: %v", err)
	}

	now := time.Now()
	counter := hashCashSolver.DoWork(hashcash.Work{
		Target: challenge.Target,
		Hash:   challenge.Puzzle,
	})
	log.Printf("challenge solve took %s", time.Since(now))

	// send response to the server
	response := transport.HashcashResponse{Counter: counter[:]}
	if err := transport.WriteSlice(conn, response.Serialize()); err != nil {
		log.Fatalf("cant write hashcash challenge response: %v", err)
	}
}
