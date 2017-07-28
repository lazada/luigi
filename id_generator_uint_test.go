package luigi

import (
	"sync"
	"testing"
	"time"
)

func Test_GeneratorUintIsUnique(c *testing.T) {
	var (
		uid  uint64
		uids = make(map[uint64]struct{})
		err  error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	for n := 0; n < 1000000; n++ {
		uid, err = uidGenerator.GenerateUint()
		if err != nil {
			panic(err)
		}

		if _, ok := uids[uid]; ok {
			c.Logf("Got non-unique ID: %v\n", uid)
			c.Fail()
		} else {
			uids[uid] = struct{}{}
		}
	}
}

func Test_GeneratorUintChannelIsUnique(c *testing.T) {
	const iterationsCount = 1000000

	var (
		uid     uint64
		uids    = make(chan uint64, iterationsCount)
		uidsMap = make(map[uint64]struct{})
		err     error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	errsCh := uidGenerator.FillChannelUint(uids)

	for n := 0; n < iterationsCount; n++ {
		select {
		case uid = <-uids:
		//Nothing to do
		case err = <-errsCh:
			panic(err)
		}

		if _, ok := uidsMap[uid]; ok {
			c.Logf("Got non-unique ID: %q\n", uid)
			c.Fail()
		} else {
			uidsMap[uid] = struct{}{}
		}
	}
}

func Test_GeneratorUintSliceIsUnique(c *testing.T) {
	const iterationsCount = 1000000

	var (
		uid     uint64
		uidsMap = make(map[uint64]struct{})
		err     error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	uids, err := uidGenerator.GenerateSliceUint(iterationsCount)
	if err != nil {
		panic(err)
	}

	for n := 0; n < iterationsCount; n++ {
		uid = uids[n]
		if _, ok := uidsMap[uid]; ok {
			c.Logf("Got non-unique ID: %q\n", uid)
			c.Fail()
		} else {
			uidsMap[uid] = struct{}{}
		}
	}
}

func BenchmarkSet_GenerateUintUniqueID(b *testing.B) {
	var (
		uid uint64
		err error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uid, err = uidGenerator.GenerateUint()
		if err != nil {
			panic(err)
		}

		_ = uid
	}
}

func BenchmarkSet_GenerateUintUniqueIDConcurrent(b *testing.B) {
	var (
		uid uint64
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
			uid, err = uidGenerator.GenerateUint()
			if err != nil {
				panic(err)
			}

			_ = uid

			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkSet_GenerateUintUniqueIDFillChannel(b *testing.B) {
	var (
		uid uint64
		err error
		ch  = make(chan uint64, 10000000)
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	errChan := uidGenerator.FillChannelUint(ch)

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

func BenchmarkSet_GenerateUintUniqueIDGetSlice(b *testing.B) {
	var err error

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	b.ResetTimer()
	uids, err := uidGenerator.GenerateSliceUint(uint32(b.N))
	if err != nil {
		panic(err)
	}
	b.StopTimer()

	_ = uids
}
