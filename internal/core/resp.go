package core

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/upinmcSE/godis/internal/constant"
)

const CRLF string = "\r\n"

// +OK\r\n => OK, 5
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

// :123\r\n => 123
func readInt64(data []byte) (int64, int, error) {
	var res int64 = 0
	pos := 1
	var sign int64 = 1
	if data[pos] == '-' {
		sign = -1
		pos++
	}
	if data[pos] == '+' {
		pos++
	}
	for data[pos] != '\r' {
		res = res*10 + int64(data[pos]-'0')
		pos++
	}
	return sign * res, pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// $5\r\nhello\r\n => 5, 4
func readLen(data []byte) (int, int) {
	res, pos, _ := readInt64(data)
	return int(res), pos
}

// $5\r\nhello\r\n => "hello"
func readBulkString(data []byte) (string, int, error) {
	length, pos := readLen(data)
	return string(data[pos:(pos + length)]), pos + length + 2, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n => {"hello", "world"}
func readArray(data []byte) (any, int, error) {
	length, pos := readLen(data)
	var res []any = make([]any, length)

	for i := range res {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		res[i] = elem
		pos += delta
	}
	return res, pos, nil
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("No data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case ':':
		return readInt64(data)
	case '-':
		return readError(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

// RESP format data => raw data
func Decode(data []byte) (any, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

func encodeString(value string, isSimpleString bool) []byte {
	if isSimpleString {
		return fmt.Appendf(nil, "+%s%s", value, CRLF)
	}
	return fmt.Appendf(nil, "$%d%s%s%s", len(value), CRLF, value, CRLF)
}

func encodeStringArray(sa []string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeString(s, false))
	}
	return fmt.Appendf(nil, "*%d\r\n%s", len(sa), buf.Bytes())
}

// Raw data => RESP format data
func Encode(value any, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		return encodeString(v, isSimpleString)
	case int64, int32, int16, int8, int:
		return fmt.Appendf(nil, ":%d\r\n", v)
	case error:
		return fmt.Appendf(nil, "-%s\r\n", v)
	case []string:
		return encodeStringArray(value.([]string))
	case [][]string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, sa := range value.([][]string) {
			buf.Write(encodeStringArray(sa))
		}
		return fmt.Appendf(nil, "*%d\r\n%s", len(value.([][]string)), buf.Bytes())
	case []any:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range value.([]any) {
			buf.Write(Encode(x, false))
		}
		return fmt.Appendf(nil, "*%d\r\n%s", len(value.([]any)), buf.Bytes())
	default:
		return constant.RespNil
	}
}
