package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var poeFuncMap map[string]poeFunc = map[string]poeFunc{
	"register":  register,
	"existence": existence,
}

type (
	// Proof of Existence Service（存在性证明服务）
	PoeService struct{}
	poeFunc    func(shim.ChaincodeStubInterface, []string) ([]byte, error)
)

func (this *PoeService) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, errors.New("invalid function name, 'deploy' only")
	}
	logger.Info("deploy Lepuscoin successfully")
	return nil, nil
}

func (this *PoeService) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if fn, ok := poeFuncMap[function]; ok {
		return fn(stub, args)
	} else {
		return nil, errors.New("unsupported operation")
	}
}

func (this *PoeService) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if fn, ok := poeFuncMap[function]; ok {
		return fn(stub, args)
	} else {
		return nil, errors.New("unsupported operation")
	}
}

// 注册键值
// args : 参数第一个元素为业务系统标记，默认为base
func register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var (
		cfgSys  *ConfigSystem
		sysName string = args[0]
		hkey    string
		e       error
	)
	// 验证非空
	if len(args) == 0 {
		return nil, errors.New("func <register> Parameter is not valid,Cannot be empty or contain null characters")
	}
	if cfgSys, e = configSystem(sysName); e != nil {
		return nil, errors.New("func <register> error:" + e.Error())
	}
	for i := 1; i < len(args); i++ {
		if len(strings.TrimSpace(args[i])) == 0 {
			continue
		}
		if hkey, e = hashKey(cfgSys, args[i]); e != nil {
			return nil, errors.New("func <register> error:" + e.Error())
		}
		if e = stub.PutState(hkey, []byte{1}); e != nil {
			return nil, errors.New("func <register> error:" + e.Error())
		}
	}
	return nil, nil
}

// 检索键值是否存在
// args : 参数第一个元素为业务系统标记，默认为base
func existence(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var (
		list    []QueryResult
		cfgSys  *ConfigSystem
		sysName string = args[0]
		hkey    string
		data    []byte
		e       error
	)
	// 验证非空
	if len(args) == 0 {
		return nil, errors.New("func <existence> Parameter is not valid,Cannot be empty or contain null characters")
	}
	if cfgSys, e = configSystem(sysName); e != nil {
		return nil, e
	}
	for i := 1; i < len(args); i++ {
		if len(strings.TrimSpace(args[i])) == 0 {
			continue
		}
		if hkey, e = hashKey(cfgSys, args[i]); e != nil {
			return nil, errors.New("func <existence> error:" + e.Error())
		}
		if data, e = stub.GetState(hkey); e != nil {
			return nil, errors.New("func <existence> error:" + e.Error())
		}
		if len(data) > 0 {
			m := QueryResult{}
			m.Key = args[i]
			m.Exist = true
			list = append(list, m)
		}
	}
	if data, e = json.Marshal(&list); e != nil {
		return nil, errors.New("func <existence> error:" + e.Error())
	}
	return data, nil
}

// 计算键值哈希
// cfgSys: 业务系统配置
// key: 需要加密的字符串
func hashKey(cfgSys *ConfigSystem, key string) (string, error) {
	var (
		strategy cryptoStrategyFunc
		data     []byte
		e        error
	)
	if strategy = cryptoStrategyMap[cfgSys.CryptoStrategy]; strategy == nil {
		strategy = cryptoStrategyMap["default"]
	}
	if data, e = strategy(cfgSys).algorithm([]byte(key)); e != nil {
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
