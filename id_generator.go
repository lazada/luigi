package luigi

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

// TODO: make inherit generators
type UIDGenerator struct {
	timestampRightPadding uint64
	// nodeID - 32 bits of hostID + 16 bits of processID
	nodeID uint64
	// sequence number - 16 bits, auto-increment for prevent sametime collisions
	sequence *uint32

	// Limits based on configured options
}

var (
	pid                   uint32
	hostID                uint16
	nodeID                uint64
	timestampRightPadding uint64

	initErrors = []error{}
)

func init() {
	DoInit()
}

func DoInit(initValues ...uint32) {
	var err error

	if len(initValues) == 0 {
		pid = uint32(os.Getpid())
	} else {
		pid = initValues[0]
	}

	if pid > maxPID {
		// Fixme: this increase collisions
		pid %= maxPID

		log.Printf(OverflowErrorMessage, "PID", pid, maxPID)
	}

	hostID, err = getHostID()
	if err != nil {
		initErrors = append(initErrors, err)
	}
	if hostID > maxHostID {
		// Fixme: this increase collisions
		hostID %= maxHostID

		log.Printf(OverflowErrorMessage, "hostID", hostID, maxHostID)
	}

	currentHostID := uint64(hostID) << defaultPIDBits
	nodeID = currentHostID | uint64(pid)
	if nodeID > maxNodeID {
		// Fixme: this increase collisions
		nodeID %= maxNodeID

		log.Printf(OverflowErrorMessage, "nodeID", nodeID, maxNodeID)
	}

	nodeID = nodeID << defaultNodeBitsRightPadding

	timestampRightPadding = uint64(math.Pow(float64(2), float64(defaultTimestampBitsRightPadding)))
}

func getHostID() (uint16, error) {
	var (
		err      error
		hostname string
	)

	// In docker hostname==container ID
	if hostname, err = os.Hostname(); err != nil {
		return 0, errors.New(fmt.Sprint("os.Hostname: ", err))
	}

	buf := md5.New().Sum([]byte(hostname))
	hostID := binary.BigEndian.Uint16(buf)

	return hostID, nil
}

func (this UIDGenerator) getTimeMicroseconds() (uint64, error) {
	now := uint64((time.Now().UnixNano() / 1000000) - GeneratorStartTimeMilli)

	if now > maxAdjustedTimestamp {
		return 0, errors.New(fmt.Sprintf(OverflowErrorMessage, "timestamp", now))
	}

	return now, nil
}

func (this UIDGenerator) getTimeNanoseconds() (uint64, error) {
	now := uint64(time.Now().UnixNano() - GeneratorStartTimeNano)

	if now > maxAdjustedTimestamp {
		return 0, errors.New(fmt.Sprintf(OverflowErrorMessage, "timestamp", now))
	}

	return now, nil
}

func NewUIDGenerator(initValues ...uint32) (*UIDGenerator, []error) {
	if len(initErrors) > 0 {
		return nil, initErrors
	}

	sequenceValue := uint32(0)

	if len(initValues) > 0 {
		DoInit(initValues...)
	}

	generator := UIDGenerator{
		sequence:              &sequenceValue,
		timestampRightPadding: timestampRightPadding,
		nodeID:                nodeID,
	}

	return &generator, nil
}
