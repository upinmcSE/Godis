package core

import (
	"errors"
	"fmt"
	"github.com/upinmcSE/godis/internal/constant"
	"strconv"
	"syscall"
	"time"
)

func cmdPING(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode("ERR wrong number of arguments for 'PING' command", false)
	}
	if len(args) == 0 {
		return Encode("PONG", false)
	} else {
		return Encode(args[0], false)
	}
	return res
}

func cmdSET(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode("ERR wrong number of arguments for 'SET' command", false)
	}

	var key, value string
	var ttlMs int64 = -1

	key, value = args[0], args[1]
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000 // convert to millisecond
	}

	dictStore.Set(key, dictStore.NewObject(key, value, ttlMs))
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

	remainMs := exp - uint64(time.Now().UnixMilli())
	if remainMs < 0 {
		return constant.TtlKeyNotExist
	}
	return Encode(int64(remainMs/1000), false)
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
	case "EXPIRE":
		res = cmdEXPIRE(cmd.Args)
	case "DEL":
		res = cmdDEL(cmd.Args)
	case "EXISTS":
		res = cmdEXISTS(cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-CMD NOT FOUND\r\n"))
	}
	_, err := syscall.Write(connFd, res)
	return err
}
