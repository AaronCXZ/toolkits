package reply

import (
	"bytes"
	"strconv"

	"github.com/Muskchen/toolkits/godis/interface/redis"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

// BulkReply 单行数据Reply接口实现
type BulkReply struct {
	Arg []byte
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

func (r *BulkReply) ToBytes() []byte {
	if len(r.Arg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

// MultiBulkReply 多行数据的Reply接口实现
type MultiBulkReply struct {
	Args [][]byte
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

func (r *MultiBulkReply) ToBytes() []byte {
	argLen := len(r.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString("$-1" + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

// MultiRawReply 嵌套的Reply接口实现
type MultiRawReply struct {
	Replies []redis.Reply
}

func MakeMUltiRawReply(replies []redis.Reply) *MultiRawReply {
	return &MultiRawReply{
		Replies: replies,
	}
}

func (r *MultiRawReply) ToBytes() []byte {
	argLen := len(r.Replies)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range r.Replies {
		buf.Write(arg.ToBytes())
	}
	return buf.Bytes()
}

// StatusReply 简单字符串的Reply接口实现
type StatusReply struct {
	Status string
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

// IntReply Int类型的Reply接口实现
type IntReply struct {
	Code int64
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

// ErrorReply 错误数据的Reply接口
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// StandardErrReply 错误接口实现
type StandardErrReply struct {
	Status string
}

func MakeErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

func IsErrorReply(reply redis.Reply) bool {
	return reply.ToBytes()[0] == '-'
}

func (r *StandardErrReply) Error() string {
	return r.Status
}

func (r *StandardErrReply) ToBytes() []byte {
	return []byte("-" + r.Status + CRLF)
}
