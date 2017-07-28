package luigi

import (
	"sync"
	"time"
	"testing"
)

func Test_GeneratorStringIsUnique(c *testing.T) {
	var (
		uid  string
		uids = make(map[string]struct{})
		err  error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	for n := 0; n < 1000000; n++ {
		uid, err = uidGenerator.GenerateString()
		if err != nil {
			panic(err)
		}

		if _, ok := uids[uid]; ok {
			c.Logf("Got non-unique ID: %q\n", uid)
			c.Fail()
		} else {
			uids[uid] = struct{}{}
		}
	}
}

func Test_GeneratorStringChannelIsUnique(c *testing.T) {
	const iterationsCount = 1000000

	var (
		uid     string
		uids    = make(chan string, iterationsCount)
		uidsMap = make(map[string]struct{})
		err     error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	errsCh := uidGenerator.FillChannelString(uids)

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

func Test_GeneratorStringSliceIsUnique(c *testing.T) {
	const (
		iterationsCount = 10000
		repeat          = 5
	)

	var (
		uid     string
		uidsMap = make(map[string]struct{})
	)

	for i := 0; i < repeat; i++ {
		uidGenerator, errs := NewUIDGenerator()
		if len(errs) != 0 {
			panic(errs)
		}

		uids, err := uidGenerator.GenerateSliceString(iterationsCount)
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
}

func BenchmarkSet_GenerateStringUniqueID(b *testing.B) {
	var (
		uid string
		err error
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uid, err = uidGenerator.GenerateString()
		if err != nil {
			panic(err)
		}

		_ = uid
	}
}

func BenchmarkSet_GenerateStringUniqueIDConcurrent(b *testing.B) {
	var (
		uid string
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
			uid, err = uidGenerator.GenerateString()
			if err != nil {
				panic(err)
			}

			_ = uid

			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkSet_GenerateStringUniqueIDFillChannel(b *testing.B) {
	var (
		uid string
		err error
		ch  = make(chan string, 10000000)
	)

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	errChan := uidGenerator.FillChannelString(ch)

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

func BenchmarkSet_GenerateStringUniqueIDGetSlice(b *testing.B) {
	var err error

	uidGenerator, errs := NewUIDGenerator()
	if len(errs) != 0 {
		panic(errs)
	}

	b.ResetTimer()
	uids, err := uidGenerator.GenerateSliceString(uint32(b.N))
	if err != nil {
		panic(err)
	}
	b.StopTimer()

	_ = uids
}
