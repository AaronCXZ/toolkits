package redis

type Connection interface {
	Write([]byte) error
	SetPassword(string)
	GetPassword() string

	Subscribe(channel string)
	UnSubscribe(channel string)
	SubsCount() int
	GetChannels() []string

	InMultiState() bool
	SetMultiState(bool)
	GetQueuedCmdLine() [][][]byte
	EnqueueCmd([][]byte)
	ClearQueuedCmds()
}
