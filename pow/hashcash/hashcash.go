package hashcash

import (
	"encoding/binary"
	"hash"
	"log"
	"math"
	"math/big"
)

/*
Generally `difficulty` should occupy less than `hashSize` in bits.
`hashSize-difficulty` is the target client will have to be close to.
Increasing our difficulty will increase the runtime of our algorithm.
*/

// uint64 in byte slice representation in binary.BigEndian ordering
type UInt64 [8]byte

// Verifier Implements the Hashcash algo as a Proof Of Work technic
type Verifier struct {
	Doer

	target *big.Int
	salt   []byte

	rand randGenerator
}

type Doer struct {
	hasher func() hash.Hash
}

type DoerConfig struct {
	Hasher func() hash.Hash
}

func NewDoer(cfg DoerConfig) *Doer {
	return &Doer{
		hasher: cfg.Hasher,
	}
}

// randGenerator encapsulates random number and bytes generation
type randGenerator interface {
	uint64() uint64
	bytes(n uint32) []byte
}

type Config struct {
	DoerConfig

	Difficulty int
}

func New(cfg Config) *Verifier {
	// the size of the resulting hash algo operates
	hashSizeBits := cfg.Hasher().Size() * 8
	// shift which leads to first `difficulty` zeros
	targetBitShift := uint(hashSizeBits - cfg.Difficulty)

	target := big.NewInt(1)
	target.Lsh(target, targetBitShift)

	rnd := newRandGen()
	return &Verifier{
		Doer: Doer{hasher: cfg.Hasher},

		target: target,
		salt:   rnd.bytes(8),

		rand: rnd,
	}
}

type Work struct {
	Nonce  UInt64
	Target *big.Int
	Hash   []byte
}

func (h *Verifier) NewWork() Work {
	var nonce UInt64
	// gen new rand uint64 and fill it to [8]byte array
	binary.BigEndian.PutUint64(nonce[:], h.rand.uint64())

	hasher := h.hasher()
	hasher.Write(h.salt)
	hasher.Write(nonce[:])

	return Work{
		Nonce:  nonce,
		Target: h.target,
		Hash:   hasher.Sum(nil),
	}
}

func (h *Verifier) VerifyWorkDone(workResult []byte, nonce UInt64) bool {
	hasher := h.hasher()
	hasher.Write(h.salt)
	hasher.Write(nonce[:])
	work := hasher.Sum(nil)

	hasher.Reset()
	return verifyWork(work, workResult, h.target, hasher)
}

func (h *Doer) DoWork(work Work) UInt64 {
	var counter UInt64
	hasher := h.hasher()
	for i := uint64(0); i < math.MaxUint64; i++ {
		attempt := counter[:]
		binary.BigEndian.PutUint64(attempt, i)

		hasher.Reset()

		if verifyWork(work.Hash[:], attempt, work.Target, hasher) {
			log.Printf("work done in %d iterations\n", i)
			break
		}
	}

	return counter
}

func verifyWork(work, workResult []byte, target *big.Int, hasher hash.Hash) bool {
	hasher.Write(work)
	hasher.Write(workResult)

	var result big.Int
	result.SetBytes(hasher.Sum(nil))

	return result.Cmp(target) != 1
}
