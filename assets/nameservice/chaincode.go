package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type (
	// 命名服务
	NameService struct{}
	Handler     func(stub shim.ChaincodeStubInterface, args []string) ([]byte, error)
)

var (
	handlerMap map[string]Handler = map[string]Handler{
		"add":   add,
		"query": query,
	}
)

// error
var (
	err_unsupported_operation      = fmt.Errorf("unsupported operation")
	err_invalid_param_empty        = fmt.Errorf("parameter is empty")
	err_invalid_param_canNotEmpty  = fmt.Errorf("arguments can not be null characters")
	err_invalid_param_surplus      = fmt.Errorf("only one piece of data can be processed at a time")
	err_invalid_param_split        = fmt.Errorf("split symbol ':' can contain only one")
	err_invalid_param_containEmpty = fmt.Errorf("can not contain null characters")
)

func (self *NameService) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, err_unsupported_operation
	}
	logger.Debug("deploy poe successfully")
	return []byte("200"), nil
}

func (self *NameService) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	handler, ok := handlerMap[function]
	if !ok {
		logger.Debug("unsupported operation %s", function)
		return nil, err_unsupported_operation
	}
	return handler(stub, args)
}

func (self *NameService) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	handler, ok := handlerMap[function]
	if !ok {
		logger.Debug("unsupported operation %s", function)
		return nil, err_unsupported_operation
	}
	return handler(stub, args)
}

// 添加一条数据
// 参数规范： 不能包含空字符;每次只处理一条数据;kv分割符号':' 只能出现一次
// 参数格式： k:v
// 存储会存储两份，第二次存储k和v交换存储
// 存储完成触发‘Add_Complete’ 事件
func add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if len(args) > 1 {
		logger.Warningf("invalid param : %v ,only one piece of data can be processed at a time", args)
		return nil, err_invalid_param_surplus
	}
	if strings.Count(args[0], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	if strings.Count(args[0], ":") != 1 {
		logger.Warningf("invalid param: %v ,split symbol ':' can contain only one", args)
		return nil, err_invalid_param_split
	}
	array := strings.Split(args[0], ":")
	if err := stub.PutState(array[0], []byte(array[1])); err != nil {
		logger.Errorf("put state into blockchain return err: %v", err)
		return nil, err
	}
	if err := stub.PutState(array[1], []byte(array[0])); err != nil {
		logger.Errorf("put state into blockchain return err: %v", err)
		return nil, err
	}
	if err := stub.SetEvent("Add_Complete", []byte(args[0])); err != nil {
		logger.Errorf("set event 'Add_Complete' return err: %v", err)
		return nil, err
	}
	return []byte(args[0]), nil
}

// 查询一条数据
// 参数规范： 不能包含空字符;每次只处理一条数据
func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if len(args) > 1 {
		logger.Warningf("invalid param : %v ,only one piece of data at a time", args)
		return nil, err_invalid_param_surplus
	}
	if len(strings.TrimSpace(args[0])) == 0 {
		logger.Warningf("invalid param : %v , arguments can not be null characters", args)
		return nil, err_invalid_param_canNotEmpty
	}
	data, err := stub.GetState(args[0])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	return data, nil
}
