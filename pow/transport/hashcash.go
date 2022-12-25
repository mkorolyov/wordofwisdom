package transport

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
)

// HashcashChallenge send from server to client
type HashcashChallenge struct {
	Target *big.Int
	Puzzle []byte
}

func (p HashcashChallenge) Serialize() []byte {
	return bytes.Join([][]byte{
		SerializeSlice(p.Target.Bytes()),
		SerializeSlice(p.Puzzle),
	}, nil,
	)
}

func (p *HashcashChallenge) Deserialize(r io.Reader) error {
	targetBytes, err := DeserializeSlice(r)
	if err != nil {
		return fmt.Errorf("read target bytes: %w", err)
	}
	target := big.NewInt(0)
	target.SetBytes(targetBytes)
	p.Target = target

	puzzeBytes, err := DeserializeSlice(r)
	if err != nil {
		return fmt.Errorf("read puzzle bytes: %w", err)
	}
	p.Puzzle = puzzeBytes

	return nil
}

// HashcashResponse sent back to server from client as response to the challenge
type HashcashResponse struct {
	Counter []byte
}

func (p HashcashResponse) Serialize() []byte {
	return SerializeSlice(p.Counter)
}

func (p *HashcashResponse) Deserialize(r io.Reader) error {
	counterBytes, err := DeserializeSlice(r)
	if err != nil {
		return fmt.Errorf("read hashcash response: %w", err)
	}
	p.Counter = counterBytes

	return nil
}
