package luigi

import "sync/atomic"

// Format: [timestamp 64bit][hostID 32bit][processID 16bit][counter 16bit]
func (this UIDGenerator) GenerateUint() (uint64, error) {
	now, err := this.getTimeMicroseconds()
	if err != nil {
		return 0, err
	}

	currentID := this.getNextUintID()
	currentTS := now << defaultTimestampBitsRightPadding

	return currentTS | this.nodeID | currentID, nil
}

func (this UIDGenerator) GenerateSliceUint(n uint32) ([]uint64, error) {
	uids := make([]uint64, n, n)
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
		if now, err = this.getTimeMicroseconds(); err != nil {
			return nil, err
		}

		currentNodeID = uint64(this.nodeID | currentID)
		currentTS = now * this.timestampRightPadding

		uids[i] = currentTS + currentNodeID + currentID
	}

	return uids, nil
}

func (this UIDGenerator) FillChannelUint(ch chan<- uint64) chan error {
	var (
		err        error
		errChannel = make(chan error)
		uid        uint64
	)

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				close(errChannel)
				return
			}
		}()

		for {
			uid, err = this.GenerateUint()
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

func (this UIDGenerator) getNextUintID() uint64 {
	id := atomic.AddUint32(this.sequence, 1)

	return uint64(id & maxSequence)
}
