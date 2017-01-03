package main

import (
	"fmt"
	"reflect"
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
		"addoto": addoto,
		"deloto": deloto,
		"addotm": addotm,
		"addmtm": addmtm,
		"query":  query,
	}
)

// error
var (
	err_unsupported_operation      = fmt.Errorf("unsupported operation")
	err_invalid_param_empty        = fmt.Errorf("parameter is empty")
	err_invalid_param_canNotEmpty  = fmt.Errorf("arguments can not be null characters")
	err_invalid_param_surplus      = fmt.Errorf("only one piece of data can be processed at a time")
	err_invalid_param_split        = fmt.Errorf("split symbol ',' can contain only one")
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
	data, err := handler(stub, args)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (self *NameService) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	handler, ok := handlerMap[function]
	if !ok {
		logger.Debug("unsupported operation %s", function)
		return nil, err_unsupported_operation
	}
	return handler(stub, args)
}

// 添加一条数据(一对一)
// 参数规范： 不能包含空字符;每次只处理一条数据;args[0] 为key, args[1] 为value
// 存储会存储两份，第二次存储k和v交换存储
// kv 都不存在或k存在v 不存在时存储，存储后，原v 值将被覆盖， kv交换 逻辑一样
func addoto(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if len(args) > 2 {
		logger.Warningf("invalid param : %v ,only one piece of data can be processed at a time", args)
		return nil, err_invalid_param_surplus
	}
	if len(strings.TrimSpace(args[0])) == 0 {
		logger.Warningf("invalid param : %v , arguments can not be null characters", args)
		return nil, err_invalid_param_canNotEmpty
	}
	if len((strings.TrimSpace(args[1]))) == 0 {
		logger.Warningf("invalid param : %v , arguments can not be null characters", args)
		return nil, err_invalid_param_canNotEmpty
	}
	if strings.Count(args[0], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	if strings.Count(args[1], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	val1, err := stub.GetState(args[0])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	if reflect.DeepEqual(val1, []byte(args[1])) { //kv 已经存在
		logger.Debugf("%s : %s is exists", args[0], args[1])
		return []byte(fmt.Sprintf("%s:%s", args[0], args[1])), nil
	} else { //kv 都不存在或k 存在，v不存在
		if err := stub.PutState(args[0], []byte(args[1])); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	}
	// 交换存储
	val2, err := stub.GetState(args[1])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	if !reflect.DeepEqual(val2, []byte(args[0])) { //kv 都不存在或k 存在，v 不存在
		if err := stub.PutState(args[1], []byte(args[0])); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	}
	return []byte(fmt.Sprintf("%s:%s", args[0], args[1])), nil
}

// 删除一条数据
// 先根据k 删除，再用查出得v 作为新k删除
func deloto(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if len(args) > 1 {
		logger.Warningf("invalid param : %v ,only one piece of data at a time", args)
		return nil, err_invalid_param_surplus
	}
	if strings.Count(args[0], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	if len(strings.TrimSpace(args[0])) == 0 {
		logger.Warningf("invalid param : %v , arguments can not be null characters", args)
		return nil, err_invalid_param_canNotEmpty
	}
	val1, err := stub.GetState(args[0])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	if len(val1) == 0 {
		return []byte("-1"), nil
	}
	err = stub.DelState(args[0])
	if err != nil {
		logger.Errorf("del state into blockchain return err: %v", err)
		return nil, err
	}
	val2, err := stub.GetState(string(val1))
	if err != nil {
		logger.Errorf("del state into blockchain return err: %v", err)
		return nil, err
	}
	if len(val2) > 0 {
		err = stub.DelState(string(val1))
		if err != nil {
			logger.Errorf("del state into blockchain return err: %v", err)
			return nil, err
		}
	}
	return []byte("200"), nil
}

// 添加一条数据(一对多)
// 参数规范： 不能包含空字符;每次只处理一条数据;args[0] 为key, args[1] 为value
// 存储会存储两份，第二次存储k和v交换存储
// kv 都不存在或k存在v 不存在时存储，存储后，原来的v = 原有v,现有v， kv交换逻辑： 当kv都不存在或者k 存在v 不存在时存储，存储后，原v 值将被覆盖
func addotm(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if len(args) > 2 {
		logger.Warningf("invalid param : %v ,only one piece of data can be processed at a time", args)
		return nil, err_invalid_param_surplus
	}
	if len(strings.TrimSpace(args[0])) == 0 {
		logger.Warningf("invalid param : %v , arguments can not be null characters", args)
		return nil, err_invalid_param_canNotEmpty
	}
	if len((strings.TrimSpace(args[1]))) == 0 {
		logger.Warningf("invalid param : %v , arguments can not be null characters", args)
		return nil, err_invalid_param_canNotEmpty
	}
	if strings.Count(args[0], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	if strings.Count(args[1], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	val1, err := stub.GetState(args[0])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	if reflect.DeepEqual(val1, []byte(args[1])) { //kv 已经存在
		logger.Debugf("%s : %s is exists", args[0], args[1])
		return []byte(fmt.Sprintf("%s:%s", args[0], args[1])), nil
	} else if len(val1) > 0 { //k 存在，v 不存在
		valStr1 := string(val1) + "," + args[1]
		if err := stub.PutState(args[0], []byte(valStr1)); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	} else { //kv 都不存在
		if err := stub.PutState(args[0], []byte(args[1])); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	}
	// 交换存储
	val2, err := stub.GetState(args[1])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	if !reflect.DeepEqual(val2, []byte(args[0])) { //kv 都不存在或k 存在，v 不存在
		if err := stub.PutState(args[1], []byte(args[0])); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	}
	return []byte(fmt.Sprintf("%s:%s", args[0], args[1])), nil
}

// 添加一条数据(多对多)
// 参数规范： 不能包含空字符;每次只处理一条数据;args[0] 为key, args[1] 为value
// 存储会存储两份，第二次存储k和v交换存储
// kv 都不存在或k存在v 不存在时存储，存储后，原来的v = 原有v,现有v， kv交换逻辑一样
func addmtm(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if len(args) > 2 {
		logger.Warningf("invalid param : %v ,only one piece of data can be processed at a time", args)
		return nil, err_invalid_param_surplus
	}
	if len(strings.TrimSpace(args[0])) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if len((strings.TrimSpace(args[1]))) == 0 {
		logger.Warningf("invalid param : %v ,parameter is empty", args)
		return nil, err_invalid_param_empty
	}
	if strings.Count(args[0], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	if strings.Count(args[1], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
	}
	val1, err := stub.GetState(args[0])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	if reflect.DeepEqual(val1, []byte(args[1])) { //kv 已经存在
		logger.Debugf("%s : %s is exists", args[0], args[1])
		return []byte(fmt.Sprintf("%s:%s", args[0], args[1])), nil
	} else if len(val1) > 0 { //k 存在，v 不存在
		valStr1 := string(val1) + "," + args[1]
		if err := stub.PutState(args[0], []byte(valStr1)); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	} else { //kv 都不存在
		if err := stub.PutState(args[0], []byte(args[1])); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	}
	// 交换存储
	val2, err := stub.GetState(args[1])
	if err != nil {
		logger.Errorf("get state into blockchain return err: %v", err)
		return nil, err
	}
	if !reflect.DeepEqual(val2, []byte(args[0])) && len(val2) > 0 {
		valStr2 := string(val2) + "," + args[0]
		if err := stub.PutState(args[1], []byte(valStr2)); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	} else {
		if err := stub.PutState(args[1], []byte(args[0])); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}
	}
	return []byte(fmt.Sprintf("%s:%s", args[0], args[1])), nil
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
	if strings.Count(args[0], " ") > 0 {
		logger.Warningf("invalid param : %v , can not contain null characters", args)
		return nil, err_invalid_param_containEmpty
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
