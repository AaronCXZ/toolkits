package reply

type PongReply struct{}

var pongBytes = []byte("+PONG\r\n")

func (r *PongReply) ToBytes() []byte {
	return pongBytes
}

type OkReply struct{}

var okBytes = []byte("+OK\r\n")
var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

type NoReply struct{}

var noBytes = []byte("")

func (r *NoReply) ToBytes() []byte {
	return noBytes
}

type QueuedReply struct{}

var queuedBytes = []byte("+QUEUED\r\n")
var theQueuedReply = new(QueuedReply)

func MakeQueuedReply() *QueuedReply {
	return theQueuedReply
}

func (r *QueuedReply) ToBytes() []byte {
	return queuedBytes
}
