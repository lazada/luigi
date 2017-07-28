package luigi

const (
	GeneratorStartTimeMilli = 1455967437000
	GeneratorStartTimeNano  = 1455967437000000000

	defaultSequenceBits = 16

	defaultPIDBits             = 17
	defaultPIDBitsRightPadding = defaultSequenceBits
	defaultPIDRightPadding     = 10 * defaultSequenceBits

	defaultHostBits             = 16
	defaultHostBitsRightPadding = defaultSequenceBits + defaultPIDBitsRightPadding
	defaultHostRightPadding     = 10 * defaultHostBitsRightPadding

	defaultNodeBits             = defaultHostBits + defaultPIDBits
	defaultNodeBitsRightPadding = defaultSequenceBits
	defaultNodeRightPadding     = 10 * defaultSequenceBits

	defaultTimestampBits             = 63
	defaultTimestampBitsRightPadding = defaultNodeBits + defaultSequenceBits
	defaultTimestampRightPadding     = 10 * defaultTimestampBitsRightPadding

	defaultUIDBits = defaultTimestampBits + defaultNodeBits + defaultSequenceBits

	// Set up limits based on configured options
	maxSequence = (1 << defaultSequenceBits) - 1 // sequence mask
	maxPID      = (1 << defaultPIDBits) - 1      // pid mask
	maxHostID   = (1 << defaultHostBits) - 1     // host mask
	maxNodeID   = (1 << defaultNodeBits) - 1     // hid mask

	//Use: int64((uint64(1<<defaultTimestampBits) - 1)) in case defaultTimestampBits < 64
	//Use: uint64(int64(^uint64(0) >> 1)) in case defaultTimestampBits = 64
	maxAdjustedTimestamp = uint64(1<<defaultTimestampBits) - 1 //  timestamp mask

	OverflowErrorMessage = "Too big %s: %+v. Max: %v"
)
