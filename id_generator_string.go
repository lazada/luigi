package luigi

import (
	"strconv"
	"strings"
	"sync/atomic"
)

func (this UIDGenerator) GenerateString() (string, error) {
	now, err := this.getTimeNanoseconds()
	if err != nil {
		return "", err
	}

	currentID := this.getNextUintID()

	return this.getUIDString(now, currentID), nil
}

func (this UIDGenerator) GenerateSliceString(n uint32) ([]string, error) {
	uids := make([]string, n, n)
	id := atomic.AddUint32(this.sequence, n)

	var (
		currentID uint64
		now       uint64
		err       error
	)

	for i := uint32(0); i < n; i++ {
		currentID = uint64((id - (n - i - 1)) & maxSequence)
		if now, err = this.getTimeNanoseconds(); err != nil {
			return nil, err
		}

		uids[i] = this.getUIDString(now, currentID)
	}

	return uids, nil
}

func (this UIDGenerator) FillChannelString(ch chan<- string) chan error {
	var (
		err        error
		errChannel = make(chan error)
		uid        string
	)

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				close(errChannel)
				return
			}
		}()

		for {
			uid, err = this.GenerateString()
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

func (this UIDGenerator) getUIDString(now, currentID uint64) string {
	currentNodeID := uint64(this.nodeID | currentID)
	return strings.Join([]string{strconv.FormatUint(now, 10), strconv.FormatUint(currentNodeID, 10)}, "")
}
