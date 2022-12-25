package hashcash

import (
	"crypto/sha256"
	"testing"
)

func Test_HashcashPOW(t *testing.T) {
	hashcash := New(Config{
		Difficulty: 20,
		DoerConfig: DoerConfig{
			Hasher: sha256.New,
		},
	})

	work := hashcash.NewWork()
	counter := hashcash.DoWork(work)
	if !hashcash.VerifyWorkDone(counter[:], work.Nonce) {
		t.Fail()
	}
}
