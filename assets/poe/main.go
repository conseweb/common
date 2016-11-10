package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Proof of Existence Service（存在性证明服务）
type PoeService struct{}

func (this *PoeService) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, errors.New("invalid function name, 'deploy' only")
	}
	logger.Info("deploy Lepuscoin successfully")
	return nil, nil
}

func (this *PoeService) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case INVOKE_FUNC_REGISTER:
		return this.register(stub, args)
	default:
		return nil, errors.New("unsupported operation")
	}
}

func (this *PoeService) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case QUERY_FUNC_EXISTENCE:
		return this.existence(stub, args)
	default:
		return nil, errors.New("unsupported operation")
	}
}

//注册键值
func (this *PoeService) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 验证非空
	if len(args) == 0 {
		return nil, errors.New("func <register> Parameter is not valid,Cannot be empty or contain null characters")
	}
	var sysName string = args[0]
	if len(strings.TrimSpace(sysName)) == 0 {
		sysName = "base"
	}
	for i := 1; i < len(args); i++ {
		if len(strings.TrimSpace(args[i])) > 0 {
			hkey, e := hashKey(sysName, args[i])
			if e != nil {
				return nil, e
			}
			e = stub.PutState(hkey, []byte{1})
			if e != nil {
				return nil, e
			}
		}
	}
	return nil, nil
}

//检索键值是否存在
func (this *PoeService) existence(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var list []QueryResult
	// 验证非空
	if len(args) == 0 {
		return nil, errors.New("func <existence> Parameter is not valid,Cannot be empty or contain null characters")
	}
	var sysName string = args[0]
	if len(strings.TrimSpace(sysName)) == 0 {
		sysName = "base"
	}
	for i := 1; i < len(args); i++ {
		if len(strings.TrimSpace(args[i])) > 0 {
			hkey, e := hashKey(sysName, args[i])
			if e != nil {
				return nil, errors.New("func <existence> error:" + e.Error())
			}
			d, e := stub.GetState(hkey)
			if e != nil {
				return nil, errors.New("func <existence> error:" + e.Error())
			}
			if d != nil && len(d) > 0 {
				m := QueryResult{}
				m.Key = args[i]
				m.Exist = true
				list = append(list, m)
			}
		}
	}
	data, e := json.Marshal(&list)
	if e != nil {
		return nil, errors.New("func <existence> error:" + e.Error())
	}
	return data, nil
}

// 计算键值哈希
// sysName: 业务系统名称
// key: 需要加密的字符串
func hashKey(sysName, key string) (r string, e error) {
	var (
		configSys *ConfigSystem
		strategy  cryptoStrategy
		data      []byte
	)
	if configSys, e = configSystem(sysName); e != nil {
		return "", e
	}
	strategy = cryptoStrategyMap[configSys.CryptoStrategy](sysName)
	if strategy == nil {
		strategy = cryptoStrategyMap["default"](sysName)
	}
	if data, e = strategy.algorithm([]byte(key)); e != nil {
		return "", e
	}
	return string(data), nil
}

// 程序入口
func main() {
	if e := shim.Start(new(PoeService)); e != nil {
		logger.Errorf("start PoeService return error: %v, exiting\n", e)
	}
}
