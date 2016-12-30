package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/op/go-logging"
)

// 日志记录器
var logger = logging.MustGetLogger("namesrvc")

// 程序入口
func main() {
	if e := shim.Start(new(NameService)); e != nil {
		logger.Panicf("start NameService return error: %v, exiting\n", e)
	}
}
