package handler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/mkorolyov/wordofwisdom/pow/transport"
	"github.com/mkorolyov/wordofwisdom/server"
)

func RandomQuotes(source io.Reader) server.Handler {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	quotes, err := ScanQuotes(source)
	if err != nil {
		log.Fatalf("cant build quotes Handler: %v", err)
	}

	return func(conn net.Conn) error {
		idx := r.Intn(len(quotes))

		if err := transport.WriteSlice(conn, transport.SerializeSlice(quotes[idx])); err != nil {
			return fmt.Errorf("send quote to client: %w", err)
		}

		return nil
	}
}

func ScanQuotes(source io.Reader) ([][]byte, error) {
	var quotes [][]byte
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		quotes = append(quotes, scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan quotes: %w", err)
	}

	return quotes, nil
}
