package luigi

import (
	"math/big"
	"sync/atomic"
)

// Format for 128bit: [timestamp 64bit][hostID 32bit][processID 16bit][counter 16bit]
func (this UIDGenerator) Generate() (*big.Int, error) {
	now, err := this.getTimeNanoseconds()
	if err != nil {
		return nil, err
	}

	id := atomic.AddUint32(this.sequence, 1)
	currentID := uint64(id & maxSequence)
	currentNodeID := uint64(this.nodeID | currentID)
	currentTS := now * this.timestampRightPadding

	var (
		bigCurrentNodeID = &big.Int{}
		bigCurrentTS     = &big.Int{}
		bigUid           = &big.Int{}
	)
	bigCurrentNodeID.SetUint64(currentNodeID)
	bigCurrentTS.SetUint64(currentTS)
	bigUid.Add(bigCurrentNodeID, bigCurrentTS)

	return bigUid, nil
}

func (this UIDGenerator) GenerateSlice(n uint32) ([]big.Int, error) {
	uids := make([]big.Int, n, n)
	id := atomic.AddUint32(this.sequence, n)

	var (
		currentID     uint64
		currentNodeID uint64
		now           uint64
		currentTS     uint64
		err           error
	)

	for i := uint32(0); i < n; i++ {
		currentID = uint64((id - (n - i - 1)) & maxSequence)
		if now, err = this.getTimeNanoseconds(); err != nil {
			return nil, err
		}

		currentNodeID = uint64(this.nodeID | currentID)
		currentTS = now * this.timestampRightPadding

		var (
			bigCurrentNodeID = &big.Int{}
			bigCurrentTS     = &big.Int{}
		)
		bigCurrentNodeID.SetUint64(currentNodeID)
		bigCurrentTS.SetUint64(currentTS)

		uids[i].Add(bigCurrentNodeID, bigCurrentTS)
	}

	return uids, nil
}

func (this UIDGenerator) FillChannel(ch chan<- *big.Int) chan error {
	var (
		err        error
		errChannel = make(chan error)
		uid        *big.Int
	)

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				close(errChannel)
				return
			}
		}()

		for {
			uid, err = this.Generate()
			if err != nil {
				errChannel <- err
				close(errChannel)
				close(ch)
				return
			}

			ch <- uid
		}
	}()

	return errChannel
}
