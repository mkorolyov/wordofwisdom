package main

import (
	"context"
	"flag"
	"log"

	"github.com/mkorolyov/wordofwisdom/handler"
	"github.com/mkorolyov/wordofwisdom/server"
)

var (
	addr = flag.String("addr", ":9991", "")
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	cfg := server.Config{
		Addr:                *addr,
		ConcurrentConnLimit: 0,
	}

	srv := server.NewTCPServer(cfg)
	ctx := context.Background()

	srv.Serve(ctx, handler.DDoSProtection(handler.RandomQuotes()))
}
