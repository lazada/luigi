# Luigi
## Lazada unique identifier generator

![Luigi](http://38.media.tumblr.com/4428256ec89da177407fbd6b916d4798/tumblr_inline_nl0r52henW1r1t6o9.jpg)

Algorithm:

1. NodeID = ((HostID << PIDLengthInBits) | PID) << counterLengthInBits

2. Timestamp = Now() - ToTimestamp(20.02.2016, 14:23:57)

3. UniqueID = Timestamp | NodeID | Counter

Counter - increment atomic counter less than predefined defaultSequenceBits const

Luigi provides 3 basic methods for IDs generate:

* Generate() - single value
* GenerateSlice(n) - slice of UIDs with n values
* FillChannelFillChannel(ch chan<- *big.Int) chan error - continuously fill ch channel with UIDs

Luigi can generate UIDs of 3 types: **uint64**(_supports millisecond timestamp_), **string** and **BigInt**(_supports nanosecond timestamp_)

Luigi loves you!