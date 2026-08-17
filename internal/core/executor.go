package core

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"time"

	"github.com/upinmcSE/godis/internal/constant"
)

func cmdPING(args []string) []byte {
	if len(args) > 1 {
		return Encode("ERR wrong number of argments for 'PING' command", false)
	}
	if len(args) == 0 {
		return Encode("PONG", false)
	} else {
		return Encode(args[0], false)
	}
}

func cmdSET(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode("ERR wrong number of arguments for 'SET' command", false)
	}

	var key, value string
	var tllMs int64 = -1
	key, value = args[0], args[1]
	if len(args) > 2 {
		tllSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		tllMs = tllSec * 1000 // convert to millisecond
	}

	dictStore.Set(key, dictStore.NewObject(key, value, tllMs))
	return constant.RespOk
}

func cmdGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}

	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}
	if dictStore.HasExpired(key) {
		return constant.RespNil
	}
	return Encode(obj.Value, false)
}

// check ttl of key
func cmdTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'TTL' command"), false)
	}

	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.TtlKeyNotExist
	}

	exp, isExpirySet := dictStore.GetExpiry(key)
	if !isExpirySet {
		return constant.TtlKeyExistNoExpire
	}

	remainMs := exp - time.Now().UnixMilli()
	if remainMs < 0 {
		return constant.TtlKeyNotExist
	}
	return Encode(remainMs/1000, false)
}

func cmdEXPIRE(args []string) []byte {
	return []byte{}
}

func cmdDEL(args []string) []byte {
	return []byte{}
}

func cmdEXISTS(args []string) []byte {
	return []byte{}
}

// ExecuteAndResponse given a Command, executes it and responses
func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	default:
		res = fmt.Appendf(nil, "-CMD NOT FOUND\r\n")
	}
	_, err := syscall.Write(connFd, res)
	return err
}
