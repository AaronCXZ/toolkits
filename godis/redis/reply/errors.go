package reply

// UnKnownErrReply 未知错误
type UnKnownErrReply struct{}

var unKnownErrBytes = []byte("-Err unknown\r\n")

func (r *UnKnownErrReply) ToBytes() []byte {
	return unKnownErrBytes
}

func (r *UnKnownErrReply) Error() string {
	return "Err unknown"
}

// ArgNumErrReply 参数错误
type ArgNumErrReply struct {
	Cmd string
}

func MakeArgUnmErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

func (r *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + r.Cmd + "' command"
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

// SyntaxErrReply 语法错误
type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}
func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

// WrongTypeErrReply 类型错误
type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

// ProtocolErrReply 协议错误
type ProtocolErrReply struct {
	Msg string
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + r.Msg
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}
