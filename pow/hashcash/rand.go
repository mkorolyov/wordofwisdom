package hashcash

import (
	crypto_rand "crypto/rand"
	"log"
	"math/rand"
	"time"
)

type randGen struct {
	rand *rand.Rand
}

func (r *randGen) uint64() uint64 {
	return r.rand.Uint64()
}

func (r randGen) bytes(n uint32) []byte {
	b := make([]byte, n)
	_, err := crypto_rand.Read(b)
	if err != nil {
		log.Fatalf("cant generate salt for proof of work: %v", err)
	}

	return b
}

func newRandGen() *randGen {
	return &randGen{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}
