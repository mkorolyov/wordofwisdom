package hashcash

import (
	"crypto/sha256"
	"hash"
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
	if !hashcash.VerifyWorkDone(counter[:], work.Nonce, work.Timestamp) {
		t.Error("work done was not verified")
	}
}

func Benchmark_HashcashDoWork(b *testing.B) {
	bench := func(b *testing.B, hashFunc func() hash.Hash, difficulty int) {
		pow := New(Config{
			DoerConfig: DoerConfig{Hasher: hashFunc},
			Difficulty: difficulty,
		})
		work := pow.NewWork()

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			pow.DoWork(work)
		}
	}

	b.Run("sha256", func(b *testing.B) {
		b.Run("difficulty=1", func(b *testing.B) {
			bench(b, sha256.New, 1)
		})

		b.Run("difficulty=10", func(b *testing.B) {
			bench(b, sha256.New, 10)
		})

		b.Run("difficulty=15", func(b *testing.B) {
			bench(b, sha256.New, 15)
		})
	})
}
