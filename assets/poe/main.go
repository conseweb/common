package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

// 日志记录器
var logger = logging.MustGetLogger("poe")

// 程序入口
func main() {
	if e := shim.Start(new(PoeService)); e != nil {
		logger.Panicf("start PoeService return error: %v, exiting\n", e)
	}
}
