package handler

import (
	"crypto/sha256"
	"fmt"
	"net"
	"time"

	"github.com/mkorolyov/wordofwisdom/pow/hashcash"
	"github.com/mkorolyov/wordofwisdom/pow/transport"
	"github.com/mkorolyov/wordofwisdom/server"
)

func DDoSProtection(h server.Handler) server.Handler {
	// We are using Hashcash algo here for Proof of Work protection.
	// 20 bit (near 1kk attempts for the search to succeed, like email spam filters configured)
	// e.g. bitcoin uses ~67.5(its algo constantly increases difficulty to limit new blocks gen speed to desired number) bits which led to 200kkkk attempts.
	// this also means that first `difficulty` bits of the hash in the result must be zeros.
	hashcashPOW := hashcash.New(hashcash.Config{
		Difficulty: 20,
		DoerConfig: hashcash.DoerConfig{Hasher: sha256.New},
	})

	return func(conn net.Conn) error {
		work := hashcashPOW.NewWork()

		if err := conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
			return fmt.Errorf("cant limit time to read the challenge by client: %v", err)
		}

		if err := sendHashcashChallenge(conn, work); err != nil {
			return err
		}

		// TODO this deadline could be calculated from difficulty with some reosanable time gap for network roundtrip
		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			return fmt.Errorf("cant limit time to solve the challenge by client: %v", err)
		}

		if err := verifyHashcashResponse(conn, hashcashPOW, work); err != nil {
			return err
		}

		if err := conn.SetDeadline(time.Time{}); err != nil {
			return fmt.Errorf("cant remove previously set deadlines: %v", err)
		}

		return h(conn)
	}
}

func sendHashcashChallenge(conn net.Conn, work hashcash.Work) error {
	challenge := transport.HashcashChallenge{
		Target: work.Target,
		Puzzle: work.Hash,
	}

	if err := transport.WriteSlice(conn, challenge.Serialize()); err != nil {
		return fmt.Errorf("cant write hashcashVerifier client puzzle to tcp conn: %w", err)
	}

	return nil
}

func verifyHashcashResponse(conn net.Conn, hashcashVerifier *hashcash.Verifier, work hashcash.Work) error {
	var clientResponse transport.HashcashResponse
	if err := clientResponse.Deserialize(conn); err != nil {
		return fmt.Errorf("cant read hashcashVerifier client response: %w", err)
	}

	if !hashcashVerifier.VerifyWorkDone(clientResponse.Counter, work.Nonce, work.Timestamp) {
		// no need to notify bad client and spend more resources on it.
		// TODO introduce typed err and handle like expected scenario
		return fmt.Errorf("client didnt solve the puzzle")
	}

	return nil
}
