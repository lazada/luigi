package luigi

import (
	"math/big"
	"sync"
	"time"
	"testing"
)

func BenchmarkSet_GenerateUniqueID(b *testing.B) {
	var (
		uid *big.Int
		err error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uid, err = uidGenerator.Generate()
		if err != nil {
			panic(err)
		}

		_ = uid
	}
}

func BenchmarkSet_GenerateUniqueIDConcurrent(b *testing.B) {
	var (
		uid *big.Int
		err error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	wg := sync.WaitGroup{}
	wg.Add(b.N)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		go func() {
			uid, err = uidGenerator.Generate()
			if err != nil {
				panic(err)
			}

			_ = uid

			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkSet_GenerateUniqueIDFillChannel(b *testing.B) {
	var (
		uid *big.Int
		err error
		ch  = make(chan *big.Int, 10000000)
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	errChan := uidGenerator.FillChannel(ch)

	time.Sleep(1 * time.Second)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uid = <-ch
		_ = uid
	}
	b.StopTimer()

	close(ch)
	for err = range errChan {
		panic(err)
	}
}

func BenchmarkSet_GenerateUniqueIDGetSlice(b *testing.B) {
	var err error

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	b.ResetTimer()
	uids, err := uidGenerator.GenerateSlice(uint32(b.N))
	if err != nil {
		panic(err)
	}
	b.StopTimer()

	_ = uids
}
