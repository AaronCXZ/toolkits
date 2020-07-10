package errors

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func New(msg string) error {
	return errors.New(msg)
}

// 自定义error
type Frame struct {
	Err  error
	Pkg  string
	Func string
	Args []interface{}
	Code string
	File string
	Line int
}

func NewFrame(err error, code, file string, line int, pkg, fn string, args ...interface{}) *Frame {
	return &Frame{
		Err:  err,
		Pkg:  pkg,
		Func: fn,
		Args: args,
		Code: code,
		File: file,
		Line: line,
	}
}

func (p *Frame) Error() string {
	return string(errorDetail(make([]byte, 0, 32), p))
}

func Err(err error) error {
	if e, ok := err.(*Frame); ok {
		return Err(e.Err)
	}
	return err
}

func errorDetail(b []byte, p *Frame) []byte {
	if f, ok := p.Err.(*Frame); ok {
		b = errorDetail(b, f)
	} else {
		b = append(b, p.Err.Error()...)
		b = append(b, "\n\n===> errors stack:\n"...)
	}
	b = append(b, p.Pkg...)
	b = append(b, '.')
	b = append(b, p.Func...)
	b = append(b, '(')
	b = funcArgsDetail(b, p.Args)
	b = append(b, ")\n\t"...)
	b = append(b, p.File...)
	b = append(b, ':')
	b = strconv.AppendInt(b, int64(p.Line), 10)
	b = append(b, ' ')
	b = append(b, p.Code...)
	b = append(b, '\n')
	return b
}

func funcArgsDetail(b []byte, args []interface{}) []byte {
	nlast := len(args) - 1
	for i, arg := range args {
		b = appendValue(b, arg)
		if i != nlast {
			b = append(b, ',', ' ')
		}
	}
	return b
}

func appendValue(b []byte, arg interface{}) []byte {
	if arg == nil {
		return append(b, "nil"...)
	}
	v := reflect.ValueOf(arg)
	kind := v.Kind()
	if kind >= reflect.Bool && kind <= reflect.Complex128 {
		return append(b, fmt.Sprint(arg)...)
	}
	if kind == reflect.String {
		val := arg.(string)
		if len(val) > 16 {
			val = val[:16] + "..."
		}
		return strconv.AppendQuote(b, val)
	}
	if kind == reflect.Array {
		return append(b, "Array"...)
	}
	if kind == reflect.Struct {
		return append(b, "Struct"...)
	}
	val := v.Pointer()
	b = append(b, '0', 'x')
	return strconv.AppendInt(b, int64(val), 16)
}

func (p *Frame) Unwrap() error {
	return p.Err
}

func (p *Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, p.Error())
	case 's':
		io.WriteString(s, Err(p.Err).Error())
	case 'q':
		fmt.Fprintf(s, "%q", Err(p.Err).Error())
	}
}

type ErrorInfo = Frame

func (p *ErrorInfo) Detail(err error) *ErrorInfo {
	p.Code = err.Error()
	return p
}

func (p *ErrorInfo) NestedObject() interface{} {
	return p.Err
}

func (p *ErrorInfo) ErrorDetail() string {
	return p.Error()
}

func (p *ErrorInfo) AppendErrorDetail(b []byte) []byte {
	return errorDetail(b, p)
}

func (p *ErrorInfo) SummaryErr() error {
	return p.Err
}

func Info(err error, cmd ...interface{}) *ErrorInfo {
	return &Frame{Err: err, Args: cmd}
}

func InfoEx(calldepth int, err error, cmd ...interface{}) *ErrorInfo {
	return &Frame{Err: err, Args: cmd}
}

func Detail(err error) string {
	return err.Error()
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}
